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
	"code.superseriousbusiness.org/gotosocial/testrig"
	"github.com/stretchr/testify/suite"
)

type FollowingTestSuite struct {
	FederatingDBTestSuite
}

func (suite *FollowingTestSuite) TestGetFollowing() {
	// Set up our test structs + tear down on finish.
	testStructs := testrig.SetupTestStructs(rMediaPath, rTemplatePath)
	defer testrig.TearDownTestStructs(testStructs)

	// Clean up test context when done.
	ctx, cncl := context.WithCancel(suite.T().Context())
	defer cncl()

	testAccount := suite.testAccounts["local_account_1"]

	f, err := testStructs.Federator.FederatingDB().Following(
		ctx, testrig.URLMustParse(testAccount.URI),
	)
	if err != nil {
		suite.FailNow(err.Error())
	}

	fi, err := ap.Serialize(f)
	if err != nil {
		suite.FailNow(err.Error())
	}

	// zork follows admin
	// and local_account_1
	out := testrig.MustJSONString(fi)
	suite.Equal(`{
  "@context": "https://www.w3.org/ns/activitystreams",
  "items": [
    "http://localhost:8080/users/admin",
    "http://localhost:8080/users/1happyturtle"
  ],
  "type": "Collection"
}`, out)
}

func TestFollowingTestSuite(t *testing.T) {
	suite.Run(t, &FollowingTestSuite{})
}
