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
		if err := migrateFollowlike(ctx, db, "follows"); err != nil {
			return err
		}

		// Migrate follows requests.
		if err := migrateFollowlike(ctx, db, "follow_requests"); err != nil {
			return err
		}

		return db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
			// Drop incorrect indexes for follow (requests) tables;
			// they were using updated_at instead of created_at.
			if err := dropIndex(ctx, tx, "follows_account_id_idx"); err != nil {
				return err
			}

			if err := dropIndex(ctx, tx, "follows_target_account_id_idx"); err != nil {
				return err
			}

			if err := dropIndex(ctx, tx, "follow_requests_account_id_idx"); err != nil {
				return err
			}

			if err := dropIndex(ctx, tx, "follow_requests_target_account_id_idx"); err != nil {
				return err
			}

			// Now re-create them.
			if err := createIndex(ctx, tx,
				"follows_account_id_idx",
				"follows",
				dbpkg.BunExpr{
					"?, ? DESC",
					dbpkg.Idents(
						"account_id",
						"created_at",
					)},
			); err != nil {
				return err
			}

			if err := createIndex(ctx, tx,
				"follows_target_account_id_idx",
				"follows",
				dbpkg.BunExpr{
					"?, ? DESC",
					dbpkg.Idents(
						"target_account_id",
						"created_at",
					)},
			); err != nil {
				return err
			}

			if err := createIndex(ctx, tx,
				"follow_requests_account_id_idx",
				"follow_requests",
				dbpkg.BunExpr{
					"?, ? DESC",
					dbpkg.Idents(
						"account_id",
						"created_at",
					)},
			); err != nil {
				return err
			}

			if err := createIndex(ctx, tx,
				"follow_requests_target_account_id_idx",
				"follow_requests",
				dbpkg.BunExpr{
					"?, ? DESC",
					dbpkg.Idents(
						"target_account_id",
						"created_at",
					)},
			); err != nil {
				return err
			}

			// Drop unused columns from database.
			for _, field := range []string{
				"ShowReblogs",
				"Notify",
				"UpdatedAt",
			} {
				for _, model := range []any{
					(*oldmodel.Follow)(nil),
					(*oldmodel.FollowRequest)(nil),
				} {
					if err := dropColumn(ctx, tx,
						model,
						field,
					); err != nil {
						return err
					}
				}
			}

			return nil
		})
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
	db *bun.DB,
	table string,
) error {
	log.Infof(ctx, "migrating %s flags to new column", table)

	// Open initial transaction.
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	// Start at largest
	// possible ULID value.
	maxID := id.Highest

	var ids []string
	for {

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
