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

package search

import (
	"context"
	"errors"
	"net/url"
	"strings"
	"unicode"

	"code.superseriousbusiness.org/gotosocial/internal/db"
	"code.superseriousbusiness.org/gotosocial/internal/gtscontext"
	"code.superseriousbusiness.org/gotosocial/internal/gtserror"
	"code.superseriousbusiness.org/gotosocial/internal/gtsmodel"
	"code.superseriousbusiness.org/gotosocial/internal/paging"
)

// queryStatuses queries for a page of statuses for requesting account
// according to given arguments. if the query is parseable as a URL,
// it attempts a direct lookup and the 'direct' return flag will be set.
//
// IMPORTANT NOTE: this DOES NOT perform all required
// visibility checks on statuses for requesting account,
// though basic ones are performed at the database layer.
// the caller is expected to do this and use the full
// slice pre-checks to calculate paging parameters.
func (p *Processor) queryStatuses(
	ctx context.Context,
	requester *gtsmodel.Account,
	query string,
	resolve bool,
	follows bool,
	authorID string,
	page *paging.Page,
) (
	statuses []*gtsmodel.Status,
	direct bool,
	errWithCode gtserror.WithCode,
) {
	// If this contains spaces,
	// it cannot be a valid URL.
	if strings.ContainsFunc(query,
		unicode.IsSpace) { //nolint

		// Check if this is a valid status URL / URI.
	} else if url := tryParseURL(query); url != nil {

		// Attempt to lookup status by URL.
		status, err := p.getStatusByURI(ctx,
			requester,
			url,
			query,
			resolve,
		)
		if err != nil {
			return nil, false, gtserror.NewErrorInternalError(err)
		}

		if status != nil {
			// Return directly fetched status.
			return []*gtsmodel.Status{status},
				true, nil
		}
	}

	// Parse query string and argumentss as status query type.
	statusQuery, err := parseSearchQuery(query, follows, authorID)
	switch err {
	case nil:
	case errNoOutput:
		return nil, false, nil
	default:
		return nil, false, gtserror.NewErrorBadRequest(
			err, // time parsing error of some kind
			err.Error(),
		)
	}

	// If the 'from:' query operator was set,
	// we need to parse as namestring and fetch
	// it from the database to get an author ID.
	if from := statusQuery.From; //
	from != "" && statusQuery.AuthorID == "" {
		var ok bool
		var username, domain string

		// If this looks like a namestring,
		// attempt to separate into parts.
		if strings.Contains(from, "@") {
			username, domain, ok = tryParseUsernameDomain(from)
		}

		if !ok {
			// Else query for
			// local username.
			username = from
			domain = ""
		}

		// Try lookup the 'from' account by its username @ domain.
		account, err := p.state.DB.GetAccountByUsernameDomain(
			gtscontext.SetBarebones(ctx),
			username,
			domain,
		)
		if err != nil && !errors.Is(err, db.ErrNoEntries) {
			err := gtserror.Newf("db error getting account %s: %w", from, err)
			return nil, false, gtserror.NewErrorInternalError(err)
		}

		// If no account was found,
		// there will be no results.
		if account == nil {
			return
		}

		// Set the account as author ID.
		statusQuery.AuthorID = account.ID
	}

	// Search for statuses by query parameters.
	statuses, err = p.dbSearchForStatuses(ctx,
		requester,
		statusQuery,
		page,
	)
	if err != nil {
		err := gtserror.Newf("error searching for status by %q: %w", query, err)
		return nil, false, gtserror.NewErrorInternalError(err)
	}

	return statuses, false, nil
}

// getStatusByURI fetches a status matching given URI (or
// URL), dereferencing where necessary if resolve is set.
func (p *Processor) getStatusByURI(
	ctx context.Context,
	requester *gtsmodel.Account,
	url *url.URL,
	str string,
	resolve bool,
) (
	status *gtsmodel.Status,
	err error,
) {
	if resolve {
		// Resolve was set, we can dereference by URI / URL.
		status, _, err = p.federator.GetStatusByURI(ctx,
			requester,
			url,
		)
		if err != nil {

			// For known network / request failure type
			// errors, i.e. account doesn't exist or the
			// request failed etc, simply return none.
			if gtserror.IsUnretrievable(err) ||
				gtserror.IsWrongType(err) ||
				gtserror.IsMalformed(err) {
				return nil, nil
			}

			return nil, gtserror.Newf("error dereferencing status %s: %w", str, err)
		}
	} else {

		// Not allowed to resolve, check db by URI string.
		status, err = p.state.DB.GetStatusByURI(ctx, str)
		if err != nil && !errors.Is(err, db.ErrNoEntries) {
			return nil, gtserror.Newf("db error getting status by uri %s: %w", str, err)
		}

		if status == nil {
			// Next search for status by its URL string.
			status, err = p.state.DB.GetStatusByURL(ctx, str)
			if err != nil && !errors.Is(err, db.ErrNoEntries) {
				return nil, gtserror.Newf("db error getting status by url %s: %w", str, err)
			}
		}
	}

	return status, nil
}
