package s3

import (
	"strings"

	"codeberg.org/gruf/go-storage"
	"codeberg.org/gruf/go-storage/internal"
	"github.com/minio/minio-go/v7"
)

// CachedErrorResponse can be returned
// when an S3 is configured with caching,
// and the basic details of an error
// response have been stored in the cache.
type CachedErrorResponse struct {
	Code string
	Key  string
}

func (err *CachedErrorResponse) Error() string {
	return "cached '" + err.Code + "' response for key:" + err.Key
}

func (err *CachedErrorResponse) Is(other error) bool {
	switch other {
	case storage.ErrNotFound:
		return err.Code == "NoSuchKey"
	case storage.ErrAlreadyExists:
		return err.Code == "Conflict"
	default:
		return false
	}
}

func isNotFoundError(err error) bool {
	errRsp, ok := err.(minio.ErrorResponse)
	return ok && errRsp.Code == "NoSuchKey"
}

func isConflictError(err error) bool {
	errRsp, ok := err.(minio.ErrorResponse)
	return ok && errRsp.Code == "Conflict"
}

func isObjectNameError(err error) bool {
	return strings.HasPrefix(err.Error(), "Object name ")
}

func cachedNotFoundError(key string) error {
	err := CachedErrorResponse{Code: "NoSuchKey", Key: key}
	return internal.WrapErr(&err, storage.ErrNotFound)
}
