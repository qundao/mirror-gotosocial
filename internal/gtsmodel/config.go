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

package gtsmodel

import (
	"time"

	"codeberg.org/gruf/go-longdur"
	"github.com/hashicorp/cronexpr"
	"github.com/uptrace/bun"
)

// CronExpression is a wrapper for cronexpr.Expression
// to allow parsing by CLI "flag"-like utilities.
type CronExpression struct {
	*cronexpr.Expression
	Expr string
}

func MustParseCron(expr string) (cron CronExpression) {
	if err := cron.Set(expr); err != nil {
		panic(err)
	}
	return
}

func (expr *CronExpression) Set(in string) (err error) {
	if in == "" {
		return
	}
	expr.Expr = in // set the raw expression string
	expr.Expression, err = cronexpr.Parse(in)
	return
}

func (expr *CronExpression) MarshalText() ([]byte, error) {
	return []byte(expr.Expr), nil
}

func (expr *CronExpression) UnmarshalText(text []byte) error {
	return expr.Set(string(text))
}

func (expr *CronExpression) String() string {
	return expr.Expr
}

// RuntimeConfig represents RUNTIME configuration of the GtS
// instance, ie., things that can safely be changed at runtime
// without requiring a restart of the entire service.
//
// This is in contrast to config in /internal/config which defines
// configuration that must be set BEFORE running the service,
// and which cannot be changed while the service is running.
type RuntimeConfig struct {
	bun.BaseModel `bun:"table:runtime_config"`

	// DB ID of the runtime config, always
	// set to 0 to ensure only 1 ever stored.
	ID uint8 `bun:",pk,notnull,default:0"`

	// Duration defining how long to
	// locally cache media from remote
	// instances (zero keeps indefinitely).
	//
	// Default 1 day.
	MediaRemoteCacheDuration longdur.Duration `bun:",notnull,default:6.048e+14"`

	// Default every night at midnight.
	MediaCleanupCron CronExpression `bun:",nullzero,notnull,default:0 0 * * *"`

	// Duration defining status
	// age beyond which to clean
	//
	// Default 0 (disabled).
	StatusesCleanupRemoteOlderThan longdur.Duration `bun:",notnull,default:0"`

	// Default every Sunday at 1am.
	StatusesCleanupCron CronExpression `bun:",nullzero,notnull,default:0 1 * * 0"`

	// Default every night at 11pm.
	InstanceSubscriptionsProcessCron CronExpression `bun:",nullzero,notnull,default:0 23 * * *"`
}

func (rc *RuntimeConfig) GetMediaRemoteCacheOlderThanTime(now time.Time) time.Time {
	_, dur := rc.MediaRemoteCacheDuration.Duration()
	return now.Add(-dur)
}

func (rc *RuntimeConfig) GetStatusesCleanupRemoteOlderThanTime(now time.Time) time.Time {
	_, dur := rc.StatusesCleanupRemoteOlderThan.Duration()
	return now.Add(-dur)
}
