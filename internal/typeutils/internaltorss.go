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

package typeutils

import (
	"context"
	"strconv"
	"strings"

	"code.superseriousbusiness.org/gotosocial/internal/config"
	"code.superseriousbusiness.org/gotosocial/internal/gtserror"
	"code.superseriousbusiness.org/gotosocial/internal/gtsmodel"
	"code.superseriousbusiness.org/gotosocial/internal/text"
	"github.com/gorilla/feeds"
)

const (
	rssTitleMaxRunes       = 128
	rssDescriptionMaxRunes = 256
)

// see https://cyber.harvard.edu/rss/rss.html
func (c *Converter) StatusToRSSItem(ctx context.Context, status *gtsmodel.Status) (*feeds.Item, error) {
	var err error

	// Ensure account populated.
	if status.Account == nil {
		status.Account, err = c.state.DB.GetAccountByID(ctx, status.AccountID)
		if err != nil {
			return nil, gtserror.Newf("db error getting status author: %w", err)
		}
	}

	// Get first attachment if present.
	var media0 *gtsmodel.MediaAttachment
	if status.AttachmentsPopulated() && len(status.Attachments) > 0 {
		media0 = status.Attachments[0]
	} else if len(status.AttachmentIDs) > 0 {
		media0, err = c.state.DB.GetAttachmentByID(ctx, status.AttachmentIDs[0])
		if err != nil {
			return nil, gtserror.Newf("db error getting status attachment: %w", err)
		}
	}

	// Title -- The title of the item.
	// example: Venice Film Festival Tries to Quit Sinking
	var title string
	if status.ContentWarning != "" {
		title = trimTo(status.ContentWarning, rssTitleMaxRunes)
	} else {
		title = trimTo(status.Text, rssTitleMaxRunes)
	}

	// Generate author name string for status.
	authorName := "@" + status.Account.Username +
		"@" + config.GetAccountDomain()

	var buf strings.Builder
	buf.Grow(512)

	// Description -- The item synopsis.
	// example: Some of the most heated chatter at the Venice Film Festival this week was
	// about the way that the arrival of the stars at the Palazzo del Cinema was being staged.
	buf.WriteString(authorName + " ")
	switch l := len(status.AttachmentIDs); {
	case l > 1:
		buf.WriteString("posted [")
		buf.WriteString(strconv.Itoa(l))
		buf.WriteString("] attachments")
	case l == 1:
		buf.WriteString("posted 1 attachment")
	default:
		buf.WriteString("made a new post")
	}
	if status.Text != "" {
		buf.WriteString(": \"")
		buf.WriteString(status.Text)
		buf.WriteString("\"")
	}
	description := trimTo(buf.String(), rssDescriptionMaxRunes)

	// Enclosure, describes a media object
	// that is attached to the item.
	var enclosure *feeds.Enclosure

	// Set media details.
	if media0 != nil {
		enclosure = new(feeds.Enclosure)
		enclosure.Type = media0.File.ContentType
		enclosure.Length = strconv.Itoa(media0.File.FileSize)
		enclosure.Url = media0.URL
	}

	// Generate emojified content.
	apiEmojis := c.emojisToAPI(ctx, status.Emojis, status.EmojiIDs)
	content := text.EmojifyRSS(apiEmojis, status.Content)

	return &feeds.Item{
		// we specifcally do not set the author, as a lot
		// of feed readers rely on the RSS standard of the
		// author being an email with optional name. but
		// our @username@domain identifiers break this.
		//
		// attribution is handled in the title/description.

		// ID -- A string that uniquely identifies the item.
		// example: http://inessential.com/2002/09/01.php#a2
		Id: status.URL,

		// Source -- The RSS channel that the item came from.
		Source: &feeds.Link{Href: status.Account.URL + "/feed.rss"},

		// Link -- The URL of the item.
		// example: http://nytimes.com/2004/12/07FEST.html
		Link: &feeds.Link{Href: status.URL},

		Title:       title,
		Description: description,
		IsPermaLink: "true",
		Updated:     status.EditedAt,
		Created:     status.CreatedAt,
		Enclosure:   enclosure,
		Content:     content,
	}, nil
}

// trimTo trims the given `in` string to
// the length `to`, measured in runes.
//
// The reason for using runes is to avoid
// cutting off UTF-8 characters in the
// middle, and generating garbled bytes.
//
// If trimming was necessary, the returned
// string will be suffixed with ellipsis
// (`...`) to indicate omission.
func trimTo(in string, to int) string {
	var (
		runes    = []rune(in)
		runesLen = len(runes)
	)

	if runesLen <= to {
		// Fine as-is.
		return in
	}

	return string(runes[:to-3]) + "..."
}
