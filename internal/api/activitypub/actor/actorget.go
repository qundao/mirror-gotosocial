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

package actor

import (
	"net/http"

	"code.superseriousbusiness.org/gopkg/httputil"
	"code.superseriousbusiness.org/gotosocial/internal/api/activitypub"
	apiutil "code.superseriousbusiness.org/gotosocial/internal/api/util"
	"code.superseriousbusiness.org/gotosocial/internal/gtserror"
)

// ActorGETHandler should be served at https://example.org/{users|relays}/:username.
// It returns the AP model of the given actor, requiring signed GET.
func (m *Module) ActorGETHandler(c *httputil.Context) {
	_, username, contentType, errWithCode := activitypub.ParseCommon(c)
	if errWithCode != nil {
		apiutil.ErrorHandler(c, m.templates, errWithCode)
		return
	}

	if contentType == apiutil.TextHTML {
		// Redirect to account web view.
		httputil.Redirect(c, http.StatusSeeOther, "/@"+username)
		return
	}

	resp, errWithCode := m.processor.Fedi().ActorGet(
		c,
		username,
	)
	if errWithCode != nil {
		apiutil.ErrorHandler(c, m.templates, errWithCode)
		return
	}

	httputil.JSONType(c, http.StatusOK, contentType, resp)
}

// InstanceActorGETHandler should be served at https://[hostname]/users/[hostname].
// It returns the AP model of the instance account, *without* requiring a signed GET.
func (m *Module) InstanceActorGETHandler(c *httputil.Context) {
	contentType, err := apiutil.NegotiateAccept(c, apiutil.ActivityPubHeaders...)
	if err != nil {
		apiutil.ErrorHandler(c, m.templates, gtserror.NewErrorNotAcceptable(err, err.Error()))
		return
	}

	resp, errWithCode := m.processor.Fedi().InstanceActorGet(c)
	if errWithCode != nil {
		apiutil.ErrorHandler(c, m.templates, errWithCode)
		return
	}

	httputil.JSONType(c, http.StatusOK, contentType, resp)
}
