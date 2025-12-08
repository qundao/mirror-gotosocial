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

package workers

import (
	"context"
	"errors"
	"slices"

	"code.superseriousbusiness.org/gopkg/log"
	"code.superseriousbusiness.org/gopkg/xslices"
	apimodel "code.superseriousbusiness.org/gotosocial/internal/api/model"
	"code.superseriousbusiness.org/gotosocial/internal/db"
	"code.superseriousbusiness.org/gotosocial/internal/filter/visibility"
	"code.superseriousbusiness.org/gotosocial/internal/gtscontext"
	"code.superseriousbusiness.org/gotosocial/internal/gtserror"
	"code.superseriousbusiness.org/gotosocial/internal/gtsmodel"
	"code.superseriousbusiness.org/gotosocial/internal/stream"
	"code.superseriousbusiness.org/gotosocial/internal/util"
)

// timelineAndNotifyStatus handles streaming a create event for the given status model to HOME, LIST, LOCAL
// and PUBLIC timelines, as well as adding to their relevant in-memory caches. It also handles sending any
// relevant notifications for the received status, e.g. mentions, followers with notify flag, conversations.
func (s *Surface) timelineAndNotifyStatus(ctx context.Context, status *gtsmodel.Status) error {

	// Ensure status fully populated; including account, mentions, etc.
	if err := s.State.DB.PopulateStatus(ctx, status); err != nil {
		return gtserror.Newf("error populating status with id %s: %w", status.ID, err)
	}

	// Local and public timeline caches
	// are global, i.e. *not* per-user,
	// so we only want to insert once.
	var localOnce, publicOnce bool

	// Timeline the status for local users
	// on the public and local timelines.
	s.timelineStatusForPublic(ctx, status,

		// local timelining and streaming function
		func(account *gtsmodel.Account, apiStatus *apimodel.Status) {
			if !localOnce {
				localOnce = true

				// Insert the status into the local timeline cache.
				_ = s.State.Caches.Timelines.Local.InsertOne(status)
			}

			// Stream the status model as local timeline update event.
			s.Stream.Update(ctx, account, apiStatus, stream.TimelineLocal)
		},

		// public timelining and streaming function
		func(account *gtsmodel.Account, apiStatus *apimodel.Status) {
			if !publicOnce {
				publicOnce = true

				// Insert the status into the public timeline cache.
				_ = s.State.Caches.Timelines.Public.InsertOne(status)
			}

			// Stream the status model as public timeline update event.
			s.Stream.Update(ctx, account, apiStatus, stream.TimelinePublic)
		},
	)

	// Timeline the status for each local follower of account, and each
	// local follower of any hashtags attached to status. This will also
	// handle notifying any followers with the 'notify' flag set.
	s.timelineAndNotifyStatusForFollowers(ctx, status,

		// home timelining and streaming function
		func(account *gtsmodel.Account, apiStatus *apimodel.Status) {

			// Insert this new status into the relevant list timeline cache.
			repeatBoost := s.State.Caches.Timelines.Home.InsertOne(account.ID, status)

			if !repeatBoost {
				// Only stream if not repeated boost of recent status.
				s.Stream.Update(ctx, account, apiStatus, stream.TimelineHome)
			}
		},

		// list timelining and streaming function
		func(list *gtsmodel.List, account *gtsmodel.Account, apiStatus *apimodel.Status) {

			// Insert this new status into the relevant list timeline cache.
			repeatBoost := s.State.Caches.Timelines.List.InsertOne(list.ID, status)

			if !repeatBoost {
				// Only stream if not repeated boost of recent status.
				streamType := stream.TimelineList + ":" + list.ID
				s.Stream.Update(ctx, account, apiStatus, streamType)
			}
		},

		// notify status for account function
		func(account *gtsmodel.Account) {
			if err := s.Notify(ctx,
				gtsmodel.NotificationStatus,
				account,
				status.Account,
				status,
				nil,
			); err != nil {
				log.Errorf(ctx, "error notifying status for account %s: %v", account.URI, err)
			}
		},
	)

	// Append to any tag timelines.
	s.timelineStatusForTags(status)

	// Notify each local account mentioned by status.
	if err := s.notifyMentions(ctx, status); err != nil {
		return gtserror.Newf("error notifying status mentions for status %s: %w", status.URI, err)
	}

	// Update conversations containing this status, and get notifications for them.
	notifications, err := s.Conversations.UpdateConversationsForStatus(ctx, status)
	if err != nil {
		return gtserror.Newf("error updating conversations for status %s: %w", status.URI, err)
	}

	// Stream these conversation notfications.
	for _, notification := range notifications {
		s.Stream.Conversation(ctx, notification.AccountID, notification.Conversation)
	}

	return nil
}

// timelineAndNotifyStatusUpdate handles streaming an update event for the given status model to HOME,
// LIST, LOCAL and PUBLIC timelines, as well as adding to their relevant in-memory caches. It also handles
// sending any relevant notifications for the received updated status, e.g. new mentions, poll closing,
// followers with the notify flag, and edits to a status that anyone local has previously interacted with.
func (s *Surface) timelineAndNotifyStatusUpdate(ctx context.Context, status *gtsmodel.Status) error {

	// Ensure fully populated; including account, mentions, etc.
	if err := s.State.DB.PopulateStatus(ctx, status); err != nil {
		return gtserror.Newf("error populating status with id %s: %w", status.ID, err)
	}

	if status.Poll != nil && status.Poll.Closing {
		// If latest status has a newly closed poll, at least compared
		// to the existing version, then notify poll close to all voters.
		if err := s.notifyPollClose(ctx, status); err != nil {
			log.Errorf(ctx, "error notifying poll close for status %s: %v", status.URI, err)
		}
	}

	// Notify account function for this status update
	// event. ONLY set if an edit was received AND we
	// successfully populated them from the database.
	var notifyAccount func(*gtsmodel.Account)

	// Ensure edits are fully populated for this status before anything.
	if err := s.State.DB.PopulateStatusEdits(ctx, status); err != nil {

		// we can still continue from here, just without
		// notifying local followers for it below here.
		log.Error(ctx, "error populating updated status edits: %v")

	} else if len(status.Edits) > 0 {
		// Track accounts we've already notified this
		// status for, as we can notify for both those
		// having interacted with status, AND those
		// that follow account with 'notify' flag set.
		notified := make(map[string]struct{})

		// Don't ever notify the status author.
		notified[status.AccountID] = struct{}{}

		// Get latest edit and notify for passed account.
		latestEdit := status.Edits[len(status.Edits)-1]
		notifyAccount = func(account *gtsmodel.Account) {
			if _, ok := notified[account.ID]; ok {
				return
			}

			// Mark account has already notified.
			notified[account.ID] = struct{}{}

			// Send notif for account.
			if err := s.Notify(ctx,
				gtsmodel.NotificationUpdate,
				account,
				status.Account,
				status,
				latestEdit,
			); err != nil {
				log.Errorf(ctx, "error notifying edit for account %s: %v", account.URI, err)
			}
		}
	}

	// Timeline the status for local users
	// on the public and local timelines.
	s.timelineStatusForPublic(ctx, status,

		// local timelining and streaming function
		func(account *gtsmodel.Account, apiStatus *apimodel.Status) {
			// NOTE: timeline invalidation is handled separately
			// as we don't need to perform it per user account.
			s.Stream.StatusUpdate(ctx, account, apiStatus, stream.TimelineLocal)
		},

		// public timelining and streaming function
		func(account *gtsmodel.Account, apiStatus *apimodel.Status) {
			// NOTE: timeline invalidation is handled separately
			// as we don't need to perform it per user account.
			s.Stream.StatusUpdate(ctx, account, apiStatus, stream.TimelinePublic)
		},
	)

	// Timeline the status update for each local follower of account,
	// and each local follower of any hashtags attached to status.
	s.timelineAndNotifyStatusForFollowers(ctx, status,

		// home timelining and streaming function
		func(account *gtsmodel.Account, apiStatus *apimodel.Status) {
			// NOTE: timeline invalidation is handled separately
			// as we don't need to perform it per account or list.
			s.Stream.StatusUpdate(ctx, account, apiStatus, stream.TimelineHome)
		},

		// list timelining and streaming function
		func(list *gtsmodel.List, account *gtsmodel.Account, apiStatus *apimodel.Status) {
			// NOTE: timeline invalidation is handled separately
			// as we don't need to perform it per account or list.
			streamType := stream.TimelineList + ":" + list.ID
			s.Stream.StatusUpdate(ctx, account, apiStatus, streamType)
		},

		// notify status for
		// account function
		notifyAccount,
	)

	// Notify any *new* mentions added by editor.
	for _, mention := range status.Mentions {

		// Check if we've seen
		// this mention already.
		if !mention.IsNew {
			continue
		}

		// Haven't seen this mention
		// yet, notify it if necessary.
		mention.Status = status
		if err := s.notifyMention(ctx, mention); err != nil {
			log.Errorf(ctx, "error notifying mention for status %s: %v", status.URI, err)
		}
	}

	if notifyAccount == nil {
		// We can only continue with further notification
		// of status edit if function was set, else return.
		return nil
	}

	// Get local-only interactions (we can't notify remotes).
	interactions, err := s.State.DB.GetStatusInteractions(ctx,
		status.ID,
		true, // local-only
	)
	if err != nil && !errors.Is(err, db.ErrNoEntries) {
		return gtserror.Newf("db error getting interactions for status %s: %w", status.URI, err)
	}

	// Notify all status interactees.
	for _, i := range interactions {
		targetAcct := i.GetAccount()
		notifyAccount(targetAcct)
	}

	return nil
}

// timelineStatusForPublic timelines the given status
// to LOCAL and PUBLIC (i.e. federated) timelines.
//
// much of the core logic is handled by functions passed as arguments
// to usage of this function with both creation and update events.
func (s *Surface) timelineStatusForPublic(
	ctx context.Context,
	status *gtsmodel.Status,
	localTimelineFn func(*gtsmodel.Account, *apimodel.Status),
	publicTimelineFn func(*gtsmodel.Account, *apimodel.Status),
) {
	if localTimelineFn == nil || publicTimelineFn == nil {
		panic("nil timeline func(s)")
	}

	if status.Visibility != gtsmodel.VisibilityPublic ||
		status.BoostOfID != "" {
		// Fast code path, if it's not "public"
		// or a boost, don't public timeline it.
		return
	}

	// Get a list of all local instance users.
	users, err := s.State.DB.GetAllUsers(ctx)
	if err != nil {
		log.Errorf(ctx, "db error getting local users: %v", err)
		return
	}

	// Iterate our list of users.
	isLocal := status.IsLocal()
	for _, user := range users {
		// Try to prepare status for timelining for local user account.
		apiStatus, timelineable, err := s.prepareStatusForTimeline(ctx,
			user.Account,
			status,
			gtsmodel.FilterContextPublic,
			(*visibility.Filter).StatusPublicTimelineable,
		)
		if err != nil {
			log.Errorf(ctx, "error preparing status %s for local user %s: %v", status.URI, user.Account.URI, err)
			continue
		}

		if !timelineable {
			continue
		}

		if isLocal {
			// This is local status, send it to local.
			localTimelineFn(user.Account, apiStatus)
		}

		// Both local and remote get sent to public.
		publicTimelineFn(user.Account, apiStatus)
	}
}

// timelineAndNotifyStatusForFollowers timelines and notifies (where
// appropriate) of the given status for all the local followers of the author
// account, and all the local followers the tags contained in status.
//
// much of the core logic is handled by functions passed as arguments
// to usage of this function with both creation and update events.
func (s *Surface) timelineAndNotifyStatusForFollowers(
	ctx context.Context,
	status *gtsmodel.Status,
	homeTimelineFn func(*gtsmodel.Account, *apimodel.Status),
	listTimelineFn func(*gtsmodel.List, *gtsmodel.Account, *apimodel.Status),
	notifyFn func(*gtsmodel.Account), // optional
) {
	if homeTimelineFn == nil || listTimelineFn == nil {
		panic("nil timeline func(s)")
	}

	// Get all local followers of the account that posted the status.
	follows, err := s.State.DB.GetAccountLocalFollowers(ctx, status.AccountID)
	if err != nil {
		log.Errorf(ctx, "db error getting local followers of account %s: %v", status.Account.URI, err)
		return
	}

	// If the poster is also local, add a fake entry for them
	// so they can see their own status in their timeline.
	if status.Account.IsLocal() {
		follows = append(follows, &gtsmodel.Follow{
			AccountID:   status.AccountID,
			Account:     status.Account,
			Notify:      util.Ptr(false), // Account shouldn't notify itself.
			ShowReblogs: util.Ptr(true),  // Account should show own reblogs.
		})
	}

	var (
		boost = (status.BoostOfID != "")
		reply = (status.InReplyToURI != "")
	)

	// Store a map of accounts the status has already
	// had home timeline processing applied for. This
	// is later used when determining whether to include
	// status in home timeline according to followed tags.
	//
	// Results:
	// - exists => already timelined OR not visible / muted
	// - empty  => not yet processed for home timeline
	processed := make(map[string]struct{}, len(follows))

	for _, follow := range follows {
		// Try to prepare this status for timelining for follow's account.
		apiStatus, timelineable, err := s.prepareStatusForTimeline(ctx,
			follow.Account,
			status,
			gtsmodel.FilterContextHome,
			(*visibility.Filter).StatusHomeTimelineable,
		)
		if err != nil {
			log.Error(ctx, err)
			continue
		}

		if !timelineable {
			// This status should not be timelined
			// for this account's home timelines.
			processed[follow.AccountID] = struct{}{}
			continue
		}

		// Get all lists that contain this given follow.
		lists, err := s.State.DB.GetListsContainingFollowID(
			gtscontext.SetBarebones(ctx), // no sub-models
			follow.ID,
		)
		if err != nil {
			log.Errorf(ctx, "error getting lists for follow %s: %v", follow.URI, err)
			continue
		}

		var exclusive bool
		for _, list := range lists {
			// Check whether list is eligible for this status.
			eligible, err := s.isListEligible(ctx, list, status)
			if err != nil {
				log.Errorf(ctx, "error checking list eligibility for status %s: %v", status.URI, err)
				continue
			}

			if !eligible {
				continue
			}

			// Update exclusive flag if list is so.
			exclusive = exclusive || *list.Exclusive

			// Timeline this status into account's list.
			listTimelineFn(list, follow.Account, apiStatus)
		}

		// If this was timelined into
		// list with exclusive flag set,
		// don't add to home timeline.
		if !exclusive {

			// Add status to account's home timeline.
			homeTimelineFn(follow.Account, apiStatus)

			// Mark as processed for home timeline in map.
			processed[follow.AccountID] = struct{}{}
		}

		if !*follow.Notify {
			// This follower doesn't have notifs
			// set for this account's new posts.
			continue
		}

		if boost || reply {
			// Don't notify for
			// boosts or replies.
			continue
		}

		if notifyFn != nil {
			// Notify for this follow.
			notifyFn(follow.Account)
		}
	}

	// From here, status has been sent to home and
	// list timelines based on follow relationships.
	// We now need to process for hashtag follows.
	tagStatus := status
	if tagStatus.BoostOf != nil {

		// Unwrap boost and work
		// with the original status.
		tagStatus = tagStatus.BoostOf
	}

	if tagStatus.Visibility != gtsmodel.VisibilityPublic {
		// Only public statuses are eligible
		// for hashtag follow inclusion.
		return
	}

	// Gather *useable* tag IDs from tag status.
	tagIDs := xslices.GatherIf(nil, tagStatus.Tags,
		func(tag *gtsmodel.Tag) (string, bool) {
			return tag.ID, (*tag.Useable)
		})

	if len(tagIDs) == 0 {
		// No tags to
		// act on.
		return
	}

	// Get the list of account IDs following determined useable tag IDs.
	accountIDs, err := s.State.DB.GetAccountIDsFollowingTagIDs(ctx, tagIDs)
	if err != nil {
		log.Errorf(ctx, "db error getting tag followers: %v", err)
		return
	}

	// Filter follower account IDs by home timelining
	// results, where any result indicates it has
	// already been processed for home timelineability.
	accountIDs = slices.DeleteFunc(accountIDs,
		func(accountID string) bool {
			_, ok := processed[accountID]
			return ok
		})

	if len(accountIDs) == 0 {
		// No accounts to
		// timeline for.
		return
	}

	// Fetch account models for enumerated IDs.
	accounts, err := s.State.DB.GetAccountsByIDs(
		gtscontext.SetBarebones(ctx),
		accountIDs,
	)
	if err != nil {
		log.Errorf(ctx, "db error getting accounts: %v", err)
		return
	}

	for _, account := range accounts {
		// Try to prepare status for timelining for tag follow's account.
		apiStatus, timelineable, err := s.prepareStatusForTimeline(ctx,
			account,
			status,
			gtsmodel.FilterContextHome,
			(*visibility.Filter).StatusVisible,
		)
		if err != nil {
			log.Errorf(ctx, "error preparing status %s for tag follower %s: %v", status.URI, account.URI, err)
			continue
		}

		if !timelineable {
			continue
		}

		// Add to account's home timeline.
		homeTimelineFn(account, apiStatus)
	}
}

// timelineStatusForTags attempts to insert given status into relevant tag timeline caches.
func (s *Surface) timelineStatusForTags(status *gtsmodel.Status) {
	if status.Visibility != gtsmodel.VisibilityPublic ||
		status.BoostOfID != "" {
		// Only include "public" non-boost
		// statuses in tag timelines.
		return
	}

	// Gather timelineable tag IDs from status.
	tagIDs := xslices.GatherIf(nil, status.Tags,
		func(tag *gtsmodel.Tag) (string, bool) {
			return tag.ID, (*tag.Useable) &&
				(*tag.Listable)
		})

	if len(tagIDs) == 0 {
		// No tags to
		// act on.
		return
	}

	for _, tagID := range tagIDs {
		// Insert new status into the relevant tag timeline cache.
		_ = s.State.Caches.Timelines.Tag.InsertOne(tagID, status)
	}
}

// prepareStatusForTimeline attempts to prepare the given status for
// a timeline owned by the given account, first passing it through
// appropriate visibility function, mute checks and status filtering
// checks applicable in the given filter context. finally, it will
// return a prepared frontend API model for timeline insertion.
func (s *Surface) prepareStatusForTimeline(
	ctx context.Context,
	account *gtsmodel.Account,
	status *gtsmodel.Status,
	filterCtx gtsmodel.FilterContext,
	isVisibleFn func(*visibility.Filter, context.Context, *gtsmodel.Account, *gtsmodel.Status) (bool, error),
) (
	apiStatus *apimodel.Status,
	timelineable bool,
	err error,
) {
	// Check status visibility for account's appropriate timeline.
	visible, err := isVisibleFn(s.VisFilter, ctx, account, status)
	if err != nil {
		return nil, false, gtserror.Newf("error checking status %s visibility: %w", status.URI, err)
	}

	if !visible {
		return nil, false, nil
	}

	// Check if the status muted by this account.
	muted, err := s.MuteFilter.StatusMuted(ctx,
		account,
		status,
	)
	if err != nil {
		return nil, false, gtserror.Newf("error checking status %s mute: %w", status.URI, err)
	}

	if muted {
		return nil, false, nil
	}

	// Check whether status is filtered in this context by timeline account.
	filtered, hide, err := s.StatusFilter.StatusFilterResultsInContext(ctx,
		account,
		status,
		filterCtx,
	)
	if err != nil {
		return nil, false, gtserror.Newf("error filtering status %s: %w", status.URI, err)
	}

	if hide {
		return nil, false, nil
	}

	// Attempt to convert status to frontend API model.
	apiStatus, err = s.Converter.StatusToAPIStatus(ctx,
		status,
		account,
	)
	if err != nil {
		log.Error(ctx, "error converting status %s to frontend: %v", status.URI, err)
	} else {

		// Attach any filter results.
		apiStatus.Filtered = filtered
	}

	return apiStatus, true, nil
}

// listEligible checks if the given status is eligible
// for inclusion in the list that that the given listEntry
// belongs to, based on the replies policy of the list.
func (s *Surface) isListEligible(
	ctx context.Context,
	list *gtsmodel.List,
	status *gtsmodel.Status,
) (bool, error) {
	if status.InReplyToURI == "" {
		// If status is not a reply,
		// then it's all gravy baby.
		return true, nil
	}

	if status.InReplyToID == "" {
		// Status is a reply but we don't
		// have the replied-to account!
		return false, nil
	}

	switch list.RepliesPolicy {
	case gtsmodel.RepliesPolicyNone:
		// This list should not show
		// replies at all, so skip it.
		return false, nil

	case gtsmodel.RepliesPolicyList:
		// This list should show replies
		// only to other people in the list.
		//
		// Check if replied-to account is
		// also included in this list.
		in, err := s.State.DB.IsAccountInList(ctx,
			list.ID,
			status.InReplyToAccountID,
		)
		if err != nil {
			err := gtserror.Newf("db error checking if account in list: %w", err)
			return false, err
		}
		return in, nil

	case gtsmodel.RepliesPolicyFollowed:
		// This list should show replies
		// only to people that the list
		// owner also follows.
		//
		// Check if replied-to account is
		// followed by list owner account.
		follows, err := s.State.DB.IsFollowing(ctx,
			list.AccountID,
			status.InReplyToAccountID,
		)
		if err != nil {
			err := gtserror.Newf("db error checking if account followed: %w", err)
			return false, err
		}
		return follows, nil

	default:
		panic("unknown reply policy: " + list.RepliesPolicy)
	}
}

// deleteStatusFromTimelines completely removes the given status from all timelines.
// It will also stream deletion of the status to all open streams.
func (s *Surface) deleteStatusFromTimelines(ctx context.Context, statusID string) {
	s.State.Caches.Timelines.Public.RemoveByStatusIDs(statusID)
	s.State.Caches.Timelines.Local.RemoveByStatusIDs(statusID)
	s.State.Caches.Timelines.Home.RemoveByStatusIDs(statusID)
	s.State.Caches.Timelines.List.RemoveByStatusIDs(statusID)
	s.State.Caches.Timelines.Tag.RemoveByStatusIDs(statusID)
	s.Stream.Delete(ctx, statusID)
}

// removeTimelineEntriesByAccount removes all cached timeline entries authored by account ID.
func (s *Surface) removeTimelineEntriesByAccount(accountID string) {
	s.State.Caches.Timelines.Public.RemoveByAccountIDs(accountID)
	s.State.Caches.Timelines.Local.RemoveByAccountIDs(accountID)
	s.State.Caches.Timelines.Home.RemoveByAccountIDs(accountID)
	s.State.Caches.Timelines.List.RemoveByAccountIDs(accountID)
	s.State.Caches.Timelines.Tag.RemoveByAccountIDs(accountID)
}

func (s *Surface) removeRelationshipFromTimelines(ctx context.Context, timelineAccountID string, targetAccountID string) {
	// Remove all statuses by target account
	// from given account's home timeline.
	s.State.Caches.Timelines.Home.
		MustGet(timelineAccountID).
		RemoveByAccountIDs(targetAccountID)

	// Get the IDs of all the lists owned by the given account ID.
	listIDs, err := s.State.DB.GetListIDsByAccountID(ctx, timelineAccountID)
	if err != nil {
		log.Errorf(ctx, "error getting lists for account %s: %v", timelineAccountID, err)
	}

	for _, listID := range listIDs {
		// Remove all statuses by target account
		// from given account's list timelines.
		s.State.Caches.Timelines.List.MustGet(listID).
			RemoveByAccountIDs(targetAccountID)
	}
}
