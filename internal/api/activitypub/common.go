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

package activitypub

import (
	"strings"

	"code.superseriousbusiness.org/gopkg/httputil"
	"code.superseriousbusiness.org/gotosocial/internal/gtserror"
	"code.superseriousbusiness.org/gotosocial/internal/uris"

	apiutil "code.superseriousbusiness.org/gotosocial/internal/api/util"
)

// ParseCommon parses out path variables and content type
// from a request to one of our ActivityPub API endpoints.
func ParseCommon(c *httputil.Context) (
	relay bool,
	username string,
	contentType string,
	errWithCode gtserror.WithCode,
) {
	// Get username from request params.
	username, errWithCode = apiutil.ParseUsername(c.PathValue(apiutil.UsernameKey))
	if errWithCode != nil {
		return
	}

	// If this is a relay path, prefix the username with
	// `relay.` to differentiate it from an ordinary user.
	const relayPrefix = "/" + uris.RelaysPath
	if strings.HasPrefix(c.R.URL.Path, relayPrefix) {
		relay = true
		username = uris.RelayUsernamePrefix + username
	}

	// Get content type.
	var err error
	contentType, err = apiutil.NegotiateAccept(c, apiutil.ActivityPubOrHTMLHeaders...)
	if err != nil {
		errWithCode = gtserror.NewErrorNotAcceptable(err, err.Error())
		return
	}

	return
}

// ParseCommonWithID is like ParseCommon,
// but includes a ULID ID found in the path.
func ParseCommonWithID(c *httputil.Context) (
	relay bool,
	username string,
	id string,
	contentType string,
	errWithCode gtserror.WithCode,
) {
	// Do parsecommon to get username + content type.
	relay, username, contentType, errWithCode = ParseCommon(c)
	if errWithCode != nil {
		return
	}

	// Additionally get ID from request params.
	id, errWithCode = apiutil.ParseID(c.PathValue(apiutil.IDKey))
	if errWithCode != nil {
		return
	}

	return
}
