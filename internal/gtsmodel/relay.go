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

import (
	"fmt"
	"regexp"
	"strconv"

	"codeberg.org/gruf/go-byteutil"
)

// RelayConnection is the interface for
// a RelayPush or RelaySubscription model.
type RelayConnection interface {
	GetID() string
	GetAccountID() string
	GetRelayActorURI() string
	GetFlags() RelayFlags
	GetMatcherIDs() []string
	GetMatchers() []*RelayMatcher
}

// Small table to keep track of which relay actors
// have relayed which posts, using which announce
// wrapper URIs (if relayed using an Announce).
type RelayedURI struct {
	ID        string `bun:"type:CHAR(26),pk,nullzero,notnull,unique"`
	RelayURI  string `bun:",notnull,nullzero,unique:relayed_uris_relay_uri_status_uri_uniq"`
	StatusURI string `bun:",notnull,nullzero,unique:relayed_uris_relay_uri_status_uri_uniq"`
	BoostURI  string `bun:",nullzero"`
}

// RelayPush represents a user-level push
// connection targeting a remote relay actor.
type RelayPush struct {
	// ID of this item in the database.
	// Creation time is encoded in the ID.
	ID string `bun:"type:CHAR(26),pk,nullzero,notnull,unique"`

	// ID of the account to which this
	// relay push connection belongs.
	AccountID string `bun:"type:CHAR(26),notnull,nullzero"`

	// ActivityPub URI of the relay service actor,
	// eg., `https://relay.activitypub.ca/actor`
	RelayActorURI string `bun:",notnull,nullzero"`

	// Flags contains numerous boolean
	// flags for this relay connection.
	// Default = public.
	Flags RelayFlags `bun:",notnull,default:2"`

	// IDs of matchers that apply
	// to this push connection.
	MatcherIDs []string `bun:"matchers,array"`

	// Matchers corresponding to MatcherIDs.
	//
	// Not stored in the database.
	Matchers []*RelayMatcher `bun:"-"`
}

func (r *RelayPush) GetID() string {
	return r.ID
}

func (r *RelayPush) GetAccountID() string {
	return r.AccountID
}

func (r *RelayPush) GetRelayActorURI() string {
	return r.RelayActorURI
}

func (r *RelayPush) GetFlags() RelayFlags {
	return r.Flags
}

func (r *RelayPush) GetMatcherIDs() []string {
	return r.MatcherIDs
}

func (r *RelayPush) GetMatchers() []*RelayMatcher {
	return r.Matchers
}

// RelaySubscription represents an admin-created
// subscription to a remote relay actor.
type RelaySubscription struct {
	// ID of this item in the database.
	// Creation time is encoded in the ID.
	ID string `bun:"type:CHAR(26),pk,nullzero,notnull,unique"`

	// ID of the account that created this relay subscription.
	AccountID string `bun:"type:CHAR(26),notnull,nullzero"`

	// ActivityPub URI of the relay service actor,
	// eg., `https://relay.activitypub.ca/actor`
	RelayActorURI string `bun:",notnull,nullzero"`

	// Flags contains numerous boolean
	// flags for this relay connection.
	// Default = relay public posts.
	Flags RelayFlags `bun:",notnull,default:2"`

	// IDs of matchers that apply
	// to this relay subscription.
	MatcherIDs []string `bun:"matchers,array"`

	// Matchers corresponding to MatcherIDs.
	//
	// Not stored in the database.
	Matchers []*RelayMatcher `bun:"-"`
}

func (r *RelaySubscription) GetID() string {
	return r.ID
}

func (r *RelaySubscription) GetAccountID() string {
	return r.AccountID
}

func (r *RelaySubscription) GetRelayActorURI() string {
	return r.RelayActorURI
}

func (r *RelaySubscription) GetFlags() RelayFlags {
	return r.Flags
}

func (r *RelaySubscription) GetMatcherIDs() []string {
	return r.MatcherIDs
}

func (r *RelaySubscription) GetMatchers() []*RelayMatcher {
	return r.Matchers
}

// RelayActor represents a *local* relay actor on
// this instance, created by an instance admin.
type RelayActor struct {
	// ID of this item in the database.
	// Creation time is encoded in the ID.
	ID string `bun:"type:CHAR(26),pk,nullzero,notnull,unique"`

	// ID of the account that created this relay actor.
	CreatedByAccountID string `bun:"type:CHAR(26),notnull,nullzero"`

	// ActivityPub URI of the relay actor,
	// eg., `https://example.org/relays/some_relay_username`
	URI string `bun:",notnull,nullzero,unique"`

	// Account corresponding to URI.
	//
	// Not stored in the database.
	ActorAccount *Account `bun:"-"`

	// Flags contains numerous boolean
	// flags for this relay actor.
	//
	// Default = relay public posts.
	Flags RelayFlags `bun:",notnull,default:2"`

	// Boolean flags for config
	// of the relay actor account.
	//
	// Default = WebShowFollowers true.
	ActorAccountFlags RelayActorAccountFlags `bun:",notnull,default:2"`

	// IDs of matchers that apply
	// to this relay actor.
	MatcherIDs []string `bun:"matchers,array"`

	// Matchers corresponding to MatcherIDs.
	//
	// Not stored in the database.
	Matchers []*RelayMatcher `bun:"-"`
}

// RelayActorAccountFlag is the bit type for
// individual RelayActorAccountFlag members.
type RelayActorAccountFlag bitFieldType

const (
	// NOTE: THE FOLLOWING VALUES SHOULD NEVER
	// BE CHANGED WITHOUT PERFORMING A DATABASE
	// MIGRATION TO UPDATE OLD -> NEW BIT VALUES.

	// Show followers of this relay actor on
	// the web view of the relay actor account.
	RelayActorAccountFlagWebShowFollowers = 1 << 1
)

// String returns a human-readable form of RelayActorAccountFlag.
func (f RelayActorAccountFlag) String() string {
	switch f {
	case 0:
		return "unset"
	case RelayActorAccountFlagWebShowFollowers:
		return "web_show_followers"
	default:
		panic(fmt.Sprintf("invalid relay flag: %d", f))
	}
}

// RelayActorAccountFlags uses smallint bit field type
// to store a variety of different boolean
// flags for attached relay entity.
type RelayActorAccountFlags bitFieldType

// WebShowFollowers returns whether RelayActorAccountFlagWebShowFollowers is set.
func (f RelayActorAccountFlags) WebShowFollowers() bool {
	return f&RelayActorAccountFlags(RelayActorAccountFlagWebShowFollowers) != 0
}

// SetWebShowFollowers sets / unsets the RelayActorAccountFlagWebShowFollowers bit.
func (f *RelayActorAccountFlags) SetWebShowFollowers(ok bool) {
	if ok {
		*f |= RelayActorAccountFlags(RelayActorAccountFlagWebShowFollowers)
	} else {
		*f &= ^RelayActorAccountFlags(RelayActorAccountFlagWebShowFollowers)
	}
}

// String returns a single human-readable form of RelayActorAccountFlags.
func (f RelayActorAccountFlags) String() string {
	var buf byteutil.Buffer
	buf.B = append(buf.B, '{')
	buf.B = append(buf.B, "web_show_followers="...)
	buf.B = strconv.AppendBool(buf.B, f.WebShowFollowers())
	buf.B = append(buf.B, '}')
	return buf.String()
}

// RelayFlag is the bit type for
// individual RelayFlags members.
type RelayFlag bitFieldType

const (
	// NOTE: THE FOLLOWING VALUES SHOULD NEVER
	// BE CHANGED WITHOUT PERFORMING A DATABASE
	// MIGRATION TO UPDATE OLD -> NEW BIT VALUES.

	// RelayFlagPublic controls whether a relay
	// entity should include public statuses
	//
	// Default is true (include public statuses).
	RelayFlagPublic RelayFlag = 1 << 1

	// RelayFlagUnlisted controls whether a relay
	// entity should include unlisted statuses.
	//
	// Default is false (don't include unlisted statuses).
	RelayFlagUnlisted RelayFlag = 1 << 2

	// RelayFlagMatchByDefault controls whether a relay entity
	// should match included, non-ignored statuses by default.
	//
	// If set true, and no "exclude"-type matchers are set on the relay
	// entity, then all included, non-ignored statuses will be relayed.
	//
	// Default is false (don't match by default).
	RelayFlagMatchByDefault RelayFlag = 1 << 3

	// RelayFlagIgnoreSensitive controls whether a relay
	// entity should ignore statuses designated as
	// sensitive via a content warning or sensitive flag.
	//
	// Default is false (don't ignore sensitive).
	RelayFlagIgnoreSensitive RelayFlag = 1 << 4

	// RelayFlagIgnoreMedia controls whether a relay
	// entity should ignore statuses with media.
	//
	// Default is false (don't ignore statuses with media).
	RelayFlagIgnoreMedia RelayFlag = 1 << 5

	// RelayFlagIgnoreReplies controls whether a relay entity
	// should ignore replies that aren't self-replies in a thread.
	//
	// Default is false (don't ignore replies).
	RelayFlagIgnoreReplies RelayFlag = 1 << 6
)

// String returns a human-readable form of RelayFlag.
func (f RelayFlag) String() string {
	switch f {
	case 0:
		return "unset"
	case RelayFlagPublic:
		return "public"
	case RelayFlagUnlisted:
		return "unlisted"
	case RelayFlagMatchByDefault:
		return "match_by_default"
	case RelayFlagIgnoreSensitive:
		return "ignore_sensitive"
	case RelayFlagIgnoreMedia:
		return "ignore_media"
	case RelayFlagIgnoreReplies:
		return "ignore_replies"
	default:
		panic(fmt.Sprintf("invalid relay flag: %d", f))
	}
}

// RelayFlags uses smallint bit field type
// to store a variety of different boolean
// flags for attached relay entity.
type RelayFlags bitFieldType

// Public returns whether RelayFlagPublic is set.
func (f RelayFlags) Public() bool {
	return f&RelayFlags(RelayFlagPublic) != 0
}

// SetPublic sets / unsets the RelayFlagPublic bit.
func (f *RelayFlags) SetPublic(ok bool) {
	if ok {
		*f |= RelayFlags(RelayFlagPublic)
	} else {
		*f &= ^RelayFlags(RelayFlagPublic)
	}
}

// Unlisted returns whether RelayFlagUnlisted is set.
func (f RelayFlags) Unlisted() bool {
	return f&RelayFlags(RelayFlagUnlisted) != 0
}

// SetUnlisted sets / unsets the RelayFlagUnlisted bit.
func (f *RelayFlags) SetUnlisted(ok bool) {
	if ok {
		*f |= RelayFlags(RelayFlagUnlisted)
	} else {
		*f &= ^RelayFlags(RelayFlagUnlisted)
	}
}

// MatchByDefault returns whether RelayFlagMatchByDefault is set.
func (f RelayFlags) MatchByDefault() bool {
	return f&RelayFlags(RelayFlagMatchByDefault) != 0
}

// SetMatchByDefault sets / unsets the RelayFlagMatchByDefault bit.
func (f *RelayFlags) SetMatchByDefault(ok bool) {
	if ok {
		*f |= RelayFlags(RelayFlagMatchByDefault)
	} else {
		*f &= ^RelayFlags(RelayFlagMatchByDefault)
	}
}

// IgnoreSensitive returns whether RelayFlagIgnoreSensitive is set.
func (f RelayFlags) IgnoreSensitive() bool {
	return f&RelayFlags(RelayFlagIgnoreSensitive) != 0
}

// SetIgnoreSensitive sets / unsets the RelayFlagIgnoreSensitive bit.
func (f *RelayFlags) SetIgnoreSensitive(ok bool) {
	if ok {
		*f |= RelayFlags(RelayFlagIgnoreSensitive)
	} else {
		*f &= ^RelayFlags(RelayFlagIgnoreSensitive)
	}
}

// IgnoreMedia returns whether RelayFlagIgnoreMedia is set.
func (f RelayFlags) IgnoreMedia() bool {
	return f&RelayFlags(RelayFlagIgnoreMedia) != 0
}

// SetIgnoreMedia sets / unsets the RelayFlagIgnoreMedia bit.
func (f *RelayFlags) SetIgnoreMedia(ok bool) {
	if ok {
		*f |= RelayFlags(RelayFlagIgnoreMedia)
	} else {
		*f &= ^RelayFlags(RelayFlagIgnoreMedia)
	}
}

// IgnoreReplies returns whether RelayFlagIgnoreReplies is set.
func (f RelayFlags) IgnoreReplies() bool {
	return f&RelayFlags(RelayFlagIgnoreReplies) != 0
}

// SetIgnoreReplies sets / unsets the RelayFlagIgnoreReplies bit.
func (f *RelayFlags) SetIgnoreReplies(ok bool) {
	if ok {
		*f |= RelayFlags(RelayFlagIgnoreReplies)
	} else {
		*f &= ^RelayFlags(RelayFlagIgnoreReplies)
	}
}

// String returns a single human-readable form of RelayFlags.
func (f RelayFlags) String() string {
	var buf byteutil.Buffer
	buf.B = append(buf.B, '{')
	buf.B = append(buf.B, "public="...)
	buf.B = strconv.AppendBool(buf.B, f.Public())
	buf.B = append(buf.B, ',')
	buf.B = append(buf.B, "unlisted="...)
	buf.B = strconv.AppendBool(buf.B, f.Unlisted())
	buf.B = append(buf.B, ',')
	buf.B = append(buf.B, "match_by_default="...)
	buf.B = strconv.AppendBool(buf.B, f.MatchByDefault())
	buf.B = append(buf.B, ',')
	buf.B = append(buf.B, "ignore_sensitive="...)
	buf.B = strconv.AppendBool(buf.B, f.IgnoreSensitive())
	buf.B = append(buf.B, ',')
	buf.B = append(buf.B, "ignore_media="...)
	buf.B = strconv.AppendBool(buf.B, f.IgnoreMedia())
	buf.B = append(buf.B, ',')
	buf.B = append(buf.B, "ignore_replies="...)
	buf.B = strconv.AppendBool(buf.B, f.IgnoreReplies())
	buf.B = append(buf.B, '}')
	return buf.String()
}

type RelayMatcher struct {
	// ID of this item in the database.
	// Creation time is encoded in the ID.
	ID string `bun:"type:CHAR(26),pk,nullzero,notnull,unique"`

	// ID of the relay entity (push, subscription,
	// or actor) to which this matcher applies.
	RelayID string `bun:"type:CHAR(26),notnull,nullzero,unique:relay_matchers_relay_id_keyword_uniq"`

	// Flags contains numerous boolean
	// flags for this relay matcher.
	Flags RelayMatcherFlags `bun:",notnull,default:0"`

	// The keyword or phrase to match.
	Keyword string `bun:",nullzero,notnull,unique:relay_matchers_relay_id_keyword_uniq"`

	// Precompiled regular expression.
	Regexp *regexp.Regexp `bun:"-"`
}

// RelayMatcherFlag is the bit type for
// individual RelayMatcherFlag members.
type RelayMatcherFlag bitFieldType

const (
	// NOTE: THE FOLLOWING VALUES SHOULD NEVER
	// BE CHANGED WITHOUT PERFORMING A DATABASE
	// MIGRATION TO UPDATE OLD -> NEW BIT VALUES.

	// If set, this relay matcher will cause matches to be
	// EXCLUDED from relaying rather than INCLUDED in relaying.
	RelayMatcherFlagExclude RelayMatcherFlag = 1 << 1

	// Consider word boundaries when matching.
	RelayMatcherFlagWholeWord RelayMatcherFlag = 1 << 2
)

// String returns a human-readable form of RelayMatcherFlag.
func (f RelayMatcherFlag) String() string {
	switch f {
	case 0:
		return "unset"
	case RelayMatcherFlagExclude:
		return "exclude"
	case RelayMatcherFlagWholeWord:
		return "whole_word"
	default:
		panic(fmt.Sprintf("invalid relay matcher flag: %d", f))
	}
}

// RelayMatcherFlags uses smallint bit field
// type to store a variety of different boolean
// flags for attached relay matcher.
type RelayMatcherFlags bitFieldType

// Exclude returns whether RelayMatcherFlagExclude is set.
func (f RelayMatcherFlags) Exclude() bool {
	return f&RelayMatcherFlags(RelayMatcherFlagExclude) != 0
}

// SetExclude sets / unsets the RelayMatcherFlagExclude bit.
func (f *RelayMatcherFlags) SetExclude(ok bool) {
	if ok {
		*f |= RelayMatcherFlags(RelayMatcherFlagExclude)
	} else {
		*f &= ^RelayMatcherFlags(RelayMatcherFlagExclude)
	}
}

// WholeWord returns whether RelayMatcherFlagWholeWord is set.
func (f RelayMatcherFlags) WholeWord() bool {
	return f&RelayMatcherFlags(RelayMatcherFlagWholeWord) != 0
}

// SetWholeWord sets / unsets the RelayMatcherFlagWholeWord bit.
func (f *RelayMatcherFlags) SetWholeWord(ok bool) {
	if ok {
		*f |= RelayMatcherFlags(RelayMatcherFlagWholeWord)
	} else {
		*f &= ^RelayMatcherFlags(RelayMatcherFlagWholeWord)
	}
}

// String returns a single human-readable form of RelayMatcherFlags.
func (f RelayMatcherFlags) String() string {
	var buf byteutil.Buffer
	buf.B = append(buf.B, '{')
	buf.B = append(buf.B, "exclude="...)
	buf.B = strconv.AppendBool(buf.B, f.Exclude())
	buf.B = append(buf.B, ',')
	buf.B = append(buf.B, "whole_word="...)
	buf.B = strconv.AppendBool(buf.B, f.WholeWord())
	buf.B = append(buf.B, '}')
	return buf.String()
}

// Compile will compile this RelayMatcher
// as a prepared regular expression.
func (r *RelayMatcher) Compile() (err error) {
	var (
		wordBreakStart string
		wordBreakEnd   string
	)

	if r.Flags.WholeWord() {
		// Either word boundary or
		// whitespace or start of line.
		wordBreakStart = `(?:\b|\s|^)`

		// Either word boundary or
		// whitespace or end of line.
		wordBreakEnd = `(?:\b|\s|$)`
	}

	// Compile keyword regexp.
	quoted := regexp.QuoteMeta(r.Keyword)
	r.Regexp, err = regexp.Compile(`(?i)` + wordBreakStart + quoted + wordBreakEnd)
	return // caller is expected to wrap this error
}
