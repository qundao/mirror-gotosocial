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

	"code.superseriousbusiness.org/gopkg/httputil"
	apiutil "code.superseriousbusiness.org/gotosocial/internal/api/util"
	"code.superseriousbusiness.org/gotosocial/internal/gtserror"
	"code.superseriousbusiness.org/gotosocial/internal/gtsmodel"
	"code.superseriousbusiness.org/gotosocial/internal/oauth"
	"github.com/gorilla/sessions"
)

func (m *Module) mustClearSession(c *httputil.Context, s *sessions.Session) {
	clear(s.Values)
	m.mustSaveSession(c, s)
}

func (m *Module) mustSaveSession(c *httputil.Context, s *sessions.Session) {
	if err := s.Save(c.R, &c.W); err != nil {
		panic(err)
	}
}

// mustUserFromSession returns a *gtsmodel.User by checking
// the session for a user id stored under the given key,
// and fetching the user from the database.
//
// On failure, the function clears session state, writes an internal
// error to the response writer, and returns nil. Callers should always
// return immediately if receiving nil back from this function!
func (m *Module) mustUserFromSession(
	c *httputil.Context,
	s *sessions.Session,
	key string,
) *gtsmodel.User {
	userID, _ := s.Values[key].(string)
	if userID == "" {
		const safe = "neither userid nor userid_awaiting_2fa keys found in session"
		m.clearSessionWithInternalError(c, s, errors.New(safe), safe, oauth.HelpfulAdvice)
		return nil
	}

	user, err := m.state.DB.GetUserByID(c, userID)
	if err != nil {
		safe := "db error getting user " + userID
		m.clearSessionWithInternalError(c, s, err, safe, oauth.HelpfulAdvice)
		return nil
	}

	return user
}

// mustAppFromSession returns a *gtsmodel.Application by checking the
// session for an application keyid and fetching the app from the database.
//
// On failure, the function clears session state, writes an internal
// error to the response writer, and returns nil. Callers should always
// return immediately if receiving nil back from this function!
func (m *Module) mustAppFromSession(
	c *httputil.Context,
	s *sessions.Session,
) *gtsmodel.Application {
	clientID, ok := s.Values[sessionClientID].(string)
	if !ok {
		const safe = "key client_id not found in session"
		m.clearSessionWithInternalError(c, s, errors.New(safe), safe, oauth.HelpfulAdvice)
		return nil
	}

	app, err := m.state.DB.GetApplicationByClientID(c, clientID)
	if err != nil {
		safe := "db error getting app for clientID " + clientID
		m.clearSessionWithInternalError(c, s, err, safe, oauth.HelpfulAdvice)
		return nil
	}

	return app
}

// mustStringFromSession returns the string value
// corresponding to the given session key, if any is set.
//
// On failure (nothing set), the function clears session
// state, writes an internal error to the response writer,
// and returns nil. Callers should always return immediately
// if receiving nil back from this function!
func (m *Module) mustStringFromSession(
	c *httputil.Context,
	s *sessions.Session,
	key string,
) string {
	v, ok := s.Values[key].(string)
	if !ok {
		safe := "key " + key + " not found in session"
		m.clearSessionWithInternalError(c, s, errors.New(safe), safe, oauth.HelpfulAdvice)
		return ""
	}
	return v
}

func (m *Module) clearSessionWithInternalError(
	c *httputil.Context,
	s *sessions.Session,
	err error,
	helpText ...string,
) {
	m.mustClearSession(c, s)
	errWithCode := gtserror.NewErrorInternalError(err, helpText...)
	apiutil.ErrorHandler(c, m.templates, errWithCode)
}

func (m *Module) clearSessionWithBadRequest(
	c *httputil.Context,
	s *sessions.Session,
	err error,
	helpText ...string,
) {
	m.mustClearSession(c, s)
	errWithCode := gtserror.NewErrorBadRequest(err, helpText...)
	apiutil.ErrorHandler(c, m.templates, errWithCode)
}
