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
	"time"

	"code.superseriousbusiness.org/gopkg/log"
	"code.superseriousbusiness.org/gopkg/xslices"
	dbpkg "code.superseriousbusiness.org/gotosocial/internal/db"
	newmodel "code.superseriousbusiness.org/gotosocial/internal/db/bundb/migrations/20260826130603_follow_migration/newmodel"
	oldmodel "code.superseriousbusiness.org/gotosocial/internal/db/bundb/migrations/20260826130603_follow_migration/oldmodel"
	"code.superseriousbusiness.org/gotosocial/internal/gtserror"
	"code.superseriousbusiness.org/gotosocial/internal/id"
	"github.com/uptrace/bun"
)

func init() {
	up := func(ctx context.Context, db *bun.DB) error {
		// Add new flags column to follows
		// and follow_requests tables.
		for _, model := range []any{
			(*newmodel.Follow)(nil),
			(*newmodel.FollowRequest)(nil),
		} {
			if err := addColumn(ctx, db, model, "Flags"); err != nil {
				return err
			}
		}

		// Migrate follows.
		log.Warnf(ctx, "migrating follow flags to new column")
		totalFollows, err := db.NewSelect().Table("follows").Count(ctx)
		if err != nil {
			return gtserror.Newf("error getting follow table count: %w", err)
		}
		if err := migrateFollowlike(ctx, totalFollows, db, "follows"); err != nil {
			return err
		}

		// Migrate follows requests.
		log.Warnf(ctx, "migrating follow_request flags to new column")
		totalFollowReqs, err := db.NewSelect().Table("follow_requests").Count(ctx)
		if err != nil {
			return gtserror.Newf("error getting follow table count: %w", err)
		}
		if err := migrateFollowlike(ctx, totalFollowReqs, db, "follow_requests"); err != nil {
			return err
		}

		// Drop unused columns from database.
		for _, field := range []string{
			"ShowReblogs",
			"Notify",
		} {
			for _, model := range []any{
				(*oldmodel.Follow)(nil),
				(*oldmodel.FollowRequest)(nil),
			} {
				if err := dropColumn(ctx, db,
					model,
					field,
				); err != nil {
					return err
				}

				// WAL merge after each drop to minimize WAL size.
				if err := doWALCheckpoint(ctx, db); err != nil {
					return err
				}
			}
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

func migrateFollowlike(
	ctx context.Context,
	total int,
	db *bun.DB,
	table string,
) error {
	// Open initial transaction.
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	// Start at largest
	// possible ULID value.
	maxID := id.Highest

	// Total updated count.
	var updatedTotal int64

	var items []oldmodel.WithID
	var ids []string
	for i := 1; ; i++ {

		// Reset slices.
		clear(items)
		clear(ids)
		items = items[:0]
		ids = ids[:0]

		// Mark batch start.
		start := time.Now()

		// Select from items.
		if err := tx.NewSelect().
			Table(table).
			Column("id").
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

		} else if i%10 == 0 {
			// Begin a new transaction every
			// 10 batches (~5,000 items),
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

		// Set next maxID value from items.
		maxID = items[len(items)-1].GetID()

		// Gather all IDs from selected items.
		ids = xslices.Gather(ids,
			items, func(s oldmodel.WithID) string {
				return s.GetID()
			})

		// IDs as a bun.List() value.
		inIDs := bun.List(ids)

		// Perform an UPDATE query for each new possible
		// follow flag bit field value, performing a bitwise
		// OR on "flags" to set the bit for matching WHERE clause.
		for _, q := range []struct {
			Bit   newmodel.FollowFlag
			Where dbpkg.BunExpr
		}{
			{Bit: newmodel.FollowFlagShowReblogs, Where: dbpkg.BunExpr{"? = true", dbpkg.Idents("show_reblogs")}},
			{Bit: newmodel.FollowFlagNotify, Where: dbpkg.BunExpr{"? = true", dbpkg.Idents("notify")}},
		} {
			if _, err := tx.NewUpdate().
				Table(table).

				// Only operating on item IDs in selected batch.
				Where("? IN (?)", bun.Ident("id"), inIDs).

				// Updating "flags" via OR to set the current 'bit' flag value.
				Set("? = (?|?)", bun.Ident("flags"), bun.Ident("flags"), q.Bit).

				// Only on given WHERE clause.
				Where(q.Where.Fmt, q.Where.Arg...).
				Exec(ctx); err != nil {
				return gtserror.Newf("error setting \"flags\" value = %s: %w", q.Bit.String(), err)
			}
		}

		// Increment updated total by ID count.
		updatedTotal += int64(len(items))

		// Calculate rows / second tx speed.
		timeTaken := time.Since(start).Seconds()
		secsPerRow := float64(timeTaken) / float64(len(items))
		rowsPerSec := float64(1) / float64(secsPerRow)

		// Calculate percentage of all items updated so far.
		perc := (float64(updatedTotal) / float64(total)) * 100

		log.Infof(ctx, "[~%.2f%% done; ~%.0f rows/s] migrating follow flags",
			perc, rowsPerSec)
	}

	// Close the final db transaction.
	if err := tx.Commit(); err != nil {
		return err
	}

	// Merge WAL file to try minimize its size.
	return doWALCheckpoint(ctx, db)
}
