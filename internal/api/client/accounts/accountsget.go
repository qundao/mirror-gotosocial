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

package accounts

import (
	"errors"
	"net/http"
	"slices"

	"code.superseriousbusiness.org/gopkg/httputil"
	apimodel "code.superseriousbusiness.org/gotosocial/internal/api/model"
	apiutil "code.superseriousbusiness.org/gotosocial/internal/api/util"
	"code.superseriousbusiness.org/gotosocial/internal/gtserror"
)

// AccountsGETHandler swagger:operation GET /api/v1/accounts accountsGet
//
// Get information about multiple accounts with the given IDs.
//
// Accounts that do not exist are omitted from the response.
// If no IDs are given, an empty array is returned.
//
//	---
//	tags:
//	- accounts
//
//	produces:
//	- application/json
//
//	parameters:
//	-
//		name: id[]
//		type: array
//		items:
//			type: string
//		description: Target account IDs (max 40).
//		in: query
//		collectionFormat: multi
//		required: true
//
//	security:
//	- OAuth2 Bearer:
//		- read:accounts
//
//	responses:
//		'200':
//			description: "The requested accounts."
//			schema:
//				type: array
//				items:
//					"$ref": "#/definitions/account"
//		'400':
//			schema:
//				"$ref": "#/definitions/error"
//			description: bad request
//		'401':
//			schema:
//				"$ref": "#/definitions/error"
//			description: unauthorized
//		'403':
//			schema:
//				"$ref": "#/definitions/error"
//			description: forbidden
//		'406':
//			schema:
//				"$ref": "#/definitions/error"
//			description: not acceptable
//		'422':
//			schema:
//				"$ref": "#/definitions/error"
//			description: unprocessable entity
//		'500':
//			schema:
//				"$ref": "#/definitions/error"
//			description: internal server error
func (m *Module) AccountsGETHandler(c *httputil.Context) {
	authed, errWithCode := apiutil.TokenAuth(c, apiutil.AuthRequirements{
		Token:   true,
		App:     true,
		User:    true,
		Account: true,
		Scope:   []apiutil.Scope{apiutil.ScopeReadAccounts},
	})
	if errWithCode != nil {
		apiutil.ErrorHandler(c, m.templates, errWithCode)
		return
	}

	if _, errWithCode := apiutil.NegotiateAccept(c, apiutil.JSONAcceptHeaders...); errWithCode != nil {
		apiutil.ErrorHandler(c, m.templates, errWithCode)
		return
	}

	ids := c.QueryArray("id[]")
	if len(ids) == 0 {
		ids = c.QueryArray("id")
	}

	ids = slices.DeleteFunc(ids, func(id string) bool {
		return id == ""
	})

	if len(ids) == 0 {
		// Return an empty array if no IDs given as Mastodon.
		httputil.JSON(c, http.StatusOK, []*apimodel.Account{})
		return
	}

	if len(ids) > 40 {
		// Same limit as Mastodon.
		const text = "too many account ids specified"
		errWithCode := gtserror.NewErrorUnprocessableEntity(errors.New(text), text)
		apiutil.ErrorHandler(c, m.templates, errWithCode)
		return
	}

	apiAccounts, errWithCode := m.processor.Account().GetMultiple(c, authed.Account, ids)
	if errWithCode != nil {
		apiutil.ErrorHandler(c, m.templates, errWithCode)
		return
	}

	httputil.JSON(c, http.StatusOK, apiAccounts)
}
