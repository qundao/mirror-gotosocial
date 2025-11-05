package mmap

import (
	"io/fs"
	"path"
	"syscall"
	"time"
)

type fileStat struct {
	syscall.Stat_t
	name string
	mode fs.FileMode
}

func (s *fileStat) Name() string       { return s.name }
func (s *fileStat) IsDir() bool        { return s.mode.IsDir() }
func (s *fileStat) Mode() fs.FileMode  { return s.mode }
func (s *fileStat) Size() int64        { return s.Stat_t.Size }
func (s *fileStat) ModTime() time.Time { return time.Unix(s.Stat_t.Mtim.Unix()) }
func (s *fileStat) Sys() any           { return &s.Stat_t }

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
		return syscall.Stat(filepath, &stat.Stat_t)
	})
	if err != nil {
		return nil, err
	}
	stat.name = path.Base(filepath)
	stat.mode = fs.FileMode(stat.Stat_t.Mode & 0777)
	switch stat.Stat_t.Mode & syscall.S_IFMT {
	case syscall.S_IFBLK:
		stat.mode |= fs.ModeDevice
	case syscall.S_IFCHR:
		stat.mode |= fs.ModeDevice | fs.ModeCharDevice
	case syscall.S_IFDIR:
		stat.mode |= fs.ModeDir
	case syscall.S_IFIFO:
		stat.mode |= fs.ModeNamedPipe
	case syscall.S_IFLNK:
		stat.mode |= fs.ModeSymlink
	case syscall.S_IFREG:
		// nothing to do
	case syscall.S_IFSOCK:
		stat.mode |= fs.ModeSocket
	}
	if stat.Stat_t.Mode&syscall.S_ISGID != 0 {
		stat.mode |= fs.ModeSetgid
	}
	if stat.Stat_t.Mode&syscall.S_ISUID != 0 {
		stat.mode |= fs.ModeSetuid
	}
	if stat.Stat_t.Mode&syscall.S_ISVTX != 0 {
		stat.mode |= fs.ModeSticky
	}
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
