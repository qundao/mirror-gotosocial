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

package cache

import (
	"time"

	"code.superseriousbusiness.org/gopkg/log"
	"code.superseriousbusiness.org/gotosocial/internal/cache/timeline"
)

type TimelineCaches struct {
	// Public provides an instance-level
	// cache of the public status timeline.
	Public timeline.StatusTimeline

	// Local provides an instance-level
	// cache of the local status timeline.
	Local timeline.StatusTimeline

	// Home provides a concurrency-safe map of status timeline
	// caches for home timelines, keyed by home's account ID.
	Home timeline.StatusTimelines

	// List provides a concurrency-safe map of status
	// timeline caches for lists, keyed by list ID.
	List timeline.StatusTimelines

	// Tag provides a concurrency-safe map of status
	// timeline caches for tags, keyed by tag ID.
	Tag timeline.StatusTimelines
}

func (c *Caches) initPublicTimeline() {
	// TODO: configurable
	cap := 800

	log.Infof(nil, "cache size = %d", cap)

	c.Timelines.Public.Init(cap)
}

func (c *Caches) initLocalTimeline() {
	// TODO: configurable
	cap := 800

	log.Infof(nil, "cache size = %d", cap)

	c.Timelines.Local.Init(cap)
}

func (c *Caches) initHomeTimelines() {
	// TODO: configurable
	timeout := 30 * time.Minute
	cap := 800

	log.Infof(nil, "cache size = %d", cap)

	c.Timelines.Home.Init(cap, timeout)
}

func (c *Caches) initListTimelines() {
	// TODO: configurable
	timeout := 30 * time.Minute
	cap := 800

	log.Infof(nil, "cache size = %d", cap)

	c.Timelines.List.Init(cap, timeout)
}

func (c *Caches) initTagTimelines() {
	// TODO: configurable
	timeout := 10 * time.Minute
	cap := 400

	log.Infof(nil, "cache size = %d", cap)

	c.Timelines.Tag.Init(cap, timeout)
}
