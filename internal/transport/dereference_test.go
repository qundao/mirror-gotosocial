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

package transport_test

import (
	"io"
	"net/http"
	"testing"

	apiutil "code.superseriousbusiness.org/gotosocial/internal/api/util"
	"code.superseriousbusiness.org/gotosocial/testrig"
	"github.com/stretchr/testify/suite"
)

type DereferenceTestSuite struct {
	TransportTestSuite
}

func (suite *DereferenceTestSuite) TestDerefLocalUser() {
	iri := testrig.URLMustParse(suite.testAccounts["local_account_1"].URI)

	resp, err := suite.transport.Dereference(suite.T().Context(), iri)
	if err != nil {
		suite.FailNow(err.Error())
	}
	defer resp.Body.Close()

	suite.Equal(http.StatusOK, resp.StatusCode)
	suite.EqualValues(2146, resp.ContentLength)
	suite.Equal("2146", resp.Header.Get("Content-Length"))
	suite.Equal(apiutil.AppActivityLDJSON, resp.Header.Get("Content-Type"))

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		suite.FailNow(err.Error())
	}

	out := testrig.MustJSONStringFromBytes(b)
	suite.Equal(`{
  "@context": [
    "https://gotosocial.org/ns",
    "https://w3id.org/security/v1",
    "https://www.w3.org/ns/activitystreams",
    {
      "discoverable": "toot:discoverable",
      "featured": {
        "@id": "toot:featured",
        "@type": "@id"
      },
      "indexable": "toot:indexable",
      "manuallyApprovesFollowers": "as:manuallyApprovesFollowers",
      "toot": "http://joinmastodon.org/ns#"
    }
  ],
  "discoverable": true,
  "featured": "http://localhost:8080/users/the_mighty_zork/collections/featured",
  "followers": "http://localhost:8080/users/the_mighty_zork/followers",
  "following": "http://localhost:8080/users/the_mighty_zork/following",
  "hidesCcPublicFromUnauthedWeb": false,
  "hidesToPublicFromUnauthedWeb": false,
  "icon": {
    "mediaType": "image/jpeg",
    "name": "a green goblin looking nasty",
    "type": "Image",
    "url": "http://localhost:8080/fileserver/01F8MH1H7YV1Z7D2C8K2730QBF/avatar/original/01F8MH58A357CV5K7R7TJMSH6S.jpg"
  },
  "id": "http://localhost:8080/users/the_mighty_zork",
  "image": {
    "mediaType": "image/jpeg",
    "name": "A very old-school screenshot of the original team fortress mod for quake",
    "type": "Image",
    "url": "http://localhost:8080/fileserver/01F8MH1H7YV1Z7D2C8K2730QBF/header/original/01PFPMWK2FF0D9WMHEJHR07C3Q.jpg"
  },
  "inbox": "http://localhost:8080/users/the_mighty_zork/inbox",
  "indexable": true,
  "manuallyApprovesFollowers": false,
  "name": "original zork (he/they)",
  "outbox": "http://localhost:8080/users/the_mighty_zork/outbox",
  "preferredUsername": "the_mighty_zork",
  "publicKey": {
    "id": "http://localhost:8080/users/the_mighty_zork/main-key",
    "owner": "http://localhost:8080/users/the_mighty_zork",
    "publicKeyPem": "-----BEGIN PUBLIC KEY-----\nMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA0iLMdE7FAnkh1w2D81Yq\nkcAf8GYQZI+X9DK3PkNMqIFJok3Lm2EHH/8Y/Szr7QKyfvP66K4tI/GhV9e42arR\ncgwyDcoSAk/ouW5NWSj7f7h/x+aBaY7bC/mwqDhfBeAFHhCxMR9HTt7sJfL/Xz2W\ndwGRBo+lAekbtdyzje3yh7fiU+rPYzbVCFKR1A4NhmWL/YCxRgw5vR/dWHq75fMh\nelVmyvu6XFcoZc+cKh0f6jVIslF4Yonvr3oiXPYqQlO0a4jRLnobxddnd60SDiv8\nEbBQBuC8bnyUvEobvFSazgZSs7Ln6ow2bZ2W/Eq02NBIyyabJTH+u80Qw9ZBA6au\nIwIDAQAB\n-----END PUBLIC KEY-----\n"
  },
  "published": "2022-05-20T11:09:18Z",
  "summary": "<p>hey yo this is my profile!</p>",
  "type": "Person",
  "url": "http://localhost:8080/@the_mighty_zork"
}`, out)
}

func (suite *DereferenceTestSuite) TestDerefLocalStatus() {
	iri := testrig.URLMustParse(suite.testStatuses["local_account_1_status_1"].URI)

	resp, err := suite.transport.Dereference(suite.T().Context(), iri)
	if err != nil {
		suite.FailNow(err.Error())
	}
	defer resp.Body.Close()

	suite.Equal(http.StatusOK, resp.StatusCode)
	suite.EqualValues(1525, resp.ContentLength)
	suite.Equal("1525", resp.Header.Get("Content-Length"))
	suite.Equal(apiutil.AppActivityLDJSON, resp.Header.Get("Content-Type"))

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		suite.FailNow(err.Error())
	}

	out := testrig.MustJSONStringFromBytes(b)
	suite.Equal(`{
  "@context": [
    "https://gotosocial.org/ns",
    "https://www.w3.org/ns/activitystreams",
    {
      "sensitive": "as:sensitive"
    }
  ],
  "attributedTo": "http://localhost:8080/users/the_mighty_zork",
  "cc": "http://localhost:8080/users/the_mighty_zork/followers",
  "content": "<p>hello everyone!</p>",
  "contentMap": {
    "en": "<p>hello everyone!</p>"
  },
  "id": "http://localhost:8080/users/the_mighty_zork/statuses/01F8MHAMCHF6Y650WCRSCP4WMY",
  "interactionPolicy": {
    "canAnnounce": {
      "automaticApproval": [
        "https://www.w3.org/ns/activitystreams#Public"
      ]
    },
    "canLike": {
      "automaticApproval": [
        "https://www.w3.org/ns/activitystreams#Public"
      ]
    },
    "canQuote": {
      "automaticApproval": [
        "http://localhost:8080/users/the_mighty_zork"
      ]
    },
    "canReply": {
      "automaticApproval": [
        "https://www.w3.org/ns/activitystreams#Public"
      ]
    }
  },
  "published": "2021-10-20T10:40:37Z",
  "replies": {
    "first": {
      "id": "http://localhost:8080/users/the_mighty_zork/statuses/01F8MHAMCHF6Y650WCRSCP4WMY/replies?page=true",
      "next": "http://localhost:8080/users/the_mighty_zork/statuses/01F8MHAMCHF6Y650WCRSCP4WMY/replies?page=true&only_other_accounts=false",
      "partOf": "http://localhost:8080/users/the_mighty_zork/statuses/01F8MHAMCHF6Y650WCRSCP4WMY/replies",
      "type": "CollectionPage"
    },
    "id": "http://localhost:8080/users/the_mighty_zork/statuses/01F8MHAMCHF6Y650WCRSCP4WMY/replies",
    "type": "Collection"
  },
  "sensitive": true,
  "summary": "introduction post",
  "to": "https://www.w3.org/ns/activitystreams#Public",
  "type": "Note",
  "url": "http://localhost:8080/@the_mighty_zork/statuses/01F8MHAMCHF6Y650WCRSCP4WMY"
}`, out)
}

func (suite *DereferenceTestSuite) TestDerefLocalFollowers() {
	iri := testrig.URLMustParse(suite.testAccounts["local_account_1"].FollowersURI)

	resp, err := suite.transport.Dereference(suite.T().Context(), iri)
	if err != nil {
		suite.FailNow(err.Error())
	}
	defer resp.Body.Close()

	suite.Equal(http.StatusOK, resp.StatusCode)
	suite.EqualValues(161, resp.ContentLength)
	suite.Equal("161", resp.Header.Get("Content-Length"))
	suite.Equal(apiutil.AppActivityLDJSON, resp.Header.Get("Content-Type"))

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		suite.FailNow(err.Error())
	}

	out := testrig.MustJSONStringFromBytes(b)
	suite.Equal(`{
  "@context": "https://www.w3.org/ns/activitystreams",
  "items": [
    "http://localhost:8080/users/1happyturtle",
    "http://localhost:8080/users/admin"
  ],
  "type": "Collection"
}`, out)
}

func (suite *DereferenceTestSuite) TestDerefLocalFollowing() {
	iri := testrig.URLMustParse(suite.testAccounts["local_account_1"].FollowingURI)

	resp, err := suite.transport.Dereference(suite.T().Context(), iri)
	if err != nil {
		suite.FailNow(err.Error())
	}
	defer resp.Body.Close()

	suite.Equal(http.StatusOK, resp.StatusCode)
	suite.EqualValues(161, resp.ContentLength)
	suite.Equal("161", resp.Header.Get("Content-Length"))
	suite.Equal(apiutil.AppActivityLDJSON, resp.Header.Get("Content-Type"))

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		suite.FailNow(err.Error())
	}

	out := testrig.MustJSONStringFromBytes(b)
	suite.Equal(`{
  "@context": "https://www.w3.org/ns/activitystreams",
  "items": [
    "http://localhost:8080/users/admin",
    "http://localhost:8080/users/1happyturtle"
  ],
  "type": "Collection"
}`, out)
}

func TestDereferenceTestSuite(t *testing.T) {
	suite.Run(t, new(DereferenceTestSuite))
}
