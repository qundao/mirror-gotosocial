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

package bundb_test

import (
	"testing"

	"code.superseriousbusiness.org/gotosocial/internal/gtsmodel"
	"github.com/stretchr/testify/suite"
)

type RuntimeConfigTestSuite struct {
	BunDBStandardTestSuite
}

func (suite *RuntimeConfigTestSuite) TestCreateGetUpdateRuntimeConfig() {
	ctx := suite.T().Context()

	// Make sure this doesn't panic.
	var rtConf *gtsmodel.RuntimeConfig
	if !suite.NotPanics(func() {
		rtConf = suite.state.DB.RuntimeConfig(ctx)
	}, "panicked!") {
		return
	}

	// Check database defaults are set.
	suite.EqualValues(1, rtConf.ID)
	suite.Equal("1 day", rtConf.MediaRemoteCacheDuration.String())
	suite.Equal("0 0 * * *", rtConf.MediaCleanupCron.Expr)
	suite.Equal("0 sec", rtConf.StatusesCleanupRemoteOlderThan.String())
	suite.Equal("0 1 * * 0", rtConf.StatusesCleanupCron.Expr)
	suite.Equal("0 23 * * *", rtConf.InstanceSubscriptionsProcessCron.Expr)

	// Update statuses cleanup to every Monday instead of every Sunday.
	if err := rtConf.StatusesCleanupCron.Set("0 2 * * 0"); err != nil {
		suite.FailNow(err.Error())
	}

	// Set statuses cleanup remote older than to 1 year.
	if err := rtConf.StatusesCleanupRemoteOlderThan.Set("1 year"); err != nil {
		suite.FailNow(err.Error())
	}

	// Update the config in the db.
	if err := suite.state.DB.UpdateRuntimeConfig(ctx,
		rtConf,
		"statuses_cleanup_cron",
		"statuses_cleanup_remote_older_than",
	); err != nil {
		suite.FailNow(err.Error())
	}
}

func TestRuntimeConfigTestSuite(t *testing.T) {
	suite.Run(t, new(RuntimeConfigTestSuite))
}
