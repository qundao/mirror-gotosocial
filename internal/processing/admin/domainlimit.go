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

package admin

import (
	"context"
	"errors"
	"slices"
	"strings"

	"code.superseriousbusiness.org/gopkg/xslices"
	apimodel "code.superseriousbusiness.org/gotosocial/internal/api/model"
	"code.superseriousbusiness.org/gotosocial/internal/db"
	"code.superseriousbusiness.org/gotosocial/internal/gtserror"
	"code.superseriousbusiness.org/gotosocial/internal/gtsmodel"
	"code.superseriousbusiness.org/gotosocial/internal/id"
	"code.superseriousbusiness.org/gotosocial/internal/paging"
	"code.superseriousbusiness.org/gotosocial/internal/typeutils"
)

func (p *Processor) DomainLimitsGet(ctx context.Context, page *paging.Page) (*apimodel.PageableResponse, gtserror.WithCode) {
	// Get domain limits.
	domainLimits, err := p.state.DB.GetDomainLimits(ctx, page)
	if err != nil && !errors.Is(err, db.ErrNoEntries) {
		err := gtserror.Newf("db error: %w", err)
		return nil, gtserror.NewErrorInternalError(err)
	}

	count := len(domainLimits)
	if count == 0 {
		return paging.EmptyResponse(), nil
	}

	// Convert each domain
	// limit to API model.
	items := make([]*apimodel.DomainLimit, count)
	for i, domainLimit := range domainLimits {
		apiDomainLimit, err := p.converter.DomainLimitToAPIDomainLimit(ctx, domainLimit)
		if err != nil {
			err := gtserror.Newf("error converting domain limit: %w", err)
			return nil, gtserror.NewErrorInternalError(err)
		}
		items[i] = apiDomainLimit
	}

	var lo, hi string
	if !page.Paging() {
		// If not paging, sort
		// items alphabetically.
		slices.SortFunc(
			items,
			func(a *apimodel.DomainLimit, b *apimodel.DomainLimit) int {
				return strings.Compare(a.Domain, b.Domain)
			},
		)
	} else {
		// If paging, leave sorted by ID,
		// and assemble next/prev queries.
		lo = domainLimits[count-1].ID
		hi = domainLimits[0].ID
	}

	return paging.PackageResponse(paging.ResponseParams{
		Items: xslices.ToAny(items),
		Path:  "/api/v1/admin/domain_limits",
		Next:  page.Next(lo, hi),
		Prev:  page.Prev(lo, hi),
	}), nil
}

func (p *Processor) DomainLimitGet(ctx context.Context, id string) (*apimodel.DomainLimit, gtserror.WithCode) {
	// Get single domain limit.
	domainLimit, err := p.state.DB.GetDomainLimitByID(ctx, id)
	if err != nil && !errors.Is(err, db.ErrNoEntries) {
		err := gtserror.Newf("db error: %w", err)
		return nil, gtserror.NewErrorInternalError(err)
	}

	if domainLimit == nil {
		err := gtserror.Newf("domain limit %s not found", id)
		return nil, gtserror.NewErrorNotFound(err)
	}

	apiDomainLimit, err := p.converter.DomainLimitToAPIDomainLimit(ctx, domainLimit)
	if err != nil {
		err := gtserror.Newf("error converting domain limit: %w", err)
		return nil, gtserror.NewErrorInternalError(err)
	}

	return apiDomainLimit, nil
}

func (p *Processor) DomainLimitCreate(
	ctx context.Context,
	acct *gtsmodel.Account,
	domain string,
	mediaPolicy apimodel.MediaPolicy,
	followsPolicy apimodel.FollowsPolicy,
	statusesPolicy apimodel.StatusesPolicy,
	accountsPolicy apimodel.AccountsPolicy,
	contentWarning string,
	publicComment string,
	privateComment string,
) (*apimodel.DomainLimit, gtserror.WithCode) {

	// Parse policies.
	mp, errWithCode := parseMediaPolicy(mediaPolicy)
	if errWithCode != nil {
		return nil, errWithCode
	}

	fp, errWithCode := parseFollowsPolicy(followsPolicy)
	if errWithCode != nil {
		return nil, errWithCode
	}

	sp, errWithCode := parseStatusesPolicy(statusesPolicy)
	if errWithCode != nil {
		return nil, errWithCode
	}

	ap, errWithCode := parseAccountsPolicy(accountsPolicy)
	if errWithCode != nil {
		return nil, errWithCode
	}

	// Create + store domain limit.
	domainLimit := &gtsmodel.DomainLimit{
		ID:                 id.NewULID(),
		Domain:             domain,
		CreatedByAccountID: acct.ID,
		CreatedByAccount:   acct,
		PrivateComment:     privateComment,
		PublicComment:      publicComment,
		MediaPolicy:        mp,
		FollowsPolicy:      fp,
		StatusesPolicy:     sp,
		AccountsPolicy:     ap,
		ContentWarning:     contentWarning,
	}

	switch err := p.state.DB.PutDomainLimit(ctx, domainLimit); {
	case err == nil:
		// No problem.

	case errors.Is(err, db.ErrAlreadyExists):
		text := "limit with domain " + domain + " already exists"
		return nil, gtserror.NewErrorConflict(errors.New(text), text)

	default:
		err := gtserror.Newf("db error storing domain limit: %w", err)
		return nil, gtserror.NewErrorInternalError(err)
	}

	apiDomainLimit, err := p.converter.DomainLimitToAPIDomainLimit(ctx, domainLimit)
	if err != nil {
		err := gtserror.Newf("error converting domain limit: %w", err)
		return nil, gtserror.NewErrorInternalError(err)
	}

	return apiDomainLimit, nil
}

func (p *Processor) DomainLimitUpdate(
	ctx context.Context,
	id string,
	mediaPolicy *apimodel.MediaPolicy,
	followsPolicy *apimodel.FollowsPolicy,
	statusesPolicy *apimodel.StatusesPolicy,
	accountsPolicy *apimodel.AccountsPolicy,
	contentWarning *string,
	publicComment *string,
	privateComment *string,
) (*apimodel.DomainLimit, gtserror.WithCode) {
	domainLimit, err := p.state.DB.GetDomainLimitByID(ctx, id)
	if err != nil && !errors.Is(err, db.ErrNoEntries) {
		err := gtserror.Newf("db error: %w", err)
		return nil, gtserror.NewErrorInternalError(err)
	}

	if domainLimit == nil {
		err := gtserror.Newf("domain limit %s not found", id)
		return nil, gtserror.NewErrorNotFound(err)
	}

	// Prepare db update columns
	// for selective updating.
	var columns []string

	// Parse policies (if set).
	if mediaPolicy != nil {
		mp, errWithCode := parseMediaPolicy(*mediaPolicy)
		if errWithCode != nil {
			return nil, errWithCode
		}

		domainLimit.MediaPolicy = mp
		columns = append(columns, "media_policy")
	}

	if followsPolicy != nil {
		fp, errWithCode := parseFollowsPolicy(*followsPolicy)
		if errWithCode != nil {
			return nil, errWithCode
		}

		domainLimit.FollowsPolicy = fp
		columns = append(columns, "follows_policy")
	}

	if statusesPolicy != nil {
		sp, errWithCode := parseStatusesPolicy(*statusesPolicy)
		if errWithCode != nil {
			return nil, errWithCode
		}

		domainLimit.StatusesPolicy = sp
		columns = append(columns, "statuses_policy")
	}

	if accountsPolicy != nil {
		ap, errWithCode := parseAccountsPolicy(*accountsPolicy)
		if errWithCode != nil {
			return nil, errWithCode
		}

		domainLimit.AccountsPolicy = ap
		columns = append(columns, "accounts_policy")
	}

	// Parse other nillable fields.
	if contentWarning != nil {
		domainLimit.ContentWarning = *contentWarning
		columns = append(columns, "content_warning")
	}

	if publicComment != nil {
		domainLimit.PublicComment = *publicComment
		columns = append(columns, "public_comment")
	}

	if privateComment != nil {
		domainLimit.PrivateComment = *privateComment
		columns = append(columns, "private_comment")
	}

	// Do the update.
	err = p.state.DB.UpdateDomainLimit(ctx, domainLimit, columns...)
	if err != nil {
		err := gtserror.Newf("db error updating domain limit: %w", err)
		return nil, gtserror.NewErrorInternalError(err)
	}

	apiDomainLimit, err := p.converter.DomainLimitToAPIDomainLimit(ctx, domainLimit)
	if err != nil {
		err := gtserror.Newf("error converting domain limit: %w", err)
		return nil, gtserror.NewErrorInternalError(err)
	}

	return apiDomainLimit, nil
}

func (p *Processor) DomainLimitDelete(
	ctx context.Context,
	id string,
) (*apimodel.DomainLimit, gtserror.WithCode) {
	domainLimit, err := p.state.DB.GetDomainLimitByID(ctx, id)
	if err != nil && !errors.Is(err, db.ErrNoEntries) {
		err := gtserror.Newf("db error: %w", err)
		return nil, gtserror.NewErrorInternalError(err)
	}

	if domainLimit == nil {
		err := gtserror.Newf("domain limit %s not found", id)
		return nil, gtserror.NewErrorNotFound(err)
	}

	// Convert the domain limit to
	// API model before the delete.
	apiDomainLimit, err := p.converter.DomainLimitToAPIDomainLimit(ctx, domainLimit)
	if err != nil {
		err := gtserror.Newf("error converting domain limit: %w", err)
		return nil, gtserror.NewErrorInternalError(err)
	}

	// Do the delete.
	if err := p.state.DB.DeleteDomainLimit(ctx, id); err != nil {
		err := gtserror.Newf("db error deleting domain limit: %w", err)
		return nil, gtserror.NewErrorInternalError(err)
	}

	return apiDomainLimit, nil
}

func parseMediaPolicy(mediaPolicy apimodel.MediaPolicy) (gtsmodel.MediaPolicy, gtserror.WithCode) {
	mp := typeutils.APIMediaPolicyToMediaPolicy(mediaPolicy)
	if mp != gtsmodel.MediaPolicyUnknown {
		return mp, nil
	}

	const text = "media_policy unknown, must be one of no_action (default), mark_sensitive, or reject"
	errWithCode := gtserror.NewErrorBadRequest(errors.New(text), text)
	return 0, errWithCode
}

func parseFollowsPolicy(followsPolicy apimodel.FollowsPolicy) (gtsmodel.FollowsPolicy, gtserror.WithCode) {
	fp := typeutils.APIFollowsPolicyToFollowsPolicy(followsPolicy)
	if fp != gtsmodel.FollowsPolicyUnknown {
		return fp, nil
	}

	const text = "follows_policy unknown, must be one of no_action (default), manual_approval, reject_non_mutual, or reject_all"
	errWithCode := gtserror.NewErrorBadRequest(errors.New(text), text)
	return 0, errWithCode
}

func parseStatusesPolicy(statusesPolicy apimodel.StatusesPolicy) (gtsmodel.StatusesPolicy, gtserror.WithCode) {
	sp := typeutils.APIStatusesPolicyToStatusesPolicy(statusesPolicy)
	if sp != gtsmodel.StatusesPolicyUnknown {
		return sp, nil
	}

	const text = "statuses_policy unknown, must be one of no_action (default), filter_warn, or filter_hide"
	errWithCode := gtserror.NewErrorBadRequest(errors.New(text), text)
	return 0, errWithCode
}

func parseAccountsPolicy(accountsPolicy apimodel.AccountsPolicy) (gtsmodel.AccountsPolicy, gtserror.WithCode) {
	ap := typeutils.APIAccountsPolicyToAccountsPolicy(accountsPolicy)
	if ap != gtsmodel.AccountsPolicyUnknown {
		return ap, nil
	}

	const text = "accounts_policy unknown, must be one of no_action (default), or mute"
	errWithCode := gtserror.NewErrorBadRequest(errors.New(text), text)
	return 0, errWithCode
}
