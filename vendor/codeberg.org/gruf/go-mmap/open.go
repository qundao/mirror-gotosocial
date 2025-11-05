package mmap

import (
	"io"
	"io/fs"
	"os"
	"runtime"
	"syscall"
)

// MmapThreshold defines the threshold file size (in bytes) for which
// a call to OpenRead() will deem as big enough for a file to be worth
// opening using an `mmap` syscall. This is a runtime initialized number
// based on the number of available CPUs as in concurrent conditions Go
// can make optimizations for blocking `read` syscalls which scales with
// the number of available goroutines it can have running at once.
var MmapThreshold = int64(runtime.NumCPU() * syscall.Getpagesize())

// FileReader defines the base interface
// of a readable file, whether accessed
// via `read` or `mmap` syscalls.
type FileReader interface {
	fs.File
	io.ReaderAt
	io.WriterTo
	io.Seeker
	Name() string
}

// Threshold is a receiving type for OpenRead()
// that allows defining a custom MmapThreshold.
type Threshold struct{ At int64 }

// OpenRead: see mmap.OpenRead().
func (t Threshold) OpenRead(path string) (FileReader, error) {
	stat, err := stat(path)
	if err != nil {
		return nil, err
	}
	if stat.Size() >= t.At {
		return openMmap(path, stat)
	} else {
		return os.OpenFile(path, syscall.O_RDONLY, 0)
	}
}

// OpenRead will open the file as read only (erroring if it does
// not already exist). If the file at path is beyond 'MmapThreshold'
// it will be opened for reads using an `mmap` syscall, by calling
// MmappedRead(path). Else, it will be opened using os.OpenFile().
//
// Please note that the reader returned by this function is not
// guaranteed to be concurrency-safe. Calls returned by os.OpenFile()
// follow the usual standard library concurrency guarantees, but the
// reader returned by MmappedRead() provides no concurrent protection.
//
// Also note that this may not always be faster! If the file you need
// to open will be immediately drained to another file, TCP or Unix
// connection, then the standard library will used optimized syscalls.
func OpenRead(path string) (FileReader, error) {
	return Threshold{MmapThreshold}.OpenRead(path)
}
