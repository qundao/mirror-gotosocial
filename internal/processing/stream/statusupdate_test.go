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

package stream_test

import (
	"testing"

	"code.superseriousbusiness.org/gotosocial/internal/stream"
	"code.superseriousbusiness.org/gotosocial/internal/typeutils"
	"code.superseriousbusiness.org/gotosocial/testrig"
	"github.com/stretchr/testify/suite"
)

type StatusUpdateTestSuite struct {
	StreamTestSuite
}

func (suite *StatusUpdateTestSuite) TestStreamNotification() {
	account := suite.testAccounts["local_account_1"]

	openStream := suite.streamProcessor.Open(suite.T().Context(), account, "user")

	editedStatus := suite.testStatuses["remote_account_1_status_1"]
	apiStatus, err := typeutils.NewConverter(suite.state).StatusToAPIStatus(suite.T().Context(), editedStatus, account)
	suite.NoError(err)

	suite.streamProcessor.StatusUpdate(suite.T().Context(), account, apiStatus, stream.TimelineHome)

	msg, ok := openStream.Recv(suite.T().Context())
	suite.True(ok)

	out := testrig.MustJSONStringFromString(msg.Payload)
	suite.Equal(`{
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
  "bookmarked": false,
  "card": null,
  "content": "<p>dark souls status bot: \"thoughts of dog\"</p>",
  "created_at": "2021-09-20T10:40:37.000Z",
  "edited_at": null,
  "emojis": [],
  "favourited": false,
  "favourites_count": 0,
  "id": "01FVW7JHQFSFK166WWKR8CBA6M",
  "in_reply_to_account_id": null,
  "in_reply_to_id": null,
  "interaction_policy": {
    "can_favourite": {
      "automatic_approval": [
        "public",
        "me"
      ],
      "manual_approval": []
    },
    "can_reblog": {
      "automatic_approval": [
        "public",
        "me"
      ],
      "manual_approval": []
    },
    "can_reply": {
      "automatic_approval": [
        "public",
        "me"
      ],
      "manual_approval": []
    }
  },
  "language": "en",
  "media_attachments": [
    {
      "blurhash": "L3Q9_@4n9E?axW4mD$Mx~q00Di%L",
      "description": "tweet from thoughts of dog: i drank. all the water. in my bowl. earlier. but just now. i returned. to the same bowl. and it was. full again.. the bowl. is haunted",
      "id": "01FVW7RXPQ8YJHTEXYPE7Q8ZY0",
      "meta": {
        "focus": {
          "x": 0,
          "y": 0
        },
        "original": {
          "aspect": 1.6219932,
          "height": 291,
          "size": "472x291",
          "width": 472
        },
        "small": {
          "aspect": 1.6219932,
          "height": 291,
          "size": "472x291",
          "width": 472
        }
      },
      "preview_remote_url": null,
      "preview_url": "http://localhost:8080/fileserver/01F8MH5ZK5VRH73AKHQM6Y9VNX/attachment/small/01FVW7RXPQ8YJHTEXYPE7Q8ZY0.webp",
      "remote_url": "http://fossbros-anonymous.io/attachments/original/13bbc3f8-2b5e-46ea-9531-40b4974d9912.jpg",
      "text_url": "http://localhost:8080/fileserver/01F8MH5ZK5VRH73AKHQM6Y9VNX/attachment/original/01FVW7RXPQ8YJHTEXYPE7Q8ZY0.jpg",
      "type": "image",
      "url": "http://localhost:8080/fileserver/01F8MH5ZK5VRH73AKHQM6Y9VNX/attachment/original/01FVW7RXPQ8YJHTEXYPE7Q8ZY0.jpg"
    }
  ],
  "mentions": [],
  "muted": false,
  "pinned": false,
  "poll": null,
  "reblog": null,
  "reblogged": false,
  "reblogs_count": 0,
  "replies_count": 0,
  "sensitive": true,
  "spoiler_text": "potentially annoying post ahead",
  "tags": [],
  "uri": "http://fossbros-anonymous.io/users/foss_satan/statuses/01FVW7JHQFSFK166WWKR8CBA6M",
  "url": "http://fossbros-anonymous.io/@foss_satan/statuses/01FVW7JHQFSFK166WWKR8CBA6M",
  "visibility": "unlisted"
}`, out)
	suite.Equal(msg.Event, "status.update")
}

func TestStatusUpdateTestSuite(t *testing.T) {
	suite.Run(t, &StatusUpdateTestSuite{})
}
