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

package statuses

import (
	"context"
	"errors"
	"fmt"
	"time"

	"code.superseriousbusiness.org/gopkg/log"
	"code.superseriousbusiness.org/gotosocial/cmd/gotosocial/action"
	"code.superseriousbusiness.org/gotosocial/internal/cleaner"
	"code.superseriousbusiness.org/gotosocial/internal/db/bundb"
	"code.superseriousbusiness.org/gotosocial/internal/state"
)

// check function conformance.
var _ action.GTSAction = PruneLeafStubs
var _ action.GTSAction = PruneOldRemote

func PruneLeafStubs(ctx context.Context) error {
	return do(ctx, func(p *pruner) error {
		olderThan := time.Now().Add(-7 * 24 * time.Hour)
		p.cleaner.Status().LogPruneLeafStubs(ctx, olderThan, 0)
		return nil
	})
}

func PruneOldRemote(ctx context.Context) error {
	return do(ctx, func(p *pruner) error {
		rtConf := p.state.DB.RuntimeConfig(ctx)
		_, dur := rtConf.StatusesCleanupRemoteOlderThan.Duration()
		if dur == 0 {
			return errors.New("StatusesCleanupRemoteOlderThan = 0, no statuses to cleanup")
		}
		olderThan := time.Now().Add(-dur)
		p.cleaner.Status().LogPruneOldRemote(ctx, olderThan, 0)
		return nil
	})
}

type pruner struct {
	state   *state.State
	cleaner *cleaner.Cleaner
}

func do(ctx context.Context, do func(*pruner) error) error {
	var state state.State

	// Setup logging
	// for CLI output.
	log.SetOutputCLI()

	defer func() {
		if state.DB != nil {
			// Lastly, if database service was started,
			// ensure it gets closed now all else stopped.
			if err := state.DB.Close(); err != nil {
				log.Errorf(ctx, "error stopping database: %v", err)
			}
		}

		// Finally reached end of shutdown.
		log.Info(ctx, "done! exiting...")
	}()

	// Initialize caches
	state.Caches.Init()
	err := state.Caches.Start()
	if err != nil {
		return fmt.Errorf("error starting caches: %w", err)
	}

	log.Info(ctx, "starting db service...")

	// Open conn to database now caches are started.
	state.DB, err = bundb.NewBunDBService(ctx, &state)
	if err != nil {
		return fmt.Errorf("error creating dbservice: %w", err)
	}

	// Create common receiver type.
	cleaner := cleaner.New(&state)
	pruner := &pruner{&state, cleaner}

	// Perform provided CLI function.
	if err := do(pruner); err != nil {
		return fmt.Errorf("error performing action: %w", err)
	}

	return nil
}
