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
	"math"
	"sync"
	"unsafe"

	"codeberg.org/gruf/go-byteutil"
	"codeberg.org/gruf/go-mempool"
)

var (
	// global map of memory pool instances.
	buffers = make(map[uint32]*MemoryPool, 8)

	// global map lock.
	mutex sync.Mutex
)

// MemoryPool is a memory pool of
// byte buffers, of a predefined size.
//
// This will be a shared global instance that
// any callers of buffers.Pool($sz) can access.
type MemoryPool struct {
	p mempool.UnsafePool
	s uint32
}

// Pool returns a shared MemoryPool instance
// of requested size (rounded to nearest 2^n).
//
// NOTE: this acquires a lock on a global mutex
// instance and should generally only be called
// on package init to store a global reference.
func Pool(sz uint32) *MemoryPool {
	n := math.Log2(float64(sz))
	n = math.Round(n)
	n = math.Exp2(max(8, n))
	sz = uint32(n)
	mutex.Lock()
	p := buffers[sz]
	if p == nil {
		p = new(MemoryPool)
		p.s = sz
		buffers[sz] = p
	}
	mutex.Unlock()
	return p
}

// Get returns a byteutil.Buffer{} instance from pool.
func (p *MemoryPool) Get() *byteutil.Buffer {
	buf := (*byteutil.Buffer)(p.p.Get())
	if buf == nil {
		buf = new(byteutil.Buffer)
		buf.B = make([]byte, p.s)
	} else {
		clear(buf.B[0:cap(buf.B)])
	}
	buf.B = buf.B[:0]
	return buf
}

// Put replaces byteutil.Buffer{} instance in pool.
func (p *MemoryPool) Put(buf *byteutil.Buffer) {
	if buf == nil {
		return
	}
	if cap(buf.B) < int(p.s) || cap(buf.B) > 2*int(p.s) {
		return // drop buffers outside size range
	}
	p.p.Put(unsafe.Pointer(buf))
}
