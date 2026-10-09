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
	"context"
	"errors"
	"fmt"
	"net/http"
	"slices"
	"strings"

	"code.superseriousbusiness.org/gopkg/httputil"
	"code.superseriousbusiness.org/gopkg/httputil/binding"
	apiutil "code.superseriousbusiness.org/gotosocial/internal/api/util"
	"code.superseriousbusiness.org/gotosocial/internal/config"
	"code.superseriousbusiness.org/gotosocial/internal/gtserror"
	"code.superseriousbusiness.org/gotosocial/internal/gtsmodel"
	"code.superseriousbusiness.org/gotosocial/internal/middleware"
	"code.superseriousbusiness.org/gotosocial/internal/oauth"
	"code.superseriousbusiness.org/gotosocial/internal/templates"
	"codeberg.org/gruf/go-byteutil"
	"github.com/pquerna/otp/totp"
	"golang.org/x/crypto/bcrypt"
)

// SignInGETHandler should be served at
// GET https://example.org/auth/sign_in.
//
// The idea is to present a friendly sign-in
// page to the user, where they can enter their
// username and password.
//
// When submitted, the form will POST to the sign-
// in page, which will be handled by SignInPOSTHandler.
//
// If an idp provider is set, then the user will
// be redirected to that to do their sign in.
func (m *Module) SignInGETHandler(c *httputil.Context) {
	if _, errWithCode := apiutil.NegotiateAccept(c, apiutil.HTMLAcceptHeaders...); errWithCode != nil {
		apiutil.ErrorHandler(c, m.templates, errWithCode)
		return
	}

	if config.GetOIDCEnabled() {
		// IDP provider is in use, so redirect to it
		// instead of serving our own sign in page.
		//
		// We need the internal state to know where
		// to redirect to.
		internalState := m.mustStringFromSession(
			c,
			middleware.GetSession(c),
			sessionInternalState,
		)
		if internalState == "" {
			// Error already
			// written.
			return
		}

		httputil.Redirect(c, http.StatusSeeOther, m.idp.AuthCodeURL(internalState))
		return
	}

	m.templates.RenderPage(c, http.StatusOK,
		templates.WebPage{Template: "sign-in.tmpl"})
}

// SignInPOSTHandler should be served at
// POST https://example.org/auth/sign_in.
//
// The handler will check the submitted credentials,
// then redirect either to the 2fa form, or straight
// to the authorize page served at /oauth/authorize.
func (m *Module) SignInPOSTHandler(c *httputil.Context) {
	s := middleware.GetSession(c)

	// Parse email
	// + password.
	form := &struct {
		Email    string `form:"username" binding:"required"`
		Password string `form:"password" binding:"required"`
	}{}
	if err := binding.ShouldBind(c, form, int64(config.GetHTTPServerMaxMultipartMemory())); err != nil { // nolint
		m.clearSessionWithBadRequest(c, s, err, oauth.HelpfulAdvice)
		return
	}

	user, errWithCode := m.validatePassword(
		c,
		form.Email,
		form.Password,
	)
	if errWithCode != nil {
		// Don't clear session here yet, so the user
		// can just press back and try again if they
		// accidentally gave the wrong password, without
		// having to do the whole sign in flow again!
		apiutil.ErrorHandler(c, m.templates, errWithCode)
		return
	}

	// Whether or not 2fa is enabled, we want
	// to save the session when we're done here.
	defer m.mustSaveSession(c, s)

	if user.TwoFactorEnabled() {
		// If this user has 2FA enabled, redirect
		// to the 2FA page and have them submit
		// a code from their authenticator app.
		s.Values[sessionUserIDAwaiting2FA] = user.ID
		httputil.Redirect(c, http.StatusFound, "/auth"+Auth2FAPath)
		return
	}

	// If the user doesn't have 2fa enabled,
	// redirect straight to the OAuth authorize page.
	s.Values[sessionUserID] = user.ID
	httputil.Redirect(c, http.StatusFound, "/oauth"+OauthAuthorizePath)
}

// validatePassword takes an email address and a password.
// The func authenticates the password against the one for
// that email address stored in the database.
//
// If OK, it returns the user, so that it can be used in
// further OAuth flows to generate a token etc.
func (m *Module) validatePassword(
	ctx context.Context,
	email string,
	password string,
) (
	*gtsmodel.User,
	gtserror.WithCode,
) {
	if email == "" || password == "" {
		err := errors.New("email or password was not provided")
		return incorrectPassword(err)
	}

	user, err := m.state.DB.GetUserByEmailAddress(ctx, email)
	if err != nil {
		err := fmt.Errorf("user %s was not retrievable from db during oauth authorization attempt: %s", email, err)
		return incorrectPassword(err)
	}

	if user.EncryptedPassword == "" {
		err := fmt.Errorf("encrypted password for user %s was empty for some reason", user.Email)
		return incorrectPassword(err)
	}

	if err := bcrypt.CompareHashAndPassword(
		byteutil.S2B(user.EncryptedPassword),
		byteutil.S2B(password),
	); err != nil {
		err := fmt.Errorf("password hash didn't match for user %s during sign in attempt: %s", user.Email, err)
		return incorrectPassword(err)
	}

	return user, nil
}

// incorrectPassword wraps the given error in a gtserror.WithCode, and returns
// only a generic 'safe' error message to the user, to not give any info away.
func incorrectPassword(err error) (*gtsmodel.User, gtserror.WithCode) {
	const errText = "password/email combination was incorrect"
	return nil, gtserror.NewErrorUnauthorized(err, errText, oauth.HelpfulAdvice)
}

// TwoFactorCodeGETHandler should be served at
// GET https://example.org/auth/2fa.
//
// The 2fa template displays a simple form asking the
// user to input a code from their authenticator app.
func (m *Module) TwoFactorCodeGETHandler(c *httputil.Context) {
	s := middleware.GetSession(c)

	// User must be awaiting 2fa.
	user := m.mustUserFromSession(c, s, sessionUserIDAwaiting2FA)
	if user == nil {
		// Error already
		// written.
		return
	}

	m.templates.RenderPage(c, http.StatusOK,
		templates.WebPage{
			Template: "2fa.tmpl",
			Extra: map[string]any{
				"user": user.Account.Username,
			},
		},
	)
}

// TwoFactorCodePOSTHandler should be served at
// POST https://example.org/auth/2fa.
//
// The idea is to handle a submitted 2fa code, validate it,
// and if valid redirect to the /oauth/authorize page that
// the user would get to if they didn't have 2fa enabled.
func (m *Module) TwoFactorCodePOSTHandler(c *httputil.Context) {
	s := middleware.GetSession(c)

	// User must be awaiting 2fa.
	user := m.mustUserFromSession(c, s, sessionUserIDAwaiting2FA)
	if user == nil {
		// Error already
		// written.
		return
	}

	// Parse 2fa code.
	form := &struct {
		Code string `form:"code" binding:"required"`
	}{}
	if err := binding.ShouldBind(c, form, int64(config.GetHTTPServerMaxMultipartMemory())); err != nil { // nolint
		m.clearSessionWithBadRequest(c, s, err, oauth.HelpfulAdvice)
		return
	}

	valid, err := m.validate2FACode(c, user, form.Code)
	if err != nil {
		m.clearSessionWithInternalError(c, s, err, oauth.HelpfulAdvice)
		return
	}

	if !valid {
		// Don't clear session here yet, so the user
		// can just press back and try again if they
		// accidentally gave the wrong code, without
		// having to do the whole sign in flow again!
		const errText = "2fa code invalid or timed out, press back and try again; " +
			"if issues persist, pester your instance admin to check the server clock"
		errWithCode := gtserror.NewErrorBadRequest(errors.New(errText), errText)
		apiutil.ErrorHandler(c, m.templates, errWithCode)
		return
	}

	// Code looks good! Redirect
	// to the OAuth authorize page.
	s.Values[sessionUserID] = user.ID
	m.mustSaveSession(c, s)
	httputil.Redirect(c, http.StatusFound, "/oauth"+OauthAuthorizePath)
}

func (m *Module) validate2FACode(c *httputil.Context, user *gtsmodel.User, code string) (bool, error) {
	code = strings.TrimSpace(code)
	if len(code) <= 6 {
		// This is a normal authenticator
		// app code, just try to validate it.
		return totp.Validate(code, user.TwoFactorSecret), nil
	}

	// This is a one-time recovery code.
	// Check against the user's stored codes.
	for i := 0; i < len(user.TwoFactorBackups); i++ {
		err := bcrypt.CompareHashAndPassword(
			byteutil.S2B(user.TwoFactorBackups[i]),
			byteutil.S2B(code),
		)
		if err != nil {
			// Doesn't match,
			// try next.
			continue
		}

		// We have a match.
		// Remove this one-time code from the user's backups.
		user.TwoFactorBackups = slices.Delete(user.TwoFactorBackups, i, i+1)
		if err := m.state.DB.UpdateUser(
			c,
			user,
			"two_factor_backups",
		); err != nil {
			return false, err
		}

		// So valid bestie!
		return true, nil
	}

	// Not a valid one-time
	// recovery code.
	return false, nil
}
