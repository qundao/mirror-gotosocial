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
	"code.superseriousbusiness.org/gopkg/httputil"
	apiutil "code.superseriousbusiness.org/gotosocial/internal/api/util"
	"code.superseriousbusiness.org/gotosocial/internal/processing"
	"code.superseriousbusiness.org/gotosocial/internal/state"
	"code.superseriousbusiness.org/gotosocial/internal/templates"
)

const (
	BasePath = "/v1/admin"
	WithID   = "/:" + apiutil.IDKey

	EmojiPath           = BasePath + "/custom_emojis"
	EmojiPathWithID     = EmojiPath + WithID
	EmojiCategoriesPath = EmojiPath + "/categories"

	DomainBlocksPath       = BasePath + "/domain_blocks"
	DomainBlocksPathWithID = DomainBlocksPath + WithID
	DomainAllowsPath       = BasePath + "/domain_allows"
	DomainAllowsPathWithID = DomainAllowsPath + WithID
	DomainLimitsPath       = BasePath + "/domain_limits"
	DomainLimitsPathWithID = DomainLimitsPath + WithID

	DomainPermissionDraftsPath               = BasePath + "/domain_permission_drafts"
	DomainPermissionDraftsPathWithID         = DomainPermissionDraftsPath + WithID
	DomainPermissionDraftAcceptPath          = DomainPermissionDraftsPathWithID + "/accept"
	DomainPermissionDraftRemovePath          = DomainPermissionDraftsPathWithID + "/remove"
	DomainPermissionExcludesPath             = BasePath + "/domain_permission_excludes"
	DomainPermissionExcludesPathWithID       = DomainPermissionExcludesPath + WithID
	DomainPermissionSubscriptionsPath        = BasePath + "/domain_permission_subscriptions"
	DomainPermissionSubscriptionsPathWithID  = DomainPermissionSubscriptionsPath + WithID
	DomainPermissionSubscriptionsPreviewPath = DomainPermissionSubscriptionsPath + "/preview"
	DomainPermissionSubscriptionRemovePath   = DomainPermissionSubscriptionsPathWithID + "/remove"
	DomainPermissionSubscriptionTestPath     = DomainPermissionSubscriptionsPathWithID + "/test"

	DomainKeysExpirePath = BasePath + "/domain_keys_expire"

	HeaderAllowsPath       = BasePath + "/header_allows"
	HeaderAllowsPathWithID = HeaderAllowsPath + WithID
	HeaderBlocksPath       = BasePath + "/header_blocks"
	HeaderBlocksPathWithID = HeaderBlocksPath + WithID

	AccountsV1Path      = BasePath + "/accounts"
	AccountsV2Path      = "/v2/admin/accounts"
	AccountsPathWithID  = AccountsV1Path + WithID
	AccountsActionPath  = AccountsPathWithID + "/action"
	AccountsApprovePath = AccountsPathWithID + "/approve"
	AccountsRejectPath  = AccountsPathWithID + "/reject"

	MediaCleanupPath = BasePath + "/media_cleanup"
	MediaPurgePath   = BasePath + "/media_purge"
	MediaRefetchPath = BasePath + "/media_refetch"

	ReportsPath        = BasePath + "/reports"
	ReportsPathWithID  = ReportsPath + WithID
	ReportsResolvePath = ReportsPathWithID + "/resolve"

	EmailPath     = BasePath + "/email"
	EmailTestPath = EmailPath + "/test"

	InstanceRulesPath               = BasePath + "/instance/rules"
	InstanceRulesPathWithID         = InstanceRulesPath + WithID
	InstancesPath                   = BasePath + "/instances"
	InstancesPathWithID             = InstancesPath + WithID
	InstanceClearDeliveryErrorsPath = InstancesPathWithID + "/clear_delivery_errors"

	RelaySubscriptionsPath                     = BasePath + "/relay_subscriptions"
	RelaySubscriptionsPathWithID               = RelaySubscriptionsPath + WithID
	RelaySubscriptionMatchersPath              = RelaySubscriptionsPathWithID + "/matchers"
	RelaySubscriptionMatchersPathWithMatcherID = RelaySubscriptionMatchersPath + "/:" + apiutil.RelayMatcherIDKey

	RelayActorsPath                      = BasePath + "/relay_actors"
	RelayActorsPathWithID                = RelayActorsPath + WithID
	RelayActorHeaderPath                 = RelayActorsPathWithID + "/profile/header"
	RelayActorAvatarPath                 = RelayActorsPathWithID + "/profile/avatar"
	RelayActorFollowRequestsPath         = RelayActorsPathWithID + "/follow_requests"
	RelayActorFollowRequestAuthorizePath = RelayActorFollowRequestsPath + "/:" + apiutil.TargetAccountIDKey + "/authorize"
	RelayActorFollowRequestRejectPath    = RelayActorFollowRequestsPath + "/:" + apiutil.TargetAccountIDKey + "/reject"
	RelayActorFollowersPath              = RelayActorsPathWithID + "/followers"
	RelayActorAccountsPath               = RelayActorsPathWithID + "/accounts"
	RelayActorRemoveFromFollowersPath    = RelayActorAccountsPath + "/:" + apiutil.TargetAccountIDKey + "/remove_from_followers"
	RelayActorBlocksPath                 = RelayActorsPathWithID + "/blocks"
	RelayActorBlockPath                  = RelayActorAccountsPath + "/:" + apiutil.TargetAccountIDKey + "/block"
	RelayActorUnblockPath                = RelayActorAccountsPath + "/:" + apiutil.TargetAccountIDKey + "/unblock"
	RelayActorMatchersPath               = RelayActorsPathWithID + "/matchers"
	RelayActorMatchersPathWithMatcherID  = RelayActorMatchersPath + "/:" + apiutil.RelayMatcherIDKey

	RuntimeConfigPath = BasePath + "/runtime_config"

	FilterQueryKey        = "filter"
	MaxShortcodeDomainKey = "max_shortcode_domain"
	MinShortcodeDomainKey = "min_shortcode_domain"
	DomainQueryKey        = "domain"
)

type Module struct {
	templates *templates.Templates
	processor *processing.Processor
	state     *state.State
}

func New(state *state.State, processor *processing.Processor, templates *templates.Templates) *Module {
	return &Module{
		templates: templates,
		processor: processor,
		state:     state,
	}
}

func (m *Module) Route(g *httputil.RouteGroup) {
	// emoji stuff
	g.POST(EmojiPath, m.EmojiCreatePOSTHandler)
	g.GET(EmojiPath, m.EmojisGETHandler)
	g.DELETE(EmojiPathWithID, m.EmojiDELETEHandler)
	g.GET(EmojiPathWithID, m.EmojiGETHandler)
	g.PATCH(EmojiPathWithID, m.EmojiPATCHHandler)
	g.GET(EmojiCategoriesPath, m.EmojiCategoriesGETHandler)

	// domain block stuff
	g.POST(DomainBlocksPath, m.DomainBlocksPOSTHandler)
	g.GET(DomainBlocksPath, m.DomainBlocksGETHandler)
	g.GET(DomainBlocksPathWithID, m.DomainBlockGETHandler)
	g.PUT(DomainBlocksPathWithID, m.DomainBlockUpdatePUTHandler)
	g.DELETE(DomainBlocksPathWithID, m.DomainBlockDELETEHandler)

	// domain allow stuff
	g.POST(DomainAllowsPath, m.DomainAllowsPOSTHandler)
	g.GET(DomainAllowsPath, m.DomainAllowsGETHandler)
	g.GET(DomainAllowsPathWithID, m.DomainAllowGETHandler)
	g.PUT(DomainAllowsPathWithID, m.DomainAllowUpdatePUTHandler)
	g.DELETE(DomainAllowsPathWithID, m.DomainAllowDELETEHandler)

	// domain limits stuff
	g.GET(DomainLimitsPath, m.DomainLimitsGETHandler)
	g.POST(DomainLimitsPath, m.DomainLimitsPOSTHandler)
	g.PUT(DomainLimitsPathWithID, m.DomainLimitPUTHandler)
	g.DELETE(DomainLimitsPathWithID, m.DomainLimitDELETEHandler)

	// domain permission draft stuff
	g.POST(DomainPermissionDraftsPath, m.DomainPermissionDraftsPOSTHandler)
	g.GET(DomainPermissionDraftsPath, m.DomainPermissionDraftsGETHandler)
	g.GET(DomainPermissionDraftsPathWithID, m.DomainPermissionDraftGETHandler)
	g.POST(DomainPermissionDraftAcceptPath, m.DomainPermissionDraftAcceptPOSTHandler)
	g.POST(DomainPermissionDraftRemovePath, m.DomainPermissionDraftRemovePOSTHandler)

	// domain permission excludes stuff
	g.POST(DomainPermissionExcludesPath, m.DomainPermissionExcludesPOSTHandler)
	g.GET(DomainPermissionExcludesPath, m.DomainPermissionExcludesGETHandler)
	g.GET(DomainPermissionExcludesPathWithID, m.DomainPermissionExcludeGETHandler)
	g.DELETE(DomainPermissionExcludesPathWithID, m.DomainPermissionExcludeDELETEHandler)

	// domain permission subscriptions stuff
	g.POST(DomainPermissionSubscriptionsPath, m.DomainPermissionSubscriptionPOSTHandler)
	g.GET(DomainPermissionSubscriptionsPath, m.DomainPermissionSubscriptionsGETHandler)
	g.GET(DomainPermissionSubscriptionsPreviewPath, m.DomainPermissionSubscriptionsPreviewGETHandler)
	g.GET(DomainPermissionSubscriptionsPathWithID, m.DomainPermissionSubscriptionGETHandler)
	g.PATCH(DomainPermissionSubscriptionsPathWithID, m.DomainPermissionSubscriptionPATCHHandler)
	g.POST(DomainPermissionSubscriptionRemovePath, m.DomainPermissionSubscriptionRemovePOSTHandler)
	g.POST(DomainPermissionSubscriptionTestPath, m.DomainPermissionSubscriptionTestPOSTHandler)

	// header filtering administration routes
	g.GET(HeaderAllowsPathWithID, m.HeaderFilterAllowGET)
	g.GET(HeaderBlocksPathWithID, m.HeaderFilterBlockGET)
	g.GET(HeaderAllowsPath, m.HeaderFilterAllowsGET)
	g.GET(HeaderBlocksPath, m.HeaderFilterBlocksGET)
	g.POST(HeaderAllowsPath, m.HeaderFilterAllowPOST)
	g.POST(HeaderBlocksPath, m.HeaderFilterBlockPOST)
	g.DELETE(HeaderAllowsPathWithID, m.HeaderFilterAllowDELETE)
	g.DELETE(HeaderBlocksPathWithID, m.HeaderFilterBlockDELETE)

	// domain maintenance stuff
	g.POST(DomainKeysExpirePath, m.DomainKeysExpirePOSTHandler)

	// accounts stuff
	g.GET(AccountsV1Path, m.AccountsGETV1Handler)
	g.GET(AccountsV2Path, m.AccountsGETV2Handler)
	g.GET(AccountsPathWithID, m.AccountGETHandler)
	g.POST(AccountsActionPath, m.AccountActionPOSTHandler)
	g.POST(AccountsApprovePath, m.AccountApprovePOSTHandler)
	g.POST(AccountsRejectPath, m.AccountRejectPOSTHandler)

	// media stuff
	g.POST(MediaCleanupPath, m.MediaCleanupPOSTHandler)
	g.POST(MediaPurgePath, m.MediaPurgePOSTHandler)
	g.POST(MediaRefetchPath, m.MediaRefetchPOSTHandler)

	// reports stuff
	g.GET(ReportsPath, m.ReportsGETHandler)
	g.GET(ReportsPathWithID, m.ReportGETHandler)
	g.POST(ReportsResolvePath, m.ReportResolvePOSTHandler)

	// email stuff
	g.POST(EmailTestPath, m.EmailTestPOSTHandler)

	// instance rules stuff
	g.GET(InstanceRulesPath, m.RulesGETHandler)
	g.GET(InstanceRulesPathWithID, m.RuleGETHandler)
	g.POST(InstanceRulesPath, m.RulePOSTHandler)
	g.PATCH(InstanceRulesPathWithID, m.RulePATCHHandler)
	g.DELETE(InstanceRulesPathWithID, m.RuleDELETEHandler)

	// instances stuff
	g.GET(InstancesPath, m.InstancesGETHandler)
	g.GET(InstancesPathWithID, m.InstanceGETHandler)
	g.POST(InstanceClearDeliveryErrorsPath, m.InstanceClearDeliveryErrorsPOSTHandler)

	// relay subscriptions stuff
	g.GET(RelaySubscriptionsPath, m.RelaySubscriptionsGETHandler)
	g.GET(RelaySubscriptionsPathWithID, m.RelaySubscriptionGETHandler)
	g.POST(RelaySubscriptionsPath, m.RelaySubscriptionPOSTHandler)
	g.PUT(RelaySubscriptionsPathWithID, m.RelaySubscriptionPUTHandler)
	g.DELETE(RelaySubscriptionsPathWithID, m.RelaySubscriptionDELETEHandler)
	g.POST(RelaySubscriptionMatchersPath, m.RelaySubscriptionMatcherPOSTHandler)
	g.DELETE(RelaySubscriptionMatchersPathWithMatcherID, m.RelaySubscriptionMatcherDELETEHandler)
	g.PUT(RelaySubscriptionMatchersPathWithMatcherID, m.RelaySubscriptionMatcherPUTHandler)

	// relay actors stuff
	g.GET(RelayActorsPath, m.RelayActorsGETHandler)
	g.GET(RelayActorsPathWithID, m.RelayActorGETHandler)
	g.POST(RelayActorsPath, m.RelayActorPOSTHandler)
	g.PUT(RelayActorsPathWithID, m.RelayActorPUTHandler)
	g.DELETE(RelayActorHeaderPath, m.RelayActorHeaderDELETEHandler)
	g.DELETE(RelayActorAvatarPath, m.RelayActorAvatarDELETEHandler)
	g.DELETE(RelayActorsPathWithID, m.RelayActorDELETEHandler)
	g.POST(RelayActorMatchersPath, m.RelayActorMatcherPOSTHandler)
	g.DELETE(RelayActorMatchersPathWithMatcherID, m.RelayActorMatcherDELETEHandler)
	g.PUT(RelayActorMatchersPathWithMatcherID, m.RelayActorMatcherPUTHandler)

	// relay actor account moderation stuff
	// follow (requests) management
	g.GET(RelayActorFollowRequestsPath, m.RelayActorFollowRequestsGETHandler)
	g.POST(RelayActorFollowRequestAuthorizePath, m.RelayActorFollowRequestAuthorizePOSTHandler)
	g.POST(RelayActorFollowRequestRejectPath, m.RelayActorFollowRequestRejectPOSTHandler)
	g.GET(RelayActorFollowersPath, m.RelayActorFollowersGETHandler)
	g.POST(RelayActorRemoveFromFollowersPath, m.RelayActorRemoveFromFollowersPOSTHandler)
	// block management
	g.GET(RelayActorBlocksPath, m.RelayActorBlocksGETHandler)
	g.POST(RelayActorBlockPath, m.RelayActorBlockPOSTHandler)
	g.POST(RelayActorUnblockPath, m.RelayActorUnblockPOSTHandler)

	// runtime config stuff
	g.GET(RuntimeConfigPath, m.RuntimeConfigGETHandler)
	g.PATCH(RuntimeConfigPath, m.RuntimeConfigPATCHHandler)
}
