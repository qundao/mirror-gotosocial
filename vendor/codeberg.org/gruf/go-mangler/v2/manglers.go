package mangler

import (
	"unsafe"
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
	if s == nil {
		return buf
	}
	buf = append_uint(buf, uint(len(s)))
	for _, s := range s {
		buf = append_uint(buf, uint(len(s)))
		buf = append(buf, s...)
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
	if s == nil {
		return buf
	}
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

func mangle_int(buf []byte, ptr unsafe.Pointer) []byte {
	return append_int(buf, *(*int)(ptr))
}

func mangle_int_slice(buf []byte, ptr unsafe.Pointer) []byte {
	s := *(*[]int)(ptr)
	if s == nil {
		return buf
	}
	buf = append_uint(buf, uint(len(s)))
	for _, v := range s {
		buf = append_int(buf, v)
	}
	return buf
}

func mangle_uint(buf []byte, ptr unsafe.Pointer) []byte {
	return append_uint(buf, *(*uint)(ptr))
}

func mangle_uint_slice(buf []byte, ptr unsafe.Pointer) []byte {
	s := *(*[]uint)(ptr)
	if s == nil {
		return buf
	}
	buf = append_uint(buf, uint(len(s)))
	for _, v := range s {
		buf = append_uint(buf, v)
	}
	return buf
}

func mangle_int16(buf []byte, ptr unsafe.Pointer) []byte {
	return append_int(buf, *(*int16)(ptr))
}

func mangle_int16_slice(buf []byte, ptr unsafe.Pointer) []byte {
	s := *(*[]int16)(ptr)
	if s == nil {
		return buf
	}
	buf = append_uint(buf, uint(len(s)))
	for _, v := range s {
		buf = append_int(buf, v)
	}
	return buf
}

func mangle_uint16(buf []byte, ptr unsafe.Pointer) []byte {
	return append_uint(buf, *(*uint16)(ptr))
}

func mangle_uint16_slice(buf []byte, ptr unsafe.Pointer) []byte {
	s := *(*[]uint16)(ptr)
	if s == nil {
		return buf
	}
	buf = append_uint(buf, uint(len(s)))
	for _, v := range s {
		buf = append_uint(buf, v)
	}
	return buf
}

func mangle_int32(buf []byte, ptr unsafe.Pointer) []byte {
	return append_int(buf, *(*int32)(ptr))
}

func mangle_int32_slice(buf []byte, ptr unsafe.Pointer) []byte {
	s := *(*[]int32)(ptr)
	if s == nil {
		return buf
	}
	buf = append_uint(buf, uint(len(s)))
	for _, v := range s {
		buf = append_int(buf, v)
	}
	return buf
}

func mangle_uint32(buf []byte, ptr unsafe.Pointer) []byte {
	return append_uint(buf, *(*uint32)(ptr))
}

func mangle_uint32_slice(buf []byte, ptr unsafe.Pointer) []byte {
	s := *(*[]uint32)(ptr)
	if s == nil {
		return buf
	}
	buf = append_uint(buf, uint(len(s)))
	for _, v := range s {
		buf = append_uint(buf, v)
	}
	return buf
}

func mangle_int64(buf []byte, ptr unsafe.Pointer) []byte {
	return append_int(buf, *(*int64)(ptr))
}

func mangle_int64_slice(buf []byte, ptr unsafe.Pointer) []byte {
	s := *(*[]int64)(ptr)
	if s == nil {
		return buf
	}
	buf = append_uint(buf, uint(len(s)))
	for _, v := range s {
		buf = append_int(buf, v)
	}
	return buf
}

func mangle_uint64(buf []byte, ptr unsafe.Pointer) []byte {
	return append_uint(buf, *(*uint64)(ptr))
}

func mangle_uint64_slice(buf []byte, ptr unsafe.Pointer) []byte {
	s := *(*[]uint64)(ptr)
	if s == nil {
		return buf
	}
	buf = append_uint(buf, uint(len(s)))
	for _, v := range s {
		buf = append_uint(buf, v)
	}
	return buf
}

func mangle_8bit(buf []byte, ptr unsafe.Pointer) []byte {
	return append(buf, *(*uint8)(ptr))
}

func mangle_8bit_slice(buf []byte, ptr unsafe.Pointer) []byte {
	s := *(*[]uint8)(ptr)
	if s == nil {
		return buf
	}
	buf = append_uint(buf, uint(len(s)))
	return append(buf, s...)
}
