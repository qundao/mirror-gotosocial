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

package db

import (
	"context"

	"code.superseriousbusiness.org/gotosocial/internal/gtsmodel"
)

type Link interface {
	// GetLinkByID gets one link with the given ID.
	GetLinkByID(ctx context.Context, id string) (*gtsmodel.Link, error)

	// GetLinksByIDs gets links with the given IDs.
	GetLinksByIDs(ctx context.Context, ids []string) ([]*gtsmodel.Link, error)

	// GetLinkByURL gets one link with the given URL.
	GetLinkByURL(ctx context.Context, url string) (*gtsmodel.Link, error)

	// GetLinksByURLs gets links with the given URLs.
	GetLinksByURLs(ctx context.Context, urls []string) ([]*gtsmodel.Link, error)

	// PopulateLink populates the given link.
	PopulateLink(ctx context.Context, link *gtsmodel.Link) error

	// PutOrGetLink tries to store the given link in the database,
	// returning either the existing link or the newly created one.
	PutOrGetLink(ctx context.Context, link *gtsmodel.Link) (*gtsmodel.Link, error)
}
