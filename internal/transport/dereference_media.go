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

package transport

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/url"

	"code.superseriousbusiness.org/gotosocial/internal/gtserror"
	"codeberg.org/gruf/go-iotools"
)

// DereferenceMedia fetches the given media attachment IRI, returning the reader limited to given max.
func (t *Transport) DereferenceMedia(ctx context.Context, iri *url.URL, maxsz int64) (io.ReadCloser, error) {
	if maxsz <= 0 {
		// Max size is zero, just return.
		return emptyLimitedReader(), nil
	}

	// Build IRI just once.
	iriStr := iri.String()

	// Prepare HTTP request to this media's IRI.
	req, err := http.NewRequestWithContext(ctx,
		"GET",
		iriStr,
		nil,
	)
	if err != nil {
		return nil, err
	}

	// We don't know what kind of
	// media we're going to get here.
	req.Header.Add("Accept", "*/*")

	// Perform the HTTP request.
	rsp, err := t.GET(req)
	if err != nil {
		return nil, err
	}

	// Check for an expected status code.
	if rsp.StatusCode != http.StatusOK {
		return nil, gtserror.NewFromResponse(rsp)
	}

	// Check media within size limit.
	if rsp.ContentLength > maxsz {
		_ = rsp.Body.Close() // close early.
		return emptyLimitedReader(), nil
	}

	// Update response body with maximum supported media size.
	rsp.Body, _, _ = iotools.UpdateReadCloserLimit(rsp.Body, maxsz)

	return rsp.Body, nil
}

var newline = []byte{'\n'}

func noop() error { return nil }

// emptyLimitReader returns an io.ReadCloser that itself
// is wrapped in a limited reader with zero length left
// in the read. Ensuring media error details passed along.
func emptyLimitedReader() io.ReadCloser {
	r := (io.Reader)(bytes.NewReader(newline[:])) //nolint
	r = &io.LimitedReader{R: r, N: 0}
	return &iotools.ReadCloserType{
		Reader: r,
		Closer: iotools.CloserFunc(noop),
	}
}
