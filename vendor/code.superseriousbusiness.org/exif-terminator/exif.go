/*
   exif-terminator
   Copyright (C) 2022 SuperSeriousBusiness admin@gotosocial.org

   This program is free software: you can redistribute it and/or modify
   it under the terms of the GNU Affero General Public License as published by
   the Free Software Foundation, either version 3 of the License, or
   (at your option) any later version.

   This program is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
   GNU Affero General Public License for more details.

   You should have received a copy of the GNU Affero General Public License
   along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/

package terminator

import (
	"errors"
	"fmt"
	"io"
	"strings"

	exif "github.com/dsoprea/go-exif/v3"
	exifcommon "github.com/dsoprea/go-exif/v3/common"
)

type withEXIF interface {
	Exif() (rootIfd *exif.Ifd, data []byte, err error)
	SetExif(ib *exif.IfdBuilder) (err error)
}

func terminateEXIF(data withEXIF) (err error) {
	var ifd *exif.Ifd

	// Read EXIF from data chunk.
	ifd, _, err = data.Exif()
	switch {
	case err == nil:

	// Catch and replace no exif errors with nil return.
	case strings.Contains(err.Error(), "no exif data"):
		return nil

	// Weird exif chunk at
	// end of file, ignore.
	case errors.Is(err, io.EOF):
		return nil

	default:
		return err
	}

	// Exif builder to (over)write
	// new exif data into the file.
	var ifdb *exif.IfdBuilder

	// Try to get IB chain from EXIF IFD chain (can be nil).
	ifdb, err = newIfdBuilderFromExistingChain(ifd)
	switch {
	case err == nil:
		// No issue.

	case ifd != nil:
		// Parsing issue or something.
		// Fall back to basic ifd builder.
		ifdb = exif.NewIfdBuilderWithExistingIfd(ifd)

	default:
		// Parsing issue and no
		// ifd available to build
		// from, error out.
		return err
	}

	if ifd != nil {
		// Search for existing orientation data in EXIF chunk.
		orientation, _ := ifdb.FindTagWithName("Orientation")

		// Start new IFD chain from fesh mapping and indices.
		im, ti := exifcommon.NewIfdMapping(), exif.NewTagIndex()
		ifdb = exif.NewIfdBuilder(im, ti, ifd.IfdIdentity(), ifd.ByteOrder())

		if orientation != nil {
			// Set old orientation.
			ifdb.Add(orientation)
		}
	}

	// Set new empty IFD
	// chain on EXIF chunk.
	err = data.SetExif(ifdb)
	return
}

// newIfdBuilderFromExistingChain wraps exif.NewIfdBuilderFromExistingChain(), recovering uncaught panics.
func newIfdBuilderFromExistingChain(ifd *exif.Ifd) (ifdb *exif.IfdBuilder, err error) {
	defer func() {
		switch r := recover().(type) {
		case nil:
		case error:
			err = r
		default:
			err = fmt.Errorf("recovered panic: %v", r)
		}
	}()
	ifdb = exif.NewIfdBuilderFromExistingChain(ifd)
	return
}
