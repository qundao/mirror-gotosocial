package mangler

import (
	"reflect"
	"unsafe"

	"codeberg.org/gruf/go-xunsafe"
)

// loadOrStore first checks the cache for a Mangler
// function, else generates one by calling get().
// note: this does store generated funcs in cache.
func loadOrStore(t xunsafe.TypeIter) Mangler {

	// Get cache key.
	key := t.TypeInfo

	// Check cache for func.
	fn := manglers.Get(key)

	if fn == nil {
		// Generate new mangler
		// func for this type.
		fn = get(t)
		if fn == nil {
			return nil
		}

		// Store func in cache.
		manglers.Put(key, fn)
	}

	return fn
}

// loadOrGet first checks the cache for a Mangler
// function, else generates one by calling get().
// note: it does not store the function in cache.
func loadOrGet(t xunsafe.TypeIter) Mangler {

	// Check cache for mangler func.
	fn := manglers.Get(t.TypeInfo)

	if fn == nil {
		// Generate new mangler
		// func for this type.
		fn = get(t)
	}

	return fn
}

var (
	// reflectTypeType is the reflected type of the reflect type,
	// used in fmt.get() to prevent iter of internal ABI structs.
	reflectTypeType = reflect.TypeOf(reflect.TypeOf(0))
)

// get attempts to generate a new Mangler function
// capable of mangling a ptr of given type information.
func get(t xunsafe.TypeIter) (fn Mangler) {
	defer func() {
		if fn == nil {
			// nothing more
			// we can do.
			return
		}

		if t.Parent != nil {
			// We're only interested
			// in wrapping top-level.
			return
		}

		// Get reflected type ptr for prefix.
		ptr := xunsafe.ReflectTypeData(t.Type)
		uptr := uintptr(ptr)

		// Outer fn.
		mng := fn

		// Wrap the mangler func to prepend type pointer.
		fn = func(buf []byte, ptr unsafe.Pointer) []byte {
			buf = append_uint(buf, uptr)
			return mng(buf, ptr)
		}
	}()

	if t.Type == nil {
		// nil type.
		return nil
	}

	if t.Type == reflectTypeType {
		// DO NOT iterate down internal ABI
		// types, some are in non-GC memory.
		return nil
	}

	// Check supports known method receiver.
	if fn := getMethodType(t); fn != nil {
		return fn
	}

	if !visit(t) {
		// don't bother for recursive types,
		// we either serialize the pointer value
		// itself which isn't deterministic, or
		// we add recursion checks to serialize
		// hot paths which i'd rather not do.
		return nil
	}

	// Get func for kind.
	switch t.Type.Kind() {
	case reflect.Pointer:
		return derefPointerType(t)
	case reflect.Struct:
		return iterStructType(t)
	case reflect.Array:
		return iterArrayType(t)
	case reflect.Slice:
		return iterSliceType(t)
	case reflect.String:
		return mangle_string
	case reflect.Bool:
		return mangle_bool
	case reflect.Int:
		return mangle_int
	case reflect.Uint:
		return mangle_uint
	case reflect.Uintptr,
		reflect.UnsafePointer:
		// mangling pointers without being
		// able to follow and know what they
		// contain is not deterministic.
		return nil
	case reflect.Int8, reflect.Uint8:
		return mangle_8bit
	case reflect.Int16:
		return mangle_int16
	case reflect.Uint16:
		return mangle_uint16
	case reflect.Int32:
		return mangle_int32
	case reflect.Uint32:
		return mangle_uint32
	case reflect.Int64:
		return mangle_int64
	case reflect.Uint64:
		return mangle_uint64
	case reflect.Float32:
		return mangle_uint32
	case reflect.Float64:
		return mangle_uint64
	default:
		return nil
	}
}
