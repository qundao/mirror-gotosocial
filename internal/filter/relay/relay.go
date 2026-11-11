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

package relay

import (
	"context"
	"errors"
	"net/url"

	"code.superseriousbusiness.org/gopkg/log"
	"code.superseriousbusiness.org/gopkg/xslices"
	"code.superseriousbusiness.org/gotosocial/internal/db"
	"code.superseriousbusiness.org/gotosocial/internal/filter"
	"code.superseriousbusiness.org/gotosocial/internal/gtserror"
	"code.superseriousbusiness.org/gotosocial/internal/gtsmodel"
	"code.superseriousbusiness.org/gotosocial/internal/state"
)

// Filter packages logic for checking whether
// given statuses should be permitted to relay.
type Filter struct {
	state *state.State
}

func NewFilter(state *state.State) *Filter {
	return &Filter{state}
}

// MatchedByActor checks whether the given
// status is matched by the given relay actor,
// including checks of that actor's relay flags
// and relay matchers against the status text.
func (f *Filter) MatchedByActor(
	ctx context.Context,
	relayActor *gtsmodel.RelayActor,
	status *gtsmodel.Status,
	inReplyToAccountURI string,
) bool {
	// Check status
	// against flags.
	if !flagsOK(
		status,
		inReplyToAccountURI,
		relayActor.Flags,
	) {
		log.Trace(ctx, "flags not OK")
		return false
	}

	// Convert text to filterable fields,
	// check that against matchers.
	fields := filter.GetFilterableFields(status)
	if !matchersOK(
		fields,
		relayActor.Flags,
		relayActor.Matchers,
	) {
		log.Trace(ctx, "matchers not OK")
		return false
	}

	return true
}

// MatchedBySubscription checks whether the given status is
// matched by a relay subscription, and returns the first relay
// subscription that it matches with, for the caller's interest.
//
// In case of no match, then nil nil is returned.
func (f *Filter) MatchedBySubscription(
	ctx context.Context,
	relayAcct *gtsmodel.Account,
	status *gtsmodel.Status,
	inReplyToAccountURI string,
) (*gtsmodel.RelaySubscription, error) {
	// Get all relay subscriptions that target the relay account's URI.
	subscriptions, err := f.state.DB.GetRelaySubscriptionsByActorURI(ctx, relayAcct.URI)
	if err != nil && !errors.Is(err, db.ErrNoEntries) {
		err := gtserror.Newf("db error getting relay subscriptions: %w", err)
		return nil, err
	}

	if len(subscriptions) == 0 {
		// No subscriptions means
		// definitely not permitted.
		return nil, nil
	}

	// Convert text to filterable
	// fields once outside the loop.
	fields := filter.GetFilterableFields(status)

	// Check each subscription to find first match.
	for _, subscription := range subscriptions {
		if matchedByConnection(
			status,
			inReplyToAccountURI,
			subscription,
			fields,
		) {
			// It's a match for
			// this subscription!
			//
			// Return early, no need
			// for further checks.
			return subscription, nil
		}
	}

	// No match found
	// among subscriptions.
	return nil, nil
}

// FilteredPushActorURIs returns a deduplicated slice of relay actor URIs
// that this status should be pushed to, after checking each push
// connection owned by status author for a match with the given status.
//
// If status author has no relay push connections, or no relay
// push connections match the given status, this will return nil.
func (f *Filter) FilteredPushActorURIs(
	ctx context.Context,
	status *gtsmodel.Status,
) ([]*url.URL, error) {
	// Get all push connections owned by this account.
	pushes, err := f.state.DB.GetRelayPushesForAccountID(ctx, status.AccountID)
	if err != nil && !errors.Is(err, db.ErrNoEntries) {
		return nil, gtserror.Newf("db error checking relay pushes: %w", err)
	}

	l := len(pushes)
	if l == 0 {
		// No pushes for this account,
		// so definitely no actor URIs.
		return nil, nil
	}

	// Convert text to filterable
	// fields once outside the loop.
	fields := filter.GetFilterableFields(status)

	var inReplyToAccountURI string
	if status.InReplyToAccount != nil {
		inReplyToAccountURI = status.InReplyToAccount.URI
	}

	// Prepare map to deduplicate actor URIs.
	pushActorURIStrs := make(map[string]struct{}, l)

	// Slice for returning parsed URIs.
	pushActorURIs := make([]*url.URL, 0, l)

	// Check each push for a match.
	for _, push := range pushes {
		if _, ok := pushActorURIStrs[push.RelayActorURI]; ok {
			// Already pushing
			// to this actor.
			continue
		}

		if matchedByConnection(
			status,
			inReplyToAccountURI,
			push,
			fields,
		) {
			// It's a match, store this actor URI.
			pushActorURIStrs[push.RelayActorURI] = struct{}{}
			pushActorURI, err := url.Parse(push.RelayActorURI)
			if err != nil {
				err := gtserror.Newf("error parsing push actor URI: %w", err)
				return nil, err
			}
			pushActorURIs = append(pushActorURIs, pushActorURI)
		}
	}

	return pushActorURIs, nil
}

// matchedByConnection checks flags
// and matchers for the given relay
// connection (push or subscription).
func matchedByConnection(
	status *gtsmodel.Status,
	inReplyToAccountURI string,
	rc gtsmodel.RelayConnection,
	fields []string,
) bool {
	// Check status
	// against flags.
	flags := rc.GetFlags()
	if !flagsOK(
		status,
		inReplyToAccountURI,
		flags,
	) {
		return false
	}

	// Check fields
	// against matchers.
	matchers := rc.GetMatchers()
	return matchersOK(
		fields,
		flags,
		matchers,
	)
}

// flagsOK returns true if the given flags
// do not forbid relaying the given status
// with the given inReplyToAccountURI.
func flagsOK(
	status *gtsmodel.Status,
	inReplyToAccountURI string,
	flags gtsmodel.RelayFlags,
) bool {
	// Ensure valid visibility.
	switch vis := status.Visibility; vis {
	case gtsmodel.VisibilityPublic:
		if !flags.Public() {
			// Public status, but
			// public not included
			// in this subscription.
			return false
		}
	case gtsmodel.VisibilityUnlocked:
		if !flags.Unlisted() {
			// Unlisted status, but
			// unlisted not included
			// in this subscription.
			return false
		}
	default:
		// Only public and unlisted
		// statuses can be relayed.
		return false
	}

	sensitive := status.Flags.Sensitive()
	if sensitive && flags.IgnoreSensitive() {
		// Sensitive status ignored
		// by this subscription.
		return false
	}

	hasMedia := len(status.Attachments) != 0
	if hasMedia && flags.IgnoreMedia() {
		// Status with media ignored
		// by this subscription.
		return false
	}

	isNonSelfReply := inReplyToAccountURI != "" &&
		(inReplyToAccountURI != status.AccountURI)
	if isNonSelfReply && flags.IgnoreReplies() {
		// Non self-replies ignored
		// by this subscription.
		return false
	}

	// None of the flags
	// prevent relaying
	// so must be OK.
	return true
}

// matchersOK returns true if the given
// matchers permit matching the given
// status, for the given relayFlags.
func matchersOK(
	fields []string,
	flags gtsmodel.RelayFlags,
	matchers []*gtsmodel.RelayMatcher,
) bool {
	// Check exclude matchers first, as an exclude
	// match means we don't need to check anything else.
	matchersLen := len(matchers)

	excludeMatchers := make([]*gtsmodel.RelayMatcher, 0, matchersLen)
	excludeMatchers = xslices.GatherIf(
		excludeMatchers,
		matchers,
		func(m *gtsmodel.RelayMatcher) (*gtsmodel.RelayMatcher, bool) {
			return m, m.Flags.Exclude()
		},
	)

	// If any exclude matcher matches
	// any field in the status, it's
	// a no for this relay connection.
	for _, matcher := range excludeMatchers {
		for _, field := range fields {
			if matcher.Regexp.MatchString(field) {
				return false
			}
		}
	}

	// If there's no exclude match,
	// and this sub matches all by
	// default, then it's a match.
	if flags.MatchByDefault() {
		return true
	}

	// Check if there's an include
	// match that matches the fields.
	includeMatchers := make([]*gtsmodel.RelayMatcher, 0, matchersLen)
	includeMatchers = xslices.GatherIf(
		includeMatchers,
		matchers,
		func(m *gtsmodel.RelayMatcher) (*gtsmodel.RelayMatcher, bool) {
			return m, !m.Flags.Exclude()
		},
	)

	// If any include matcher matches
	// any field in the status, it's
	// a yes for this relay connection.
	for _, matcher := range includeMatchers {
		for _, field := range fields {
			if matcher.Regexp.MatchString(field) {
				return true
			}
		}
	}

	// No matching include matcher
	// and we've exhausted all other
	// ways of matching this status,
	// so there's no match.
	return false
}
