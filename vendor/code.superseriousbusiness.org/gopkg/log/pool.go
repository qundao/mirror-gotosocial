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

package log

import (
	"unsafe"

	"codeberg.org/gruf/go-byteutil"
	"codeberg.org/gruf/go-mempool"
)

// memory pool of log buffers.
var bufpool mempool.UnsafePool

// getBuf acquires a buffer from memory pool.
func getBuf() *byteutil.Buffer {
	buf := (*byteutil.Buffer)(bufpool.Get())
	if buf == nil {
		buf = new(byteutil.Buffer)
		buf.B = make([]byte, 0, 512)
	}
	return buf
}

// putBuf places (after resetting) buffer back in
// memory pool, dropping if capacity too large.
func putBuf(buf *byteutil.Buffer) {
	if cap(buf.B) > int(^uint16(0)) {
		return // drop large buffer
	}
	buf.B = buf.B[:0]
	ptr := unsafe.Pointer(buf)
	bufpool.Put(ptr)
}
