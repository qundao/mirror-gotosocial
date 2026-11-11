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

package bundb

import (
	"context"
	"net/url"
	"slices"

	"code.superseriousbusiness.org/gopkg/log"
	"code.superseriousbusiness.org/gopkg/xslices"
	"code.superseriousbusiness.org/gotosocial/internal/db"
	"code.superseriousbusiness.org/gotosocial/internal/gtscontext"
	"code.superseriousbusiness.org/gotosocial/internal/gtserror"
	"code.superseriousbusiness.org/gotosocial/internal/gtsmodel"
	"code.superseriousbusiness.org/gotosocial/internal/state"
	"github.com/uptrace/bun"
)

type linkDB struct {
	db    *bun.DB
	state *state.State
}

func (l *linkDB) GetLinkByID(ctx context.Context, id string) (*gtsmodel.Link, error) {
	return l.getLink(
		ctx,
		"ID",
		func(link *gtsmodel.Link) error {
			return l.db.NewSelect().
				Model(link).
				Where("? = ?", bun.Ident("link.id"), id).
				Scan(ctx)
		},
		id,
	)
}

func (l *linkDB) GetLinksByIDs(ctx context.Context, ids []string) ([]*gtsmodel.Link, error) {
	// Load all link IDs via cache loader callbacks.
	links, err := l.state.Caches.DB.Link.LoadIDs("ID",
		ids,
		func(uncached []string) ([]*gtsmodel.Link, error) {
			// Preallocate expected length of uncached links.
			links := make([]*gtsmodel.Link, 0, len(uncached))

			// Perform database query scanning
			// the remaining (uncached) IDs.
			if err := l.db.NewSelect().
				Model(&links).
				Where("? IN (?)", bun.Ident("id"), bun.List(uncached)).
				Scan(ctx); err != nil {
				return nil, err
			}

			return links, nil
		},
	)
	if err != nil {
		return nil, err
	}

	// Reorder the links by given IDs
	// slice to ensure in correct order.
	getID := func(l *gtsmodel.Link) string { return l.ID }
	xslices.OrderBy(links, ids, getID)

	if gtscontext.Barebones(ctx) {
		// no need to fully populate.
		return links, nil
	}

	// Populate all loaded links, removing those we fail to
	// populate (removes needing so many nil checks everywhere).
	links = slices.DeleteFunc(links, func(link *gtsmodel.Link) bool {
		if err := l.PopulateLink(ctx, link); err != nil {
			log.Errorf(ctx, "error populating link %s: %v", link.URL, err)
			return true
		}
		return false
	})

	return links, nil
}

func (l *linkDB) GetLinkByURL(ctx context.Context, url string) (*gtsmodel.Link, error) {
	return l.getLink(
		ctx,
		"URL",
		func(link *gtsmodel.Link) error {
			return l.db.NewSelect().
				Model(link).
				Where("? = ?", bun.Ident("link.url"), url).
				Scan(ctx)
		},
		url,
	)
}

func (l *linkDB) GetLinksByURLs(ctx context.Context, urls []string) ([]*gtsmodel.Link, error) {
	// Load all link URLs via cache loader callbacks.
	links, err := l.state.Caches.DB.Link.LoadIDs("URL",
		urls,
		func(uncached []string) ([]*gtsmodel.Link, error) {
			// Preallocate expected length of uncached links.
			links := make([]*gtsmodel.Link, 0, len(uncached))

			// Perform database query scanning
			// the remaining (uncached) urls.
			if err := l.db.NewSelect().
				Model(&links).
				Where("? IN (?)", bun.Ident("url"), bun.List(uncached)).
				Scan(ctx); err != nil {
				return nil, err
			}

			return links, nil
		},
	)
	if err != nil {
		return nil, err
	}

	// Reorder the links by given urls
	// slice to ensure in correct order.
	getID := func(l *gtsmodel.Link) string { return l.URL }
	xslices.OrderBy(links, urls, getID)

	if gtscontext.Barebones(ctx) {
		// no need to fully populate.
		return links, nil
	}

	// Populate all loaded links, removing those we fail to
	// populate (removes needing so many nil checks everywhere).
	links = slices.DeleteFunc(links, func(link *gtsmodel.Link) bool {
		if err := l.PopulateLink(ctx, link); err != nil {
			log.Errorf(ctx, "error populating link %s: %v", link.URL, err)
			return true
		}
		return false
	})

	return links, nil
}

func (l *linkDB) getLink(ctx context.Context, lookup string, dbQuery func(*gtsmodel.Link) error, keyParts ...any) (*gtsmodel.Link, error) {
	link, err := l.state.Caches.DB.Link.LoadOne(lookup, func() (*gtsmodel.Link, error) {
		var link gtsmodel.Link

		// Not cached! Perform database query.
		if err := dbQuery(&link); err != nil {
			return nil, err
		}

		return &link, nil
	}, keyParts...)
	if err != nil {
		// already processed
		return nil, err
	}

	if gtscontext.Barebones(ctx) {
		// Only a barebones model was requested.
		return link, nil
	}

	if err := l.PopulateLink(ctx, link); err != nil {
		return nil, err
	}

	return link, nil
}

func (l *linkDB) PopulateLink(ctx context.Context, link *gtsmodel.Link) error {
	var err error

	if link.URLParsed == nil {
		link.URLParsed, err = url.Parse(link.URL)
		if err != nil {
			return gtserror.Newf("error populating link: %w", err)
		}
	}

	return nil
}

func (l *linkDB) PutOrGetLink(ctx context.Context, link *gtsmodel.Link) (*gtsmodel.Link, error) {
	// Try to store the link.
	err := l.state.Caches.DB.Link.Store(link, func() error {
		res, err := l.db.
			NewInsert().
			Model(link).
			On("CONFLICT (?) DO NOTHING", bun.Ident("url")).
			Exec(ctx)
		if err != nil {
			return err
		}

		ra, err := res.RowsAffected()
		if err != nil {
			return err
		}

		if ra == 0 {
			return db.ErrAlreadyExists
		}

		return nil
	})

	switch err {
	case nil:
		// Return freshly stored link,
		// populating if desired.
		if gtscontext.Barebones(ctx) {
			return link, nil
		} else { //nolint
			return link, l.PopulateLink(ctx, link)
		}

	case db.ErrAlreadyExists:
		// Link already stored, try to return it.
		return l.GetLinkByURL(ctx, link.URL)

	default:
		// Real error.
		return nil, err
	}
}
