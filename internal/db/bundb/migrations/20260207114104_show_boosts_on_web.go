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

	"code.superseriousbusiness.org/gopkg/log"
	gtsmodel "code.superseriousbusiness.org/gotosocial/internal/db/bundb/migrations/20260207114104_show_boosts_on_web/newmodel"
	"github.com/uptrace/bun"
)

func init() {
	up := func(ctx context.Context, db *bun.DB) error {
		return db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {

			// Add column to AccountSettings table. Its default of false is safe.
			if err := addColumn(ctx, tx, (*gtsmodel.AccountSettings)(nil), "WebIncludeBoosts"); err != nil {
				return err
			}

			log.Info(ctx, "recreating statuses web view indexes, this may take a little while...")

			// Remove existing web index.
			if _, err := tx.
				NewDropIndex().
				Index("statuses_profile_web_view_idx").
				IfExists().
				Exec(ctx); err != nil {
				return err
			}

			// Create index for standard web view
			// (local, no boosts, no replies).
			if _, err := tx.NewCreateIndex().
				Table("statuses").
				Index("statuses_profile_web_view_idx").
				Column(
					"local",
					"account_id",
					"visibility",
					"in_reply_to_uri",
					"boost_of_id",
					"federated",
				).
				ColumnExpr("? DESC", bun.Ident("id")).
				Where("? = ?", bun.Ident("local"), true).
				Where("? IS NULL", bun.Ident("in_reply_to_uri")).
				Where("? IS NULL", bun.Ident("boost_of_id")).
				Where("? = ?", bun.Ident("federated"), true).
				IfNotExists().
				Exec(ctx); err != nil {
				return err
			}

			// Create index for including boosts in web view.
			// This is the same as the above index but ignores
			// "boost_of_id" column as this may or may not be set.
			if _, err := tx.NewCreateIndex().
				Table("statuses").
				Index("statuses_profile_web_view_including_boosts_idx").
				Column(
					"local",
					"account_id",
					"visibility",
					"in_reply_to_uri",
					"federated",
				).
				ColumnExpr("? DESC", bun.Ident("id")).
				Where("? = ?", bun.Ident("local"), true).
				Where("? IS NULL", bun.Ident("in_reply_to_uri")).
				Where("? = ?", bun.Ident("federated"), true).
				IfNotExists().
				Exec(ctx); err != nil {
				return err
			}

			return nil
		})
	}

	down := func(ctx context.Context, db *bun.DB) error {
		return db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
			return nil
		})
	}

	if err := Migrations.Register(up, down); err != nil {
		panic(err)
	}
}
