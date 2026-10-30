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
				TableExpr("? AS ?", bun.Ident("statuses"), bun.Ident("status")).
				Set("? = NULL", bun.Ident("status.url")).
				Where("? IS NOT NULL", bun.Ident("status.url")).
				Where("? = ?", bun.Ident("status.url"), bun.Ident("status.uri")).
				Exec(ctx); err != nil {
				return gtserror.Newf("error zeroing status urls: %w", err)
			}

			// Zero all account URLs that are the same as the URI.
			log.Info(ctx, "removing duplicate url from accounts, this may be a large query...")
			if _, err := tx.
				NewUpdate().
				TableExpr("? AS ?", bun.Ident("accounts"), bun.Ident("account")).
				Set("? = NULL", bun.Ident("account.url")).
				Where("? IS NOT NULL", bun.Ident("account.url")).
				Where("? = ?", bun.Ident("account.url"), bun.Ident("account.uri")).
				Exec(ctx); err != nil {
				return gtserror.Newf("error zeroing account urls: %w", err)
			}

			return nil
		}); err != nil {
			return err
		}

		if err := doWALCheckpoint(ctx, db); err != nil {
			return err
		}

		// Iterate through statuses and ensure we
		// only store url when url.host = uri.host.
		type UrlUri struct {
			ID  string
			URI string
			URL string
		}

		// Start at largest
		// possible ULID value.
		maxID := id.Highest

		log.Warnf(ctx, "migrating status flags to new column, this may take a *long* time")

		// Open initial transaction.
		tx, err := db.BeginTx(ctx, nil)
		if err != nil {
			return err
		}

		// Total updated count.
		var updatedTotal int64

		var items []UrlUri

		for i := 1; ; i++ {

			// Reset items.
			clear(items)
			items = items[:0]

			// Mark batch start.
			start := time.Now()

			// Select from items.
			if err := tx.NewSelect().
				Table("statuses").
				Column("id", "uri", "url").
				Where("? < ?", bun.Ident("id"), maxID).
				OrderExpr("? DESC", bun.Ident("id")).
				Limit(500).
				Scan(ctx); err != nil && !errors.Is(err, sql.ErrNoRows) {
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

			// Increment updated total by ID count.
			updatedTotal += int64(len(items))

			// Calculate rows / second tx speed.
			timeTaken := time.Since(start).Seconds()
			secsPerRow := float64(timeTaken) / float64(len(items))
			rowsPerSec := float64(1) / float64(secsPerRow)

			// Calculate percentage of all statuses updated so far.
			perc := (float64(updatedTotal) / float64(total)) * 100

			log.Infof(ctx, "[~%.2f%% done; ~%.0f rows/s] migrating status urls",
				perc, rowsPerSec)
		}

		// Close the final db transaction.
		if err := tx.Commit(); err != nil {
			return err
		}

		// Merge WAL file to try minimize its size.
		if err := doWALCheckpoint(ctx, db); err != nil {
			return err
		}

		return nil
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
	type UrlUri struct {
		ID  string
		URI string
		URL string
	}

	// Start at largest
	// possible ULID value.
	maxID := id.Highest

	log.Warnf(ctx, "migrating status flags to new column, this may take a *long* time")

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
	var updatedTotal int64

	var items []UrlUri

	for i := 1; ; i++ {

		// Reset items.
		clear(items)
		items = items[:0]

		// Mark batch start.
		start := time.Now()

		// Select from items.
		if err := tx.NewSelect().
			Table("statuses").
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

		// For each item, check if URL is
		// parseable and on the same host
		// as URI. If not, remove it.
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

		}

		// Increment updated total by ID count.
		updatedTotal += int64(len(items))

		// Calculate rows / second tx speed.
		timeTaken := time.Since(start).Seconds()
		secsPerRow := float64(timeTaken) / float64(len(items))
		rowsPerSec := float64(1) / float64(secsPerRow)

		// Calculate percentage of all statuses updated so far.
		perc := (float64(updatedTotal) / float64(total)) * 100

		log.Infof(ctx, "[~%.2f%% done; ~%.0f rows/s] migrating status urls",
			perc, rowsPerSec)
	}

	// Close the final db transaction.
	if err := tx.Commit(); err != nil {
		return err
	}

	// Merge WAL file to try minimize its size.
	if err := doWALCheckpoint(ctx, db); err != nil {
		return err
	}

	return nil
}
