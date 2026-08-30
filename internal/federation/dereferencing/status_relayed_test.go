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

package dereferencing_test

import (
	"context"
	"testing"

	"code.superseriousbusiness.org/gotosocial/internal/db"
	"code.superseriousbusiness.org/gotosocial/internal/gtserror"
	"code.superseriousbusiness.org/gotosocial/internal/gtsmodel"
	"code.superseriousbusiness.org/gotosocial/testrig"
	"github.com/stretchr/testify/suite"
)

type StatusRelayedTestSuite struct {
	DereferencerTestSuite
}

const irrelevantErrStr = "getOrDerefRelayableStatus: status https://unknown-instance.com/users/brand_new_person/statuses/01FE4NTHKWW7THT67EF10EB839 not matched"

func (suite *StatusRelayedTestSuite) TestGetStatusFromRelaySubscription() {
	// Set up our test structs + tear down on finish.
	testStructs := testrig.SetupTestStructs(rMediaPath, rTemplatePath)
	defer testrig.TearDownTestStructs(testStructs)

	// Clean up test context when done.
	ctx, cncl := context.WithCancel(suite.T().Context())
	defer cncl()

	// Use our instance account to deref.
	instanceAcct := suite.testAccounts["instance_account"]
	// Pretend remote_account_1 is a relay.
	relaySubAcct := suite.testAccounts["remote_account_1"]
	// Pretend status has been relayed to us.
	uri := testrig.URLMustParse("https://unknown-instance.com/users/brand_new_person/statuses/01FE4NTHKWW7THT67EF10EB839")

	// Try to deref the given status URI as though
	// it was forwarded to us by a relay subscription.
	_, err := testStructs.Federator.Dereferencer.GetStatusFromRelaySubscription(ctx,
		instanceAcct, relaySubAcct, uri,
	)

	// Since we don't have a subscription
	// targeting remote_account_1, this will not work.
	suite.True(gtserror.IsNotRelevant(err))
	suite.EqualError(err, irrelevantErrStr)

	// Create a relay subscription targeting remote_account_1.
	// This is obviously a bit silly but it'll do for now.
	relaySub := &gtsmodel.RelaySubscription{
		ID:            "01M030PDD1RS66DFW81KM78MGG",
		AccountID:     suite.testAccounts["admin_account"].ID,
		RelayActorURI: relaySubAcct.URI,
	}
	if err := testStructs.State.DB.PutRelaySubscription(ctx, relaySub); err != nil {
		suite.FailNow(err.Error())
	}

	// Try to deref the given status URI again.
	_, err = testStructs.Federator.Dereferencer.GetStatusFromRelaySubscription(ctx,
		instanceAcct, relaySubAcct, uri,
	)

	// This still won't work
	// because there's no matcher!
	suite.True(gtserror.IsNotRelevant(err))
	suite.EqualError(err, irrelevantErrStr)

	// Update the relay subscription
	// to match everything by default.
	relaySub.Flags.SetMatchByDefault(true)
	if err := testStructs.State.DB.UpdateRelaySubscription(ctx, relaySub); err != nil {
		suite.FailNow(err.Error())
	}

	// Status should actually match now.
	status, err := testStructs.Federator.Dereferencer.GetStatusFromRelaySubscription(ctx,
		instanceAcct, relaySubAcct, uri,
	)
	if err != nil {
		suite.FailNow(err.Error())
	}

	// Status should be in the database.
	dbStatus, err := testStructs.State.DB.GetStatusByURI(ctx, status.URI)
	if err != nil {
		suite.FailNow(err.Error())
	}
	suite.Equal(status.ID, dbStatus.ID)
}

func (suite *StatusRelayedTestSuite) TestGetStatusForRelayActor() {
	// Set up our test structs + tear down on finish.
	testStructs := testrig.SetupTestStructs(rMediaPath, rTemplatePath)
	defer testrig.TearDownTestStructs(testStructs)

	// Clean up test context when done.
	ctx, cncl := context.WithCancel(suite.T().Context())
	defer cncl()

	// Use our relay actor account to deref.
	relayActorAcct := suite.testAccounts["relay_actor_1"]
	// Pretend remote_account_1 is relaying a post to us.
	relayingAcct := suite.testAccounts["remote_account_1"]
	// Pretend status has been relayed to us.
	uri := testrig.URLMustParse("https://unknown-instance.com/users/brand_new_person/statuses/01FE4NTHKWW7THT67EF10EB839")

	// Try to deref the given status URI as
	// though it was forwarded to our relay actor.
	_, _, err := testStructs.Federator.Dereferencer.GetStatusForRelayActor(ctx,
		relayActorAcct, relayingAcct, uri,
	)

	// This won't work since remote_account_1 is trying
	// to relay a status to us from someone else's domain!
	suite.True(gtserror.IsNotRelevant(err))
	suite.EqualError(err, irrelevantErrStr)

	// Relay a status to us from the same domain as the relayer.
	uri = testrig.URLMustParse("http://fossbros-anonymous.io/users/foss_satan/statuses/106221634728637552")
	status, _, err := testStructs.Federator.Dereferencer.GetStatusForRelayActor(ctx,
		relayActorAcct, relayingAcct, uri,
	)
	if err != nil {
		suite.FailNow(err.Error())
	}

	// Status should be defined.
	suite.NotNil(status)

	// But status should NOT be in the database,
	// as we're only relaying it, not storing it.
	_, err = testStructs.State.DB.GetStatusByURI(ctx, status.URI)
	suite.ErrorIs(err, db.ErrNoEntries)
}

func TestStatusRelayedTestSuite(t *testing.T) {
	suite.Run(t, new(StatusRelayedTestSuite))
}
