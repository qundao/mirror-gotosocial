package format

import (
	"unsafe"

	"codeberg.org/gruf/go-xunsafe"
)

// iterMapType returns a FormatFunc capable of iterating
// and formatting the given map type currently in TypeIter{}.
// note this will fetch sub-FormatFuncs for key / value types.
func (fmt *Formatter) iterMapType(t xunsafe.TypeIter) FormatFunc {

	// Get nested k / v.
	kt := t.MapKey()
	vt := t.MapElem()
	kt.Flag |= flagKeyType

	// Get key format func.
	kfn := fmt.loadOrGet(kt)
	if kfn == nil {
		panic("unreachable")
	}

	// Get value format func.
	vfn := fmt.loadOrGet(vt)
	if vfn == nil {
		panic("unreachable")
	}

	// Final map type.
	rtype := t.Type
	flags := t.Flag

	if !needs_typestr(t) {
		return func(s *State) {
			if s.P == nil || *(*unsafe.Pointer)(s.P) == nil {
				// Append nil.
				appendNil(s)
				return
			}

			// Build reflect value, and then a map iterator.
			v := xunsafe.BuildReflectValue(rtype, s.P, flags)
			i := xunsafe.GetMapIter(v)

			// Prepend object brace.
			s.B = append(s.B, '{')

			// Before len.
			l := len(s.B)

			for i.Next() {
				// Pass to map key func.
				s.P = xunsafe.Map_Key(i)
				kfn(s)

				// Add key seperator.
				s.B = append(s.B, '=')

				// Pass to map elem func.
				s.P = xunsafe.Map_Elem(i)
				vfn(s)

				// Add comma pair seperator.
				s.B = append(s.B, ',', ' ')
			}

			if len(s.B) != l {
				// Drop final ", ".
				s.B = s.B[:len(s.B)-2]
			}

			// Append object brace.
			s.B = append(s.B, '}')
		}
	}

	// Map type string with ptrs / refs.
	typestrPtrs := typestr_with_ptrs(t)
	typestrRefs := typestr_with_refs(t)

	return func(s *State) {
		if s.P == nil || *(*unsafe.Pointer)(s.P) == nil {
			// Append nil value with type.
			appendNilType(s, typestrPtrs)
			return
		}

		// Build reflect value, and then a map iter.
		v := xunsafe.BuildReflectValue(rtype, s.P, flags)
		i := xunsafe.GetMapIter(v)

		// Include type info.
		if s.A.WithType() {
			s.B = append(s.B, typestrRefs...)
		}

		// Prepend object brace.
		s.B = append(s.B, '{')

		// Before len.
		l := len(s.B)

		for i.Next() {
			// Pass to map key func.
			s.P = xunsafe.Map_Key(i)
			kfn(s)

			// Add key seperator.
			s.B = append(s.B, '=')

			// Pass to map elem func.
			s.P = xunsafe.Map_Elem(i)
			vfn(s)

			// Add comma pair seperator.
			s.B = append(s.B, ',', ' ')
		}

		if len(s.B) != l {
			// Drop final ", ".
			s.B = s.B[:len(s.B)-2]
		}

		// Append object brace.
		s.B = append(s.B, '}')
	}
}
