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

package account

import (
	"context"
	"errors"
	"time"

	"code.superseriousbusiness.org/gotosocial/internal/config"
	"code.superseriousbusiness.org/gotosocial/internal/db"
	"code.superseriousbusiness.org/gotosocial/internal/gtserror"
	"code.superseriousbusiness.org/gotosocial/internal/gtsmodel"
	"code.superseriousbusiness.org/gotosocial/internal/paging"
	"github.com/gorilla/feeds"
)

var never time.Time

type GetRSSFeed func() (*feeds.Feed, gtserror.WithCode)

// GetRSSFeedForUsername returns a function to return the RSS feed of a local account
// with the given username, and the last-modified time (time that the account last
// posted a status eligible to be included in the rss feed).
//
// To save db calls, callers to this function should only call the returned GetRSSFeed
// func if the last-modified time is newer than the last-modified time they have cached.
//
// If the account has not yet posted an RSS-eligible status, the returned last-modified
// time will be zero, and the GetRSSFeed func will return a valid RSS xml with no items.
func (p *Processor) GetRSSFeedForUsername(ctx context.Context, username string, page *paging.Page) (GetRSSFeed, time.Time, gtserror.WithCode) {

	// Fetch local (i.e. empty domain) account from database by username.
	account, err := p.state.DB.GetAccountByUsernameDomain(ctx, username, "")
	if err != nil {
		err := gtserror.Newf("db error getting account %s: %w", username, err)
		return nil, never, gtserror.NewErrorInternalError(err)
	}

	// Check if exists.
	if account == nil {
		err := gtserror.New("account not found")
		return nil, never, gtserror.NewErrorNotFound(err)
	}

	// Ensure account has rss feed enabled.
	if !*account.Settings.EnableRSS {
		err := gtserror.New("account RSS feed not enabled")
		return nil, never, gtserror.NewErrorNotFound(err)
	}

	// Ensure account stats populated for last status fetch information.
	if err := p.state.DB.PopulateAccountStats(ctx, account); err != nil {
		err := gtserror.Newf("db error getting account stats %s: %w", username, err)
		return nil, never, gtserror.NewErrorInternalError(err)
	}

	// LastModified time is needed by callers to check freshness for cacheing.
	// This might be a zero time.Time if account has never posted a status that's
	// eligible to appear in the RSS feed; that's fine.
	lastPostAt := account.Stats.LastStatusAt

	return func() (*feeds.Feed, gtserror.WithCode) {
		var image *feeds.Image

		// Assemble author namestring.
		author := "@" + account.Username +
			"@" + config.GetAccountDomain()

		// Check if account has an avatar media attachment.
		if id := account.AvatarMediaAttachmentID; id != "" {
			if account.AvatarMediaAttachment == nil {
				var err error

				// Populate the account's avatar media attachment from database by its ID.
				account.AvatarMediaAttachment, err = p.state.DB.GetAttachmentByID(ctx, id)
				if err != nil && !errors.Is(err, db.ErrNoEntries) {
					err := gtserror.Newf("db error getting account avatar: %w", err)
					return nil, gtserror.NewErrorInternalError(err)
				}
			}

			// If avatar is found, use as feed image.
			if account.AvatarMediaAttachment != nil {
				image = &feeds.Image{
					Title: "Avatar for " + author,
					Url:   account.AvatarMediaAttachment.Thumbnail.URL,
					Link:  account.URL,
				}
			}
		}

		// Start creating feed.
		feed := &feeds.Feed{
			// we specifcally do not set the author, as a lot
			// of feed readers rely on the RSS standard of the
			// author being an email with optional name. but
			// our @username@domain identifiers break this.
			//
			// attribution is handled in the title/description.

			Title:       "Posts from " + author,
			Description: "Posts from " + author,
			Link:        &feeds.Link{Href: account.URL},
			Image:       image,
		}

		// If the account has never posted anything, just use
		// account creation time as Updated value for the feed;
		// we could use time.Now() here but this would likely
		// mess up cacheing; we want something determinate.
		//
		// We can also return early rather than wasting a db call,
		// since we already know there's no eligible statuses.
		if lastPostAt.IsZero() {
			feed.Updated = account.CreatedAt
			return feed, nil
		}

		// Account has posted at least one status that's
		// eligible to appear in the RSS feed.
		//
		// Reuse the lastPostAt value for feed.Updated.
		feed.Updated = lastPostAt

		// Retrieve latest statuses as they'd be shown
		// on the web view of the account profile.
		//
		// Take into account whether the user wants
		// their web view laid out in gallery mode.
		mediaOnly := (account.Settings != nil &&
			account.Settings.WebLayout == gtsmodel.WebLayoutGallery)
		statuses, err := p.state.DB.GetAccountWebStatuses(
			ctx,
			account,
			page,
			mediaOnly,
			false, // don't include boosts
		)
		if err != nil && !errors.Is(err, db.ErrNoEntries) {
			err := gtserror.Newf("db error getting account web statuses: %w", err)
			return nil, gtserror.NewErrorInternalError(err)
		}

		// Check for no statuses.
		if len(statuses) == 0 {
			return feed, nil
		}

		// Get next / prev paging parameters.
		lo := statuses[len(statuses)-1].ID
		hi := statuses[0].ID
		next := page.Next(lo, hi)
		prev := page.Prev(lo, hi)

		// Add each status to the rss feed.
		for _, status := range statuses {
			item, err := p.converter.StatusToRSSItem(ctx, status)
			if err != nil {
				err := gtserror.Newf("error converting status to feed item: %w", err)
				return nil, gtserror.NewErrorInternalError(err)
			}
			feed.Add(item)
		}

		// TODO: when we have some manner of supporting
		// atom:link in RSS (and Atom), set the paging
		// parameters for next / prev feed pages here.
		_, _ = next, prev

		return feed, nil
	}, lastPostAt, nil
}
