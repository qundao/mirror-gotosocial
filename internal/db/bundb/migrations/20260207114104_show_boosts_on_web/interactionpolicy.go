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

package gtsmodel

type PolicyValue string

const (
	PolicyValuePublic    PolicyValue = "public"
	PolicyValueFollowers PolicyValue = "followers"
	PolicyValueFollowing PolicyValue = "following"
	PolicyValueMutuals   PolicyValue = "mutuals"
	PolicyValueMentioned PolicyValue = "mentioned"
	PolicyValueAuthor    PolicyValue = "author"
)

type PolicyValues []PolicyValue

type PolicyPermission int

const (
	PolicyPermissionForbidden PolicyPermission = iota
	PolicyPermissionManualApproval
	PolicyPermissionAutomaticApproval
)

type PolicyCheckResult struct {
	Permission          PolicyPermission
	PermissionMatchedOn *PolicyValue
}

type InteractionPolicy struct {
	CanLike     *PolicyRules
	CanReply    *PolicyRules
	CanAnnounce *PolicyRules
}

type PolicyRules struct {
	AutomaticApproval PolicyValues `json:"Always,omitempty"`
	ManualApproval    PolicyValues `json:"WithApproval,omitempty"`
}
