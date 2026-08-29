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
	"strconv"
	"time"

	"codeberg.org/gruf/go-byteutil"
)

// smallint is the largest size supported
// by a PostgreSQL SMALLINT, since an SQLite
// SMALLINT is actually variable in size.
type smallint int16

// bitFieldType is the type we use
// for database int bit fields, at
// least where the smallest int size
// will suffice for number of fields.
type bitFieldType smallint

// FollowFlag is the bit type for
// individual FollowFlags members.
type FollowFlag bitFieldType

const (
	// NOTE: THE FOLLOWING VALUES SHOULD NEVER
	// BE CHANGED WITHOUT PERFORMING A DATABASE
	// MIGRATION TO UPDATE OLD -> NEW BIT VALUES.

	// FollowFlagShowReblogs controls whether reblogs should
	// be shown in the home timeline for followed account.
	//
	// Default is true.
	FollowFlagShowReblogs FollowFlag = 1 << 1

	// FollowFlagNotify controls whether the following
	// account should be notified when followed account posts.
	//
	// Default is false.
	FollowFlagNotify FollowFlag = 1 << 2

	// FollowFlagUsePublicURI indicates that the follow
	// (request) originally used the ActivityPub Public
	// URI as the object of the Follow, rather than the
	// AP ID/URI of the Followed actor.
	//
	// This is used for relay follows from Masto / *key.
	//
	// Default is false.
	FollowFlagUsePublicURI FollowFlag = 1 << 3
)

// String returns a human-readable form of FollowFlag.
func (f FollowFlag) String() string {
	switch f {
	case 0:
		return "unset"
	case FollowFlagShowReblogs:
		return "show_reblogs"
	case FollowFlagNotify:
		return "notify"
	case FollowFlagUsePublicURI:
		return "use_public_uri"
	default:
		panic(fmt.Sprintf("invalid follow flag: %d", f))
	}
}

// FollowFlags uses smallint bit field type
// to store a variety of different boolean
// flags for attached follow (request).
type FollowFlags bitFieldType

// ShowReblogs returns whether FollowFlagShowReblogs is set.
func (f FollowFlags) ShowReblogs() bool {
	return f&FollowFlags(FollowFlagShowReblogs) != 0
}

// SetShowReblogs sets / unsets the FollowFlagShowReblogs bit.
func (f *FollowFlags) SetShowReblogs(ok bool) {
	if ok {
		*f |= FollowFlags(FollowFlagShowReblogs)
	} else {
		*f &= ^FollowFlags(FollowFlagShowReblogs)
	}
}

// Notify returns whether FollowFlagNotify is set.
func (f FollowFlags) Notify() bool {
	return f&FollowFlags(FollowFlagNotify) != 0
}

// SetNotify sets / unsets the FollowFlagNotify bit.
func (f *FollowFlags) SetNotify(ok bool) {
	if ok {
		*f |= FollowFlags(FollowFlagNotify)
	} else {
		*f &= ^FollowFlags(FollowFlagNotify)
	}
}

// UsePublicURI returns whether FollowFlagUsePublicURI is set.
func (f FollowFlags) UsePublicURI() bool {
	return f&FollowFlags(FollowFlagUsePublicURI) != 0
}

// SetUsePublicURI sets / unsets the FollowFlagUsePublicURI bit.
func (f *FollowFlags) SetUsePublicURI(ok bool) {
	if ok {
		*f |= FollowFlags(FollowFlagUsePublicURI)
	} else {
		*f &= ^FollowFlags(FollowFlagUsePublicURI)
	}
}

// String returns a single human-readable form of FollowFlags.
func (f FollowFlags) String() string {
	var buf byteutil.Buffer
	buf.B = append(buf.B, '{')
	buf.B = append(buf.B, "show_reblogs="...)
	buf.B = strconv.AppendBool(buf.B, f.ShowReblogs())
	buf.B = append(buf.B, ',')
	buf.B = append(buf.B, "notify="...)
	buf.B = strconv.AppendBool(buf.B, f.Notify())
	buf.B = append(buf.B, ',')
	buf.B = append(buf.B, "use_public_uri="...)
	buf.B = strconv.AppendBool(buf.B, f.UsePublicURI())
	buf.B = append(buf.B, '}')
	return buf.String()
}

// Follow represents one account following another,
// and the metadata around that follow.
type Follow struct {
	// ID of this item in the database.
	ID string `bun:"type:CHAR(26),pk,nullzero,notnull,unique"`

	// Time when the item was created.
	CreatedAt time.Time `bun:"type:timestamptz,nullzero,notnull,default:current_timestamp"`

	// Time when the item was last updated.
	UpdatedAt time.Time `bun:"type:timestamptz,nullzero,notnull,default:current_timestamp"`

	// URI of the ActivityPub Follow.
	URI string `bun:",notnull,nullzero,unique"`

	// ID of the follow origin account.
	AccountID string `bun:"type:CHAR(26),unique:srctarget,notnull,nullzero"`

	// ID of the follow target account.
	TargetAccountID string `bun:"type:CHAR(26),unique:srctarget,notnull,nullzero"`

	// Flags controlling Follow behavior.
	Flags FollowFlags `bun:",notnull,default:2"`
}

// FollowRequest represents one account requesting to
// follow another, and the metadata around that request.
type FollowRequest struct {
	// ID of this item in the database.
	ID string `bun:"type:CHAR(26),pk,nullzero,notnull,unique"`

	// Time when the item was created.
	CreatedAt time.Time `bun:"type:timestamptz,nullzero,notnull,default:current_timestamp"`

	// Time when the item was last updated.
	UpdatedAt time.Time `bun:"type:timestamptz,nullzero,notnull,default:current_timestamp"`

	// URI of the ActivityPub Follow.
	URI string `bun:",notnull,nullzero,unique"`

	// ID of the follow request origin account.
	AccountID string `bun:"type:CHAR(26),unique:frsrctarget,notnull,nullzero"`

	// ID of the follow request target account.
	TargetAccountID string `bun:"type:CHAR(26),unique:frsrctarget,notnull,nullzero"`

	// Flags controlling Follow behavior.
	Flags FollowFlags `bun:",notnull,default:2"`
}
