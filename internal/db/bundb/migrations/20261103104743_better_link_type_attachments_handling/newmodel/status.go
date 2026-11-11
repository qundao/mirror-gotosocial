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

import "time"

type Status struct {
	ID                       string    `bun:"type:CHAR(26),pk,nullzero,notnull,unique"`
	CreatedAt                time.Time `bun:"type:timestamptz,nullzero,notnull,default:current_timestamp"`
	EditedAt                 time.Time `bun:"type:timestamptz,nullzero"`
	FetchedAt                time.Time `bun:"type:timestamptz,nullzero"`
	URI                      string    `bun:",unique,nullzero,notnull"`
	URL                      string    `bun:",nullzero"`
	Content                  string    `bun:""`
	AttachmentIDs            []string  `bun:"attachments,array"`
	LinkIDs                  []LinkID  `bun:"links"` // <----- added in this migration
	TagIDs                   []string  `bun:"tags,array"`
	MentionIDs               []string  `bun:"mentions,array"`
	EmojiIDs                 []string  `bun:"emojis,array"`
	AccountID                string    `bun:"type:CHAR(26),nullzero,notnull"`
	AccountURI               string    `bun:",nullzero,notnull"`
	InReplyToID              string    `bun:"type:CHAR(26),nullzero"`
	InReplyToURI             string    `bun:",nullzero"`
	InReplyToAccountID       string    `bun:"type:CHAR(26),nullzero"`
	BoostOfID                string    `bun:"type:CHAR(26),nullzero"`
	BoostOfAccountID         string    `bun:"type:CHAR(26),nullzero"`
	ThreadID                 string    `bun:"type:CHAR(26),nullzero,notnull,default:'00000000000000000000000000'"`
	EditIDs                  []string  `bun:"edits,array"`
	PollID                   string    `bun:"type:CHAR(26),nullzero"`
	ContentWarning           string    `bun:",nullzero"`
	ContentWarningText       string    `bun:""`
	Flags                    int16     `bun:",notnull,default:0"`
	Visibility               int16     `bun:",nullzero,notnull"`
	Language                 string    `bun:",nullzero"`
	CreatedWithApplicationID string    `bun:"type:CHAR(26),nullzero"`
	ActivityStreamsType      string    `bun:",nullzero,notnull"`
	Text                     string    `bun:""`
	ContentType              int16     `bun:",nullzero"`
	InteractionPolicy        *struct{} `bun:""`
	ApprovedByURI            string    `bun:",nullzero"`
}

// StatusToLink is an intermediate struct to facilitate the
// many2many relationship between a status and one or more links.
type StatusToLink struct {
	StatusID string  `bun:"type:CHAR(26),unique:statuslink,nullzero,notnull"`
	Status   *Status `bun:"rel:belongs-to"`
	LinkID   string  `bun:"type:CHAR(26),unique:statuslink,nullzero,notnull"`
	Link     *Link   `bun:"rel:belongs-to"`
}

// LinkID encapsulates a link ID alongside an
// indication of whether the link was attached
// to the status or in the status body.
type LinkID struct {
	// ID of a gtsmodel.Link
	// entry in the database.
	ID string `json:"id"`

	// 1 if link was included as separate
	// AP `attachment` rather than being
	// parsed from the body text, else 0.
	Attached uint8 `json:"a,omitempty"`
}

// InBody returns true if link was in body
// text rather than separate attachment.
func (l LinkID) InBody() bool {
	return l.Attached == 0
}
