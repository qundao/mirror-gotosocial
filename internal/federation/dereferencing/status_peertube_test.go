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

package dereferencing_test

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"testing"

	"code.superseriousbusiness.org/gotosocial/internal/federation/dereferencing"
	"code.superseriousbusiness.org/gotosocial/internal/filter/interaction"
	"code.superseriousbusiness.org/gotosocial/internal/filter/relay"
	"code.superseriousbusiness.org/gotosocial/internal/filter/visibility"
	"code.superseriousbusiness.org/gotosocial/internal/media"
	"code.superseriousbusiness.org/gotosocial/testrig"
	"github.com/stretchr/testify/suite"
)

// Peertube video objects are quite
// unique, so this test gets its own file.
const peertubeVideoJSON = `{
  "@context": [
    "https://www.w3.org/ns/activitystreams",
    "https://w3id.org/security/v1",
    {
      "RsaSignature2017": "https://w3id.org/security#RsaSignature2017"
    },
    {
      "pt": "https://joinpeertube.org/ns#",
      "sc": "http://schema.org/",
      "Hashtag": "as:Hashtag",
      "category": "sc:category",
      "licence": "sc:license",
      "subtitleLanguage": "sc:subtitleLanguage",
      "automaticallyGenerated": "pt:automaticallyGenerated",
      "sensitive": "as:sensitive",
      "language": "sc:inLanguage",
      "identifier": "sc:identifier",
      "isLiveBroadcast": "sc:isLiveBroadcast",
      "liveSaveReplay": {
        "@type": "sc:Boolean",
        "@id": "pt:liveSaveReplay"
      },
      "permanentLive": {
        "@type": "sc:Boolean",
        "@id": "pt:permanentLive"
      },
      "latencyMode": {
        "@type": "sc:Number",
        "@id": "pt:latencyMode"
      },
      "dvrWindow": {
        "@type": "sc:Duration",
        "@id": "pt:dvrWindow"
      },
      "Infohash": "pt:Infohash",
      "SensitiveTag": "pt:SensitiveTag",
      "tileWidth": {
        "@type": "sc:Number",
        "@id": "pt:tileWidth"
      },
      "tileHeight": {
        "@type": "sc:Number",
        "@id": "pt:tileHeight"
      },
      "tileDuration": {
        "@type": "sc:Number",
        "@id": "pt:tileDuration"
      },
      "aspectRatio": {
        "@type": "sc:Float",
        "@id": "pt:aspectRatio"
      },
      "uuid": {
        "@type": "sc:identifier",
        "@id": "pt:uuid"
      },
      "originallyPublishedAt": "sc:datePublished",
      "schedules": "sc:eventSchedule",
      "startDate": "sc:startDate",
      "embedUrl": "sc:embedUrl",
      "uploadDate": "sc:uploadDate",
      "hasParts": "sc:hasParts",
      "playerSettings": "pt:playerSettings",
      "views": {
        "@type": "sc:Number",
        "@id": "pt:views"
      },
      "downloads": {
        "@type": "sc:Number",
        "@id": "pt:downloads"
      },
      "state": {
        "@type": "sc:Number",
        "@id": "pt:state"
      },
      "size": {
        "@type": "sc:Number",
        "@id": "pt:size"
      },
      "fps": {
        "@type": "sc:Number",
        "@id": "pt:fps"
      },
      "canReply": "pt:canReply",
      "commentsPolicy": {
        "@type": "sc:Number",
        "@id": "pt:commentsPolicy"
      },
      "downloadEnabled": {
        "@type": "sc:Boolean",
        "@id": "pt:downloadEnabled"
      },
      "waitTranscoding": {
        "@type": "sc:Boolean",
        "@id": "pt:waitTranscoding"
      },
      "support": {
        "@type": "sc:Text",
        "@id": "pt:support"
      },
      "likes": {
        "@id": "as:likes",
        "@type": "@id"
      },
      "dislikes": {
        "@id": "as:dislikes",
        "@type": "@id"
      },
      "shares": {
        "@id": "as:shares",
        "@type": "@id"
      },
      "comments": {
        "@id": "as:comments",
        "@type": "@id"
      },
      "PropertyValue": "sc:PropertyValue",
      "value": "sc:value"
    },
    {
      "ptlc": "urn:peertube-plugin-livechat",
      "peertubeLiveChat": {
        "@id": "ptlc:peertubeLiveChat",
        "@type": "@json"
      }
    }
  ],
  "to": [
    "https://www.w3.org/ns/activitystreams#Public",
    "https://spectra.video/video-channels/rumirum"
  ],
  "cc": [
    "https://spectra.video/accounts/mberos/followers"
  ],
  "type": "Video",
  "id": "https://spectra.video/videos/watch/4e06b13a-b7bb-4530-a653-c82d47a7bfb8",
  "name": "Tape Recording!!!",
  "duration": "PT87S",
  "uuid": "4e06b13a-b7bb-4530-a653-c82d47a7bfb8",
  "views": 60,
  "downloads": 89,
  "sensitive": false,
  "summary": null,
  "waitTranscoding": false,
  "state": 1,
  "canReply": null,
  "commentsPolicy": 1,
  "downloadEnabled": true,
  "published": "2026-08-18T18:50:24.225Z",
  "originallyPublishedAt": null,
  "schedules": [],
  "updated": "2026-10-27T17:02:04.768Z",
  "tag": [],
  "mediaType": "text/markdown",
  "content": "I recently found an old 90s tape recorder (Marantz PMD222) and fixed the belts, so I can can use it like this!! Pushed the gain waaaay too much on this, so I'd have to learn how to find the sweetspot for that analogue tape saturation. ",
  "support": null,
  "subtitleLanguage": [
    {
      "identifier": "cy",
      "name": "Welsh",
      "automaticallyGenerated": true,
      "url": [
        {
          "type": "Link",
          "mediaType": "text/vtt",
          "href": "https://spectra-prod.us-east-1.linodeobjects.com/captions/92837ca4-2496-46f5-a189-d07722c07629-cy.vtt"
        },
        {
          "type": "Link",
          "mediaType": "application/x-mpegURL",
          "href": "https://spectra-prod.us-east-1.linodeobjects.com/streaming-playlistshls/4e06b13a-b7bb-4530-a653-c82d47a7bfb8/92837ca4-2496-46f5-a189-d07722c07629-cy.m3u8"
        }
      ]
    }
  ],
  "icon": [
    {
      "type": "Image",
      "url": "https://spectra.video/lazy-static/thumbnails/2971dd32-8243-45a9-8bf7-71d1f45fea77.jpg",
      "mediaType": "image/jpeg",
      "width": 280,
      "height": 157
    },
    {
      "type": "Image",
      "url": "https://spectra.video/lazy-static/thumbnails/6b0c5695-b5c3-4780-a00e-0c16d9689399.jpg",
      "mediaType": "image/jpeg",
      "width": 850,
      "height": 480
    }
  ],
  "preview": [
    {
      "type": "Image",
      "rel": [
        "storyboard"
      ],
      "url": [
        {
          "mediaType": "image/jpeg",
          "href": "https://spectra.video/lazy-static/storyboards/0b510407-e709-4cc2-a0b8-bf1aa894eedf.jpg",
          "width": 2112,
          "height": 864,
          "tileWidth": 192,
          "tileHeight": 108,
          "tileDuration": "PT1S"
        }
      ]
    }
  ],
  "aspectRatio": 1.7778,
  "url": [
    {
      "type": "Link",
      "mediaType": "text/html",
      "href": "https://spectra.video/w/aCQ9f8rMn14CW8jBzCTiA9"
    },
    {
      "type": "Link",
      "mediaType": "text/html",
      "href": "https://spectra.video/videos/watch/4e06b13a-b7bb-4530-a653-c82d47a7bfb8"
    },
    {
      "type": "Link",
      "mediaType": "video/mp4",
      "href": "https://spectra-prod.us-east-1.linodeobjects.com/videos/83af9e16-f4f7-4570-b251-361247693de0-0.mp4",
      "height": 0,
      "width": 0,
      "size": 1753703,
      "fps": 0,
      "attachment": [
        {
          "type": "PropertyValue",
          "name": "ffprobe_codec_type",
          "value": "audio"
        },
        {
          "type": "PropertyValue",
          "name": "peertube_format_flag",
          "value": "web-video"
        }
      ]
    },
    {
      "type": "Link",
      "rel": [
        "metadata",
        "video/mp4"
      ],
      "mediaType": "application/json",
      "href": "https://spectra.video/api/v1/videos/4e06b13a-b7bb-4530-a653-c82d47a7bfb8/metadata/1193828",
      "height": 0,
      "width": 0,
      "fps": 0
    },
    {
      "type": "Link",
      "mediaType": "application/x-bittorrent",
      "href": "https://spectra.video/lazy-static/torrents/409f1127-4a32-4aa8-94bf-ed721b5933a2-0.torrent",
      "height": 0,
      "width": 0,
      "fps": 0
    },
    {
      "type": "Link",
      "mediaType": "application/x-bittorrent;x-scheme-handler/magnet",
      "href": "magnet:?xs=https%3A%2F%2Fspectra.video%2Flazy-static%2Ftorrents%2F409f1127-4a32-4aa8-94bf-ed721b5933a2-0.torrent&xt=urn:btih:2f4ff5e41e12f0ab0c96482b52465d2edb247829&dn=Tape+Recording!!!&tr=https%3A%2F%2Fspectra.video%2Ftracker%2Fannounce&tr=wss%3A%2F%2Fspectra.video%3A443%2Ftracker%2Fsocket&ws=https%3A%2F%2Fspectra-prod.us-east-1.linodeobjects.com%2Fvideos%2F83af9e16-f4f7-4570-b251-361247693de0-0.mp4",
      "height": 0,
      "width": 0,
      "fps": 0
    },
    {
      "type": "Link",
      "mediaType": "application/x-mpegURL",
      "href": "https://spectra-prod.us-east-1.linodeobjects.com/streaming-playlistshls/4e06b13a-b7bb-4530-a653-c82d47a7bfb8/4f0101a1-3791-46b9-963b-43c444fd2930-master.m3u8",
      "tag": [
        {
          "type": "Infohash",
          "name": "rGDjZMArUP2AecH1IfpS"
        },
        {
          "type": "Infohash",
          "name": "dIM2f411lgJRQior+HfU"
        },
        {
          "type": "Link",
          "name": "sha256",
          "mediaType": "application/json",
          "href": "https://spectra-prod.us-east-1.linodeobjects.com/streaming-playlistshls/4e06b13a-b7bb-4530-a653-c82d47a7bfb8/06f952d7-6045-40b5-8c58-eb4b1d9fa95c-segments-sha256.json"
        },
        {
          "type": "Link",
          "mediaType": "video/mp4",
          "href": "https://spectra-prod.us-east-1.linodeobjects.com/streaming-playlistshls/4e06b13a-b7bb-4530-a653-c82d47a7bfb8/97226d01-15e3-4db6-a6ef-d62e960254b3-720-fragmented.mp4",
          "height": 720,
          "width": 1280,
          "size": 24268457,
          "fps": 25,
          "attachment": [
            {
              "type": "PropertyValue",
              "name": "ffprobe_codec_type",
              "value": "audio"
            },
            {
              "type": "PropertyValue",
              "name": "ffprobe_codec_type",
              "value": "video"
            },
            {
              "type": "PropertyValue",
              "name": "peertube_format_flag",
              "value": "fragmented"
            }
          ]
        },
        {
          "type": "Link",
          "rel": [
            "metadata",
            "video/mp4"
          ],
          "mediaType": "application/json",
          "href": "https://spectra.video/api/v1/videos/4e06b13a-b7bb-4530-a653-c82d47a7bfb8/metadata/1193827",
          "height": 720,
          "width": 1280,
          "fps": 25
        },
        {
          "type": "Link",
          "mediaType": "application/x-bittorrent",
          "href": "https://spectra.video/lazy-static/torrents/c677d4f9-e55d-46cf-bfb8-4a62bd6b0c66-720-hls.torrent",
          "height": 720,
          "width": 1280,
          "fps": 25
        },
        {
          "type": "Link",
          "mediaType": "application/x-bittorrent;x-scheme-handler/magnet",
          "href": "magnet:?xs=https%3A%2F%2Fspectra.video%2Flazy-static%2Ftorrents%2Fc677d4f9-e55d-46cf-bfb8-4a62bd6b0c66-720-hls.torrent&xt=urn:btih:673e14a16d0d4ebd4f9b900530041dfe3d5e34ba&dn=Tape+Recording!!!&tr=https%3A%2F%2Fspectra.video%2Ftracker%2Fannounce&tr=wss%3A%2F%2Fspectra.video%3A443%2Ftracker%2Fsocket&ws=https%3A%2F%2Fspectra-prod.us-east-1.linodeobjects.com%2Fstreaming-playlistshls%2F4e06b13a-b7bb-4530-a653-c82d47a7bfb8%2F97226d01-15e3-4db6-a6ef-d62e960254b3-720-fragmented.mp4",
          "height": 720,
          "width": 1280,
          "fps": 25
        },
        {
          "type": "Link",
          "mediaType": "video/mp4",
          "href": "https://spectra-prod.us-east-1.linodeobjects.com/streaming-playlistshls/4e06b13a-b7bb-4530-a653-c82d47a7bfb8/20c6372a-da1d-43dd-96e1-4abe931d8b0f-360-fragmented.mp4",
          "height": 360,
          "width": 640,
          "size": 5118078,
          "fps": 25,
          "attachment": [
            {
              "type": "PropertyValue",
              "name": "ffprobe_codec_type",
              "value": "audio"
            },
            {
              "type": "PropertyValue",
              "name": "ffprobe_codec_type",
              "value": "video"
            },
            {
              "type": "PropertyValue",
              "name": "peertube_format_flag",
              "value": "fragmented"
            }
          ]
        },
        {
          "type": "Link",
          "rel": [
            "metadata",
            "video/mp4"
          ],
          "mediaType": "application/json",
          "href": "https://spectra.video/api/v1/videos/4e06b13a-b7bb-4530-a653-c82d47a7bfb8/metadata/1193829",
          "height": 360,
          "width": 640,
          "fps": 25
        },
        {
          "type": "Link",
          "mediaType": "application/x-bittorrent",
          "href": "https://spectra.video/lazy-static/torrents/70205937-fa8a-413d-8b83-d9c215b38d6f-360-hls.torrent",
          "height": 360,
          "width": 640,
          "fps": 25
        },
        {
          "type": "Link",
          "mediaType": "application/x-bittorrent;x-scheme-handler/magnet",
          "href": "magnet:?xs=https%3A%2F%2Fspectra.video%2Flazy-static%2Ftorrents%2F70205937-fa8a-413d-8b83-d9c215b38d6f-360-hls.torrent&xt=urn:btih:e30036bdc3c7ed84175aceffd68d2a5202375834&dn=Tape+Recording!!!&tr=https%3A%2F%2Fspectra.video%2Ftracker%2Fannounce&tr=wss%3A%2F%2Fspectra.video%3A443%2Ftracker%2Fsocket&ws=https%3A%2F%2Fspectra-prod.us-east-1.linodeobjects.com%2Fstreaming-playlistshls%2F4e06b13a-b7bb-4530-a653-c82d47a7bfb8%2F20c6372a-da1d-43dd-96e1-4abe931d8b0f-360-fragmented.mp4",
          "height": 360,
          "width": 640,
          "fps": 25
        }
      ]
    },
    {
      "type": "Link",
      "name": "tracker-http",
      "rel": [
        "tracker",
        "http"
      ],
      "href": "https://spectra.video/tracker/announce"
    },
    {
      "type": "Link",
      "name": "tracker-websocket",
      "rel": [
        "tracker",
        "websocket"
      ],
      "href": "wss://spectra.video:443/tracker/socket"
    }
  ],
  "likes": "https://spectra.video/videos/watch/4e06b13a-b7bb-4530-a653-c82d47a7bfb8/likes",
  "dislikes": "https://spectra.video/videos/watch/4e06b13a-b7bb-4530-a653-c82d47a7bfb8/dislikes",
  "shares": "https://spectra.video/videos/watch/4e06b13a-b7bb-4530-a653-c82d47a7bfb8/announces",
  "comments": "https://spectra.video/videos/watch/4e06b13a-b7bb-4530-a653-c82d47a7bfb8/comments",
  "hasParts": "https://spectra.video/videos/watch/4e06b13a-b7bb-4530-a653-c82d47a7bfb8/chapters",
  "playerSettings": "https://spectra.video/videos/watch/4e06b13a-b7bb-4530-a653-c82d47a7bfb8/player-settings",
  "embedUrl": "https://spectra.video/videos/embed/aCQ9f8rMn14CW8jBzCTiA9",
  "attributedTo": [
    {
      "type": "Person",
      "id": "https://spectra.video/accounts/mberos"
    },
    {
      "type": "Group",
      "id": "https://spectra.video/video-channels/rumirum"
    }
  ],
  "audience": "https://spectra.video/video-channels/rumirum",
  "isLiveBroadcast": false,
  "liveSaveReplay": null,
  "permanentLive": null,
  "latencyMode": null,
  "dvrWindow": null
}`

const peertubePersonJSON = `{
  "@context": [
    "https://www.w3.org/ns/activitystreams",
    "https://w3id.org/security/v1",
    {
      "RsaSignature2017": "https://w3id.org/security#RsaSignature2017"
    },
    {
      "pt": "https://joinpeertube.org/ns#",
      "sc": "http://schema.org/",
      "playerSettings": "pt:playerSettings",
      "playlists": {
        "@id": "pt:playlists",
        "@type": "@id"
      },
      "support": {
        "@type": "sc:Text",
        "@id": "pt:support"
      },
      "lemmy": "https://join-lemmy.org/ns#",
      "postingRestrictedToMods": "lemmy:postingRestrictedToMods",
      "toot": "http://joinmastodon.org/ns#",
      "indexable": "toot:indexable",
      "discoverable": "toot:discoverable",
      "email": "sc:email"
    }
  ],
  "type": "Person",
  "id": "https://spectra.video/accounts/mberos",
  "following": "https://spectra.video/accounts/mberos/following",
  "followers": "https://spectra.video/accounts/mberos/followers",
  "playlists": "https://spectra.video/accounts/mberos/playlists",
  "inbox": "https://spectra.video/accounts/mberos/inbox",
  "outbox": "https://spectra.video/accounts/mberos/outbox",
  "preferredUsername": "mberos",
  "name": "Mathias Berg Rosendal",
  "endpoints": {
    "sharedInbox": "https://spectra.video/inbox"
  },
  "publicKey": {
    "id": "https://spectra.video/accounts/mberos#main-key",
    "owner": "https://spectra.video/accounts/mberos",
    "publicKeyPem": "-----BEGIN PUBLIC KEY-----\nMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAnCn+BykDS0yjCHl5vlIH\nyQ1TOXyrIJ8u33ID3EcQobGFIrsheWIzS9Enaz3RVYavfPgnSTtL9sD8CxhmRKW3\nUSxbMpHh7L/N6ztr3GvEtMbhWjTGz2OskDCXpkEFvtlnDmbvvFfKakTV8I7UrrrG\ncbQG5qWUf43P7cw2u8xlDz/XM/YT3qnoys9pPXWc5lDFYFkGH++e4EmNjIDuUWCt\nTwRdgspDHycTYFlteI0w+qfsl7Z/C6bpkKxEFcf7e9lI/y6L46gC1KwSffKj14B6\nOfdZy17iXfZefOMiIgGQKV/EAV79hiAuBvkdH3ntFlNUrVA3Et2XR0/Dv1nGZCRZ\nbQIDAQAB\n-----END PUBLIC KEY-----\n"
  },
  "published": "2026-06-23T11:19:42.920Z",
  "indexable": true,
  "discoverable": true,
  "icon": [
    {
      "type": "Image",
      "mediaType": "image/jpeg",
      "height": 120,
      "width": 120,
      "url": "https://spectra.video/lazy-static/avatars/58664243-7313-414a-8073-0b90fb2f7608.jpg"
    },
    {
      "type": "Image",
      "mediaType": "image/jpeg",
      "height": 48,
      "width": 48,
      "url": "https://spectra.video/lazy-static/avatars/b63937c8-757d-4c7a-b1a7-deafb9dd7632.jpg"
    },
    {
      "type": "Image",
      "mediaType": "image/jpeg",
      "height": 600,
      "width": 600,
      "url": "https://spectra.video/lazy-static/avatars/f5b02a08-ad8e-4ac3-9b3f-c2e643b8c8e4.jpg"
    },
    {
      "type": "Image",
      "mediaType": "image/jpeg",
      "height": 1500,
      "width": 1500,
      "url": "https://spectra.video/lazy-static/avatars/8014a1b8-b684-4e2b-a272-ace4e6043c72.jpg"
    }
  ],
  "url": [
    {
      "type": "Link",
      "mediaType": "text/html",
      "href": "https://spectra.video/a/mberos/video-channels"
    },
    {
      "type": "Link",
      "mediaType": "text/html",
      "href": "https://spectra.video/a/mberos"
    },
    {
      "type": "Link",
      "mediaType": "text/html",
      "href": "https://spectra.video/accounts/mberos"
    }
  ],
  "summary": "Website: https://rumirum.eu\nMastodon: https://mastodon.social/@mberos\nPixelfed: https://pixelfed.social/mberos\n"
}`

const peertubeWebfingerJSON = `{
  "subject": "acct:mberos@spectra.video",
  "aliases": [
    "https://spectra.video/accounts/mberos"
  ],
  "links": [
    {
      "rel": "self",
      "type": "application/activity+json",
      "href": "https://spectra.video/accounts/mberos"
    },
    {
      "rel": "http://ostatus.org/schema/1.0/subscribe",
      "template": "https://spectra.video/remote-interaction?uri={uri}"
    }
  ]
}`

type StatusPeertubeTestSuite struct {
	DereferencerTestSuite
}

func (suite *StatusPeertubeTestSuite) TestPeertubeVideoDereference() {
	// Set up our test structs + tear down on finish.
	testStructs := testrig.SetupTestStructs(rMediaPath, rTemplatePath)
	defer testrig.TearDownTestStructs(testStructs)

	// Clean up test context when done.
	ctx, cncl := context.WithCancel(suite.T().Context())
	defer cncl()

	// Replace the HTTP client with a version that
	// responds to requests for peertube objects.
	do := func(req *http.Request) (*http.Response, error) {
		var (
			responseCode        int
			responseBytes       []byte
			responseContentType string
		)

		switch req.URL.String() {
		case "https://spectra.video/videos/watch/4e06b13a-b7bb-4530-a653-c82d47a7bfb8":
			// Video.
			responseCode = http.StatusOK
			responseBytes = []byte(peertubeVideoJSON)
			responseContentType = "application/activity+json"

		case "https://spectra.video/.well-known/webfinger?resource=acct%3Amberos%40spectra.video":
			// Webfinger.
			responseCode = http.StatusOK
			responseBytes = []byte(peertubeWebfingerJSON)
			responseContentType = "application/json; charset=utf-8"

		case "https://spectra.video/accounts/mberos":
			// Account.
			responseCode = http.StatusOK
			responseBytes = []byte(peertubePersonJSON)
			responseContentType = "application/activity+json"

		default:
			// Unknown.
			responseCode = http.StatusNotFound
			responseBytes = []byte(`{"error":"404 not found"}`)
			responseContentType = "application/activity+json"
		}

		reader := bytes.NewReader(responseBytes)
		readCloser := io.NopCloser(reader)
		return &http.Response{
			Request:       req,
			StatusCode:    responseCode,
			Body:          readCloser,
			ContentLength: int64(len(responseBytes)),
			Header: http.Header{
				"Content-Type": {responseContentType},
			},
		}, nil
	}
	testStructs.HTTPClient = testrig.NewMockHTTPClient(do, rMediaPath)

	// Update dereferencer to use new test HTTP client.
	testStructs.Federator.Dereferencer = dereferencing.NewDereferencer(
		testStructs.State,
		testStructs.TypeConverter,
		testrig.NewTestTransportController(testStructs.State, testStructs.HTTPClient),
		visibility.NewFilter(testStructs.State),
		interaction.NewFilter(testStructs.State),
		relay.NewFilter(testStructs.State),
		media.NewManager(testStructs.State),
	)

	// The local account we will be fetching statuses as.
	fetchingAccount := suite.testAccounts["local_account_1"]

	status, _, err := testStructs.Federator.Dereferencer.GetStatusByURI(ctx,
		fetchingAccount,
		testrig.URLMustParse("https://spectra.video/videos/watch/4e06b13a-b7bb-4530-a653-c82d47a7bfb8"),
	)
	if err != nil {
		suite.FailNow(err.Error())
	}
	
	suite.Equal("https://spectra.video/videos/watch/4e06b13a-b7bb-4530-a653-c82d47a7bfb8", status.URI)
	suite.Equal("https://spectra.video/w/aCQ9f8rMn14CW8jBzCTiA9", status.URL)
	suite.Equal("https://spectra.video/accounts/mberos", status.Account.URI)
	suite.Equal("https://spectra.video/a/mberos/video-channels", status.Account.URL)
}

func TestStatusPeertubeTestSuite(t *testing.T) {
	suite.Run(t, new(StatusPeertubeTestSuite))
}
