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

package admin

import (
	"context"

	apimodel "code.superseriousbusiness.org/gotosocial/internal/api/model"
	"code.superseriousbusiness.org/gotosocial/internal/gtserror"
)

func (p *Processor) RuntimeConfigGet(ctx context.Context) *apimodel.AdminRuntimeConfig {
	rtConf := p.state.DB.RuntimeConfig(ctx)
	return p.converter.RuntimeConfigToAPIRuntimeConfig(rtConf)
}

func (p *Processor) RuntimeConfigUpdate(
	ctx context.Context,
	form *apimodel.AdminRuntimeConfig,
) (*apimodel.AdminRuntimeConfig, gtserror.WithCode) {
	// Fetch config to update it.
	rtConf := p.state.DB.RuntimeConfig(ctx)

	// DB columns to update.
	columns := make([]string, 0, 5)

	// True if config changes necessitate
	// rescheduling instance blocklist
	// subscription job.
	var rescheduleSubscriptions bool

	// True if config changes necessitate
	// rescheduling media cleanup job.
	var rescheduleCleaner bool

	if v := form.InstanceSubscriptionsProcessCron; v != nil {
		// InstanceSubscriptionsProcessCron was provided on form.
		const col = "instance_subscriptions_process_cron"

		// Only update if expression changed.
		if *v != rtConf.InstanceSubscriptionsProcessCron.Expr {
			err := rtConf.InstanceSubscriptionsProcessCron.Set(*v)
			if err != nil {
				err := gtserror.Newf("bad value for %s: %w", col, err)
				return nil, gtserror.NewErrorBadRequest(err)
			}
			columns = append(columns, col)

			// Must reschedule
			// subscriptions for this.
			rescheduleSubscriptions = true
		}
	}

	if v := form.MediaCleanupCron; v != nil {
		// MediaCleanupCron was provided on form.
		const col = "media_cleanup_cron"

		// Only update if expression changed.
		if *v != rtConf.MediaCleanupCron.Expr {
			err := rtConf.MediaCleanupCron.Set(*v)
			if err != nil {
				err := gtserror.Newf("bad value for %s: %w", col, err)
				return nil, gtserror.NewErrorBadRequest(err)
			}
			columns = append(columns, col)

			// Must reschedule
			// cleaning for this.
			rescheduleCleaner = true
		}
	}

	if v := form.MediaRemoteCacheDuration; v != nil {
		// MediaRemoteCacheDuration was provided on form.
		const col = "media_remote_cache_duration"

		// Only update if duration changed.
		if v.Duration != rtConf.MediaRemoteCacheDuration {
			rtConf.MediaRemoteCacheDuration = v.Duration
			columns = append(columns, col)
		}
	}

	if v := form.StatusesCleanupCron; v != nil {
		// StatusesCleanupCron was provided on form.
		const col = "statuses_cleanup_cron"

		// Only update if expression changed.
		if *v != rtConf.StatusesCleanupCron.Expr {
			err := rtConf.StatusesCleanupCron.Set(*v)
			if err != nil {
				err := gtserror.Newf("bad value for %s: %w", col, err)
				return nil, gtserror.NewErrorBadRequest(err)
			}
			columns = append(columns, col)

			// Must reschedule cleaning for this,
			// as scheduling of statuses cleanup
			// depends on duration not being <= 0.
			rescheduleCleaner = true
		}
	}

	if v := form.StatusesCleanupRemoteOlderThan; v != nil {
		// StatusesCleanupRemoteOlderThan was provided on form.
		const col = "statuses_cleanup_remote_older_than"

		// Only update if duration changed.
		if v.Duration != rtConf.StatusesCleanupRemoteOlderThan {
			rtConf.StatusesCleanupRemoteOlderThan = v.Duration
			columns = append(columns, col)

			// Must reschedule
			// cleaning for this.
			rescheduleCleaner = true
		}
	}

	// Only update/reschedule if
	// a value actually changed.
	if len(columns) != 0 {
		// Update config in the db.
		if err := p.state.DB.UpdateRuntimeConfig(ctx, rtConf, columns...); err != nil {
			err := gtserror.Newf("db error updating runtime config: %w", err)
			return nil, gtserror.NewErrorInternalError(err)
		}

		// Reschedule instance blocklist
		// subscriptions if necessary.
		if rescheduleSubscriptions {
			if err := p.subscriptions.ScheduleJobs(ctx); err != nil {
				err := gtserror.Newf("error rescheduling blocklist subscription: %w", err)
				return nil, gtserror.NewErrorInternalError(err, err.Error())
			}
		}

		// Reschedule media
		// cleanup if necessary.
		if rescheduleCleaner {
			if err := p.cleaner.ScheduleJobs(ctx); err != nil {
				err := gtserror.Newf("error rescheduling media cleanup: %w", err)
				return nil, gtserror.NewErrorInternalError(err, err.Error())
			}
		}
	}

	return p.converter.RuntimeConfigToAPIRuntimeConfig(rtConf), nil
}
