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
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"

	"code.superseriousbusiness.org/activity/pub"
	"code.superseriousbusiness.org/gotosocial/internal/ap"
	apiutil "code.superseriousbusiness.org/gotosocial/internal/api/util"
	"code.superseriousbusiness.org/gotosocial/internal/config"
	"code.superseriousbusiness.org/gotosocial/internal/db"
	"code.superseriousbusiness.org/gotosocial/internal/federation/federatingdb"
	"code.superseriousbusiness.org/gotosocial/internal/gtserror"
	"code.superseriousbusiness.org/gotosocial/internal/gtsmodel"
	"code.superseriousbusiness.org/gotosocial/internal/state"
	"code.superseriousbusiness.org/gotosocial/internal/util"
	"codeberg.org/gruf/go-byteutil"
	"codeberg.org/gruf/go-cache/v3/simple"
)

type Controller struct {
	state     *state.State
	fedDB     *federatingdb.DB
	client    pub.HttpClient
	cache     *simple.Cache[string, *Transport] // transport cache
	instance  *Transport                        // single instance-account transport
	userAgent string
}

func NewController(ctx context.Context, state *state.State, federatingDB *federatingdb.DB, client pub.HttpClient) *Controller {
	c := &Controller{
		state:    state,
		fedDB:    federatingDB,
		client:   client,
		cache:    simple.New[string, *Transport](0, 100),
		instance: nil,
		userAgent: fmt.Sprintf("gotosocial/%s (+%s://%s)",
			config.GetSoftwareVersion(),
			config.GetProtocol(),
			config.GetHost(),
		),
	}

	// FIXME: a hacky fix for the testrig database not containing the instance
	// account at the time of Controller{} instantiation. truly a heinous code crime.
	instance, _ := ctx.Value("testrig-transport-controller-key").(*gtsmodel.Account) //nolint
	if instance == nil {
		var err error

		// Get account associated with this local instance.
		instance, err = state.DB.GetInstanceAccount(ctx, "")
		if err != nil {
			panic(gtserror.Newf("error getting instance account: %w", err))
		}
	}

	// Pre-prepare the instance account transport.
	c.instance = c.GetTransportForAccount(instance)

	return c
}

func (c *Controller) NewTransport(pubKeyID string, privkey *rsa.PrivateKey) (*Transport, error) {
	return c.GetTransport(pubKeyID, privkey), nil
}

func (c *Controller) GetTransport(pubKeyID string, privkey *rsa.PrivateKey) *Transport {
	// Generate public key string for cache key
	//
	// NOTE: it is safe to use the public key as the cache
	// key here as we are generating it ourselves from the
	// private key. If we were simply using a public key
	// provided as argument that would absolutely NOT be safe.
	pubStr := privkeyToPublicStr(privkey)

	// First check for cached transport.
	transp, ok := c.cache.Get(pubStr)
	if ok {
		return transp
	}

	// Create transport.
	transp = &Transport{
		controller: c,
		pubKeyID:   pubKeyID,
		privkey:    privkey,
	}

	// Cache transport under pubkey
	if !c.cache.Add(pubStr, transp) {
		var cached *Transport

		cached, ok = c.cache.Get(pubStr)
		if !ok {
			// Some ridiculous race cond.
			c.cache.Set(pubStr, transp)
		} else {
			// Use already cached
			transp = cached
		}
	}

	return transp
}

// GetTransportForAccount will return the Transport{} instance for given local account.
// A nil account model will instead return a Transport{} for the local instance account.
func (c *Controller) GetTransportForAccount(account *gtsmodel.Account) *Transport {
	if account == nil {

		// nil account indicates this should
		// go through the instance account
		// transport, return pre-prepared.
		return c.instance
	}

	// Ensure they've provided us a valid local account.
	if !account.IsLocal() || account.PrivateKey == nil {
		panic(gtserror.New("cannot create transport for account"))
	}

	// Return a transport for local account's private key + its URI.
	return c.GetTransport(account.PublicKeyURI, account.PrivateKey)
}

// dereferenceLocal is a shortcut to try dereferencing
// something on this instance without making any http calls.
//
// Will return an error if nothing could be found, indicating that
// the calling transport should continue with an http call anyway.
//
// It should only be invoked when the iri.Host == this host.
func (c *Controller) dereferenceLocal(
	ctx context.Context,
	uri *url.URL,
) (*http.Response, error) {

	// Try fetch via federating DB.
	t, err := c.fedDB.Get(ctx, uri)

	switch {
	// No problem.
	case err == nil:

	// Catch and handle objects not found.
	case errors.Is(err, db.ErrNoEntries):
		return &http.Response{
			Request:    &http.Request{URL: uri},
			Status:     http.StatusText(http.StatusNotFound),
			StatusCode: http.StatusNotFound,
			Header: map[string][]string{
				"Content-Type":   {apiutil.AppActivityLDJSON},
				"Content-Length": {"0"},
			},
		}, nil

	// Any other.
	default:
		return nil, gtserror.Newf("error getting: %w", err)
	}

	if util.IsNil(t) {
		// Assert this should never happen.
		panic(gtserror.New("nil vocab.Type"))
	}

	// Serialize type to JSON map.
	m, err := ap.Serialize(t)
	if err != nil {
		return nil, err
	}

	// Marshal JSON to bytes.
	b, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}

	// Return a response
	// with AS data as body.
	contentLength := len(b)
	rsp := &http.Response{
		Request:       &http.Request{URL: uri},
		Status:        http.StatusText(http.StatusOK),
		StatusCode:    http.StatusOK,
		Body:          io.NopCloser(bytes.NewReader(b)),
		ContentLength: int64(contentLength),
		Header: map[string][]string{
			"Content-Type":   {apiutil.AppActivityLDJSON},
			"Content-Length": {strconv.Itoa(contentLength)},
		},
	}

	return rsp, nil
}

// privkeyToPublicStr will create a string representation of RSA public key from private.
func privkeyToPublicStr(privkey *rsa.PrivateKey) string {
	b := x509.MarshalPKCS1PublicKey(&privkey.PublicKey)
	return byteutil.B2S(b)
}
