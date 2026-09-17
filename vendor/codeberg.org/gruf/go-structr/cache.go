package structr

import (
	"context"
	"errors"
	"unsafe"

	"codeberg.org/gruf/go-byteutil"
)

// DefaultIgnoreErr is the default function used to
// ignore (i.e. not cache) incoming error results during
// Load() calls. By default ignores context pkg errors.
func DefaultIgnoreErr(err error) bool {
	return errors.Is(err, context.Canceled) ||
		errors.Is(err, context.DeadlineExceeded)
}

// CacheConfig defines config vars
// for initializing a Cache{} type.
type CacheConfig[StructType any] struct {

	// IgnoreErr defines which errors to
	// ignore (i.e. not cache) returned
	// from load function callback calls.
	// This may be left as nil, on which
	// DefaultIgnoreErr will be used.
	IgnoreErr func(error) bool

	// Copy provides a means of copying
	// cached values, to ensure returned values
	// do not share memory with those in cache.
	Copy func(StructType) StructType

	// Invalidate is called when cache values
	// (NOT errors) are invalidated, either
	// as the values passed to Put() / Store(),
	// or by the keys by calls to Invalidate().
	Invalidate func(StructType)

	// Indices defines indices to create
	// in the Cache for the receiving
	// generic struct type parameter.
	Indices []IndexConfig

	// MaxSize defines the maximum number
	// of items allowed in the Cache at
	// one time, before old items start
	// getting evicted.
	MaxSize int
}

// Cache provides a structure cache with automated
// indexing and lookups by any initialization-defined
// combination of fields. This also supports caching
// of negative results (errors!) returned by LoadOne().
type Cache[StructType any] struct {

	// hook functions.
	ignore  func(error) bool
	copy    func(StructType) StructType
	invalid func(StructType)

	// keeps track of all indexed items,
	// in order of last recently used (LRU).
	lru list

	// indices used in storing passed struct
	// types by user defined sets of fields.
	indices []Index

	// max cache size, imposes size
	// limit on the lru list in order
	// to evict old entries.
	maxSize int

	// protective mutex, guards:
	// - Cache{}.*
	// - Index{}.data
	mutex mutex
}

// Init initializes the cache with given configuration
// including struct fields to index, and necessary fns.
func (c *Cache[T]) Init(config CacheConfig[T]) {
	t := get_type_iter[T]()

	if len(config.Indices) == 0 {
		panic("no indices provided")
	}

	if config.IgnoreErr == nil {
		config.IgnoreErr = DefaultIgnoreErr
	}

	if config.Copy == nil {
		panic("copy function must be provided")
	}

	if config.MaxSize < 2 {
		panic("minimum cache size is 2 for LRU to work")
	}

	// Acquire lock.
	c.mutex.Lock()
	defer c.mutex.Unlock()

	// Ensure not setup.
	if c.indices != nil {
		panic("already initialized")
	}

	// Prepare each of configured struct indices.
	c.indices = make([]Index, len(config.Indices))
	for i, cfg := range config.Indices {
		c.indices[i].ptr = unsafe.Pointer(c)
		c.indices[i].init(t, cfg, config.MaxSize)
	}

	// Copy over functions.
	c.ignore = config.IgnoreErr
	c.copy = config.Copy

	// Copy over size configs.
	c.invalid = config.Invalidate
	c.maxSize = config.MaxSize
}

// Index selects index with given name from cache, else panics.
func (c *Cache[T]) Index(name string) *Index {
	for i := range c.indices {
		if c.indices[i].name == name {
			return &(c.indices[i])
		}
	}
	panic("unknown index: " + name)
}

// GetOne fetches value from cache stored under index, using precalculated index key.
func (c *Cache[T]) GetOne(index *Index, key Key) (T, bool) {
	values := c.Get(index, key)
	if len(values) == 0 {
		var zero T
		return zero, false
	}
	return values[0], true
}

// Get fetches values from the cache stored under index, using precalculated index keys.
func (c *Cache[T]) Get(index *Index, keys ...Key) []T {
	if index == nil {
		panic("no index given")
	} else if index.ptr != unsafe.Pointer(c) {
		panic("invalid index for cache")
	}

	// Preallocate expected ret slice.
	values := make([]T, 0, len(keys))

	// Acquire lock.
	c.mutex.Lock()
	defer c.mutex.Unlock()

	// Check cache init.
	if c.copy == nil {
		panic("not initialized")
	}

	for i := range keys {
		// Concatenate all *values* from cached items.
		index.get(keys[i].key, func(item *indexed_item) {
			if value, ok := item.data.(T); ok {

				// Append value COPY.
				value = c.copy(value)
				values = append(values, value)

				// Push to front of LRU list, USING
				// THE ITEM'S LRU ENTRY, NOT THE
				// INDEX KEY ENTRY. VERY IMPORTANT!!
				c.lru.move_front(&item.elem)
			}
		})
	}

	return values
}

// Put will insert the given values into cache,
// calling any invalidate hook on each value.
func (c *Cache[T]) Put(values ...T) {
	if len(values) < 1 {
		return
	}

	// Acquire buffer.
	buf := new_buffer()

	// Acquire lock w/ safe unlock.
	unlock := c.mutex.SafeLock()
	defer unlock()

	// Check cache init.
	if c.copy == nil {
		panic("not initialized")
	}

	// Store all passed values.
	for i := range values {
		c.store_value(
			buf,
			nil,
			"",
			values[i],
		)
	}

	// Get func ptrs.
	invalid := c.invalid

	// Done w/
	// lock.
	unlock()

	// Done w/ buffer.
	free_buffer(buf)

	if invalid != nil {
		// Pass all invalidated values
		// to given user hook (if set).
		for _, value := range values {
			invalid(value)
		}
	}
}

// LoadOne fetches one result from the cache stored under index, using precalculated index key.
// In the case that no result is found, provided load callback will be used to hydrate the cache.
func (c *Cache[T]) LoadOne(index *Index, key Key, load func() (T, error)) (T, error) {
	if index == nil {
		panic("no index given")
	} else if index.ptr != unsafe.Pointer(c) {
		panic("invalid index for cache")
	} else if !is_unique(index.flags) {
		panic("cannot get one by non-unique index")
	}

	var (
		// whether an item was found
		// (and so val / err are set).
		ok bool

		// separate value / error ptrs
		// as the item is liable to
		// change outside of lock.
		val T
		err error
	)

	// Acquire lock w/ safe unlock,
	// note the wrapping function.
	unlock := c.mutex.SafeLock()
	defer func() { unlock() }()

	// Check init'd.
	if c.copy == nil ||
		c.ignore == nil {
		panic("not initialized")
	}

	// Get item indexed at key.
	item := index.get_one(key)

	if ok = (item != nil); ok {
		var is bool

		if val, is = item.data.(T); is {
			// Set value COPY.
			val = c.copy(val)

		} else {

			// Attempt to return error.
			err, _ = item.data.(error)
		}

		// Push to front of LRU list, USING
		// THE ITEM'S LRU ENTRY, NOT THE
		// INDEX KEY ENTRY. VERY IMPORTANT!!
		c.lru.move_front(&item.elem)
	}

	// Get func ptrs.
	ignore := c.ignore

	// Done w/
	// lock.
	unlock()

	if ok {
		// item found!
		return val, err
	}

	// Load new result.
	val, err = load()

	// Check for ignored error types.
	if err != nil && ignore(err) {
		return val, err
	}

	// Acquire buffer.
	buf := new_buffer()

	// Reacquire the mutex lock.
	unlock = c.mutex.SafeLock()

	// Index this new loaded item.
	// Note this handles copying of
	// the provided value, so it is
	// safe for us to return as-is.
	if err != nil {
		c.store_error(index, key.key, err)
	} else {
		c.store_value(buf, index, key.key, val)
	}

	// Done w/
	// lock.
	unlock()

	// Done w/ buffer.
	free_buffer(buf)

	return val, err
}

// Load fetches values from the cache stored under index, using precalculated index keys. The cache will attempt to
// results with values stored under keys, passing keys with uncached results to the provider load callback to further
// hydrate the cache with missing results. Cached error results not included or returned by this function.
func (c *Cache[T]) Load(index *Index, keys []Key, load func([]Key) ([]T, error)) ([]T, error) {
	if index == nil {
		panic("no index given")
	} else if index.ptr != unsafe.Pointer(c) {
		panic("invalid index for cache")
	} else if len(keys) < 1 {
		return nil, nil
	}

	// Preallocate expected ret slice.
	values := make([]T, 0, len(keys))

	// Acquire lock w/ safe unlock,
	// note the wrapping function.
	unlock := c.mutex.SafeLock()
	defer func() { unlock() }()

	// Check init'd.
	if c.copy == nil {
		panic("not initialized")
	}

	// Iterate keys and catch uncached.
	toLoad := make([]Key, 0, len(keys))
	for _, key := range keys {

		// Value length before
		// any below appends.
		before := len(values)

		// Concatenate all *values* from cached items.
		index.get(key.key, func(item *indexed_item) {
			if value, ok := item.data.(T); ok {
				// Append value COPY.
				value = c.copy(value)
				values = append(values, value)

				// Push to front of LRU list, USING
				// THE ITEM'S LRU ENTRY, NOT THE
				// INDEX KEY ENTRY. VERY IMPORTANT!!
				c.lru.move_front(&item.elem)
			}
		})

		// Only if values changed did
		// we actually find anything.
		if len(values) == before {
			toLoad = append(toLoad, key)
		}
	}

	// Done w/
	// lock.
	unlock()

	if len(toLoad) == 0 {
		// We loaded everything!
		return values, nil
	}

	// Load uncached key values.
	uncached, err := load(toLoad)
	if err != nil {
		return values, err
	}

	// Acquire buffer.
	buf := new_buffer()

	// Reacquire the mutex lock.
	unlock = c.mutex.SafeLock()

	// Store all uncached values.
	for i := range uncached {
		c.store_value(
			buf,
			nil,
			"",
			uncached[i],
		)
	}

	// Done w/
	// lock.
	unlock()

	// Done w/ buffer.
	free_buffer(buf)

	// Append uncached to return values.
	values = append(values, uncached...)

	return values, nil
}

// Store will call the given store callback, on non-error then
// passing the provided value to the Put() function. On error
// return the value is still passed to stored invalidate hook.
func (c *Cache[T]) Store(value T, store func() error) error {
	// Store value.
	err := store()
	if err != nil {

		// Get func ptrs.
		c.mutex.Lock()
		invalid := c.invalid
		c.mutex.Unlock()

		// On error don't store
		// value, but still pass
		// to invalidate hook.
		if invalid != nil {
			invalid(value)
		}

		return err
	}

	// Store value.
	c.Put(value)

	return nil
}

// Invalidate invalidates all results stored under index keys.
// Note that if set, this will call the invalidate hook on each.
func (c *Cache[T]) Invalidate(index *Index, keys ...Key) {
	if index == nil {
		panic("no index given")
	} else if index.ptr != unsafe.Pointer(c) {
		panic("invalid index for cache")
	}

	// Acquire lock w/ safe unlock.
	unlock := c.mutex.SafeLock()
	defer unlock()

	// Preallocate expected ret slice.
	values := make([]T, 0, len(keys))

	for i := range keys {
		// Delete all items under key from index, collecting
		// value items and dropping them from all their indices.
		index.delete(keys[i].key, func(item *indexed_item) {

			if value, ok := item.data.(T); ok {
				// No need to copy, as item
				// being deleted from cache.
				values = append(values, value)
			}

			// Delete item.
			c.delete(item)
		})
	}

	// Get func ptrs.
	invalid := c.invalid

	// Done w/
	// lock.
	unlock()

	if invalid != nil {
		// Pass all invalidated values
		// to given user hook (if set).
		for _, value := range values {
			invalid(value)
		}
	}
}

// Trim will truncate the cache to ensure it
// stays within given percentage of MaxSize.
func (c *Cache[T]) Trim(perc float64) {

	// Acquire lock.
	c.mutex.Lock()
	defer c.mutex.Unlock()

	// Calculate number of cache items to drop.
	max := (perc / 100) * float64(c.maxSize)
	diff := c.lru.len - int(max)
	if diff <= 0 {

		// Trim not
		// needed.
		return
	}

	// Iterate over 'diff' items
	// from back (oldest) of cache.
	for i := 0; i < diff; i++ {

		// Get oldest LRU elem.
		oldest := c.lru.tail
		if oldest == nil {

			// reached
			// end.
			break
		}

		// Drop oldest item from cache.
		item := (*indexed_item)(oldest.data)
		c.delete(item)
	}

	// Compact index data stores.
	for i := range c.indices {
		c.indices[i].data.Compact()
	}
}

// Clear empties the cache by calling .Trim(0).
func (c *Cache[T]) Clear() { c.Trim(0) }

// Len returns the current length of cache.
func (c *Cache[T]) Len() int {
	c.mutex.Lock()
	l := c.lru.len
	c.mutex.Unlock()
	return l
}

// Debug returns debug stats about cache.
func (c *Cache[T]) Debug() map[string]any {
	m := make(map[string]any, 2)
	c.mutex.Lock()
	m["lru"] = c.lru.len
	indices := make(map[string]any, len(c.indices))
	m["indices"] = indices
	for _, idx := range c.indices {
		var n uint64
		for _, l := range idx.data.m {
			n += uint64(l.len)
		}
		indices[idx.name] = n
	}
	c.mutex.Unlock()
	return m
}

// Cap returns the maximum capacity (size) of cache.
func (c *Cache[T]) Cap() int {
	c.mutex.Lock()
	m := c.maxSize
	c.mutex.Unlock()
	return m
}

func (c *Cache[T]) store_value(
	// key generation buffer.
	buf *byteutil.Buffer,

	// optional index to
	// already store in.
	index *Index,
	key string,

	// value being
	// stored.
	value T,
) {
	// Alloc new index item.
	item := new_indexed_item()
	if cap(item.indexed) < len(c.indices) {

		// Preallocate item indices slice to prevent Go auto
		// allocating overlying large slices we don't need.
		item.indexed = make([]*index_entry, 0, len(c.indices))
	}

	// Create COPY of value.
	value = c.copy(value)
	item.data = value

	if index != nil {
		// Append item to index a key
		// was already generated for.
		evicted := index.append(key, item)
		if evicted != nil {

			// This item is no longer
			// indexed, remove from list.
			c.lru.remove(&evicted.elem)
			free_indexed_item(evicted)
		}
	}

	// Get pointer to value data.
	ptr := unsafe.Pointer(&value)
	for i := range c.indices {

		// Get current index ptr.
		idx := (&c.indices[i])
		if idx == index {

			// Already stored under
			// this index, ignore.
			continue
		}

		// Calculate key for this index
		// for this value by extracting
		// its struct fields and mangling.
		key := idx.extract_key(buf, ptr)
		if key == "" {
			continue
		}

		// Append item to this index.
		evicted := idx.append(key, item)
		if evicted != nil {

			// This item is no longer
			// indexed, remove from list.
			c.lru.remove(&evicted.elem)
			free_indexed_item(evicted)
		}
	}

	if len(item.indexed) == 0 {
		// Item was not stored under
		// any index. Drop this item.
		free_indexed_item(item)
		return
	}

	// Add item to main lru list.
	c.lru.push_front(&item.elem)

	if c.lru.len > c.maxSize {
		// Cache has hit max size!
		// Drop the oldest element.
		ptr := c.lru.tail.data
		item := (*indexed_item)(ptr)
		c.delete(item)
	}
}

func (c *Cache[T]) store_error(index *Index, key string, err error) {
	if index == nil {
		// nothing we
		// can do here.
		return
	}

	// Alloc new index item.
	item := new_indexed_item()
	if cap(item.indexed) < len(c.indices) {

		// Preallocate item indices slice to prevent Go auto
		// allocating overlying large slices we don't need.
		item.indexed = make([]*index_entry, 0, len(c.indices))
	}

	// Set error val.
	item.data = err

	// Append item to index a key
	// was already generated for.
	evicted := index.append(key, item)
	if evicted != nil {

		// This item is no longer
		// indexed, remove from list.
		c.lru.remove(&evicted.elem)
		free_indexed_item(evicted)
	}

	// Add item to main lru list.
	c.lru.push_front(&item.elem)

	if c.lru.len > c.maxSize {
		// Cache has hit max size!
		// Drop the oldest element.
		ptr := c.lru.tail.data
		item := (*indexed_item)(ptr)
		c.delete(item)
	}
}

func (c *Cache[T]) delete(i *indexed_item) {
	for len(i.indexed) > 0 {
		// Pop last indexed entry from list.
		entry := i.indexed[len(i.indexed)-1]
		i.indexed[len(i.indexed)-1] = nil
		i.indexed = i.indexed[:len(i.indexed)-1]

		// Delete entry from
		// its storing index.
		entry.delete_self()

		// Free entry to pool.
		free_index_entry(entry)
	}

	// Drop from lru list.
	c.lru.remove(&i.elem)

	// Free unused item.
	free_indexed_item(i)
}
