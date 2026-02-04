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
	"github.com/uptrace/bun"
)

func init() {
	up := func(ctx context.Context, db *bun.DB) error {
		log.Info(ctx, "creating \"statuses_polls_scheduler_index\", this may take a minute...")

		// CREATE INDEX IF NOT EXISTS "statuses_polls_scheduler_index"
		// ON "statuses" ("poll_id", "local")
		// WHERE ("poll_id" IS NOT NULL) AND ("local" = TRUE)
		_, err := db.NewCreateIndex().
			Table("statuses").
			Index("statuses_polls_scheduler_index").
			Column("poll_id", "local").
			Where("? IS NOT NULL", bun.Ident("poll_id")).
			Where("? = ?", bun.Ident("local"), true).
			IfNotExists().
			Exec(ctx)

		return err
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
