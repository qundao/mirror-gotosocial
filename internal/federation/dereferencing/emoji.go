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
	"io"
	"net/url"

	"code.superseriousbusiness.org/gopkg/log"
	"code.superseriousbusiness.org/gotosocial/internal/config"
	"code.superseriousbusiness.org/gotosocial/internal/db"
	"code.superseriousbusiness.org/gotosocial/internal/gtserror"
	"code.superseriousbusiness.org/gotosocial/internal/gtsmodel"
	"code.superseriousbusiness.org/gotosocial/internal/media"
	"code.superseriousbusiness.org/gotosocial/internal/util"
)

// GetEmoji fetches the emoji with given shortcode,
// domain and remote URL to dereference it by. This
// handles the case of existing emojis by passing them
// to RefreshEmoji(), which in the case of a local
// emoji will be a no-op. If the emoji does not yet
// exist it will be newly inserted into the database
// followed by dereferencing the actual media file.
//
// If RejectMedia is set to true on the given
// info struct, the media will only be stubbed.
//
// Please note that even if an error is returned,
// an emoji model may still be returned if the error
// was only encountered during actual dereferencing.
// In this case, it will act as a placeholder.
func (d *Dereferencer) GetEmoji(
	ctx context.Context,
	shortcode string,
	domain string,
	remoteURL string,
	info media.AdditionalEmojiInfo,
	refresh bool,
	async bool,
) (
	*gtsmodel.Emoji,
	error,
) {
	// Look for an existing emoji with shortcode domain.
	emoji, err := d.state.DB.GetEmojiByShortcodeDomain(ctx,
		shortcode,
		domain,
	)
	if err != nil && !errors.Is(err, db.ErrNoEntries) {
		return nil, gtserror.Newf("error fetching emoji from db: %w", err)
	}

	if emoji != nil {
		// This was an existing emoji, pass to refresh func.
		return d.RefreshEmoji(ctx, emoji, info, refresh, async)
	}

	if domain == "" {
		// failed local lookup, will be db.ErrNoEntries.
		return nil, gtserror.SetUnretrievable(err)
	}

	// Get maximum supported remote emoji size.
	maxsz := int64(config.GetMediaEmojiRemoteMaxSize()) // #nosec G115 -- Already validated.

	// Generate shortcode@domain for locks + logging.
	shortcodeDomain := shortcode + "@" + domain

	// Pass along for safe processing.
	return d.processEmojiSafely(ctx,
		shortcodeDomain,
		func() (*media.ProcessingEmoji, *gtsmodel.Emoji, error) {

			// Reload emoji within lock to check if *just* cached.
			emoji, err := d.state.DB.GetEmojiByShortcodeDomain(ctx,
				shortcode,
				domain,
			)
			if err != nil && !errors.Is(err, db.ErrNoEntries) {
				return nil, nil, gtserror.Newf("error fetching emoji from db: %w", err)
			}

			if emoji != nil {
				// *just* fetched.
				return nil, emoji, nil
			}

			// Ensure we have a valid remote URL.
			url, err := url.Parse(remoteURL)
			if err != nil {
				err := gtserror.Newf("invalid image remote url %s for emoji %s@%s: %w", remoteURL, shortcode, domain, err)
				return nil, nil, err
			}

			// Acquire new instance account transport for emoji dereferencing.
			tsport, err := d.transportController.NewTransportForUsername(ctx, "")
			if err != nil {
				err := gtserror.Newf("error getting instance transport: %w", err)
				return nil, nil, err
			}

			// Prepare data function to dereference remote emoji media.
			data := func(context.Context) (io.ReadCloser, error) {
				return tsport.DereferenceMedia(ctx, url, maxsz)
			}

			// Start emoji create operation with prepared info.
			processing, err := d.mediaManager.CreateEmoji(ctx,
				shortcode,
				domain,
				data,
				info,
			)
			return processing, nil, err
		},
		async,
	)
}

// RefreshEmoji ensures that the given emoji is
// up-to-date, both in terms of being cached in
// in local instance storage, and compared to extra
// information provided in media.AdditionEmojiInfo{}.
// (note that is a no-op to pass in a local emoji).
//
// If RejectMedia is set to true on the given
// info struct, the emoji will only be stubbed.
//
// Please note that even if an error is returned,
// an emoji model may still be returned if the error
// was only encountered during actual dereferencing.
// In this case, it will act as a placeholder.
func (d *Dereferencer) RefreshEmoji(
	ctx context.Context,
	emoji *gtsmodel.Emoji,
	info media.AdditionalEmojiInfo,
	force bool,
	async bool,
) (
	*gtsmodel.Emoji,
	error,
) {
	// Can't refresh local.
	if emoji.IsLocal() {
		return emoji, nil
	}

	// Check if emoji is up-to-date.
	if !emojiUpToDate(emoji, info) {
		force = true
	}

	if !force {
		// We still want to make sure
		// the emoji is cached. Simply
		// check whether emoji is cached.
		return d.RecacheEmoji(ctx, emoji, info, async)
	}

	// Generate emoji shortcode@domain
	// for locks + logging when processing.
	shortcodeDomain := emoji.ShortcodeDomain()

	// Ensure we have a valid image remote URL.
	url, err := url.Parse(emoji.ImageRemoteURL)
	if err != nil {
		err := gtserror.Newf("invalid image remote url %s for emoji %s: %w", emoji.ImageRemoteURL, shortcodeDomain, err)
		return nil, err
	}

	// Pass along for safe processing.
	return d.processEmojiSafely(ctx,
		shortcodeDomain,
		func() (*media.ProcessingEmoji, *gtsmodel.Emoji, error) {

			// Reload emoji within lock to check if *just* refreshed.
			emoji, err := d.state.DB.GetEmojiByID(ctx, emoji.ID)
			if err != nil {
				return nil, nil, gtserror.Newf("error fetching emoji from db: %w", err)
			}

			// Check if still needs refresh.
			if emojiUpToDate(emoji, info) {
				return nil, emoji, nil
			}

			// Acquire new instance account transport for emoji dereferencing.
			tsport, err := d.transportController.NewTransportForUsername(ctx, "")
			if err != nil {
				err := gtserror.Newf("error getting instance transport: %w", err)
				return nil, nil, err
			}

			// Get maximum supported remote emoji size.
			maxsz := int64(config.GetMediaEmojiRemoteMaxSize()) // #nosec G115 -- Already validated.

			// Prepare data function to dereference remote emoji media.
			data := func(context.Context) (io.ReadCloser, error) {
				return tsport.DereferenceMedia(ctx, url, maxsz)
			}

			// Start emoji update operation with prepared info.
			processing, err := d.mediaManager.UpdateEmoji(ctx,
				emoji,
				data,
				info,
			)
			return processing, nil, err
		},
		async,
	)
}

// RecacheEmoji handles the simplest case which is that
// of an existing emoji that only needs to be recached.
// It handles the case of both local emojis, and those
// already cached as no-ops.
//
// Please note that even if an error is returned,
// an emoji model may still be returned if the error
// was only encountered during actual dereferencing.
// In this case, it will act as a placeholder.
func (d *Dereferencer) RecacheEmoji(
	ctx context.Context,
	emoji *gtsmodel.Emoji,
	info media.AdditionalEmojiInfo,
	async bool,
) (
	*gtsmodel.Emoji,
	error,
) {
	// Can't recache local.
	if emoji.IsLocal() {
		return emoji, nil
	}

	if emoji.Cached() {
		// Already cached.
		return emoji, nil
	}

	// Generate emoji shortcode@domain
	// for locks + logging when processing.
	shortcodeDomain := emoji.ShortcodeDomain()

	// Ensure we have a valid image remote URL.
	url, err := url.Parse(emoji.ImageRemoteURL)
	if err != nil {
		err := gtserror.Newf("invalid image remote url %s for emoji %s: %w", emoji.ImageRemoteURL, shortcodeDomain, err)
		return nil, err
	}

	// Pass along for safe processing.
	return d.processEmojiSafely(ctx,
		shortcodeDomain,
		func() (*media.ProcessingEmoji, *gtsmodel.Emoji, error) {

			// Reload emoji within lock to check if *just* cached.
			emoji, err := d.state.DB.GetEmojiByID(ctx, emoji.ID)
			if err != nil && !errors.Is(err, db.ErrNoEntries) {
				return nil, nil, gtserror.Newf("error fetching emoji from db: %w", err)
			}

			if emoji != nil && emoji.Cached() {
				// This was *just* cached.
				return nil, emoji, nil
			}

			// Acquire new instance account transport for emoji dereferencing.
			tsport, err := d.transportController.NewTransportForUsername(ctx, "")
			if err != nil {
				err := gtserror.Newf("error getting instance transport: %w", err)
				return nil, nil, err
			}

			// Get maximum supported remote emoji size.
			maxsz := int64(config.GetMediaEmojiRemoteMaxSize()) // #nosec G115 -- Already validated.

			// Prepare data function to dereference remote emoji media.
			data := func(context.Context) (io.ReadCloser, error) {
				return tsport.DereferenceMedia(ctx, url, maxsz)
			}

			// Start emoji recache operation with prepared info.
			processing, err := d.mediaManager.CacheEmoji(ctx,
				emoji,
				data,
				info,
			)
			return processing, nil, err
		},
		async,
	)
}

// WaitOnEmoji blocks until all given emoji have finished processing, finally updating the slice model ptrs.
func (d *Dereferencer) WaitOnEmoji(ctx context.Context, emojis []*gtsmodel.Emoji) {
	type uncachedEmoji struct {

		// Ptr to currently processing
		// emoji to block on, if any.
		Ptr *media.ProcessingEmoji

		// Database ID.
		ID string

		// shortcode@dommain key.
		ShortcodeDomain string

		// Index in status
		// attachment slice.
		Idx int
	}

	// Append emoji in status emoji that isn't yet cached.
	uncached := make([]uncachedEmoji, 0, len(emojis))
	for i, emoji := range emojis {
		if !emoji.Cached() {
			uncached = append(uncached, uncachedEmoji{
				ID:              emoji.ID,
				ShortcodeDomain: emoji.ShortcodeDomain(),
				Idx:             i,
			})
		}
	}

	// Check if any uncached.
	if len(uncached) == 0 {
		return
	}

	// Acquire processing
	// emojis slice lock.
	d.derefEmojisMu.Lock()

	for i, entry := range uncached {
		// Check for processing emoji by its shortcode@domain.
		processing := d.derefEmojis.get(entry.ShortcodeDomain)
		uncached[i].Ptr = processing
	}

	// Done with mutex lock.
	d.derefEmojisMu.Unlock()

	for _, entry := range uncached {
		if entry.Ptr != nil {
			// If emoji was processing, block
			// until finished loading. We don't
			// care about error return as async
			// thread will handle logging it,
			// and emoji is always non-nil.
			emoji, _ := entry.Ptr.Load(ctx)

			// Set latest emoji model.
			emojis[entry.Idx] = emoji
		} else {

			// Emoji had finished processing, get latest from db.
			emoji, err := d.state.DB.GetEmojiByID(ctx, entry.ID)
			if err != nil {
				log.Errorf(ctx, "error getting latest emoji %s: %v", entry.ShortcodeDomain, err)
				continue
			}

			// Set latest emoji model.
			emojis[entry.Idx] = emoji
		}
	}
}

// processingEmojiSafely provides concurrency-safe processing of
// an emoji with given shortcode+domain. if a copy of the emoji is
// not already being processed, the given 'process' callback will
// be used to generate new *media.ProcessingEmoji{} instance. async
// determines whether to load it immediately, or in the background.
// the provided process function can also optionally return ready
// emoji model directly for catching just-cached race conditions.
func (d *Dereferencer) processEmojiSafely(
	ctx context.Context,
	shortcodeDomain string,
	process func() (*media.ProcessingEmoji, *gtsmodel.Emoji, error),
	async bool,
) (
	emoji *gtsmodel.Emoji,
	err error,
) {
	var existing bool

	// Acquire map lock.
	d.derefEmojisMu.Lock()

	// Ensure unlock only done once.
	unlock := d.derefEmojisMu.Unlock
	unlock = util.DoOnce(unlock)
	defer unlock()

	// Look for an existing dereference in progress.
	processing := d.derefEmojis.get(shortcodeDomain)
	if existing = (processing != nil); !existing {
		var emoji *gtsmodel.Emoji

		// Start new processing of emoji.
		processing, emoji, err = process()
		if err != nil {
			return nil, err
		}

		if emoji != nil {
			// The caller must have caught
			// a race condition and found the
			// emoji model was *just* cached.
			return emoji, nil
		}

		// Add processing emoji media entry to the list.
		d.derefEmojis.put(shortcodeDomain, processing)
	}

	// Unlock map.
	unlock()

	if async {
		// Acquire a placeholder to return.
		emoji = processing.Placeholder()

		// Enqueue the processing emoji load logic for background processing.
		d.state.Workers.Dereference.Queue.Push(func(ctx context.Context) {
			if !existing {
				defer func() {
					// We started the processing,
					// remove from list on finish.
					d.derefEmojisMu.Lock()
					d.derefEmojis.delete(shortcodeDomain)
					d.derefEmojisMu.Unlock()
				}()
			}

			// Perform emoji load operation.
			emoji, err = processing.Load(ctx)
			if err != nil {
				log.Errorf(ctx, "error loading emoji %s: %v", shortcodeDomain, err)
			}

			// Pass to its dereferencer hook.
			d.onEmojiDereference(ctx, emoji)
		})
	} else {
		if !existing {
			defer func() {
				// We started the processing,
				// remove from list on finish.
				d.derefEmojisMu.Lock()
				d.derefEmojis.delete(shortcodeDomain)
				d.derefEmojisMu.Unlock()
			}()
		}

		// Perform emoji load operation,
		// falling back to asynchronous
		// operation on context cancelled.
		emoji, err = processing.MustLoad(ctx)
		if err != nil {

			// TODO: in time we should return checkable flags by gtserror.Is___()
			// which can determine if loading error should allow remaining placeholder.
			err = gtserror.Newf("error loading emoji %s: %w", shortcodeDomain, err)
		}

		// Pass to its dereferencer hook.
		d.onEmojiDereference(ctx, emoji)
	}

	return
}

func (d *Dereferencer) fetchEmojis(
	ctx context.Context,
	existing []*gtsmodel.Emoji,
	emojis []*gtsmodel.Emoji, // newly dereferenced
	rejectReason *gtsmodel.MediaErrorDetails, // optional reason to reject media with
) (
	[]*gtsmodel.Emoji,
	bool, // any changes?
	error,
) {
	// Track any changes.
	changed := false

	for i, placeholder := range emojis {
		// Look for an existing emoji with shortcode + domain.
		existing, ok := getEmojiByShortcodeDomain(existing,
			placeholder.Shortcode,
			placeholder.Domain,
		)
		if ok && existing.ID != "" {

			// Check for any emoji changes that
			// indicate we should force a refresh.
			force := emojiChanged(existing, placeholder)

			// Set latest values from placeholder.
			info := media.AdditionalEmojiInfo{
				URI:                  &placeholder.URI,
				ImageRemoteURL:       &placeholder.ImageRemoteURL,
				ImageStaticRemoteURL: &placeholder.ImageStaticRemoteURL,

				// Pass reject reason ptr, which
				// will skip downloading if set.
				RejectReason: rejectReason,
			}

			// Ensure that the existing emoji
			// model is up-to-date and cached.
			existing, err := d.RefreshEmoji(
				ctx,
				existing,
				info,
				force,
				true, // async
			)
			if err != nil {
				log.Errorf(ctx, "error refreshing emoji: %v", err)

				// specifically do NOT continue here,
				// we already have a model, we don't
				// want to drop it from the slice, just
				// log that an update for it failed.
			}

			// Set existing emoji.
			emojis[i] = existing
			continue
		}

		// Emojis changed!
		changed = true

		// Prepare emoji info, including the
		// rejectMedia flag if necessary.
		info := media.AdditionalEmojiInfo{
			URI:                  &placeholder.URI,
			ImageRemoteURL:       &placeholder.ImageRemoteURL,
			ImageStaticRemoteURL: &placeholder.ImageStaticRemoteURL,

			// Pass reject reason ptr, which
			// will skip downloading if set.
			RejectReason: rejectReason,
		}

		// Fetch this newly added emoji,
		// this function handles the case
		// of existing cached emojis and
		// new ones requiring dereference.
		emoji, err := d.GetEmoji(ctx,
			placeholder.Shortcode,
			placeholder.Domain,
			placeholder.ImageRemoteURL,
			info,
			false, // refresh
			true,  // async
		)
		if err != nil {
			if emoji == nil {
				log.Errorf(ctx, "error loading emoji %s: %v", placeholder.ImageRemoteURL, err)
				continue
			}

			// non-fatal error occurred during loading, still use it.
			log.Warnf(ctx, "partially loaded emoji: %v", err)
		}

		// Set updated emoji.
		emojis[i] = emoji
	}

	for i := 0; i < len(emojis); {
		if emojis[i].ID == "" {
			// Remove failed emoji populations.
			copy(emojis[i:], emojis[i+1:])
			emojis = emojis[:len(emojis)-1]
			continue
		}
		i++
	}

	return emojis, changed, nil
}
