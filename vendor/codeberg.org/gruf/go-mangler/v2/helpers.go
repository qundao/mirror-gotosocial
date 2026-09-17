package mangler

import (
	"unsafe"
)

// unsigned integer leb128 variable length encoding, much
// less character output at little processing regression.
//
// note that uint8 types are expected to be handled as regular bytes.
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

// signed integer leb128 variable length encoding, much
// less character output at little processing regression.
//
// note that int8 types are expected to be handled as regular bytes.
func append_int[Int int | int16 | int32 | int64](b []byte, v Int) []byte {
	var c, s uint8
	for {
		c = uint8(v & 0x7f)
		s = uint8(v & 0x40)
		v >>= 7
		if (v != -1 || s == 0) &&
			(v != 0 || s != 0) {
			c |= 0x80
		}
		b = append(b, c)
		if c&0x80 == 0 {
			break
		}
	}
	return b
}

func empty_mangler(buf []byte, _ unsafe.Pointer) []byte {
	return buf
}
