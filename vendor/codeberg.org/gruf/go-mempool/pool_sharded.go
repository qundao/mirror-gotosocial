package mempool

// const (
// 	// platform CPU cache line size to avoid false sharing.
// 	cache_line_bytes = unsafe.Sizeof(cpu.CacheLinePad{})
// )

// // UnsafePool provides a form of UnsafeSimplePool
// // with the addition of concurrency safety.
// type UnsafePool struct {
// 	pool_internal
// 	_ [cache_line_bytes - unsafe.Sizeof(pool_internal{})]byte
// }

// func NewUnsafePool(check func(current, victim int) bool) UnsafePool {
// 	return UnsafePool{pool_internal: pool_internal{
// 		check: unsafe.Pointer(&check),
// 	}}
// }

// type pool_internal struct {
// 	shard unsafe.Pointer // *shards
// 	check unsafe.Pointer // *func(current, victim int) bool
// 	index atomic.Uint32
// }

// func (p *pool_internal) Check(check func(current, victim int) bool) func(current, victim int) bool {
// 	if check == nil {
// 		check = p.load_check()
// 		if check == nil {
// 			check = defaultCheck
// 		}
// 		return check
// 	}
// 	ptr := unsafe.Pointer(&check)
// 	atomic.StorePointer(&p.check, ptr)
// 	return check
// }

// func (p *pool_internal) Get() unsafe.Pointer {
// 	shards := p.load_shards()
// 	idx := p.index.Add(1) % uint32(len(shards))
// 	if ptr, ok := shards[idx].TryGet(); ptr != nil {
// 		return ptr
// 	} else if ok {
// 		idx++
// 	}
// 	for i := idx; i < uint32(len(shards)); i++ {
// 		if ptr := shards[i].Get(); ptr != nil {
// 			return ptr
// 		}
// 	}
// 	for i := uint32(0); i < idx; i++ {
// 		if ptr := shards[i].Get(); ptr != nil {
// 			return ptr
// 		}
// 	}
// 	return nil
// }

// func (p *pool_internal) Put(ptr unsafe.Pointer) {
// 	shards := p.load_shards()
// 	idx := p.index.Add(1) % uint32(len(shards))
// 	shards[idx].Put(ptr)
// }

// func (p *pool_internal) GC() {
// 	shards := p.load_shards()
// 	for i := range shards {
// 		for j := range shards[i].priv {
// 			atomic.StorePointer(&shards[i].priv[j], nil)
// 		}
// 	}
// 	for i := range shards {
// 		shards[i].GC()
// 	}
// }

// func (p *pool_internal) Size() (sz int) {
// 	shards := p.load_shards()
// 	for i := range shards {
// 		sz += shards[i].Size()
// 	}
// 	return
// }

// func (p *pool_internal) Clear() {
// 	atomic.StorePointer(&p.shard, nil)
// }

// // load_shards ...
// func (p *pool_internal) load_shards() []pool_shard {
// 	for {
// 		// Try load existing shards pointer.
// 		ptr := atomic.LoadPointer(&p.shard)
// 		shards := (*[]pool_shard)(ptr)
// 		if ptr != nil {
// 			return *shards
// 		}

// 		// Load check function.
// 		check := p.load_check()

// 		// Allocate new shards.
// 		shards = new([]pool_shard)
// 		(*shards) = make([]pool_shard, runtime.GOMAXPROCS(0))
// 		for i := range *shards {
// 			(*shards)[i].pool.Check = check
// 		}

// 		// Attempt to set the new shards pointer.
// 		if atomic.CompareAndSwapPointer(&p.shard,
// 			ptr,
// 			unsafe.Pointer(shards),
// 		) {
// 			return *shards
// 		}
// 	}
// }

// // load_check ...
// func (p *pool_internal) load_check() (check func(current, victim int) bool) {
// 	if ptr := atomic.LoadPointer(&p.check); ptr != nil {
// 		check = *(*func(int, int) bool)(ptr)
// 	}
// 	return
// }

// type pool_shard struct {
// 	pool_shard_internal
// 	_ [cache_line_bytes - unsafe.Sizeof(pool_shard_internal{})%cache_line_bytes]byte
// }

// type pool_shard_internal struct {
// 	priv [4]unsafe.Pointer
// 	pool UnsafeSimplePool
// 	lock sync.Mutex
// }

// func (p *pool_shard_internal) TryGet() (ptr unsafe.Pointer, locked bool) {
// 	if ptr = atomic.SwapPointer(&p.priv[0], nil); ptr != nil {
// 		return ptr, false
// 	}
// 	if ptr = atomic.SwapPointer(&p.priv[1], nil); ptr != nil {
// 		return ptr, false
// 	}
// 	if ptr = atomic.SwapPointer(&p.priv[2], nil); ptr != nil {
// 		return ptr, false
// 	}
// 	if ptr = atomic.SwapPointer(&p.priv[3], nil); ptr != nil {
// 		return ptr, false
// 	}
// 	if !p.lock.TryLock() {
// 		return nil, false
// 	}
// 	ptr = p.pool.Get()
// 	p.lock.Unlock()
// 	return ptr, true
// }

// func (p *pool_shard_internal) Get() unsafe.Pointer {
// 	p.lock.Lock()
// 	ptr := p.pool.Get()
// 	p.lock.Unlock()
// 	return ptr
// }

// func (p *pool_shard_internal) Put(ptr unsafe.Pointer) {
// 	if ptr = atomic.SwapPointer(&p.priv[0], ptr); ptr == nil {
// 		return
// 	}
// 	if ptr = atomic.SwapPointer(&p.priv[1], ptr); ptr == nil {
// 		return
// 	}
// 	if ptr = atomic.SwapPointer(&p.priv[2], ptr); ptr == nil {
// 		return
// 	}
// 	if ptr = atomic.SwapPointer(&p.priv[3], ptr); ptr == nil {
// 		return
// 	}
// 	p.lock.Lock()
// 	p.pool.Put(ptr)
// 	p.lock.Unlock()
// }

// func (p *pool_shard_internal) GC() {
// 	p.lock.Lock()
// 	p.pool.GC()
// 	p.lock.Unlock()
// }

// func (p *pool_shard_internal) Size() int {
// 	p.lock.Lock()
// 	sz := p.pool.Size()
// 	p.lock.Unlock()
// 	return sz
// }
