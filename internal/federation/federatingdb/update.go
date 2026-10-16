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

package federatingdb

import (
	"context"
	"errors"

	"code.superseriousbusiness.org/activity/streams/vocab"
	"code.superseriousbusiness.org/gopkg/log"
	"code.superseriousbusiness.org/gotosocial/internal/ap"
	"code.superseriousbusiness.org/gotosocial/internal/config"
	"code.superseriousbusiness.org/gotosocial/internal/db"
	"code.superseriousbusiness.org/gotosocial/internal/gtserror"
	"code.superseriousbusiness.org/gotosocial/internal/gtsmodel"
	"code.superseriousbusiness.org/gotosocial/internal/messages"
)

// Update sets an existing entry to the database based on the value's
// id.
//
// Note that Activity values received from federated peers may also be
// updated in the database this way if the Federating Protocol is
// enabled. The client may freely decide to store only the id instead of
// the entire value.
//
// The library makes this call only after acquiring a lock first.
func (f *DB) Update(ctx context.Context, asType vocab.Type) error {
	log.DebugKV(ctx, "update", Serialize{asType})

	// Mark activity as handled.
	f.storeActivityID(asType)

	// Extract relevant values from passed ctx.
	activityContext := getActivityContext(ctx)
	if activityContext.internal {
		return nil // Already processed.
	}

	requesting := activityContext.requesting
	receiving := activityContext.receiving

	if accountable, ok := ap.ToAccountable(asType); ok {
		return f.updateAccountable(ctx, requesting, receiving, accountable)
	}

	if statusable, ok := ap.ToStatusable(asType); ok {
		return f.updateStatusable(ctx, requesting, receiving, statusable)
	}

	log.Debugf(ctx, "unhandled object type: %T", asType)
	return nil
}

func (f *DB) updateAccountable(
	ctx context.Context,
	requesting *gtsmodel.Account,
	receiving *gtsmodel.Account,
	accountable ap.Accountable,
) error {
	// Extract AP URI of the updated Accountable model.
	idProp := accountable.GetJSONLDId()
	if idProp == nil || !idProp.IsIRI() {
		return gtserror.New("invalid id prop")
	}

	// Get the account URI string for checks
	accountURI := idProp.GetIRI()
	accountURIStr := accountURI.String()

	// Don't try to update local accounts.
	if accountURI.Host == config.GetHost() {
		return nil
	}

	// Check that Update was by
	// the account owner themself.
	if accountURIStr != requesting.URI {
		// If the receiver was our instance account, this is
		// likely a json-ld-signed forwarded Update from a relay
		// (the only type of actor our instance account follows).
		//
		// Since we don't do anything with json-ld signatures,
		// we can just throw this away and we'll catch the update
		// next time we dereference the updated account.
		if receiving.IsInstance() {
			log.Debug(ctx, "throwing away Update delivered to instance account")
			return nil
		}

		// This is someone delivering an Update for someone
		// else's account, return 403 Forbidden for this.
		text := "update for " + accountURIStr + " was not delivered by owner"
		return gtserror.NewErrorForbidden(errors.New(text), text)
	}

	// Pass in to the processor the existing version of the requesting
	// account that we have, plus the Accountable representation that
	// was delivered along with the Update, for further asynchronous
	// updating of eg., avatar/header, emojis, etc. The actual db
	// inserts/updates will take place there.
	f.state.Workers.Federator.Queue.Push(&messages.FromFediAPI{
		APObjectType:   ap.ActorPerson,
		APActivityType: ap.ActivityUpdate,
		GTSModel:       requesting,
		APObject:       accountable,
		Requesting:     requesting,
		Receiving:      receiving,
	})

	return nil
}

func (f *DB) updateStatusable(
	ctx context.Context,
	requesting *gtsmodel.Account,
	receiving *gtsmodel.Account,
	statusable ap.Statusable,
) error {
	// Extract AP URI of the updated model.
	idProp := statusable.GetJSONLDId()
	if idProp == nil || !idProp.IsIRI() {
		return gtserror.New("invalid id prop")
	}

	// Get the status URI string for lookups.
	statusURI := idProp.GetIRI()
	statusURIStr := statusURI.String()

	// Don't try to update local statuses.
	if statusURI.Host == config.GetHost() {
		return nil
	}

	// Check if this is a forwarded object, i.e. did
	// the account making the request also create this?
	forwarded := !isSender(statusable, requesting)

	// Get the status we have on file for this URI string.
	status, err := f.state.DB.GetStatusByURI(ctx, statusURIStr)
	if err != nil && !errors.Is(err, db.ErrNoEntries) {
		return gtserror.Newf("error fetching status from db: %w", err)
	}

	if status == nil {
		// We haven't seen this status before, be
		// lenient and handle as a CREATE event.
		return f.createStatusable(ctx,
			requesting,
			receiving,
			statusable,
			forwarded,
		)
	}

	if forwarded {
		// For forwarded updates, set a nil AS
		// status to force refresh from remote.
		statusable = nil
	}

	// Queue an UPDATE NOTE activity to our fedi API worker,
	// this will handle necessary database insertions, etc.
	f.state.Workers.Federator.Queue.Push(&messages.FromFediAPI{
		APObjectType:   ap.ObjectNote,
		APActivityType: ap.ActivityUpdate,
		GTSModel:       status, // original status
		APObject:       (ap.Statusable)(statusable),
		Requesting:     requesting,
		Receiving:      receiving,
	})

	return nil
}
