package mempool

import (
	"sync"
	"sync/atomic"
	"unsafe"

	"golang.org/x/sys/cpu"
)

// Pool provides a form of SimplePool
// with the addition of concurrency safety.
type Pool[T any] struct {
	UnsafePool

	// New is an optionally provided
	// allocator used when no value
	// is available for use in pool.
	New func() T

	// Reset is an optionally provided
	// value resetting function called
	// on passed value to Put().
	Reset func(T) bool
}

func NewPool[T any](new func() T, reset func(T) bool, check func(current, victim int) bool) Pool[T] {
	return Pool[T]{
		New:        new,
		Reset:      reset,
		UnsafePool: NewUnsafePool(check),
	}
}

func (p *Pool[T]) Get() T {
	if ptr := p.UnsafePool.Get(); ptr != nil {
		return *(*T)(ptr)
	}
	var t T
	if p.New != nil {
		t = p.New()
	}
	return t
}

func (p *Pool[T]) Put(t T) {
	if p.Reset != nil && !p.Reset(t) {
		return
	}
	ptr := unsafe.Pointer(&t)
	p.UnsafePool.Put(ptr)
}

// UnsafePool provides a form of UnsafeSimplePool
// with the addition of concurrency safety.
type UnsafePool struct {
	internal
	_ [cache_line_bytes - unsafe.Sizeof(internal{})%cache_line_bytes]byte
}

func NewUnsafePool(check func(current, victim int) bool) UnsafePool {
	return UnsafePool{internal: internal{
		pool: UnsafeSimplePool{Check: check},
	}}
}

const (
	// platform CPU cache line size to avoid false sharing.
	cache_line_bytes = unsafe.Sizeof(cpu.CacheLinePad{})
)

type internal struct {
	// fast-access ring-buffer of
	// pointers accessible by index.
	//
	// if Go ever exposes goroutine IDs
	// to us we can make this a lot faster.
	ring  [20]unsafe.Pointer
	index atomic.Uint32

	// underlying pool and
	// slow mutex protection.
	pool  UnsafeSimplePool
	mutex sync.Mutex
}

func (p *internal) Check(fn func(current, victim int) bool) func(current, victim int) bool {
	p.mutex.Lock()
	if fn == nil {
		if p.pool.Check == nil {
			fn = defaultCheck
		} else {
			fn = p.pool.Check
		}
	} else {
		p.pool.Check = fn
	}
	p.mutex.Unlock()
	return fn
}

func (p *internal) Get() unsafe.Pointer {
	const cap = uint32(cap(p.ring))
	idx := p.index.Load() % cap
	if ptr := atomic.SwapPointer(&p.ring[idx], nil); ptr != nil {
		p.index.Add(^uint32(0)) // i.e. -1
		return ptr
	}
	p.mutex.Lock()
	ptr := p.pool.Get()
	p.mutex.Unlock()
	return ptr
}

func (p *internal) Put(ptr unsafe.Pointer) {
	const cap = uint32(cap(p.ring))
	idx := p.index.Add(1) % cap
	if ptr := atomic.SwapPointer(&p.ring[idx], ptr); ptr != nil {
		p.mutex.Lock()
		p.pool.Put(ptr)
		p.mutex.Unlock()
	}
}

func (p *internal) GC() {
	for i := range p.ring {
		atomic.StorePointer(&p.ring[i], nil)
	}
	p.mutex.Lock()
	p.pool.GC()
	p.mutex.Unlock()
}

func (p *internal) Size() int {
	p.mutex.Lock()
	sz := p.pool.Size()
	p.mutex.Unlock()
	return sz
}
