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

package buffers

import (
	"io"
	"unsafe"

	"codeberg.org/gruf/go-byteutil"
	"codeberg.org/gruf/go-fastcopy"
)

// global fastcopy instance.
var copypool = fastcopy.Global()

func init() {
	// Get new pool.
	p := Pool(16384)

	// Set fastcopy.CopyPool{} getter and setter functions, ensuring buffer is set to max size.
	copypool.Get = func() *[]byte { buf := p.Get(); buf.B = buf.Full(); return toBytes(buf) }
	copypool.Put = func(b *[]byte) { p.Put(toBuffer(b)) }
}

// Copy: see fastcopy.CopyPool{}.Copy().
func Copy(dst io.Writer, src io.Reader) (int64, error) {
	return copypool.Copy(dst, src)
}

// CopyN: see fastcopy.CopyPool{}.CopyN().
func CopyN(dst io.Writer, src io.Reader, n int64) (int64, error) {
	return copypool.CopyN(dst, src, n)
}

// CopyBuffer: see fastcopy.CopyBuffer().
func CopyBuffer(dst io.Writer, src io.Reader, buf *byteutil.Buffer) (int64, error) {
	return fastcopy.CopyBuffer(dst, src, toBytes(buf))
}

// toBytes casts buffer to bytes with compile-time assertion,
// this allows us to use our byteutil.Buffer{} pool for []bytes too.
// this is only possible due to the confirmed memory semantics below.
func toBytes(buf *byteutil.Buffer) (b *[]byte) {
	if unsafe.Sizeof(buf) != unsafe.Sizeof(b) ||
		unsafe.Offsetof(buf.B) != 0 {
		panic("compile time assertion")
	}
	b = (*[]byte)(unsafe.Pointer(buf))
	return
}

// toBuffer casts bytes to buffer with compile-time assertion,
// this allows us to use our byteutil.Buffer{} pool for []bytes too.
// this is only possible due to the confirmed memory semantics below.
func toBuffer(b *[]byte) (buf *byteutil.Buffer) {
	if unsafe.Sizeof(buf) != unsafe.Sizeof(b) ||
		unsafe.Offsetof(buf.B) != 0 {
		panic("compile time assertion")
	}
	buf = (*byteutil.Buffer)(unsafe.Pointer(b))
	return
}
