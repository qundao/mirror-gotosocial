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
	// Get relayed status from DB or fetch it from remote.
	status, statusable, isNew, err := d.safelyGetOrDerefRelayableStatus(
		ctx,
		instanceAcct,
		relaySubAcct,
		uri,
		d.matchesRelaySubscription,
		// Within the lock, enrich
		// and store relayed status.
		func(
			ctx context.Context,
			status *gtsmodel.Status,
			statusable ap.Statusable,
			_ string,
			isNew bool,
		) error {
			if !isNew {
				// Nothing
				// to do.
				return nil
			}

			// Note: set a new status ID first as some
			// of the peripherals need to refer to it
			// (eg., mention/attachment.status_id etc).
			status.ID = id.NewULIDFromTime(status.CreatedAt)
			if _, err := d.handleStatusPeripherals(ctx,
				instanceAcct.Username,
				uri,
				&gtsmodel.Status{URI: status.URI},
				status,
			); err != nil {
				err := gtserror.Newf("error handling peripheral dereferencing: %w", err)
				return err
			}

			// Store the enriched status.
			if err := d.state.DB.PutStatus(ctx, status); err != nil {
				err := gtserror.Newf("error inserting new status %s: %w", uri, err)
				return err
			}

			return nil
		},
	)
	if err != nil {
		// Couldn't deref or get;
		// error already wrapped.
		return nil, err
	}

	// Deref parents + children.
	d.dereferenceThread(ctx,
		instanceAcct.Username,
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

// GetBoostedStatusFromRelaySubscription is like GetStatusFromRelaySubscription,
// but for Announces sent to us, instead of Creates forwarded to us.
func (d *Dereferencer) GetBoostedStatusFromRelaySubscription(
	ctx context.Context,
	instanceAcct *gtsmodel.Account,
	relaySubAcct *gtsmodel.Account,
	boostWrapper *gtsmodel.Status,
) (*gtsmodel.Status, error) {
	boostOfURI := boostWrapper.BoostOfURI
	boostOfURIStr := boostWrapper.BoostOfURIStr

	// Get boosted status from DB or fetch it from remote.
	status, statusable, isNew, err := d.safelyGetOrDerefRelayableStatus(
		ctx,
		instanceAcct,
		relaySubAcct,
		boostOfURI,
		d.matchesRelaySubscription,
		// Within the lock, check the boost
		// wrapper status for permissivity,
		// and enrich + store boosted status.
		func(
			ctx context.Context,
			status *gtsmodel.Status,
			statusable ap.Statusable,
			_ string,
			isNew bool,
		) error {
			if !isNew {
				// Nothing
				// to do.
				return nil
			}

			// Add the dereffed status to the
			// incoming boost wrapper we were passed,
			// and then check whether relaySubAcct
			// had permission to do the boost.
			boostWrapper.BoostOf = status
			if err := d.allowedIncomingBoost(ctx,
				boostWrapper,
				instanceAcct,
			); err != nil {
				// Already
				// wrapped.
				return err
			}

			// Note: set a new status ID first as some
			// of the peripherals need to refer to it
			// (eg., mention/attachment.status_id etc).
			status.ID = id.NewULIDFromTime(status.CreatedAt)
			if _, err := d.handleStatusPeripherals(ctx,
				instanceAcct.Username,
				boostOfURI,
				&gtsmodel.Status{URI: boostOfURIStr},
				status,
			); err != nil {
				err := gtserror.Newf("error handling peripheral dereferencing: %w", err)
				return err
			}

			// Store the enriched status.
			//
			// Note: we don't bother storing
			// the boost wrapper from relays.
			if err := d.state.DB.PutStatus(ctx, status); err != nil {
				err := gtserror.Newf("error inserting new status %s: %w", boostOfURIStr, err)
				return err
			}

			return nil
		},
	)
	if err != nil {
		// Couldn't deref or get;
		// error already wrapped.
		return nil, err
	}

	// Deref parents + children.
	d.dereferenceThread(ctx,
		instanceAcct.Username,
		boostOfURI,
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

// GetStatusForRelayActor retrieves status with the given URI
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
func (d *Dereferencer) GetStatusForRelayActor(
	ctx context.Context,
	relayActorAcct *gtsmodel.Account,
	relayingAcct *gtsmodel.Account,
	uri *url.URL,
) (*gtsmodel.Status, ap.Statusable, error) {
	status, statusable, _, err := d.safelyGetOrDerefRelayableStatus(ctx,
		relayActorAcct,
		relayingAcct,
		uri,
		d.matchesRelayActorAcct,
		func(
			ctx context.Context,
			status *gtsmodel.Status,
			statusable ap.Statusable,
			inReplyToAccountURI string,
			isNew bool,
		) error {
			// Make sure we have permission to
			// boost the status as well, since
			// we're going to try to relay it.
			if err := d.allowedOutgoingBoost(ctx,
				status,
				relayActorAcct,
			); err != nil {
				// Already
				// wrapped.
				return err
			}

			if !isNew {
				// Nothing
				// else to do.
				return nil
			}

			// Since this is a new status, check if
			// we're subscribed to our own relay
			// actor and therefore need to store it.
			return d.handleLocalSubscriptionMatch(ctx,
				status,
				inReplyToAccountURI,
				relayActorAcct,
			)
		},
	)
	if err != nil {
		// Couldn't deref or get;
		// error already wrapped.
		return nil, nil, err
	}

	return status, statusable, nil
}

// GetBoostedStatusForRelayActor is like GetStatusForRelayActor,
// but for Announces sent to us, instead of Creates forwarded to us.
func (d *Dereferencer) GetBoostedStatusForRelayActor(
	ctx context.Context,
	relayActorAcct *gtsmodel.Account,
	relayingAcct *gtsmodel.Account,
	boostWrapper *gtsmodel.Status,
) (*gtsmodel.Status, ap.Statusable, error) {
	status, statusable, _, err := d.safelyGetOrDerefRelayableStatus(ctx,
		relayActorAcct,
		relayingAcct,
		boostWrapper.BoostOfURI,
		d.matchesRelayActorAcct,
		func(
			ctx context.Context,
			status *gtsmodel.Status,
			statusable ap.Statusable,
			inReplyToAccountURI string,
			isNew bool,
		) error {
			// Add the dereffed status to the
			// incoming boost wrapper we were passed,
			// and then check whether relayingAcct
			// had permission to do the boost.
			boostWrapper.BoostOf = status
			if err := d.allowedIncomingBoost(ctx,
				boostWrapper,
				relayActorAcct,
			); err != nil {
				// Already
				// wrapped.
				return err
			}

			// Make sure we have permission to
			// boost the status as well, since
			// we're going to try to relay it.
			if err := d.allowedOutgoingBoost(ctx,
				status,
				relayActorAcct,
			); err != nil {
				// Already
				// wrapped.
				return err
			}

			if !isNew {
				// Nothing
				// else to do.
				return nil
			}

			// Since this is a new status, check if
			// we're subscribed to our own relay
			// actor and therefore need to store it.
			return d.handleLocalSubscriptionMatch(ctx,
				status,
				inReplyToAccountURI,
				relayActorAcct,
			)
		},
	)
	if err != nil {
		// Couldn't deref or get;
		// error already wrapped.
		return nil, nil, err
	}

	return status, statusable, nil
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
	// approval without any fannying around as we don't
	// bother sending AnnounceRequests from relay actors.
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

// handleLocalSubscription match checks if the given *new*
// status matches a local subscription to *our own* relay
// actor account. If so, it will be enriched + stored in the db.
//
// DO NOT PASS STATUSES TO THIS FUNCTION IF THEY'RE NOT NEW.
func (d *Dereferencer) handleLocalSubscriptionMatch(
	ctx context.Context,
	newStatus *gtsmodel.Status,
	inReplyToAccountURI string,
	relayActorAcct *gtsmodel.Account,
) error {
	localSubMatch, err := d.matchesRelaySubscription(ctx,
		nil, // not relevant.
		relayActorAcct,
		newStatus,
		inReplyToAccountURI,
	)
	if err != nil {
		// Already
		// wrapped.
		return err
	}

	if !localSubMatch {
		return nil
	}

	// We have a subscription to our own relay actor, woah.
	// We need the instance account to do further processing.
	instanceAcct, err := d.state.DB.GetInstanceAccount(ctx, "")
	if err != nil {
		return gtserror.Newf("db error getting instance account: %w", err)
	}

	// Make sure the instance account -- which we use
	// to do relay subscriptions -- follows our own
	// relay actor, ie., the subscription is approved.
	follows, err := d.state.DB.IsFollowing(ctx,
		instanceAcct.ID,
		relayActorAcct.ID,
	)
	if err != nil {
		err := gtserror.Newf("db error checking follow: %w", err)
		return err
	}

	if !follows {
		// Doesn't follow,
		// so don't store.
		return nil
	}

	// Need to reparse the status
	// URI for handleStatusPeripherals.
	newStatusURI, err := url.Parse(newStatus.URI)
	if err != nil {
		err := gtserror.Newf("error parsing new status URI: %w", err)
		return err
	}

	// Note: set a new status ID first as some
	// of the peripherals need to refer to it
	// (eg., mention/attachment.status_id etc).
	newStatus.ID = id.NewULIDFromTime(newStatus.CreatedAt)
	if _, err := d.handleStatusPeripherals(ctx,
		instanceAcct.Username,
		newStatusURI,
		&gtsmodel.Status{URI: newStatus.URI},
		newStatus,
	); err != nil {
		err := gtserror.Newf("error handling peripheral dereferencing: %w", err)
		return err
	}

	// Store the enriched status.
	if err := d.state.DB.PutStatus(ctx, newStatus); err != nil {
		err := gtserror.Newf("error inserting new status %s: %w", newStatus.URI, err)
		return err
	}

	return nil
}

// safelyGetOrDerefRelayableStatus dereferences or retrieves a
// relayable status from the db using a lock on the status's URI
// and performing permissivity and relay relevance checks.
//
// This function specifically does *not* store the status in the db.
// If necessary, the callback function can be used for that.
func (d *Dereferencer) safelyGetOrDerefRelayableStatus(
	ctx context.Context,
	receiver *gtsmodel.Account,
	relayingAcct *gtsmodel.Account,
	uri *url.URL,
	matchesRelay matchesRelayF,
	inLockCallback func(
		ctx context.Context,
		status *gtsmodel.Status,
		statusable ap.Statusable,
		inReplyToAccountURI string,
		isNew bool,
	) error,
) (
	status *gtsmodel.Status,
	statusable ap.Statusable,
	isNew bool,
	err error,
) {
	// Stringify URI once.
	uriStr := uri.String()

	// Acquire per-URI deref lock.
	unlock := d.state.FedLocks.Lock(uriStr)
	defer unlock()

	// Get relay status from
	// DB or fetch from remote.
	var inReplyToAccountURI string
	status, statusable, inReplyToAccountURI, isNew, err = d.getOrDerefRelayableStatus(
		ctx,
		receiver,
		relayingAcct,
		uri, uriStr,
		matchesRelay,
	)
	if err != nil {
		return
	}

	// Execute callback within
	// the lock if it's provided.
	if inLockCallback != nil {
		err = inLockCallback(ctx,
			status,
			statusable,
			inReplyToAccountURI,
			isNew,
		)
	}

	return
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
) (
	status *gtsmodel.Status,
	statusable ap.Statusable,
	inReplyToAccountURI string,
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

	// Ensure this status matches
	// appropriate relay entity.
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
