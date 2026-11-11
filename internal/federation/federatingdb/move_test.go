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

type MoveTestSuite struct {
	FederatingDBTestSuite
}

func (suite *MoveTestSuite) move(
	ctx context.Context,
	testStructs *testrig.TestStructs,
	receiver *gtsmodel.Account,
	requester *gtsmodel.Account,
	moveStr string,
) error {
	ctx = createTestContext(ctx, requester, receiver)

	rawMove := make(map[string]interface{})
	if err := json.Unmarshal([]byte(moveStr), &rawMove); err != nil {
		suite.FailNow(err.Error())
	}

	t, err := streams.ToType(ctx, rawMove)
	if err != nil {
		suite.FailNow(err.Error())
	}

	move, ok := t.(vocab.ActivityStreamsMove)
	if !ok {
		suite.FailNow("", "couldn't cast %T to Move", t)
	}

	return testStructs.Federator.FederatingDB().Move(ctx, move)
}

func (suite *MoveTestSuite) TestMove() {
	// Set up our test structs + tear down on finish.
	testStructs := testrig.SetupTestStructs(rMediaPath, rTemplatePath)
	defer testrig.TearDownTestStructs(testStructs)

	// Clean up test context when done.
	ctx, cncl := context.WithCancel(suite.T().Context())
	defer cncl()

	var (
		receivingAcct  = suite.testAccounts["local_account_1"]
		requestingAcct = suite.testAccounts["remote_account_1"]
		moveURI        = "http://fossbros-anonymous.io/users/foss_satan/moves/01HR9FDFCAGM7JYPMWNTFRDQE9"
		moveActor      = requestingAcct.URI
		moveTarget     = "https://turnip.farm/users/turniplover6969"
		moveStr1       = `{
  "@context": "https://www.w3.org/ns/activitystreams",
  "id": "` + moveURI + `",
  "actor": "` + moveActor + `",
  "type": "Move",
  "object": "` + moveActor + `",
  "target": "` + moveTarget + `",
  "to": "http://fossbros-anonymous.io/users/foss_satan/followers"
}`
	)

	// Trigger the move.
	suite.move(ctx,
		testStructs,
		receivingAcct,
		requestingAcct,
		moveStr1,
	)

	// Wait for the move attempt to be put in the db. We only care
	// here about the move passing initial checks and being inserted,
	// not whether it succeeds or not, as that part is tested in
	// internal/processing/workers/fromfediapi_test.go.
	var move *gtsmodel.Move
	if !testrig.WaitFor(func() bool {
		var err error
		move, err = testStructs.State.DB.GetMoveByURI(ctx, moveURI)
		return err == nil && move != nil
	}) {
		suite.FailNow("waiting for move")
	}
	suite.Equal(moveTarget, move.TargetURI)
	suite.Equal(moveActor, move.OriginURI)
}

func (suite *MoveTestSuite) TestBadMoves() {
	// Set up our test structs + tear down on finish.
	testStructs := testrig.SetupTestStructs(rMediaPath, rTemplatePath)
	defer testrig.TearDownTestStructs(testStructs)

	// Clean up test context when done.
	ctx, cncl := context.WithCancel(suite.T().Context())
	defer cncl()

	var (
		receivingAcct  = suite.testAccounts["local_account_1"]
		requestingAcct = suite.testAccounts["remote_account_1"]
	)

	type testStruct struct {
		moveStr string
		err     string
	}

	for _, t := range []testStruct{
		{
			// Move signed by someone else.
			moveStr: `{
  "@context": "https://www.w3.org/ns/activitystreams",
  "id": "http://fossbros-anonymous.io/users/foss_satan/moves/01HR9FDFCAGM7JYPMWNTFRDQE9",
  "actor": "http://fossbros-anonymous.io/users/someone_else",
  "type": "Move",
  "object": "http://fossbros-anonymous.io/users/foss_satan",
  "target": "https://turnip.farm/users/turniplover6969",
  "to": "http://fossbros-anonymous.io/users/foss_satan/followers"
}`,
			err: "Move was signed by http://fossbros-anonymous.io/users/foss_satan but actor was http://fossbros-anonymous.io/users/someone_else",
		},
		{
			// Actor and object not the same.
			moveStr: `{
  "@context": "https://www.w3.org/ns/activitystreams",
  "id": "http://fossbros-anonymous.io/users/foss_satan/moves/01HR9FDFCAGM7JYPMWNTFRDQE9",
  "actor": "http://fossbros-anonymous.io/users/foss_satan",
  "type": "Move",
  "object": "http://fossbros-anonymous.io/users/someone_else",
  "target": "https://turnip.farm/users/turniplover6969",
  "to": "http://fossbros-anonymous.io/users/foss_satan/followers"
}`,
			err: "Move was signed by http://fossbros-anonymous.io/users/foss_satan but object was http://fossbros-anonymous.io/users/someone_else",
		},
		{
			// Object and target the same.
			moveStr: `{
  "@context": "https://www.w3.org/ns/activitystreams",
  "id": "http://fossbros-anonymous.io/users/foss_satan/moves/01HR9FDFCAGM7JYPMWNTFRDQE9",
  "actor": "http://fossbros-anonymous.io/users/foss_satan",
  "type": "Move",
  "object": "http://fossbros-anonymous.io/users/foss_satan",
  "target": "http://fossbros-anonymous.io/users/foss_satan",
  "to": "http://fossbros-anonymous.io/users/foss_satan/followers"
}`,
			err: "Move target and origin were the same (http://fossbros-anonymous.io/users/foss_satan)",
		},
	} {
		// Trigger the move.
		err := suite.move(ctx,
			testStructs,
			receivingAcct,
			requestingAcct,
			t.moveStr,
		)
		if t.err != "" {
			suite.EqualError(err, t.err)
		}
	}
}

func TestMoveTestSuite(t *testing.T) {
	suite.Run(t, &MoveTestSuite{})
}
