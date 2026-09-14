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
	var runtimeConfig *gtsmodel.RuntimeConfig
	suite.NotPanics(func() {
		runtimeConfig = suite.state.DB.RuntimeConfig(ctx)
	}, "panicked!")

	// Check database defaults are set.
	suite.EqualValues(0, int(runtimeConfig.ID))
	suite.Equal("1 week", runtimeConfig.MediaRemoteCacheDuration.String())
}

func TestRuntimeConfigTestSuite(t *testing.T) {
	suite.Run(t, new(RuntimeConfigTestSuite))
}
