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

package cleaner

import (
	"context"
	"time"
	"unsafe"

	"code.superseriousbusiness.org/gopkg/log"
	"code.superseriousbusiness.org/gotosocial/internal/gtscontext"
	"code.superseriousbusiness.org/gotosocial/internal/gtserror"
	"code.superseriousbusiness.org/gotosocial/internal/state"
	"code.superseriousbusiness.org/gotosocial/internal/storage"
	"codeberg.org/gruf/go-longdur"
)

const selectLimit = 50
const stamp = "Jan _2 2006 15:04:05"

type Cleaner struct{ state *state.State }

func New(state *state.State) *Cleaner {
	c := new(Cleaner)
	c.state = state
	return c
}

// Emoji returns the emoji set of cleaner utilities.
func (c *Cleaner) Emoji() *Emoji {
	if unsafe.Sizeof(Emoji{}) != unsafe.Sizeof(Cleaner{}) ||
		unsafe.Offsetof(Emoji{}.Cleaner) != 0 {
		panic(gtserror.New("compile time unsafe pointer assertion"))
	}
	return (*Emoji)(unsafe.Pointer(c))
}

// Media returns the media set of cleaner utilities.
func (c *Cleaner) Media() *Media {
	if unsafe.Sizeof(Media{}) != unsafe.Sizeof(Cleaner{}) ||
		unsafe.Offsetof(Media{}.Cleaner) != 0 {
		panic(gtserror.New("compile time unsafe pointer assertion"))
	}
	return (*Media)(unsafe.Pointer(c))
}

// Status returns the status set of cleaner utilities.
func (c *Cleaner) Status() *Status {
	if unsafe.Sizeof(Status{}) != unsafe.Sizeof(Cleaner{}) ||
		unsafe.Offsetof(Status{}.Cleaner) != 0 {
		panic(gtserror.New("compile time unsafe pointer assertion"))
	}
	return (*Status)(unsafe.Pointer(c))
}

// haveFiles returns whether all of the provided files exist within current storage.
func (c *Cleaner) haveFiles(ctx context.Context, files ...string) (bool, error) {
	for _, path := range files {
		if path == "" {
			// File not stored.
			return false, nil
		}

		// Check whether each file exists in storage.
		have, err := c.state.Storage.Has(ctx, path)
		if err != nil {
			return false, gtserror.Newf("error checking storage for %s: %w", path, err)
		}

		if !have {
			// Missing file(s).
			return false, nil
		}
	}
	return true, nil
}

// removeFiles removes the provided files, returning the number of them returned.
func (c *Cleaner) removeFiles(ctx context.Context, files ...string) {
	if gtscontext.DryRun(ctx) {
		// Dry run,
		// do nothing.
		return
	}

	for _, path := range files {
		if path == "" {
			// not stored.
			continue
		}

		// Remove each provided storage path.
		log.Debugf(ctx, "removing file: %s", path)
		err := c.state.Storage.Delete(ctx, path)
		if err != nil && !storage.IsNotFound(err) {
			log.Errorf(ctx, "error removing %s: %v", path, err)
			continue
		}
	}
}

// ScheduleJobs schedules cleaning
// jobs using configured parameters.
func (c *Cleaner) ScheduleJobs(ctx context.Context) error {
	// Schedule media cleanup.
	rtConf := c.state.DB.RuntimeConfig(ctx)
	expr := rtConf.MediaCleanupCron

	// Cancel any existing job.
	const mcTaskID = "@mediacleanup"
	c.state.Workers.Scheduler.Cancel(mcTaskID)

	// Schedule new job.
	log.Infof(ctx, "scheduling media cleanup: %s", expr.Expr)
	if !c.state.Workers.Scheduler.Add(mcTaskID, c.cleanMedia, expr) {
		return gtserror.New("failed to schedule " + mcTaskID)
	}

	// Schedule statuses cleanup.
	// Cancel any existing job.
	const scTaskID = "@statuscleanup"
	c.state.Workers.Scheduler.Cancel(scTaskID)
	if _, dur := rtConf.StatusesCleanupRemoteOlderThan.Duration(); dur <= 0 {
		log.Infof(ctx, "skipping statuses cleanup scheduling, as statuses-cleanup-remote-older-than <= 0")
		return nil
	}

	// Schedule new job.
	expr = rtConf.StatusesCleanupCron
	log.Infof(ctx, "scheduling statuses cleanup: %s", expr.Expr)
	if !c.state.Workers.Scheduler.Add(scTaskID, c.cleanStatuses, expr) {
		return gtserror.New("failed to schedule " + scTaskID)
	}

	return nil
}

func (c *Cleaner) cleanMedia(ctx context.Context, start time.Time) {
	log.Info(ctx, "starting")
	rtConf := c.state.DB.RuntimeConfig(ctx)
	c.Media().All(ctx, start, rtConf.MediaRemoteCacheDuration)
	c.Emoji().All(ctx, start, rtConf.MediaRemoteCacheDuration)
	log.Infof(ctx, "finished after %s", time.Since(start))
}

func (c *Cleaner) cleanStatuses(ctx context.Context, start time.Time) {
	log.Info(ctx, "starting")
	rtConf := c.state.DB.RuntimeConfig(ctx)
	maxRemoteAge := rtConf.StatusesCleanupRemoteOlderThan
	c.Status().All(ctx, start, 7*longdur.Day, maxRemoteAge)
	log.Infof(ctx, "finished after %s", time.Since(start))
}
