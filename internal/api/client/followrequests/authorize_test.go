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

package followrequests_test

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	apiutil "code.superseriousbusiness.org/gotosocial/internal/api/util"
	"code.superseriousbusiness.org/gotosocial/internal/gtsmodel"
	"code.superseriousbusiness.org/gotosocial/testrig"
	"github.com/stretchr/testify/suite"
)

type AuthorizeTestSuite struct {
	FollowRequestStandardTestSuite
}

func (suite *AuthorizeTestSuite) TestAuthorize() {
	requestingAccount := suite.testAccounts["remote_account_2"]
	targetAccount := suite.testAccounts["local_account_1"]

	// put a follow request in the database
	fr := &gtsmodel.FollowRequest{
		ID:              "01FJ1S8DX3STJJ6CEYPMZ1M0R3",
		CreatedAt:       time.Now(),
		URI:             fmt.Sprintf("%s/follow/01FJ1S8DX3STJJ6CEYPMZ1M0R3", requestingAccount.URI),
		AccountID:       requestingAccount.ID,
		TargetAccountID: targetAccount.ID,
	}

	err := suite.db.Put(suite.T().Context(), fr)
	suite.NoError(err)

	recorder := httptest.NewRecorder()
	c := suite.newContext(recorder, http.MethodPost, []byte{}, fmt.Sprintf("/api/v1/follow_requests/%s/authorize", requestingAccount.ID), "")
	c.SetPathValue(apiutil.IDKey, requestingAccount.ID)

	// call the handler
	suite.followRequestModule.FollowRequestAuthorizePOSTHandler(c)

	// 1. we should have OK because our request was valid
	suite.Equal(http.StatusOK, recorder.Code)

	// 2. we should have no error message in the result body
	result := recorder.Result()
	defer result.Body.Close()

	// check the response
	b, err := io.ReadAll(result.Body)
	if err != nil {
		suite.FailNow(err.Error())
	}

	out := testrig.MustJSONStringFromBytes(b)
	suite.Equal(`{
  "blocked_by": false,
  "blocking": false,
  "domain_blocking": false,
  "endorsed": false,
  "followed_by": true,
  "following": false,
  "id": "01FHMQX3GAABWSM0S2VZEC2SWC",
  "muting": false,
  "muting_notifications": false,
  "note": "",
  "notifying": false,
  "requested": false,
  "requested_by": false,
  "showing_reblogs": false
}`, out)
}

func (suite *AuthorizeTestSuite) TestAuthorizeNoFR() {
	requestingAccount := suite.testAccounts["remote_account_2"]

	recorder := httptest.NewRecorder()
	c := suite.newContext(recorder, http.MethodPost, []byte{}, fmt.Sprintf("/api/v1/follow_requests/%s/authorize", requestingAccount.ID), "")
	c.SetPathValue(apiutil.IDKey, requestingAccount.ID)

	// call the handler
	suite.followRequestModule.FollowRequestAuthorizePOSTHandler(c)

	suite.Equal(http.StatusNotFound, recorder.Code)

	result := recorder.Result()
	defer result.Body.Close()

	// check the response
	b, err := io.ReadAll(result.Body)
	suite.NoError(err)

	suite.Equal(`{"error":"Not Found"}`, string(b))
}

func TestAuthorizeTestSuite(t *testing.T) {
	suite.Run(t, &AuthorizeTestSuite{})
}
