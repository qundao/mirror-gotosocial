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

	"code.superseriousbusiness.org/gotosocial/internal/config"
	gtsmodel "code.superseriousbusiness.org/gotosocial/internal/db/bundb/migrations/20251201113019_swap_application_and_service_actor_types"
	"github.com/uptrace/bun"
)

func init() {
	up := func(ctx context.Context, db *bun.DB) error {
		return db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {

			// If an instance account has already been created,
			// swap its actor type from Service to Application.
			// See: https://codeberg.org/superseriousbusiness/gotosocial/issues/4565
			if _, err := tx.NewUpdate().
				Table("accounts").
				Set("? = ?", bun.Ident("actor_type"), gtsmodel.AccountActorTypeApplication).
				Where("? = ?", bun.Ident("username"), config.GetHost()).
				Exec(ctx); err != nil {
				return err
			}

			// If we have any local bot accounts, swap
			// their actor type from Application to Service.
			//
			// To ensure we only update our own accounts,
			// and not other remote accounts or the instance
			// account, join on the users table.
			if _, err := tx.NewUpdate().
				Table("accounts").
				Set("? = ?", bun.Ident("actor_type"), gtsmodel.AccountActorTypeService).
				TableExpr("? AS ?", bun.Ident("accounts"), bun.Ident("account")).
				Join(
					"INNER JOIN ? AS ? ON ? = ?",
					bun.Ident("users"), bun.Ident("user"),
					bun.Ident("account.id"), bun.Ident("user.account_id"),
				).
				Where("? = ?", bun.Ident("account.actor_type"), gtsmodel.AccountActorTypeApplication).
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
