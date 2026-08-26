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

package account

import (
	"context"
	"errors"

	"code.superseriousbusiness.org/gotosocial/internal/ap"
	apimodel "code.superseriousbusiness.org/gotosocial/internal/api/model"
	"code.superseriousbusiness.org/gotosocial/internal/db"
	"code.superseriousbusiness.org/gotosocial/internal/gtscontext"
	"code.superseriousbusiness.org/gotosocial/internal/gtserror"
	"code.superseriousbusiness.org/gotosocial/internal/gtsmodel"
	"code.superseriousbusiness.org/gotosocial/internal/id"
	"code.superseriousbusiness.org/gotosocial/internal/messages"
	"code.superseriousbusiness.org/gotosocial/internal/uris"
	"code.superseriousbusiness.org/gotosocial/internal/util"
)

// FollowCreate handles a follow request
// to an account, either remote or local.
func (p *Processor) FollowCreate(
	ctx context.Context,
	requestingAccount *gtsmodel.Account,
	form *apimodel.AccountFollowRequest,
) (*apimodel.Relationship, gtserror.WithCode) {
	targetAccount, errWithCode := p.getFollowTarget(ctx, requestingAccount, form.ID)
	if errWithCode != nil {
		return nil, errWithCode
	}

	// Check if a follow exists already.
	if follow, err := p.state.DB.GetFollow(
		gtscontext.SetBarebones(ctx),
		requestingAccount.ID,
		targetAccount.ID,
	); err != nil && !errors.Is(err, db.ErrNoEntries) {
		err = gtserror.Newf("db error checking existing follow: %w", err)
		return nil, gtserror.NewErrorInternalError(err)
	} else if follow != nil {
		// Already follows, update if necessary + return relationship.
		return p.updateFollow(
			ctx,
			requestingAccount,
			form,
			&follow.Flags,
			func() error { return p.state.DB.UpdateFollow(ctx, follow, "flags") },
		)
	}

	// Check if a follow request exists already.
	if followRequest, err := p.state.DB.GetFollowRequest(
		gtscontext.SetBarebones(ctx),
		requestingAccount.ID,
		targetAccount.ID,
	); err != nil && !errors.Is(err, db.ErrNoEntries) {
		err = gtserror.Newf("db error checking existing follow request: %w", err)
		return nil, gtserror.NewErrorInternalError(err)
	} else if followRequest != nil {
		// Already requested, update if necessary + return relationship.
		return p.updateFollow(
			ctx,
			requestingAccount,
			form,
			&followRequest.Flags,
			func() error { return p.state.DB.UpdateFollowRequest(ctx, followRequest, "flags") },
		)
	}

	// Neither follows nor follow requests, so
	// create and store a new follow request.
	followID := id.NewRandomULID()
	followURI := uris.GenerateURIForFollow(
		requestingAccount.PathPrefix(),
		requestingAccount.Username,
		followID,
	)

	var flags gtsmodel.FollowFlags
	flags.SetShowReblogs(util.PtrOrValue(form.Reblogs, true))
	flags.SetNotify(util.PtrOrValue(form.Notify, false))
	fr := &gtsmodel.FollowRequest{
		ID:              followID,
		URI:             followURI,
		AccountID:       requestingAccount.ID,
		Account:         requestingAccount,
		TargetAccountID: form.ID,
		TargetAccount:   targetAccount,
		Flags:           flags,
	}

	// Insert the new follow request.
	if err := p.state.DB.PutFollowRequest(ctx, fr); err != nil {
		err = gtserror.Newf("error creating follow request in db: %s", err)
		return nil, gtserror.NewErrorInternalError(err)
	}

	// And get the new relationship state.
	rel, errWithCode := p.c.APIRelationship(ctx, requestingAccount, form.ID)
	if errWithCode != nil {
		return nil, errWithCode
	}

	// For unlocked accounts on the same instance,
	// we can already optimistically show the follow
	// request as accepted in the returned relationship.
	if targetAccount.IsLocal() && !*targetAccount.Locked {
		rel.Requested = false
		rel.Following = true
		rel.ShowingReblogs = fr.Flags.ShowReblogs()
		rel.Notifying = fr.Flags.Notify()
	}

	// Handle side effects async.
	p.state.Workers.Client.Queue.Push(&messages.FromClientAPI{
		APObjectType:   ap.ActivityFollow,
		APActivityType: ap.ActivityCreate,
		GTSModel:       fr,
		Origin:         requestingAccount,
		Target:         targetAccount,
	})

	return rel, nil
}

// FollowRemove handles the removal of a follow/follow
// request to an account, either remote or local.
func (p *Processor) FollowRemove(
	ctx context.Context,
	requestingAccount *gtsmodel.Account,
	targetAccountID string,
) (*apimodel.Relationship, gtserror.WithCode) {
	targetAccount, errWithCode := p.getFollowTarget(ctx, requestingAccount, targetAccountID)
	if errWithCode != nil {
		return nil, errWithCode
	}

	// Unfollow and deal
	// with side effects.
	if err := p.c.Unfollow(ctx,
		requestingAccount,
		targetAccount,
		true, // sideEffects
	); err != nil {
		err := gtserror.Newf("db error removing follow (req): %w", err)
		return nil, gtserror.NewErrorNotFound(err)
	}

	return p.c.APIRelationship(ctx, requestingAccount, targetAccountID)
}

/*
	Utility functions.
*/

// updateFollow is a utility function for updating an existing
// follow or followRequest with the parameters provided in the
// given form. If nothing changes, this function is a no-op and
// will just return the existing relationship between follow
// origin and follow target account.
func (p *Processor) updateFollow(
	ctx context.Context,
	requestingAccount *gtsmodel.Account,
	form *apimodel.AccountFollowRequest,
	flags *gtsmodel.FollowFlags,
	update func() error,
) (*apimodel.Relationship, gtserror.WithCode) {
	if form.Reblogs == nil && form.Notify == nil {
		// There's nothing to update.
		return p.c.APIRelationship(ctx, requestingAccount, form.ID)
	}

	// Check what we need to
	// update (if anything).
	var updatingFlags bool

	if newReblogs := form.Reblogs; newReblogs != nil && *newReblogs != flags.ShowReblogs() {
		flags.SetShowReblogs(*newReblogs)
		updatingFlags = true
	}

	if newNotify := form.Notify; newNotify != nil && *newNotify != flags.Notify() {
		flags.SetNotify(*newNotify)
		updatingFlags = true
	}

	if !updatingFlags {
		// Nothing actually changed.
		return p.c.APIRelationship(ctx, requestingAccount, form.ID)
	}

	if err := update(); err != nil {
		err = gtserror.Newf("error updating existing follow (request): %w", err)
		return nil, gtserror.NewErrorInternalError(err)
	}

	return p.c.APIRelationship(ctx, requestingAccount, form.ID)
}

// getFollowTarget is a convenience function which:
//   - Checks if account is trying to follow/unfollow itself.
//   - Returns not found if target should not be visible to requester.
//   - Returns target account according to its id.
func (p *Processor) getFollowTarget(
	ctx context.Context,
	requester *gtsmodel.Account,
	targetID string,
) (*gtsmodel.Account, gtserror.WithCode) {
	// Check for requester.
	if requester == nil {
		err := errors.New("no authorized user")
		return nil, gtserror.NewErrorUnauthorized(err)
	}

	// Account can't follow or unfollow itself.
	if requester.ID == targetID {
		err := errors.New("account can't follow or unfollow itself")
		return nil, gtserror.NewErrorNotAcceptable(err)
	}

	// Fetch the target account for requesting user account.
	return p.c.GetVisibleTargetAccount(ctx, requester, targetID)
}
