package mangler

import (
	"unsafe"

	"codeberg.org/gruf/go-xunsafe"
)

// iterArrayType returns a Mangler capable of iterating
// and mangling the given array type currently in TypeIter{}.
// note this will fetch sub-Mangler for array element type.
func iterArrayType(t xunsafe.TypeIter) Mangler {

	// Get nested elem.
	et := t.ArrayElem()

	// Get elem mangler.
	fn := loadOrGet(et)
	if fn == nil {
		return nil
	}

	// Array element in-memory size.
	esz := t.Type.Elem().Size()

	// No of elements.
	n := t.Type.Len()
	switch n {
	case 0:
		return empty_mangler
	case 1:
		return fn
	default:
		return func(buf []byte, ptr unsafe.Pointer) []byte {
			var offset uintptr
			for range n {
				// Mangle data at array index.
				eptr := unsafe.Add(ptr, offset)
				buf = fn(buf, eptr)
				buf = append(buf, ',')
				offset += esz
			}

			// Drop final comma.
			return buf[:len(buf)-1]
		}
	}
}
