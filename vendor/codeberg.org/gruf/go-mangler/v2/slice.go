package mangler

import (
	"reflect"
	"unsafe"

	"codeberg.org/gruf/go-xunsafe"
)

// iterSliceType returns a Mangler capable of iterating
// and mangling the given slice type currently in TypeIter{}.
// note this will fetch sub-Mangler for slice element type.
func iterSliceType(t xunsafe.TypeIter) Mangler {

	// Get nested elem.
	et := t.SliceElem()
	esz := et.Type.Size()

	// Prefer to use a known slice mangler func.
	if fn := mangleKnownSlice(et); fn != nil {
		return fn
	}

	// Get elem mangler.
	fn := loadOrGet(et)
	if fn == nil {
		return nil
	}

	return func(buf []byte, ptr unsafe.Pointer) []byte {
		// Cast data as unsafe slice header type.
		hdr := (*xunsafe.Unsafeheader_Slice)(ptr)
		if hdr == nil || hdr.Data == nil {
			return buf
		}

		// Append slice length.
		buf = append_uint(buf,
			uint(hdr.Len))

		var offset uintptr
		for range hdr.Len {
			// Mangle data at slice index.
			ptr = unsafe.Add(hdr.Data, offset)
			buf = fn(buf, ptr)
			buf = append(buf, ',')
			offset += esz
		}

		if hdr.Len > 0 {
			// Drop final comma.
			buf = buf[:len(buf)-1]
		}

		return buf
	}
}

// mangleKnownSlice loads a Mangler function for a
// known slice-of-element type (in this case, primtives).
func mangleKnownSlice(t xunsafe.TypeIter) Mangler {
	switch t.Type.Kind() {
	case reflect.String:
		return mangle_string_slice
	case reflect.Bool:
		return mangle_bool_slice
	case reflect.Int:
		return mangle_int_slice
	case reflect.Uint:
		return mangle_uint_slice
	case reflect.Int8,
		reflect.Uint8:
		return mangle_8bit_slice
	case reflect.Int16:
		return mangle_int16_slice
	case reflect.Uint16:
		return mangle_uint16_slice
	case reflect.Int32:
		return mangle_int32_slice
	case reflect.Uint32:
		return mangle_uint32_slice
	case reflect.Int64:
		return mangle_int64_slice
	case reflect.Uint64:
		return mangle_uint64_slice
	case reflect.Float32:
		return mangle_int32_slice
	case reflect.Float64:
		return mangle_int64_slice
	default:
		return nil
	}
}
