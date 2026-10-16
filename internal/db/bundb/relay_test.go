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

package bundb_test

import (
	"testing"

	"code.superseriousbusiness.org/gotosocial/internal/gtsmodel"
	"github.com/stretchr/testify/suite"
)

type RelayTestSuite struct {
	BunDBStandardTestSuite
}

func (suite *RelayTestSuite) TestCreateGetDeleteRelayPush() {
	var (
		ctx           = suite.T().Context()
		testAccountID = suite.testAccounts["local_account_1"].ID
		relayPushID   = "01KMMVSB3CBX7HBG7RA52ASMMC"
		relayActorURI = "https://relay.example.org/actor"
		relayPush     = &gtsmodel.RelayPush{
			ID:            relayPushID,
			AccountID:     testAccountID,
			RelayActorURI: relayActorURI,
		}
	)

	if err := suite.state.DB.PutRelayPush(ctx, relayPush); err != nil {
		suite.FailNow(err.Error())
	}

	// Ensure default flag values set.
	suite.True(relayPush.Flags.Public())
	suite.False(relayPush.Flags.Unlisted())
	suite.False(relayPush.Flags.IgnoreSensitive())
	suite.False(relayPush.Flags.IgnoreMedia())
	suite.False(relayPush.Flags.IgnoreReplies())

	// Get relay pushes by account ID.
	relayPushes, err := suite.state.DB.GetRelayPushesForAccountID(ctx, testAccountID)
	if err != nil {
		suite.FailNow(err.Error())
	}
	suite.Len(relayPushes, 1)

	// Get relay pushes by actor URI.
	relayPushes, err = suite.state.DB.GetRelayPushesByActorURI(ctx, relayActorURI)
	if err != nil {
		suite.FailNow(err.Error())
	}
	suite.Len(relayPushes, 1)

	// Create some matchers for this relay push.
	relayMatchers := []*gtsmodel.RelayMatcher{
		{
			ID:      "01KMMZ3G858KVAD3JN65KEJN5V",
			RelayID: relayPushID,
			Keyword: "#HEALTH",
		},
		{
			ID:      "01KMMZ2ABEX3SDARVFXND6Z375",
			RelayID: relayPushID,
			Keyword: "#Boobs",
		},
		{
			ID:      "01KMMZ148TZE5XNEB3GVPGEASF",
			RelayID: relayPushID,
			Keyword: "#MoreBoobs",
		},
	}
	for _, matcher := range relayMatchers {
		if err := suite.state.DB.PutRelayMatcher(ctx, matcher); err != nil {
			suite.FailNow(err.Error())
		}
		relayPush.MatcherIDs = append(relayPush.MatcherIDs, matcher.ID)
	}

	// Update relay push with matcher IDs.
	if err := suite.state.DB.UpdateRelayPush(ctx, relayPush, "matchers"); err != nil {
		suite.FailNow(err.Error())
	}

	// Get push from DB and ensure matchers populated.
	relayPush, err = suite.state.DB.GetRelayPushByID(ctx, relayPushID)
	if err != nil {
		suite.FailNow(err.Error())
	}
	suite.Len(relayPush.Matchers, 3)
	for _, matcher := range relayPush.Matchers {
		if matcher.Regexp == nil {
			suite.FailNow("expected not nil matcher regexp")
		}
	}

	// Delete push.
	if err := suite.state.DB.DeleteRelayPush(ctx, relayPush); err != nil {
		suite.FailNow(err.Error())
	}
}

func (suite *RelayTestSuite) TestCreateGetDeleteRelaySubscription() {
	var (
		ctx                 = suite.T().Context()
		testAccountID       = suite.testAccounts["local_account_1"].ID
		relaySubscriptionID = "01KMMVSB3CBX7HBG7RA52ASMMC"
		relayActorURI       = "https://relay.example.org/actor"
		relaySubscription   = &gtsmodel.RelaySubscription{
			ID:            relaySubscriptionID,
			AccountID:     testAccountID,
			RelayActorURI: relayActorURI,
		}
	)

	if err := suite.state.DB.PutRelaySubscription(ctx, relaySubscription); err != nil {
		suite.FailNow(err.Error())
	}

	// Ensure default flag values set.
	suite.True(relaySubscription.Flags.Public())
	suite.False(relaySubscription.Flags.Unlisted())
	suite.False(relaySubscription.Flags.IgnoreSensitive())
	suite.False(relaySubscription.Flags.IgnoreMedia())
	suite.False(relaySubscription.Flags.IgnoreReplies())

	// Get relay subscriptions by actor URI.
	relaySubscriptions, err := suite.state.DB.GetRelaySubscriptionsByActorURI(ctx, relayActorURI)
	if err != nil {
		suite.FailNow(err.Error())
	}
	suite.Len(relaySubscriptions, 1)

	// Create some matchers for this relay subscription.
	relayMatchers := []*gtsmodel.RelayMatcher{
		{
			ID:      "01KMMZ3G858KVAD3JN65KEJN5V",
			RelayID: relaySubscriptionID,
			Keyword: "#HEALTH",
		},
		{
			ID:      "01KMMZ2ABEX3SDARVFXND6Z375",
			RelayID: relaySubscriptionID,
			Keyword: "#Boobs",
		},
		{
			ID:      "01KMMZ148TZE5XNEB3GVPGEASF",
			RelayID: relaySubscriptionID,
			Keyword: "#MoreBoobs",
		},
	}
	for _, matcher := range relayMatchers {
		if err := suite.state.DB.PutRelayMatcher(ctx, matcher); err != nil {
			suite.FailNow(err.Error())
		}
		relaySubscription.MatcherIDs = append(relaySubscription.MatcherIDs, matcher.ID)
	}

	// Update relay subscription with matcher IDs.
	if err := suite.state.DB.UpdateRelaySubscription(ctx, relaySubscription, "matchers"); err != nil {
		suite.FailNow(err.Error())
	}

	// Get subscription from DB and ensure matchers populated.
	relaySubscription, err = suite.state.DB.GetRelaySubscriptionByID(ctx, relaySubscriptionID)
	if err != nil {
		suite.FailNow(err.Error())
	}
	suite.Len(relaySubscription.Matchers, 3)
	for _, matcher := range relaySubscription.Matchers {
		if matcher.Regexp == nil {
			suite.FailNow("expected not nil matcher regexp")
		}
	}

	// Delete subscription.
	if err := suite.state.DB.DeleteRelaySubscription(ctx, relaySubscription); err != nil {
		suite.FailNow(err.Error())
	}
}

func (suite *RelayTestSuite) TestPutRelayedURI() {
	ctx := suite.T().Context()
	relayedURI := &gtsmodel.RelayedURI{
		ID:        "01M52CWKYJ3Z1T8V75T0ZDCWJ1",
		RelayURI:  "http://localhost:8080/relays/some_relay_actor",
		StatusURI: "https://example.org/statuses/some_status",
	}

	// Insert relayedURI.
	inserted, err := suite.state.DB.PutRelayedURI(ctx, relayedURI)
	if err != nil {
		suite.FailNow(err.Error())
	}
	if !inserted {
		suite.FailNow("should have been inserted")
	}

	// Try to insert relayedURI again.
	inserted, err = suite.state.DB.PutRelayedURI(ctx, relayedURI)
	if err != nil {
		suite.FailNow(err.Error())
	}
	if inserted {
		suite.FailNow("should not have been inserted")
	}
}

func TestRelayTestSuite(t *testing.T) {
	suite.Run(t, new(RelayTestSuite))
}
