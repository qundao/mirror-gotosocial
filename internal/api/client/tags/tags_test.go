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

package tags_test

import (
	"encoding/json"
	"io"
	"net/http/httptest"
	"strings"
	"testing"

	"code.superseriousbusiness.org/gopkg/httputil"
	"code.superseriousbusiness.org/gotosocial/internal/admin"
	"code.superseriousbusiness.org/gotosocial/internal/api/client/tags"
	apimodel "code.superseriousbusiness.org/gotosocial/internal/api/model"
	"code.superseriousbusiness.org/gotosocial/internal/config"
	"code.superseriousbusiness.org/gotosocial/internal/db"
	"code.superseriousbusiness.org/gotosocial/internal/email"
	"code.superseriousbusiness.org/gotosocial/internal/federation"
	"code.superseriousbusiness.org/gotosocial/internal/gtserror"
	"code.superseriousbusiness.org/gotosocial/internal/gtsmodel"
	"code.superseriousbusiness.org/gotosocial/internal/media"
	"code.superseriousbusiness.org/gotosocial/internal/oauth"
	"code.superseriousbusiness.org/gotosocial/internal/processing"
	"code.superseriousbusiness.org/gotosocial/internal/state"
	"code.superseriousbusiness.org/gotosocial/internal/storage"
	"code.superseriousbusiness.org/gotosocial/testrig"
	"github.com/stretchr/testify/suite"
)

type TagsTestSuite struct {
	suite.Suite
	db           db.DB
	storage      *storage.Driver
	mediaManager *media.Manager
	federator    *federation.Federator
	processor    *processing.Processor
	emailSender  email.Sender
	sentEmails   map[string]string
	state        *state.State

	// standard suite models
	testTokens       map[string]*gtsmodel.Token
	testApplications map[string]*gtsmodel.Application
	testUsers        map[string]*gtsmodel.User
	testAccounts     map[string]*gtsmodel.Account
	testTags         map[string]*gtsmodel.Tag

	// module being tested
	tagsModule *tags.Module
}

func (suite *TagsTestSuite) SetupSuite() {
	suite.testTokens = testrig.NewTestTokens()
	suite.testApplications = testrig.NewTestApplications()
	suite.testUsers = testrig.NewTestUsers()
	suite.testAccounts = testrig.NewTestAccounts()
	suite.testTags = testrig.NewTestTags()
}

func (suite *TagsTestSuite) SetupTest() {
	suite.state = new(state.State)
	testrig.StartNoopWorkers(suite.state)

	testrig.InitTestConfig()
	config.Config(func(cfg *config.Configuration) {
		cfg.WebAssetBaseDir = "../../../../web/assets/"
		cfg.WebTemplateBaseDir = "../../../../web/templates/"
	})
	testrig.InitTestLog()

	suite.db = testrig.NewTestDB(suite.state)
	suite.state.DB = suite.db
	suite.state.AdminActions = admin.New(suite.state.DB, &suite.state.Workers)
	suite.storage = testrig.NewInMemoryStorage()
	suite.state.Storage = suite.storage

	suite.mediaManager = testrig.NewTestMediaManager(suite.state)
	suite.federator = testrig.NewTestFederator(suite.state, testrig.NewTestTransportController(suite.state, testrig.NewMockHTTPClient(nil, "../../../../testrig/media")), suite.mediaManager)
	suite.sentEmails = make(map[string]string)
	suite.emailSender = testrig.NewEmailSender("../../../../web/template/", suite.sentEmails)
	suite.processor = testrig.NewTestProcessor(
		suite.state,
		suite.federator,
		suite.emailSender,
		testrig.NewNoopWebPushSender(),
		suite.mediaManager,
	)
	suite.tagsModule = tags.New(suite.processor, testrig.LoadTemplates(suite.state, ""))

	testrig.StandardDBSetup(suite.db, nil)
	testrig.StandardStorageSetup(suite.storage, "../../../../testrig/media")
}

func (suite *TagsTestSuite) TearDownTest() {
	testrig.StandardDBTeardown(suite.db)
	testrig.StandardStorageTeardown(suite.storage)
	testrig.StopWorkers(suite.state)
}

// tagAction gets, follows, or unfollows a tag, returning the tag.
func (suite *TagsTestSuite) tagAction(
	accountFixtureName string,
	tagName string,
	method string,
	path string,
	handler func(c *httputil.Context),
	expectedHTTPStatus int,
	expectedBody string,
) (*apimodel.Tag, error) {
	// create the request
	url := config.GetProtocol() + "://" + config.GetHost() + "/api/" + path
	req := httptest.NewRequest(method, strings.Replace(url, ":tag_name", tagName, 1), nil)
	req.Header.Set("accept", "application/json")

	// instantiate recorder + test context
	recorder := httptest.NewRecorder()
	c := httputil.ToContext(recorder, req)
	c.V.Set(oauth.SessionAuthorizedAccount, suite.testAccounts[accountFixtureName])
	c.V.Set(oauth.SessionAuthorizedToken, oauth.DBTokenToToken(suite.testTokens[accountFixtureName]))
	c.V.Set(oauth.SessionAuthorizedApplication, suite.testApplications["application_1"])
	c.V.Set(oauth.SessionAuthorizedUser, suite.testUsers[accountFixtureName])
	c.SetPathValue("tag_name", tagName)

	// trigger
	// the handler
	handler(c)

	// read the response
	result := recorder.Result()
	defer result.Body.Close()

	b, err := io.ReadAll(result.Body)
	if err != nil {
		return nil, err
	}

	errs := gtserror.NewMultiError(2)

	// check code + body
	if resultCode := recorder.Code; expectedHTTPStatus != resultCode {
		errs.Appendf("expected %d got %d", expectedHTTPStatus, resultCode)
		if expectedBody == "" {
			return nil, errs.Combine()
		}
	}

	// if we got an expected body, return early
	if expectedBody != "" {
		if string(b) != expectedBody {
			errs.Appendf("expected %s got %s", expectedBody, string(b))
		}
		return nil, errs.Combine()
	}

	resp := &apimodel.Tag{}
	if err := json.Unmarshal(b, resp); err != nil {
		return nil, err
	}

	return resp, nil
}

func TestTagsTestSuite(t *testing.T) {
	suite.Run(t, new(TagsTestSuite))
}
