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

import "net/url"

// Link represents one link seen by this
// instance, stored in the database table `links`.
//
// TODO: More fields here once we add OpenGraph
// link previews and link preview thumbnails.
type Link struct {
	// Punified URL of the link.
	URL string `bun:",pk,nullzero,notnull,unique"`

	// Parsed version of URL,
	// not stored in the db.
	URLParsed *url.URL `bun:"-"`
}

// LinkAttachment is a bit like Link but doesn't
// have its own DB table, and is specifically
// for embedding in status database entries.
type LinkAttachment struct {
	// Punified URL of the link attachment.
	//
	// Points back to a gtsmodel.Link
	// entry in the database.
	URL string `json:"url"`

	// True if link was present in
	// AP `status.content` rather
	// than simply being included
	// as separate AP `attachment`.
	InBody bool `json:"in_body,omitempty"`
}
