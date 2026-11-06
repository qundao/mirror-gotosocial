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

package admin

import (
	"context"
	"time"

	"code.superseriousbusiness.org/gopkg/log"
	"code.superseriousbusiness.org/gotosocial/internal/gtscontext"
	"code.superseriousbusiness.org/gotosocial/internal/gtserror"
	"code.superseriousbusiness.org/gotosocial/internal/gtsmodel"
	"codeberg.org/gruf/go-longdur"
)

// MediaRefetch forces a refetch of remote emojis.
func (p *Processor) MediaRefetch(ctx context.Context, requester *gtsmodel.Account, domain string) {
	go func() {
		// Emoji media should always use instance transport.
		tsport := p.transport.GetTransportForAccount(nil)

		// Prepare background context with values of requester.
		ctx := gtscontext.WithValues(context.Background(), ctx)

		log.Info(ctx, "starting emoji refetch")
		refetched, err := p.media.RefetchEmojis(ctx, domain, tsport.DereferenceMedia)
		if err != nil {
			log.Errorf(ctx, "error refetching emojis: %s", err)
		} else {
			log.Infof(ctx, "refetched %d emojis from remote", refetched)
		}
	}()
}

// MediaPrune triggers a non-blocking prune of unused media, orphaned, uncaching remote and fixing cache states.
func (p *Processor) MediaPrune(ctx context.Context, remoteCacheAge longdur.Duration) gtserror.WithCode {

	// Start background task
	// performing media cleanup.
	go func() {
		now := time.Now()
		ctx := gtscontext.WithValues(context.Background(), ctx)
		p.cleaner.Media().AllAndFix(ctx, now, remoteCacheAge)
		p.cleaner.Emoji().AllAndFix(ctx, now, remoteCacheAge)
	}()

	return nil
}

// MediaPurge triggers a non-blocking purge of all
// media attachments + emojis from the given domain.
func (p *Processor) MediaPurge(
	ctx context.Context,
	domain string,
) gtserror.WithCode {

	// Start background task
	// performing media purge.
	go func() {
		ctx := gtscontext.WithValues(context.Background(), ctx)
		p.cleaner.Media().LogPurgeRemote(ctx, domain)
		p.cleaner.Emoji().LogPurgeRemote(ctx, domain)
	}()

	return nil
}
