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

	"code.superseriousbusiness.org/activity/streams/vocab"
	"code.superseriousbusiness.org/gotosocial/internal/ap"
	"code.superseriousbusiness.org/gotosocial/internal/gtsmodel"
	"code.superseriousbusiness.org/gotosocial/testrig"
	"github.com/stretchr/testify/suite"
)

type AnnounceTestSuite struct {
	FederatingDBTestSuite
}

func (suite *AnnounceTestSuite) TestNewAnnounce() {
	// Set up our test structs + tear down on finish.
	testStructs := testrig.SetupTestStructs(rMediaPath, rTemplatePath)
	defer testrig.TearDownTestStructs(testStructs)

	// Clean up test context when done.
	ctx, cncl := context.WithCancel(suite.T().Context())
	defer cncl()

	receivingAccount1 := suite.testAccounts["local_account_1"]
	announcingAccount := suite.testAccounts["remote_account_1"]

	// Call the Announce function.
	ctx = createTestContext(ctx, announcingAccount, receivingAccount1)
	announce := suite.testActivities["announce_forwarded_1_zork"].Activity.(vocab.ActivityStreamsAnnounce)
	if err := testStructs.Federator.FederatingDB().Announce(ctx, announce); err != nil {
		suite.FailNow(err.Error())
	}

	// Wait for the boost wrapper
	// status to appear in the db.
	announceURI := ap.GetJSONLDId(announce).String()
	var boost *gtsmodel.Status
	if !testrig.WaitFor(func() bool {
		var err error
		boost, err = testStructs.State.DB.GetStatusByURI(ctx, announceURI)
		return err == nil && boost != nil
	}) {
		suite.FailNow("timed out waiting for boost")
	}

	// Boost should refer to the Announced status.
	boostOfURI, err := ap.GetOneObjectIRI(announce)
	if err != nil {
		suite.FailNow(err.Error())
	}
	suite.Equal(boostOfURI.String(), boost.BoostOf.URI)

	// Boost should be attributed to the announcing account.
	suite.Equal(announcingAccount.ID, boost.AccountID)
}

func TestAnnounceTestSuite(t *testing.T) {
	suite.Run(t, &AnnounceTestSuite{})
}
