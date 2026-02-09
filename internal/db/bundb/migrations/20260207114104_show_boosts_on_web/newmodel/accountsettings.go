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

package gtsmodel

import (
	"time"

	commonmodel "code.superseriousbusiness.org/gotosocial/internal/db/bundb/migrations/20260207114104_show_boosts_on_web"
)

type AccountSettings struct {
	AccountID                      string                         `bun:"type:CHAR(26),pk,nullzero,notnull,unique"`
	CreatedAt                      time.Time                      `bun:"type:timestamptz,nullzero,notnull,default:current_timestamp"`
	UpdatedAt                      time.Time                      `bun:"type:timestamptz,nullzero,notnull,default:current_timestamp"`
	Privacy                        commonmodel.Visibility         `bun:",nullzero,default:3"`
	Sensitive                      *bool                          `bun:",nullzero,notnull,default:false"`
	Language                       string                         `bun:",nullzero,notnull,default:'en'"`
	StatusContentType              string                         `bun:",nullzero"`
	Theme                          string                         `bun:",nullzero"`
	CustomCSS                      string                         `bun:",nullzero"`
	EnableRSS                      *bool                          `bun:",nullzero,notnull,default:false"`
	HideCollections                *bool                          `bun:",nullzero,notnull,default:false"`
	WebLayout                      WebLayout                      `bun:",nullzero,notnull,default:1"`
	InteractionPolicyDirect        *commonmodel.InteractionPolicy `bun:""`
	InteractionPolicyMutualsOnly   *commonmodel.InteractionPolicy `bun:""`
	InteractionPolicyFollowersOnly *commonmodel.InteractionPolicy `bun:""`
	InteractionPolicyUnlocked      *commonmodel.InteractionPolicy `bun:""`
	InteractionPolicyPublic        *commonmodel.InteractionPolicy `bun:""`

	// Added in this migration.
	WebIncludeBoosts *bool `bun:",nullzero,notnull,default:false"`
}

type WebLayout commonmodel.EnumType
