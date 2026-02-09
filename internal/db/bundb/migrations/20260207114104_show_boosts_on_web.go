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
	"reflect"

	gtsmodel "code.superseriousbusiness.org/gotosocial/internal/db/bundb/migrations/20260207114104_show_boosts_on_web/newmodel"
	"github.com/uptrace/bun"
)

func init() {
	up := func(ctx context.Context, db *bun.DB) error {
		return db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {

			// Generate new column definition from bun.
			accountSettingsType := reflect.TypeOf((*gtsmodel.AccountSettings)(nil))
			colDef, err := getBunColumnDef(tx, accountSettingsType, "WebIncludeBoosts")
			if err != nil {
				return err
			}

			// Add column to AccountSettings table.
			// Its default of false is safe.
			if _, err = tx.NewAddColumn().
				Model((*gtsmodel.AccountSettings)(nil)).
				ColumnExpr(colDef).
				Exec(ctx); // nocollapse
			err != nil {
				return err
			}

			// Create index for including boosts in web view.
			// This is the same as the existing index but
			// doesn't include the "boost_of_id" column.
			err = createIndex(ctx, tx,
				"statuses_profile_web_view_including_boosts_idx",
				"statuses",
				"?, ?, ?, ?, ? DESC",
				bun.Ident("account_id"),
				bun.Ident("visibility"),
				bun.Ident("in_reply_to_uri"),
				bun.Ident("federated"),
				bun.Ident("id"),
			)

			return err
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
