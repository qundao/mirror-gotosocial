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
)

// StatusGETHandler serves the target status as an activitystreams NOTE so that other AP servers can parse it.
func (m *Module) StatusGETHandler(c *httputil.Context) {
	_, username, statusID, contentType, errWithCode := activitypub.ParseCommonWithID(c)
	if errWithCode != nil {
		apiutil.ErrorHandler(c, m.templates, errWithCode)
		return
	}

	if contentType == apiutil.TextHTML {
		// Redirect to status web view.
		httputil.Redirect(c, http.StatusSeeOther, "/@"+username+"/statuses/"+statusID)
		return
	}

	resp, errWithCode := m.processor.Fedi().StatusGet(c, username, statusID)
	if errWithCode != nil {
		apiutil.ErrorHandler(c, m.templates, errWithCode)
		return
	}

	httputil.JSONType(c, http.StatusOK, contentType, resp)
}
