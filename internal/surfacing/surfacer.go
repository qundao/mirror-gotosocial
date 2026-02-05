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

package surfacing

import (
	"code.superseriousbusiness.org/gotosocial/internal/email"
	"code.superseriousbusiness.org/gotosocial/internal/filter/mutes"
	"code.superseriousbusiness.org/gotosocial/internal/filter/status"
	"code.superseriousbusiness.org/gotosocial/internal/filter/visibility"
	"code.superseriousbusiness.org/gotosocial/internal/processing/conversations"
	"code.superseriousbusiness.org/gotosocial/internal/processing/stream"
	"code.superseriousbusiness.org/gotosocial/internal/state"
	"code.superseriousbusiness.org/gotosocial/internal/typeutils"
	"code.superseriousbusiness.org/gotosocial/internal/webpush"
)

// Surfacer wraps functions for 'surfacing' the result
// of ingesting a message into the server, eg:
//   - timelining a status
//   - removing a status from timelines
//   - sending a notification to a user
//   - sending an email
type Surfacer struct {
	state         *state.State
	converter     *typeutils.Converter
	stream        *stream.Processor
	visFilter     *visibility.Filter
	muteFilter    *mutes.Filter
	statusFilter  *status.Filter
	emailSender   email.Sender
	webPushSender webpush.Sender
	conversations *conversations.Processor
}

// New returns a pointer
// to a new surfacer struct.
func New(
	state *state.State,
	converter *typeutils.Converter,
	stream *stream.Processor,
	visFilter *visibility.Filter,
	muteFilter *mutes.Filter,
	statusFilter *status.Filter,
	emailSender email.Sender,
	webPushSender webpush.Sender,
	conversations *conversations.Processor,
) *Surfacer {
	return &Surfacer{
		state:         state,
		converter:     converter,
		stream:        stream,
		visFilter:     visFilter,
		muteFilter:    muteFilter,
		statusFilter:  statusFilter,
		emailSender:   emailSender,
		webPushSender: webPushSender,
		conversations: conversations,
	}
}
