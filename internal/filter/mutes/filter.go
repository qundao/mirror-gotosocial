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

package mutes

import (
	"time"

	"code.superseriousbusiness.org/gotosocial/internal/state"
)

type muteDetails struct {
	// mute flags.
	mute  bool
	notif bool

	// mute expiry times.
	muteExpiry  expiryTime
	notifExpiry expiryTime
}

// noauth is a placeholder ID used in cache lookups
// when there is no authorized account ID to use.
const noauth = "noauth"

// Filter packages up a bunch of logic for checking whether
// given statuses or accounts are muted by a requester (user).
type Filter struct{ state *state.State }

// NewFilter returns a new Filter interface that will use the provided state.
func NewFilter(state *state.State) *Filter { return &Filter{state: state} }

// expiryTime wraps a time.Time{}
// to also handle the case of zero
// value indicating "never expires",
// and tracking this as appropriate.
type expiryTime struct {
	time.Time
	never bool
}

// Update will update the expiryTime{} according
// to passed time, handling the case of an already
// set 'never' flag, newly setting it if necessary,
// or simply incrementing if 't' is after current.
func (e *expiryTime) Update(t time.Time) {
	if e.never {
		return
	}
	if e.Time.IsZero() {
		e.Never()
		return
	}
	switch {
	case e.never:
		// never expires

	case t.IsZero():
		// set time to
		// ever expire
		e.Never()

	case t.After(e.Time):
		// increment
		// expiry time
		e.Time = t
	}
}

// Never will set the 'never' flag
// in expiryTime{} to true preventing
// any further time incrementation.
func (e *expiryTime) Never() {
	e.Time = time.Time{}
	e.never = true
}
