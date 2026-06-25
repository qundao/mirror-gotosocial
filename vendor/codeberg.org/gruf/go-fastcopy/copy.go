package fastcopy

import (
	"errors"
	"io"
)

var (
	// global pool.
	global CopyPool

	// errInvalidWrite means that a write returned an impossible count.
	errInvalidWrite = errors.New("invalid write result")
)

// Global returns global CopyPool{} instance.
func Global() *CopyPool { return &global }

// See CopyPool.CopyN().
func CopyN(dst io.Writer, src io.Reader, n int64) (int64, error) {
	return global.CopyN(dst, src, n)
}

// See CopyPool.Copy().
func Copy(dst io.Writer, src io.Reader) (int64, error) {
	return global.Copy(dst, src)
}

// CopyPool provides a memory pool of byte
// buffers for io copies from read -> writer.
type CopyPool struct {
	Get func() *[]byte
	Put func(*[]byte)
}

// get attempts to fetch buffer from
// currently-set CopyPool{}.Get() func.
func (cp *CopyPool) get() *[]byte {
	if cp.Get != nil {
		return cp.Get()
	}
	return nil
}

// put attempts to replace buffer into
// currently-set CopyPool{}.Put() func.
func (cp *CopyPool) put(buf *[]byte) {
	if cp.Put != nil {
		cp.Put(buf)
	}
}

// CopyN performs the same logic as io.CopyN(), with the difference
// being that the byte buffer is acquired from a memory pool if required.
func (cp *CopyPool) CopyN(dst io.Writer, src io.Reader, n int64) (int64, error) {
	written, err := cp.Copy(dst, io.LimitReader(src, n))
	if written == n {
		return n, nil
	}
	if written < n && err == nil {
		// src stopped early;
		// must have been EOF.
		err = io.EOF
	}
	return written, err
}

// Copy performs the same logic as io.Copy(), with the difference
// being that the byte buffer is acquired from a memory pool if required.
func (cp *CopyPool) Copy(dst io.Writer, src io.Reader) (int64, error) {

	// Prefer using io.WriterTo to do
	// the copy (avoids alloc + copy).
	if wt, ok := src.(io.WriterTo); ok {
		return wt.WriteTo(dst)
	}

	// Prefer using io.ReaderFrom to copy.
	if rt, ok := dst.(io.ReaderFrom); ok {
		return rt.ReadFrom(src)
	}

	// Acquire buf.
	buf := cp.get()

	// Perform copy operation with buf.
	n, err := CopyBuffer(dst, src, buf)

	// Release.
	cp.put(buf)
	return n, err
}

// CopyBuffer performs the same logic as io.Copy(), without any optimized
// checks for io.WriterTo{} or io.ReaderFrom{}. Data will definitively be
// copied from source to destination using the given byte buffer. This allows
// callers to wrap this with flexible buffer sourcing, knowing it will be used.
func CopyBuffer(dst io.Writer, src io.Reader, buf *[]byte) (int64, error) {
	switch {
	case dst == nil:
		panic("nil dst")
	case src == nil:
		panic("nil src")
	case buf == nil || len(*buf) == 0:
		new := make([]byte, 4096)
		buf = &new
	default:
		// Ensure full buf available.
		(*buf) = (*buf)[0:cap(*buf)]
	}
	var n int64
	for {
		// Perform read into buf.
		nr, err := src.Read(*buf)
		if nr > 0 {

			// We error check AFTER checking
			// no. read bytes so incomplete
			// read still gets written up to nr.

			// Perform next write from buf.
			nw, ew := dst.Write((*buf)[0:nr])

			// Check for valid write.
			if nw < 0 || nr < nw {
				if ew == nil {
					ew = errInvalidWrite
				}
				return n, ew
			}

			// Incr total.
			n += int64(nw)

			// Check for
			// write error
			if ew != nil {
				return n, ew
			}

			// Check unequal
			// read / writes.
			if nr != nw {
				return n, io.ErrShortWrite
			}
		}

		// Return on
		// any set error.
		if err != nil {
			if err == io.EOF {
				err = nil // expected
			}
			return n, err
		}
	}
}
