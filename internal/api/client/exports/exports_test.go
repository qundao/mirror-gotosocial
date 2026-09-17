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

package exports_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"code.superseriousbusiness.org/gopkg/httputil"
	"code.superseriousbusiness.org/gotosocial/internal/api/client/exports"
	apiutil "code.superseriousbusiness.org/gotosocial/internal/api/util"
	"code.superseriousbusiness.org/gotosocial/internal/gtsmodel"
	"code.superseriousbusiness.org/gotosocial/internal/oauth"
	"code.superseriousbusiness.org/gotosocial/internal/state"
	"code.superseriousbusiness.org/gotosocial/testrig"
	"github.com/stretchr/testify/suite"
)

type ExportsTestSuite struct {
	// Suite interfaces
	suite.Suite
	state *state.State

	// standard suite models
	testTokens       map[string]*gtsmodel.Token
	testApplications map[string]*gtsmodel.Application
	testUsers        map[string]*gtsmodel.User
	testAccounts     map[string]*gtsmodel.Account

	// module being tested
	exportsModule *exports.Module
}

func (suite *ExportsTestSuite) SetupSuite() {
	suite.testTokens = testrig.NewTestTokens()
	suite.testApplications = testrig.NewTestApplications()
	suite.testUsers = testrig.NewTestUsers()
	suite.testAccounts = testrig.NewTestAccounts()
}

func (suite *ExportsTestSuite) SetupTest() {
	suite.state = new(state.State)
	testrig.StartNoopWorkers(suite.state)

	testrig.InitTestConfig()
	testrig.InitTestLog()

	suite.state.DB = testrig.NewTestDB(suite.state)
	suite.state.Storage = testrig.NewInMemoryStorage()

	testrig.StandardDBSetup(suite.state.DB, nil)
	testrig.StandardStorageSetup(suite.state.Storage, "../../../../testrig/media")

	mediaManager := testrig.NewTestMediaManager(suite.state)

	federator := testrig.NewTestFederator(
		suite.state,
		testrig.NewTestTransportController(
			suite.state,
			testrig.NewMockHTTPClient(nil, "../../../../testrig/media"),
		),
		mediaManager,
	)

	processor := testrig.NewTestProcessor(
		suite.state,
		federator,
		testrig.NewEmailSender("../../../../web/template/", nil),
		testrig.NewNoopWebPushSender(),
		mediaManager,
	)

	suite.exportsModule = exports.New(processor, testrig.LoadTemplates(suite.state, ""))
}

func (suite *ExportsTestSuite) TriggerHandler(
	handler httputil.HandlerFunc,
	path string,
	contentType string,
	application *gtsmodel.Application,
	token *gtsmodel.Token,
	user *gtsmodel.User,
	account *gtsmodel.Account,
) *httptest.ResponseRecorder {
	// Set up request.
	target := "http://localhost:8080/api" + path
	req := httptest.NewRequest(http.MethodGet, target, nil)
	req.Header.Set("Accept", contentType)
	recorder := httptest.NewRecorder()
	c := httputil.ToContext(recorder, req)

	// Authorize the request ctx as though it
	// had passed through API auth handlers.
	c.V.Set(oauth.SessionAuthorizedApplication, application)
	c.V.Set(oauth.SessionAuthorizedToken, oauth.DBTokenToToken(token))
	c.V.Set(oauth.SessionAuthorizedUser, user)
	c.V.Set(oauth.SessionAuthorizedAccount, account)

	// Trigger handler.
	handler(c)

	return recorder
}

func (suite *ExportsTestSuite) TearDownTest() {
	testrig.StandardDBTeardown(suite.state.DB)
	testrig.StandardStorageTeardown(suite.state.Storage)
	testrig.StopWorkers(suite.state)
}

func (suite *ExportsTestSuite) TestExports() {
	type testCase struct {
		handler     httputil.HandlerFunc
		path        string
		contentType string
		application *gtsmodel.Application
		token       *gtsmodel.Token
		user        *gtsmodel.User
		account     *gtsmodel.Account
		expect      string
	}

	testCases := []testCase{
		// Export Following
		{
			handler:     suite.exportsModule.ExportFollowingGETHandler,
			path:        exports.FollowingPath,
			contentType: apiutil.TextCSV,
			application: suite.testApplications["application_1"],
			token:       suite.testTokens["local_account_1"],
			user:        suite.testUsers["local_account_1"],
			account:     suite.testAccounts["local_account_1"],
			expect: `Account address,Show boosts,Notify on new posts,Languages
1happyturtle@localhost:8080,true,false,
admin@localhost:8080,true,false,
`,
		},
		// Export Followers.
		{
			handler:     suite.exportsModule.ExportFollowersGETHandler,
			path:        exports.FollowingPath,
			contentType: apiutil.TextCSV,
			application: suite.testApplications["application_1"],
			token:       suite.testTokens["local_account_1"],
			user:        suite.testUsers["local_account_1"],
			account:     suite.testAccounts["local_account_1"],
			expect: `Account address
1happyturtle@localhost:8080
admin@localhost:8080
`,
		},
		// Export Lists.
		{
			handler:     suite.exportsModule.ExportListsGETHandler,
			path:        exports.ListsPath,
			contentType: apiutil.TextCSV,
			application: suite.testApplications["application_1"],
			token:       suite.testTokens["local_account_1"],
			user:        suite.testUsers["local_account_1"],
			account:     suite.testAccounts["local_account_1"],
			expect: `Cool Ass Posters From This Instance,1happyturtle@localhost:8080
Cool Ass Posters From This Instance,admin@localhost:8080
`,
		},
		// Export Mutes.
		{
			handler:     suite.exportsModule.ExportMutesGETHandler,
			path:        exports.MutesPath,
			contentType: apiutil.TextCSV,
			application: suite.testApplications["application_1"],
			token:       suite.testTokens["local_account_1"],
			user:        suite.testUsers["local_account_1"],
			account:     suite.testAccounts["local_account_1"],
			expect: `Account address,Hide notifications
`,
		},
		// Export Blocks.
		{
			handler:     suite.exportsModule.ExportBlocksGETHandler,
			path:        exports.BlocksPath,
			contentType: apiutil.TextCSV,
			application: suite.testApplications["application_1"],
			token:       suite.testTokens["local_account_2"],
			user:        suite.testUsers["local_account_2"],
			account:     suite.testAccounts["local_account_2"],
			expect: `foss_satan@fossbros-anonymous.io
`,
		},
		// Export Stats.
		{
			handler:     suite.exportsModule.ExportStatsGETHandler,
			path:        exports.StatsPath,
			contentType: apiutil.AppJSON,
			application: suite.testApplications["application_1"],
			token:       suite.testTokens["local_account_1"],
			user:        suite.testUsers["local_account_1"],
			account:     suite.testAccounts["local_account_1"],
			expect: `{
  "blocks_count": 0,
  "followers_count": 2,
  "following_count": 2,
  "lists_count": 1,
  "media_storage": "",
  "mutes_count": 0,
  "statuses_count": 9
}`,
		},
	}

	for _, test := range testCases {
		recorder := suite.TriggerHandler(
			test.handler,
			test.path,
			test.contentType,
			test.application,
			test.token,
			test.user,
			test.account,
		)

		// Check response code.
		suite.EqualValues(http.StatusOK, recorder.Code)

		// Check response body.
		b, err := io.ReadAll(recorder.Body)
		if err != nil {
			suite.FailNow(err.Error())
		}

		// If json response, indent it nicely.
		if recorder.Result().Header.Get("Content-Type") == "application/json" {
			b = testrig.MustJSONBytesFromBytes(b)
		}

		suite.Equal(test.expect, string(b))
	}
}

func TestExportsTestSuite(t *testing.T) {
	suite.Run(t, new(ExportsTestSuite))
}
