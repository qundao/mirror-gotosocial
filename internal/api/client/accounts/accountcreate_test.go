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

package accounts_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"code.superseriousbusiness.org/gotosocial/internal/api/client/accounts"
	"code.superseriousbusiness.org/gotosocial/internal/config"
	"code.superseriousbusiness.org/gotosocial/internal/language"
	"code.superseriousbusiness.org/gotosocial/testrig"
	"github.com/stretchr/testify/suite"
)

type AccountCreateTestSuite struct {
	AccountStandardTestSuite
}

func (suite *AccountCreateTestSuite) TestAccountCreatePOSTHandlerLocale() {
	for _, testStruct := range []struct {
		username       string
		password       string
		email          string
		agreement      string
		reason         string
		locale         string
		expectedLocale string
		before         func()
	}{
		{
			username:       "brand_new_stinkin_account",
			password:       "passwordpasswordpasswordpasswordpassword",
			email:          "someone@example.org",
			agreement:      "true",
			reason:         "i tell you what i want what i really really want so tell me what you want what you really really want",
			locale:         "be",
			expectedLocale: "be",
		},
		{
			username:  "cacapoopoo",
			password:  "passwordpasswordpasswordpasswordpassword",
			email:     "someone_else@example.org",
			agreement: "true",
			reason:    "i tell you what i want what i really really want so tell me what you want what you really really want",
			locale:    "",
			// Locale falls back to "en"
			// when nothing else is set.
			expectedLocale: "en",
			before: func() {
				// Clear instance languages.
				config.SetInstanceLanguages(language.Languages{})
			},
		},
		{
			username:  "peepeeintheweewee",
			password:  "passwordpasswordpasswordpasswordpassword",
			email:     "another_someone_else@example.org",
			agreement: "true",
			reason:    "i tell you what i want what i really really want so tell me what you want what you really really want",
			locale:    "",
			// Locale falls back to first
			// instance language when not provided.
			expectedLocale: "nl",
			before: func() {
				// Reset instance languages back to test defaults.
				config.SetInstanceLanguages(language.Languages{{TagStr: "nl"}, {TagStr: "en-gb"}})
			},
		},
	} {
		if testStruct.before != nil {
			testStruct.before()
		}

		// Set up the request
		requestBody, w, err := testrig.CreateMultipartFormData(
			nil,
			map[string][]string{
				"username":  {testStruct.username},
				"password":  {testStruct.password},
				"email":     {testStruct.email},
				"agreement": {testStruct.agreement},
				"reason":    {testStruct.reason},
				"locale":    {testStruct.locale},
			})
		if err != nil {
			panic(err)
		}
		bodyBytes := requestBody.Bytes()
		recorder := httptest.NewRecorder()
		ctx := suite.newContext(recorder, http.MethodPost, bodyBytes, accounts.BasePath, w.FormDataContentType())

		// Call the handler
		suite.accountsModule.AccountCreatePOSTHandler(ctx)

		// We should have status OK
		// because our request was valid.
		if recorder.Code != http.StatusOK {
			b, err := io.ReadAll(recorder.Result().Body)
			if err != nil {
				suite.FailNow(err.Error())
			}
			suite.FailNow("", "expected code 202, got %d: %s", recorder.Code, string(b))
		}

		// There should be an account in the database now.
		acct, err := suite.state.DB.GetAccountByUsernameDomain(ctx, testStruct.username, "")
		if err != nil {
			suite.FailNow(err.Error())
		}

		// There should be a user in the database now.
		user, err := suite.state.DB.GetUserByAccountID(ctx, acct.ID)
		if err != nil {
			suite.FailNow(err.Error())
		}

		// Account preferred language
		// should be what we set above.
		suite.Equal(testStruct.expectedLocale, acct.Settings.Language)

		// User locale should be what we set above.
		suite.Equal(testStruct.expectedLocale, user.Locale)
	}
}

func TestAccountCreateTestSuite(t *testing.T) {
	suite.Run(t, new(AccountCreateTestSuite))
}
