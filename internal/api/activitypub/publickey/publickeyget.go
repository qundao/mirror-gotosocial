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

package publickey

import (
	"net/http"
	"strings"

	"code.superseriousbusiness.org/gopkg/httputil"
	apiutil "code.superseriousbusiness.org/gotosocial/internal/api/util"
	"code.superseriousbusiness.org/gotosocial/internal/gtserror"
	"code.superseriousbusiness.org/gotosocial/internal/uris"
)

// PublicKeyGETHandler should be served at eg https://example.org/{users|relays}/:username/main-key.
//
// The goal here is to return a MINIMAL activitypub representation of an account
// in the form of a vocab.ActivityStreamsPerson. The account will only contain the id,
// public key, username, and type of the account.
func (m *Module) PublicKeyGETHandler(c *httputil.Context) {

	// Get username from request params.
	username, errWithCode := apiutil.ParseUsername(c.PathValue(apiutil.UsernameKey))
	if errWithCode != nil {
		return
	}

	// If this is a relay path, prefix the username with
	// `relay.` to differentiate it from an ordinary user.
	pathTrimmed := strings.TrimPrefix(c.R.URL.Path, "/")
	if strings.HasPrefix(pathTrimmed, uris.RelaysPath) {
		username = uris.RelayUsernamePrefix + username
	}

	contentType, err := apiutil.NegotiateAccept(c, apiutil.ActivityPubOrHTMLHeaders...)
	if err != nil {
		apiutil.ErrorHandler(c, m.templates, gtserror.NewErrorNotAcceptable(err, err.Error()))
		return
	}

	// If HTML is requested, redirect to profile.
	if contentType == string(apiutil.TextHTML) {
		httputil.Redirect(c, http.StatusSeeOther, "/@"+username)
		return
	}

	resp, errWithCode := m.processor.Fedi().UserGetMinimal(
		c,
		username,
	)
	if errWithCode != nil {
		apiutil.ErrorHandler(c, m.templates, errWithCode)
		return
	}

	// Encode JSON response.
	httputil.JSONType(c,
		http.StatusOK,
		contentType,
		resp,
	)
}
