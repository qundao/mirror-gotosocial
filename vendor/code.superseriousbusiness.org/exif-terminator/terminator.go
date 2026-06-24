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
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"

	jpegstructure "code.superseriousbusiness.org/go-jpeg-image-structure/v2"
	pngstructure "code.superseriousbusiness.org/go-png-image-structure/v2"
)

// ErrUnsupported is returned when a media-type
// passed to Terminate{,Into}() is not supported.
var ErrUnsupported = errors.New("unsupported media type")

// ErrOnExifTermination is returned as a "prefix" error, wrapping error(s)
// that occurred during exif termination. These are deferred and returned
// at the end, so exif data ahead of encountered error may still be cleared.
var ErrOnExifTermination = errors.New("error(s) during exif termination")

// Terminate will attempt to strip EXIF from image of 'mediaType' contained
// in reader, returning a reader that streams the resulting cleaned image.
//
// NOTE: You should prefer using TerminateInto() which is more performant.
func Terminate(in io.Reader, mediaType string) (io.Reader, error) {

	// To avoid keeping too much stuff
	// in memory we want to pipe data
	// directly to the reader.
	pr, pw := io.Pipe()

	// Setup scanner to terminate exif data into our pipe writer.
	scanner, deferred, err := terminatingScanner(pw, in, mediaType)
	if err != nil {
		_ = pw.Close()
		return nil, err
	}

	go func() {
		var err error

		defer func() {
			// Try recover any error, since
			// disproprea libraries do much
			// error handling via panicking :|
			switch r := recover().(type) {
			case nil:
			case error:
				err = r
			default:
				err = fmt.Errorf("recovered panic: %v", r)
			}

			// Always close writer, using returned
			// scanner error (if any). If err is nil
			// then the standard io.EOF will be used.
			// (this will not overwrite existing).
			pw.CloseWithError(err)
		}()

		// Scan through input.
		for scanner.Scan() {
		}

		// Set error on return.
		err = scanner.Err()
		if err != nil {
			return
		}

		// Else check deferred exif termination errs.
		if deferred != nil && len(*deferred) > 0 {
			err = fmt.Errorf("%w: %v", ErrOnExifTermination, errors.Join(*deferred...))
		}
	}()

	return pr, nil
}

// TerminateInto will attempt to strip EXIF from image of 'mediaType' contained in reader, writing to output.
func TerminateInto(out io.Writer, in io.Reader, mediaType string) error {

	// Setup scanner to terminate exif data from 'in' to 'out'.
	scanner, deferred, err := terminatingScanner(out, in, mediaType)
	if err != nil {
		return err
	}

	defer func() {
		// Try recover any error, since
		// disproprea libraries do much
		// error handling via panicking :|
		switch r := recover().(type) {
		case nil:
		case error:
			err = r
		default:
			err = fmt.Errorf("recovered panic: %v", r)
		}
	}()

	// Scan through input.
	for scanner.Scan() {
	}

	// Check and return any scan errors.
	if err := scanner.Err(); err != nil {
		return err
	}

	// Else check deferred exif termination errs.
	if deferred != nil && len(*deferred) > 0 {
		return fmt.Errorf("%w: %v", ErrOnExifTermination, errors.Join(*deferred...))
	}

	return nil
}

func terminatingScanner(out io.Writer, in io.Reader, mediaType string) (scanner *bufio.Scanner, deferred *[]error, err error) {
	scanner = bufio.NewScanner(in)

	// 40mb buffer size should be enough
	// to scan through most file chunks
	// without running into issues, they're
	// usually chunked smaller than this...
	scanner.Buffer(nil, 40*1024*1024)

	switch mediaType {
	case "image/jpeg", "jpeg", "jpg":
		v := &jpegVisitor{write: out}
		deferred = &v.errs

		// Provide the visitor to the splitter so
		// that it triggers on every section scan.
		js := jpegstructure.NewJpegSplitter(v)

		// The visitor also needs to read back the
		// list of segments: for this it needs to
		// know what jpeg splitter it's attached to,
		// so give it a pointer to the splitter.
		v.split = js

		// Jpeg visitor's 'split' function
		// satisfies bufio.SplitFunc{}.
		scanner.Split(js.Split)

	case "image/webp", "webp":
		// Webp visitor's 'split' function
		// satisfies bufio.SplitFunc{}.
		scanner.Split((&webpVisitor{
			writer: out,
		}).Split)

	case "image/png", "png":
		// For pngs we need to skip the header bytes, so read
		// them in and check we're really dealing with a png.
		header := make([]byte, len(pngstructure.PngSignature))
		if _, headerError := in.Read(header); headerError != nil {
			return nil, nil, headerError
		} else if !bytes.Equal(header, pngstructure.PngSignature[:]) {
			return nil, nil, errors.New("could not decode png: invalid header")
		}

		// Don't bother checking CRC;
		// we're overwriting it anyway.
		ps := pngstructure.NewPngSplitter()
		ps.DoCheckCrc(false)

		// Png visitor's 'Split' function
		// satisfies bufio.SplitFunc{}.
		v := &pngVisitor{
			split: ps,
			write: out,
			last:  -1,
		}
		deferred = &v.errs
		scanner.Split(v.Split)

	default:
		return nil, nil, ErrUnsupported
	}

	return
}
