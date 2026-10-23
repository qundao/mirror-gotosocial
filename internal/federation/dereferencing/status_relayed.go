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

	"code.superseriousbusiness.org/gopkg/log"
	"code.superseriousbusiness.org/gotosocial/internal/ap"
	"code.superseriousbusiness.org/gotosocial/internal/config"
	"code.superseriousbusiness.org/gotosocial/internal/gtserror"
	"code.superseriousbusiness.org/gotosocial/internal/gtsmodel"
	"code.superseriousbusiness.org/gotosocial/internal/transport"
)

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
	instance *gtsmodel.Account,
	relaySubscriber *gtsmodel.Account,
	uri *url.URL,
) (
	*gtsmodel.Status,
	error,
) {
	return d.getRelayStatusByURI(ctx,
		instance,
		relaySubscriber,
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
	instance *gtsmodel.Account,
	relaySubscriber *gtsmodel.Account,
	boost *gtsmodel.Status,
) (
	*gtsmodel.Status,
	error,
) {
	return d.getRelayStatusByURI(ctx,
		instance,
		relaySubscriber,
		boost.BoostOfURI,
		d.matchesRelaySubscription,
		func(ctx context.Context, relayActor *gtsmodel.Account, status *gtsmodel.Status) (bool, error) {
			// Add the dereffed status to the
			// incoming boost wrapper we were passed,
			// and then check whether relayingAcct
			// had permission to do the boost.
			//
			// Any error here is already wrapped.
			boost.BoostOf = status
			return d.allowedIncomingBoost(ctx,
				instance,
				boost,
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
	relayActor *gtsmodel.Account,
	relaying *gtsmodel.Account,
	uri *url.URL,
) (
	*gtsmodel.Status,
	error,
) {
	return d.getRelayStatusByURI(ctx,
		relayActor,
		relaying,
		uri,
		d.matchesRelayActor,
		func(ctx context.Context, relayActor *gtsmodel.Account, status *gtsmodel.Status) (bool, error) {
			// Make sure we'll have permission to boost
			// the status, since we're going to relay it.
			//
			// Any error here is already wrapped.
			return d.allowedOutgoingBoost(ctx,
				relayActor,
				status,
			)
		},
	)
}

// GetBoostedStatusForRelaying is like GetStatusForRelaying,
// but for Announces sent to us, instead of Creates forwarded to us.
func (d *Dereferencer) GetBoostedStatusForRelaying(
	ctx context.Context,
	relayActor *gtsmodel.Account,
	relaying *gtsmodel.Account,
	boost *gtsmodel.Status,
) (
	*gtsmodel.Status,
	error,
) {
	return d.getRelayStatusByURI(ctx,
		relayActor,
		relaying,
		boost.BoostOfURI,
		d.matchesRelayActor,
		func(ctx context.Context, relayActor *gtsmodel.Account, status *gtsmodel.Status) (bool, error) {
			// Add the dereffed status to the
			// incoming boost wrapper we were passed,
			// and then check whether relayingAcct
			// had permission to do the boost.
			//
			// Any error here is already wrapped.
			boost.BoostOf = status
			if ok, err := d.allowedIncomingBoost(ctx,
				relayActor,
				boost,
			); !ok { // includes err != nil
				return false, err
			}

			// Make sure we'll have permission to boost
			// the status, since we're going to relay it.
			//
			// Any error here is already wrapped.
			return d.allowedOutgoingBoost(ctx,
				relayActor,
				status,
			)
		},
	)
}

// Checks if this status matches at least one relay
// subscription targeting relay account relaySubTarget.
func (d *Dereferencer) matchesRelaySubscription(
	ctx context.Context,
	_ *gtsmodel.Account,
	relaySubTarget *gtsmodel.Account,
	status *gtsmodel.Status,
	inReplyToAccountURI string,
) (bool, error) {
	sub, err := d.relaying.MatchedBySubscription(ctx,
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
// and matchers of the given relayActor.
func (d *Dereferencer) matchesRelayActor(
	ctx context.Context,
	relayActor *gtsmodel.Account,
	relaying *gtsmodel.Account,
	status *gtsmodel.Status,
	inReplyToAccountURI string,
) (bool, error) {
	{
		relayActor, err := d.state.DB.GetRelayActorByURI(ctx, relayActor.URI)
		if err != nil {
			return false, gtserror.Newf("db error getting relay actor: %w", err)
		}

		// You can't relay a status to us from
		// someone else's domain, that's creepy.
		if status.Account.Domain != relaying.Domain {
			return false, nil
		}

		return d.relaying.MatchedByActor(ctx,
			relayActor,
			status,

			// this may be empty,
			// but that is fine.
			inReplyToAccountURI,
		), nil
	}
}

// allowedIncomingBoost checks if the given boostWrapper
// status was automatically approved and can therefore be
// considered relayable. Returns wrapped error if not,
// or unwrapped error if something goes wrong.
func (d *Dereferencer) allowedIncomingBoost(
	ctx context.Context,
	receiver *gtsmodel.Account,
	boost *gtsmodel.Status,
) (
	bool,
	error,
) {
	boostPermit, err := d.isPermittedStatus(ctx,
		receiver,
		boost,
	)
	if err != nil {
		return false, gtserror.Newf("error checking permissibility for incoming boost of %s: %w",
			boost.BoostOfURI, err)
	} else if !boostPermit {
		return false, nil
	}

	if boost.Flags.PendingApproval() {
		// If the boost wrapper is permitted, but only
		// pending approval, we can't be arsed with it.
		return false, nil
	}

	return true, nil
}

// allowedOutgoingBoost checks whether the relayActorAccount
// is unconditionally permitted to boost the given status,
// ie., it has automatic approval to do so, not conditional
// on presence in a followers/following collection.
func (d *Dereferencer) allowedOutgoingBoost(
	ctx context.Context,
	relayActor *gtsmodel.Account,
	status *gtsmodel.Status,
) (
	bool,
	error,
) {
	outgoingBoostPermit, err := d.interaction.StatusBoostable(ctx,
		relayActor,
		status,
	)
	if err != nil {
		return false, gtserror.Newf("error checking boostability of %s: %w",
			status.URI, err)
	}

	// For relay actors we want straight up automatic
	// approval without any fannying around as we don't want
	// to bother sending AnnounceRequests from relay actors.
	if ok := outgoingBoostPermit.AutomaticApproval() &&
		!outgoingBoostPermit.MatchedOnCollection(); !ok {
		return false, nil
	}

	return true, nil
}

func (d *Dereferencer) getRelayStatusByURI(
	ctx context.Context,
	relayActor *gtsmodel.Account,
	relaying *gtsmodel.Account,
	uri *url.URL,

	matches func(ctx context.Context,
		relayActor *gtsmodel.Account,
		relaying *gtsmodel.Account,
		status *gtsmodel.Status,
		inReplyToAccountURI string,
	) (
		bool,
		error,
	),

	relayable func(
		ctx context.Context,
		relayActor *gtsmodel.Account,
		status *gtsmodel.Status,
	) (
		bool,
		error,
	),
) (
	*gtsmodel.Status,
	error,
) {
	// Search the database
	// for existing status.
	uriStr := uri.String()
	status, err := d.getStatusFromDB(ctx, uriStr)
	if err != nil {
		return nil, err
	}

	// relevant is the hook passed to enrichAndStoreStatus() below, and also used
	// to check whether an existing up-to-date status is relevant for relay return.
	relevant := func(ctx context.Context, status *gtsmodel.Status) (bool, error) {
		if status.Flags.PendingApproval() {
			// If the status is pending approval we don't really
			// care about it, we'll get it from somewhere else at
			// some point perhaps. We only care about relayed posts
			// that are already approved or don't require approval.
			log.Debugf(ctx, "status pending approval: %s", status.URI)
			return false, nil
		}

		// For relay purposes we're only interested
		// in statuses that are public or unlisted.
		// Ignore followers-only + direct.
		if vis := status.Visibility; //
		!(vis == gtsmodel.VisibilityPublic ||
			vis == gtsmodel.VisibilityUnlocked) {
			log.Debugf(ctx, "status not public/unlisted: %s", status.URI)
			return false, nil
		}

		// For match checks we have to
		// know the author of the parent
		// status (if there is one).
		var inReplyToAccountURI string
		if status.InReplyToURI != "" {
			if status.InReplyTo != nil {

				// Parent has already been dereferenced, use it.
				inReplyToAccountURI = status.InReplyTo.AccountURI

			} else {
				var err error
				var tsport transport.Transport

				// Get transport on behalf of passed relay actor.
				ctx, tsport, err = d.getTransport(ctx, relayActor)
				if err != nil {
					return false, err
				}

				// Parent not yet dereferenced, try to fetch it from remote.
				inReplyToAccountURI, err = d.retrieveInReplyToAccountURI(ctx,
					tsport, status)
				if err != nil {
					err := gtserror.Newf("error retrieving inReplyToAccountURI: %w", err)
					return false, gtserror.SetUnretrievable(err)
				}
			}
		}

		if matches != nil {
			// Check if status matches
			// defined matching for relay.
			if ok, err := matches(ctx,
				relayActor,
				relaying,
				status,
				inReplyToAccountURI,
			); !ok { // includes err != nil
				log.Debugf(ctx, "status does not match: %s", status.URI)
				return false, err
			}
		}

		if relayable != nil {
			// Check if status is
			// actually relayable.
			if ok, err := relayable(ctx,
				relayActor,
				status,
			); !ok { // includes err != nil
				log.Debugf(ctx, "status not relayable: %s", status.URI)
				return false, err
			}
		}

		return true, nil
	}

	if status == nil {
		// Ensure not a failed search for a local
		// status, if so we know it doesn't exist.
		if uri.Host == config.GetHost() ||
			uri.Host == config.GetAccountDomain() {
			return nil, gtserror.SetUnretrievable(err)
		}

		// Create a bare-bones model for deref.
		status = &gtsmodel.Status{URI: uriStr}

	} else if statusFresh(status, &DefaultStatusFreshness) {

		// This is an existing status that is up-to-date,
		// before returning ensure it is fully populated.
		if err := d.state.DB.PopulateStatus(ctx, status); //
		err != nil {
			log.Errorf(ctx, "error populating existing status: %v", err)
		}

		// Before return, check if existing status
		// is actually relevant for relay usage.
		if ok, err := relevant(ctx, status); !ok {
			return nil, err // includes err != nil
		}

		return status, nil
	}

	// Attempt to enrich and store the existing / newly prepared status.
	status, statusable, isNew, err := d.enrichAndStoreStatusSafely(ctx,
		relayActor,
		uri,
		status,
		nil, // statusable
		relevant,
	)

	if err != nil {
		log.Errorf(ctx, "error getting status %s: %v", uri, err)

		if status != nil {
			err = nil

			// Before return, check if existing status
			// is actually relevant for relay usage.
			if ok, err := relevant(ctx, status); !ok {
				return nil, err // includes err != nil
			}
		}

	} else if statusable != nil {

		// Deref parents + children.
		d.dereferenceThread(ctx,
			relayActor,
			uri,
			status,
			statusable,
			isNew,
		)
	}

	return status, err
}

// retrieveInReplyToAccountURI dereferences the parent
// of the given status to retrieve the URI of the account.
func (d *Dereferencer) retrieveInReplyToAccountURI(
	ctx context.Context,
	tsport transport.Transport,
	status *gtsmodel.Status,
) (
	string,
	error,
) {
	if status.InReplyTo != nil {
		// Had parent stored already.
		return status.InReplyTo.AccountURI, nil
	}

	// We'll have to dereference.
	inReplyToURI := status.InReplyToURI
	parentURI, err := url.Parse(inReplyToURI)
	if err != nil {
		return "", gtserror.Newf("error parsing parent URI: %w", err)
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
		return "", gtserror.Newf("db error checking blocked domain: %w", err)
	}

	if blocked {
		// Just return without dereffing, kiss my ass.
		log.Debugf(ctx, "blocked domain: %s", parentURI.Host)
		return "", nil
	}

	// We don't have the parent stored, try to fetch but *don't* store it, only check for now.
	parentStatusable, parentStatus, _, err := d.retrieveStatusable(ctx, tsport, parentURI)
	if err != nil {
		return "", gtserror.Newf("error retrieving %s: %w", inReplyToURI, err)
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
