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
	"io"

	pngstructure "code.superseriousbusiness.org/go-png-image-structure/v2"
)

type pngVisitor struct {
	split *pngstructure.PngSplitter
	write io.Writer
	last  int

	// exif termination
	// errors, handled
	// externally.
	errs []error
}

func (v *pngVisitor) Split(data []byte, atEOF bool) (int, []byte, error) {
	// execute the ps split function to read in data.
	advance, token, err := v.split.Split(data, atEOF)
	if err != nil {
		return advance, token, err
	}

	// if we haven't written anything at all yet, then
	// write the png header back into the writer first.
	if v.last == -1 {
		if _, err := v.write.Write(pngstructure.PngSignature[:]); err != nil {
			return advance, token, err
		}
	}

	// Check if the splitter now has
	// any new chunks in it for us.
	chunkSlice, err := v.split.Chunks()
	if err != nil {
		return advance, token, err
	}

	// Extract chunks from slice.
	chunks := chunkSlice.Chunks()

	// Iterate, terminate and write chunks, from last written.
	for i := v.last + 1; i < len(chunks); i++ {
		chunk := chunks[i]

		if chunk.Type == pngstructure.EXifChunkType {
			// Finally, some exif data! Attempt to terminate!
			if err := terminateEXIF(chunkSlice); err != nil {
				v.errs = append(v.errs, err)
			}

			// Update chunk crc.
			chunk.UpdateCrc32()
		}

		// Write this particular chunk to underlying writer.
		if _, err := chunk.WriteTo(v.write); err != nil {
			return advance, token, err
		}

		// Update
		// chunk
		// index.
		v.last = i

		// Zero data; here you
		// go garbage collector.
		chunk.Data = nil
	}

	return advance, token, err
}
