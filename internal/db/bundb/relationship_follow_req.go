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

package bundb

import (
	"context"
	"errors"
	"slices"
	"time"

	"code.superseriousbusiness.org/gopkg/log"
	"code.superseriousbusiness.org/gopkg/xslices"
	"code.superseriousbusiness.org/gotosocial/internal/db"
	"code.superseriousbusiness.org/gotosocial/internal/gtscontext"
	"code.superseriousbusiness.org/gotosocial/internal/gtserror"
	"code.superseriousbusiness.org/gotosocial/internal/gtsmodel"
	"github.com/uptrace/bun"
)

func (r *relationshipDB) GetFollowRequestByID(ctx context.Context, id string) (*gtsmodel.FollowRequest, error) {
	return r.getFollowRequest(
		ctx,
		"ID",
		func(followReq *gtsmodel.FollowRequest) error {
			return r.db.NewSelect().
				Model(followReq).
				Where("? = ?", bun.Ident("id"), id).
				Scan(ctx)
		},
		id,
	)
}

func (r *relationshipDB) GetFollowRequestByURI(ctx context.Context, uri string) (*gtsmodel.FollowRequest, error) {
	return r.getFollowRequest(
		ctx,
		"URI",
		func(followReq *gtsmodel.FollowRequest) error {
			return r.db.NewSelect().
				Model(followReq).
				Where("? = ?", bun.Ident("uri"), uri).
				Scan(ctx)
		},
		uri,
	)
}

func (r *relationshipDB) GetFollowRequest(ctx context.Context, sourceAccountID string, targetAccountID string) (*gtsmodel.FollowRequest, error) {
	return r.getFollowRequest(
		ctx,
		"AccountID,TargetAccountID",
		func(followReq *gtsmodel.FollowRequest) error {
			return r.db.NewSelect().
				Model(followReq).
				Where("? = ?", bun.Ident("account_id"), sourceAccountID).
				Where("? = ?", bun.Ident("target_account_id"), targetAccountID).
				Scan(ctx)
		},
		sourceAccountID,
		targetAccountID,
	)
}

func (r *relationshipDB) GetFollowRequestsByIDs(ctx context.Context, ids []string) ([]*gtsmodel.FollowRequest, error) {
	// Load all follow IDs via cache loader callbacks.
	follows, err := r.state.Caches.DB.FollowRequest.LoadIDs("ID",
		ids,
		func(uncached []string) ([]*gtsmodel.FollowRequest, error) {
			// Preallocate expected length of uncached followReqs.
			follows := make([]*gtsmodel.FollowRequest, 0, len(uncached))

			// Perform database query scanning
			// the remaining (uncached) IDs.
			if err := r.db.NewSelect().
				Model(&follows).
				Where("? IN (?)", bun.Ident("id"), bun.In(uncached)).
				Scan(ctx); err != nil {
				return nil, err
			}

			return follows, nil
		},
	)
	if err != nil {
		return nil, err
	}

	// Reorder the requests by their
	// IDs to ensure in correct order.
	getID := func(f *gtsmodel.FollowRequest) string { return f.ID }
	xslices.OrderBy(follows, ids, getID)

	if gtscontext.Barebones(ctx) {
		// no need to fully populate.
		return follows, nil
	}

	// Populate all loaded followreqs, removing those we fail to
	// populate (removes needing so many nil checks everywhere).
	follows = slices.DeleteFunc(follows, func(follow *gtsmodel.FollowRequest) bool {
		if err := r.PopulateFollowRequest(ctx, follow); err != nil {
			log.Errorf(ctx, "error populating follow request %s: %v", follow.ID, err)
			return true
		}
		return false
	})

	return follows, nil
}

func (r *relationshipDB) IsFollowRequested(ctx context.Context, sourceAccountID string, targetAccountID string) (bool, error) {
	followReq, err := r.GetFollowRequest(
		gtscontext.SetBarebones(ctx),
		sourceAccountID,
		targetAccountID,
	)
	if err != nil && !errors.Is(err, db.ErrNoEntries) {
		return false, err
	}
	return (followReq != nil), nil
}

func (r *relationshipDB) getFollowRequest(ctx context.Context, lookup string, dbQuery func(*gtsmodel.FollowRequest) error, keyParts ...any) (*gtsmodel.FollowRequest, error) {
	// Fetch follow request from database cache with loader callback
	followReq, err := r.state.Caches.DB.FollowRequest.LoadOne(lookup, func() (*gtsmodel.FollowRequest, error) {
		var followReq gtsmodel.FollowRequest

		// Not cached! Perform database query
		if err := dbQuery(&followReq); err != nil {
			return nil, err
		}

		return &followReq, nil
	}, keyParts...)
	if err != nil {
		// error already processed
		return nil, err
	}

	if gtscontext.Barebones(ctx) {
		// Only a barebones model was requested.
		return followReq, nil
	}

	if err := r.state.DB.PopulateFollowRequest(ctx, followReq); err != nil {
		return nil, err
	}

	return followReq, nil
}

func (r *relationshipDB) PopulateFollowRequest(ctx context.Context, follow *gtsmodel.FollowRequest) error {
	var (
		err  error
		errs gtserror.MultiError
	)

	if follow.Account == nil {
		// Follow account is not set, fetch from the database.
		follow.Account, err = r.state.DB.GetAccountByID(
			gtscontext.SetBarebones(ctx),
			follow.AccountID,
		)
		if err != nil {
			errs.Appendf("error populating follow request account: %w", err)
		}
	}

	if follow.TargetAccount == nil {
		// Follow target account is not set, fetch from the database.
		follow.TargetAccount, err = r.state.DB.GetAccountByID(
			gtscontext.SetBarebones(ctx),
			follow.TargetAccountID,
		)
		if err != nil {
			errs.Appendf("error populating follow target request account: %w", err)
		}
	}

	return errs.Combine()
}

func (r *relationshipDB) PutFollowRequest(ctx context.Context, follow *gtsmodel.FollowRequest) error {
	return r.insertFollowRequest(ctx, follow, func(tx bun.Tx) error {
		_, err := tx.NewInsert().
			Model(follow).
			Exec(ctx)
		return err
	})
}

func (r *relationshipDB) UpdateFollowRequest(ctx context.Context, followRequest *gtsmodel.FollowRequest, columns ...string) error {
	followRequest.UpdatedAt = time.Now()
	if len(columns) > 0 {
		// If we're updating by column, ensure "updated_at" is included.
		columns = append(columns, "updated_at")
	}
	return r.state.Caches.DB.FollowRequest.Store(followRequest, func() error {
		_, err := r.db.NewUpdate().
			Model(followRequest).
			Where("? = ?", bun.Ident("follow_request.id"), followRequest.ID).
			Column(columns...).
			Exec(ctx)
		return err
	})
}

func (r *relationshipDB) AcceptFollowRequest(ctx context.Context, sourceAccountID string, targetAccountID string) (*gtsmodel.Follow, error) {
	followReq, err := r.GetFollowRequest(ctx, sourceAccountID, targetAccountID)
	if err != nil {
		return nil, err
	}

	// Create a new follow to 'replace'
	// the original follow request with.
	follow := &gtsmodel.Follow{
		ID:              followReq.ID,
		AccountID:       sourceAccountID,
		Account:         followReq.Account,
		TargetAccountID: targetAccountID,
		TargetAccount:   followReq.TargetAccount,
		URI:             followReq.URI,
		ShowReblogs:     followReq.ShowReblogs,
		Notify:          followReq.Notify,
	}

	// Insert the new follow modelled after request into database.
	if err := r.insertFollow(ctx, follow, func(tx bun.Tx) error {
		_, err := tx.NewInsert().
			Model(follow).
			On("CONFLICT (?,?) DO UPDATE set ? = ?", bun.Ident("account_id"), bun.Ident("target_account_id"), bun.Ident("uri"), follow.URI).
			Exec(ctx)
		return err
	}); err != nil {
		return nil, err
	}

	// Delete the follow request now that it's accepted and not needed.
	if err := r.DeleteFollowRequestByID(ctx, followReq.ID); err != nil {
		return nil, err
	}

	// Delete original follow request notification
	if err := r.state.DB.DeleteNotifications(ctx, []gtsmodel.NotificationType{
		gtsmodel.NotificationFollowRequest,
	}, targetAccountID, sourceAccountID); err != nil {
		return nil, err
	}

	return follow, nil
}

func (r *relationshipDB) RejectFollowRequest(ctx context.Context, sourceAccountID string, targetAccountID string) error {
	if err := r.DeleteFollowRequest(ctx, sourceAccountID, targetAccountID); err != nil {
		return err
	}
	return r.state.DB.DeleteNotifications(ctx, []gtsmodel.NotificationType{
		gtsmodel.NotificationFollowRequest,
	}, targetAccountID, sourceAccountID)
}

func (r *relationshipDB) DeleteFollowRequest(
	ctx context.Context,
	sourceAccountID string,
	targetAccountID string,
) error {
	return r.deleteFollowRequest(ctx, func(tx bun.Tx) (*gtsmodel.FollowRequest, error) {
		var deleted gtsmodel.FollowRequest
		deleted.AccountID = sourceAccountID
		deleted.TargetAccountID = targetAccountID

		if _, err := tx.NewDelete().
			Model(&deleted).
			Where("? = ?", bun.Ident("account_id"), sourceAccountID).
			Where("? = ?", bun.Ident("target_account_id"), targetAccountID).
			Returning("?", bun.Ident("id")).
			Exec(ctx); err != nil {
			return nil, err
		}

		return &deleted, nil
	})
}

func (r *relationshipDB) DeleteFollowRequestByID(ctx context.Context, id string) error {
	return r.deleteFollowRequest(ctx, func(tx bun.Tx) (*gtsmodel.FollowRequest, error) {
		var deleted gtsmodel.FollowRequest
		deleted.ID = id

		if _, err := tx.NewDelete().
			Model(&deleted).
			Where("? = ?", bun.Ident("id"), id).
			Returning("?, ?",
				bun.Ident("account_id"),
				bun.Ident("target_account_id"),
			).
			Exec(ctx); err != nil {
			return nil, err
		}

		return &deleted, nil
	})
}

func (r *relationshipDB) DeleteFollowRequestByURI(ctx context.Context, uri string) error {
	return r.deleteFollowRequest(ctx, func(tx bun.Tx) (*gtsmodel.FollowRequest, error) {
		var deleted gtsmodel.FollowRequest
		deleted.URI = uri

		if _, err := tx.NewDelete().
			Model(&deleted).
			Where("? = ?", bun.Ident("uri"), uri).
			Returning("?, ?, ?",
				bun.Ident("id"),
				bun.Ident("account_id"),
				bun.Ident("target_account_id"),
			).
			Exec(ctx); err != nil {
			return nil, err
		}

		return &deleted, nil
	})
}

func (r *relationshipDB) DeleteAccountFollowRequests(ctx context.Context, accountID string) error {
	// Gather necessary fields from
	// deleted for cache invaliation.
	var deleted []*gtsmodel.FollowRequest

	if err := r.db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		// Delete all follows either from
		// account, or targeting account,
		// returning the deleted models.
		if _, err := tx.NewDelete().
			Model(&deleted).
			WhereOr("? = ? OR ? = ?",
				bun.Ident("account_id"),
				accountID,
				bun.Ident("target_account_id"),
				accountID,
			).
			Returning("?, ?, ?",
				bun.Ident("id"),
				bun.Ident("account_id"),
				bun.Ident("target_account_id"),
			).
			Exec(ctx); err != nil {

			// the RETURNING here will cause an ErrNoRows
			// to be returned on DELETE, which is caught
			// outside this RunInTx() func, and ensures we
			// return early here to *not* update statistics.
			return err
		}

		for _, follow := range deleted {
			// Decrement target follow requests count.
			if err := decrementAccountStats(ctx, tx,
				"follow_requests_count",
				follow.TargetAccountID,
			); err != nil {
				return err
			}
		}

		return nil
	}); err != nil && !errors.Is(err, db.ErrNoEntries) {
		return err
	}

	// Invalidate all account's incoming / outoing follows requests.
	r.state.Caches.DB.FollowRequest.Invalidate("AccountID", accountID)
	r.state.Caches.DB.FollowRequest.Invalidate("TargetAccountID", accountID)

	// In case not all follow were in
	// cache, manually call invalidate hooks.
	for _, followReq := range deleted {
		r.state.Caches.OnInvalidateFollowRequest(followReq)
	}

	return nil
}

func (r *relationshipDB) insertFollowRequest(ctx context.Context, follow *gtsmodel.FollowRequest, insert func(bun.Tx) error) error {
	return r.state.Caches.DB.FollowRequest.Store(follow, func() error {
		return r.db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
			// Perform the insert operation.
			if err := insert(tx); err != nil {
				return gtserror.Newf("error inserting follow request: %w", err)
			}

			// Increment target follow requests count.
			return incrementAccountStats(ctx, tx,
				"follow_requests_count",
				follow.TargetAccountID,
			)
		})
	})
}

func (r *relationshipDB) deleteFollowRequest(ctx context.Context, delete func(bun.Tx) (*gtsmodel.FollowRequest, error)) error {
	var follow *gtsmodel.FollowRequest

	if err := r.db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) (err error) {
		// Perform delete operation.
		follow, err = delete(tx)
		if err != nil {

			// the RETURNING here will cause an ErrNoRows
			// to be returned on DELETE, which is caught
			// outside this RunInTx() func, and ensures we
			// return early here to *not* update statistics.
			return err
		}

		// Decrement target follow requests count.
		return decrementAccountStats(ctx, tx,
			"follow_requests_count",
			follow.TargetAccountID,
		)
	}); err != nil && !errors.Is(err, db.ErrNoEntries) {
		return err
	}

	if follow == nil {
		return nil
	}

	// Invalidate cached follow with ID, manually
	// call invalidate hook in case not cached.
	r.state.Caches.DB.FollowRequest.Invalidate("ID", follow.ID)
	r.state.Caches.OnInvalidateFollowRequest(follow)

	return nil
}
