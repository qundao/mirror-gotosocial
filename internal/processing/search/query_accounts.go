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

	"code.superseriousbusiness.org/gotosocial/internal/config"
	"code.superseriousbusiness.org/gotosocial/internal/db"
	"code.superseriousbusiness.org/gotosocial/internal/gtserror"
	"code.superseriousbusiness.org/gotosocial/internal/gtsmodel"
	"code.superseriousbusiness.org/gotosocial/internal/paging"

	errorsv2 "codeberg.org/gruf/go-errors/v2"
)

// queryAccounts queries for a page of accounts for requesting account
// according to given arguments. if the query is parseable as a URL, or
// as a username@domain namestring, it attempts a direct lookup and the
// 'direct' return flag will be set.
//
// IMPORTANT NOTE: this DOES NOT perform all required
// visibility checks on accounts for requesting account,
// though basic ones are performed at the database layer.
// the caller is expected to do this and use the full
// slice pre-checks to calculate paging parameters.
func (p *Processor) queryAccounts(
	ctx context.Context,
	requester *gtsmodel.Account,
	query string,
	resolve bool,
	follows bool,
	page *paging.Page,
) (
	accounts []*gtsmodel.Account,
	direct bool,
	errWithCode gtserror.WithCode,
) {
	// If this contains spaces, it cannot
	// be a valid URL or username@domain.
	if strings.ContainsFunc(query,
		unicode.IsSpace) { //nolint

		// Check if this is a valid account URL / URI.
	} else if url := tryParseURL(query); url != nil {

		// Attempt to lookup account by the URL.
		account, err := p.getAccountByURI(ctx,
			requester,
			url,
			query,
			resolve,
		)
		if err != nil {
			return nil, false, gtserror.NewErrorInternalError(err)
		}

		if account != nil {
			// Return directly fetched account.
			return []*gtsmodel.Account{account},
				true, nil
		}

		// Next, check whether we've been given a username@domain string.
	} else if username, domain, ok := tryParseUsernameDomain(query); ok {

		// Try lookup the account by its username @ domain.
		account, err := p.getAccountByUsernameDomain(ctx,
			requester,
			username,
			domain,
			resolve,
		)
		if err != nil {
			return nil, false, gtserror.NewErrorInternalError(err)
		}

		if account != nil {
			// Return directly fetched account.
			return []*gtsmodel.Account{account},
				true, nil
		}

		// Search for accounts ~username | ~domain.
		accounts, err = p.dbSearchForAccounts(ctx,
			requester,
			&searchQuery{
				CommonSearchParams: db.CommonSearchParams{
					RequesterID: requester.ID,
					Relationship: func() db.RelationType {
						if follows {
							return db.RelationFollowing
						} else { //nolint
							return 0
						}
					}(),
				},
				UsernameQuery: username,
				DomainQuery:   domain,
			},
			page,
		)
		if err != nil {
			err := gtserror.Newf("error searching for account by %q@%q: %w", username, domain, err)
			return nil, false, gtserror.NewErrorInternalError(err)
		}

		return accounts, false, nil
	}

	// Parse query string and args as account query type.
	accountQuery, err := parseSearchQuery(query, follows, "")
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

	// Search for accounts by query parameters.
	accounts, err = p.dbSearchForAccounts(ctx,
		requester,
		accountQuery,
		page,
	)
	if err != nil {
		err := gtserror.Newf("error searching for account by %q: %w", query, err)
		return nil, false, gtserror.NewErrorInternalError(err)
	}

	return accounts, false, nil
}

// getAccountByURI fetches an account matching given URI
// (or URL), dereferencing where necessary if resolve if set.
func (p *Processor) getAccountByURI(
	ctx context.Context,
	requester *gtsmodel.Account,
	uri *url.URL,
	str string, // uri string
	resolve bool,
) (
	account *gtsmodel.Account,
	err error,
) {
	if resolve {
		// Resolve was set, we can dereference by URI / URL.
		account, _, err = p.federator.GetAccountByURI(ctx,
			requester,
			uri,
			true,
		)
		switch {
		case err == nil:

		// Not-found during a search is fine.
		case errors.Is(err, db.ErrNoEntries):
			return nil, nil

		// For known network / request failure type
		// errors, i.e. account doesn't exist or the
		// request failed etc, simply return none.
		case gtserror.IsUnretrievable(err) ||
			gtserror.IsWrongType(err) ||
			gtserror.IsMalformed(err) ||
			gtserror.IsNotFound(err):
			return nil, nil

		default:
			return nil, gtserror.Newf("error dereferencing account %s: %w", str, err)
		}
	} else {

		// Not allowed to resolve, check db for URI first.
		account, err = p.state.DB.GetAccountByURI(ctx, str)
		if err != nil && !errors.Is(err, db.ErrNoEntries) {
			return nil, gtserror.Newf("db error getting account by uri %s: %w", str, err)
		}

		if account == nil {
			// Next see if this is a direct account URL match.
			account, err = p.state.DB.GetOneAccountByURL(ctx, str)
			if err != nil && !errorsv2.IsV2(err,
				db.ErrNoEntries,
				db.ErrMultipleEntries,
			) {
				return nil, gtserror.Newf("db error getting account by url %s: %w", str, err)
			}
		}
	}

	return account, nil
}

// getAccountByUsernameDomain fetches an account matching given
// username@domain, dereferencing where necessary if resolve is set.
func (p *Processor) getAccountByUsernameDomain(
	ctx context.Context,
	requester *gtsmodel.Account,
	username string,
	domain string,
	resolve bool,
) (
	account *gtsmodel.Account,
	err error,
) {
	switch domain {
	case config.GetHost(),
		config.GetAccountDomain():
		// Normalize domain according
		// to database column values.
		domain = ""

		// It also means we do not
		// need to bother resolving.
		resolve = false
	}

	if resolve {
		// Resolve was set, we can dereference by username @ domain.
		account, _, err = p.federator.GetAccountByUsernameDomain(ctx,
			requester,
			username,
			domain,
		)
		if err != nil && !errors.Is(err, db.ErrNoEntries) {

			// For known network / request failure type
			// errors, i.e. account doesn't exist or the
			// request failed etc, simply return none.
			if gtserror.IsUnretrievable(err) ||
				gtserror.IsWrongType(err) ||
				gtserror.IsMalformed(err) ||
				gtserror.IsNotFound(err) {
				return nil, nil
			}

			return nil, gtserror.Newf("error dereferencing account %s@%s: %w", username, domain, err)
		}
	} else {

		// Not allowed to resolve, check db for username@domain.
		account, err = p.state.DB.GetAccountByUsernameDomain(ctx,
			username,
			domain,
		)
		if err != nil && !errors.Is(err, db.ErrNoEntries) {
			return nil, gtserror.Newf("db error getting account %s@%s: %w", username, domain, err)
		}
	}

	return account, nil
}
