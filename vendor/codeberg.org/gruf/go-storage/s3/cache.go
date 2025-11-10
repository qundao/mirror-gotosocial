package s3

import (
	"time"

	"github.com/minio/minio-go/v7"
)

// EntryCache should provide a cache of
// S3 object information for speeding up
// Get(), Stat() and Remove() operations.
type EntryCache interface {

	// Get should return 'found' = true when information is cached,
	// with 'info' optionally being nilable to allow caching errors.
	Get(key string) (info *CachedObjectInfo, found bool)

	// Put should cache the given information under key, with
	// nil CachedObjectInfo{} meaning a 'not found' response.
	Put(key string, info *CachedObjectInfo)
}

// CachedObjectInfo provides the minimum cacheable
// set of S3 object information that may be returned
// from a Get() or Stat() operation, or on Put().
type CachedObjectInfo struct {
	Key          string
	ETag         string
	Size         int64
	ContentType  string
	LastModified time.Time
	VersionID    string
}

// ToObjectInfo converts CachedObjectInfo to returnable minio.ObjectInfo.
func (info *CachedObjectInfo) ToObjectInfo() minio.ObjectInfo {
	return minio.ObjectInfo{
		Key:          info.Key,
		ETag:         info.ETag,
		Size:         info.Size,
		ContentType:  info.ContentType,
		LastModified: info.LastModified,
		VersionID:    info.VersionID,
	}
}

// cacheGet wraps cache.Get() operations to check if cache is nil.
func cacheGet(cache EntryCache, key string) (*CachedObjectInfo, bool) {
	if cache != nil {
		return cache.Get(key)
	}
	return nil, false
}

// objectToCachedObjectInfo converts minio.ObjectInfo to CachedObjectInfo for caching.
func objectToCachedObjectInfo(info minio.ObjectInfo) *CachedObjectInfo {
	return &CachedObjectInfo{
		Key:          info.Key,
		ETag:         info.ETag,
		Size:         info.Size,
		ContentType:  info.ContentType,
		LastModified: info.LastModified,
		VersionID:    info.VersionID,
	}
}

// cacheObject wraps cache.Put() operations to check if cache is nil.
func cacheObject(cache EntryCache, key string, info minio.ObjectInfo) {
	if cache != nil {
		cache.Put(key, objectToCachedObjectInfo(info))
	}
}

// cacheUpload wraps cache.Put() operations to check if cache is nil, uses ContentType from given opts.
func cacheUpload(cache EntryCache, key string, info minio.UploadInfo, opts minio.PutObjectOptions) {
	if cache != nil {
		cache.Put(key, &CachedObjectInfo{
			Key:          info.Key,
			ETag:         info.ETag,
			Size:         info.Size,
			ContentType:  opts.ContentType,
			LastModified: info.LastModified,
			VersionID:    info.VersionID,
		})
	}
}

// cacheNotFound wraps cache.Put() to check if cache is
// nil, storing a nil entry (i.e. not found) in cache.
func cacheNotFound(cache EntryCache, key string) {
	if cache != nil {
		cache.Put(key, nil)
	}
}
