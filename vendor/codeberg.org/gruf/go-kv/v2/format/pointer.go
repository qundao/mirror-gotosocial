package format

import (
	"reflect"
	"unsafe"

	"codeberg.org/gruf/go-xunsafe"
)

// derefPointerType returns a FormatFunc capable of dereferencing
// and formatting the given pointer type currently in TypeIter{}.
// note this will fetch a sub-FormatFunc for resulting value type.
func (fmt *Formatter) derefPointerType(t xunsafe.TypeIter) FormatFunc {
	var n int
	vt := t

	// Iteratively dereference pointer types.
	for vt.Type.Kind() == reflect.Pointer {

		// If this actual indirect memory,
		// increase dereferences counter.
		if vt.Indirect() {
			n++
		}

		// Get next elem type.
		vt = vt.PointerElem()
	}

	// Get value format func.
	fn := fmt.loadOrGet(vt)
	if fn == nil {
		panic("unreachable")
	}

	if !needs_typestr(t) {
		if n <= 0 {
			// No derefs are needed.
			return func(s *State) {
				if s.P == nil {
					// Final check.
					appendNil(s)
					return
				}

				// Format
				// final
				// value.
				fn(s)
			}
		}

		return func(s *State) {
			// Deref n number times.
			for i := n; i > 0; i-- {

				if s.P == nil {
					// Nil check.
					appendNil(s)
					return
				}

				// Further deref pointer value.
				s.P = *(*unsafe.Pointer)(s.P)
			}

			if s.P == nil {
				// Final check.
				appendNil(s)
				return
			}

			// Format
			// final
			// value.
			fn(s)
		}
	}

	// Final type string with ptrs.
	typestr := typestr_with_ptrs(t)

	if n <= 0 {
		// No derefs are needed.
		return func(s *State) {
			if s.P == nil {
				// Final nil value check.
				appendNilType(s, typestr)
				return
			}

			// Format
			// final
			// value.
			fn(s)
		}
	}

	return func(s *State) {
		// Deref n number times.
		for i := n; i > 0; i-- {
			if s.P == nil {
				// Check for nil value.
				appendNilType(s, typestr)
				return
			}

			// Further deref pointer value.
			s.P = *(*unsafe.Pointer)(s.P)
		}

		if s.P == nil {
			// Final nil value check.
			appendNilType(s, typestr)
			return
		}

		// Format
		// final
		// value.
		fn(s)
	}
}
