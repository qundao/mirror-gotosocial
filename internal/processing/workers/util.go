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

package workers

import (
	"context"
	"errors"

	"code.superseriousbusiness.org/gopkg/log"
	apimodel "code.superseriousbusiness.org/gotosocial/internal/api/model"
	"code.superseriousbusiness.org/gotosocial/internal/db"
	"code.superseriousbusiness.org/gotosocial/internal/gtscontext"
	"code.superseriousbusiness.org/gotosocial/internal/gtserror"
	"code.superseriousbusiness.org/gotosocial/internal/gtsmodel"
	"code.superseriousbusiness.org/gotosocial/internal/processing/account"
	"code.superseriousbusiness.org/gotosocial/internal/processing/media"
	"code.superseriousbusiness.org/gotosocial/internal/state"
	"code.superseriousbusiness.org/gotosocial/internal/surfacing"
	"code.superseriousbusiness.org/gotosocial/internal/typeutils"
	"code.superseriousbusiness.org/gotosocial/internal/util"
)

// util provides util functions used by both
// the fromClientAPI and fromFediAPI functions.
type utils struct {
	state     *state.State
	media     *media.Processor
	account   *account.Processor
	surfacer  *surfacing.Surfacer
	converter *typeutils.Converter
}

// deleteBoost encapsulates common logic to delete a boost wrapper
// status from the database, caches and any timeline streams.
func (u *utils) deleteBoost(
	ctx context.Context,
	boost *gtsmodel.Status,
) error {

	// Remove the boost from timeline caches / streams.
	u.surfacer.DeleteStatusFromTimelines(ctx, boost.ID)

	// Finally, delete boost wrapper status itself.
	if err := u.state.DB.DeleteStatusBoost(ctx, boost); //
	err != nil {
		return gtserror.Newf("db error deleting boost %s: %w", boost.URI, err)
	}

	return nil
}

// deleteStatus encapsulates common logic used to
// delete a status and any related models from the
// database, caches and timeline streams.
//
// if 'attachments' = true, then attached status
// media attachments will also be deleted, otherwise
// just detached and left in the database.
//
// if 'sinBin' = true, then the status will be copied
// to `sin_bin_statuses` table prior to deletion.
//
// if 'wipe' = true, then status won't just be
// stubbed out in `statuses`, but fully removed
func (u *utils) deleteStatus(
	ctx context.Context,
	status *gtsmodel.Status,
	attachments bool, // delete attached media
	sinBin bool, // copy to sinbin
	wipe bool, // totally wipe, no stub
) error {
	log := log.New().
		WithContext(ctx).
		WithField("uri", status.URI)

	if sinBin {
		// Copy this status to the sin bin before properly deleting it.
		sbStatus, err := u.converter.StatusToSinBinStatus(ctx, status)
		if err != nil {
			log.Errorf("error converting to sinBinStatus: %v", err)
		} else {
			if err := u.state.DB.PutSinBinStatus(ctx, sbStatus); err != nil {
				log.Errorf("db error storing sinBinStatus: %v", err)
			}
		}
	}

	// Get IDs of any boosts referencing this status.
	boostIDs, err := u.state.DB.GetStatusBoostIDs(ctx,
		status.ID)
	if err != nil {
		log.Errorf("db error getting boosts: %v", err)
	}

	if wipe {
		// Fully delete status and related models from database.
		err := u.state.DB.DeleteStatus(ctx, status, attachments)
		if err != nil {
			return gtserror.Newf("db error deleting status %s: %w", status.URI, err)
		}
	} else {
		// Stub out the status and related model to delete it.
		err := u.state.DB.StubStatus(ctx, status, attachments)
		if err != nil {
			return gtserror.Newf("db error stubbing status %s: %w", status.URI, err)
		}
	}

	for _, id := range boostIDs {
		// Remove the boost from any and all timelines.
		u.surfacer.DeleteStatusFromTimelines(ctx, id)
	}

	// Remove the status from timeline caches / streams.
	u.surfacer.DeleteStatusFromTimelines(ctx, status.ID)

	return nil
}

// redirectFollowers redirects all local
// followers of originAcct to targetAcct.
//
// Both accounts must be fully dereferenced
// already, and the Move must be valid.
//
// Return bool will be true if all goes OK.
func (u *utils) redirectFollowers(
	ctx context.Context,
	originAcct *gtsmodel.Account,
	targetAcct *gtsmodel.Account,
) bool {
	// Any local followers of originAcct should
	// send follow requests to targetAcct instead,
	// and have followers of originAcct removed.
	//
	// Select local followers with barebones, since
	// we only need follow.Account and we can get
	// that ourselves.
	followers, err := u.state.DB.GetAccountLocalFollowers(
		gtscontext.SetBarebones(ctx),
		originAcct.ID,
	)
	if err != nil && !errors.Is(err, db.ErrNoEntries) {
		log.Errorf(ctx,
			"db error getting follows targeting originAcct: %v",
			err)
		return false
	}

	for _, follow := range followers {
		// Fetch the local account that
		// owns the follow targeting originAcct.
		if follow.Account, err = u.state.DB.GetAccountByID(
			gtscontext.SetBarebones(ctx),
			follow.AccountID,
		); err != nil {
			log.Errorf(ctx,
				"db error getting follow account %s: %v",
				follow.AccountID, err,
			)
			return false
		}

		// Use the account processor FollowCreate
		// function to send off the new follow,
		// carrying over the Reblogs and Notify
		// values from the old follow to the new.
		//
		// This will also handle cases where our
		// account has already followed the target
		// account, by just updating the existing
		// follow of target account.
		//
		// Also, ensure new follow wouldn't be a
		// self follow, since that will error.
		if follow.AccountID != targetAcct.ID {
			if _, err := u.account.FollowCreate(
				ctx,
				follow.Account,
				&apimodel.AccountFollowRequest{
					ID:      targetAcct.ID,
					Reblogs: util.Ptr(follow.Flags.ShowReblogs()),
					Notify:  util.Ptr(follow.Flags.Notify()),
				},
			); err != nil {
				log.Errorf(ctx,
					"error creating new follow for account %s: %v",
					follow.AccountID, err,
				)
				return false
			}
		}

		// New follow is in the process of
		// sending, remove the existing follow.
		// This will send out an Undo Activity for each Follow.
		if _, err := u.account.FollowRemove(
			ctx,
			follow.Account,
			follow.TargetAccountID,
		); err != nil {
			log.Errorf(ctx,
				"error removing old follow for account %s: %v",
				follow.AccountID, err,
			)
			return false
		}
	}

	return true
}

// storeInteractionRequest ensures that
// the given interaction request for the
// given interaction is stored in the db.
func (u *utils) storeInteractionRequest(
	ctx context.Context,
	intReq *gtsmodel.InteractionRequest,
) error {
	// Lock on the interaction URI.
	unlock := u.state.ProcessingLocks.Lock(intReq.InteractionURI)
	defer unlock()

	// Ensure no req with this URI exists already.
	req, err := u.state.DB.GetInteractionRequestByInteractionURI(ctx, intReq.InteractionURI)
	if err != nil && !errors.Is(err, db.ErrNoEntries) {
		return gtserror.Newf("db error checking for existing interaction request: %w", err)
	}

	if req != nil {
		// Interaction req already exists,
		// no need to do anything else.
		return nil
	}

	// Store interaction request.
	if err := u.state.DB.PutInteractionRequest(ctx, intReq); err != nil {
		return gtserror.Newf("db error storing interaction request: %w", err)
	}

	return nil
}
