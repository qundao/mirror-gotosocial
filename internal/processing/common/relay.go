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

package common

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"code.superseriousbusiness.org/gotosocial/internal/ap"
	apimodel "code.superseriousbusiness.org/gotosocial/internal/api/model"
	"code.superseriousbusiness.org/gotosocial/internal/db"
	"code.superseriousbusiness.org/gotosocial/internal/gtscontext"
	"code.superseriousbusiness.org/gotosocial/internal/gtserror"
	gtsmodel "code.superseriousbusiness.org/gotosocial/internal/gtsmodel"
	"code.superseriousbusiness.org/gotosocial/internal/id"
	"code.superseriousbusiness.org/gotosocial/internal/messages"
	"code.superseriousbusiness.org/gotosocial/internal/uris"
)

// FollowRelayActor checks if a follow or follow request exists from our instance account towards
// the given relayActor account. If not, it creates a follow request and instigates side effects.
func (p *Processor) FollowRelayActor(ctx context.Context, relayActor *gtsmodel.Account) error {
	// Get our instance account from the db.
	iAcct, err := p.state.DB.GetInstanceAccount(ctx, "")
	if err != nil {
		return gtserror.NewfAt(3, "db error getting instance account: %w", err)
	}

	// Check if already following.
	follows, err := p.state.DB.IsFollowing(ctx, iAcct.ID, relayActor.ID)
	if err != nil {
		return gtserror.NewfAt(3, "db error checking follow: %w", err)
	}
	if follows {
		// Nothing to do.
		return nil
	}

	// Check if already follow requested.
	followRequests, err := p.state.DB.IsFollowRequested(ctx, iAcct.ID, relayActor.ID)
	if err != nil {
		return gtserror.NewfAt(3, "db error checking follow request: %w", err)
	}
	if followRequests {
		// Nothing to do.
		return nil
	}

	// No follow or follow
	// request, create one.
	followID := id.NewULID()
	fr := &gtsmodel.FollowRequest{
		ID: followID,
		URI: uris.GenerateURIForFollow(
			uris.UsersPath,
			iAcct.Username,
			followID,
		),
		AccountID:       iAcct.ID,
		Account:         iAcct,
		TargetAccountID: relayActor.ID,
		TargetAccount:   relayActor,
		Flags: func() (flags gtsmodel.FollowFlags) {
			flags.SetShowReblogs(false)
			flags.SetNotify(false)
			return
		}(),
	}

	// Insert the follow request.
	if err := p.state.DB.PutFollowRequest(ctx, fr); err != nil {
		return gtserror.NewfAt(3, "db error putting follow request: %w", err)
	}

	// Handle side effects async.
	p.state.Workers.Client.Queue.Push(&messages.FromClientAPI{
		APObjectType:   ap.ActivityFollow,
		APActivityType: ap.ActivityCreate,
		GTSModel:       fr,
		Origin:         iAcct,
		Target:         relayActor,
	})

	return nil
}

// UnfollowRelayActorIfNoConnections checks if the given relayActor is targeted
// by any relay pushes or relay subscriptions remaining in the database. If not,
// an Undo Follow will be sent from the instance account targeting the relay actor,
// essentially stopping the connection.
func (p *Processor) UnfollowRelayActorIfNoConnections(ctx context.Context, relayActor *gtsmodel.Account) error {
	relayActorURI := relayActor.URI

	// Check if there's any relay pushes with this actor URI.
	pushes, err := p.state.DB.GetRelayPushesByActorURI(
		// Just use barebones, we're
		// only checking if exists.
		gtscontext.SetBarebones(ctx),
		relayActorURI,
	)
	if err != nil && !errors.Is(err, db.ErrNoEntries) {
		return gtserror.NewfAt(3, "db error getting relay pushes: %w", err)
	}

	if len(pushes) != 0 {
		// Pushes still
		// exist, noop.
		return nil
	}

	// Check if there's any relay subs with this actor URI.
	subs, err := p.state.DB.GetRelaySubscriptionsByActorURI(
		// Just use barebones, we're
		// only checking if exists.
		gtscontext.SetBarebones(ctx),
		relayActorURI,
	)
	if err != nil && !errors.Is(err, db.ErrNoEntries) {
		return gtserror.NewfAt(3, "db error getting relay subscriptions: %w", err)
	}

	if len(subs) != 0 {
		// Subs still
		// exist, noop.
		return nil
	}

	// No pushes or subscriptions remain targeting the
	// given relay actor, so we should Undo any Follow
	// from the instance account targeting them.

	// Get our instance account from the db.
	iAcct, err := p.state.DB.GetInstanceAccount(ctx, "")
	if err != nil {
		return gtserror.NewfAt(3, "db error getting instance account: %w", err)
	}

	// Unfollow the relay actor, sending
	// side effects to the client worker
	// so the Undo Follow gets federated.
	if err := p.Unfollow(ctx,
		iAcct,
		relayActor,
		true, // sideEffects
	); err != nil {
		return gtserror.NewfAt(3, "error unfollowing: %w", err)
	}

	return nil
}

func (p *Processor) GetRelayMatcher(
	ctx context.Context,
	id string,
) (*gtsmodel.RelayMatcher, gtserror.WithCode) {
	matcher, err := p.state.DB.GetRelayMatcher(ctx, id)
	if err != nil && !errors.Is(err, db.ErrNoEntries) {
		err := gtserror.Newf("db error getting relay matcher: %w", err)
		return nil, gtserror.NewErrorInternalError(err)
	}

	if matcher == nil {
		err := gtserror.New("relay matcher not found")
		return nil, gtserror.NewErrorNotFound(err)
	}

	return matcher, nil
}

func (p *Processor) DereferenceRelayActorURI(
	ctx context.Context,
	relayActorURI *url.URL,
) (*gtsmodel.Account, gtserror.WithCode) {
	// Ensure relay actor host is not blocked.
	blocked, err := p.state.DB.IsDomainBlocked(ctx, relayActorURI.Host)
	if err != nil {
		err = gtserror.Newf("db error checking for domain block: %w", err)
		return nil, gtserror.NewErrorInternalError(err, err.Error())
	}

	if blocked {
		err = gtserror.Newf("target domain %s is blocked", relayActorURI.Host)
		return nil, gtserror.NewErrorUnprocessableEntity(err, err.Error())
	}

	// Dereference relay actor using our instance account.
	instanceAcct, err := p.state.DB.GetInstanceAccount(ctx, "")
	if err != nil {
		err = gtserror.Newf("db error getting instance account: %w", err)
		return nil, gtserror.NewErrorInternalError(err, err.Error())
	}

	relayActor, _, err := p.federator.Dereferencer.GetAccountByURI(ctx,
		instanceAcct.Username,
		relayActorURI,
		true, // tryURL
	)
	if err != nil {
		err := fmt.Errorf(
			"error dereferencing relay_actor_uri (%s) account: %w",
			relayActorURI.String(), err,
		)
		return nil, gtserror.NewErrorUnprocessableEntity(err, err.Error())
	}

	return relayActor, nil
}

func (p *Processor) APIRelayConnection(
	ctx context.Context,
	relayConnection gtsmodel.RelayConnection,
) (*apimodel.RelayConnection, gtserror.WithCode) {
	item, err := p.converter.RelayConnectionToAPIRelayConnection(ctx, relayConnection)
	if err != nil {
		err := gtserror.Newf("error converting relay connection to api model: %w", err)
		return nil, gtserror.NewErrorInternalError(err)
	}

	return item, nil
}

func (p *Processor) UpdateRelayMatcher(
	ctx context.Context,
	matcher *gtsmodel.RelayMatcher,
	keyword *string,
	wholeWord *bool,
	exclude *bool,
) gtserror.WithCode {
	// Db columns to be updated.
	columns := make([]string, 0, 3)

	// Check/update keyword.
	if keyword != nil {
		matcher.Keyword = *keyword
		if matcher.Keyword == "" {
			const errText = "keyword cannot be empty"
			return gtserror.NewErrorBadRequest(errors.New(errText), errText)
		}
		columns = append(columns, "keyword")
	}

	// Check/update whole word.
	if wholeWord != nil {
		matcher.Flags.SetWholeWord(*wholeWord)
		columns = append(columns, "whole_word")
	}

	// Check/update exclude.
	if exclude != nil {
		matcher.Flags.SetExclude(*exclude)
		columns = append(columns, "exclude")
	}

	// Ensure matcher can be compiled.
	if err := matcher.Compile(); err != nil {
		err := gtserror.Newf("matcher could not be compiled: %w", err)
		return gtserror.NewErrorUnprocessableEntity(err, err.Error())
	}

	// Store it.
	switch err := p.state.DB.UpdateRelayMatcher(ctx, matcher, columns...); {
	case err == nil:
		// No problemo!
		return nil

	case errors.Is(err, db.ErrAlreadyExists):
		const text = "duplicate keyword"
		return gtserror.NewWithCode(http.StatusConflict, text)

	default:
		err := gtserror.Newf("db error updating matcher: %w", err)
		return gtserror.NewErrorInternalError(err)
	}
}
