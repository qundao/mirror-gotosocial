package cache

import (
	"time"

	"codeberg.org/gruf/go-cache/v3/simple"
	"codeberg.org/gruf/go-cache/v3/ttl"
	"codeberg.org/gruf/go-storage/s3"
)

// check interface conformity.
var _ s3.EntryCache = &EntryCache{}
var _ s3.EntryCache = &EntryTTLCache{}

// EntryCache provides a basic implementation
// of an s3.EntryCache{}. Under the hood it is
// a mutex locked ordered map with max capacity.
type EntryCache struct {
	simple.Cache[string, *s3.CachedObjectInfo]
}

func New(len, cap int) *EntryCache {
	var cache EntryCache
	cache.Init(len, cap)
	return &cache
}

func (c *EntryCache) Put(key string, info *s3.CachedObjectInfo) {
	c.Cache.Set(key, info)
}

type EntryTTLCache struct {
	ttl.Cache[string, *s3.CachedObjectInfo]
}

func NewTTL(len, cap int, ttl time.Duration) *EntryTTLCache {
	var cache EntryTTLCache
	cache.Init(len, cap, ttl)
	return &cache
}

func (c *EntryTTLCache) Put(key string, info *s3.CachedObjectInfo) {
	c.Cache.Set(key, info)
}
