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

package admin_test

import (
	"bytes"
	"fmt"
	"net/http/httptest"

	"code.superseriousbusiness.org/gopkg/httputil"
	adminactions "code.superseriousbusiness.org/gotosocial/internal/admin"
	"code.superseriousbusiness.org/gotosocial/internal/api/client/admin"
	"code.superseriousbusiness.org/gotosocial/internal/config"
	"code.superseriousbusiness.org/gotosocial/internal/gtsmodel"
	"code.superseriousbusiness.org/gotosocial/internal/oauth"
	"code.superseriousbusiness.org/gotosocial/internal/processing"
	"code.superseriousbusiness.org/gotosocial/internal/state"
	"code.superseriousbusiness.org/gotosocial/testrig"
	"github.com/stretchr/testify/suite"
)

type AdminStandardTestSuite struct {
	// standard suite interfaces
	suite.Suite
	processor  *processing.Processor
	sentEmails map[string]string
	state      *state.State

	// standard suite models
	testTokens          map[string]*gtsmodel.Token
	testApplications    map[string]*gtsmodel.Application
	testUsers           map[string]*gtsmodel.User
	testAccounts        map[string]*gtsmodel.Account
	testAttachments     map[string]*gtsmodel.MediaAttachment
	testStatuses        map[string]*gtsmodel.Status
	testEmojis          map[string]*gtsmodel.Emoji
	testEmojiCategories map[string]*gtsmodel.EmojiCategory
	testReports         map[string]*gtsmodel.Report

	// module being tested
	adminModule *admin.Module
}

func (suite *AdminStandardTestSuite) SetupSuite() {
	suite.testTokens = testrig.NewTestTokens()
	suite.testApplications = testrig.NewTestApplications()
	suite.testUsers = testrig.NewTestUsers()
	suite.testAccounts = testrig.NewTestAccounts()
	suite.testAttachments = testrig.NewTestAttachments()
	suite.testStatuses = testrig.NewTestStatuses()
	suite.testEmojis = testrig.NewTestEmojis()
	suite.testEmojiCategories = testrig.NewTestEmojiCategories()
	suite.testReports = testrig.NewTestReports()
}

func (suite *AdminStandardTestSuite) SetupTest() {
	suite.state = new(state.State)

	testrig.InitTestConfig()
	testrig.InitTestLog()

	suite.state.DB = testrig.NewTestDB(suite.state)
	suite.state.Storage = testrig.NewInMemoryStorage()

	testrig.StandardDBSetup(suite.state.DB, nil)
	testrig.StandardStorageSetup(suite.state.Storage, "../../../../testrig/media")

	suite.state.AdminActions = adminactions.New(suite.state.DB, &suite.state.Workers)

	suite.sentEmails = make(map[string]string)
	suite.processor = testrig.NewTestProcessor(
		suite.state,
		testrig.NewTestFederator(
			suite.state,
			testrig.NewTestTransportController(
				suite.state,
				testrig.NewMockHTTPClient(nil, "../../../../testrig/media"),
			),
			testrig.NewTestMediaManager(suite.state),
		),
		testrig.NewEmailSender("../../../../web/template/", suite.sentEmails),
		testrig.NewNoopWebPushSender(),
		testrig.NewTestMediaManager(suite.state),
	)
	testrig.StartWorkers(suite.state, suite.processor.Workers())
	suite.adminModule = admin.New(suite.state, suite.processor, testrig.LoadTemplates(suite.state, ""))
}

func (suite *AdminStandardTestSuite) TearDownTest() {
	testrig.StandardDBTeardown(suite.state.DB)
	testrig.StandardStorageTeardown(suite.state.Storage)
	testrig.StopWorkers(suite.state)
}

func (suite *AdminStandardTestSuite) newContext(recorder *httptest.ResponseRecorder, requestMethod string, requestBody []byte, requestPath string, bodyContentType string) *httputil.Context {
	protocol := config.GetProtocol()
	host := config.GetHost()

	baseURI := fmt.Sprintf("%s://%s", protocol, host)
	requestURI := fmt.Sprintf("%s/%s", baseURI, requestPath)

	req := httptest.NewRequest(requestMethod, requestURI, bytes.NewReader(requestBody)) // the endpoint we're hitting

	c := httputil.ToContext(recorder, req)

	c.V.Set(oauth.SessionAuthorizedAccount, suite.testAccounts["admin_account"])
	c.V.Set(oauth.SessionAuthorizedToken, oauth.DBTokenToToken(suite.testTokens["admin_account"]))
	c.V.Set(oauth.SessionAuthorizedApplication, suite.testApplications["admin_account"])
	c.V.Set(oauth.SessionAuthorizedUser, suite.testUsers["admin_account"])

	if bodyContentType != "" {
		c.R.Header.Set("Content-Type", bodyContentType)
	}

	c.R.Header.Set("accept", "application/json")

	return c
}
