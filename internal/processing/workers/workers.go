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
	"code.superseriousbusiness.org/gotosocial/internal/federation"
	"code.superseriousbusiness.org/gotosocial/internal/processing/account"
	"code.superseriousbusiness.org/gotosocial/internal/processing/media"
	"code.superseriousbusiness.org/gotosocial/internal/state"
	"code.superseriousbusiness.org/gotosocial/internal/surfacing"
	"code.superseriousbusiness.org/gotosocial/internal/typeutils"
	"code.superseriousbusiness.org/gotosocial/internal/workers"
)

type Processor struct {
	clientAPI clientAPI
	fediAPI   fediAPI
	workers   *workers.Workers
}

func New(
	state *state.State,
	federator *federation.Federator,
	converter *typeutils.Converter,
	surfacer *surfacing.Surfacer,
	account *account.Processor,
	media *media.Processor,
) Processor {
	// Init federate logic
	// wrapper struct.
	federate := &federate{
		Federator: federator,
		state:     state,
		converter: converter,
	}

	// Init shared util funcs.
	utils := &utils{
		state:     state,
		media:     media,
		account:   account,
		surfacer:  surfacer,
		converter: converter,
	}

	return Processor{
		workers: &state.Workers,
		clientAPI: clientAPI{
			state:    state,
			surfacer: surfacer,
			federate: federate,
			account:  account,
			utils:    utils,
		},
		fediAPI: fediAPI{
			state:    state,
			surfacer: surfacer,
			federate: federate,
			account:  account,
			utils:    utils,
		},
	}
}
