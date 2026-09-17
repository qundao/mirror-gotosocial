# go-structr

A library with a series of performant data types with automated struct value indexing. Indexing is supported via arbitrary combinations of fields, and in the case of the cache type, negative results (errors!) are also supported.

Under the hood, go-structr maintains a hashmap per index, where each hashmap is keyed by serialized input key. This is handled by the incredibly performant serialization library [go-mangler/v2](https://codeberg.org/gruf/go-mangler), which at this point in time supports all concrete types, so feel free to index by by *almost* anything!

See the [docs](https://pkg.go.dev/codeberg.org/gruf/go-structr) for more API information.

## Notes

If I don't reach a v2 of this library beforehand, some nice performance improvements could be found with:
- speed-up multi-index key generation by building a trie at init time so we don't need to re-iterate over indices with shared field prefixes
- exploring alternate data structures for the underlying types, especially for the queue and timeline types for which hashmaps probably aren't ideal
- exploring ways of minimising GC time with so many pointers in the linked lists by either hiding them with uintptrs, or linking to them with fixed indices in an underlying backing data array

This is a core underpinning of [GoToSocial](https://github.com/superseriousbusiness/gotosocial)'s performance.