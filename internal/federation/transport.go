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

package federation

import (
	"context"
	"net/url"

	"code.superseriousbusiness.org/activity/pub"
	"code.superseriousbusiness.org/gotosocial/internal/gtserror"
	"code.superseriousbusiness.org/gotosocial/internal/uris"
)

// NewTransport returns a new Transport on behalf of a specific actor.
//
// The actorBoxIRI will be either the inbox or outbox of an actor who is
// attempting to do the dereferencing or delivery. Any authentication
// scheme applied on the request must be based on this actor. The
// request must contain some sort of credential of the user, such as a
// HTTP Signature.
//
// The gofedAgent passed in should be used by the Transport
// implementation in the User-Agent, as well as the application-specific
// user agent string. The gofedAgent will indicate this library's use as
// well as the library's version number.
//
// Any server-wide rate-limiting that needs to occur should happen in a
// Transport implementation. This factory function allows this to be
// created, so peer servers are not DOS'd.
//
// Any retry logic should also be handled by the Transport
// implementation.
//
// Note that the library will not maintain a long-lived pointer to the
// returned Transport so that any private credentials are able to be
// garbage collected.
func (f *Federator) NewTransport(ctx context.Context, actorBoxIRI *url.URL, _ string) (pub.Transport, error) {
	var username string
	var err error

	switch {
	case uris.IsInboxPath(actorBoxIRI):
		_, username, err = uris.ParseInboxPath(actorBoxIRI)
		if err != nil {
			return nil, gtserror.Newf("error parsing inbox path %s: %w", actorBoxIRI, err)
		}

	case uris.IsOutboxPath(actorBoxIRI):
		_, username, err = uris.ParseOutboxPath(actorBoxIRI)
		if err != nil {
			return nil, gtserror.Newf("error parsing outbox path %s: %w", actorBoxIRI, err)
		}

	default:
		return nil, gtserror.Newf("unexpected iri path type: %s", actorBoxIRI)
	}

	// Look for the provided local account by username.
	account, err := f.db.GetAccountByUsernameDomain(ctx,
		username,
		"",
	)
	if err != nil {
		return nil, gtserror.Newf("error getting local account with username %s: %w", username, err)
	}

	// Return a transport instance for the fetched account.
	return f.transport.GetTransportForAccount(account), nil
}
