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

	"code.superseriousbusiness.org/gopkg/xslices"
	apimodel "code.superseriousbusiness.org/gotosocial/internal/api/model"
	"code.superseriousbusiness.org/gotosocial/internal/gtserror"
	"code.superseriousbusiness.org/gotosocial/internal/gtsmodel"
	"code.superseriousbusiness.org/gotosocial/internal/typeutils"
)

// Get gets the given status, taking account of privacy settings and blocks etc.
func (p *Processor) Get(ctx context.Context, requester *gtsmodel.Account, statusID string) (*apimodel.Status, gtserror.WithCode) {
	target, errWithCode := p.c.GetVisibleTargetStatus(ctx,
		requester,
		statusID,
		nil, // default freshness
	)
	if errWithCode != nil {
		return nil, errWithCode
	}
	return p.c.GetAPIStatus(ctx, requester, target)
}

// GetMultiple gets the given statuses with the same semantics as Get. Missing or invisible statuses are omitted.
func (p *Processor) GetMultiple(ctx context.Context, requester *gtsmodel.Account, statusIDs []string) ([]*apimodel.Status, gtserror.WithCode) {

	// Without auth, just
	// return equivalent of
	// 404 not found for all.
	if requester == nil {
		return nil, nil
	}

	// Ensure we've only got unique statuses.
	statusIDs = xslices.Deduplicate(statusIDs)

	// Fetch the requested statues by IDs from the database.
	statuses, err := p.state.DB.GetStatusesByIDs(ctx, statusIDs)
	if err != nil {
		err := gtserror.Newf("db error getting status(es): %w", err)
		return nil, gtserror.NewErrorInternalError(err)
	}

	// Check for empty return.
	if len(statuses) == 0 {
		return nil, nil
	}

	// Enqueue refresh of all statuses.
	for _, status := range statuses {
		p.federator.Dereferencer.RefreshStatusAsync(ctx,
			requester, status, nil, nil)
	}

	// Perform visibility checks and return appropriate API models.
	return p.c.GetVisibleAPIStatuses(ctx, requester, statuses, 0), nil
}

// SourceGet returns the *apimodel.StatusSource version of the targetStatusID.
// Status must belong to the requester, and must not be a boost.
func (p *Processor) SourceGet(ctx context.Context, requester *gtsmodel.Account, statusID string) (*apimodel.StatusSource, gtserror.WithCode) {
	status, errWithCode := p.c.GetOwnStatus(ctx, requester, statusID)
	if errWithCode != nil {
		return nil, errWithCode
	}
	if status.BoostOfID != "" {
		return nil, gtserror.NewErrorNotFound(
			errors.New("status is a boost wrapper"),
			"target status not found",
		)
	}

	// Try to use unparsed content
	// warning text if available,
	// fall back to parsed cw html.
	var spoilerText string
	if status.ContentWarningText != "" {
		spoilerText = status.ContentWarningText
	} else {
		spoilerText = status.ContentWarning
	}

	return &apimodel.StatusSource{
		ID:          status.ID,
		Text:        status.Text,
		SpoilerText: spoilerText,
		ContentType: typeutils.ContentTypeToAPIContentType(status.ContentType),
	}, nil
}
