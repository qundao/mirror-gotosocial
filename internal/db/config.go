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

package db

import (
	"context"

	"code.superseriousbusiness.org/gotosocial/internal/gtsmodel"
)

type Config interface {
	// RuntimeConfig retrieves the service runtime
	// config or panics, since being unable to retrieve
	// runtime config is an "oh fuck" type of event.
	RuntimeConfig(ctx context.Context) *gtsmodel.RuntimeConfig

	// UpdateRuntimeConfig updates the given columns of the service runtime config.
	UpdateRuntimeConfig(ctx context.Context, config *gtsmodel.RuntimeConfig, columns ...string) error
}
