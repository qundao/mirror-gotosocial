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
	var derefs int8
	var indirects uint64
	vt := t

	// Iteratively dereference pointer types.
	for vt.Type.Kind() == reflect.Pointer {

		// Only if this is actual indirect memory do we
		// perform a derefence, otherwise we just skip over
		// and increase the dereference indicator, i.e. '1'.
		if vt.Indirect() {
			indirects |= (1 << derefs)
		}
		derefs++

		// Get next elem type.
		vt = vt.PointerElem()
	}

	// Ensure this is a reasonable number of derefs.
	if derefs > 4*int8(unsafe.Sizeof(indirects)) {
		return nil
	}

	// Get value mangler.
	fn := loadOrGet(vt)
	if fn == nil {
		return nil
	}

	return func(buf []byte, ptr unsafe.Pointer) []byte {
		var i int8

		for ; i < derefs; i++ {
			switch {
			case indirects&(1<<i) == 0:

			case ptr == nil:
				// Nil value, return here.
				buf = append(buf, uint8(i))
				return buf

			default:
				// Further dereference ptr.
				ptr = *(*unsafe.Pointer)(ptr)
			}
		}

		if ptr == nil {
			// Final nil value check.
			buf = append(buf, uint8(i))
			return buf
		}

		// Mangle fully dereferenced.
		buf = append(buf, uint8(i))
		buf = fn(buf, ptr)
		return buf
	}
}
