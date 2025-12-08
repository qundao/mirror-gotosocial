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
	"errors"
	"slices"

	"code.superseriousbusiness.org/gopkg/log"
	"code.superseriousbusiness.org/gotosocial/internal/config"
	"code.superseriousbusiness.org/gotosocial/internal/db"
	"code.superseriousbusiness.org/gotosocial/internal/gtscontext"
	"code.superseriousbusiness.org/gotosocial/internal/gtserror"
	"code.superseriousbusiness.org/gotosocial/internal/gtsmodel"
	"code.superseriousbusiness.org/gotosocial/internal/paging"
	"code.superseriousbusiness.org/gotosocial/internal/util"
	"github.com/uptrace/bun"
)

func (d *domainDB) getDomainLimit(
	ctx context.Context,
	lookup string,
	dbQuery func(*gtsmodel.DomainLimit) error,
	keyParts ...any,
) (*gtsmodel.DomainLimit, error) {
	// Fetch perm limit from database cache with loader callback.
	permLimit, err := d.state.Caches.DB.DomainLimit.LoadOne(
		lookup,
		// Only called if not cached.
		func() (*gtsmodel.DomainLimit, error) {
			var permLimit gtsmodel.DomainLimit
			if err := dbQuery(&permLimit); err != nil {
				return nil, err
			}
			return &permLimit, nil
		},
		keyParts...,
	)
	if err != nil {
		return nil, err
	}

	if gtscontext.Barebones(ctx) {
		// No need to fully populate.
		return permLimit, nil
	}

	if permLimit.CreatedByAccount == nil {
		// Not set, fetch from database.
		permLimit.CreatedByAccount, err = d.state.DB.GetAccountByID(
			gtscontext.SetBarebones(ctx),
			permLimit.CreatedByAccountID,
		)
		if err != nil {
			return nil, gtserror.Newf("error populating created by account: %w", err)
		}
	}

	return permLimit, nil
}

func (d *domainDB) GetDomainLimitByID(
	ctx context.Context,
	id string,
) (*gtsmodel.DomainLimit, error) {
	return d.getDomainLimit(
		ctx,
		"ID",
		func(permLimit *gtsmodel.DomainLimit) error {
			return d.db.
				NewSelect().
				Model(permLimit).
				Where("? = ?", bun.Ident("domain_limit.id"), id).
				Scan(ctx)
		},
		id,
	)
}

func (d *domainDB) GetDomainLimitByDomain(
	ctx context.Context,
	domain string,
) (*gtsmodel.DomainLimit, error) {
	// Normalize domain as punycode for lookup.
	domain, err := util.Punify(domain)
	if err != nil {
		return nil, gtserror.Newf("error punifying domain %s: %w", domain, err)
	}

	// Check for easy case, domain referencing *us*
	if domain == "" || domain == config.GetAccountDomain() ||
		domain == config.GetHost() {
		return nil, db.ErrNoEntries
	}

	return d.getDomainLimit(
		ctx,
		"Domain",
		func(permLimit *gtsmodel.DomainLimit) error {
			return d.db.
				NewSelect().
				Model(permLimit).
				Where("? = ?", bun.Ident("domain_limit.domain"), domain).
				Scan(ctx)
		},
		domain,
	)
}

func (d *domainDB) GetDomainLimits(
	ctx context.Context,
	page *paging.Page,
) (
	[]*gtsmodel.DomainLimit,
	error,
) {
	var (
		// Get paging params.
		minID = page.GetMin()
		maxID = page.GetMax()
		limit = page.GetLimit()
		order = page.GetOrder()

		// Make educated guess for slice size
		permLimitIDs = make([]string, 0, limit)
	)

	q := d.db.
		NewSelect().
		TableExpr(
			"? AS ?",
			bun.Ident("domain_limits"),
			bun.Ident("domain_limit"),
		).
		// Select only IDs from table
		Column("domain_limit.id")

	// Return only items with id
	// lower than provided maxID.
	if maxID != "" {
		q = q.Where(
			"? < ?",
			bun.Ident("domain_limit.id"),
			maxID,
		)
	}

	// Return only items with id
	// greater than provided minID.
	if minID != "" {
		q = q.Where(
			"? > ?",
			bun.Ident("domain_limit.id"),
			minID,
		)
	}

	if limit > 0 {
		// Limit amount of
		// items returned.
		q = q.Limit(limit)
	}

	if order == paging.OrderAscending {
		// Page up.
		q = q.OrderExpr(
			"? ASC",
			bun.Ident("domain_limit.id"),
		)
	} else {
		// Page down.
		q = q.OrderExpr(
			"? DESC",
			bun.Ident("domain_limit.id"),
		)
	}

	if err := q.Scan(ctx, &permLimitIDs); err != nil {
		return nil, err
	}

	// Catch case of no items early
	if len(permLimitIDs) == 0 {
		return nil, db.ErrNoEntries
	}

	// If we're paging up, we still want items
	// to be sorted by ID desc, so reverse slice.
	if order == paging.OrderAscending {
		slices.Reverse(permLimitIDs)
	}

	// Allocate return slice (will be at most len permLimitIDs)
	permLimits := make([]*gtsmodel.DomainLimit, 0, len(permLimitIDs))
	for _, id := range permLimitIDs {
		permLimit, err := d.GetDomainLimitByID(ctx, id)
		if err != nil {
			log.Errorf(ctx, "error getting domain limit %q: %v", id, err)
			continue
		}

		// Append to return slice
		permLimits = append(permLimits, permLimit)
	}

	return permLimits, nil
}

func (d *domainDB) MatchDomainLimit(
	ctx context.Context,
	domain string,
) (*gtsmodel.DomainLimit, error) {
	// Normalize domain as punycode for lookup.
	domain, err := util.Punify(domain)
	if err != nil {
		return nil, gtserror.Newf("error punifying domain %s: %w", domain, err)
	}

	// Domain referencing *us* cannot be limited.
	if domain == "" || domain == config.GetAccountDomain() ||
		domain == config.GetHost() {
		return nil, nil
	}

	// Check the domain limited cache for a limit covering the given
	// domain, hydrating the cache with the load function if needed.
	matchedOn, err := d.state.Caches.DB.DomainLimited.MatchesOn(
		domain,
		func() ([]string, error) {
			var domains []string

			// Scan list of all
			// limited domains from DB
			q := d.db.NewSelect().
				Table("domain_limits").
				Column("domain")
			if err := q.Scan(ctx, &domains); err != nil {
				return nil, err
			}

			return domains, nil
		},
	)
	if err != nil {
		return nil, gtserror.Newf("error matching domain %s: %w", domain, err)
	}

	if matchedOn == "" {
		// No match!
		return nil, nil
	}

	// Match was found, fetch the domain limit entry from
	// the database so the caller can do stuff with it.
	return d.GetDomainLimitByDomain(ctx, matchedOn)
}

func (d *domainDB) PutDomainLimit(
	ctx context.Context,
	limit *gtsmodel.DomainLimit,
) error {
	var err error

	// Normalize the domain as punycode, note the extra
	// validation step for domain name write operations.
	limit.Domain, err = util.PunifySafely(limit.Domain)
	if err != nil {
		return gtserror.Newf("error punifying domain %s: %w", limit.Domain, err)
	}

	// Store the domain limit using cache.
	if err := d.state.Caches.DB.DomainLimit.Store(
		limit,
		func() error {
			_, err := d.db.
				NewInsert().
				Model(limit).
				Exec(ctx)
			return err
		},
	); err != nil {
		return err
	}

	// Clear the domain limited cache,
	// will be reloaded later on demand.
	d.state.Caches.DB.DomainLimited.Clear()

	return nil
}

func (d *domainDB) UpdateDomainLimit(
	ctx context.Context,
	limit *gtsmodel.DomainLimit,
	columns ...string,
) error {
	var err error

	// Normalize the domain as punycode, note the extra
	// validation step for domain name write operations.
	limit.Domain, err = util.PunifySafely(limit.Domain)
	if err != nil {
		return gtserror.Newf("error punifying domain %s: %w", limit.Domain, err)
	}

	// Update the domain limit using cache.
	if err := d.state.Caches.DB.DomainLimit.Store(
		limit,
		func() error {
			_, err := d.db.
				NewUpdate().
				Model(limit).
				Column(columns...).
				Where("? = ?", bun.Ident("domain_limit.id"), limit.ID).
				Exec(ctx)
			return err
		},
	); err != nil {
		return err
	}

	// Clear the domain limited cache,
	// will be reloaded later on demand.
	d.state.Caches.DB.DomainLimited.Clear()

	return nil
}

func (d *domainDB) DeleteDomainLimit(
	ctx context.Context,
	id string,
) error {
	// Delete the permLimit from DB.
	q := d.db.NewDelete().
		TableExpr(
			"? AS ?",
			bun.Ident("domain_limits"),
			bun.Ident("domain_limit"),
		).
		Where(
			"? = ?",
			bun.Ident("domain_limit.id"),
			id,
		)

	_, err := q.Exec(ctx)
	if err != nil && !errors.Is(err, db.ErrNoEntries) {
		return err
	}

	// Invalidate any cached model by ID.
	d.state.Caches.DB.DomainLimit.Invalidate("ID", id)

	// Clear the domain limited cache,
	// will be reloaded later on demand.
	d.state.Caches.DB.DomainLimited.Clear()

	return nil
}
