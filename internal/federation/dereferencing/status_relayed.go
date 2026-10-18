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

package dereferencing

import (
	"context"
	"net/url"

	"code.superseriousbusiness.org/gotosocial/internal/ap"
	"code.superseriousbusiness.org/gotosocial/internal/config"
	"code.superseriousbusiness.org/gotosocial/internal/gtserror"
	"code.superseriousbusiness.org/gotosocial/internal/gtsmodel"
	"code.superseriousbusiness.org/gotosocial/internal/id"
	"code.superseriousbusiness.org/gotosocial/internal/transport"
	"code.superseriousbusiness.org/gotosocial/internal/util"
)

// function for checking whether a
// relayed status is matched by a
// relay entity and should be handled.
//
// "receiver" should be an account
// on our instance; "relayingAcct"
// should be a remote account.
type matchesRelayF func(
	ctx context.Context,
	receiver *gtsmodel.Account,
	relayingAcct *gtsmodel.Account,
	status *gtsmodel.Status,
	inReplyToAccountURI string,
) (bool, error)

// Checks if this status matches at least one relay
// subscription targeting relay account relaySubTarget.
//
// Fulfils matchesRelayF.
func (d *Dereferencer) matchesRelaySubscription(
	ctx context.Context,
	_ *gtsmodel.Account,
	relaySubTarget *gtsmodel.Account,
	status *gtsmodel.Status,
	inReplyToAccountURI string,
) (bool, error) {
	sub, err := d.relayFilter.MatchedBySubscription(ctx,
		relaySubTarget,
		status,
		// may be empty string
		inReplyToAccountURI,
	)
	if err != nil {
		// Err already wrapped
		// in MatchedBySubscription.
		return false, err
	}

	return sub != nil, nil
}

// Checks if this status matches the flags
// and matchers of the given relayActorAcct.
//
// Fulfils matchesRelayF.
func (d *Dereferencer) matchesRelayActorAcct(
	ctx context.Context,
	relayActorAcct *gtsmodel.Account,
	relayingAcct *gtsmodel.Account,
	status *gtsmodel.Status,
	inReplyToAccountURI string,
) (bool, error) {
	relayActor, err := d.state.DB.GetRelayActorByURI(ctx, relayActorAcct.URI)
	if err != nil {
		err := gtserror.Newf("db error getting relay actor: %w", err)
		return false, err
	}

	// You can't relay a status to us from
	// someone else's domain, that's creepy.
	if status.Account.Domain != relayingAcct.Domain {
		return false, nil
	}

	return d.relayFilter.MatchedByActor(ctx,
		relayActor,
		status,
		// may be empty string
		inReplyToAccountURI,
	), nil
}

// relayableF is called by safelyGetOrDerefRelayableStatus, after
// matcher checks, to evaluate whether a status should be relayed.
//
// If an error is returned, the status will not be relayed or stored.
type relayableF func(
	ctx context.Context,
	status *gtsmodel.Status,
) error

// GetStatusFromRelaySubscription dereferences and returns
// a status with the given URI only if it is permitted to be
// relayed by matching at least one active relay subscription.
//
// If the status matches, then it will be inserted into
// the db and the rest of its thread will be dereferenced,
// calling newThreadEntryCallback for each parent or child.
// The status will then be returned for further processing.
//
// If the status doesn't match against at least one relay
// subscription, or the relayed status is not relayable
// due to visibility or permissivity then a checkable
// error type will be returned.
//
// All necessary dereferencing will be done using the passed
// instance service account, as that is the actor that
// receives messages from subscribed relay actors.
func (d *Dereferencer) GetStatusFromRelaySubscription(
	ctx context.Context,
	instanceAcct *gtsmodel.Account,
	relaySubAcct *gtsmodel.Account,
	uri *url.URL,
) (*gtsmodel.Status, error) {
	return d.safelyGetOrDerefRelayableStatus(
		ctx,
		instanceAcct,
		relaySubAcct,
		uri,
		d.matchesRelaySubscription,
		// No need to do any further checks
		// here, if it matches the relay
		// subscription that's good enough.
		nil,
	)
}

// GetBoostedStatusFromRelaySubscription is like GetStatusFromRelaySubscription,
// but for Announces sent to us, instead of Creates forwarded to us.
func (d *Dereferencer) GetBoostedStatusFromRelaySubscription(
	ctx context.Context,
	instanceAcct *gtsmodel.Account,
	relaySubAcct *gtsmodel.Account,
	boostWrapper *gtsmodel.Status,
) (*gtsmodel.Status, error) {
	return d.safelyGetOrDerefRelayableStatus(
		ctx,
		instanceAcct,
		relaySubAcct,
		boostWrapper.BoostOfURI,
		d.matchesRelaySubscription,
		func(
			ctx context.Context,
			status *gtsmodel.Status,
		) error {
			// Add the dereffed status to the
			// incoming boost wrapper we were passed,
			// and then check whether relayingAcct
			// had permission to do the boost.
			//
			// Any error here is already wrapped.
			boostWrapper.BoostOf = status
			return d.allowedIncomingBoost(ctx,
				boostWrapper,
				instanceAcct,
			)
		},
	)
}

// GetStatusForRelaying retrieves status with the given URI
// only if it is permitted to be relayed by matching the given
// local relay actor account, and passing permissivity checks.
// This includes a check on whether our relayActorAcct would
// be allowed to boost this status, since that's how our relay
// actors relay statuses to followers.
//
// If there's an error or the retrieved status is not relevant
// or not permitted, a checkable error type will be returned.
//
// If a matching relay subscription exists on this instance targeting
// our own relay actor, then the status will be stored locally too.
func (d *Dereferencer) GetStatusForRelaying(
	ctx context.Context,
	relayActorAcct *gtsmodel.Account,
	relayingAcct *gtsmodel.Account,
	uri *url.URL,
) (*gtsmodel.Status, error) {
	return d.safelyGetOrDerefRelayableStatus(ctx,
		relayActorAcct,
		relayingAcct,
		uri,
		d.matchesRelayActorAcct,
		func(
			ctx context.Context,
			status *gtsmodel.Status,
		) error {
			// Make sure we'll have permission to boost
			// the status, since we're going to relay it.
			//
			// Any error here is already wrapped.
			return d.allowedOutgoingBoost(ctx,
				status,
				relayActorAcct,
			)
		},
	)
}

// GetBoostedStatusForRelaying is like GetStatusForRelaying,
// but for Announces sent to us, instead of Creates forwarded to us.
func (d *Dereferencer) GetBoostedStatusForRelaying(
	ctx context.Context,
	relayActorAcct *gtsmodel.Account,
	relayingAcct *gtsmodel.Account,
	boostWrapper *gtsmodel.Status,
) (*gtsmodel.Status, error) {
	return d.safelyGetOrDerefRelayableStatus(ctx,
		relayActorAcct,
		relayingAcct,
		boostWrapper.BoostOfURI,
		d.matchesRelayActorAcct,
		func(
			ctx context.Context,
			status *gtsmodel.Status,
		) error {
			// Add the dereffed status to the
			// incoming boost wrapper we were passed,
			// and then check whether relayingAcct
			// had permission to do the boost.
			//
			// Any error here is already wrapped.
			boostWrapper.BoostOf = status
			if err := d.allowedIncomingBoost(ctx,
				boostWrapper,
				relayActorAcct,
			); err != nil {
				return err
			}

			// Make sure we'll have permission to boost
			// the status, since we're going to relay it.
			//
			// Any error here is already wrapped.
			return d.allowedOutgoingBoost(ctx,
				status,
				relayActorAcct,
			)
		},
	)
}

// allowedIncomingBoost checks if the given boostWrapper
// status was automatically approved and can therefore be
// considered relayable. Returns wrapped error if not,
// or unwrapped error if something goes wrong.
func (d *Dereferencer) allowedIncomingBoost(
	ctx context.Context,
	boostWrapper *gtsmodel.Status,
	receiver *gtsmodel.Account,
) error {
	boostPermit, err := d.isPermittedStatus(ctx,
		receiver.Username,
		boostWrapper,
	)
	if err != nil {
		err := gtserror.Newf(
			"error checking permissibility for incoming boost of %s: %w",
			boostWrapper.BoostOfURI, err,
		)
		return err
	}

	if !boostPermit {
		err := gtserror.Newf(
			"dropping unpermitted boost: %s",
			boostWrapper.BoostOfURI,
		)
		return gtserror.SetNotPermitted(err)
	}

	if boostWrapper.Flags.PendingApproval() {
		// If the boost wrapper is permitted, but only
		// pending approval, we can't be arsed with it.
		err := gtserror.Newf(
			"dropping pending-approval boost: %s",
			boostWrapper.BoostOfURI,
		)
		return gtserror.SetNotRelevant(err)
	}

	return nil
}

// allowedOutgoingBoost checks whether the relayActorAccount
// is unconditionally permitted to boost the given status,
// ie., it has automatic approval to do so, not conditional
// on presence in a followers/following collection.
func (d *Dereferencer) allowedOutgoingBoost(
	ctx context.Context,
	status *gtsmodel.Status,
	relayActorAcct *gtsmodel.Account,
) error {
	outgoingBoostPermit, err := d.intFilter.StatusBoostable(ctx,
		relayActorAcct,
		status,
	)
	if err != nil {
		err := gtserror.Newf(
			"error checking boostability of %s: %w",
			status.URI, err,
		)
		return err
	}

	// For relay actors we want straight up automatic
	// approval without any fannying around as we don't want
	// to bother sending AnnounceRequests from relay actors.
	if ok := outgoingBoostPermit.AutomaticApproval() &&
		!outgoingBoostPermit.MatchedOnCollection(); !ok {
		err := gtserror.Newf(
			"dropping status we're not permitted to boost: %s",
			status.URI,
		)
		return gtserror.SetNotPermitted(err)
	}

	return nil
}

// safelyGetOrDerefRelayableStatus dereferences or retrieves a
// relayable status from the db using a lock on the status's URI
// and performing permissivity and relay relevance checks.
func (d *Dereferencer) safelyGetOrDerefRelayableStatus(
	ctx context.Context,
	receiver *gtsmodel.Account,
	relayingAcct *gtsmodel.Account,
	uri *url.URL,
	matchesRelay matchesRelayF,
	relayable relayableF,
) (*gtsmodel.Status, error) {
	// Stringify URI once.
	uriStr := uri.String()

	// Acquire per-URI deref lock, wrapping
	// it in doOnce to allow early unlock.
	unlock := d.state.FedLocks.Lock(uriStr)
	unlock = util.DoOnce(unlock)
	defer unlock()

	// Get status from DB or fetch from remote.
	status, statusable, isNew, err := d.getOrDerefRelayableStatus(
		ctx,
		receiver,
		relayingAcct,
		uri, uriStr,
		matchesRelay,
		relayable,
	)
	if err != nil {
		return nil, err
	}

	// If this is a new
	// status, store it.
	if isNew {
		// Note: set a new status ID first as some
		// of the peripherals need to refer to it
		// (eg., mention/attachment.status_id etc).
		status.ID = id.NewULIDFromTime(status.CreatedAt)
		if _, err := d.handleStatusPeripherals(ctx,
			receiver.Username,
			uri,
			&gtsmodel.Status{URI: uriStr},
			status,
		); err != nil {
			err := gtserror.Newf("error handling peripheral dereferencing: %w", err)
			return nil, err
		}

		// Store the enriched status.
		if err := d.state.DB.PutStatus(ctx, status); err != nil {
			err := gtserror.Newf("error inserting new status %s: %w", uriStr, err)
			return nil, err
		}
	}

	// We're done getting,
	// dereffing, and/or
	// storing this status.
	unlock()

	// Outside of the lock,
	// deref parents + children.
	d.dereferenceThread(ctx,
		receiver.Username,
		uri,
		status,
		statusable,
		isNew,
	)

	// Pass status to dereferencer
	// hook for timelining etc.
	d.onStatusDereference(ctx,
		status,
		isNew,
	)

	return status, nil
}

// getOrDerefRelayableStatus gets a status from the db
// by URI, or dereferences the status fresh from remote.
//
// Regardless of how the status was obtained, the provided
// relevance check will be done on the status to ensure that
// it's relayable. If not, a wrapped error type will be returned.
func (d *Dereferencer) getOrDerefRelayableStatus(
	ctx context.Context,
	receiver *gtsmodel.Account,
	relayingAcct *gtsmodel.Account,
	uri *url.URL,
	uriStr string,
	matchesRelay matchesRelayF,
	relayable relayableF,
) (
	status *gtsmodel.Status,
	statusable ap.Statusable,
	isNew bool,
	err error,
) {
	// Check whether the URI is a
	// blocked domain / subdomain.
	var blocked bool
	blocked, err = d.state.DB.IsDomainBlocked(ctx, uri.Host)
	if err != nil {
		err = gtserror.Newf("error checking blocked domain: %w", err)
		return
	}
	if blocked {
		err = gtserror.Newf("%s is blocked", uri.Host)
		err = gtserror.SetUnretrievable(err)
		return
	}

	// We might need a transport to deref the status
	// and/or its parent, so instantiate that once here.
	tsport, tsportErr := d.transportController.NewTransport(
		receiver.PublicKeyURI,
		receiver.PrivateKey,
	)
	if tsportErr != nil {
		err = gtserror.Newf("couldn't create transport: %w", tsportErr)
		return
	}

	// Search the database for an
	// existing status under URI / URL.
	var existing *gtsmodel.Status
	existing, err = d.getStatusFromDB(ctx, uriStr)
	if err != nil {
		// Err already
		// wrapped.
		return
	}

	if existing != nil {
		// We had it stored in the
		// db already, populate it.
		err = d.state.DB.PopulateStatus(ctx, existing)
		if err != nil {
			return
		}
		status = existing
	}

	if status == nil {
		// We don't have the status stored in the db under the given URI.
		// Dereference statusable from remote, checking if we already had this
		// status stored under a different URI (ie., final URI after redirects).
		statusable, existing, uri, err = d.retrieveStatusable(ctx, tsport, uri)
		if err != nil {
			// Err already
			// wrapped.
			return
		}

		if existing != nil {
			// We had the status
			// stored at the final
			// URI after redirects.
			// Just use this.
			status = existing
		} else {
			// We didn't have
			// the status yet.
			isNew = true

			// Convert the statusable to GtS format.
			status, err = d.convertStatusable(ctx,
				receiver.Username,
				uri,
				statusable,
			)
			if err != nil {
				// Err already
				// wrapped.
				return
			}
		}
	}

	// For relay purposes we're only interested
	// in statuses that are public or unlisted.
	// Ignore followers-only + direct.
	if vis := status.Visibility; //
	!(vis == gtsmodel.VisibilityPublic ||
		vis == gtsmodel.VisibilityUnlocked) {
		// Return checkable error type.
		err = gtserror.Newf("dropping non-public/non-unlisted status: %s", uriStr)
		err = gtserror.SetNotRelevant(err)
		return
	}

	// For match checks we have to
	// know the author of the parent
	// status (if there is one).
	var inReplyToAccountURI string
	if status.InReplyToURI != "" {
		if status.InReplyTo != nil {
			// Parent already dereffed, use it.
			inReplyToAccountURI = status.InReplyTo.AccountURI
		} else {
			// Parent not dereffed yet, try to get it.
			inReplyToAccountURI, err = d.retrieveInReplyToAccountURI(ctx, tsport, status)
			if err != nil {
				err = gtserror.Newf("error retrieving inReplyToAccountURI: %w", err)
				err = gtserror.SetUnretrievable(err)
				return
			}
		}
	}

	// Ensure this status matches appropriate
	// relay entity, taking account of matchers,
	// visibility, reply status etc.
	var match bool
	match, err = matchesRelay(ctx,
		receiver,
		relayingAcct,
		status,
		// May be empty string.
		inReplyToAccountURI,
	)
	if err != nil {
		return
	}

	if !match {
		// Return checkable error type.
		err = gtserror.Newf("status %s not matched", uriStr)
		err = gtserror.SetNotRelevant(err)
		return
	}

	// Check if this is a permitted status,
	// bearing in mind interaction policy
	// of replied-to status (if present).
	var permit bool
	permit, err = d.isPermittedStatus(ctx,
		receiver.Username,
		status,
	)
	if err != nil {
		err = gtserror.Newf("error checking permissibility for status %s: %w", uriStr, err)
		return
	}

	if !permit {
		// Return checkable error type.
		err = gtserror.Newf("status %s unpermitted", uriStr)
		err = gtserror.SetNotRelevant(err)
		return
	}

	if status.Flags.PendingApproval() {
		// If the status is pending approval we don't really
		// care about it, we'll get it from somewhere else at
		// some point perhaps. We only care about relayed posts
		// that are already approved or don't require approval.
		err = gtserror.Newf("status %s pending approval", uriStr)
		err = gtserror.SetNotRelevant(err)
		return
	}

	// Do any further relayable checks,
	// eg., was status allowed to be boosted
	// to us, can we boost it out again, etc.
	if relayable != nil {
		if err = relayable(ctx, status); err != nil {
			// Already
			// wrapped.
			return
		}
	}

	// Everything
	// looks OK!
	return
}

// retrieveInReplyToAccountURI dereferences the parent
// of the given status to retrieve the URI of the account.
func (d *Dereferencer) retrieveInReplyToAccountURI(
	ctx context.Context,
	tsport transport.Transport,
	status *gtsmodel.Status,
) (string, error) {
	if status.InReplyTo != nil {
		// Had parent stored already.
		return status.InReplyTo.AccountURI, nil
	}

	// We'll have to dereference.
	inReplyToURI := status.InReplyToURI
	parentURI, err := url.Parse(inReplyToURI)
	if err != nil {
		err := gtserror.Newf("error parsing parent URI: %w", err)
		return "", err
	}

	// Make sure it's not replying to something on our
	// domain that may have been deleted from the db.
	if parentURI.Host == config.GetHost() ||
		parentURI.Host == config.GetAccountDomain() {
		return "", nil
	}

	// Check whether the parent status URI is a blocked
	// domain, we don't wanna be making requests to baddies.
	blocked, err := d.state.DB.IsDomainBlocked(ctx, parentURI.Host)
	if err != nil {
		err := gtserror.Newf("db error checking blocked domain: %w", err)
		return "", err
	}

	if blocked {
		// Just return without
		// dereffing, kiss my ass.
		return "", nil
	}

	// We don't have the parent stored, try to fetch but *don't* store it, only check for now.
	parentStatusable, parentStatus, _, err := d.retrieveStatusable(ctx, tsport, parentURI)
	if err != nil {
		err := gtserror.Newf("error retrieving %s: %w", inReplyToURI, err)
		return "", err
	}

	if parentStatus != nil {
		// We had the parent status
		// stored at a different URI!
		return parentStatus.AccountURI, nil
	}

	// Get attributedTo URI from the dereferenced parent status.
	attributedTo, err := ap.GetOneAttributedTo(parentStatusable)
	if err != nil {
		return "", gtserror.SetMalformed(err)
	}

	return attributedTo.String(), nil
}
