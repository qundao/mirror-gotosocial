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

package storage

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/url"
	"time"

	"code.superseriousbusiness.org/gotosocial/internal/config"
	"codeberg.org/gruf/go-bytesize"
	"codeberg.org/gruf/go-fastcopy"
	"codeberg.org/gruf/go-storage"
	"codeberg.org/gruf/go-storage/disk"
)

// PresignedURL represents a pre signed S3 URL with
// an expiry time.
type PresignedURL struct {
	*url.URL
	Expiry time.Time // link expires at this time
}

// IsInvalidKey returns whether error is an invalid-key
// type error returned by the underlying storage library.
func IsInvalidKey(err error) bool {
	return errors.Is(err, storage.ErrInvalidKey)
}

// IsAlreadyExist returns whether error is an already-exists
// type error returned by the underlying storage library.
func IsAlreadyExist(err error) bool {
	return errors.Is(err, storage.ErrAlreadyExists)
}

// IsNotFound returns whether error is a not-found error
// type returned by the underlying storage library.
func IsNotFound(err error) bool {
	return errors.Is(err, storage.ErrNotFound)
}

// Get returns the byte value for key in storage.
func (d *Driver) Get(ctx context.Context, key string) ([]byte, error) {
	return d.Storage.ReadBytes(ctx, key)
}

// GetStream returns an io.ReadCloser for the value bytes at key in the storage.
func (d *Driver) GetStream(ctx context.Context, key string) (io.ReadCloser, error) {
	return d.Storage.ReadStream(ctx, key)
}

// Put writes the supplied value bytes at key in the storage
func (d *Driver) Put(ctx context.Context, key string, value []byte) (int, error) {
	return d.Storage.WriteBytes(ctx, key, value)
}

// Delete attempts to remove the supplied key (and corresponding value) from storage.
func (d *Driver) Delete(ctx context.Context, key string) error {
	return d.Storage.Remove(ctx, key)
}

// Has checks if the supplied key is in the storage.
func (d *Driver) Has(ctx context.Context, key string) (bool, error) {
	stat, err := d.Storage.Stat(ctx, key)
	return (stat != nil), err
}

// WalkKeys walks the keys in the storage.
func (d *Driver) WalkKeys(ctx context.Context, walk func(string) error) error {
	return d.Storage.WalkKeys(ctx, storage.WalkKeysOpts{
		Step: func(entry storage.Entry) error {
			return walk(entry.Key)
		},
	})
}

func AutoConfig() (*Driver, error) {
	switch backend := config.GetStorageBackend(); backend {
	case "s3":
		return NewS3Storage()
	case "local":
		return NewFileStorage()
	default:
		return nil, fmt.Errorf("invalid storage backend: %s", backend)
	}
}

func NewFileStorage() (*Driver, error) {
	// Load runtime configuration
	basePath := config.GetStorageLocalBasePath()

	// Update fastcopy global buffer pool
	// to use our requested buffer size.
	const bufsize = 16 * bytesize.KiB
	fastcopy.Buffer(int(bufsize))

	// Use default disk config with
	// increased write buffer size.
	diskCfg := disk.DefaultConfig()
	diskCfg.CopyFn = fastcopy.Copy

	// Open the disk storage implementation
	disk, err := disk.Open(basePath, &diskCfg)
	if err != nil {
		return nil, fmt.Errorf("error opening disk storage: %w", err)
	}

	return &Driver{Storage: disk}, nil
}
