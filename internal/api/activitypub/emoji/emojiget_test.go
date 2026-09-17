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

package emoji_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"code.superseriousbusiness.org/gopkg/httputil"
	"code.superseriousbusiness.org/gotosocial/internal/admin"
	"code.superseriousbusiness.org/gotosocial/internal/api/activitypub/emoji"
	apiutil "code.superseriousbusiness.org/gotosocial/internal/api/util"
	"code.superseriousbusiness.org/gotosocial/internal/db"
	"code.superseriousbusiness.org/gotosocial/internal/email"
	"code.superseriousbusiness.org/gotosocial/internal/federation"
	"code.superseriousbusiness.org/gotosocial/internal/gtsmodel"
	"code.superseriousbusiness.org/gotosocial/internal/media"
	"code.superseriousbusiness.org/gotosocial/internal/middleware"
	"code.superseriousbusiness.org/gotosocial/internal/processing"
	"code.superseriousbusiness.org/gotosocial/internal/state"
	"code.superseriousbusiness.org/gotosocial/internal/storage"
	"code.superseriousbusiness.org/gotosocial/internal/typeutils"
	"code.superseriousbusiness.org/gotosocial/testrig"
	"github.com/stretchr/testify/suite"
)

type EmojiGetTestSuite struct {
	suite.Suite
	db           db.DB
	tc           *typeutils.Converter
	mediaManager *media.Manager
	federator    *federation.Federator
	emailSender  email.Sender
	processor    *processing.Processor
	storage      *storage.Driver
	state        *state.State

	testEmojis   map[string]*gtsmodel.Emoji
	testAccounts map[string]*gtsmodel.Account

	emojiModule *emoji.Module

	signatureCheck httputil.Middleware
}

func (suite *EmojiGetTestSuite) SetupSuite() {
	suite.testAccounts = testrig.NewTestAccounts()
	suite.testEmojis = testrig.NewTestEmojis()
}

func (suite *EmojiGetTestSuite) SetupTest() {
	suite.state = new(state.State)
	testrig.StartNoopWorkers(suite.state)

	testrig.InitTestConfig()
	testrig.InitTestLog()

	suite.db = testrig.NewTestDB(suite.state)
	suite.state.DB = suite.db
	suite.state.AdminActions = admin.New(suite.state.DB, &suite.state.Workers)
	suite.storage = testrig.NewInMemoryStorage()
	suite.state.Storage = suite.storage
	suite.tc = typeutils.NewConverter(suite.state)

	suite.mediaManager = testrig.NewTestMediaManager(suite.state)
	suite.federator = testrig.NewTestFederator(suite.state, testrig.NewTestTransportController(suite.state, testrig.NewMockHTTPClient(nil, "../../../../testrig/media")), suite.mediaManager)
	suite.emailSender = testrig.NewEmailSender("../../../../web/template/", nil)
	suite.processor = testrig.NewTestProcessor(
		suite.state,
		suite.federator,
		suite.emailSender,
		testrig.NewNoopWebPushSender(),
		suite.mediaManager,
	)
	suite.emojiModule = emoji.New(suite.processor, testrig.LoadTemplates(suite.state, ""))
	testrig.StandardDBSetup(suite.db, suite.testAccounts)
	testrig.StandardStorageSetup(suite.storage, "../../../../testrig/media")

	suite.signatureCheck = middleware.ExtractSignature(suite.db.IsURIBlocked)
}

func (suite *EmojiGetTestSuite) TearDownTest() {
	testrig.StandardDBTeardown(suite.db)
	testrig.StandardStorageTeardown(suite.storage)
	testrig.StopWorkers(suite.state)
}

func (suite *EmojiGetTestSuite) TestGetEmoji() {
	// the dereference we're gonna use
	derefRequests := testrig.NewTestDereferenceRequests(suite.testAccounts)
	signedRequest := derefRequests["foss_satan_dereference_emoji"]
	targetEmoji := suite.testEmojis["rainbow"]

	// setup request
	req := httptest.NewRequest(http.MethodGet, targetEmoji.URI, nil) // the endpoint we're hitting
	recorder := httptest.NewRecorder()
	c := httputil.ToContext(recorder, req)
	c.R.Header.Set("accept", "application/activity+json")
	c.R.Header.Set("Signature", signedRequest.SignatureHeader)
	c.R.Header.Set("Date", signedRequest.DateHeader)

	// normally the router would populate these params from the path values,
	// but because we're calling the function directly, we need to set them manually.
	c.SetPathValue(apiutil.IDKey, targetEmoji.ID)

	// trigger the function being tested, first passing through sigcheck.
	suite.signatureCheck.Compile(suite.emojiModule.EmojiGETHandler)(c)

	// check response
	suite.EqualValues(http.StatusOK, recorder.Code)

	result := recorder.Result()
	defer result.Body.Close()
	b, err := io.ReadAll(result.Body)
	suite.NoError(err)

	suite.Contains(string(b), `"icon":{"mediaType":"image/png","type":"Image","url":"http://localhost:8080/fileserver/01AY6P665V14JJR0AFVRT7311Y/emoji/original/01F8MH9H8E4VG3KDYJR9EGPXCQ.png"},"id":"http://localhost:8080/emoji/01F8MH9H8E4VG3KDYJR9EGPXCQ","name":":rainbow:","type":"Emoji"`)
}

func TestEmojiGetTestSuite(t *testing.T) {
	suite.Run(t, new(EmojiGetTestSuite))
}
