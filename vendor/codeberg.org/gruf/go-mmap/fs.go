package mmap

import (
	"io/fs"
	"syscall"
	"time"
)

type fileStat struct {
	name    string
	size    int64
	mode    fs.FileMode
	modTime time.Time
	sys     syscall.Stat_t
}

func (fs *fileStat) Name() string       { return fs.name }
func (fs *fileStat) Size() int64        { return fs.size }
func (fs *fileStat) IsDir() bool        { return fs.mode.IsDir() }
func (fs *fileStat) Mode() fs.FileMode  { return fs.mode }
func (fs *fileStat) ModTime() time.Time { return fs.modTime }
func (fs *fileStat) Sys() any           { return &fs.sys }

// open is a simple wrapper around syscall.Open().
func open(filepath string, mode int, perm uint32) (fd int, err error) {
	err = retryOnEINTR(func() (err error) {
		fd, err = syscall.Open(filepath, mode, perm)
		return
	})
	return
}

// stat is a simple wrapper around syscall.Stat().
func stat(filepath string) (*fileStat, error) {
	var stat fileStat
	err := retryOnEINTR(func() error {
		return syscall.Stat(filepath, &stat.sys)
	})
	if err != nil {
		return nil, err
	}
	fillFileStatFromSys(&stat, filepath)
	return &stat, nil
}

// mmap is a simple wrapper around syscall.Mmap().
func mmap(fd int, offset int64, length int, prot int, flags int) (b []byte, err error) {
	err = retryOnEINTR(func() error {
		b, err = syscall.Mmap(fd, offset, length, prot, flags)
		return err
	})
	return
}

// munmap is a simple wrapper around syscall.Munmap().
func munmap(b []byte) error {
	return retryOnEINTR(func() error {
		return syscall.Munmap(b)
	})
}

// close_ is a simple wrapper around syscall.Close().
func close_(fd int) error {
	return retryOnEINTR(func() error {
		return syscall.Close(fd)
	})
}

// retryOnEINTR is a low-level filesystem function
// for retrying syscalls on O_EINTR received.
func retryOnEINTR(do func() error) error {
	for {
		err := do()
		if err == syscall.EINTR {
			continue
		}
		return err
	}
}
