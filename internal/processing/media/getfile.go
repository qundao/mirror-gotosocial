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

package media

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	apimodel "code.superseriousbusiness.org/gotosocial/internal/api/model"
	"code.superseriousbusiness.org/gotosocial/internal/db"
	"code.superseriousbusiness.org/gotosocial/internal/gtserror"
	"code.superseriousbusiness.org/gotosocial/internal/gtsmodel"
	"code.superseriousbusiness.org/gotosocial/internal/media"
	"code.superseriousbusiness.org/gotosocial/internal/regexes"
	"code.superseriousbusiness.org/gotosocial/internal/storage"
	"code.superseriousbusiness.org/gotosocial/internal/uris"
	"code.superseriousbusiness.org/gotosocial/internal/util"
)

// GetFile retrieves a file from storage and streams it back
// to the caller via an io.reader embedded in *apimodel.Content.
func (p *Processor) GetFile(
	ctx context.Context,
	requester *gtsmodel.Account,
	form *apimodel.GetContentRequestForm,
) (*apimodel.Content, gtserror.WithCode) {
	// Parse media size (small, static, original).
	mediaSize, err := parseSize(form.MediaSize)
	if err != nil {
		err := gtserror.Newf("media size %s not valid", form.MediaSize)
		return nil, gtserror.NewErrorNotFound(err)
	}

	// Parse media type (emoji, header, avatar, attachment).
	mediaType, err := parseType(form.MediaType)
	if err != nil {
		err := gtserror.Newf("media type %s not valid", form.MediaType)
		return nil, gtserror.NewErrorNotFound(err)
	}

	// Parse media ID from file name.
	mediaID, _, err := parseFileName(form.FileName)
	if err != nil {
		err := gtserror.Newf("media file name %s not valid", form.FileName)
		return nil, gtserror.NewErrorNotFound(err)
	}

	// Get the account that owns the media
	// and make sure it's not suspended.
	acctID := form.AccountID
	acct, err := p.state.DB.GetAccountByID(ctx, acctID)
	if err != nil {
		err := gtserror.Newf("db error getting account %s: %w", acctID, err)
		return nil, gtserror.NewErrorNotFound(err)
	}

	if acct.IsSuspended() {
		err := gtserror.Newf("account %s is suspended", acctID)
		return nil, gtserror.NewErrorNotFound(err)
	}

	// If requester was authenticated, ensure media
	// owner and requester don't block each other.
	if requester != nil {
		blocked, err := p.state.DB.IsEitherBlocked(ctx, requester.ID, acctID)
		if err != nil {
			err := gtserror.Newf("db error checking block between %s and %s: %w", acctID, requester.ID, err)
			return nil, gtserror.NewErrorNotFound(err)
		}

		if blocked {
			err := gtserror.Newf("block exists between %s and %s", acctID, requester.ID)
			return nil, gtserror.NewErrorNotFound(err)
		}
	}

	// The way we store emojis is a bit different
	// from the way we store other attachments,
	// so we need to take different steps depending
	// on the media type being requested.
	switch mediaType {

	case media.TypeEmoji:
		return p.getEmojiContent(ctx,
			acctID,
			mediaSize,
			mediaID,
		)

	case media.TypeAttachment, media.TypeHeader, media.TypeAvatar:
		return p.getAttachmentContent(ctx,
			requester,
			acctID,
			mediaSize,
			mediaID,
		)

	default:
		err := gtserror.Newf("media type %s not recognized", mediaType)
		return nil, gtserror.NewErrorNotFound(err)
	}
}

func (p *Processor) getAttachmentContent(
	ctx context.Context,
	requester *gtsmodel.Account,
	acctID string,
	sizeStr media.Size,
	mediaID string,
) (
	*apimodel.Content,
	gtserror.WithCode,
) {
	// Get attachment with given ID from the database.
	attach, err := p.state.DB.GetAttachmentByID(ctx, mediaID)
	if err != nil && !errors.Is(err, db.ErrNoEntries) {
		err := gtserror.Newf("db error getting attachment %s: %w", mediaID, err)
		return nil, gtserror.NewErrorInternalError(err)
	}

	if attach == nil {
		const text = "media not found"
		return nil, gtserror.NewErrorNotFound(errors.New(text), text)
	}

	// Ensure the account
	// actually owns the media.
	if attach.AccountID != acctID {
		const text = "media was not owned by passed account id"
		return nil, gtserror.NewErrorNotFound(errors.New(text) /* no help text! */)
	}

	// Unknown file types indicate no *locally*
	// stored data we can serve. Handle separately.
	if attach.Type == gtsmodel.FileTypeUnknown {
		return handleUnknown(attach)
	}

	// If requester was provided, use their username
	// to create a transport to potentially re-fetch
	// the media. Else falls back to instance account.
	var requestUser string
	if requester != nil {
		requestUser = requester.Username
	}

	// Start preparing API content model and other
	// values depending on requested media size.
	var content apimodel.Content
	var mediaPath string
	switch sizeStr {

	// Original media size.
	case media.SizeOriginal:
		content.ContentType = attach.File.ContentType
		content.ContentLength = int64(attach.File.FileSize)
		mediaPath = attach.File.Path

	// Thumbnail media size.
	case media.SizeSmall:
		content.ContentType = attach.Thumbnail.ContentType
		content.ContentLength = int64(attach.Thumbnail.FileSize)
		mediaPath = attach.Thumbnail.Path

	default:
		const text = "invalid media size"
		return nil, gtserror.NewErrorBadRequest(
			errors.New(text),
			text,
		)
	}

	// Attachment file
	// stream from storage.
	var rc io.ReadCloser

	// Check media is meant
	// to be cached locally.
	if *attach.Cached {

		// Check storage for media at determined path.
		rc, err = p.state.Storage.GetStream(ctx, mediaPath)
		if err != nil && !storage.IsNotFound(err) {
			err := gtserror.Newf("storage error getting media %s: %w", attach.URL, err)
			return nil, gtserror.NewErrorInternalError(err)
		}
	}

	if rc == nil {
		// This is local media without
		// a cached attachment, unfulfillable!
		if attach.IsLocal() {
			return nil, gtserror.NewfWithCode(http.StatusNotFound,
				"local media file not found: %s", attach.URL)
		}

		// Whether the cached flag was set or
		// not, we know it isn't in storage.
		attach.Cached = util.Ptr(false)

		// Attempt to recache this remote media.
		attach, err = p.federator.RefreshMedia(ctx,
			requestUser,
			attach,
			media.AdditionalMediaInfo{},
			false,
		)
		if err != nil {
			err := gtserror.Newf("error recaching media %s: %w", attach.URL, err)
			return nil, gtserror.NewErrorNotFound(err)
		}

		// Check storage for media at determined path.
		rc, err = p.state.Storage.GetStream(ctx, mediaPath)
		if err != nil && !storage.IsNotFound(err) {
			err := gtserror.Newf("storage error getting media %s: %w", attach.URL, err)
			return nil, gtserror.NewErrorInternalError(err)
		} else if rc == nil {
			return nil, gtserror.NewfWithCode(http.StatusNotFound,
				"remote media file not found: %s", attach.URL)
		}
	}

	// If running on S3 storage with proxying disabled,
	// just fetch a pre-signed URL instead of the content.
	if url := p.state.Storage.URL(ctx, mediaPath); url != nil {
		_ = rc.Close() // close storage stream
		content.URL = url
		return &content, nil
	}

	// Return with stream.
	content.Content = rc
	return &content, nil
}

func (p *Processor) getEmojiContent(
	ctx context.Context,
	acctID string,
	sizeStr media.Size,
	emojiID string,
) (
	*apimodel.Content,
	gtserror.WithCode,
) {
	// Reconstruct static emoji image URL to search for it.
	// As refreshed emojis use a newly generated path ID to
	// differentiate them (cache-wise) from the original.
	staticURL := uris.URIForAttachment(
		acctID,
		string(media.TypeEmoji),
		string(media.SizeStatic),
		emojiID,
		"png",
	)

	// Search for emoji with given static URL in the database.
	emoji, err := p.state.DB.GetEmojiByStaticURL(ctx, staticURL)
	if err != nil && !errors.Is(err, db.ErrNoEntries) {
		err := gtserror.Newf("error fetching emoji from database: %w", err)
		return nil, gtserror.NewErrorInternalError(err)
	}

	if emoji == nil {
		const text = "emoji not found"
		return nil, gtserror.NewErrorNotFound(errors.New(text), text)
	}

	if *emoji.Disabled {
		const text = "emoji has been disabled"
		return nil, gtserror.NewErrorNotFound(errors.New(text), text)
	}

	// Start preparing API content model and other
	// values depending on requested media size.
	var content apimodel.Content
	var emojiPath string
	switch sizeStr {

	// Original emoji image.
	case media.SizeOriginal:
		content.ContentType = emoji.ImageContentType
		content.ContentLength = int64(emoji.ImageFileSize)
		emojiPath = emoji.ImagePath

	// Static emoji image.
	case media.SizeStatic:
		content.ContentType = emoji.ImageStaticContentType
		content.ContentLength = int64(emoji.ImageStaticFileSize)
		emojiPath = emoji.ImageStaticPath

	default:
		const text = "invalid emoji size"
		return nil, gtserror.NewErrorBadRequest(
			errors.New(text),
			text,
		)
	}

	// Emoji image file
	// stream from storage.
	var rc io.ReadCloser

	// Check emoji is meant
	// to be cached locally.
	if *emoji.Cached {

		// Check storage for emoji at determined image path.
		rc, err = p.state.Storage.GetStream(ctx, emojiPath)
		if err != nil && !storage.IsNotFound(err) {
			err := gtserror.Newf("storage error getting emoji %s: %w", emoji.URI, err)
			return nil, gtserror.NewErrorInternalError(err)
		}
	}

	if rc == nil {
		// This is a local emoji without
		// a cached image, unfulfillable!
		if emoji.IsLocal() {
			return nil, gtserror.NewfWithCode(http.StatusNotFound,
				"local emoji image not found: %s", emoji.URI)
		}

		// Whether the cached flag was set or
		// not, we know it isn't in storage.
		emoji.Cached = util.Ptr(false)

		// Attempt to recache this remote emoji.
		emoji, err = p.federator.RecacheEmoji(ctx,
			emoji,
			false,
		)
		if err != nil {
			err := gtserror.Newf("error recaching emoji %s: %w", emoji.URI, err)
			return nil, gtserror.NewErrorNotFound(err)
		}

		// Check storage for emoji at determined image path.
		rc, err = p.state.Storage.GetStream(ctx, emojiPath)
		if err != nil && !storage.IsNotFound(err) {
			err := gtserror.Newf("storage error getting emoji %s after recache: %w", emoji.URI, err)
			return nil, gtserror.NewErrorInternalError(err)
		} else if rc == nil {
			return nil, gtserror.NewfWithCode(http.StatusNotFound,
				"remote emoji image not found: %s", emoji.URI)
		}
	}

	// If running on S3 storage with proxying disabled,
	// just fetch a pre-signed URL instead of the content.
	if url := p.state.Storage.URL(ctx, emojiPath); url != nil {
		_ = rc.Close() // close storage stream
		content.URL = url
		return &content, nil
	}

	// Return with stream.
	content.Content = rc
	return &content, nil
}

// handles serving Content for "unknown" file
// type, ie., a file we couldn't cache (this time).
func handleUnknown(
	attach *gtsmodel.MediaAttachment,
) (*apimodel.Content, gtserror.WithCode) {
	if attach.RemoteURL == "" {
		err := gtserror.Newf("empty remote url for %s", attach.ID)
		return nil, gtserror.NewErrorInternalError(err)
	}

	// Parse media remote URL to valid URL object.
	remoteURL, err := url.Parse(attach.RemoteURL)
	if err != nil {
		err := gtserror.Newf("invalid remote url for %s: %w", attach.ID, err)
		return nil, gtserror.NewErrorInternalError(err)
	}

	if remoteURL == nil {
		err := gtserror.Newf("nil remote url for %s", attach.ID)
		return nil, gtserror.NewErrorInternalError(err)
	}

	// Just forward the request to the remote URL,
	// since this is a type we couldn't process.
	url := &storage.PresignedURL{
		URL: remoteURL,

		// We might manage to cache the media
		// at some point, so set a low-ish expiry.
		Expiry: time.Now().Add(2 * time.Hour),
	}

	return &apimodel.Content{URL: url}, nil
}

func parseType(s string) (media.Type, error) {
	switch s {
	case string(media.TypeAttachment):
		return media.TypeAttachment, nil
	case string(media.TypeHeader):
		return media.TypeHeader, nil
	case string(media.TypeAvatar):
		return media.TypeAvatar, nil
	case string(media.TypeEmoji):
		return media.TypeEmoji, nil
	}
	return "", fmt.Errorf("%s not a recognized media.Type", s)
}

func parseSize(s string) (media.Size, error) {
	switch s {
	case string(media.SizeSmall):
		return media.SizeSmall, nil
	case string(media.SizeOriginal):
		return media.SizeOriginal, nil
	case string(media.SizeStatic):
		return media.SizeStatic, nil
	}
	return "", fmt.Errorf("%s not a recognized media.Size", s)
}

// Extract the mediaID and file extension from
// a string like "01J3CTH8CZ6ATDNMG6CPRC36XE.gif"
func parseFileName(s string) (string, string, error) {
	spl := strings.Split(s, ".")
	if len(spl) != 2 || spl[0] == "" || spl[1] == "" {
		return "", "", errors.New("file name not splittable on '.'")
	}

	var (
		mediaID  = spl[0]
		mediaExt = spl[1]
	)

	if !regexes.ULID.MatchString(mediaID) {
		return "", "", fmt.Errorf("%s not a valid ULID", mediaID)
	}

	return mediaID, mediaExt, nil
}
