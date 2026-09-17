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

package webfinger_test

import (
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"code.superseriousbusiness.org/gopkg/httputil"
	apiutil "code.superseriousbusiness.org/gotosocial/internal/api/util"
	"code.superseriousbusiness.org/gotosocial/internal/api/wellknown/webfinger"
	"code.superseriousbusiness.org/gotosocial/internal/config"
	"code.superseriousbusiness.org/gotosocial/internal/gtsmodel"
	"code.superseriousbusiness.org/gotosocial/testrig"
	"github.com/stretchr/testify/suite"
)

type WebfingerGetTestSuite struct {
	WebfingerStandardTestSuite
}

func (suite *WebfingerGetTestSuite) finger(requestPath string) string {
	// Set up the request.
	req := httptest.NewRequest(http.MethodGet, requestPath, nil)
	req.Header.Set("accept", "application/jrd+json")
	recorder := httptest.NewRecorder()
	c := httputil.ToContext(recorder, req)

	// Trigger the handler.
	suite.webfingerModule.WebfingerGETRequest(c)

	// Read the result + return it
	// as nicely indented JSON.
	result := recorder.Result()
	defer result.Body.Close()

	// Result should always use the
	// webfinger content-type.
	if ct := result.Header.Get("content-type"); ct != string(apiutil.AppJRDJSON) {
		suite.FailNow("", "expected content type %s, got %s", apiutil.AppJRDJSON, ct)
	}

	b, err := io.ReadAll(result.Body)
	if err != nil {
		suite.FailNow(err.Error())
	}

	return testrig.MustJSONStringFromBytes(b)
}

func (suite *WebfingerGetTestSuite) funkifyAccountDomain(host string, accountDomain string) *gtsmodel.Account {
	// Reset suite structs + config
	// to new host + account domain.
	config.SetHost(host)
	config.SetAccountDomain(accountDomain)
	suite.TearDownTest()
	suite.setupTest()

	// Generate a new account for the
	// tester, which uses the new host.
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		panic(err)
	}
	publicKey := &privateKey.PublicKey

	targetAccount := &gtsmodel.Account{
		ID:                    "01FG1K8EA7SYHEC7V6XKVNC4ZA",
		Username:              "new_account_domain_user",
		URI:                   "http://" + host + "/users/new_account_domain_user",
		URL:                   "http://" + host + "/@new_account_domain_user",
		InboxURI:              "http://" + host + "/users/new_account_domain_user/inbox",
		OutboxURI:             "http://" + host + "/users/new_account_domain_user/outbox",
		FollowingURI:          "http://" + host + "/users/new_account_domain_user/following",
		FollowersURI:          "http://" + host + "/users/new_account_domain_user/followers",
		FeaturedCollectionURI: "http://" + host + "/users/new_account_domain_user/collections/featured",
		ActorType:             gtsmodel.AccountActorTypePerson,
		PrivateKey:            privateKey,
		PublicKey:             publicKey,
		PublicKeyURI:          "http://" + host + "/users/new_account_domain_user/main-key",
	}

	if err := suite.db.PutAccount(suite.T().Context(), targetAccount); err != nil {
		suite.FailNow(err.Error())
	}

	if err := suite.db.PutAccountSettings(suite.T().Context(), &gtsmodel.AccountSettings{AccountID: targetAccount.ID}); err != nil {
		suite.FailNow(err.Error())
	}

	return targetAccount
}

func (suite *WebfingerGetTestSuite) TestFingerUser() {
	targetAccount := suite.testAccounts["local_account_1"]
	requestPath := fmt.Sprintf("/%s?resource=acct:%s@%s", webfinger.WebfingerBasePath, targetAccount.Username, config.GetHost())

	resp := suite.finger(requestPath)
	suite.Equal(`{
  "aliases": [
    "http://localhost:8080/users/the_mighty_zork",
    "http://localhost:8080/@the_mighty_zork"
  ],
  "links": [
    {
      "href": "http://localhost:8080/@the_mighty_zork",
      "rel": "http://webfinger.net/rel/profile-page",
      "type": "text/html"
    },
    {
      "href": "http://localhost:8080/users/the_mighty_zork",
      "rel": "self",
      "type": "application/activity+json"
    }
  ],
  "subject": "acct:the_mighty_zork@localhost:8080"
}`, resp)
}

func (suite *WebfingerGetTestSuite) TestFingerUserActorURI() {
	targetAccount := suite.testAccounts["local_account_1"]
	host := config.GetHost()

	tests := []struct {
		resource string
	}{
		{resource: fmt.Sprintf("https://%s/@%s", host, targetAccount.Username)},
		{resource: fmt.Sprintf("https://%s/users/%s", host, targetAccount.Username)},
	}

	for _, tt := range tests {
		tt := tt
		suite.Run(tt.resource, func() {
			requestPath := fmt.Sprintf("/%s?resource=%s", webfinger.WebfingerBasePath, tt.resource)
			resp := suite.finger(requestPath)
			suite.Equal(`{
  "aliases": [
    "http://localhost:8080/users/the_mighty_zork",
    "http://localhost:8080/@the_mighty_zork"
  ],
  "links": [
    {
      "href": "http://localhost:8080/@the_mighty_zork",
      "rel": "http://webfinger.net/rel/profile-page",
      "type": "text/html"
    },
    {
      "href": "http://localhost:8080/users/the_mighty_zork",
      "rel": "self",
      "type": "application/activity+json"
    }
  ],
  "subject": "acct:the_mighty_zork@localhost:8080"
}`, resp)
		})
	}
}

func (suite *WebfingerGetTestSuite) TestFingerUserWithDifferentAccountDomainByHost() {
	targetAccount := suite.funkifyAccountDomain("gts.example.org", "example.org")
	requestPath := fmt.Sprintf("/%s?resource=acct:%s@%s", webfinger.WebfingerBasePath, targetAccount.Username, config.GetHost())

	resp := suite.finger(requestPath)
	suite.Equal(`{
  "aliases": [
    "http://gts.example.org/users/new_account_domain_user",
    "http://gts.example.org/@new_account_domain_user"
  ],
  "links": [
    {
      "href": "http://gts.example.org/@new_account_domain_user",
      "rel": "http://webfinger.net/rel/profile-page",
      "type": "text/html"
    },
    {
      "href": "http://gts.example.org/users/new_account_domain_user",
      "rel": "self",
      "type": "application/activity+json"
    }
  ],
  "subject": "acct:new_account_domain_user@example.org"
}`, resp)
}

func (suite *WebfingerGetTestSuite) TestFingerUserWithDifferentAccountDomainByAccountDomain() {
	targetAccount := suite.funkifyAccountDomain("gts.example.org", "example.org")
	requestPath := fmt.Sprintf("/%s?resource=acct:%s@%s", webfinger.WebfingerBasePath, targetAccount.Username, config.GetAccountDomain())

	resp := suite.finger(requestPath)
	suite.Equal(`{
  "aliases": [
    "http://gts.example.org/users/new_account_domain_user",
    "http://gts.example.org/@new_account_domain_user"
  ],
  "links": [
    {
      "href": "http://gts.example.org/@new_account_domain_user",
      "rel": "http://webfinger.net/rel/profile-page",
      "type": "text/html"
    },
    {
      "href": "http://gts.example.org/users/new_account_domain_user",
      "rel": "self",
      "type": "application/activity+json"
    }
  ],
  "subject": "acct:new_account_domain_user@example.org"
}`, resp)
}

func (suite *WebfingerGetTestSuite) TestFingerUserWithoutAcct() {
	// Leave out the 'acct:' part in the request path;
	// the handler should be generous + still work OK.
	targetAccount := suite.testAccounts["local_account_1"]
	requestPath := fmt.Sprintf("/%s?resource=%s@%s", webfinger.WebfingerBasePath, targetAccount.Username, config.GetHost())

	resp := suite.finger(requestPath)
	suite.Equal(`{
  "aliases": [
    "http://localhost:8080/users/the_mighty_zork",
    "http://localhost:8080/@the_mighty_zork"
  ],
  "links": [
    {
      "href": "http://localhost:8080/@the_mighty_zork",
      "rel": "http://webfinger.net/rel/profile-page",
      "type": "text/html"
    },
    {
      "href": "http://localhost:8080/users/the_mighty_zork",
      "rel": "self",
      "type": "application/activity+json"
    }
  ],
  "subject": "acct:the_mighty_zork@localhost:8080"
}`, resp)
}

func TestWebfingerGetTestSuite(t *testing.T) {
	suite.Run(t, new(WebfingerGetTestSuite))
}
