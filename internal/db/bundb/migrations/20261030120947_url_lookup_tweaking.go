// GoToSocial
// Copyright (C) GoToSocial Authors admin@gotosocial.org
// SPDX-License-Identifier: AGPL-3.0-or-later
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package migrations

import (
	"context"
	"database/sql"
	"errors"
	"net/url"
	"time"

	"code.superseriousbusiness.org/gopkg/log"
	"code.superseriousbusiness.org/gotosocial/internal/gtserror"
	"code.superseriousbusiness.org/gotosocial/internal/id"
	"code.superseriousbusiness.org/gotosocial/internal/util"
	"github.com/uptrace/bun"
)

func init() {
	up := func(ctx context.Context, db *bun.DB) error {
		if err := db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {

			// Zero all status URLs that are the same as the URI.
			log.Info(ctx, "removing duplicate url from statuses, this may be a large query...")
			if _, err := tx.
				NewUpdate().
				Table("statuses").
				Set("? = NULL", bun.Ident("url")).
				Where("? IS NOT NULL", bun.Ident("url")).
				Where("? = ?", bun.Ident("url"), bun.Ident("uri")).
				Exec(ctx); err != nil {
				return gtserror.Newf("error zeroing status urls: %w", err)
			}

			// Zero all account URLs that are the same as the URI.
			log.Info(ctx, "removing duplicate url from accounts, this may be a large query...")
			if _, err := tx.
				NewUpdate().
				Table("accounts").
				Set("? = NULL", bun.Ident("url")).
				Where("? IS NOT NULL", bun.Ident("url")).
				Where("? = ?", bun.Ident("url"), bun.Ident("uri")).
				Exec(ctx); err != nil {
				return gtserror.Newf("error zeroing account urls: %w", err)
			}

			return nil
		}); err != nil {
			return err
		}

		if err := migrateUrls(ctx, db, "statuses"); err != nil {
			return err
		}

		return migrateUrls(ctx, db, "accounts")
	}

	down := func(ctx context.Context, db *bun.DB) error {
		return nil
	}

	if err := Migrations.Register(up, down); err != nil {
		panic(err)
	}
}

func migrateUrls(
	ctx context.Context,
	db *bun.DB,
	table string,
) error {

	// Start at largest
	// possible ULID value.
	maxID := id.Highest

	log.Warnf(ctx, "migrating `url` column on %s, this may take a long time", table)

	// Open initial transaction.
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	// Count all entries
	// we'll need to look at.
	total, err := tx.
		NewSelect().
		Table(table).
		Where("? IS NOT NULL", bun.Ident("url")).
		Count(ctx)
	if err != nil {
		return err
	}

	// Total updated count.
	var migratedTotal int64

	var items []struct {
		ID  string
		URI string
		URL string
	}
	for i := 1; ; i++ {

		// Reset items.
		clear(items)
		items = items[:0]

		// Mark batch start.
		start := time.Now()

		// Select from items.
		if err := tx.NewSelect().
			Table(table).
			Column("id", "uri", "url").
			Where("? IS NOT NULL", bun.Ident("url")).
			Where("? < ?", bun.Ident("id"), maxID).
			OrderExpr("? DESC", bun.Ident("id")).
			Limit(500).
			Scan(ctx, &items); err != nil && !errors.Is(err, sql.ErrNoRows) {
			return gtserror.Newf("error selecting items: %w", err)
		}

		if len(items) == 0 {
			// No more items!
			//
			// Transaction will be closed
			// after leaving the loop.
			break

		} else if i%200 == 0 {
			// Begin a new transaction every
			// 200 batches (~100,000 items),
			// to avoid massive commits.

			// Close existing db transaction.
			if err := tx.Commit(); err != nil {
				return err
			}

			// Merge WAL file to try minimize its size.
			if err := doWALCheckpoint(ctx, db); err != nil {
				return err
			}

			// Start a new db transaction.
			tx, err = db.BeginTx(ctx, nil)
			if err != nil {
				return err
			}
		}

		// Set next max ID for paging.
		maxID = items[len(items)-1].ID

		// For each item, check if URL is parseable and on
		// same host as URI. If not, mark URL for removal.
		var ids []string
		for _, item := range items {
			uri, err := url.Parse(item.URI)
			if err != nil {
				// This should never happen since we should
				// parse the URI of every status and account
				// before we store them, but you never know...
				log.Errorf(ctx, "error parsing status uri: %v", err)
				continue
			}

			url, err := url.Parse(item.URL)
			var ok bool
			if err == nil {
				ok, _ = util.SameHost(url, uri)
			}
			if ok {
				// This item is
				// fine, carry on.
				continue
			}

			// Item is not okay either
			// because URL didn't parse,
			// or host was not same as URI.
			ids = append(ids, item.ID)
		}

		// Set url = NULL
		// on gathered IDs.
		if len(ids) != 0 {
			_, err := tx.NewUpdate().
				Table(table).
				Set("? = NULL", bun.Ident("url")).
				Where("? IN (?)", bun.Ident("id"), bun.List(ids)).
				Exec(ctx)
			if err != nil {
				return gtserror.Newf("error updating: %w", err)
			}
		}

		// Increment migrated total.
		migratedTotal += int64(len(items))

		// Calculate rows / second tx speed.
		timeTaken := time.Since(start).Seconds()
		secsPerRow := float64(timeTaken) / float64(len(items))
		rowsPerSec := float64(1) / float64(secsPerRow)

		// Calculate percentage of all statuses migrated so far.
		perc := (float64(migratedTotal) / float64(total)) * 100

		log.Infof(ctx, "[~%.2f%% done; ~%.0f rows/s] migrating %s.url",
			perc, rowsPerSec, table)
	}

	// Close the final db transaction.
	if err := tx.Commit(); err != nil {
		return err
	}

	// Merge WAL file to try minimize its size.
	return doWALCheckpoint(ctx, db)
}
