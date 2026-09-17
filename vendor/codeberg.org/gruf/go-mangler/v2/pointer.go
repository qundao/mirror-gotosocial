package mangler

import (
	"reflect"
	"unsafe"

	"codeberg.org/gruf/go-xunsafe"
)

// derefPointerType returns a Mangler capable of dereferencing
// and formatting the given pointer type currently in TypeIter{}.
// note this will fetch a sub-Mangler for resulting value type.
func derefPointerType(t xunsafe.TypeIter) Mangler {
	var derefs uint8
	vt := t

	// Iteratively dereference pointer types.
	for vt.Type.Kind() == reflect.Pointer {

		// Only if this is actual indirect memory do we
		// perform a derefence, otherwise we just skip over
		// and increase the dereference indicator, i.e. '1'.
		if vt.Indirect() {

			// Once we overflow our uint8 w/ derefs
			// don't bother. This is getting absurd.
			if d := derefs + 1; d < derefs {
				return nil
			} else {
				derefs = d
			}
		}

		// Get next elem type.
		vt = vt.PointerElem()
	}

	// Get value mangler.
	fn := loadOrGet(vt)
	if fn == nil {
		return nil
	}

	return func(buf []byte, ptr unsafe.Pointer) []byte {
		var i uint8

		for ; i < derefs; i++ {
			if ptr == nil {
				break
			}

			// Further dereference ptr.
			ptr = *(*unsafe.Pointer)(ptr)
		}

		// Dereference count.
		buf = append(buf, i)

		// Final check.
		if ptr == nil {
			return buf
		}

		// Mangle deref'd.
		buf = fn(buf, ptr)
		return buf
	}
}
