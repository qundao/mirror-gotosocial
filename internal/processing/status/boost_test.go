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

package status_test

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type StatusBoostTestSuite struct {
	StatusStandardTestSuite
}

func (suite *StatusBoostTestSuite) TestBoostOfBoost() {
	ctx := suite.T().Context()

	// first boost a status, no big deal
	boostingAccount1 := suite.testAccounts["local_account_1"]
	application1 := suite.testApplications["application_1"]
	targetStatus1 := suite.testStatuses["admin_account_status_1"]

	boost1, err := suite.status.BoostCreate(ctx, boostingAccount1, application1, targetStatus1.ID)
	suite.NoError(err)
	suite.NotNil(boost1)
	suite.Equal(targetStatus1.ID, boost1.Reblog.ID)

	// now take another account and boost that boost
	boostingAccount2 := suite.testAccounts["local_account_2"]
	application2 := suite.testApplications["application_2"]
	targetStatus2ID := boost1.ID

	boost2, err := suite.status.BoostCreate(ctx, boostingAccount2, application2, targetStatus2ID)
	suite.NoError(err)
	suite.NotNil(boost2)
	// the boosted status should not be the boost,
	// but the original status that was boosted
	suite.Equal(targetStatus1.ID, boost2.Reblog.ID)
}

func (suite *StatusBoostTestSuite) TestBoostUnboost() {
	var (
		ctx      = suite.T().Context()
		acct     = suite.testAccounts["local_account_1"]
		app      = suite.testApplications["application_1"]
		statusID = suite.testStatuses["admin_account_status_1"].ID
	)

	// Boost the status.
	_, err := suite.status.BoostCreate(ctx, acct, app, statusID)
	if err != nil {
		suite.FailNow(err.Error())
	}

	// Check status is boosted by the caller.
	status, err := suite.status.Get(ctx, acct, statusID)
	if err != nil {
		suite.FailNow(err.Error())
	}
	if !status.Reblogged {
		suite.FailNow("", "expected reblogged=true, got false")
	}

	// Unboost that status.
	// Check it's not boosted by the caller.
	status, err = suite.status.BoostRemove(ctx, acct, app, statusID)
	if err != nil {
		suite.FailNow(err.Error())
	}
	if status.Reblogged {
		suite.FailNow("", "expected reblogged=false, got true")
	}
}

func TestStatusBoostTestSuite(t *testing.T) {
	suite.Run(t, new(StatusBoostTestSuite))
}
