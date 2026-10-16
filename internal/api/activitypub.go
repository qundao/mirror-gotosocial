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

package api

import (
	"code.superseriousbusiness.org/gopkg/httputil"
	"code.superseriousbusiness.org/gotosocial/internal/api/activitypub/actor"
	"code.superseriousbusiness.org/gotosocial/internal/api/activitypub/emoji"
	"code.superseriousbusiness.org/gotosocial/internal/api/activitypub/publickey"
	"code.superseriousbusiness.org/gotosocial/internal/config"
	"code.superseriousbusiness.org/gotosocial/internal/db"
	"code.superseriousbusiness.org/gotosocial/internal/middleware"
	"code.superseriousbusiness.org/gotosocial/internal/processing"
	"code.superseriousbusiness.org/gotosocial/internal/router"
	"code.superseriousbusiness.org/gotosocial/internal/templates"
)

type ActivityPub struct {
	emoji     *emoji.Module
	actor     *actor.Module
	publicKey *publickey.Module
	sigcheck  httputil.Middleware
}

func (a *ActivityPub) Route(r *router.Router, m ...httputil.Middleware) {
	// create groupings for the
	// 'emoji' and 'users|relays' prefixes
	emojiGroup := r.Group("emoji")
	usersGroup := r.Group("users")
	relaysGroup := r.Group("relays")

	// Use provided middlewares.
	emojiGroup.Use(m...)
	usersGroup.Use(m...)
	relaysGroup.Use(m...)

	// Attach cache control middleware.
	ccMiddleware := middleware.CacheControl(middleware.CacheControlConfig{
		Directives: []string{"no-store"},
	})
	emojiGroup.Use(ccMiddleware)
	usersGroup.Use(ccMiddleware)
	relaysGroup.Use(ccMiddleware)

	// Route the instance actor endpoint first
	// so it doesn't require signature auth.
	usersGroup.GET(config.GetHost(), a.actor.InstanceActorGETHandler)

	// *Now* add sig checking.
	emojiGroup.Use(a.sigcheck)
	usersGroup.Use(a.sigcheck)
	relaysGroup.Use(a.sigcheck)

	a.emoji.Route(emojiGroup)
	a.actor.Route(usersGroup)
	a.actor.Route(relaysGroup)
}

// Public key endpoint requires different middleware + cache policies from other AP endpoints.
func (a *ActivityPub) RoutePublicKey(r *router.Router, m ...httputil.Middleware) {

	// Create grouping for the
	// '{users|relays}/[username]/main-key' prefixes
	usersPublicKeyGroup := r.Group(publickey.UsersPublicKeyPath)
	relaysPublicKeyGroup := r.Group(publickey.RelaysPublicKeyPath)

	// Attach middleware allowing public cacheing of main-key.
	ccMiddleware := middleware.CacheControl(middleware.CacheControlConfig{
		Directives: []string{"public", "max-age=604800"},
		Vary:       []string{"Accept", "Accept-Encoding"},
	})
	usersPublicKeyGroup.Use(m...)
	relaysPublicKeyGroup.Use(m...)
	usersPublicKeyGroup.Use(a.sigcheck, ccMiddleware)
	relaysPublicKeyGroup.Use(a.sigcheck, ccMiddleware)

	a.publicKey.Route(usersPublicKeyGroup)
	a.publicKey.Route(relaysPublicKeyGroup)
}

func NewActivityPub(db db.DB, processor *processing.Processor, templates *templates.Templates) *ActivityPub {
	return &ActivityPub{
		emoji:     emoji.New(processor, templates),
		actor:     actor.New(processor, templates),
		publicKey: publickey.New(processor, templates),
		sigcheck:  middleware.ExtractSignature(db.IsURIBlocked),
	}
}
