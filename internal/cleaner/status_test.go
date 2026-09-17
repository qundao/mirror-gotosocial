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

package cleaner_test

import (
	"testing"
	"time"

	"code.superseriousbusiness.org/gotosocial/internal/ap"
	"code.superseriousbusiness.org/gotosocial/internal/cleaner"
	"code.superseriousbusiness.org/gotosocial/internal/gtsmodel"
	"code.superseriousbusiness.org/gotosocial/internal/id"
	"code.superseriousbusiness.org/gotosocial/internal/state"
	"code.superseriousbusiness.org/gotosocial/testrig"
	"codeberg.org/gruf/go-longdur"
	"github.com/stretchr/testify/suite"
)

type StatusTestSuite struct {
	suite.Suite

	state        *state.State
	cleaner      *cleaner.Cleaner
	testStatuses map[string]*gtsmodel.Status
}

func TestStatusTestSuite(t *testing.T) {
	suite.Run(t, &StatusTestSuite{})
}

func (suite *StatusTestSuite) SetupTest() {
	testrig.InitTestConfig()
	testrig.InitTestLog()

	suite.state = new(state.State)
	testrig.StartNoopWorkers(suite.state)

	_ = testrig.NewTestDB(suite.state)
	testrig.StandardDBSetup(suite.state.DB, nil)

	suite.cleaner = cleaner.New(suite.state)

	suite.testStatuses = testrig.NewTestStatuses()
}

func (suite *StatusTestSuite) TearDownTest() {
	testrig.StandardDBTeardown(suite.state.DB)
	testrig.StopWorkers(suite.state)
}

func (suite *StatusTestSuite) TestPruneOldRemote() {
	suite.pruneOldRemote(4 * longdur.YearApprox)
	suite.pruneOldRemote(5 * longdur.YearApprox)
	suite.pruneOldRemote(6 * longdur.YearApprox)
}

func (suite *StatusTestSuite) pruneOldRemote(olderThan longdur.Duration) {
	ctx := suite.T().Context()
	now := time.Now()

	// Get older than as concrete time.
	_, dur := olderThan.Duration()
	olderThanTime := now.Add(-dur)

	var threadTimes []time.Time
	var statuses []*gtsmodel.Status

	// Generate a number of thread creation times
	// based on duration before olderThanTime.
	for _, d := range []longdur.Duration{
		longdur.Hour,
		longdur.Day,
		longdur.Week,
		longdur.MonthApprox,
		longdur.YearApprox,
	} {
		_, dur := d.Duration()
		time := olderThanTime.Add(-dur)
		threadTimes = append(threadTimes, time)
	}

	// Generate a number of remote statuses based on thread
	// times, which we actually expect to be later purged.
	statuses = suite.generateRemoteStatuses(olderThanTime, threadTimes)
	expect := len(statuses)

	// Generate a new group of remote statuses based on thread
	// times, which are also within age range to be later purged.
	statuses = suite.generateRemoteStatuses(olderThanTime, threadTimes)

	// Insert a boost at current date
	// for this second batch of statuses,
	// this should prevent being purged.
	for _, status := range statuses {

		// Generate status ID for current time.
		statusID := id.NewULIDFromTime(now)
		statusURI := "https://google.com/s/" + statusID

		// Use just some random account ID.
		accountID := "some_boostin_dude"
		accountURI := "https://google.com/u/" + accountID

		// Set boost of details from status.
		boostOfAccountID := status.AccountID
		boostOfID := status.ID

		// Create the status model.
		boost := &gtsmodel.Status{
			ID:  statusID,
			URI: statusURI,

			AccountID:  accountID,
			AccountURI: accountURI,

			BoostOfID:        boostOfID,
			BoostOfAccountID: boostOfAccountID,

			ThreadID: status.ThreadID,

			// This needs a proper fetched_at
			// setting so it gets appropriately
			// selected for its age.
			FetchedAt: now,

			ActivityStreamsType: ap.ObjectNote,

			Visibility: status.Visibility,
		}

		// Insert the new status in database,
		// specifically bypassing our typical
		// status threading logic confusing it.
		err := suite.state.DB.Put(ctx, boost)
		suite.NoError(err)
	}

	// Generate a new group of remote statuses based on thread
	// times, which are also within age range to be later purged.
	statuses = suite.generateRemoteStatuses(olderThanTime, threadTimes)

	// Insert a favourite by a local account for each status
	// into the database, this should prevent their purging.
	localAccountID := testrig.NewTestAccounts()["admin_account"].ID
	for _, status := range statuses {

		// Generate a new favourite ID.
		faveID := id.NewRandomULID()
		faveURI := "https://google.com/f/" + faveID

		// Create the favourite model.
		fave := &gtsmodel.StatusFave{
			ID:  faveID,
			URI: faveURI,

			AccountID:       localAccountID,
			TargetAccountID: status.AccountID,

			StatusID: status.ID,
		}

		// Insert the new favourite into the database.
		err := suite.state.DB.PutStatusFave(ctx, fave)
		suite.NoError(err)
	}

	// Generate a new group of remote statuses based on thread
	// times, which are also within age range to be later purged.
	statuses = suite.generateRemoteStatuses(olderThanTime, threadTimes)

	// Insert a bookmark by a local account for each status
	// into the database, this should prevent their purging.
	for _, status := range statuses {

		// Generate a new bookmark ID.
		bookmarkID := id.NewRandomULID()

		// Create the status bookmark model.
		bookmark := &gtsmodel.StatusBookmark{
			ID: bookmarkID,

			AccountID:       localAccountID,
			TargetAccountID: status.AccountID,

			StatusID: status.ID,
		}

		// Insert the new status bookmark into the database.
		err := suite.state.DB.PutStatusBookmark(ctx, bookmark)
		suite.NoError(err)
	}

	// Perform the status pruning for those older than time.
	count, err := suite.cleaner.Status().PruneOldRemote(ctx,
		olderThanTime, 0)
	suite.NoError(err)
	suite.Equal(expect, count)
}

func (suite *StatusTestSuite) generateRemoteStatuses(olderThan time.Time, threadTimes []time.Time) []*gtsmodel.Status {
	var statuses []*gtsmodel.Status
	ctx := suite.T().Context()

	// Generate a number of thread IDs
	// based on generated thread times.
	for _, t := range threadTimes {
		threadID := id.NewULIDFromTime(t)

		// Determine difference between
		// thread time and 'olderThan'.
		diff := olderThan.Sub(t)

		// Iterate between thread time and older than, creating
		// statuses of differing ages between these two points.
		for i := time.Duration(0); i < diff; i += (diff / 10) {

			// Generate new status ID from this time.
			statusID := id.NewULIDFromTime(t.Add(i))
			statusURI := "https://google.com/s/" + statusID

			// Use some random account ID.
			accountID := "some_random_dude"
			accountURI := "https://google.com/u/" + accountID

			// Create the status model.
			status := &gtsmodel.Status{
				ID:  statusID,
				URI: statusURI,

				AccountID:  accountID,
				AccountURI: accountURI,

				ThreadID: threadID,

				ActivityStreamsType: ap.ObjectNote,

				Visibility: gtsmodel.VisibilityDefault,
			}

			// Insert the new status in database,
			// specifically bypassing our typical
			// status threading logic confusing it.
			err := suite.state.DB.Put(ctx, status)
			suite.NoError(err)

			// Append status to the return slice.
			statuses = append(statuses, status)
		}
	}

	return statuses
}
