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

package media

import (
	"context"
	"time"

	"code.superseriousbusiness.org/gopkg/log"
	"code.superseriousbusiness.org/gotosocial/cmd/gotosocial/action"
	"code.superseriousbusiness.org/gotosocial/internal/config"
	"code.superseriousbusiness.org/gotosocial/internal/gtscontext"
)

// check function conformance.
var _ action.GTSAction = PruneAll

// PruneAll performs all media clean actions
func PruneAll(ctx context.Context) error {

	// Setup pruning utilities.
	prune, err := setupPrune(ctx)
	if err != nil {
		return err
	}

	defer func() {
		// Ensure pruner gets shutdown on exit.
		if err := prune.shutdown(); err != nil {
			log.Error(ctx, err)
		}
	}()

	if config.GetAdminMediaPruneDryRun() {
		log.Info(ctx, "prune DRY RUN")
		ctx = gtscontext.SetDryRun(ctx)
	}

	// Current time.
	now := time.Now()

	// Get media maximum remote cache duration.
	rtConf := prune.state.DB.RuntimeConfig(ctx)
	dur := rtConf.MediaRemoteCacheDuration

	// Perform the actual pruning with logging.
	prune.cleaner.Media().AllAndFix(ctx, now, dur)
	prune.cleaner.Emoji().AllAndFix(ctx, now, dur)

	return nil
}
