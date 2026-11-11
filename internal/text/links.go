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
	"net/url"
	"strings"

	"code.superseriousbusiness.org/gotosocial/internal/gtserror"
	"code.superseriousbusiness.org/gotosocial/internal/gtsmodel"
	"code.superseriousbusiness.org/gotosocial/internal/util"
	"github.com/yuin/goldmark/ast"
	"golang.org/x/net/html"
)

// linkFromNode checks the given HTML node to see if it's an anchor
// element with a valid href that can be extracted. Will return nil
// if not, or nil and an error if the discovered link is malformed.
func linkFromNode(node *html.Node) (*gtsmodel.Link, error) {
	if node.Type != html.ElementNode || node.Data != "a" {
		// Not an <a> element.
		return nil, nil
	}

	// Check each HTML attribute of
	// the <a> element to find href.
	var href *url.URL
	for _, attr := range node.Attr {
		switch attr.Key {
		case "class":
			classes := strings.Split(attr.Val, " ")
			for _, class := range classes {
				if class == "mention" || class == "hashtag" {
					// Not relevant, stop
					// looking at this node.
					return nil, nil
				}
			}

		case "href":
			if attr.Val == "" {
				return nil, gtserror.New("empty href")
			}

			u, err := url.Parse(attr.Val)
			if err != nil {
				return nil, gtserror.Newf("malformed link: %w", err)
			}

			if !u.IsAbs() || u.Host == "" {
				// We don't want schemaless
				// hostless / relative links.
				return nil, nil
			}

			// Found a link! But allow attribute
			// loop to continue in case we discover
			// it's a mention that should be ignored.
			href = u
		}
	}

	if href == nil {
		// There was no href set on this
		// node so it's just a weird anchor.
		return nil, nil
	}

	// Punify the URL
	// before returning.
	var err error
	href.Host, err = util.Punify(href.Host)
	if err != nil {
		return nil, gtserror.Newf("error punifying %s: %w", href.Host, err)
	}

	return &gtsmodel.Link{
		URL:       href.String(),
		URLParsed: href,
	}, nil
}

// ParseLinksFromContent extracts punified links from the hrefs of all
// non-mention/hashtag, absolute (with scheme) anchor links in the content.
//
// Because it involves parsing content to HTML and recursively looking
// through the DOM node by node, this is not a cheap function, and
// should not be called multiple times per piece of text.
//
// Note that even if this function returns an error, it might still
// return some links anyway, so check the return values carefully.
func ParseLinksFromContent(htmlStr string) ([]*gtsmodel.Link, error) {
	// Parse input as HTML.
	node, err := html.Parse(strings.NewReader(htmlStr))
	if err != nil {
		return nil, gtserror.Newf("invalid html input: %w", err)
	}

	// Find each HTML anchor node that's not
	// a mention, and take the href of it.
	var links []*gtsmodel.Link
	var errs gtserror.MultiError
	var f func(*html.Node)
	f = func(node *html.Node) {

		// Try to extract
		// link from this node.
		link, err := linkFromNode(node)
		if err != nil {
			errs.Append(err)
		} else if link != nil {
			// Valid non-mention link, append it.
			links = append(links, link)
		}

		// Tail recurse on the next node.
		for c := node.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(node)

	return links, errs.Combine()
}

// parseLinksFromDoc walks the given document with the
// given source to extract links from it where possible.
//
// Note that even if this function returns an error, it might still
// return some links anyway, so check the return values carefully.
func parseLinksFromDoc(doc ast.Node, source []byte) ([]*gtsmodel.Link, error) {
	// Extract all links
	// from the parsed document.
	var links []*gtsmodel.Link
	var errs gtserror.MultiError
	if err := ast.Walk(doc, func(node ast.Node, enter bool) (ws ast.WalkStatus, _ error) {
		// We're gonna keep
		// walking no matter
		// what so just set here.
		ws = ast.WalkContinue

		if !enter {
			// Not interested,
			// keep walking.
			return
		}

		var hrefStr string
		switch n := node.(type) {

		// Raw HTML style:
		// `<a href="https://example.org>link</a>"`
		case *ast.RawHTML:
			raw := string(n.Segments.Value(source))
			if !strings.HasPrefix(raw, "<a") {
				// Doesn't look like
				// an anchor element,
				// don't try parsing,
				// just keep walking.
				return
			}

			// ParseLinksFromContent does
			// parsing + punification for us.
			parsedLinks, err := ParseLinksFromContent(raw)
			if err != nil {
				// Store any error but don't return yet
				// as we might have got some links anyway.
				err := gtserror.Newf("error getting link from *ast.RawHTML: %w", err)
				errs.Append(err)
			}

			if len(parsedLinks) != 0 {
				links = append(links, parsedLinks...)
			}

			// Keep walking.
			return

		// Markdown link style:
		// `[link](https://example.org)`
		case *ast.Link:
			hrefStr = string(n.Destination)

		// Autolink style:
		// `https://example.org`
		case *ast.AutoLink:
			hrefStr = string(n.URL(source))
		}

		if hrefStr == "" {
			// Nothing here or
			// not a relevant
			// node type to us.
			// Keep walking.
			return
		}

		// We had an href string from a link
		// or autolink, parse + punify it.
		href, err := url.Parse(hrefStr)
		if err != nil {
			// Store err, keep walking.
			err := gtserror.Newf("malformed link: %w", err)
			errs.Append(err)
			return
		}

		if !href.IsAbs() || href.Host == "" {
			// We don't want schemaless
			// hostless / relative links.
			// Keep walking.
			return
		}

		href.Host, err = util.Punify(href.Host)
		if err != nil {
			// Store err, keep walking.
			err := gtserror.Newf("error punifying %s: %w", href.Host, err)
			errs.Append(err)
			return
		}

		// We got one captain.
		links = append(links, &gtsmodel.Link{
			URL:       href.String(),
			URLParsed: href,
		})

		return
	}); err != nil {
		return nil, err
	}

	return links, errs.Combine()
}
