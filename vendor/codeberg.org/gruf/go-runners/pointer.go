package runners

import (
	"sync/atomic"
	"unsafe"
)

// atomic_pointer wraps an unsafe.Pointer with
// receiver methods for their atomic counterparts.
type atomic_pointer struct{ p unsafe.Pointer }

func (p *atomic_pointer) Load() unsafe.Pointer {
	return atomic.LoadPointer(&p.p)
}

func (p *atomic_pointer) Store(ptr unsafe.Pointer) {
	atomic.StorePointer(&p.p, ptr)
}

func (p *atomic_pointer) CAS(old, new unsafe.Pointer) bool {
	return atomic.CompareAndSwapPointer(&p.p, old, new)
}
