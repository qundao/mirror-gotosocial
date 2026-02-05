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

package testrig

import (
	"code.superseriousbusiness.org/gotosocial/internal/email"
	"code.superseriousbusiness.org/gotosocial/internal/filter/mutes"
	"code.superseriousbusiness.org/gotosocial/internal/filter/status"
	"code.superseriousbusiness.org/gotosocial/internal/filter/visibility"
	"code.superseriousbusiness.org/gotosocial/internal/processing/conversations"
	"code.superseriousbusiness.org/gotosocial/internal/processing/stream"
	"code.superseriousbusiness.org/gotosocial/internal/state"
	"code.superseriousbusiness.org/gotosocial/internal/surfacing"
	"code.superseriousbusiness.org/gotosocial/internal/typeutils"
	"code.superseriousbusiness.org/gotosocial/internal/util"
	"code.superseriousbusiness.org/gotosocial/internal/webpush"
)

func NewTestSurfacer(
	state *state.State,
	emailSender email.Sender,
	webPushSender webpush.Sender,
) *surfacing.Surfacer {
	converter := typeutils.NewConverter(state)
	visFilter := visibility.NewFilter(state)
	muteFilter := mutes.NewFilter(state)
	statusFilter := status.NewFilter(state)

	return surfacing.New(
		state,
		converter,
		util.Ptr(stream.New(state, NewTestOauthServer(state))),
		visFilter,
		muteFilter,
		statusFilter,
		emailSender,
		webPushSender,
		util.Ptr(conversations.New(state, converter, visFilter, muteFilter, statusFilter)),
	)
}
