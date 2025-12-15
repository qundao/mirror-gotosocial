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

package util

// See:
//
//   - https://developers.google.com/search/docs/crawling-indexing/robots-meta-tag#robotsmeta
//   - https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/X-Robots-Tag
//   - https://www.rfc-editor.org/rfc/rfc9309.html
const (
	RobotsDirectivesDisallow  = "noindex, nofollow"
	RobotsDirectivesAllowSome = "nofollow, noarchive, nositelinkssearchbox, max-image-preview:standard"
	RobotsTxt                 = `# GoToSocial robots.txt -- to edit, see internal/api/util/robots.go
# More info @ https://developers.google.com/search/docs/crawling-indexing/robots/intro

# AI scrapers and the like.
# https://github.com/ai-robots-txt/ai.robots.txt/
User-agent: AddSearchBot
User-agent: AI2Bot
User-agent: AI2Bot-DeepResearchEval
User-agent: Ai2Bot-Dolma
User-agent: aiHitBot
User-agent: amazon-kendra
User-agent: Amazonbot
User-agent: AmazonBuyForMe
User-agent: Andibot
User-agent: Anomura
User-agent: anthropic-ai
User-agent: Applebot
User-agent: Applebot-Extended
User-agent: atlassian-bot
User-agent: Awario
User-agent: bedrockbot
User-agent: bigsur.ai
User-agent: Bravebot
User-agent: Brightbot 1.0
User-agent: BuddyBot
User-agent: Bytespider
User-agent: CCBot
User-agent: Channel3Bot
User-agent: ChatGLM-Spider
User-agent: ChatGPT Agent
User-agent: ChatGPT-User
User-agent: Claude-SearchBot
User-agent: Claude-User
User-agent: Claude-Web
User-agent: ClaudeBot
User-agent: Cloudflare-AutoRAG
User-agent: CloudVertexBot
User-agent: cohere-ai
User-agent: cohere-training-data-crawler
User-agent: Cotoyogi
User-agent: Crawl4AI
User-agent: Crawlspace
User-agent: Datenbank Crawler
User-agent: DeepSeekBot
User-agent: Devin
User-agent: Diffbot
User-agent: DuckAssistBot
User-agent: Echobot Bot
User-agent: EchoboxBot
User-agent: FacebookBot
User-agent: facebookexternalhit
User-agent: Factset_spyderbot
User-agent: FirecrawlAgent
User-agent: FriendlyCrawler
User-agent: Gemini-Deep-Research
User-agent: Google-CloudVertexBot
User-agent: Google-Extended
User-agent: Google-Firebase
User-agent: Google-NotebookLM
User-agent: GoogleAgent-Mariner
User-agent: GoogleOther
User-agent: GoogleOther-Image
User-agent: GoogleOther-Video
User-agent: GPTBot
User-agent: iAskBot
User-agent: iaskspider
User-agent: iaskspider/2.0
User-agent: IbouBot
User-agent: ICC-Crawler
User-agent: ImagesiftBot
User-agent: imageSpider
User-agent: img2dataset
User-agent: ISSCyberRiskCrawler
User-agent: Kangaroo Bot
User-agent: KlaviyoAIBot
User-agent: KunatoCrawler
User-agent: laion-huggingface-processor
User-agent: LAIONDownloader
User-agent: LCC
User-agent: LinerBot
User-agent: Linguee Bot
User-agent: LinkupBot
User-agent: Manus-User
User-agent: meta-externalagent
User-agent: Meta-ExternalAgent
User-agent: meta-externalfetcher
User-agent: Meta-ExternalFetcher
User-agent: meta-webindexer
User-agent: MistralAI-User
User-agent: MistralAI-User/1.0
User-agent: MyCentralAIScraperBot
User-agent: netEstate Imprint Crawler
User-agent: NotebookLM
User-agent: NovaAct
User-agent: OAI-SearchBot
User-agent: omgili
User-agent: omgilibot
User-agent: OpenAI
User-agent: Operator
User-agent: PanguBot
User-agent: Panscient
User-agent: panscient.com
User-agent: Perplexity-User
User-agent: PerplexityBot
User-agent: PetalBot
User-agent: PhindBot
User-agent: Poggio-Citations
User-agent: Poseidon Research Crawler
User-agent: QualifiedBot
User-agent: QuillBot
User-agent: quillbot.com
User-agent: SBIntuitionsBot
User-agent: Scrapy
User-agent: SemrushBot-OCOB
User-agent: SemrushBot-SWA
User-agent: ShapBot
User-agent: Sidetrade indexer bot
User-agent: Spider
User-agent: TerraCotta
User-agent: Thinkbot
User-agent: TikTokSpider
User-agent: Timpibot
User-agent: TwinAgent
User-agent: VelenPublicWebCrawler
User-agent: WARDBot
User-agent: Webzio-Extended
User-agent: webzio-extended
User-agent: wpbot
User-agent: WRTNBot
User-agent: YaK
User-agent: YandexAdditional
User-agent: YandexAdditionalBot
User-agent: YouBot
User-agent: ZanistaBot
Disallow: /

# Marketing/SEO "intelligence" data scrapers
User-agent: AwarioRssBot
User-agent: AwarioSmartBot
User-agent: DataForSeoBot
User-agent: magpie-crawler
User-agent: Meltwater
User-agent: peer39_crawler
User-agent: peer39_crawler/1.0
User-agent: PiplBot
User-agent: scoop.it
User-agent: Seekr
Disallow: /

# Well-known.dev crawler. Indexes stuff under /.well-known.
# https://well-known.dev/about/
User-agent: WellKnownBot
Disallow: /

# Rules for everything else.
User-agent: *
Crawl-delay: 500

# API endpoints.
Disallow: /api/

# Auth/Sign in endpoints.
Disallow: /auth/
Disallow: /oauth/
Disallow: /check_your_email
Disallow: /wait_for_approval
Disallow: /account_disabled
Disallow: /signup

# Fileserver/media.
Disallow: /fileserver/

# Fedi S2S API endpoints.
Disallow: /users/
Disallow: /emoji/

# Settings panels.
Disallow: /admin
Disallow: /user
Disallow: /settings/

# Domain blocklist.
Disallow: /about/suspended

# Webfinger endpoint.
Disallow: /.well-known/webfinger
`
	RobotsTxtDisallowNodeInfo = RobotsTxt + `
# Disallow nodeinfo
Disallow: /.well-known/nodeinfo
Disallow: /nodeinfo/
`

	// MD5 hash of basic robots.txt.
	RobotsTxtETag = `7b6b498f7381ac33cb3efb34c68f662d`
	// MD5 hash of robots.txt with NodeInfo disallowed.
	RobotsTxtDisallowNodeInfoETag = `6d21be573d502581a3bf7271b7e63fc8`
)
