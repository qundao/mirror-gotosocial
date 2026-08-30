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

package admin_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"code.superseriousbusiness.org/gotosocial/internal/api/client/admin"
	"code.superseriousbusiness.org/gotosocial/testrig"
	"github.com/stretchr/testify/suite"
)

type AccountsGetTestSuite struct {
	AdminStandardTestSuite
}

func (suite *AccountsGetTestSuite) TestAccountsGetFromTop() {
	recorder := httptest.NewRecorder()

	path := admin.AccountsV2Path
	ctx := suite.newContext(recorder, http.MethodGet, nil, path, "application/json")

	suite.adminModule.AccountsGETV2Handler(ctx)
	suite.Equal(http.StatusOK, recorder.Code)

	b, err := io.ReadAll(recorder.Body)
	if err != nil {
		suite.FailNow(err.Error())
	}

	link := recorder.Header().Get("Link")
	suite.Equal(`<http://localhost:8080/api/v2/admin/accounts?limit=50&max_id=xn--xample-ova.org%2F%40%C3%BCser>; rel="next", <http://localhost:8080/api/v2/admin/accounts?limit=50&min_id=%2F%401happyturtle>; rel="prev"`, link)

	out := testrig.MustJSONStringFromBytes(b)
	suite.Equal(`[
  {
    "account": {
      "acct": "1happyturtle",
      "avatar": "",
      "avatar_static": "",
      "bot": false,
      "created_at": "2022-06-04T13:12:00.000Z",
      "discoverable": false,
      "display_name": "happy little turtle :3",
      "emojis": [],
      "fields": [
        {
          "name": "should you follow me?",
          "value": "maybe!",
          "verified_at": null
        },
        {
          "name": "age",
          "value": "120",
          "verified_at": null
        }
      ],
      "followers_count": 1,
      "following_count": 1,
      "group": false,
      "header": "http://localhost:8080/assets/default_header.webp",
      "header_description": "Flat gray background (default header).",
      "header_static": "http://localhost:8080/assets/default_header.webp",
      "hide_collections": true,
      "id": "01F8MH5NBDF2MV7CTC4Q5128HF",
      "indexable": false,
      "last_status_at": "2026-01-01",
      "locked": true,
      "noindex": true,
      "note": "<p>i post about things that concern me</p>",
      "statuses_count": 10,
      "url": "http://localhost:8080/@1happyturtle",
      "username": "1happyturtle"
    },
    "approved": true,
    "confirmed": true,
    "created_at": "2022-06-04T13:12:00.000Z",
    "created_by_application_id": "01F8MGY43H3N2C8EWPR2FPYEXG",
    "disabled": false,
    "domain": null,
    "email": "tortle.dude@example.org",
    "id": "01F8MH5NBDF2MV7CTC4Q5128HF",
    "invite_request": null,
    "ip": null,
    "ips": [],
    "locale": "en",
    "role": {
      "color": "",
      "highlighted": false,
      "id": "user",
      "name": "user",
      "permissions": "0"
    },
    "silenced": false,
    "suspended": false,
    "username": "1happyturtle"
  },
  {
    "account": {
      "acct": "admin",
      "avatar": "",
      "avatar_static": "",
      "bot": false,
      "created_at": "2022-05-17T13:10:59.000Z",
      "discoverable": true,
      "display_name": "",
      "emojis": [],
      "enable_rss": true,
      "fields": [],
      "followers_count": 1,
      "following_count": 1,
      "group": false,
      "header": "http://localhost:8080/assets/default_header.webp",
      "header_description": "Flat gray background (default header).",
      "header_static": "http://localhost:8080/assets/default_header.webp",
      "id": "01F8MH17FWEB39HZJ76B6VXSKF",
      "indexable": true,
      "last_status_at": "2021-10-20",
      "locked": false,
      "noindex": false,
      "note": "",
      "roles": [
        {
          "color": "",
          "id": "admin",
          "name": "admin"
        }
      ],
      "statuses_count": 4,
      "url": "http://localhost:8080/@admin",
      "username": "admin"
    },
    "approved": true,
    "confirmed": true,
    "created_at": "2022-05-17T13:10:59.000Z",
    "created_by_application_id": "01F8MGXQRHYF5QPMTMXP78QC2F",
    "disabled": false,
    "domain": null,
    "email": "admin@example.org",
    "id": "01F8MH17FWEB39HZJ76B6VXSKF",
    "invite_request": null,
    "ip": null,
    "ips": [],
    "locale": "en",
    "role": {
      "color": "",
      "highlighted": true,
      "id": "admin",
      "name": "admin",
      "permissions": "546033"
    },
    "silenced": false,
    "suspended": false,
    "username": "admin"
  },
  {
    "account": {
      "acct": "localhost:8080",
      "avatar": "",
      "avatar_static": "",
      "bot": true,
      "created_at": "2020-05-17T13:10:59.000Z",
      "discoverable": true,
      "display_name": "",
      "emojis": [],
      "fields": [],
      "followers_count": 0,
      "following_count": 0,
      "group": false,
      "header": "http://localhost:8080/assets/default_header.webp",
      "header_description": "Flat gray background (default header).",
      "header_static": "http://localhost:8080/assets/default_header.webp",
      "id": "01AY6P665V14JJR0AFVRT7311Y",
      "indexable": true,
      "last_status_at": null,
      "locked": false,
      "noindex": false,
      "note": "",
      "statuses_count": 0,
      "url": "http://localhost:8080/@localhost:8080",
      "username": "localhost:8080"
    },
    "approved": false,
    "confirmed": false,
    "created_at": "2020-05-17T13:10:59.000Z",
    "disabled": false,
    "domain": null,
    "email": "",
    "id": "01AY6P665V14JJR0AFVRT7311Y",
    "invite_request": null,
    "ip": null,
    "ips": [],
    "locale": "",
    "role": {
      "color": "",
      "highlighted": false,
      "id": "user",
      "name": "user",
      "permissions": "0"
    },
    "silenced": false,
    "suspended": false,
    "username": "localhost:8080"
  },
  {
    "account": {
      "acct": "media_mogul",
      "avatar": "http://localhost:8080/fileserver/01JPCMD83Y4WR901094YES3QC5/avatar/original/01JPHQZ0ZHC2AXJK1JQNXRXQZN.jpeg",
      "avatar_description": "DESCRIPTION_GOES_HERE",
      "avatar_media_id": "01JPHQZ0ZHC2AXJK1JQNXRXQZN",
      "avatar_static": "http://localhost:8080/fileserver/01JPCMD83Y4WR901094YES3QC5/avatar/small/01JPHQZ0ZHC2AXJK1JQNXRXQZN.jpeg",
      "bot": false,
      "created_at": "2025-03-15T11:08:00.000Z",
      "discoverable": false,
      "display_name": "",
      "emojis": [],
      "enable_rss": true,
      "fields": [
        {
          "name": "I'm going to post a lot of",
          "value": "media!",
          "verified_at": null
        },
        {
          "name": "and there's nothing",
          "value": "you can do about it",
          "verified_at": null
        }
      ],
      "followers_count": 0,
      "following_count": 0,
      "group": false,
      "header": "http://localhost:8080/fileserver/01JPCMD83Y4WR901094YES3QC5/header/original/01JPHRB7F2RXPTEQFRYC85EPD9.png",
      "header_description": "DESCRIPTION_GOES_HERE",
      "header_media_id": "01JPHRB7F2RXPTEQFRYC85EPD9",
      "header_static": "http://localhost:8080/fileserver/01JPCMD83Y4WR901094YES3QC5/header/small/01JPHRB7F2RXPTEQFRYC85EPD9.webp",
      "id": "01JPCMD83Y4WR901094YES3QC5",
      "indexable": false,
      "last_status_at": "2025-03-15",
      "locked": false,
      "noindex": true,
      "note": "<p>I'm a test account that posts a shitload of media and I have my account rendered in \"gallery\" mode</p>",
      "statuses_count": 2,
      "url": "http://localhost:8080/@media_mogul",
      "username": "media_mogul"
    },
    "approved": true,
    "confirmed": true,
    "created_at": "2025-03-15T11:08:00.000Z",
    "created_by_application_id": "01HT5P2YHDMPAAD500NDAY8JW1",
    "disabled": false,
    "domain": null,
    "email": "media.mogul@example.org",
    "id": "01JPCMD83Y4WR901094YES3QC5",
    "invite_request": null,
    "ip": null,
    "ips": [],
    "locale": "en",
    "role": {
      "color": "",
      "highlighted": false,
      "id": "user",
      "name": "user",
      "permissions": "0"
    },
    "silenced": false,
    "suspended": false,
    "username": "media_mogul"
  },
  {
    "account": {
      "acct": "the_mighty_zork",
      "avatar": "http://localhost:8080/fileserver/01F8MH1H7YV1Z7D2C8K2730QBF/avatar/original/01F8MH58A357CV5K7R7TJMSH6S.jpg",
      "avatar_description": "a green goblin looking nasty",
      "avatar_media_id": "01F8MH58A357CV5K7R7TJMSH6S",
      "avatar_static": "http://localhost:8080/fileserver/01F8MH1H7YV1Z7D2C8K2730QBF/avatar/small/01F8MH58A357CV5K7R7TJMSH6S.webp",
      "bot": false,
      "created_at": "2022-05-20T11:09:18.000Z",
      "discoverable": true,
      "display_name": "original zork (he/they)",
      "emojis": [],
      "enable_rss": true,
      "fields": [],
      "followers_count": 2,
      "following_count": 2,
      "group": false,
      "header": "http://localhost:8080/fileserver/01F8MH1H7YV1Z7D2C8K2730QBF/header/original/01PFPMWK2FF0D9WMHEJHR07C3Q.jpg",
      "header_description": "A very old-school screenshot of the original team fortress mod for quake",
      "header_media_id": "01PFPMWK2FF0D9WMHEJHR07C3Q",
      "header_static": "http://localhost:8080/fileserver/01F8MH1H7YV1Z7D2C8K2730QBF/header/small/01PFPMWK2FF0D9WMHEJHR07C3Q.webp",
      "id": "01F8MH1H7YV1Z7D2C8K2730QBF",
      "indexable": true,
      "last_status_at": "2024-11-01",
      "locked": false,
      "noindex": false,
      "note": "<p>hey yo this is my profile!</p>",
      "statuses_count": 9,
      "url": "http://localhost:8080/@the_mighty_zork",
      "username": "the_mighty_zork"
    },
    "approved": true,
    "confirmed": true,
    "created_at": "2022-05-20T11:09:18.000Z",
    "created_by_application_id": "01F8MGY43H3N2C8EWPR2FPYEXG",
    "disabled": false,
    "domain": null,
    "email": "zork@example.org",
    "id": "01F8MH1H7YV1Z7D2C8K2730QBF",
    "invite_request": "I wanna be on this damned webbed site so bad! Please! Wow",
    "ip": null,
    "ips": [],
    "locale": "en",
    "role": {
      "color": "",
      "highlighted": false,
      "id": "user",
      "name": "user",
      "permissions": "0"
    },
    "silenced": false,
    "suspended": false,
    "username": "the_mighty_zork"
  },
  {
    "account": {
      "acct": "weed_lord420",
      "avatar": "",
      "avatar_static": "",
      "bot": false,
      "created_at": "2022-06-04T13:12:00.000Z",
      "discoverable": false,
      "display_name": "",
      "emojis": [],
      "fields": [],
      "followers_count": 0,
      "following_count": 0,
      "group": false,
      "header": "http://localhost:8080/assets/default_header.webp",
      "header_description": "Flat gray background (default header).",
      "header_static": "http://localhost:8080/assets/default_header.webp",
      "id": "01F8MH0BBE4FHXPH513MBVFHB0",
      "indexable": false,
      "last_status_at": null,
      "locked": false,
      "noindex": true,
      "note": "",
      "statuses_count": 0,
      "url": "http://localhost:8080/@weed_lord420",
      "username": "weed_lord420"
    },
    "approved": false,
    "confirmed": false,
    "created_at": "2022-06-04T13:12:00.000Z",
    "created_by_application_id": "01F8MGY43H3N2C8EWPR2FPYEXG",
    "disabled": false,
    "domain": null,
    "email": "weed_lord420@example.org",
    "id": "01F8MH0BBE4FHXPH513MBVFHB0",
    "invite_request": "hi, please let me in! I'm looking for somewhere neato bombeato to hang out.",
    "ip": "199.222.111.89",
    "ips": [],
    "locale": "en",
    "role": {
      "color": "",
      "highlighted": false,
      "id": "user",
      "name": "user",
      "permissions": "0"
    },
    "silenced": false,
    "suspended": false,
    "username": "weed_lord420"
  },
  {
    "account": {
      "acct": "Some_User@example.org",
      "avatar": "",
      "avatar_static": "",
      "bot": false,
      "created_at": "2020-08-10T12:13:28.000Z",
      "discoverable": true,
      "display_name": "some user",
      "emojis": [],
      "fields": [],
      "followers_count": 0,
      "following_count": 0,
      "group": false,
      "header": "http://localhost:8080/assets/default_header.webp",
      "header_description": "Flat gray background (default header).",
      "header_static": "http://localhost:8080/assets/default_header.webp",
      "id": "01FHMQX3GAABWSM0S2VZEC2SWC",
      "indexable": true,
      "last_status_at": "2023-11-02",
      "locked": true,
      "noindex": false,
      "note": "i'm a real son of a gun",
      "statuses_count": 1,
      "url": "http://example.org/@Some_User",
      "username": "Some_User"
    },
    "approved": false,
    "confirmed": false,
    "created_at": "2020-08-10T12:13:28.000Z",
    "disabled": false,
    "domain": "example.org",
    "email": "",
    "id": "01FHMQX3GAABWSM0S2VZEC2SWC",
    "invite_request": null,
    "ip": null,
    "ips": [],
    "locale": "",
    "role": {
      "color": "",
      "highlighted": false,
      "id": "user",
      "name": "user",
      "permissions": "0"
    },
    "silenced": false,
    "suspended": false,
    "username": "Some_User"
  },
  {
    "account": {
      "acct": "foss_satan@fossbros-anonymous.io",
      "avatar": "",
      "avatar_static": "",
      "bot": false,
      "created_at": "2021-09-26T10:52:36.000Z",
      "discoverable": true,
      "display_name": "big gerald",
      "emojis": [],
      "fields": [],
      "followers_count": 0,
      "following_count": 0,
      "group": false,
      "header": "http://localhost:8080/assets/default_header.webp",
      "header_description": "Flat gray background (default header).",
      "header_static": "http://localhost:8080/assets/default_header.webp",
      "id": "01F8MH5ZK5VRH73AKHQM6Y9VNX",
      "indexable": true,
      "last_status_at": "2024-11-01",
      "locked": false,
      "noindex": false,
      "note": "i post about like, i dunno, stuff, or whatever!!!!",
      "statuses_count": 4,
      "url": "http://fossbros-anonymous.io/@foss_satan",
      "username": "foss_satan"
    },
    "approved": false,
    "confirmed": false,
    "created_at": "2021-09-26T10:52:36.000Z",
    "disabled": false,
    "domain": "fossbros-anonymous.io",
    "email": "",
    "id": "01F8MH5ZK5VRH73AKHQM6Y9VNX",
    "invite_request": null,
    "ip": null,
    "ips": [],
    "locale": "",
    "role": {
      "color": "",
      "highlighted": false,
      "id": "user",
      "name": "user",
      "permissions": "0"
    },
    "silenced": false,
    "suspended": false,
    "username": "foss_satan"
  },
  {
    "account": {
      "acct": "her_fuckin_maj@thequeenisstillalive.technology",
      "avatar": "",
      "avatar_static": "",
      "bot": false,
      "created_at": "2020-08-10T12:13:28.000Z",
      "discoverable": true,
      "display_name": "lizzzieeeeeeeeeeee",
      "emojis": [],
      "fields": [],
      "followers_count": 0,
      "following_count": 0,
      "group": false,
      "header": "http://localhost:8080/fileserver/062G5WYKY35KKD12EMSM3F8PJ8/header/original/01G549FP8065NKWBPTWHP6Y3PD.jpg",
      "header_description": "tweet from thoughts of dog: i drank. all the water. in my bowl. earlier. but just now. i returned. to the same bowl. and it was. full again.. the bowl. is haunted",
      "header_media_id": "01G549FP8065NKWBPTWHP6Y3PD",
      "header_static": "http://localhost:8080/fileserver/062G5WYKY35KKD12EMSM3F8PJ8/header/small/01G549FP8065NKWBPTWHP6Y3PD.webp",
      "id": "062G5WYKY35KKD12EMSM3F8PJ8",
      "indexable": true,
      "last_status_at": null,
      "locked": true,
      "noindex": false,
      "note": "if i die blame charles don't let that fuck become king",
      "statuses_count": 0,
      "url": "http://thequeenisstillalive.technology/@her_fuckin_maj",
      "username": "her_fuckin_maj"
    },
    "approved": false,
    "confirmed": false,
    "created_at": "2020-08-10T12:13:28.000Z",
    "disabled": false,
    "domain": "thequeenisstillalive.technology",
    "email": "",
    "id": "062G5WYKY35KKD12EMSM3F8PJ8",
    "invite_request": null,
    "ip": null,
    "ips": [],
    "locale": "",
    "role": {
      "color": "",
      "highlighted": false,
      "id": "user",
      "name": "user",
      "permissions": "0"
    },
    "silenced": false,
    "suspended": false,
    "username": "her_fuckin_maj"
  },
  {
    "account": {
      "acct": "üser@ëxample.org",
      "avatar": "",
      "avatar_static": "",
      "bot": false,
      "created_at": "2020-08-10T12:13:28.000Z",
      "discoverable": false,
      "display_name": "",
      "emojis": [],
      "fields": [],
      "followers_count": 0,
      "following_count": 0,
      "group": false,
      "header": "http://localhost:8080/assets/default_header.webp",
      "header_description": "Flat gray background (default header).",
      "header_static": "http://localhost:8080/assets/default_header.webp",
      "id": "07GZRBAEMBNKGZ8Z9VSKSXKR98",
      "indexable": false,
      "last_status_at": null,
      "locked": false,
      "noindex": true,
      "note": "",
      "statuses_count": 0,
      "url": "https://xn--xample-ova.org/users/@%C3%BCser",
      "username": "üser"
    },
    "approved": false,
    "confirmed": false,
    "created_at": "2020-08-10T12:13:28.000Z",
    "disabled": false,
    "domain": "ëxample.org",
    "email": "",
    "id": "07GZRBAEMBNKGZ8Z9VSKSXKR98",
    "invite_request": null,
    "ip": null,
    "ips": [],
    "locale": "",
    "role": {
      "color": "",
      "highlighted": false,
      "id": "user",
      "name": "user",
      "permissions": "0"
    },
    "silenced": false,
    "suspended": false,
    "username": "üser"
  }
]`, out)
}

func (suite *AccountsGetTestSuite) TestAccountsMinID() {
	recorder := httptest.NewRecorder()

	path := admin.AccountsV2Path + "?limit=1&min_id=/@admin"
	c := suite.newContext(recorder, http.MethodGet, nil, path, "application/json")

	c.SetPathValue("min_id", "@admin")
	c.SetPathValue("limit", "1")

	suite.adminModule.AccountsGETV2Handler(c)
	suite.Equal(http.StatusOK, recorder.Code)

	b, err := io.ReadAll(recorder.Body)
	if err != nil {
		suite.FailNow(err.Error())
	}

	link := recorder.Header().Get("Link")
	suite.Equal(`<http://localhost:8080/api/v2/admin/accounts?limit=1&max_id=%2F%401happyturtle>; rel="next", <http://localhost:8080/api/v2/admin/accounts?limit=1&min_id=%2F%401happyturtle>; rel="prev"`, link)

	out := testrig.MustJSONStringFromBytes(b)
	suite.Equal(`[
  {
    "account": {
      "acct": "1happyturtle",
      "avatar": "",
      "avatar_static": "",
      "bot": false,
      "created_at": "2022-06-04T13:12:00.000Z",
      "discoverable": false,
      "display_name": "happy little turtle :3",
      "emojis": [],
      "fields": [
        {
          "name": "should you follow me?",
          "value": "maybe!",
          "verified_at": null
        },
        {
          "name": "age",
          "value": "120",
          "verified_at": null
        }
      ],
      "followers_count": 1,
      "following_count": 1,
      "group": false,
      "header": "http://localhost:8080/assets/default_header.webp",
      "header_description": "Flat gray background (default header).",
      "header_static": "http://localhost:8080/assets/default_header.webp",
      "hide_collections": true,
      "id": "01F8MH5NBDF2MV7CTC4Q5128HF",
      "indexable": false,
      "last_status_at": "2026-01-01",
      "locked": true,
      "noindex": true,
      "note": "<p>i post about things that concern me</p>",
      "statuses_count": 10,
      "url": "http://localhost:8080/@1happyturtle",
      "username": "1happyturtle"
    },
    "approved": true,
    "confirmed": true,
    "created_at": "2022-06-04T13:12:00.000Z",
    "created_by_application_id": "01F8MGY43H3N2C8EWPR2FPYEXG",
    "disabled": false,
    "domain": null,
    "email": "tortle.dude@example.org",
    "id": "01F8MH5NBDF2MV7CTC4Q5128HF",
    "invite_request": null,
    "ip": null,
    "ips": [],
    "locale": "en",
    "role": {
      "color": "",
      "highlighted": false,
      "id": "user",
      "name": "user",
      "permissions": "0"
    },
    "silenced": false,
    "suspended": false,
    "username": "1happyturtle"
  }
]`, out)
}

func TestAccountsGetTestSuite(t *testing.T) {
	suite.Run(t, &AccountsGetTestSuite{})
}
