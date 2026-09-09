package mangler

import (
	"unsafe"
	_ "unsafe"
)

// Notes:
//   the use of unsafe conversion from the direct interface values to
//   the chosen types in each of the below functions allows us to convert
//   not only those types directly, but anything type-aliased to those
//   types. e.g. `time.Duration` directly as int64.

func mangle_string(buf []byte, ptr unsafe.Pointer) []byte {
	s := *(*string)(ptr)
	buf = append_uint(buf, uint(len(s)))
	return append(buf, s...)
}

func mangle_string_slice(buf []byte, ptr unsafe.Pointer) []byte {
	s := *(*[]string)(ptr)
	buf = append_uint(buf, uint(len(s)))
	for _, s := range s {
		buf = append(buf, s...)
		buf = append(buf, ',')
	}
	if len(s) > 0 {
		buf = buf[:len(buf)-1]
	}
	return buf
}

func mangle_bool(buf []byte, ptr unsafe.Pointer) []byte {
	if *(*bool)(ptr) {
		return append(buf, '1')
	}
	return append(buf, '0')
}

func mangle_bool_slice(buf []byte, ptr unsafe.Pointer) []byte {
	s := *(*[]bool)(ptr)
	buf = append_uint(buf, uint(len(s)))
	for _, b := range s {
		if b {
			buf = append(buf, '1')
		} else {
			buf = append(buf, '0')
		}
	}
	return buf
}

func mangle_8bit(buf []byte, ptr unsafe.Pointer) []byte {
	return append(buf, *(*uint8)(ptr))
}

func mangle_8bit_slice(buf []byte, ptr unsafe.Pointer) []byte {
	s := *(*[]uint8)(ptr)
	buf = append_uint(buf, uint(len(s)))
	return append(buf, s...)
}

func mangle_16bit(buf []byte, ptr unsafe.Pointer) []byte {
	return append_uint(buf, *(*uint16)(ptr))
}

func mangle_16bit_slice(buf []byte, ptr unsafe.Pointer) []byte {
	s := *(*[]uint16)(ptr)
	buf = append_uint(buf, uint(len(s)))
	for _, u := range s {
		buf = append_uint(buf, u)
	}
	return buf
}

func mangle_32bit(buf []byte, ptr unsafe.Pointer) []byte {
	return append_uint(buf, *(*uint32)(ptr))
}

func mangle_32bit_slice(buf []byte, ptr unsafe.Pointer) []byte {
	s := *(*[]uint32)(ptr)
	buf = append_uint(buf, uint(len(s)))
	for _, u := range s {
		buf = append_uint(buf, u)
	}
	return buf
}

func mangle_64bit(buf []byte, ptr unsafe.Pointer) []byte {
	return append_uint(buf, *(*uint64)(ptr))
}

func mangle_64bit_slice(buf []byte, ptr unsafe.Pointer) []byte {
	s := *(*[]uint64)(ptr)
	buf = append_uint(buf, uint(len(s)))
	for _, u := range s {
		buf = append_uint(buf, u)
	}
	return buf
}

func mangle_int(buf []byte, ptr unsafe.Pointer) []byte {
	return append_uint(buf, *(*uint)(ptr))
}

func mangle_int_slice(buf []byte, ptr unsafe.Pointer) []byte {
	s := *(*[]uint)(ptr)
	buf = append_uint(buf, uint(len(s)))
	for _, u := range s {
		buf = append_uint(buf, u)
	}
	return buf
}

func mangle_128bit(buf []byte, ptr unsafe.Pointer) []byte {
	u2 := *(*[2]uint64)(ptr)
	buf = append_uint(buf, u2[0])
	buf = append_uint(buf, u2[1])
	return buf
}

func mangle_128bit_slice(buf []byte, ptr unsafe.Pointer) []byte {
	for _, u2 := range *(*[][2]uint64)(ptr) {
		buf = append_uint(buf, u2[0])
		buf = append_uint(buf, u2[1])
	}
	return buf
}
