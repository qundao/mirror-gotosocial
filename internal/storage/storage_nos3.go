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

//go:build nos3

package storage

import (
	"context"
	"errors"
	"os"

	"code.superseriousbusiness.org/gotosocial/internal/gtserror"
	"code.superseriousbusiness.org/gotosocial/internal/log"
	"codeberg.org/gruf/go-storage"
)

// Driver wraps a disk or memory storage.Storage
// to provide optimized write operations.
type Driver struct{ Storage storage.Storage }

// PutFile: see PutFile() in storage.go.
func (d *Driver) PutFile(ctx context.Context, key, filepath, _ string) (int64, error) {

	// Open file at path for reading.
	file, err := os.Open(filepath)
	if err != nil {
		return 0, gtserror.Newf("error opening file %s: %w", filepath, err)
	}

	// Write the file data to storage under key. Note
	// that for disk.DiskStorage{} this should end up
	// being a highly optimized Linux sendfile syscall.
	sz, err := d.Storage.WriteStream(ctx, key, file)

	// Wrap write error.
	if err != nil {
		err = gtserror.Newf("error writing file %s: %w", key, err)
	}

	// Close the file: done with it.
	if e := file.Close(); e != nil {
		log.Errorf(ctx, "error closing file %s: %v", filepath, e)
	}

	return sz, err
}

// URL: not implemented for 'nos3'.
func (d *Driver) URL(ctx context.Context, key string) *PresignedURL {
	return nil
}

// ProbeCSPUri: not implemented for 'nos3'.
func (d *Driver) ProbeCSPUri(ctx context.Context) (string, error) {
	return "", nil
}

func NewS3Storage() (*Driver, error) {
	return nil, errors.New("gotosocial was compiled without S3 storage support")
}
