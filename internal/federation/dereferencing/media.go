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
	"io"
	"net/url"

	"code.superseriousbusiness.org/gopkg/log"
	"code.superseriousbusiness.org/gotosocial/internal/config"
	"code.superseriousbusiness.org/gotosocial/internal/gtserror"
	"code.superseriousbusiness.org/gotosocial/internal/gtsmodel"
	"code.superseriousbusiness.org/gotosocial/internal/media"
	"code.superseriousbusiness.org/gotosocial/internal/util"
)

// GetMedia fetches the media at given remote URL by
// dereferencing it. The passed accountID is used to
// store it as being owned by that account. Additional
// information to set on the media attachment may also
// be provided.
//
// If RejectMedia is set to true on the given
// info struct, the media will only be stubbed.
//
// Please note that even if an error is returned,
// a media model may still be returned if the error
// was only encountered during actual dereferencing.
// In this case, it will act as a placeholder.
//
// Also note that since account / status dereferencing is
// already protected by per-uri locks, and that fediverse
// media is generally not shared between accounts (etc),
// there aren't any concurrency protections against multiple
// insertion / dereferencing of media at remoteURL. Worst
// case scenario, an extra media entry will be inserted
// and the scheduled cleaner.Cleaner{} will catch it!
func (d *Dereferencer) GetMedia(
	ctx context.Context,
	requester *gtsmodel.Account,
	accountID string, // media account owner
	remoteURL string,
	info media.AdditionalMediaInfo,
	async bool,
) (
	*gtsmodel.MediaAttachment,
	error,
) {
	// Ensure we have a valid remote URL.
	url, err := url.Parse(remoteURL)
	if err != nil {
		err := gtserror.Newf("invalid media remote url %s: %w", remoteURL, err)
		return nil, err
	}

	return d.processMediaSafely(ctx,
		remoteURL,
		func() (*media.ProcessingMedia, error) {

			// Fetch transport for the provided request user from controller.
			tsport, err := d.transportController.NewTransportForUsername(ctx,
				accountUsername(requester))
			if err != nil {
				return nil, gtserror.Newf("failed getting transport for %s: %w", accountUsername(requester), err)
			}

			// Get maximum supported remote media size.
			maxsz := int64(config.GetMediaRemoteMaxSize()) // #nosec G115 -- Already validated.

			// Prepare data function to dereference remote media.
			data := func(context.Context) (io.ReadCloser, error) {
				return tsport.DereferenceMedia(ctx, url, maxsz)
			}

			// Create media with prepared info.
			return d.mediaManager.CreateMedia(ctx,
				accountID,
				data,
				info,
			)
		},
		async,
	)
}

// RefreshMedia ensures that given media is up-to-date,
// both in terms of being cached in local instance,
// storage and compared to extra info in information
// in given gtsmodel.AdditionMediaInfo{}. This handles
// the case of local emoji by returning early.
//
// If RejectMedia is set to true on the given
// info struct, the media will only be stubbed.
//
// Please note that even if an error is returned,
// a media model may still be returned if the error
// was only encountered during actual dereferencing.
// In this case, it will act as a placeholder.
//
// Also note that since account / status dereferencing is
// already protected by per-uri locks, and that fediverse
// media is generally not shared between accounts (etc),
// there aren't any concurrency protections against multiple
// insertion / dereferencing of media at remoteURL. Worst
// case scenario, an extra media entry will be inserted
// and the scheduled cleaner.Cleaner{} will catch it!
func (d *Dereferencer) RefreshMedia(
	ctx context.Context,
	requester *gtsmodel.Account,
	attach *gtsmodel.MediaAttachment,
	info media.AdditionalMediaInfo,
	force bool,
	async bool,
) (
	*gtsmodel.MediaAttachment,
	error,
) {
	// Can't refresh local.
	if attach.IsLocal() {
		return attach, nil
	}

	// Check if media is up-to-date.
	if !mediaUpToDate(attach, info) {
		force = true
	}

	switch {
	case force:
		// Unset any previous error
		// to force a dereference.
		attach.Error = 0

	case attach.Cached() || !attach.Error.SupportsRetry():
		// Return early, is already cached or
		// error that does not support retry.
		return attach, nil
	}

	// Pass along for safe processing.
	return d.processMediaSafely(ctx,
		attach.RemoteURL,
		func() (*media.ProcessingMedia, error) {

			// Ensure we have a valid remote URL.
			url, err := url.Parse(attach.RemoteURL)
			if err != nil {
				return nil, gtserror.Newf("invalid media remote url %s: %w", attach.RemoteURL, err)
			}

			// Fetch transport for the provided request user from controller.
			tsport, err := d.transportController.NewTransportForUsername(ctx,
				accountUsername(requester))
			if err != nil {
				return nil, gtserror.Newf("failed getting transport for %s: %w", accountUsername(requester), err)
			}

			// Get maximum supported remote media size.
			maxsz := int64(config.GetMediaRemoteMaxSize()) // #nosec G115 -- Already validated.

			// Prepare data function to dereference remote media.
			data := func(context.Context) (io.ReadCloser, error) {
				return tsport.DereferenceMedia(ctx, url, maxsz)
			}

			// Recache media with prepared info,
			// this will also update media in db.
			return d.mediaManager.CacheMedia(
				attach,
				data,
				info,
			), nil
		},
		async,
	)
}

// WaitOnMedia blocks until all given media have finished processing, finally updating the slice model ptrs.
func (d *Dereferencer) WaitOnMedia(ctx context.Context, medias []*gtsmodel.MediaAttachment) {
	type uncachedMedia struct {

		// Ptr to currently processing
		// media to block on, if any.
		Ptr *media.ProcessingMedia

		// Database ID.
		ID string

		// Remote URL key.
		URL string

		// Index in status
		// medias slice.
		Idx int
	}

	// Append media in slice that isn't yet cached.
	uncached := make([]uncachedMedia, 0, len(medias))
	for i, media := range medias {
		if !media.Cached() {
			uncached = append(uncached, uncachedMedia{
				ID:  media.ID,
				URL: media.RemoteURL,
				Idx: i,
			})
		}
	}

	// Check if any uncached.
	if len(uncached) == 0 {
		return
	}

	// Acquire processing
	// medias slice lock.
	d.derefMediaMu.Lock()

	for i, entry := range medias {
		// Check for processing media by remote URL.
		processing := d.derefMedia.get(entry.URL)
		uncached[i].Ptr = processing
	}

	// Done with mutex lock.
	d.derefMediaMu.Unlock()

	for _, entry := range uncached {
		if entry.Ptr != nil {
			// If media was processing, block
			// until finished loading. We don't
			// care about error return as async
			// thread will handle logging it,
			// and media is always non-nil.
			media, _ := entry.Ptr.Load(ctx)

			// Set latest attachment.
			medias[entry.Idx] = media
		} else {

			// Media had finished processing, get latest from database.
			media, err := d.state.DB.GetAttachmentByID(ctx, entry.ID)
			if err != nil {
				log.Errorf(ctx, "error getting latest attachment %s: %v", entry.URL, err)
				continue
			}

			// Set latest attachment.
			medias[entry.Idx] = media
		}
	}
}

// processingMediaSafely provides concurrency-safe processing of
// a media with given remote URL string. if a copy of the media is
// not already being processed, the given 'process' callback will
// be used to generate new *media.ProcessingMedia{} instance. async
// determines whether to load it immediately, or in the background.
// the provided process function can also optionally return ready
// media model directly for catching just-cached race conditions.
func (d *Dereferencer) processMediaSafely(
	ctx context.Context,
	remoteURL string,
	process func() (*media.ProcessingMedia, error),
	async bool,
) (
	attach *gtsmodel.MediaAttachment,
	err error,
) {
	var existing bool

	// Acquire map lock.
	d.derefMediaMu.Lock()

	// Ensure unlock only done once.
	unlock := d.derefMediaMu.Unlock
	unlock = util.DoOnce(unlock)
	defer unlock()

	// Look for an existing deref in progress.
	processing := d.derefMedia.get(remoteURL)
	if existing = (processing != nil); !existing {

		// Start new processing media.
		processing, err = process()
		if err != nil {
			return nil, err
		}

		// Add processing media to the list.
		d.derefMedia.put(remoteURL, processing)
	}

	// Unlock map.
	unlock()

	if async {
		// Acquire a placeholder to return.
		attach = processing.Placeholder()

		// Enqueue the processing media load logic for background processing.
		d.state.Workers.Dereference.Queue.Push(func(ctx context.Context) {
			if !existing {
				defer func() {
					// We started the processing,
					// remove from list on finish.
					d.derefMediaMu.Lock()
					d.derefMedia.delete(remoteURL)
					d.derefMediaMu.Unlock()
				}()
			}

			// Perform media load operation.
			attach, err = processing.Load(ctx)
			if err != nil {
				log.Errorf(ctx, "error loading media %s: %v", remoteURL, err)
			}

			// Pass to its dereferencer hook.
			d.onMediaDereference(ctx, attach)
		})
	} else {
		if !existing {
			defer func() {
				// We started the processing,
				// remove from list on finish.
				d.derefMediaMu.Lock()
				d.derefMedia.delete(remoteURL)
				d.derefMediaMu.Unlock()
			}()
		}

		// Perform media load operation,
		// falling back to asynchronous
		// operation on context cancelled.
		attach, err = processing.MustLoad(ctx)
		if err != nil {

			// TODO: in time we should return checkable flags by gtserror.Is___()
			// which can determine if loading error should allow remaining placeholder.
			err = gtserror.Newf("error loading media %s: %w", remoteURL, err)
		}

		// Pass to its dereferencer hook.
		d.onMediaDereference(ctx, attach)
	}

	return
}
