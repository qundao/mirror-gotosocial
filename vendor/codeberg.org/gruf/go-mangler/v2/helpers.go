package mangler

import (
	"unsafe"
)

func empty_mangler(buf []byte, _ unsafe.Pointer) []byte {
	return buf
}

// unsigned integer leb128 variable length encoding, much
// less character output at little processing regression.
func append_uint[Uint uint | uintptr | uint16 | uint32 | uint64](b []byte, v Uint) []byte {
	var c uint8
	for {
		c = uint8(v & 0x7f)
		v >>= 7
		if v != 0 {
			c |= 0x80
		}
		b = append(b, c)
		if c&0x80 == 0 {
			break
		}
	}
	return b
}

// add returns the ptr addition of starting ptr and a delta.
func add(ptr unsafe.Pointer, delta uintptr) unsafe.Pointer {
	return unsafe.Pointer(uintptr(ptr) + delta)
}
