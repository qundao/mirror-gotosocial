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
		for i := int8(0); i < derefs; i++ {
			switch {
			case indirects&(1<<i) == 0:
				// No dereference needed.
				buf = append(buf, '1')

			case ptr == nil:
				// Nil value, return here.
				buf = append(buf, '0')
				return buf

			default:
				// Further deref ptr.
				buf = append(buf, '1')
				ptr = *(*unsafe.Pointer)(ptr)
			}
		}

		if ptr == nil {
			// Final nil val check.
			buf = append(buf, '0')
			return buf
		}

		// Mangle fully deref'd.
		buf = append(buf, '1')
		buf = fn(buf, ptr)
		return buf
	}
}
