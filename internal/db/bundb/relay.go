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

	"code.superseriousbusiness.org/gopkg/log"
	"code.superseriousbusiness.org/gopkg/xslices"
	"code.superseriousbusiness.org/gotosocial/internal/db"
	"code.superseriousbusiness.org/gotosocial/internal/gtscontext"
	"code.superseriousbusiness.org/gotosocial/internal/gtserror"
	"code.superseriousbusiness.org/gotosocial/internal/gtsmodel"
	"code.superseriousbusiness.org/gotosocial/internal/state"
	"github.com/uptrace/bun"
)

type relayDB struct {
	db    *bun.DB
	state *state.State
}

func (r *relayDB) GetRelayPushByID(ctx context.Context, id string) (*gtsmodel.RelayPush, error) {
	relayPush, err := r.state.Caches.DB.RelayPush.LoadOne(
		"ID",
		func() (*gtsmodel.RelayPush, error) {
			var relayPush gtsmodel.RelayPush
			err := r.db.
				NewSelect().
				Model(&relayPush).
				Where("? = ?", bun.Ident("id"), id).
				Scan(ctx)
			return &relayPush, err
		},
		id,
	)
	if err != nil {
		// already processed
		return nil, err
	}

	if !gtscontext.Barebones(ctx) {
		if err := r.PopulateRelayPush(ctx, relayPush); err != nil {
			return nil, err
		}
	}

	return relayPush, nil
}

func (r *relayDB) GetRelayPushesForAccountID(ctx context.Context, accountID string) ([]*gtsmodel.RelayPush, error) {
	ids, err := r.state.Caches.DB.RelayPushIDs.Load(accountID, func() ([]string, error) {
		var ids []string
		if err := r.db.
			NewSelect().
			Table("relay_pushes").
			Column("id").
			Where("? = ?", bun.Ident("account_id"), accountID).
			Scan(ctx, &ids); err != nil {
			return nil, err
		}
		return ids, nil
	})
	if err != nil {
		return nil, err
	}

	if len(ids) == 0 {
		return nil, nil
	}

	return r.getRelayPushesByIDs(ctx, ids)
}

func (r *relayDB) GetRelayPushesByActorURI(ctx context.Context, uri string) ([]*gtsmodel.RelayPush, error) {
	var ids []string
	if err := r.db.
		NewSelect().
		TableExpr("? AS ?", bun.Ident("relay_pushes"), bun.Ident("relay_push")).
		Column("relay_push.id").
		Where("? = ?", bun.Ident("relay_push.relay_actor_uri"), uri).
		Scan(ctx, &ids); err != nil {
		return nil, err
	}

	if len(ids) == 0 {
		return nil, nil
	}

	return r.getRelayPushesByIDs(ctx, ids)
}

func (r *relayDB) getRelayPushesByIDs(ctx context.Context, ids []string) ([]*gtsmodel.RelayPush, error) {
	relayPushes, err := r.state.Caches.DB.RelayPush.LoadIDs("ID",
		ids,
		func(uncached []string) ([]*gtsmodel.RelayPush, error) {
			// Preallocate expected length of uncached relayPushes.
			relayPushes := make([]*gtsmodel.RelayPush, 0, len(uncached))

			// Perform database query scanning
			// the remaining (uncached) IDs.
			if err := r.db.NewSelect().
				Model(&relayPushes).
				Where("? IN (?)", bun.Ident("id"), bun.List(uncached)).
				Scan(ctx); err != nil {
				return nil, err
			}

			return relayPushes, nil
		},
	)
	if err != nil {
		return nil, err
	}

	// Reorder the pushes by their
	// IDs to ensure in correct order.
	getID := func(r *gtsmodel.RelayPush) string { return r.ID }
	xslices.OrderBy(relayPushes, ids, getID)

	if gtscontext.Barebones(ctx) {
		// Return without populating.
		return relayPushes, nil
	}

	// Populate the relay pushes. Remove any that
	// we can't populate from the return slice.
	var errs gtserror.MultiError
	relayPushes = slices.DeleteFunc(relayPushes, func(relayPush *gtsmodel.RelayPush) bool {
		if err := r.PopulateRelayPush(ctx, relayPush); err != nil {
			errs.Appendf("error populating relay push %s: %w", relayPush.ID, err)
			return true
		}
		return false
	})

	return relayPushes, nil
}

func (r *relayDB) PopulateRelayPush(ctx context.Context, relayPush *gtsmodel.RelayPush) error {
	if len(relayPush.MatcherIDs) == 0 {
		// Nothing to populate.
		return nil
	}

	var err error
	relayPush.Matchers, err = r.getRelayMatchersByIDs(ctx, relayPush.MatcherIDs)
	if err != nil && !errors.Is(err, db.ErrNoEntries) {
		return gtserror.Newf("error getting relay matchers for push: %w", err)
	}

	return nil
}

func (r *relayDB) PutRelayPush(ctx context.Context, relayPush *gtsmodel.RelayPush) error {
	return r.state.Caches.DB.RelayPush.Store(relayPush, func() error {
		_, err := r.db.
			NewInsert().
			Model(relayPush).
			Exec(ctx)
		return err
	})
}

func (r *relayDB) UpdateRelayPush(ctx context.Context, relayPush *gtsmodel.RelayPush, columns ...string) error {
	return r.state.Caches.DB.RelayPush.Store(relayPush, func() error {
		_, err := r.db.
			NewUpdate().
			Model(relayPush).
			Column(columns...).
			WherePK().
			Exec(ctx)
		return err
	})
}

func (r *relayDB) DeleteRelayPush(ctx context.Context, relayPush *gtsmodel.RelayPush) error {
	if err := r.db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		// Delete all matchers both known
		// by push, and possible stragglers,
		// storing IDs in relayPush.MatcherIDs.
		if _, err := tx.
			NewDelete().
			Model((*gtsmodel.RelayMatcher)(nil)).
			Where("? = ?", bun.Ident("relay_id"), relayPush.ID).
			Returning("?", bun.Ident("id")).
			Exec(ctx, &relayPush.MatcherIDs); err != nil &&
			!errors.Is(err, db.ErrNoEntries) {
			return err
		}

		// Delete relay push itself.
		_, err := tx.
			NewDelete().
			Model((*gtsmodel.RelayPush)(nil)).
			Where("? = ?", bun.Ident("id"), relayPush.ID).
			Exec(ctx)
		return err
	}); err != nil {
		return err
	}

	// Invalidate the relay push itself, and
	// call invalidate hook in-case not cached.
	r.state.Caches.DB.RelayPush.Invalidate("ID", relayPush.ID)
	r.state.Caches.OnInvalidateRelayPush(relayPush)

	return nil
}

func (r *relayDB) GetRelayActorByID(ctx context.Context, id string) (*gtsmodel.RelayActor, error) {
	relayActor, err := r.state.Caches.DB.RelayActor.LoadOne(
		"ID",
		func() (*gtsmodel.RelayActor, error) {
			var relayActor gtsmodel.RelayActor
			err := r.db.
				NewSelect().
				Model(&relayActor).
				Where("? = ?", bun.Ident("id"), id).
				Scan(ctx)
			return &relayActor, err
		},
		id,
	)
	if err != nil {
		// already processed
		return nil, err
	}

	if !gtscontext.Barebones(ctx) {
		if err := r.PopulateRelayActor(ctx, relayActor); err != nil {
			return nil, err
		}
	}

	return relayActor, nil
}

func (r *relayDB) GetRelayActorByURI(ctx context.Context, uri string) (*gtsmodel.RelayActor, error) {
	relayActor, err := r.state.Caches.DB.RelayActor.LoadOne(
		"URI",
		func() (*gtsmodel.RelayActor, error) {
			var relayActor gtsmodel.RelayActor
			err := r.db.
				NewSelect().
				Model(&relayActor).
				Where("? = ?", bun.Ident("uri"), uri).
				Scan(ctx)
			return &relayActor, err
		},
		uri,
	)
	if err != nil {
		// already processed
		return nil, err
	}

	if !gtscontext.Barebones(ctx) {
		if err := r.PopulateRelayActor(ctx, relayActor); err != nil {
			return nil, err
		}
	}

	return relayActor, nil
}

func (r *relayDB) GetRelayActors(ctx context.Context) ([]*gtsmodel.RelayActor, error) {
	// Load all IDs.
	var ids []string
	if err := r.db.
		NewSelect().
		Table("relay_actors").
		Column("id").
		Scan(ctx, &ids); err != nil {
		return nil, err
	}

	if len(ids) == 0 {
		return nil, nil
	}

	return r.getRelayActorsByIDs(ctx, ids)
}

func (r *relayDB) getRelayActorsByIDs(ctx context.Context, ids []string) ([]*gtsmodel.RelayActor, error) {
	relayActors, err := r.state.Caches.DB.RelayActor.LoadIDs("ID",
		ids,
		func(uncached []string) ([]*gtsmodel.RelayActor, error) {
			// Preallocate expected length of uncached relayActors.
			relayActors := make([]*gtsmodel.RelayActor, 0, len(uncached))

			// Perform database query scanning
			// the remaining (uncached) IDs.
			if err := r.db.NewSelect().
				Model(&relayActors).
				Where("? IN (?)", bun.Ident("id"), bun.List(uncached)).
				Scan(ctx); err != nil {
				return nil, err
			}

			return relayActors, nil
		},
	)
	if err != nil {
		return nil, err
	}

	// Reorder the actores by their
	// IDs to ensure in correct order.
	getID := func(r *gtsmodel.RelayActor) string { return r.ID }
	xslices.OrderBy(relayActors, ids, getID)

	if gtscontext.Barebones(ctx) {
		// Return without populating.
		return relayActors, nil
	}

	// Populate the relay actores. Remove any that
	// we can't populate from the return slice.
	var errs gtserror.MultiError
	relayActors = slices.DeleteFunc(relayActors, func(relayActor *gtsmodel.RelayActor) bool {
		if err := r.PopulateRelayActor(ctx, relayActor); err != nil {
			errs.Appendf("error populating relay actor %s: %w", relayActor.ID, err)
			return true
		}
		return false
	})

	return relayActors, nil
}

func (r *relayDB) PopulateRelayActor(ctx context.Context, relayActor *gtsmodel.RelayActor) error {
	// Get relay actor's account (barebones) from the db.
	// Caller can populate it themselves if they need to.
	if relayActor.ActorAccount == nil {
		var err error
		relayActor.ActorAccount, err = r.state.DB.GetAccountByURI(
			gtscontext.SetBarebones(ctx),
			relayActor.URI,
		)
		if err != nil {
			return gtserror.Newf("error getting account for actor: %w", err)
		}
	}

	if len(relayActor.MatcherIDs) == 0 {
		// Nothing more
		// to populate.
		return nil
	}

	var err error
	relayActor.Matchers, err = r.getRelayMatchersByIDs(ctx, relayActor.MatcherIDs)
	if err != nil && !errors.Is(err, db.ErrNoEntries) {
		return gtserror.Newf("error getting relay matchers for actor: %w", err)
	}

	return nil
}

func (r *relayDB) PutRelayActor(ctx context.Context, relayActor *gtsmodel.RelayActor) error {
	return r.state.Caches.DB.RelayActor.Store(relayActor, func() error {
		_, err := r.db.
			NewInsert().
			Model(relayActor).
			Exec(ctx)
		return err
	})
}

func (r *relayDB) UpdateRelayActor(ctx context.Context, relayActor *gtsmodel.RelayActor, columns ...string) error {
	return r.state.Caches.DB.RelayActor.Store(relayActor, func() error {
		_, err := r.db.
			NewUpdate().
			Model(relayActor).
			Column(columns...).
			WherePK().
			Exec(ctx)
		return err
	})
}

func (r *relayDB) DeleteRelayActor(ctx context.Context, relayActor *gtsmodel.RelayActor) error {
	if err := r.db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		// Delete all matchers both known
		// by actor, and possible stragglers,
		// storing IDs in relayActor.MatcherIDs.
		if _, err := tx.
			NewDelete().
			Model((*gtsmodel.RelayMatcher)(nil)).
			Where("? = ?", bun.Ident("relay_id"), relayActor.ID).
			Returning("?", bun.Ident("id")).
			Exec(ctx, &relayActor.MatcherIDs); err != nil &&
			!errors.Is(err, db.ErrNoEntries) {
			return err
		}

		// Delete relay actor itself.
		_, err := tx.
			NewDelete().
			Model((*gtsmodel.RelayActor)(nil)).
			Where("? = ?", bun.Ident("id"), relayActor.ID).
			Exec(ctx)
		return err
	}); err != nil {
		return err
	}

	// Invalidate the relay actor itself, and
	// call invalidate hook in-case not cached.
	r.state.Caches.DB.RelayActor.Invalidate("ID", relayActor.ID)
	r.state.Caches.OnInvalidateRelayActor(relayActor)

	return nil
}

func (r *relayDB) GetRelaySubscriptionByID(ctx context.Context, id string) (*gtsmodel.RelaySubscription, error) {
	relaySubscription, err := r.state.Caches.DB.RelaySubscription.LoadOne(
		"ID",
		func() (*gtsmodel.RelaySubscription, error) {
			var relaySubscription gtsmodel.RelaySubscription
			err := r.db.
				NewSelect().
				Model(&relaySubscription).
				Where("? = ?", bun.Ident("id"), id).
				Scan(ctx)
			return &relaySubscription, err
		},
		id,
	)
	if err != nil {
		// already processed
		return nil, err
	}

	if !gtscontext.Barebones(ctx) {
		if err := r.PopulateRelaySubscription(ctx, relaySubscription); err != nil {
			return nil, err
		}
	}

	return relaySubscription, nil
}

func (r *relayDB) GetRelaySubscriptions(ctx context.Context) ([]*gtsmodel.RelaySubscription, error) {
	// Load all IDs.
	var ids []string
	if err := r.db.
		NewSelect().
		Table("relay_subscriptions").
		Column("id").
		Scan(ctx, &ids); err != nil {
		return nil, err
	}

	if len(ids) == 0 {
		return nil, nil
	}

	return r.getRelaySubscriptionsByIDs(ctx, ids)
}

func (r *relayDB) GetRelaySubscriptionsByActorURI(ctx context.Context, uri string) ([]*gtsmodel.RelaySubscription, error) {
	var ids []string
	if err := r.db.
		NewSelect().
		TableExpr("? AS ?", bun.Ident("relay_subscriptions"), bun.Ident("relay_subscription")).
		Column("relay_subscription.id").
		Where("? = ?", bun.Ident("relay_subscription.relay_actor_uri"), uri).
		Scan(ctx, &ids); err != nil {
		return nil, err
	}

	if len(ids) == 0 {
		return nil, nil
	}

	return r.getRelaySubscriptionsByIDs(ctx, ids)
}

func (r *relayDB) getRelaySubscriptionsByIDs(ctx context.Context, ids []string) ([]*gtsmodel.RelaySubscription, error) {
	relaySubscriptions, err := r.state.Caches.DB.RelaySubscription.LoadIDs("ID",
		ids,
		func(uncached []string) ([]*gtsmodel.RelaySubscription, error) {
			// Preallocate expected length of uncached relaySubscriptions.
			relaySubscriptions := make([]*gtsmodel.RelaySubscription, 0, len(uncached))

			// Perform database query scanning
			// the remaining (uncached) IDs.
			if err := r.db.NewSelect().
				Model(&relaySubscriptions).
				Where("? IN (?)", bun.Ident("id"), bun.List(uncached)).
				Scan(ctx); err != nil {
				return nil, err
			}

			return relaySubscriptions, nil
		},
	)
	if err != nil {
		return nil, err
	}

	// Reorder the subscriptiones by their
	// IDs to ensure in correct order.
	getID := func(r *gtsmodel.RelaySubscription) string { return r.ID }
	xslices.OrderBy(relaySubscriptions, ids, getID)

	if gtscontext.Barebones(ctx) {
		// Return without populating.
		return relaySubscriptions, nil
	}

	// Populate the relay subscriptiones. Remove any that
	// we can't populate from the return slice.
	var errs gtserror.MultiError
	relaySubscriptions = slices.DeleteFunc(relaySubscriptions, func(relaySubscription *gtsmodel.RelaySubscription) bool {
		if err := r.PopulateRelaySubscription(ctx, relaySubscription); err != nil {
			errs.Appendf("error populating relay subscription %s: %w", relaySubscription.ID, err)
			return true
		}
		return false
	})

	return relaySubscriptions, nil
}

func (r *relayDB) PopulateRelaySubscription(ctx context.Context, relaySubscription *gtsmodel.RelaySubscription) error {
	if len(relaySubscription.MatcherIDs) == 0 {
		// Nothing to populate.
		return nil
	}

	var err error
	relaySubscription.Matchers, err = r.getRelayMatchersByIDs(ctx, relaySubscription.MatcherIDs)
	if err != nil && !errors.Is(err, db.ErrNoEntries) {
		return gtserror.Newf("error getting relay matchers for subscription: %w", err)
	}

	return nil
}

func (r *relayDB) PutRelaySubscription(ctx context.Context, relaySubscription *gtsmodel.RelaySubscription) error {
	return r.state.Caches.DB.RelaySubscription.Store(relaySubscription, func() error {
		_, err := r.db.
			NewInsert().
			Model(relaySubscription).
			Exec(ctx)
		return err
	})
}

func (r *relayDB) UpdateRelaySubscription(ctx context.Context, relaySubscription *gtsmodel.RelaySubscription, columns ...string) error {
	return r.state.Caches.DB.RelaySubscription.Store(relaySubscription, func() error {
		_, err := r.db.
			NewUpdate().
			Model(relaySubscription).
			Column(columns...).
			WherePK().
			Exec(ctx)
		return err
	})
}

func (r *relayDB) DeleteRelaySubscription(ctx context.Context, relaySubscription *gtsmodel.RelaySubscription) error {
	if err := r.db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		// Delete all matchers both known
		// by subscription, and possible stragglers,
		// storing IDs in relaySubscription.MatcherIDs.
		if _, err := tx.
			NewDelete().
			Model((*gtsmodel.RelayMatcher)(nil)).
			Where("? = ?", bun.Ident("relay_id"), relaySubscription.ID).
			Returning("?", bun.Ident("id")).
			Exec(ctx, &relaySubscription.MatcherIDs); err != nil &&
			!errors.Is(err, db.ErrNoEntries) {
			return err
		}

		// Delete relay subscription itself.
		_, err := tx.
			NewDelete().
			Model((*gtsmodel.RelaySubscription)(nil)).
			Where("? = ?", bun.Ident("id"), relaySubscription.ID).
			Exec(ctx)
		return err
	}); err != nil {
		return err
	}

	// Invalidate the relay subscription itself, and
	// call invalidate hook in-case not cached.
	r.state.Caches.DB.RelaySubscription.Invalidate("ID", relaySubscription.ID)
	r.state.Caches.OnInvalidateRelaySubscription(relaySubscription)

	return nil
}

func (r *relayDB) getRelayMatchersByIDs(ctx context.Context, ids []string) ([]*gtsmodel.RelayMatcher, error) {
	relayMatchers, err := r.state.Caches.DB.RelayMatcher.LoadIDs("ID",
		ids,
		func(uncached []string) ([]*gtsmodel.RelayMatcher, error) {
			// Preallocate expected length of uncached relay matchers.
			relayMatchers := make([]*gtsmodel.RelayMatcher, 0, len(uncached))

			// Perform database query scanning
			// the remaining (uncached) IDs.
			if err := r.db.NewSelect().
				Model(&relayMatchers).
				Where("? IN (?)", bun.Ident("id"), bun.List(uncached)).
				Scan(ctx); err != nil {
				return nil, err
			}

			// Compile all the keyword regular expressions.
			relayMatchers = slices.DeleteFunc(relayMatchers, func(relayMatcher *gtsmodel.RelayMatcher) bool {
				if err := relayMatcher.Compile(); err != nil {
					log.Errorf(ctx, "error compiling relay matcher keyword regex: %v", err)
					return true
				}
				return false
			})

			return relayMatchers, nil
		},
	)
	if err != nil {
		return nil, err
	}

	// Put matcher keyword structs in the
	// same order as the matcher keyword IDs.
	getID := func(a *gtsmodel.RelayMatcher) string { return a.ID }
	xslices.OrderBy(relayMatchers, ids, getID)

	return relayMatchers, nil
}

func (r *relayDB) GetRelayMatcher(ctx context.Context, id string) (*gtsmodel.RelayMatcher, error) {
	return r.state.Caches.DB.RelayMatcher.LoadOne(
		"ID",
		func() (*gtsmodel.RelayMatcher, error) {
			var relayMatcher gtsmodel.RelayMatcher
			err := r.db.
				NewSelect().
				Model(&relayMatcher).
				Where("? = ?", bun.Ident("id"), id).
				Scan(ctx)
			return &relayMatcher, err
		},
		id,
	)
}

func (r *relayDB) PutRelayMatcher(ctx context.Context, relayMatcher *gtsmodel.RelayMatcher) error {
	if relayMatcher.Regexp == nil {
		// Ensure regexp is compiled
		// before attempted caching.
		err := relayMatcher.Compile()
		if err != nil {
			return gtserror.Newf("error compiling relay matcher regex: %w", err)
		}
	}
	return r.state.Caches.DB.RelayMatcher.Store(relayMatcher, func() error {
		_, err := r.db.
			NewInsert().
			Model(relayMatcher).
			Exec(ctx)
		return err
	})
}

func (r *relayDB) UpdateRelayMatcher(ctx context.Context, relayMatcher *gtsmodel.RelayMatcher, columns ...string) error {
	if relayMatcher.Regexp == nil {
		// Ensure regexp is compiled
		// before attempted caching.
		err := relayMatcher.Compile()
		if err != nil {
			return gtserror.Newf("error compiling relay matcher regex: %w", err)
		}
	}
	return r.state.Caches.DB.RelayMatcher.Store(relayMatcher, func() error {
		_, err := r.db.
			NewUpdate().
			Model(relayMatcher).
			Column(columns...).
			WherePK().
			Exec(ctx)
		return err
	})
}

func (r *relayDB) DeleteRelayMatcher(ctx context.Context, id string) error {
	_, err := r.db.
		NewDelete().
		Table("relay_matchers").
		Where("? = ?", bun.Ident("id"), id).
		Exec(ctx)
	if err != nil {
		return err
	}

	r.state.Caches.DB.RelayMatcher.Invalidate("ID", id)
	return nil
}

func (r *relayDB) PutRelayedURI(ctx context.Context, relayedURI *gtsmodel.RelayedURI) error {
	_, err := r.db.NewInsert().Model(relayedURI).Exec(ctx)
	return err
}

func (r *relayDB) GetRelayedURI(ctx context.Context, relayURI string, statusURI string) (*gtsmodel.RelayedURI, error) {
	relayedURI := new(gtsmodel.RelayedURI)
	if err := r.db.
		NewSelect().
		Model(relayedURI).
		Where("? = ?", bun.Ident("relay_uri"), relayURI).
		Where("? = ?", bun.Ident("status_uri"), statusURI).
		Scan(ctx); err != nil {
		return nil, err
	}
	return relayedURI, nil
}
