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

	"code.superseriousbusiness.org/gopkg/log"
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
		totalFollows, err := db.NewSelect().Table("follows").Count(ctx)
		if err != nil {
			return gtserror.Newf("error getting follow table count: %w", err)
		}

		log.Warnf(ctx, "migrating %d follow flags to new column", totalFollows)
		if err := migrateFollowlike(ctx, totalFollows, db, "follows"); err != nil {
			return err
		}

		// Migrate follows requests.
		totalFollowReqs, err := db.NewSelect().Table("follow_requests").Count(ctx)
		if err != nil {
			return gtserror.Newf("error getting follow table count: %w", err)
		}

		log.Warnf(ctx, "migrating %d follow request flags to new column", totalFollowReqs)
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

	var ids []string
	for i := 1; ; i++ {

		// Reset IDs.
		clear(ids)
		ids = ids[:0]

		// Select IDs page.
		if err := tx.
			NewSelect().
			Table(table).
			Column("id").
			Where("? < ?", bun.Ident("id"), maxID).
			OrderExpr("? DESC", bun.Ident("id")).
			Limit(500).
			Scan(ctx, &ids); err != nil && !errors.Is(err, sql.ErrNoRows) {
			return gtserror.Newf("error selecting items: %w", err)
		}

		if len(ids) == 0 {
			// No more items!
			//
			// Transaction will be closed
			// after leaving the loop.
			break
		}

		// Set next maxID value.
		maxID = ids[len(ids)-1]

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
	}

	// Close the transaction.
	if err := tx.Commit(); err != nil {
		return err
	}

	// Merge WAL file to try minimize its size.
	return doWALCheckpoint(ctx, db)
}
