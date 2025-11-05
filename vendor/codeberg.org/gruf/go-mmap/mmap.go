package mmap

import (
	"errors"
	"io"
	"io/fs"
	"runtime"
	"syscall"
)

// MmapFile maps file at path into memory using syscall.mmap(),
// and returns a protected MmapReader{} for accessing the mapped data.
// Note that the mapped memory is not concurrency safe (other than
// concurrent ReadAt() calls). Any other calls made concurrently to
// Read() or Close() (including ReadAt()) require protection.
func MmapFile(path string) (*MmappedFile, error) {

	// Stat file information.
	stat, err := stat(path)
	if err != nil {
		return nil, err
	}

	// Mmap file into memory.
	return openMmap(path, stat)
}

func openMmap(path string, stat *fileStat) (*MmappedFile, error) {
	if stat.Size() <= 0 {
		// Empty file, no-op read.
		return &MmappedFile{}, nil
	}

	// Check file data size is accessible.
	if stat.Size() != int64(int(stat.Size())) {
		return nil, errors.New("file is too large")
	}

	// Open file at path for read-only access.
	fd, err := open(path, syscall.O_RDONLY, 0)
	if err != nil {
		return nil, err
	}

	// Map this file into memory as slice.
	mem, err := mmap(fd, 0, int(stat.Size()),
		syscall.PROT_READ, syscall.MAP_PRIVATE)

	// Done with file.
	_ = close_(fd)

	if err != nil {
		return nil, err
	}

	// Return as wrapped reader type.
	return newMmapReader(mem, stat), nil
}

// newMmapReader wraps a mapped memory slice in an
// MmappedFile{}, also setting a GC finalizer function.
func newMmapReader(mem []byte, stat *fileStat) *MmappedFile {
	r := &MmappedFile{b: mem, s: stat}
	runtime.SetFinalizer(r, (*MmappedFile).Close)
	return r
}

type MmappedFile struct {
	b []byte    // mapped memory
	n int       // read index
	s *fileStat // file info
}

func (r *MmappedFile) Name() string {
	return r.s.name
}

func (r *MmappedFile) Stat() (fs.FileInfo, error) {
	return r.s, nil
}

func (r *MmappedFile) Read(b []byte) (n int, err error) {
	if r.n >= len(r.b) {
		return 0, io.EOF
	}
	n = copy(b, r.b[r.n:])
	r.n += n
	return
}

func (r *MmappedFile) ReadAt(b []byte, off int64) (n int, err error) {
	if off > int64(len(r.b)) {
		return 0, io.EOF
	}
	n = copy(b, r.b[off:])
	return n, nil
}

func (r *MmappedFile) WriteTo(w io.Writer) (int64, error) {
	if r.n >= len(r.b) {
		return 0, io.EOF
	}
	n, err := w.Write(r.b[r.n:])
	r.n += n
	return int64(n), err
}

func (r *MmappedFile) Seek(off int64, whence int) (int64, error) {
	var n int
	switch whence {
	case io.SeekCurrent:
		n = r.n + int(off)
	case io.SeekStart:
		n = 0 + int(off)
	case io.SeekEnd:
		n = len(r.b) + int(off)
	default:
		return 0, errors.New("invalid argument")
	}
	if n < 0 || n > len(r.b) {
		return 0, errors.New("invalid argument")
	}
	r.n = n
	return int64(n), nil
}

func (r *MmappedFile) Len() int {
	return len(r.b) - r.n
}

func (r *MmappedFile) Size() int64 {
	return int64(len(r.b))
}

func (r *MmappedFile) Close() error {
	if b := r.b; b != nil {
		r.b = nil
		runtime.SetFinalizer(r, nil)
		return munmap(b)
	}
	return nil
}
