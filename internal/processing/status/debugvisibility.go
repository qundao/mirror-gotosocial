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

package status

import (
	"context"
	"errors"
	"net/url"
	"time"

	"code.superseriousbusiness.org/gopkg/log"
	apimodel "code.superseriousbusiness.org/gotosocial/internal/api/model"
	"code.superseriousbusiness.org/gotosocial/internal/cache"
	"code.superseriousbusiness.org/gotosocial/internal/gtserror"
	"code.superseriousbusiness.org/gotosocial/internal/gtsmodel"
	"code.superseriousbusiness.org/gotosocial/internal/util"
)

func (p *Processor) DebugVisibilityGet(
	ctx context.Context,
	requester *gtsmodel.Account,
	statusURI string,
) (
	*apimodel.StatusVisibilityDebugResponse,
	gtserror.WithCode,
) {
	// Don't leak to no-auth, also check empty.
	if requester == nil || statusURI == "" {
		const text = "target status not found"
		return nil, gtserror.NewErrorNotFound(
			errors.New(text),
			text,
		)
	}

	// Try parse string as URL obj.
	uri, err := url.Parse(statusURI)
	if err != nil {
		return nil, gtserror.NewErrorBadRequest(
			gtserror.Newf("invalid status uri: %w", err),
			"invalid status uri",
		)
	}

	// Ensure the provided URL has an acceptable scheme.
	if uri.Scheme != "http" && uri.Scheme != "https" {
		const text = "invalid URL scheme, acceptable schemes are http or https"
		return nil, gtserror.NewErrorBadRequest(errors.New(text), text)
	}

	// Now we know we've been provided a valid URI, try fetch status.
	status, _, err := p.federator.Dereferencer.GetStatusByURI(ctx,
		requester,
		uri,
	)
	if err != nil {
		log.Errorf(ctx, "error fetching status %s: %v", uri, err)
	}

	if status == nil {
		const text = "target status not found"
		return nil, gtserror.NewErrorNotFound(
			errors.New(text),
			text,
		)
	}

	// Start building status vis response.
	var rsp apimodel.StatusVisibilityDebugResponse
	rsp.URI = status.URI
	allocFilters := func() {
		if rsp.Filters == nil {
			rsp.Filters = new(apimodel.StatusFiltersResult)
		}
	}

	// Get status filtering results according to fetching account, for *all* contexts.
	filters, now, err := p.statusFilter.StatusFilterResults(ctx, requester, status)
	if err != nil {
		err := gtserror.Newf("error getting status filter results: %w", err)
		return nil, gtserror.NewErrorInternalError(err)
	}

	// Append filters applied to status under each context to result.
	for _, filter := range filters.Results[cache.KeyContextHome] {
		allocFilters() // lazily allocate filters when needed
		rsp.Filters.Home = append(rsp.Filters.Home, toFilterResult(filter, now))
	}
	for _, filter := range filters.Results[cache.KeyContextPublic] {
		allocFilters() // lazily allocate filters when needed
		rsp.Filters.Public = append(rsp.Filters.Public, toFilterResult(filter, now))
	}
	for _, filter := range filters.Results[cache.KeyContextNotifs] {
		allocFilters() // lazily allocate filters when needed
		rsp.Filters.Notifications = append(rsp.Filters.Notifications, toFilterResult(filter, now))
	}
	for _, filter := range filters.Results[cache.KeyContextThread] {
		allocFilters() // lazily allocate filters when needed
		rsp.Filters.Thread = append(rsp.Filters.Thread, toFilterResult(filter, now))
	}
	for _, filter := range filters.Results[cache.KeyContextAccount] {
		allocFilters() // lazily allocate filters when needed
		rsp.Filters.Account = append(rsp.Filters.Account, toFilterResult(filter, now))
	}

	// Get mute details for the status according to fetching account.
	mute, err := p.muteFilter.StatusMuteDetails(ctx, requester, status)
	if err != nil {
		err := gtserror.Newf("error getting status mute results: %w", err)
		return nil, gtserror.NewErrorInternalError(err)
	}

	if mute.Mute {
		// Convert mute expiry time to a mute result obj.
		rsp.Mute = toMuteResult(mute.MuteExpiry, now)
	}

	if mute.Notifications {
		// Convert notifications expiry time to a notifications result obj.
		rsp.MuteNotifications = toMuteResult(mute.NotificationExpiry, now)
	}

	// Check whether status is generally visible to the requesting authed account.
	rsp.Visibility.General, err = p.visFilter.StatusVisible(ctx, requester, status)
	if err != nil {
		err := gtserror.Newf("error getting status visibility: %w", err)
		return nil, gtserror.NewErrorInternalError(err)
	}

	// Check whether status should be visible to authed account on their public timelines.
	rsp.Visibility.Public, err = p.visFilter.StatusPublicTimelineable(ctx, requester, status)
	if err != nil {
		err := gtserror.Newf("error getting status public visibility: %w", err)
		return nil, gtserror.NewErrorInternalError(err)
	}

	// Check whether status should be visible to authed account on their home timelines.
	rsp.Visibility.Home, err = p.visFilter.StatusHomeTimelineable(ctx, requester, status)
	if err != nil {
		err := gtserror.Newf("error getting status home visibility: %w", err)
		return nil, gtserror.NewErrorInternalError(err)
	}

	// Check whether status should be visible to authed account on any tag timelines.
	rsp.Visibility.Tag, err = p.visFilter.StatusTagTimelineable(ctx, requester, status)
	if err != nil {
		err := gtserror.Newf("error getting status tag visibility: %w", err)
		return nil, gtserror.NewErrorInternalError(err)
	}

	return &rsp, nil
}

func toFilterResult(filter cache.StatusFilterResult, now time.Time) apimodel.StatusFilterResult {
	var expiry *string
	active := true
	if !filter.Expiry.IsZero() {
		active = !filter.Expired(now)
		format := util.FormatISO8601(filter.Expiry)
		expiry = &format
	}
	return apimodel.StatusFilterResult{
		Active:  active,
		Result:  filter.Result,
		Expires: expiry,
	}
}

func toMuteResult(expiresAt time.Time, now time.Time) *apimodel.StatusMuteResult {
	var expiry *string
	active := true
	if !expiresAt.IsZero() {
		active = expiresAt.After(now)
		format := util.FormatISO8601(expiresAt)
		expiry = &format
	}
	return &apimodel.StatusMuteResult{
		Active:  active,
		Expires: expiry,
	}
}
