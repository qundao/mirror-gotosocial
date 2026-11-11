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

package dereferencing

import (
	"context"

	"code.superseriousbusiness.org/gotosocial/internal/gtsmodel"
	"code.superseriousbusiness.org/gotosocial/internal/transport"
)

// TODO: at some point i'll move all of the below code
// into gtscontext or in same place we can reduce the
// number of times we have to re-fetch transport types.
// or, alternatively, overhaul the controller type!

type transportContext struct {
	context.Context
	account *gtsmodel.Account
	tsport  *transport.Transport
}

func (ctx *transportContext) Value(key any) any {
	requester, ok := key.(*gtsmodel.Account)
	if ok && ctx.account == requester {
		return ctx.tsport
	}
	return ctx.Context.Value(key)
}

// getTransport returns a transport for the requesting account, caching it in the returned context when newly fetched.
func (d *Dereferencer) getTransport(ctx context.Context, requester *gtsmodel.Account) (context.Context, *transport.Transport) {
	if tsport, ok := ctx.Value(requester).(*transport.Transport); ok {
		return ctx, tsport
	}

	// Load new transport instance for the passed requesting account.
	tsport := d.transportController.GetTransportForAccount(requester)

	// Store the transport in
	// a new wrapping context.
	ctx = &transportContext{
		Context: ctx,
		account: requester,
		tsport:  tsport,
	}

	return ctx, tsport
}
