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

	"code.superseriousbusiness.org/gotosocial/internal/ap"
	"code.superseriousbusiness.org/gotosocial/internal/gtsmodel"
	"code.superseriousbusiness.org/gotosocial/testrig"
	"github.com/stretchr/testify/suite"
)

type CreateTestSuite struct {
	FederatingDBTestSuite
}

func (suite *CreateTestSuite) TestCreateNote() {
	// Set up our test structs + tear down on finish.
	testStructs := testrig.SetupTestStructs(rMediaPath, rTemplatePath)
	defer testrig.TearDownTestStructs(testStructs)

	// Clean up test context when done.
	ctx, cncl := context.WithCancel(suite.T().Context())
	defer cncl()

	receivingAccount := suite.testAccounts["local_account_1"]
	requestingAccount := suite.testAccounts["remote_account_1"]

	// Call the Create function with
	// a Create activity for a note.
	create := suite.testActivities["dm_for_zork"].Activity
	objProp := create.GetActivityStreamsObject()
	note := objProp.At(0).GetType().(ap.Statusable)
	noteURI := ap.GetJSONLDId(note)
	err := testStructs.Federator.FederatingDB().Create(
		createTestContext(
			ctx,
			requestingAccount,
			receivingAccount,
		),
		create,
	)
	if err != nil {
		suite.FailNow(err.Error())
	}

	// Wait for the status to
	// be created in the database.
	var status *gtsmodel.Status
	if !testrig.WaitFor(func() bool {
		var err error
		status, err = testStructs.State.DB.GetStatusByURI(ctx, noteURI.String())
		return err == nil && status != nil
	}) {
		suite.FailNow("timed out waiting for status")
	}
}

func (suite *CreateTestSuite) TestCreateBareNote() {
	// Set up our test structs + tear down on finish.
	testStructs := testrig.SetupTestStructs(rMediaPath, rTemplatePath)
	defer testrig.TearDownTestStructs(testStructs)

	// Clean up test context when done.
	ctx, cncl := context.WithCancel(suite.T().Context())
	defer cncl()

	receivingAccount := suite.testAccounts["local_account_1"]
	requestingAccount := suite.testAccounts["remote_account_1"]

	// Call the Create function with
	// the bare Note object as param.
	create := suite.testActivities["dm_for_zork"].Activity
	objProp := create.GetActivityStreamsObject()
	note := objProp.At(0).GetType().(ap.Statusable)
	noteURI := ap.GetJSONLDId(note)
	err := testStructs.Federator.FederatingDB().Create(
		createTestContext(
			ctx,
			requestingAccount,
			receivingAccount,
		),
		note,
	)
	if err != nil {
		suite.FailNow(err.Error())
	}

	// Wait for the status to
	// be created in the database.
	var status *gtsmodel.Status
	if !testrig.WaitFor(func() bool {
		var err error
		status, err = testStructs.State.DB.GetStatusByURI(ctx, noteURI.String())
		return err == nil && status != nil
	}) {
		suite.FailNow("timed out waiting for status")
	}
}

func (suite *CreateTestSuite) TestCreateNoteForward() {
	// Set up our test structs + tear down on finish.
	testStructs := testrig.SetupTestStructs(rMediaPath, rTemplatePath)
	defer testrig.TearDownTestStructs(testStructs)

	// Clean up test context when done.
	ctx, cncl := context.WithCancel(suite.T().Context())
	defer cncl()

	receivingAccount := suite.testAccounts["local_account_1"]
	requestingAccount := suite.testAccounts["remote_account_1"]

	// Ensure a follow exists between requesting
	// and receiving account, this ensures the forward
	// will be seen as "relevant" and not get dropped.
	err := testStructs.State.DB.PutFollow(ctx, &gtsmodel.Follow{
		ID:              "01M2588ZDD5NV2G3GVR811MREB",
		URI:             "https://this.is.a.url",
		AccountID:       receivingAccount.ID,
		TargetAccountID: requestingAccount.ID,
	})
	if err != nil {
		suite.FailNow(err.Error())
	}

	// Call the Create function with
	// Create of the forwarded Note.
	create := suite.testActivities["forwarded_message"].Activity
	objProp := create.GetActivityStreamsObject()
	note := objProp.At(0).GetType().(ap.Statusable)
	noteURI := ap.GetJSONLDId(note)
	err = testStructs.Federator.FederatingDB().Create(
		createTestContext(
			ctx,
			requestingAccount,
			receivingAccount,
		),
		create,
	)
	if err != nil {
		suite.FailNow(err.Error())
	}

	// Wait for the status to
	// be created in the database.
	var status *gtsmodel.Status
	if !testrig.WaitFor(func() bool {
		var err error
		status, err = testStructs.State.DB.GetStatusByURI(ctx, noteURI.String())
		return err == nil && status != nil
	}) {
		suite.FailNow("timed out waiting for status")
	}
}

func TestCreateTestSuite(t *testing.T) {
	suite.Run(t, &CreateTestSuite{})
}
