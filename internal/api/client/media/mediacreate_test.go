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

package media_test

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"code.superseriousbusiness.org/gopkg/httputil"
	mediamodule "code.superseriousbusiness.org/gotosocial/internal/api/client/media"
	apimodel "code.superseriousbusiness.org/gotosocial/internal/api/model"
	apiutil "code.superseriousbusiness.org/gotosocial/internal/api/util"
	"code.superseriousbusiness.org/gotosocial/internal/config"
	"code.superseriousbusiness.org/gotosocial/internal/db"
	"code.superseriousbusiness.org/gotosocial/internal/email"
	"code.superseriousbusiness.org/gotosocial/internal/federation"
	"code.superseriousbusiness.org/gotosocial/internal/gtsmodel"
	"code.superseriousbusiness.org/gotosocial/internal/media"
	"code.superseriousbusiness.org/gotosocial/internal/oauth"
	"code.superseriousbusiness.org/gotosocial/internal/processing"
	"code.superseriousbusiness.org/gotosocial/internal/state"
	"code.superseriousbusiness.org/gotosocial/internal/storage"
	"code.superseriousbusiness.org/gotosocial/internal/typeutils"
	"code.superseriousbusiness.org/gotosocial/testrig"
	"github.com/stretchr/testify/suite"
)

type MediaCreateTestSuite struct {
	// standard suite interfaces
	suite.Suite
	db           db.DB
	storage      *storage.Driver
	mediaManager *media.Manager
	federator    *federation.Federator
	tc           *typeutils.Converter
	oauthServer  oauth.Server
	emailSender  email.Sender
	processor    *processing.Processor
	state        *state.State

	// standard suite models
	testTokens       map[string]*gtsmodel.Token
	testApplications map[string]*gtsmodel.Application
	testUsers        map[string]*gtsmodel.User
	testAccounts     map[string]*gtsmodel.Account
	testAttachments  map[string]*gtsmodel.MediaAttachment

	// item being tested
	mediaModule *mediamodule.Module
}

/*
	TEST INFRASTRUCTURE
*/

func (suite *MediaCreateTestSuite) SetupTest() {
	suite.state = new(state.State)
	testrig.StartNoopWorkers(suite.state)

	// setup standard items
	testrig.InitTestConfig()
	testrig.InitTestLog()

	suite.storage = testrig.NewInMemoryStorage()
	suite.state.Storage = suite.storage

	suite.db = testrig.NewTestDB(suite.state)
	testrig.StandardDBSetup(suite.db, nil)
	testrig.StandardStorageSetup(suite.storage, "../../../../testrig/media")

	suite.tc = typeutils.NewConverter(suite.state)

	suite.mediaManager = testrig.NewTestMediaManager(suite.state)
	suite.oauthServer = testrig.NewTestOauthServer(suite.state)
	suite.federator = testrig.NewTestFederator(suite.state, testrig.NewTestTransportController(suite.state, testrig.NewMockHTTPClient(nil, "../../../../testrig/media")), suite.mediaManager)
	suite.emailSender = testrig.NewEmailSender("../../../../web/template/", nil)
	suite.processor = testrig.NewTestProcessor(
		suite.state,
		suite.federator,
		suite.emailSender,
		testrig.NewNoopWebPushSender(),
		suite.mediaManager,
	)

	// setup module being tested
	suite.mediaModule = mediamodule.New(suite.processor, testrig.LoadTemplates(suite.state, ""))

	// setup test data
	suite.testTokens = testrig.NewTestTokens()
	suite.testApplications = testrig.NewTestApplications()
	suite.testUsers = testrig.NewTestUsers()
	suite.testAccounts = testrig.NewTestAccounts()
	suite.testAttachments = testrig.NewTestAttachments()
}

func (suite *MediaCreateTestSuite) TearDownTest() {
	testrig.StandardDBTeardown(suite.db)
	testrig.StandardStorageTeardown(suite.storage)
	testrig.StopWorkers(suite.state)
}

/*
	ACTUAL TESTS
*/

func (suite *MediaCreateTestSuite) TestMediaCreateSuccessful() {
	// set up the context for the request
	t := suite.testTokens["local_account_1"]
	oauthToken := oauth.DBTokenToToken(t)
	recorder := httptest.NewRecorder()

	// create the request
	buf, w, err := testrig.CreateMultipartFormData(testrig.FileToDataF("file", "../../../../testrig/media/test-jpeg.jpg"), map[string][]string{
		"description": {"this is a test image -- a cool background from somewhere"},
		"focus":       {"-0.5,0.5"},
	})
	if err != nil {
		panic(err)
	}
	req := httptest.NewRequest(http.MethodPost, "http://localhost:8080/api/v1/media", bytes.NewReader(buf.Bytes())) // the endpoint we're hitting
	req.Header.Set("Content-Type", w.FormDataContentType())
	req.Header.Set("accept", "application/json")

	c := httputil.ToContext(recorder, req)
	c.V.Set(oauth.SessionAuthorizedApplication, suite.testApplications["application_1"])
	c.V.Set(oauth.SessionAuthorizedToken, oauthToken)
	c.V.Set(oauth.SessionAuthorizedUser, suite.testUsers["local_account_1"])
	c.V.Set(oauth.SessionAuthorizedAccount, suite.testAccounts["local_account_1"])

	// see what's in storage *before* the request
	var storageKeysBeforeRequest []string
	if err := suite.storage.WalkKeys(c, func(key string) error {
		storageKeysBeforeRequest = append(storageKeysBeforeRequest, key)
		return nil
	}); err != nil {
		panic(err)
	}

	c.SetPathValue(apiutil.APIVersionKey, apiutil.APIv1)

	// do the actual request
	suite.mediaModule.MediaCreatePOSTHandler(c)

	// check what's in storage *after* the request
	var storageKeysAfterRequest []string
	if err := suite.storage.WalkKeys(c, func(key string) error {
		storageKeysAfterRequest = append(storageKeysAfterRequest, key)
		return nil
	}); err != nil {
		panic(err)
	}

	// check response
	suite.EqualValues(http.StatusOK, recorder.Code)

	result := recorder.Result()
	defer result.Body.Close()
	b, err := io.ReadAll(result.Body)
	suite.NoError(err)
	fmt.Println(string(b))

	attachmentReply := &apimodel.Attachment{}
	err = json.Unmarshal(b, attachmentReply)
	suite.NoError(err)

	suite.Equal("this is a test image -- a cool background from somewhere", *attachmentReply.Description)
	suite.Equal("image", attachmentReply.Type)
	suite.EqualValues(apimodel.MediaMeta{
		Original: apimodel.MediaDimensions{
			Width:  1920,
			Height: 1080,
			Size:   "1920x1080",
			Aspect: 1.7777778,
		},
		Small: apimodel.MediaDimensions{
			Width:  512,
			Height: 288,
			Size:   "512x288",
			Aspect: 1.7777778,
		},
		Focus: &apimodel.MediaFocus{
			X: -0.5,
			Y: 0.5,
		},
	}, *attachmentReply.Meta)
	suite.Equal("LiB|W-#6RQR.~qvzRjWF_3rqV@a$", *attachmentReply.Blurhash)
	suite.NotEmpty(attachmentReply.ID)
	suite.NotEmpty(attachmentReply.URL)
	suite.NotEmpty(attachmentReply.PreviewURL)
	suite.Equal(len(storageKeysBeforeRequest)+2, len(storageKeysAfterRequest)) // 2 images should be added to storage: the original and the thumbnail
}

func (suite *MediaCreateTestSuite) TestMediaCreateSuccessfulV2() {
	// set up the context for the request
	t := suite.testTokens["local_account_1"]
	oauthToken := oauth.DBTokenToToken(t)
	recorder := httptest.NewRecorder()

	// create the request
	buf, w, err := testrig.CreateMultipartFormData(testrig.FileToDataF("file", "../../../../testrig/media/test-jpeg.jpg"), map[string][]string{
		"description": {"this is a test image -- a cool background from somewhere"},
		"focus":       {"-0.5,0.5"},
	})
	if err != nil {
		panic(err)
	}
	req := httptest.NewRequest(http.MethodPost, "http://localhost:8080/api/v2/media", bytes.NewReader(buf.Bytes())) // the endpoint we're hitting
	req.Header.Set("Content-Type", w.FormDataContentType())
	req.Header.Set("accept", "application/json")

	c := httputil.ToContext(recorder, req)
	c.V.Set(oauth.SessionAuthorizedApplication, suite.testApplications["application_1"])
	c.V.Set(oauth.SessionAuthorizedToken, oauthToken)
	c.V.Set(oauth.SessionAuthorizedUser, suite.testUsers["local_account_1"])
	c.V.Set(oauth.SessionAuthorizedAccount, suite.testAccounts["local_account_1"])

	// see what's in storage *before* the request
	var storageKeysBeforeRequest []string
	if err := suite.storage.WalkKeys(c, func(key string) error {
		storageKeysBeforeRequest = append(storageKeysBeforeRequest, key)
		return nil
	}); err != nil {
		panic(err)
	}

	c.SetPathValue(apiutil.APIVersionKey, apiutil.APIv2)

	// do the actual request
	suite.mediaModule.MediaCreatePOSTHandler(c)

	// check what's in storage *after* the request
	var storageKeysAfterRequest []string
	if err := suite.storage.WalkKeys(c, func(key string) error {
		storageKeysAfterRequest = append(storageKeysAfterRequest, key)
		return nil
	}); err != nil {
		panic(err)
	}

	// check response
	suite.EqualValues(http.StatusOK, recorder.Code)

	result := recorder.Result()
	defer result.Body.Close()
	b, err := io.ReadAll(result.Body)
	suite.NoError(err)
	fmt.Println(string(b))

	attachmentReply := &apimodel.Attachment{}
	err = json.Unmarshal(b, attachmentReply)
	suite.NoError(err)

	suite.Equal("this is a test image -- a cool background from somewhere", *attachmentReply.Description)
	suite.Equal("image", attachmentReply.Type)
	suite.EqualValues(apimodel.MediaMeta{
		Original: apimodel.MediaDimensions{
			Width:  1920,
			Height: 1080,
			Size:   "1920x1080",
			Aspect: 1.7777778,
		},
		Small: apimodel.MediaDimensions{
			Width:  512,
			Height: 288,
			Size:   "512x288",
			Aspect: 1.7777778,
		},
		Focus: &apimodel.MediaFocus{
			X: -0.5,
			Y: 0.5,
		},
	}, *attachmentReply.Meta)
	suite.Equal("LiB|W-#6RQR.~qvzRjWF_3rqV@a$", *attachmentReply.Blurhash)
	suite.NotEmpty(attachmentReply.ID)
	suite.NotEmpty(attachmentReply.URL)
	suite.NotEmpty(attachmentReply.PreviewURL)
	suite.Equal(len(storageKeysBeforeRequest)+2, len(storageKeysAfterRequest)) // 2 images should be added to storage: the original and the thumbnail
}

func (suite *MediaCreateTestSuite) TestMediaCreateLongDescription() {
	// set up the context for the request
	t := suite.testTokens["local_account_1"]
	oauthToken := oauth.DBTokenToToken(t)
	recorder := httptest.NewRecorder()

	// read a random string of a really long description
	descriptionBytes := make([]byte, 5000)
	if _, err := rand.Read(descriptionBytes); err != nil {
		panic(err)
	}
	description := base64.RawStdEncoding.EncodeToString(descriptionBytes)

	// create the request
	buf, w, err := testrig.CreateMultipartFormData(testrig.FileToDataF("file", "../../../../testrig/media/test-jpeg.jpg"), map[string][]string{
		"description": {description},
		"focus":       {"-0.5,0.5"},
	})
	if err != nil {
		panic(err)
	}
	req := httptest.NewRequest(http.MethodPost, "http://localhost:8080/api/v1/media", bytes.NewReader(buf.Bytes())) // the endpoint we're hitting
	req.Header.Set("Content-Type", w.FormDataContentType())
	req.Header.Set("accept", "application/json")

	c := httputil.ToContext(recorder, req)
	c.V.Set(oauth.SessionAuthorizedApplication, suite.testApplications["application_1"])
	c.V.Set(oauth.SessionAuthorizedToken, oauthToken)
	c.V.Set(oauth.SessionAuthorizedUser, suite.testUsers["local_account_1"])
	c.V.Set(oauth.SessionAuthorizedAccount, suite.testAccounts["local_account_1"])

	c.SetPathValue(apiutil.APIVersionKey, apiutil.APIv1)

	// do the actual request
	suite.mediaModule.MediaCreatePOSTHandler(c)

	// check response
	suite.EqualValues(http.StatusBadRequest, recorder.Code)

	result := recorder.Result()
	defer result.Body.Close()
	b, err := io.ReadAll(result.Body)
	suite.NoError(err)

	suite.Equal(`{"error":"Bad Request: image description length must be between 0 and 500 characters (inclusive), but provided image description was 6667 chars"}`, string(b))
}

func (suite *MediaCreateTestSuite) TestMediaCreateTooShortDescription() {
	// set the min description length
	config.SetMediaDescriptionMinChars(500)

	// set up the context for the request
	t := suite.testTokens["local_account_1"]
	oauthToken := oauth.DBTokenToToken(t)
	recorder := httptest.NewRecorder()

	// create the request
	buf, w, err := testrig.CreateMultipartFormData(testrig.FileToDataF("file", "../../../../testrig/media/test-jpeg.jpg"), map[string][]string{
		"description": {""}, // provide an empty description
		"focus":       {"-0.5,0.5"},
	})
	if err != nil {
		panic(err)
	}
	req := httptest.NewRequest(http.MethodPost, "http://localhost:8080/api/v1/media", bytes.NewReader(buf.Bytes())) // the endpoint we're hitting
	req.Header.Set("Content-Type", w.FormDataContentType())
	req.Header.Set("accept", "application/json")

	c := httputil.ToContext(recorder, req)
	c.V.Set(oauth.SessionAuthorizedApplication, suite.testApplications["application_1"])
	c.V.Set(oauth.SessionAuthorizedToken, oauthToken)
	c.V.Set(oauth.SessionAuthorizedUser, suite.testUsers["local_account_1"])
	c.V.Set(oauth.SessionAuthorizedAccount, suite.testAccounts["local_account_1"])

	c.SetPathValue(apiutil.APIVersionKey, apiutil.APIv1)

	// do the actual request
	suite.mediaModule.MediaCreatePOSTHandler(c)

	// check response -- there should be no error because minimum description length is checked on *UPDATE*, not initial upload
	suite.EqualValues(http.StatusOK, recorder.Code)
}

func TestMediaCreateTestSuite(t *testing.T) {
	suite.Run(t, new(MediaCreateTestSuite))
}
