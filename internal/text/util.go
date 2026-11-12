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

package text

import (
	"context"
	"unicode"

	"code.superseriousbusiness.org/gopkg/buffers"
	"code.superseriousbusiness.org/gopkg/log"
	"code.superseriousbusiness.org/gotosocial/internal/config"
	"code.superseriousbusiness.org/gotosocial/internal/gtsmodel"
	"code.superseriousbusiness.org/gotosocial/internal/language"
	"code.superseriousbusiness.org/gotosocial/internal/regexes"
	"codeberg.org/gruf/go-byteutil"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/renderer/html"
)

// bufpool is a memory pool
// of buffers for conversions.
var bufpool = buffers.Pool(2048)

// isPermittedInHashtag is true for characters that
// may be permitted anywhere in hashtag, entirety or not.
func isPermittedInHashtag(r rune) bool {
	return unicode.IsLetter(r) || isPermittedIfNotEntireHashtag(r)
}

// isPermittedIfNotEntireHashtag is true for characters that may be in a hashtag
// but are not allowed to be the only characters making up the hashtag.
func isPermittedIfNotEntireHashtag(r rune) bool {
	return unicode.IsNumber(r) || unicode.IsMark(r) || r == '_'
}

// isHashtagBoundary returns true if rune r
// is a recognized break character for before
// or after a #hashtag.
func isHashtagBoundary(r rune) bool {
	switch {

	// Zero width space.
	case r == '\u200B':
		return true

	// Zero width no-break space.
	case r == '\uFEFF':
		return true

	// Pipe character sometimes
	// used as workaround.
	case r == '|':
		return true

	// Standard Unicode white space.
	case unicode.IsSpace(r):
		return true

	// Non-underscore punctuation.
	case unicode.IsPunct(r) && r != '_':
		return true

	// Not recognized
	// hashtag boundary.
	default:
		return false
	}
}

// isMentionBoundary returns true if rune r
// is a recognized break character for before
// or after a @mention.
func isMentionBoundary(r rune) bool {
	return unicode.IsSpace(r) ||
		unicode.IsPunct(r)
}

// ContentToContentLanguage tries to
// extract a content string and language
// tag string from the given intermediary
// content.
//
// Either/both of the returned strings may
// be empty, depending on how things go.
//
// If content.MediaType = "text/markdown",
// the content will be converted to HTML
// using a basic parser, before return.
func ContentToContentLanguage(
	ctx context.Context,
	content gtsmodel.Content,
) (
	contentStr string, // content
	langTagStr string, // language
) {
	switch contentMap := content.ContentMap; {
	// Simplest case: no `contentMap`.
	// Take `content`, even if empty.
	case contentMap == nil:
		contentStr = content.Content

	// `content` and `contentMap` set.
	// Try to infer "primary" language.
	case content.Content != "":
		// Assume `content` is intended
		// primary content, and look for
		// corresponding language tag.
		contentStr = content.Content

		for t, c := range contentMap {
			if contentStr == c {
				langTagStr = t
				break
			}
		}

	// `content` not set; `contentMap`
	// is set with only one value.
	// This must be the "primary" lang.
	case len(contentMap) == 1:
		// Use an empty loop to
		// get the values we want.
		// nolint:revive
		for langTagStr, contentStr = range contentMap {
		}

	// Only `contentMap` is set, with more
	// than one value. Map order is not
	// guaranteed so we can't know the
	// "primary" language.
	//
	// Try to select content using our
	// instance's configured languages.
	//
	// In case of no hits, just take the
	// first tag and content in the map.
	default:
		instanceLangs := config.GetInstanceLanguages()
		for _, langTagStr = range instanceLangs.TagStrs() {
			if contentStr = contentMap[langTagStr]; contentStr != "" {
				// Hit!
				break
			}
		}

		// If nothing found, just take
		// the first entry we can get by
		// breaking after the first iter.
		if contentStr == "" {
			for langTagStr, contentStr = range contentMap {
				break
			}
		}
	}

	if langTagStr != "" {
		// Found a lang tag for this content,
		// make sure it's valid / parseable.
		lang, err := language.Parse(langTagStr)
		if err != nil {
			log.Warnf(
				ctx,
				"could not parse %s as BCP47 language tag in status contentMap: %v",
				langTagStr, err,
			)
		} else {
			// Inferred the language!
			// Use normalized version.
			langTagStr = lang.TagStr
		}
	}

	// If content type is text/markdown instead of
	// implicit HTML, convert it before returning.
	if contentStr != "" && content.MediaType == "text/markdown" {
		// Take a very basic
		// renderer for this.
		md := goldmark.New(
			goldmark.WithRendererOptions(
				html.WithXHTML(),
				html.WithHardWraps(),
			),
			goldmark.WithExtensions(
				// Allow turning
				// URLs into links.
				extension.NewLinkify(
					extension.WithLinkifyURLRegexp(regexes.URLLike),
				),
				extension.Strikethrough,
			),
		)

		// Convert input string to bytes.
		source := byteutil.S2B(contentStr)

		// Get byte buffer.
		buf := bufpool.Get()

		// Convert markdown to HTML.
		if err := md.Convert(source, buf); err != nil {
			log.Errorf(ctx, "error formatting markdown input to HTML: %v", err)
			buf.B = buf.B[:0] // reset buffer
		} else {
			// Copy data.
			contentStr = string(buf.B)

			// Clean and shrink HTML.
			contentStr = SanitizeHTML(contentStr)
			contentStr = MinifyHTML(contentStr)
		}

		// Release buffer.
		bufpool.Put(buf)
	}

	return contentStr, langTagStr
}
