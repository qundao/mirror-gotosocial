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
	"net/url"
	"testing"
	"time"

	"code.superseriousbusiness.org/gotosocial/internal/ap"
	"code.superseriousbusiness.org/gotosocial/internal/gtscontext"
	"code.superseriousbusiness.org/gotosocial/internal/gtsmodel"
	"code.superseriousbusiness.org/gotosocial/testrig"
	"github.com/stretchr/testify/suite"
)

type CreateRelayTestSuite struct {
	FederatingDBTestSuite
}

// TestCreateWrappedNoteToRelayActor ensures that a Create
// of a Note gets relayed when delivered to a relay actor.
func (suite *CreateRelayTestSuite) TestCreateWrappedNoteToRelayActor() {
	// Set up our test structs + tear down on finish.
	testStructs := testrig.SetupTestStructs(rMediaPath, rTemplatePath)
	defer testrig.TearDownTestStructs(testStructs)

	// Clean up test context when done.
	ctx, cncl := context.WithCancel(suite.T().Context())
	defer cncl()

	requester := suite.testAccounts["remote_account_2"]
	receiver := suite.testAccounts["relay_actor_1"]

	// Have the requester follow the relay actor
	// so the relay actor will process the note.
	if err := testStructs.State.DB.PutFollow(ctx,
		&gtsmodel.Follow{
			ID:              "01M16VMQ2GE4SZEDR0YR8CE5T1",
			URI:             "https://example.org/doesn't-matter",
			AccountID:       requester.ID,
			TargetAccountID: receiver.ID,
		},
	); err != nil {
		suite.FailNow(err.Error())
	}

	// Prepare a Note to send to the relay actor.
	status := testrig.NewTestStatuses()["remote_account_2_status_1"]
	requesterURI := testrig.URLMustParse(requester.URI)
	note := testrig.NewAPNote(&testrig.NewAPNoteParams{
		ID:           testrig.URLMustParse(status.URI),
		URL:          testrig.URLMustParse(status.URL),
		Content:      status.Content,
		Summary:      status.ContentWarning,
		AttributedTo: requesterURI,
		To:           []*url.URL{ap.PublicIRI()},
	})

	// Wrap the Note in a Create Activity.
	create := testrig.WrapAPNoteInCreate(
		testrig.URLMustParse(status.URI+"#Create"),
		requesterURI,
		time.Now(),
		note,
	)

	// Call Create.
	ctx = gtscontext.SetRequestingAccount(ctx, requester)
	ctx = gtscontext.SetReceivingAccount(ctx, receiver)
	err := testStructs.Federator.FederatingDB().Create(ctx, create)
	if err != nil {
		suite.FailNow(err.Error())
	}

	// Status should be relayed.
	if !testrig.WaitFor(func() bool {
		relayedURI, err := testStructs.State.DB.GetRelayedURI(ctx, receiver.URI, status.URI)
		return err == nil && relayedURI != nil
	}) {
		suite.FailNow("timed out waiting for relayedURI")
	}
}

// TestCreateNoteToRelayActor is like TestCreateWrappedNoteToRelayActor
// but it passes the Note object into the Create function directly.
func (suite *CreateRelayTestSuite) TestCreateNoteToRelayActor() {
	// Set up our test structs + tear down on finish.
	testStructs := testrig.SetupTestStructs(rMediaPath, rTemplatePath)
	defer testrig.TearDownTestStructs(testStructs)

	// Clean up test context when done.
	ctx, cncl := context.WithCancel(suite.T().Context())
	defer cncl()

	requester := suite.testAccounts["remote_account_2"]
	receiver := suite.testAccounts["relay_actor_1"]

	// Have the requester follow the relay actor
	// so the relay actor will process the note.
	if err := testStructs.State.DB.PutFollow(ctx,
		&gtsmodel.Follow{
			ID:              "01M16VMQ2GE4SZEDR0YR8CE5T1",
			URI:             "https://example.org/doesn't-matter",
			AccountID:       requester.ID,
			TargetAccountID: receiver.ID,
		},
	); err != nil {
		suite.FailNow(err.Error())
	}

	// Prepare a Note to send to the relay actor.
	status := testrig.NewTestStatuses()["remote_account_2_status_1"]
	requesterURI := testrig.URLMustParse(requester.URI)
	note := testrig.NewAPNote(&testrig.NewAPNoteParams{
		ID:           testrig.URLMustParse(status.URI),
		URL:          testrig.URLMustParse(status.URL),
		Content:      status.Content,
		Summary:      status.ContentWarning,
		AttributedTo: requesterURI,
		To:           []*url.URL{ap.PublicIRI()},
	})

	// Call Create.
	ctx = gtscontext.SetRequestingAccount(ctx, requester)
	ctx = gtscontext.SetReceivingAccount(ctx, receiver)
	err := testStructs.Federator.FederatingDB().Create(ctx, note)
	if err != nil {
		suite.FailNow(err.Error())
	}

	// Status should be relayed.
	if !testrig.WaitFor(func() bool {
		relayedURI, err := testStructs.State.DB.GetRelayedURI(ctx, receiver.URI, status.URI)
		return err == nil && relayedURI != nil
	}) {
		suite.FailNow("timed out waiting for relayedURI")
	}
}

// TestCreateNoteToRelayActorNoMatch ensures that a relay
// actor doesn't relay a status it doesn't match against.
func (suite *CreateRelayTestSuite) TestCreateNoteToRelayActorNoMatch() {
	// Set up our test structs + tear down on finish.
	testStructs := testrig.SetupTestStructs(rMediaPath, rTemplatePath)
	defer testrig.TearDownTestStructs(testStructs)

	// Clean up test context when done.
	ctx, cncl := context.WithCancel(suite.T().Context())
	defer cncl()

	// Update the relay actor to not match all by default.
	relayActor := new(gtsmodel.RelayActor)
	*relayActor = *suite.testRelayActors["relay_actor_1"]
	relayActor.Flags.SetMatchByDefault(false)
	if err := testStructs.State.DB.UpdateRelayActor(ctx,
		relayActor, "flags",
	); err != nil {
		suite.FailNow(err.Error())
	}

	requester := suite.testAccounts["remote_account_2"]
	receiver := suite.testAccounts["relay_actor_1"]

	// Have the requester follow the relay actor
	// so the relay actor will process the note.
	if err := testStructs.State.DB.PutFollow(ctx,
		&gtsmodel.Follow{
			ID:              "01M16VMQ2GE4SZEDR0YR8CE5T1",
			URI:             "https://example.org/doesn't-matter",
			AccountID:       requester.ID,
			TargetAccountID: receiver.ID,
		},
	); err != nil {
		suite.FailNow(err.Error())
	}

	// Prepare a Note to send to the relay actor.
	status := testrig.NewTestStatuses()["remote_account_2_status_1"]
	requesterURI := testrig.URLMustParse(requester.URI)
	note := testrig.NewAPNote(&testrig.NewAPNoteParams{
		ID:           testrig.URLMustParse(status.URI),
		URL:          testrig.URLMustParse(status.URL),
		Content:      status.Content,
		Summary:      status.ContentWarning,
		AttributedTo: requesterURI,
		To:           []*url.URL{ap.PublicIRI()},
	})

	// Wrap the Note in a Create Activity.
	create := testrig.WrapAPNoteInCreate(
		testrig.URLMustParse(status.URI+"#Create"),
		requesterURI,
		time.Now(),
		note,
	)

	// Call Create.
	ctx = gtscontext.SetRequestingAccount(ctx, requester)
	ctx = gtscontext.SetReceivingAccount(ctx, receiver)
	err := testStructs.Federator.FederatingDB().Create(ctx, create)
	if err != nil {
		suite.FailNow(err.Error())
	}

	// Status should NOT be relayed.
	if testrig.WaitFor(func() bool {
		relayedURI, err := testStructs.State.DB.GetRelayedURI(ctx, receiver.URI, status.URI)
		return err == nil && relayedURI != nil
	}) {
		suite.FailNow("note should not have been relayed")
	}
}

func TestCreateRelayTestSuite(t *testing.T) {
	suite.Run(t, new(CreateRelayTestSuite))
}
