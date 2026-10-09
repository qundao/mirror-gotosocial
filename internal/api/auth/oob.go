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

package auth

import (
	"errors"
	"net/http"

	"code.superseriousbusiness.org/gopkg/httputil"
	"code.superseriousbusiness.org/gotosocial/internal/middleware"
	"code.superseriousbusiness.org/gotosocial/internal/templates"
)

// OOBTokenGETHandler parses the OAuth code from the query
// params and serves a nice little HTML page showing the code.
func (m *Module) OOBTokenGETHandler(c *httputil.Context) {
	s := middleware.GetSession(c)

	oobToken := c.Query("code")
	if oobToken == "" {
		const errText = "no 'code' query value provided in callback redirect"
		m.clearSessionWithBadRequest(c, s, errors.New(errText), errText)
		return
	}

	// User must be fully authed.
	user := m.mustUserFromSession(c, s, sessionUserID)
	if user == nil {
		// Error already
		// written.
		return
	}

	scope := m.mustStringFromSession(c, s, sessionScope)
	if scope == "" {
		// Error already
		// written.
		return
	}

	// We're done with
	// the session now.
	m.mustClearSession(c, s)

	m.templates.RenderPage(c,
		http.StatusOK,
		templates.WebPage{
			Template: "oob.tmpl",
			Extra: map[string]any{
				"user":     user.Account.Username,
				"oobToken": oobToken,
				"scope":    scope,
			},
		},
	)
}
