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

// DomainLimit models federation
// limitations put on a domain by an admin.
type DomainLimit struct {

	// ID of this item in the database.
	ID string `bun:"type:CHAR(26),pk,nullzero,notnull,unique"`

	// Domain to limit. Eg. 'whatever.com'.
	Domain string `bun:",nullzero,notnull,unique"`

	// ID of the account that created this limit.
	CreatedByAccountID string `bun:"type:CHAR(26),nullzero,notnull"`

	// Account corresponding
	// to createdByAccountID.
	//
	// Not stored in the db.
	CreatedByAccount *Account `bun:"-"`

	// Private comment on the
	// limit, viewable by admins.
	PrivateComment string `bun:",nullzero"`

	// Public comment on this perm,
	// viewable (optionally) by everyone.
	PublicComment string `bun:",nullzero"`

	// Policy to apply to media files
	// originating from the limited domain.
	MediaPolicy MediaPolicy `bun:",nullzero,notnull,default:1"`

	// Policy to apply to follow (requests)
	// originating from the limited domain.
	FollowsPolicy FollowsPolicy `bun:",nullzero,notnull,default:1"`

	// Policy to apply to statuses from
	// non-followed accounts on the limited domain.
	StatusesPolicy StatusesPolicy `bun:",nullzero,notnull,default:1"`

	// Policy to apply to non-followed
	// accounts on the limited domain.
	AccountsPolicy AccountsPolicy `bun:",nullzero,notnull,default:1"`

	// Content warning to prepend to statuses
	// originating from the limited domain.
	ContentWarning string `bun:",nullzero"`
}

type MediaPolicy enumType

const (
	MediaPolicyUnknown MediaPolicy = 0

	// Default media behavior for
	// domains that aren't limited.
	MediaPolicyNoAction MediaPolicy = 1

	// Mark all media attachments
	// from the limited domain as sensitive.
	MediaPolicyMarkSensitive MediaPolicy = 2

	// Do not download, thumbnail, or store any media
	// files from the limited domain. A direct link to
	// any media files rejected in this way will be
	// appended to the bottom of the status.
	MediaPolicyReject MediaPolicy = 3
)

// MediaReject returns true if this domain
// limit is not nil and its MediaPolicy
// says that media should be rejected.
func (l *DomainLimit) MediaReject() bool {
	return l != nil && l.MediaPolicy == MediaPolicyReject
}

// MediaMarkSensitive returns true if this
// domain limit is not nil and its MediaPolicy
// says that media should be marked sensitive.
func (l *DomainLimit) MediaMarkSensitive() bool {
	return l != nil && l.MediaPolicy == MediaPolicyMarkSensitive
}

type FollowsPolicy enumType

const (
	FollowsPolicyUnknown FollowsPolicy = 0

	// Default follows behavior for
	// domains that aren't limited.
	FollowsPolicyNoAction FollowsPolicy = 1

	// Always require manual approval for all
	// follows issuing from the limited domain,
	// even if they target unlocked accounts.
	FollowsPolicyManualApproval FollowsPolicy = 2

	// Only process follows from the limited
	// domain when they are "follow backs",
	// ie., the followee already follows or
	// follow-requests the would-be follower.
	FollowsPolicyRejectNonMutual FollowsPolicy = 3

	// Reject all follows coming from the limited
	// domain, even if they target unlocked accounts.
	FollowsPolicyRejectAll FollowsPolicy = 4
)

type StatusesPolicy enumType

const (
	StatusesPolicyUnknown StatusesPolicy = 0

	// Default behavior for statuses from
	// accounts on domains that aren't limited.
	StatusesPolicyNoAction StatusesPolicy = 1

	// Apply a warn filter to statuses by non-
	// followed accounts from the limited domain.
	//
	// Statuses filtered in this way will also not
	// be shown on public web views of a thread.
	StatusesPolicyFilterWarn StatusesPolicy = 2

	// Apply a hide filter to statuses by non
	// -followed accounts from the limited domain.
	//
	// Statuses filtered in this way will also not
	// be shown on public web views of a thread.
	StatusesPolicyFilterHide StatusesPolicy = 3
)

// StatusesFilter returns true if this
// domain limit is not nil and its
// StatusesPolicy says that statuses
// should be filtered (warn or hide).
func (l *DomainLimit) StatusesFilter() bool {
	return l != nil && (l.StatusesPolicy == StatusesPolicyFilterWarn || l.StatusesPolicy == StatusesPolicyFilterHide)
}

type AccountsPolicy enumType

const (
	AccountsPolicyUnknown AccountsPolicy = 0

	// Default behavior for accounts
	// on domains that aren't limited.
	AccountsPolicyNoAction AccountsPolicy = 1

	// Mute aka silence non-followed
	// accounts from the limited domain.
	AccountsPolicyMute AccountsPolicy = 2
)

// AccountsMute returns true if this domain
// limit is not nil and its AccountsPolicy
// says that accounts should be muted.
func (l *DomainLimit) AccountsMute() bool {
	return l != nil && l.AccountsPolicy == AccountsPolicyMute
}
