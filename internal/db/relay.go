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

package db

import (
	"context"

	"code.superseriousbusiness.org/gotosocial/internal/gtsmodel"
)

// Relay contains functions for managing relay
// pushes, subscriptions, and matchers.
type Relay interface {
	// GetRelayPushByID gets one relay push with the given ID.
	GetRelayPushByID(ctx context.Context, id string) (*gtsmodel.RelayPush, error)

	// GetRelayPushesForAccountID gets relay pushes for given accountID.
	GetRelayPushesForAccountID(ctx context.Context, accountID string) ([]*gtsmodel.RelayPush, error)

	// GetRelayPushesByActorURI gets all relay pushes targeting given actor URI.
	GetRelayPushesByActorURI(ctx context.Context, uri string) ([]*gtsmodel.RelayPush, error)

	// PopulateRelayPush populates the given relay push.
	PopulateRelayPush(ctx context.Context, relayPush *gtsmodel.RelayPush) error

	// PutRelayPush inserts the given relay push into the db.
	PutRelayPush(ctx context.Context, relayPush *gtsmodel.RelayPush) error

	// UpdateRelayPush updates the given columns of the given relay push.
	UpdateRelayPush(ctx context.Context, relayPush *gtsmodel.RelayPush, columns ...string) error

	// DeleteRelayPush deletes the given relay push.
	DeleteRelayPush(ctx context.Context, relayPush *gtsmodel.RelayPush) error

	// GetRelayActorByID gets relay actor with the given ID.
	GetRelayActorByID(ctx context.Context, id string) (*gtsmodel.RelayActor, error)

	// GetRelayActorByURI gets relay actor with the given URI.
	GetRelayActorByURI(ctx context.Context, uri string) (*gtsmodel.RelayActor, error)

	// GetRelayActors gets all relay actors.
	GetRelayActors(ctx context.Context) ([]*gtsmodel.RelayActor, error)

	// Populate relay actors populates the given relay actor.
	PopulateRelayActor(ctx context.Context, relayActor *gtsmodel.RelayActor) error

	// PutRelayActor inserts the given relay actor into the db.
	PutRelayActor(ctx context.Context, relayActor *gtsmodel.RelayActor) error

	// UpdateRelayActor updates the given relay actor.
	UpdateRelayActor(ctx context.Context, relayActor *gtsmodel.RelayActor, columns ...string) error

	// DeleteRelayActor deletes the given relay actor.
	DeleteRelayActor(ctx context.Context, relayActor *gtsmodel.RelayActor) error

	// GetRelaySubscriptionByID gets one relay subscription with the given ID.
	GetRelaySubscriptionByID(ctx context.Context, id string) (*gtsmodel.RelaySubscription, error)

	// GetRelaySubscriptionsByActorURI gets all relay subscriptions targeting given actor URI.
	GetRelaySubscriptionsByActorURI(ctx context.Context, uri string) ([]*gtsmodel.RelaySubscription, error)

	// GetRelaySubscriptionsPage fetches all relay subscriptions.
	GetRelaySubscriptions(ctx context.Context) ([]*gtsmodel.RelaySubscription, error)

	// PopulateRelaySubscription populates the given relay subscription.
	PopulateRelaySubscription(ctx context.Context, relaySubscription *gtsmodel.RelaySubscription) error

	// PutRelaySubscription inserts the given relay subscription into the db.
	PutRelaySubscription(ctx context.Context, relaySubscription *gtsmodel.RelaySubscription) error

	// UpdateRelaySubscription updates the given columns of the given relay subscription.
	UpdateRelaySubscription(ctx context.Context, relaySubscription *gtsmodel.RelaySubscription, columns ...string) error

	// DeleteRelaySubscription deletes the given relay subscription.
	DeleteRelaySubscription(ctx context.Context, relaySubscription *gtsmodel.RelaySubscription) error

	// GetRelayMatcher gets one relay matcher with the given ID.
	GetRelayMatcher(ctx context.Context, id string) (*gtsmodel.RelayMatcher, error)

	// PutRelayMatcher inserts the given relay matcher into the db.
	PutRelayMatcher(ctx context.Context, relayMatcher *gtsmodel.RelayMatcher) error

	// UpdateRelayMatcher updates the given columns of the given relay matcher.
	UpdateRelayMatcher(ctx context.Context, relayMatcher *gtsmodel.RelayMatcher, columns ...string) error

	// DeleteRelayMatcher deletes relay matcher with the given ID.
	DeleteRelayMatcher(ctx context.Context, id string) error

	// PutRelayedURI inserts the given relayedURI into the db.
	PutRelayedURI(ctx context.Context, relayedURI *gtsmodel.RelayedURI) error

	// GetRelayedURI gets one relayedURI model by the given relay actor URI and status URI, if it exists.
	GetRelayedURI(ctx context.Context, relayURI string, statusURI string) (*gtsmodel.RelayedURI, error)

	// TODO: add a function to count relayedURIs per actor,
	// so that relay account status "created posts" thingy
	// can be populated, to provide fun stats for admins.
}
