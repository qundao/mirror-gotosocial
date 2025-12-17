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
	"context"
	"errors"
	"time"

	"code.superseriousbusiness.org/gotosocial/internal/db"
	"code.superseriousbusiness.org/gotosocial/internal/gtscontext"
	"code.superseriousbusiness.org/gotosocial/internal/gtserror"
)

// NOTE:
// we don't bother using the Mutes cache for any
// of the accounts functions below, as there's only
// a single cache load required of any UserMute.

// AccountNotificationsMuted returns whether notifications
// from target account are muted for requesting account.
func (f *Filter) AccountNotificationsMuted(
	ctx context.Context,
	requesterID string,
	targetID string,
) (bool, error) {
	// Look for mute against target.
	mute, err := f.state.DB.GetMute(
		gtscontext.SetBarebones(ctx),
		requesterID,
		targetID,
	)
	if err != nil && !errors.Is(err, db.ErrNoEntries) {
		return false, gtserror.Newf("db error getting user mute: %w", err)
	}

	if mute == nil {
		// No user mute exists!
		return false, nil
	}

	// To avoid calling time.Now(),
	// return early if this mute
	// doesn't apply to notifs.
	if !*mute.Notifications {
		return false, nil
	}

	// This mute applies to notifs.
	// If mute doesn't expire then
	// notifs are definitely muted.
	if mute.ExpiresAt.IsZero() {
		return true, nil
	}

	// The mute applies to notifs
	// and may expire. Only return
	// true if it's not expired.
	expired := time.Now().After(mute.ExpiresAt)
	return !expired, nil
}
