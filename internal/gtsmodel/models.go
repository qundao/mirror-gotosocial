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

// AllModels returns a slice of all currently
// known database models, for registration into
// bun's reflect-based ORM system.
func AllModels() []any {
	return []any{

		// m2m models must
		// be up first for
		// db registration.
		&AccountToEmoji{},
		&ConversationToStatus{},
		&StatusToEmoji{},
		&StatusToTag{},

		&Account{},
		&AccountNote{},
		&AccountSettings{},
		&AccountStats{},
		&AdminAction{},
		&AdvancedMigration{},
		&Application{},
		&Block{},
		&Conversation{},
		&DeniedUser{},
		&DomainAllow{},
		&DomainBlock{},
		&DomainLimit{},
		&DomainPermissionDraft{},
		&DomainPermissionExclude{},
		&DomainPermissionSubscription{},
		&EmailDomainBlock{},
		&Emoji{},
		&EmojiCategory{},
		&FederationError{},
		&Filter{},
		&FilterKeyword{},
		&FilterStatus{},
		&Follow{},
		&FollowedTag{},
		&FollowRequest{},
		&HeaderFilterAllow{},
		&HeaderFilterBlock{},
		&Instance{},
		&InstanceSettings{},
		&InteractionRequest{},
		&List{},
		&ListEntry{},
		&Marker{},
		&MediaAttachment{},
		&Mention{},
		&Move{},
		&Notification{},
		&Poll{},
		&PollVote{},
		&RelayActor{},
		&RelayMatcher{},
		&RelayPush{},
		&RelaySubscription{},
		&RelayedURI{},
		&Report{},
		&RouterSession{},
		&Rule{},
		&ScheduledStatus{},
		&Status{},
		&StatusBookmark{},
		&StatusEdit{},
		&StatusFave{},
		&StatusPin{},
		&Tag{},
		&Thread{},
		&ThreadMute{},
		&Token{},
		&Tombstone{},
		&User{},
		&UserMute{},
		&VAPIDKeyPair{},
		&WebPushSubscription{},
		&WorkerTask{},
	}
}
