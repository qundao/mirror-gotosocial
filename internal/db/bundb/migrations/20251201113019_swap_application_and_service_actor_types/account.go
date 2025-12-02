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

// smallint is the largest size supported
// by a PostgreSQL SMALLINT, since an SQLite
// SMALLINT is actually variable in size.
type smallint int16

// enumType is the type we (at least, should) use
// for database enum types, as smallest int size.
type enumType smallint

// AccountActorType is the ActivityStreams type of an actor.
type AccountActorType enumType

const (
	AccountActorTypeUnknown      AccountActorType = 0
	AccountActorTypeApplication  AccountActorType = 1 // https://www.w3.org/TR/activitystreams-vocabulary/#dfn-application
	AccountActorTypeGroup        AccountActorType = 2 // https://www.w3.org/TR/activitystreams-vocabulary/#dfn-group
	AccountActorTypeOrganization AccountActorType = 3 // https://www.w3.org/TR/activitystreams-vocabulary/#dfn-organization
	AccountActorTypePerson       AccountActorType = 4 // https://www.w3.org/TR/activitystreams-vocabulary/#dfn-person
	AccountActorTypeService      AccountActorType = 5 // https://www.w3.org/TR/activitystreams-vocabulary/#dfn-service
)
