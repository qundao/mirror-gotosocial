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
	"net/url"
	"slices"
	"time"

	"code.superseriousbusiness.org/activity/streams/vocab"
	"code.superseriousbusiness.org/gopkg/log"
	"code.superseriousbusiness.org/gotosocial/internal/ap"
	"code.superseriousbusiness.org/gotosocial/internal/config"
	"code.superseriousbusiness.org/gotosocial/internal/gtserror"
	"code.superseriousbusiness.org/gotosocial/internal/gtsmodel"
	"code.superseriousbusiness.org/gotosocial/internal/messages"
)

func (f *DB) Announce(ctx context.Context, announce vocab.ActivityStreamsAnnounce) error {
	log.DebugKV(ctx, "announce", Serialize{announce})

	activityContext := getActivityContext(ctx)
	if activityContext.internal {
		return nil // Already processed.
	}

	requesting := activityContext.requesting
	receiving := activityContext.receiving

	if requesting.IsMoving() {
		// A Moving account
		// can't do this.
		return nil
	}

	// Ensure requestingAccount is among
	// the Actors doing the Announce.
	//
	// We don't support Announce forwards.
	actorIRIs := ap.GetActorIRIs(announce)
	if !slices.ContainsFunc(actorIRIs, func(actorIRI *url.URL) bool {
		return actorIRI.String() == requesting.URI
	}) {
		// Just return nil (status 202) here and
		// not error, as it's not really an error
		// per se, just something we don't support.
		log.Debugf(ctx,
			"requestingAccount %s was not among Announce Actors, dropping Announce forward",
			requesting.URI,
		)
		return nil
	}

	// Via our instance account inbox we
	// only accept Announces delivered by
	// a relay actor we subscribe to.
	if receiving.IsInstance() {
		// Check if the announce originates from an actor
		// we target with at least one relay subscription.
		relaySubscribed, err := f.relaySubscribedToActor(ctx, requesting.URI)
		if err != nil {
			return err
		}

		if !relaySubscribed {
			log.Debugf(ctx, "dropping Announce delivered to our inbox account from Actor we don't subscribe to")
			return nil
		}

		// We have a subscription to the
		// requester, assume the Announce
		// was delivered by a relay we follow.
		return f.announcedFromRelaySubscription(ctx,
			requesting,
			receiving,
			announce,
		)
	}

	if receiving.IsRelayActor() {
		// If the receiver is a relay actor,
		// assume the Announce was delivered
		// by an actor connected to the relay.
		return f.announcedToRelayActor(ctx,
			requesting,
			receiving,
			announce,
		)
	}

	// Convert boost to gtsmodel,
	// checking if it's new to us.
	var isNew bool
	boost, isNew, err := f.converter.ASAnnounceToStatus(ctx, announce)
	if err != nil {
		return gtserror.Newf("error converting announce to boost: %w", err)
	}

	if !isNew {
		// We've already seen
		// and stored this boost;
		// nothing else to do here.
		return nil
	}

	// This is a new or relayed boost.
	// Process side effects asynchronously.
	f.state.Workers.Federator.Queue.Push(&messages.FromFediAPI{
		APObjectType:   ap.ActivityAnnounce,
		APActivityType: ap.ActivityCreate,
		GTSModel:       boost,
		Requesting:     requesting,
		Receiving:      receiving,
	})

	return nil
}

func (f *DB) announcedFromRelaySubscription(
	ctx context.Context,
	requesting *gtsmodel.Account,
	receiving *gtsmodel.Account,
	announce ap.Announceable,
) error {
	// Some relay software doesn't set published
	// prop on the Announce. If this is so, just set
	// time.Now() and let the typeconverter use that.
	published := ap.GetPublished(announce)
	if published.IsZero() {
		ap.SetPublished(announce, time.Now())
	}

	// Convert boost to gtsmodel.
	//
	// We never store boost wrappers from
	// relays so don't bother with isNew here.
	boost, _, err := f.converter.ASAnnounceToStatus(ctx, announce)
	if err != nil {
		return gtserror.Newf("error converting announce to boost: %w", err)
	}

	// From relay actors we don't care
	// about Announces of our own posts.
	uri := boost.BoostOfURI
	if uri.Host == config.GetHost() ||
		uri.Host == config.GetAccountDomain() {
		log.Debugf(ctx, "dropping delivery from %s (relay actor announcing one of our posts)", requesting.URI)
		return nil
	}

	// Ensure we actually follow this
	// relay actor with the instance account.
	following, err := f.state.DB.IsFollowing(ctx, receiving.ID, requesting.ID)
	if err != nil {
		return gtserror.Newf("db error checking follow of actor URI %s: %w", requesting.URI, err)
	}
	if !following {
		// No follow means we're not interested.
		log.Debugf(ctx, "dropping delivery from %s (not following this actor)", requesting.URI)
		return nil
	}

	// Process side effects asynchronously.
	f.state.Workers.Federator.Queue.Push(&messages.FromFediAPI{
		APObjectType:   ap.ActivityAnnounce,
		APActivityType: ap.ActivityCreate,
		GTSModel:       boost,
		Requesting:     requesting,
		Receiving:      receiving,
	})

	return nil
}

func (f *DB) announcedToRelayActor(
	ctx context.Context,
	requesting *gtsmodel.Account,
	receiving *gtsmodel.Account,
	announce ap.Announceable,
) error {
	// Via a relay actor account's inbox we should
	// only process an Announce if it originates from
	// a "connected" account, ie., one that follows us.
	following, err := f.state.DB.IsFollowing(ctx, requesting.ID, receiving.ID)
	if err != nil {
		return gtserror.Newf("db error checking follow of actor URI %s: %w", requesting.URI, err)
	}
	if !following {
		// No follow means we're not interested.
		log.Debugf(ctx, "dropping delivery from %s (our relay doesn't follow them)", requesting.URI)
		return nil
	}

	// Convert boost to gtsmodel.
	//
	// Ignore `isNew` check here: Announce might have been
	// sent to an actual user on our instance as well, but
	// when processing via a relay actor inbox we have special
	// side effects we want to make sure get triggered.
	boost, _, err := f.converter.ASAnnounceToStatus(ctx, announce)
	if err != nil {
		return gtserror.Newf("error converting announce to boost: %w", err)
	}

	// In a relay actor's inbox we don't do
	// anything with Announces of our own posts.
	uri := boost.BoostOfURI
	if uri.Host == config.GetHost() ||
		uri.Host == config.GetAccountDomain() {
		log.Debugf(ctx, "dropping delivery from %s (announcing one of our own posts)", requesting.URI)
		return nil
	}

	// Process side effects asynchronously.
	f.state.Workers.Federator.Queue.Push(&messages.FromFediAPI{
		APObjectType:   ap.ActivityAnnounce,
		APActivityType: ap.ActivityCreate,
		GTSModel:       boost,
		Requesting:     requesting,
		Receiving:      receiving,
	})

	return nil
}
