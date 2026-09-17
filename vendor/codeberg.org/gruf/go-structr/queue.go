package structr

import (
	"unsafe"

	"codeberg.org/gruf/go-byteutil"
)

// QueueConfig defines config vars
// for initializing a struct queue.
type QueueConfig[StructType any] struct {

	// Pop is called when queue values
	// are popped, during calls to any
	// of the Pop___() series of fns.
	Pop func(StructType)

	// Indices defines indices to create
	// in the Queue for the receiving
	// generic struct parameter type.
	Indices []IndexConfig
}

// Queue provides a structure model queue with
// automated indexing and popping by any init
// defined lookups of field combinations.
type Queue[StructType any] struct {

	// hook functions.
	pop func(StructType)

	// main underlying
	// struct item queue.
	queue list

	// indices used in storing passed struct
	// types by user defined sets of fields.
	indices []Index

	// protective mutex, guards:
	// - Queue{}.queue
	// - Index{}.data
	// - Queue{} hook fns
	mutex mutex
}

// Init initializes the queue with given configuration
// including struct fields to index, and necessary fns.
func (q *Queue[T]) Init(config QueueConfig[T]) {
	t := get_type_iter[T]()

	if len(config.Indices) == 0 {
		panic("no indices provided")
	}

	// Acquire lock.
	q.mutex.Lock()
	defer q.mutex.Unlock()

	// Ensure not setup.
	if q.indices != nil {
		panic("already initialized")
	}

	// Prepare each of configured struct indices.
	q.indices = make([]Index, len(config.Indices))
	for i, cfg := range config.Indices {
		q.indices[i].ptr = unsafe.Pointer(q)
		q.indices[i].init(t, cfg, 0)
	}

	// Copy functions.
	q.pop = config.Pop
}

// Index selects index with given name from queue, else panics.
func (q *Queue[T]) Index(name string) *Index {
	for i := range q.indices {
		if q.indices[i].name == name {
			return &(q.indices[i])
		}
	}
	panic("unknown index: " + name)
}

// PopFront pops the current value at front of the queue.
func (q *Queue[T]) PopFront() (T, bool) {
	t := q.PopFrontN(1)
	if len(t) == 0 {
		var t T
		return t, false
	}
	return t[0], true
}

// PopBack pops the current value at back of the queue.
func (q *Queue[T]) PopBack() (T, bool) {
	t := q.PopBackN(1)
	if len(t) == 0 {
		var t T
		return t, false
	}
	return t[0], true
}

// PopFrontN attempts to pop n values from front of the queue.
func (q *Queue[T]) PopFrontN(n int) []T {
	return q.pop_n(n, func() *list_elem {
		return q.queue.head
	})
}

// PopBackN attempts to pop n values from back of the queue.
func (q *Queue[T]) PopBackN(n int) []T {
	return q.pop_n(n, func() *list_elem {
		return q.queue.tail
	})
}

// Pop attempts to pop values from queue indexed under any of keys.
func (q *Queue[T]) Pop(index *Index, keys ...Key) []T {
	if index == nil {
		panic("no index given")
	} else if index.ptr != unsafe.Pointer(q) {
		panic("invalid index for queue")
	}

	// Acquire lock w/ safe unlock.
	unlock := q.mutex.SafeLock()
	defer unlock()

	// Preallocate expected ret slice.
	values := make([]T, 0, len(keys))

	for i := range keys {
		// Delete all items under key from index, collecting
		// value items and dropping them from all their indices.
		index.delete(keys[i].key, func(item *indexed_item) {

			// Append deleted to values.
			value := item.data.(T)
			values = append(values, value)

			// Delete item.
			q.delete(item)
		})
	}

	// Get func ptrs.
	pop := q.pop

	// Done w/
	// lock.
	unlock()

	if pop != nil {
		// Pass all popped values
		// to given user hook (if set).
		for _, value := range values {
			pop(value)
		}
	}

	return values
}

// PushFront pushes values to front of queue.
func (q *Queue[T]) PushFront(values ...T) {
	q.mutex.Lock()
	defer q.mutex.Unlock()
	buf := new_buffer()
	for i := range values {
		item := q.index(buf, values[i])
		q.queue.push_front(&item.elem)
	}
	free_buffer(buf)
}

// PushBack pushes values to back of queue.
func (q *Queue[T]) PushBack(values ...T) {
	q.mutex.Lock()
	defer q.mutex.Unlock()
	buf := new_buffer()
	for i := range values {
		item := q.index(buf, values[i])
		q.queue.push_back(&item.elem)
	}
	free_buffer(buf)
}

// MoveFront attempts to move values indexed under any of keys to the front of the queue.
func (q *Queue[T]) MoveFront(index *Index, keys ...Key) {
	q.mutex.Lock()
	defer q.mutex.Unlock()
	for i := range keys {
		index.get(keys[i].key, func(item *indexed_item) {
			q.queue.move_front(&item.elem)
		})
	}
}

// MoveBack attempts to move values indexed under any of keys to the back of the queue.
func (q *Queue[T]) MoveBack(index *Index, keys ...Key) {
	q.mutex.Lock()
	defer q.mutex.Unlock()
	for i := range keys {
		index.get(keys[i].key, func(item *indexed_item) {
			q.queue.move_back(&item.elem)
		})
	}
}

// Count returns how many values are indexed under any of keys.
func (q *Queue[T]) Count(index *Index, keys ...Key) (n int) {
	q.mutex.Lock()
	defer q.mutex.Unlock()
	for i := range keys {
		index.get(keys[i].key, func(*indexed_item) { n++ })
	}
	return
}

// Len returns the current length of queue.
func (q *Queue[T]) Len() int {
	q.mutex.Lock()
	l := q.queue.len
	q.mutex.Unlock()
	return l
}

// Debug returns debug stats about queue.
func (q *Queue[T]) Debug() map[string]any {
	m := make(map[string]any, 2)
	q.mutex.Lock()
	m["queue"] = q.queue.len
	indices := make(map[string]any, len(q.indices))
	m["indices"] = indices
	for _, idx := range q.indices {
		var n uint64
		for _, l := range idx.data.m {
			n += uint64(l.len)
		}
		indices[idx.name] = n
	}
	q.mutex.Unlock()
	return m
}

func (q *Queue[T]) pop_n(n int, next func() *list_elem) []T {
	if next == nil {
		panic("nil fn")
	}

	// Acquire lock w/ safe unlock.
	unlock := q.mutex.SafeLock()
	defer unlock()

	// Preallocate ret slice.
	values := make([]T, 0, n)

	// Iterate over 'n' items.
	for i := 0; i < n; i++ {

		// Get next elem.
		next := next()
		if next == nil {

			// reached
			// end.
			break
		}

		// Cast the indexed item from elem.
		item := (*indexed_item)(next.data)

		// Append deleted to values.
		value := item.data.(T)
		values = append(values, value)

		// Delete item.
		q.delete(item)
	}

	// Get func ptrs.
	pop := q.pop

	// Done w/
	// lock.
	unlock()

	if pop != nil {
		// Pass all popped values
		// to given user hook (if set).
		for _, value := range values {
			pop(value)
		}
	}

	return values
}

func (q *Queue[T]) index(buf *byteutil.Buffer, value T) *indexed_item {
	item := new_indexed_item()
	if cap(item.indexed) < len(q.indices) {

		// Preallocate item indices slice to prevent Go auto
		// allocating overlying large slices we don't need.
		item.indexed = make([]*index_entry, 0, len(q.indices))
	}

	// Set item value.
	item.data = value

	// Get pointer to value data.
	ptr := unsafe.Pointer(&value)
	for i := range q.indices {

		// Get current index ptr.
		idx := &(q.indices[i])

		// Calculate key for this index
		// for this value by extracting
		// its struct fields and mangling.
		key := idx.extract_key(buf, ptr)
		if key == "" {

			// zero value field
			// w/ !AllowZero.
			continue
		}

		// Append item to this index.
		evicted := idx.append(key, item)
		if evicted != nil {

			// This item is no longer
			// indexed, remove from list.
			q.queue.remove(&evicted.elem)
			free_indexed_item(evicted)
		}
	}

	return item
}

func (q *Queue[T]) delete(i *indexed_item) {
	for len(i.indexed) > 0 {
		// Pop last indexed entry from list.
		entry := i.indexed[len(i.indexed)-1]
		i.indexed[len(i.indexed)-1] = nil
		i.indexed = i.indexed[:len(i.indexed)-1]

		// Delete entry from
		// its storing index.
		entry.delete_self()

		// Check index load factor.
		entry.index.data.Compact()

		// Free entry to pool.
		free_index_entry(entry)
	}

	// Drop from queue list.
	q.queue.remove(&i.elem)

	// Free unused item.
	free_indexed_item(i)
}
