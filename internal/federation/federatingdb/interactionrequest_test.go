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
	"encoding/json"
	"testing"

	"code.superseriousbusiness.org/activity/streams"
	"code.superseriousbusiness.org/activity/streams/vocab"
	"code.superseriousbusiness.org/gotosocial/internal/gtsmodel"
	"code.superseriousbusiness.org/gotosocial/testrig"
	"github.com/stretchr/testify/suite"
)

type InteractionRequestTestSuite struct {
	FederatingDBTestSuite
}

func (suite *InteractionRequestTestSuite) intReq(
	ctx context.Context,
	requesting *gtsmodel.Account,
	receiving *gtsmodel.Account,
	jsonStr string,
	dbF func(ctx context.Context, req vocab.Type) error,
) error {
	ctx = createTestContext(ctx, requesting, receiving)

	raw := make(map[string]interface{})
	if err := json.Unmarshal([]byte(jsonStr), &raw); err != nil {
		suite.FailNow(err.Error())
	}

	t, err := streams.ToType(ctx, raw)
	if err != nil {
		suite.FailNow(err.Error())
	}

	return dbF(ctx, t)
}

func (suite *InteractionRequestTestSuite) TestReplyRequest() {
	// Set up our test structs + tear down on finish.
	testStructs := testrig.SetupTestStructs(rMediaPath, rTemplatePath)
	defer testrig.TearDownTestStructs(testStructs)

	// Clean up test context when done.
	ctx, cncl := context.WithCancel(suite.T().Context())
	defer cncl()

	var (
		receiving  = suite.testAccounts["admin_account"]
		requesting = suite.testAccounts["remote_account_1"]
		testStatus = suite.testStatuses["admin_account_status_1"]
		intReqURI  = "http://fossbros-anonymous.io/requests/87fb1478-ac46-406a-8463-96ce05645219"
		intURI     = "http://fossbros-anonymous.io/users/foss_satan/statuses/87fb1478-ac46-406a-8463-96ce05645219"
		jsonStr    = `{
  "@context": [
    "https://www.w3.org/ns/activitystreams",
    "https://gotosocial.org/ns"
  ],
  "type": "ReplyRequest",
  "id": "` + intReqURI + `",
  "actor": "` + requesting.URI + `",
  "object": "` + testStatus.URI + `",
  "to": "` + receiving.URI + `",
  "instrument": {
    "attributedTo": "` + requesting.URI + `",
    "cc": "` + requesting.FollowersURI + `",
    "content": "\u003cp\u003ethis is a reply!\u003c/p\u003e",
    "id": "` + intURI + `",
    "inReplyTo": "` + testStatus.URI + `",
	"tag": {
      "href": "` + receiving.URI + `",
      "name": "@` + receiving.Username + `@localhost:8080",
      "type": "Mention"
    },
    "to": "https://www.w3.org/ns/activitystreams#Public",
    "type": "Note"
  }
}`
	)

	suite.T().Logf("testing reply request:\n\n%s", jsonStr)

	// Call the federatingDB function.
	err := suite.intReq(
		ctx,
		requesting,
		receiving,
		jsonStr,
		func(ctx context.Context, req vocab.Type) error {
			replyReq := req.(vocab.GoToSocialReplyRequest)
			return testStructs.Federator.FederatingDB().ReplyRequest(ctx, replyReq)
		},
	)
	if err != nil {
		suite.FailNow(err.Error())
	}

	// There should be an interaction request in the DB now.
	var intReq *gtsmodel.InteractionRequest
	if !testrig.WaitFor(func() bool {
		intReq, err = testStructs.State.DB.GetInteractionRequestByInteractionURI(ctx, intURI)
		return err == nil && intReq != nil
	}) {
		suite.FailNow("timed out waiting for int req to appear in the db")
	}
	suite.Equal(testStatus.ID, intReq.TargetStatusID)
	suite.Equal(receiving.ID, intReq.TargetAccountID)
	suite.Equal(requesting.ID, intReq.InteractingAccountID)
	suite.Equal(intReqURI, intReq.InteractionRequestURI)
	suite.Equal(intURI, intReq.InteractionURI)
	suite.Equal(gtsmodel.InteractionReply, intReq.InteractionType)

	// Wait for approval or perish.
	if !testrig.WaitFor(func() bool {
		reply, err := testStructs.State.DB.GetStatusByURI(ctx, intURI)
		return err == nil && reply != nil && reply.ApprovedByURI != ""
	}) {
		suite.FailNow("waiting for approval")
	}
}

func (suite *InteractionRequestTestSuite) TestLikeRequest() {
	// Set up our test structs + tear down on finish.
	testStructs := testrig.SetupTestStructs(rMediaPath, rTemplatePath)
	defer testrig.TearDownTestStructs(testStructs)

	// Clean up test context when done.
	ctx, cncl := context.WithCancel(suite.T().Context())
	defer cncl()

	var (
		receiving  = suite.testAccounts["admin_account"]
		requesting = suite.testAccounts["remote_account_1"]
		testStatus = suite.testStatuses["admin_account_status_1"]
		intReqURI  = "http://fossbros-anonymous.io/requests/87fb1478-ac46-406a-8463-96ce05645219"
		intURI     = "http://fossbros-anonymous.io/users/foss_satan/statuses/87fb1478-ac46-406a-8463-96ce05645219"
		jsonStr    = `{
  "@context": [
    "https://www.w3.org/ns/activitystreams",
    "https://gotosocial.org/ns"
  ],
  "type": "LikeRequest",
  "id": "` + intReqURI + `",
  "actor": "` + requesting.URI + `",
  "object": "` + testStatus.URI + `",
  "to": "` + receiving.URI + `",
  "instrument": {
    "id": "` + intURI + `",
    "object": "` + testStatus.URI + `",
    "actor": "` + requesting.URI + `",
    "to": "` + receiving.URI + `",
    "type": "Like"
  }
}`
	)

	suite.T().Logf("testing like request:\n\n%s", jsonStr)

	// Call the federatingDB function.
	err := suite.intReq(
		ctx,
		requesting,
		receiving,
		jsonStr,
		func(ctx context.Context, req vocab.Type) error {
			likeReq := req.(vocab.GoToSocialLikeRequest)
			return testStructs.Federator.FederatingDB().LikeRequest(ctx, likeReq)
		},
	)
	if err != nil {
		suite.FailNow(err.Error())
	}

	// There should be an interaction request in the DB now.
	var intReq *gtsmodel.InteractionRequest
	if !testrig.WaitFor(func() bool {
		intReq, err = testStructs.State.DB.GetInteractionRequestByInteractionURI(ctx, intURI)
		return err == nil && intReq != nil
	}) {
		suite.FailNow("timed out waiting for int req to appear in the db")
	}
	suite.Equal(testStatus.ID, intReq.TargetStatusID)
	suite.Equal(receiving.ID, intReq.TargetAccountID)
	suite.Equal(requesting.ID, intReq.InteractingAccountID)
	suite.Equal(intReqURI, intReq.InteractionRequestURI)
	suite.Equal(intURI, intReq.InteractionURI)
	suite.Equal(gtsmodel.InteractionLike, intReq.InteractionType)

	// Wait for approval or perish.
	if !testrig.WaitFor(func() bool {
		fave, err := testStructs.State.DB.GetStatusFaveByURI(ctx, intURI)
		return err == nil && fave != nil && fave.ApprovedByURI != ""
	}) {
		suite.FailNow("waiting for approval")
	}
}

func (suite *InteractionRequestTestSuite) TestAnnounceRequest() {
	// Set up our test structs + tear down on finish.
	testStructs := testrig.SetupTestStructs(rMediaPath, rTemplatePath)
	defer testrig.TearDownTestStructs(testStructs)

	// Clean up test context when done.
	ctx, cncl := context.WithCancel(suite.T().Context())
	defer cncl()

	var (
		receiving  = suite.testAccounts["admin_account"]
		requesting = suite.testAccounts["remote_account_1"]
		testStatus = suite.testStatuses["admin_account_status_1"]
		intReqURI  = "http://fossbros-anonymous.io/requests/87fb1478-ac46-406a-8463-96ce05645219"
		intURI     = "http://fossbros-anonymous.io/users/foss_satan/statuses/87fb1478-ac46-406a-8463-96ce05645219"
		jsonStr    = `{
  "@context": [
    "https://www.w3.org/ns/activitystreams",
    "https://gotosocial.org/ns"
  ],
  "type": "AnnounceRequest",
  "id": "` + intReqURI + `",
  "actor": "` + requesting.URI + `",
  "object": "` + testStatus.URI + `",
  "to": "` + receiving.URI + `",
  "instrument": {
    "id": "` + intURI + `",
    "object": "` + testStatus.URI + `",
    "actor": "` + requesting.URI + `",
    "to": "` + requesting.FollowersURI + `",
    "cc": "` + receiving.URI + `",
    "type": "Announce"
  }
}`
	)

	suite.T().Logf("testing announce request:\n\n%s", jsonStr)

	// Call the federatingDB function.
	err := suite.intReq(
		ctx,
		requesting,
		receiving,
		jsonStr,
		func(ctx context.Context, req vocab.Type) error {
			announceReq := req.(vocab.GoToSocialAnnounceRequest)
			return testStructs.Federator.FederatingDB().AnnounceRequest(ctx, announceReq)
		},
	)
	if err != nil {
		suite.FailNow(err.Error())
	}

	// There should be an interaction request in the DB now.
	var intReq *gtsmodel.InteractionRequest
	if !testrig.WaitFor(func() bool {
		intReq, err = testStructs.State.DB.GetInteractionRequestByInteractionURI(ctx, intURI)
		return err == nil && intReq != nil
	}) {
		suite.FailNow("timed out waiting for int req to appear in the db")
	}
	suite.Equal(testStatus.ID, intReq.TargetStatusID)
	suite.Equal(receiving.ID, intReq.TargetAccountID)
	suite.Equal(requesting.ID, intReq.InteractingAccountID)
	suite.Equal(intReqURI, intReq.InteractionRequestURI)
	suite.Equal(intURI, intReq.InteractionURI)
	suite.Equal(gtsmodel.InteractionAnnounce, intReq.InteractionType)

	// Wait for approval or perish.
	if !testrig.WaitFor(func() bool {
		boost, err := testStructs.State.DB.GetStatusByURI(ctx, intURI)
		return err == nil && boost != nil && boost.ApprovedByURI != ""
	}) {
		suite.FailNow("waiting for approval")
	}
}

func TestInteractionRequestTestSuite(t *testing.T) {
	suite.Run(t, &InteractionRequestTestSuite{})
}
