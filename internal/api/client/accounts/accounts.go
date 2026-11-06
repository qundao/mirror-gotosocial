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

package accounts

import (
	"code.superseriousbusiness.org/gopkg/httputil"
	apiutil "code.superseriousbusiness.org/gotosocial/internal/api/util"
	"code.superseriousbusiness.org/gotosocial/internal/processing"
	"code.superseriousbusiness.org/gotosocial/internal/templates"
)

const (
	ExcludeReblogsKey = "exclude_reblogs"
	ExcludeRepliesKey = "exclude_replies"
	LimitKey          = "limit"
	MaxIDKey          = "max_id"
	MinIDKey          = "min_id"
	OnlyMediaKey      = "only_media"
	OnlyPublicKey     = "only_public"
	PinnedKey         = "pinned"

	BasePath       = "/v1/accounts"
	BasePathWithID = BasePath + "/:" + apiutil.IDKey

	BlockPath               = BasePathWithID + "/block"
	DeletePath              = BasePath + "/delete"
	FeaturedTagsPath        = BasePathWithID + "/featured_tags"
	FollowersPath           = BasePathWithID + "/followers"
	FollowingPath           = BasePathWithID + "/following"
	FollowPath              = BasePathWithID + "/follow"
	ListsPath               = BasePathWithID + "/lists"
	LookupPath              = BasePath + "/lookup"
	MutePath                = BasePathWithID + "/mute"
	NotePath                = BasePathWithID + "/note"
	RelationshipsPath       = BasePath + "/relationships"
	RemoveFromFollowersPath = BasePathWithID + "/remove_from_followers"
	SearchPath              = BasePath + "/search"
	StatusesPath            = BasePathWithID + "/statuses"
	UnblockPath             = BasePathWithID + "/unblock"
	UnfollowPath            = BasePathWithID + "/unfollow"
	UnmutePath              = BasePathWithID + "/unmute"
	UpdatePath              = BasePath + "/update_credentials"
	VerifyPath              = BasePath + "/verify_credentials"
	MovePath                = BasePath + "/move"
	AliasPath               = BasePath + "/alias"
	ThemesPath              = BasePath + "/themes"

	// ProfileBasePath for the profile API, an extension
	// of the account update API with a different path.
	ProfileBasePath = "/v1/profile"
	AvatarPath      = ProfileBasePath + "/avatar"
	HeaderPath      = ProfileBasePath + "/header"
)

type Module struct {
	templates *templates.Templates
	processor *processing.Processor
}

func New(processor *processing.Processor, templates *templates.Templates) *Module {
	return &Module{
		templates: templates,
		processor: processor,
	}
}

func (m *Module) Route(g *httputil.RouteGroup) {
	// create account
	g.POST(BasePath, m.AccountCreatePOSTHandler)

	// get account(s)
	g.GET(BasePath, m.AccountsGETHandler)
	g.GET(BasePathWithID, m.AccountGETHandler)

	// delete account
	g.POST(DeletePath, m.AccountDeletePOSTHandler)

	// verify account
	g.GET(VerifyPath, m.AccountVerifyGETHandler)

	// modify account
	g.PATCH(UpdatePath, m.AccountUpdateCredentialsPATCHHandler)

	// modify account profile media
	g.DELETE(AvatarPath, m.AccountAvatarDELETEHandler)
	g.DELETE(HeaderPath, m.AccountHeaderDELETEHandler)

	// get account's statuses
	g.GET(StatusesPath, m.AccountStatusesGETHandler)

	// get account's featured tags
	g.GET(FeaturedTagsPath, m.AccountFeaturedTagsGETHandler)

	// get following or followers
	g.GET(FollowersPath, m.AccountFollowersGETHandler)
	g.GET(FollowingPath, m.AccountFollowingGETHandler)

	// get relationship with account
	g.GET(RelationshipsPath, m.AccountRelationshipsGETHandler)

	// follow or unfollow account
	g.POST(FollowPath, m.AccountFollowPOSTHandler)
	g.POST(UnfollowPath, m.AccountUnfollowPOSTHandler)

	// block or unblock account
	g.POST(BlockPath, m.AccountBlockPOSTHandler)
	g.POST(UnblockPath, m.AccountUnblockPOSTHandler)

	// remove account from followers
	g.POST(RemoveFromFollowersPath, m.AccountRemoveFromFollowersPOSTHandler)

	// account lists
	g.GET(ListsPath, m.AccountListsGETHandler)

	// account note
	g.POST(NotePath, m.AccountNotePOSTHandler)

	// mute or unmute account
	g.POST(MutePath, m.AccountMutePOSTHandler)
	g.POST(UnmutePath, m.AccountUnmutePOSTHandler)

	// search for accounts
	g.GET(SearchPath, m.AccountSearchGETHandler)
	g.GET(LookupPath, m.AccountLookupGETHandler)

	// migration handlers
	g.POST(AliasPath, m.AccountAliasPOSTHandler)
	g.POST(MovePath, m.AccountMovePOSTHandler)

	// account themes
	g.GET(ThemesPath, m.AccountThemesGETHandler)
}
