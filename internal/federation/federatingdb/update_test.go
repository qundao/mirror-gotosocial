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
	"time"

	"code.superseriousbusiness.org/gotosocial/internal/ap"
	"code.superseriousbusiness.org/gotosocial/internal/gtsmodel"
	"code.superseriousbusiness.org/gotosocial/testrig"
	"github.com/stretchr/testify/suite"
)

type UpdateTestSuite struct {
	FederatingDBTestSuite
}

func (suite *UpdateTestSuite) TestUpdateNewMention() {
	// Set up our test structs + tear down on finish.
	testStructs := testrig.SetupTestStructs(rMediaPath, rTemplatePath)
	defer testrig.TearDownTestStructs(testStructs)

	// Clean up test context when done.
	ctx, cncl := context.WithCancel(suite.T().Context())
	defer cncl()

	var (
		update         = suite.testActivities["remote_account_2_status_1_update"]
		receivingAcct  = suite.testAccounts["local_account_1"]
		requestingAcct = suite.testAccounts["remote_account_2"]
	)

	// Process Update of a note that adds an extra mention.
	note := update.Activity.GetActivityStreamsObject().At(0).GetActivityStreamsNote()
	if err := testStructs.Federator.FederatingDB().Update(
		createTestContext(
			ctx,
			requestingAcct,
			receivingAcct,
		),
		note,
	); err != nil {
		suite.FailNow(err.Error())
	}

	// Wait for the update.
	noteURI := ap.GetJSONLDId(note).String()
	var status *gtsmodel.Status
	var statusEdit *gtsmodel.StatusEdit
	if !testrig.WaitFor(func() bool {
		var err error
		status, err = testStructs.State.DB.GetStatusByURI(ctx, noteURI)
		if err != nil {
			suite.FailNow(err.Error())
		}
		if len(status.EditIDs) != 1 {
			return false
		}

		if err := testStructs.State.DB.PopulateStatusEdits(ctx, status); err != nil {
			suite.FailNow(err.Error())
		}

		statusEdit = status.Edits[0]
		return true
	}) {
		suite.FailNow("waiting for update")
	}

	// Time of the status edit should be set to the snapshot
	// of the OG version of the status, counter-intuitively.
	//
	// It would be more appropriate to think of a status edit
	// as a status REVISION, ie., a snapshot of the status
	// before the point in time it was edited.
	suite.Equal(status.CreatedAt.UTC(), statusEdit.CreatedAt.UTC())

	// The status editedAt time should be time of the
	// `updated` field on the latest version of the note.
	suite.EqualValues(status.EditedAt, ap.GetUpdated(note))

	// Status fetchedAt time should be now.
	suite.WithinDuration(time.Now(), status.FetchedAt, 5*time.Second)

	// Status should have
	// an extra mention now.
	suite.Len(status.Mentions, 2)
}

func TestUpdateTestSuite(t *testing.T) {
	suite.Run(t, new(UpdateTestSuite))
}
