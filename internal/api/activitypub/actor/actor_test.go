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

package actor_test

import (
	"code.superseriousbusiness.org/gopkg/httputil"
	"code.superseriousbusiness.org/gotosocial/internal/admin"
	"code.superseriousbusiness.org/gotosocial/internal/api/activitypub/actor"
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

type ActorStandardTestSuite struct {
	// standard suite interfaces
	suite.Suite
	db           db.DB
	tc           *typeutils.Converter
	mediaManager *media.Manager
	federator    *federation.Federator
	emailSender  email.Sender
	processor    *processing.Processor
	storage      *storage.Driver
	state        *state.State

	// standard suite models
	testTokens       map[string]*gtsmodel.Token
	testApplications map[string]*gtsmodel.Application
	testUsers        map[string]*gtsmodel.User
	testAccounts     map[string]*gtsmodel.Account
	testAttachments  map[string]*gtsmodel.MediaAttachment
	testStatuses     map[string]*gtsmodel.Status
	testBlocks       map[string]*gtsmodel.Block
	testActivities   map[string]testrig.ActivityWithSignature

	// module being tested
	actorModule *actor.Module

	signatureCheck httputil.Middleware
}

func (suite *ActorStandardTestSuite) SetupSuite() {
	suite.testTokens = testrig.NewTestTokens()
	suite.testApplications = testrig.NewTestApplications()
	suite.testUsers = testrig.NewTestUsers()
	suite.testAccounts = testrig.NewTestAccounts()
	suite.testAttachments = testrig.NewTestAttachments()
	suite.testStatuses = testrig.NewTestStatuses()
	suite.testBlocks = testrig.NewTestBlocks()
	suite.testActivities = testrig.NewTestActivities(suite.testAccounts)
}

func (suite *ActorStandardTestSuite) SetupTest() {
	suite.state = new(state.State)

	testrig.InitTestConfig()
	testrig.InitTestLog()

	suite.db = testrig.NewTestDB(suite.state)
	suite.state.DB = suite.db
	suite.state.AdminActions = admin.New(suite.state.DB, &suite.state.Workers)
	suite.tc = typeutils.NewConverter(suite.state)

	suite.storage = testrig.NewInMemoryStorage()
	suite.state.Storage = suite.storage
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
	testrig.StartWorkers(suite.state, suite.processor.Workers())

	suite.actorModule = actor.New(suite.processor, testrig.LoadTemplates(suite.state, ""))
	testrig.StandardDBSetup(suite.db, suite.testAccounts)
	testrig.StandardStorageSetup(suite.storage, "../../../../testrig/media")

	suite.signatureCheck = middleware.ExtractSignature(suite.db.IsURIBlocked)
}

func (suite *ActorStandardTestSuite) TearDownTest() {
	testrig.StandardDBTeardown(suite.db)
	testrig.StandardStorageTeardown(suite.storage)
	testrig.StopWorkers(suite.state)
}
