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

package accounts_test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"code.superseriousbusiness.org/gotosocial/internal/api/client/accounts"
	apimodel "code.superseriousbusiness.org/gotosocial/internal/api/model"
	"code.superseriousbusiness.org/gotosocial/internal/id"
	"code.superseriousbusiness.org/gotosocial/internal/oauth"
	"github.com/stretchr/testify/suite"
)

type AccountsGetTestSuite struct {
	AccountStandardTestSuite
}

// getAccounts calls the get multiple accounts API method as requester
// (or the default local_account_1 if empty), with the given account IDs.
func (suite *AccountsGetTestSuite) getAccounts(requester string, ids []string, expectedHTTPStatus int, expectedBody string) []apimodel.Account {
	query := url.Values{}
	for _, id := range ids {
		query.Add("id[]", id)
	}

	recorder := httptest.NewRecorder()
	c := suite.newContext(recorder, http.MethodGet, nil, "api"+accounts.BasePath+"?"+query.Encode(), "")

	if requester != "" {
		c.V.Set(oauth.SessionAuthorizedAccount, suite.testAccounts[requester])
		c.V.Set(oauth.SessionAuthorizedToken, oauth.DBTokenToToken(suite.testTokens[requester]))
		c.V.Set(oauth.SessionAuthorizedUser, suite.testUsers[requester])
	}

	// call the handler
	suite.accountsModule.AccountsGETHandler(c)

	result := recorder.Result()
	defer result.Body.Close()

	b, err := io.ReadAll(result.Body)
	if err != nil {
		suite.FailNow(err.Error())
	}

	suite.Equal(expectedHTTPStatus, result.StatusCode, string(b))
	if expectedBody != "" {
		suite.Equal(expectedBody, string(b))
		return nil
	}

	var apiAccounts []apimodel.Account
	if err := json.Unmarshal(b, &apiAccounts); err != nil {
		suite.FailNow(err.Error())
	}

	return apiAccounts
}

func accountIDs(apiAccounts []apimodel.Account) []string {
	ids := make([]string, 0, len(apiAccounts))
	for _, a := range apiAccounts {
		ids = append(ids, a.ID)
	}
	return ids
}

// Accounts should be returned in requested order, with
// duplicates removed and missing accounts omitted.
func (suite *AccountsGetTestSuite) TestGetAccounts() {
	var (
		self   = suite.testAccounts["local_account_1"]
		local  = suite.testAccounts["local_account_2"]
		remote = suite.testAccounts["remote_account_1"]
		admin  = suite.testAccounts["admin_account"]
	)

	apiAccounts := suite.getAccounts("", []string{
		remote.ID,
		id.NewULID(), // doesn't exist
		self.ID,
		local.ID,
		remote.ID, // duplicate
		admin.ID,
	}, http.StatusOK, "")

	suite.Equal([]string{
		remote.ID,
		self.ID,
		local.ID,
		admin.ID,
	}, accountIDs(apiAccounts))

	// Only requester's own account should show source.
	for _, a := range apiAccounts {
		if a.ID == self.ID {
			suite.NotNil(a.Source)
		} else {
			suite.Nil(a.Source)
		}
	}

	// Non-blocked account should show the public view.
	suite.Equal(remote.DisplayName, apiAccounts[0].DisplayName)
}

// Blocked accounts should only show the minimal blocked view.
func (suite *AccountsGetTestSuite) TestGetAccountsBlocked() {
	var (
		blocked = suite.testAccounts["remote_account_1"]
		other   = suite.testAccounts["local_account_1"]
	)

	// local_account_2 blocks remote_account_1.
	apiAccounts := suite.getAccounts("local_account_2", []string{
		blocked.ID,
		other.ID,
	}, http.StatusOK, "")

	if !suite.Len(apiAccounts, 2) {
		suite.FailNow("")
	}

	suite.Equal(blocked.ID, apiAccounts[0].ID)
	suite.Empty(apiAccounts[0].DisplayName)
	suite.Empty(apiAccounts[0].Note)

	suite.Equal(other.ID, apiAccounts[1].ID)
	suite.Equal(other.DisplayName, apiAccounts[1].DisplayName)
}

func (suite *AccountsGetTestSuite) TestGetAccountsNoneExist() {
	apiAccounts := suite.getAccounts("", []string{
		id.NewULID(),
		id.NewULID(),
	}, http.StatusOK, "")
	suite.Empty(apiAccounts)
}

func (suite *AccountsGetTestSuite) TestGetAccountsNoIDs() {
	suite.getAccounts("", nil, http.StatusOK, `[]`)
}

func (suite *AccountsGetTestSuite) TestGetAccountsEmptyIDs() {
	suite.getAccounts("", []string{"", ""}, http.StatusOK, `[]`)
}

func (suite *AccountsGetTestSuite) TestGetAccountsTooMany() {
	ids := make([]string, 41)
	for i := range ids {
		ids[i] = id.NewULID()
	}

	suite.getAccounts("", ids,
		http.StatusUnprocessableEntity,
		`{"error":"Unprocessable Entity: too many account ids specified"}`,
	)
}

func TestAccountsGetTestSuite(t *testing.T) {
	suite.Run(t, new(AccountsGetTestSuite))
}
