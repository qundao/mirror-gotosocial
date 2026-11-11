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

package federatingdb_test

import (
	"context"
	"testing"

	"code.superseriousbusiness.org/activity/streams"
	"code.superseriousbusiness.org/activity/streams/vocab"
	"code.superseriousbusiness.org/gotosocial/internal/ap"
	"code.superseriousbusiness.org/gotosocial/internal/gtscontext"
	"code.superseriousbusiness.org/gotosocial/internal/gtsmodel"
	"code.superseriousbusiness.org/gotosocial/internal/transport/delivery"
	"code.superseriousbusiness.org/gotosocial/testrig"
	"github.com/stretchr/testify/suite"
)

type FollowTestSuite struct {
	FederatingDBTestSuite
}

// Mastodon and *key softwares (and probably some others)
// follow relay actors using the ActivityPub Public URI
// rather than the AP ID/URI of the actor. This test checks
// that our code properly handles follows like this.
func (suite *FollowTestSuite) TestFollowRelayActorInboxURI() {
	// Set up our test structs + tear down on finish.
	testStructs := testrig.SetupTestStructs(rMediaPath, rTemplatePath)
	defer testrig.TearDownTestStructs(testStructs)

	// Clean up test context when done.
	ctx, cncl := context.WithCancel(suite.T().Context())
	defer cncl()

	const followURI = "https://example.org/some_follow_uri"
	followingAcct := suite.testAccounts["remote_account_2"]
	targetAcct := suite.testAccounts["relay_actor_1"]

	followable, err := testStructs.TypeConverter.FollowToAS(ctx, &gtsmodel.Follow{
		URI:             followURI,
		AccountID:       followingAcct.ID,
		TargetAccountID: targetAcct.ID,
	})
	if err != nil {
		suite.FailNow(err.Error())
	}

	// Rewrite the object URI of
	// the Follow to the Public URI.
	followable.SetActivityStreamsObject(
		func() vocab.ActivityStreamsObjectProperty {
			objectProp := streams.NewActivityStreamsObjectProperty()
			objectProp.AppendIRI(ap.PublicIRI())
			return objectProp
		}(),
	)

	// Trigger federatingDB Follow func.
	ctx = gtscontext.SetRequestingAccount(ctx, followingAcct)
	ctx = gtscontext.SetReceivingAccount(ctx, targetAcct)
	if err := testStructs.Federator.FederatingDB().Follow(ctx, followable); err != nil {
		suite.FailNow(err.Error())
	}

	// Wait for follow in the db.
	var follow *gtsmodel.Follow
	if !testrig.WaitFor(func() bool {
		follow, err = testStructs.State.DB.GetFollowByURI(ctx, followURI)
		return err == nil && follow != nil
	}) {
		suite.FailNow("timed out waiting for follow")
	}

	// Should be true to indicate
	// a redirect from the public URI.
	suite.True(follow.Flags.UsePublicURI())

	// Wait for delivery side effect.
	var delivery *delivery.Delivery
	if !testrig.WaitFor(func() bool {
		var ok bool
		delivery, ok = testStructs.State.Workers.Delivery.Queue.Pop()
		return ok && delivery != nil
	}) {
		suite.FailNow("timed out waiting for delivery")
	}

	// Decode delivery message.
	t, err := ap.DecodeType(ctx, delivery.Request.Body)
	if err != nil {
		suite.FailNow(err.Error())
	}

	// Should be an Accept.
	accept, ok := t.(vocab.ActivityStreamsAccept)
	if !ok {
		suite.FailNow("%T not Accept", t)
	}

	// Should have one object of the Accept.
	objects := ap.ExtractObjects(accept)
	if len(objects) != 1 {
		suite.FailNow("wrong objects length")
	}

	// Accept object should be the
	// recreated original Follow.
	followable = objects[0].GetType().(vocab.ActivityStreamsFollow)

	// Accepted Follow object IRI should
	// be the ActivityPub Public URI.
	followObjectIRI := ap.GetObjectIRIs(followable)[0]
	suite.Equal(ap.PublicIRI(), followObjectIRI)
}

func TestFollowTestSuite(t *testing.T) {
	suite.Run(t, new(FollowTestSuite))
}
