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

package dereferencing

import (
	"context"
	"errors"
	"net/url"

	"code.superseriousbusiness.org/gopkg/log"
	"code.superseriousbusiness.org/gopkg/xslices"
	"code.superseriousbusiness.org/gotosocial/internal/db"
	"code.superseriousbusiness.org/gotosocial/internal/gtserror"
	"code.superseriousbusiness.org/gotosocial/internal/gtsmodel"
	"code.superseriousbusiness.org/gotosocial/internal/id"
	"code.superseriousbusiness.org/gotosocial/internal/media"
	"code.superseriousbusiness.org/gotosocial/internal/util"
)

type statusChanges struct {
	poll     bool
	mentions bool
	tags     bool
	media    bool
	emoji    bool
	policy   bool
}

// handleStatusPeripherals handles dereferencing
// and storage of peripheral stuff: polls, mentions,
// tags, media, emoji, and interaction policy.
//
// The returned struct indicates what has changed
// between existing and new status (useful for edits).
//
// "Existing" status can be a barebones model in case
// of brand new statuses, but it must *not* be nil.
func (d *Dereferencer) handleStatusPeripherals(
	ctx context.Context,
	requester *gtsmodel.Account,
	uri *url.URL,
	existing *gtsmodel.Status,
	status *gtsmodel.Status,
) (
	statusChanges,
	error,
) {
	// Insert / update any attached status poll.
	pollChanged, err := d.handleStatusPoll(ctx,
		existing,
		status,
	)
	if err != nil {
		err := gtserror.Newf("error handling poll for status %s: %w", uri, err)
		return statusChanges{}, err
	}

	// Populate mentions associated with status, passing
	// in existing status to reuse old where possible.
	// (especially important here to reduce need to dereference).
	mentionsChanged, err := d.fetchStatusMentions(ctx,
		requester,
		existing,
		status,
	)
	if err != nil {
		err := gtserror.Newf("error populating mentions for status %s: %w", uri, err)
		return statusChanges{}, err
	}

	// Populate tags associated with status, passing
	// in existing status to reuse old where possible.
	tagsChanged, err := d.fetchStatusTags(ctx,
		existing,
		status,
	)
	if err != nil {
		err := gtserror.Newf("error populating tags for status %s: %w", uri, err)
		return statusChanges{}, err
	}

	// Check if there's any limits in place for (sub)domain.
	limit, err := d.state.DB.MatchDomainLimit(ctx, uri.Host)
	if err != nil {
		err := gtserror.Newf("error matching domain limit: %w", err)
		return statusChanges{}, err
	}

	// If domain media is limited, set reject reason,
	// this gets passed to media fetching functions
	// and prevents download of attached status media.
	var rejectReason *gtsmodel.MediaErrorDetails
	if limit.MediaReject() {
		rejectReason = new(gtsmodel.MediaErrorDetails)
		*rejectReason = gtsmodel.NewMediaErrorDetails(
			gtsmodel.MediaErrorTypePolicy,
			gtsmodel.MediaErrorTypePolicy_Domain,
		)
	}

	// Populate media attachments associated with status,
	// passing in existing status to reuse old where possible
	// (especially important here to reduce need to dereference).
	mediaChanged, err := d.fetchStatusAttachments(ctx,
		requester,
		existing,
		status,
		rejectReason,
	)
	if err != nil {
		err := gtserror.Newf("error populating attachments for status %s: %w", uri, err)
		return statusChanges{}, err
	}

	// Populate emoji associated with status, passing
	// in existing status to reuse old where possible
	// (especially important here to reduce need to dereference).
	emojiChanged, err := d.fetchStatusEmojis(ctx,
		existing,
		status,
		rejectReason,
	)
	if err != nil {
		err := gtserror.Newf("error populating emojis for status %s: %w", uri, err)
		return statusChanges{}, err
	}

	// Check if interaction policy has changed between latest and existing status.
	policyChanged := status.InteractionPolicy.DifferentFrom(existing.InteractionPolicy)

	return statusChanges{
		poll:     pollChanged,
		mentions: mentionsChanged,
		tags:     tagsChanged,
		media:    mediaChanged,
		emoji:    emojiChanged,
		policy:   policyChanged,
	}, nil
}

// fetchStatusMentions populates the mentions on 'status', creating
// new where needed, or using unchanged mentions from 'existing' status.
func (d *Dereferencer) fetchStatusMentions(
	ctx context.Context,
	requester *gtsmodel.Account,
	existing *gtsmodel.Status,
	status *gtsmodel.Status,
) (
	changed bool,
	err error,
) {

	// Get most-recent modified time
	// for use in new mention ULIDs.
	updatedAt := status.UpdatedAt()

	// Allocate new slice to take the yet-to-be created mention IDs.
	status.MentionIDs = make([]string, len(status.Mentions))

	for i := range status.Mentions {
		mention := status.Mentions[i]

		// Search existing status + db for a mention already stored,
		// else ensure new mention's target account is populated.
		mention, alreadyExists, err := d.newOrExistingMention(ctx,
			requester,
			existing,
			mention,
		)
		switch {
		case err == nil:
			// No problem.

		case gtserror.IsUnretrievable(err):
			// Unretrievable means the remote is
			// probably blocked or webfinger wasn't working.
			log.Debugf(ctx, "unretrievable mention: %v", err)
			continue

		default:
			// Some kind of actual error.
			log.Errorf(ctx, "failed to derive mention: %v", err)
			continue
		}

		if alreadyExists {
			// This mention was already
			// stored, use it and continue.
			status.Mentions[i] = mention
			status.MentionIDs[i] = mention.ID
			continue
		}

		// Mark status as
		// having changed.
		changed = true

		// This mention didn't exist yet.
		// Generate new ID according to latest update.
		mention.ID = id.NewULIDFromTime(updatedAt)

		// Set further mention details.
		mention.CreatedAt = updatedAt
		mention.OriginAccount = status.Account
		mention.OriginAccountID = status.AccountID
		mention.OriginAccountURI = status.AccountURI
		mention.TargetAccountID = mention.TargetAccount.ID
		mention.TargetAccountURI = mention.TargetAccount.URI
		mention.TargetAccountURL = mention.TargetAccount.URL
		mention.StatusID = status.ID
		mention.Status = status
		mention.IsNew = true

		// Place the new mention into the database.
		if err := d.state.DB.PutMention(ctx, mention); err != nil {
			return changed, gtserror.Newf("error putting mention in database: %w", err)
		}

		// Set the *new* mention and ID.
		status.Mentions[i] = mention
		status.MentionIDs[i] = mention.ID
	}

	for i := 0; i < len(status.MentionIDs); {
		if status.MentionIDs[i] == "" {
			// This is a failed mention population, likely due
			// to invalid incoming data / now-deleted accounts.
			copy(status.Mentions[i:], status.Mentions[i+1:])
			copy(status.MentionIDs[i:], status.MentionIDs[i+1:])
			status.Mentions = status.Mentions[:len(status.Mentions)-1]
			status.MentionIDs = status.MentionIDs[:len(status.MentionIDs)-1]
			continue
		}
		i++
	}

	return changed, nil
}

// fetchStatusTags populates the tags on 'status', fetching existing
// from the database and creating new where needed. 'existing' is used
// to fetch tags that have not changed since previous stored status.
func (d *Dereferencer) fetchStatusTags(
	ctx context.Context,
	existing *gtsmodel.Status,
	status *gtsmodel.Status,
) (
	changed bool,
	err error,
) {

	// Allocate new slice to take the yet-to-be determined tag IDs.
	status.TagIDs = make([]string, len(status.Tags))

	for i := range status.Tags {
		tag := status.Tags[i]

		// Look for tag in existing status with name.
		existing, ok := existing.GetTagByName(tag.Name)
		if ok && existing.ID != "" {
			status.Tags[i] = existing
			status.TagIDs[i] = existing.ID
			continue
		}

		// Mark status as
		// having changed.
		changed = true

		// Look for existing tag with name in the database.
		existing, err := d.state.DB.GetTagByName(ctx, tag.Name)
		if err != nil && !errors.Is(err, db.ErrNoEntries) {
			return changed, gtserror.Newf("db error getting tag %s: %w", tag.Name, err)
		} else if existing != nil {
			status.Tags[i] = existing
			status.TagIDs[i] = existing.ID
			continue
		}

		// Create new ID for tag.
		tag.ID = id.NewULID()

		// Insert this tag with new name into the database.
		if err := d.state.DB.PutTag(ctx, tag); err != nil {
			log.Errorf(ctx, "db error putting tag %s: %v", tag.Name, err)
			continue
		}

		// Set new tag ID in slice.
		status.TagIDs[i] = tag.ID
	}

	// Remove any tag we couldn't get or create.
	for i := 0; i < len(status.TagIDs); {
		if status.TagIDs[i] == "" {
			// This is a failed tag population, likely due
			// to some database peculiarity / race condition.
			copy(status.Tags[i:], status.Tags[i+1:])
			copy(status.TagIDs[i:], status.TagIDs[i+1:])
			status.Tags = status.Tags[:len(status.Tags)-1]
			status.TagIDs = status.TagIDs[:len(status.TagIDs)-1]
			continue
		}
		i++
	}

	return changed, nil
}

// fetchStatusAttachments populates the attachments on 'status', creating new database
// entries where needed and dereferencing it, or using unchanged from 'existing' status.
func (d *Dereferencer) fetchStatusAttachments(
	ctx context.Context,
	requester *gtsmodel.Account,
	existing *gtsmodel.Status,
	status *gtsmodel.Status,
	rejectReason *gtsmodel.MediaErrorDetails, // optional reason to reject media with
) (
	changed bool,
	err error,
) {
	// Allocate new slice to take the yet-to-be fetched attachment IDs.
	status.AttachmentIDs = make([]string, len(status.Attachments))
	for i := range status.Attachments {
		placeholder := status.Attachments[i]

		// Look for existing media attachment with remote URL first.
		existing, ok := existing.GetAttachmentByRemoteURL(placeholder.RemoteURL)
		if ok && existing.ID != "" {

			info := media.AdditionalMediaInfo{
				// Pass reject reason ptr, which
				// will skip downloading if set.
				RejectReason: rejectReason,
			}

			// Look for any difference in stored media description.
			diff := (existing.Description != placeholder.Description)
			if diff {
				info.Description = &placeholder.Description
			}

			// If description changed,
			// we mark media as changed.
			changed = changed || diff

			// Store any attachment updates and
			// ensure media is locally cached.
			existing, err := d.RefreshMedia(ctx,
				requester,
				existing,
				info,
				diff,
				true, // async
			)
			if err != nil {
				log.Errorf(ctx, "error updating existing attachment: %v", err)

				// specifically do NOT continue here,
				// we already have a model, we don't
				// want to drop it from the status, just
				// log that an update for it failed.
			}

			// Set the existing attachment.
			status.Attachments[i] = existing
			status.AttachmentIDs[i] = existing.ID
			continue
		}

		// Mark status as
		// having changed.
		changed = true

		// Load this new media attachment.
		attachment, err := d.GetMedia(ctx,
			requester,
			status.AccountID,
			placeholder.RemoteURL,
			media.AdditionalMediaInfo{
				StatusID:    &status.ID,
				RemoteURL:   &placeholder.RemoteURL,
				Description: &placeholder.Description,
				Blurhash:    &placeholder.Blurhash,
				FocusX:      &placeholder.FileMeta.Focus.X,
				FocusY:      &placeholder.FileMeta.Focus.Y,

				// Pass reject reason ptr, which
				// will skip downloading if set.
				RejectReason: rejectReason,
			},
			false, // async
		)
		if err != nil {
			if attachment == nil {
				log.Errorf(ctx, "error loading attachment %s: %v", placeholder.RemoteURL, err)
				continue
			}

			// non-fatal error occurred during loading, still use it.
			log.Warnf(ctx, "partially loaded attachment: %v", err)
		}

		// Set the *new* attachment and ID.
		status.Attachments[i] = attachment
		status.AttachmentIDs[i] = attachment.ID
	}

	for i := 0; i < len(status.AttachmentIDs); {
		if status.AttachmentIDs[i] == "" {
			// Remove totally failed attachment populations
			copy(status.Attachments[i:], status.Attachments[i+1:])
			copy(status.AttachmentIDs[i:], status.AttachmentIDs[i+1:])
			status.Attachments = status.Attachments[:len(status.Attachments)-1]
			status.AttachmentIDs = status.AttachmentIDs[:len(status.AttachmentIDs)-1]
			continue
		}
		i++
	}

	return changed, nil
}

// fetchStatusEmojis populates the emojis on 'status', creating new database entries
// where needed and dereferencing it, or using unchanged from 'existing' status.
func (d *Dereferencer) fetchStatusEmojis(
	ctx context.Context,
	existing *gtsmodel.Status,
	status *gtsmodel.Status,
	rejectReason *gtsmodel.MediaErrorDetails, // optional reason to reject media with
) (
	changed bool,
	err error,
) {

	// Fetch the updated emojis for our status.
	emojis, changed, err := d.fetchEmojis(ctx,
		existing.Emojis,
		status.Emojis,
		rejectReason,
	)
	if err != nil {
		return changed, gtserror.Newf("error fetching emojis: %w", err)
	}

	if !changed {
		// Use existing status emoji objects.
		status.EmojiIDs = existing.EmojiIDs
		status.Emojis = existing.Emojis
		return false, nil
	}

	// Set latest emojis.
	status.Emojis = emojis

	// Extract IDs from latest slice of emojis.
	status.EmojiIDs = make([]string, len(emojis))
	for i, emoji := range emojis {
		status.EmojiIDs[i] = emoji.ID
	}

	// Combine both old and new emojis, as statuses.emojis
	// keeps track of emojis for both old and current edits.
	status.EmojiIDs = append(status.EmojiIDs, existing.EmojiIDs...)
	status.Emojis = append(status.Emojis, existing.Emojis...)
	status.EmojiIDs = xslices.Deduplicate(status.EmojiIDs)
	status.Emojis = xslices.DeduplicateFunc(status.Emojis,
		func(e *gtsmodel.Emoji) string { return e.ID },
	)

	return true, nil
}

// handleStatusPoll handles both inserting of new status poll or the
// update of an existing poll. this handles the case of simple vote
// count updates (without being classified as a change of the poll
// itself), as well as full poll changes that delete existing instance.
func (d *Dereferencer) handleStatusPoll(
	ctx context.Context,
	existing *gtsmodel.Status,
	status *gtsmodel.Status,
) (
	changed bool,
	err error,
) {
	switch {
	case existing.Poll == nil && status.Poll == nil:
		// no poll before or after, nothing to do.
		return false, nil

	case existing.Poll == nil && status.Poll != nil:
		// no previous poll, insert new status poll!
		return true, d.insertStatusPoll(ctx, status)

	case status.Poll == nil:
		// existing status poll has been deleted, remove this from the database.
		if err = d.state.DB.DeletePollByID(ctx, existing.Poll.ID); err != nil {
			err = gtserror.Newf("error deleting poll from database: %w", err)
		}
		return true, err

	case pollChanged(existing.Poll, status.Poll):
		// existing status poll has been changed, remove this from the database.
		if err = d.state.DB.DeletePollByID(ctx, existing.Poll.ID); err != nil {
			return true, gtserror.Newf("error deleting poll from database: %w", err)
		}

		// insert latest poll version into database.
		return true, d.insertStatusPoll(ctx, status)

	case pollStateUpdated(existing.Poll, status.Poll):
		// Since we last saw it, the poll has updated!
		// Whether that be stats, or close time.
		poll := existing.Poll
		poll.Closing = pollJustClosed(existing.Poll, status.Poll)
		poll.ClosedAt = status.Poll.ClosedAt
		poll.Voters = status.Poll.Voters
		poll.Votes = status.Poll.Votes

		// Update poll model in the database (specifically only the possible changed columns).
		if err = d.state.DB.UpdatePoll(ctx, poll, "closed_at", "voters", "votes"); err != nil {
			return false, gtserror.Newf("error updating poll: %w", err)
		}

		// Update poll on status.
		status.PollID = poll.ID
		status.Poll = poll
		return false, nil

	default:
		// latest and existing
		// polls are up to date.
		poll := existing.Poll
		status.PollID = poll.ID
		status.Poll = poll
		return false, nil
	}
}

// insertStatusPoll inserts an assumed new poll attached to status into the database, this
// also handles generating new ID for the poll and setting necessary fields on the status.
func (d *Dereferencer) insertStatusPoll(ctx context.Context, status *gtsmodel.Status) error {
	var err error

	// Get most-recent modified time
	// which will be poll creation time.
	createdAt := status.UpdatedAt()

	// Generate new ID for poll from createdAt.
	status.Poll.ID = id.NewULIDFromTime(createdAt)

	// Update the status<->poll links.
	status.PollID = status.Poll.ID
	status.Poll.StatusID = status.ID
	status.Poll.Status = status

	// Insert this latest poll into the database.
	err = d.state.DB.PutPoll(ctx, status.Poll)
	if err != nil {
		return gtserror.Newf("error putting poll in database: %w", err)
	}

	return nil
}

// handleStatusEdit compiles a list of changed status table columns between
// existing and latest status model, and where necessary inserts a historic
// edit of the status into the database to store its previous state. the
// returned slice is a list of columns requiring updating in the database.
func (d *Dereferencer) handleStatusEdit(
	ctx context.Context,
	existing *gtsmodel.Status,
	status *gtsmodel.Status,
	changes statusChanges,
) (
	cols []string,
	err error,
) {
	var edited bool

	// Copy previous status edit columns.
	status.EditIDs = existing.EditIDs
	status.Edits = existing.Edits

	// Preallocate max slice length.
	cols = make([]string, 1, 13)

	// Always update `fetched_at`.
	cols[0] = "fetched_at"

	// Check for edited status content.
	if existing.Content != status.Content {
		cols = append(cols, "content")
		edited = true
	}

	// Check for edited status content warning.
	if existing.ContentWarning != status.ContentWarning {
		cols = append(cols, "content_warning")
		edited = true
	}

	// Check for edited status sensitive flag.
	if existing.Flags.Sensitive() != status.Flags.Sensitive() {
		cols = append(cols, "flags")
		edited = true
	}

	// Check for edited status language tag.
	if existing.Language != status.Language {
		cols = append(cols, "language")
		edited = true
	}

	if changes.poll {
		// Attached poll was changed.
		cols = append(cols, "poll_id")
		edited = true
	}

	if changes.mentions {
		cols = append(cols, "mentions") // i.e. MentionIDs

		// Mentions changed doesn't necessarily
		// indicate an edit, it may just not have
		// been previously populated properly.
	}

	if changes.tags {
		cols = append(cols, "tags") // i.e. TagIDs

		// Tags changed doesn't necessarily
		// indicate an edit, it may just not have
		// been previously populated properly.
	}

	if changes.media {
		// Attached media was changed.
		cols = append(cols, "attachments") // i.e. AttachmentIDs
		edited = true
	}

	if changes.emoji {
		// Attached emojis changed.
		cols = append(cols, "emojis") // i.e. EmojiIDs

		// We specifically store both *new* AND *old* edit
		// revision emojis in the statuses.emojis column.
		emojiByID := func(e *gtsmodel.Emoji) string { return e.ID }
		status.Emojis = append(status.Emojis, existing.Emojis...)
		status.Emojis = xslices.DeduplicateFunc(status.Emojis, emojiByID)
		status.EmojiIDs = xslices.Gather(status.EmojiIDs[:0], status.Emojis, emojiByID)

		// Emojis changed doesn't necessarily
		// indicate an edit, it may just not have
		// been previously populated properly.
	}

	if changes.policy {
		// Interaction policy changed.
		cols = append(cols, "interaction_policy")

		// Int pol changed doesn't necessarily
		// indicate an edit, it may just not have
		// been previously populated properly.
	}

	if edited {
		// Get previous-most-recent modified time,
		// which will be this edit's creation time.
		createdAt := existing.UpdatedAt()

		// Status has been editted since last
		// we saw it, take snapshot of existing.
		var edit gtsmodel.StatusEdit
		edit.ID = id.NewULIDFromTime(createdAt)
		edit.Content = existing.Content
		edit.ContentWarning = existing.ContentWarning
		edit.Text = existing.Text
		edit.ContentType = existing.ContentType
		edit.Language = existing.Language
		edit.Sensitive = util.Ptr(existing.Flags.Sensitive())
		edit.StatusID = status.ID
		edit.CreatedAt = createdAt

		// Copy existing attachments and descriptions.
		edit.AttachmentIDs = existing.AttachmentIDs
		edit.Attachments = existing.Attachments
		if l := len(existing.Attachments); l > 0 {
			edit.AttachmentDescriptions = make([]string, l)
			for i, attach := range existing.Attachments {
				edit.AttachmentDescriptions[i] = attach.Description
			}
		}

		if existing.Poll != nil {
			// Poll only set if existing contained them.
			edit.PollOptions = existing.Poll.Options

			if changes.poll || !*existing.Poll.HideCounts ||
				!existing.Poll.ClosedAt.IsZero() {
				// If the counts are allowed to be
				// shown, or poll has changed, then
				// include poll vote counts in edit.
				edit.PollVotes = existing.Poll.Votes
			}
		}

		// Insert this new edit of existing status into database.
		if err := d.state.DB.PutStatusEdit(ctx, &edit); err != nil {
			return nil, gtserror.Newf("error putting edit in database: %w", err)
		}

		// Add edit to list of edits on the status.
		status.EditIDs = append(status.EditIDs, edit.ID)
		status.Edits = append(status.Edits, &edit)

		// Add edit to list of cols.
		cols = append(cols, "edits")
	}

	if !existing.EditedAt.Equal(status.EditedAt) {
		// Whether status edited or not,
		// edited_at column has changed.
		cols = append(cols, "edited_at")

		// Mark the status as edited for
		// appropriate handling in hooks.
		status.Edited = true
	}

	return cols, nil
}

// newOrExistingMention tries to populate the given
// mention with the correct TargetAccount and (if not
// yet set) TargetAccountURI, returning the populated
// mention.
//
// Will check on the existing status and in the db
// if the mention is already there and populated;
// if so, existing mention will be returned along
// with `true` to indicate that it already existed.
//
// Otherwise, this function will try to parse first
// the Href of the mention, and then the namestring,
// to see who it targets, and go fetch that account.
//
// Note: Ordinarily it would make sense to try the
// namestring first, as it definitely can't be a URL
// rather than a URI, but because some remotes do
// silly things like only provide `@username` instead
// of `@username@domain`, we try by URI first.
func (d *Dereferencer) newOrExistingMention(
	ctx context.Context,
	requester *gtsmodel.Account,
	existing *gtsmodel.Status,
	mention *gtsmodel.Mention,
) (
	*gtsmodel.Mention,
	bool, // True if mention already exists in the DB.
	error,
) {
	// Mentions can be created using `name` or `href`.
	//
	// Prefer `href` (TargetAccountURI), fall back to Name.
	switch {
	case mention.TargetAccountURI != "":
		// Look on the status for existing mention with target account's URI.
		existingMention, ok := existing.GetMentionByTargetURI(mention.TargetAccountURI)
		if ok && existingMention.ID != "" {
			// Already populated
			// mention, use this.
			return existingMention, true, nil
		}

		// Ensure that mention account URI is parseable.
		targetAccountURI, err := url.Parse(mention.TargetAccountURI)
		if err != nil {
			err := gtserror.Newf("invalid account uri %q: %w", mention.TargetAccountURI, err)
			return nil, false, err
		}

		// Ensure we have account of mention target dereferenced.
		mention.TargetAccount, _, err = d.getAccountByURI(ctx,
			requester,
			targetAccountURI,
			false, // tryURL: URI match only here
		)
		if err != nil {
			err := gtserror.Newf("failed to dereference account %s: %w", targetAccountURI, err)
			return nil, false, err
		}

		// Look in the db for this existing mention.
		existingMention, err = d.state.DB.GetMentionByTargetAcctStatus(
			ctx,
			mention.TargetAccount.ID,
			existing.ID,
		)
		if err != nil && !errors.Is(err, db.ErrNoEntries) {
			err := gtserror.Newf("db error looking for existing mention: %w", err)
			return nil, false, err
		}

		if existingMention != nil {
			// Already had stored
			// mention, use this.
			return existingMention, true, nil
		}

	case mention.NameString != "":
		// Href wasn't set, extract the username and domain parts from namestring.
		username, domain, err := util.ExtractNamestringParts(mention.NameString)
		if err != nil {
			err := gtserror.Newf("failed to parse namestring %s: %w", mention.NameString, err)
			return nil, false, err
		}

		// Look on the status for existing mention with username domain target.
		existingMention, ok := existing.GetMentionByUsernameDomain(username, domain)
		if ok && existingMention.ID != "" {
			// Already populated
			// mention, use this.
			return existingMention, true, nil
		}

		// Ensure we have account of mention target dereferenced.
		//
		// This might fail if remote only sets `@username` (no domain).
		mention.TargetAccount, _, err = d.getAccountByUsernameDomain(ctx,
			requester,
			username,
			domain,
		)
		if err != nil && !errors.Is(err, db.ErrNoEntries) {
			err := gtserror.Newf("failed to dereference account %s: %w", mention.NameString, err)
			return nil, false, err
		}

		if mention.TargetAccount == nil {
			// Probably failed for abovementioned
			// silly reason. Nothing we can do about it.
			err := gtserror.Newf(
				"failed to populate mention target account (badly formatted namestring?) %s: %w",
				mention.NameString, err,
			)
			return nil, false, err
		}

		// Look on the status for existing mention with target account's URI.
		existingMention, ok = existing.GetMentionByTargetURI(mention.TargetAccountURI)
		if ok && existingMention.ID != "" {
			// Already populated
			// mention, use this.
			return existingMention, true, nil
		}

		// Look in the db for this existing mention.
		existingMention, err = d.state.DB.GetMentionByTargetAcctStatus(
			ctx,
			mention.TargetAccount.ID,
			existing.ID,
		)
		if err != nil && !errors.Is(err, db.ErrNoEntries) {
			err := gtserror.Newf("db error looking for existing mention: %w", err)
			return nil, false, err
		}

		if existingMention != nil {
			// Already had stored
			// mention, use this.
			return existingMention, true, nil
		}

	default:
		const errText = "neither target uri nor namestring were set on mention, cannot process it"
		return nil, false, gtserror.New(errText)
	}

	// At this point, mention.TargetAccountURI
	// and mention.TargetAccount must be set.
	return mention, false, nil
}
