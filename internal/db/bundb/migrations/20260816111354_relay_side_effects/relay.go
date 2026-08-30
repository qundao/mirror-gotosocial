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

// Small table to keep track of which relay actors
// have relayed which posts, using which announce
// wrapper URIs (if relayed using an Announce).
type RelayedURI struct {
	ID        string `bun:"type:CHAR(26),pk,nullzero,notnull,unique"`
	RelayURI  string `bun:",notnull,nullzero,unique:relayed_uris_relay_uri_status_uri_uniq"`
	StatusURI string `bun:",notnull,nullzero,unique:relayed_uris_relay_uri_status_uri_uniq"`
	BoostURI  string `bun:",nullzero"`
}
