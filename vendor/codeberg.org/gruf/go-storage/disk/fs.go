package disk

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"syscall"

	"codeberg.org/gruf/go-fastpath/v2"
	"codeberg.org/gruf/go-storage/internal"
)

// open file for read args.
var readArgs = OpenArgs{
	Flags: syscall.O_RDONLY,
	Perms: 0,
}

// walkDir traverses the dir tree of the supplied path, performing the supplied walkFn on each entry.
func walkDir(pb *fastpath.Builder, path string, walkFn func(string, fs.DirEntry) error) error {

	// Read directory entries at path.
	entries, err := readDir(path)
	if err != nil {
		return err
	}

	// frame represents a directory entry
	// walk-loop snapshot, taken when a sub
	// directory requiring iteration is found
	type frame struct {
		path    string
		entries []fs.DirEntry
	}

	// stack contains a list of held snapshot
	// frames, representing unfinished upper
	// layers of a directory structure yet to
	// be traversed.
	var stack []frame

outer:
	for {
		if len(entries) == 0 {
			if len(stack) == 0 {
				// Reached end
				break outer
			}

			// Pop frame from stack
			frame := stack[len(stack)-1]
			stack = stack[:len(stack)-1]

			// Update loop vars
			entries = frame.entries
			path = frame.path
		}

		for len(entries) > 0 {
			// Pop next entry from queue
			entry := entries[0]
			entries = entries[1:]

			// Pass to provided walk function
			if err := walkFn(path, entry); err != nil {
				return err
			}

			if entry.IsDir() {
				// Push current frame to stack
				stack = append(stack, frame{
					path:    path,
					entries: entries,
				})

				// Update current directory path
				path = pb.Join(path, entry.Name())

				// Read next directory entries
				next, err := readDir(path)
				if err != nil {
					return err
				}

				// Set next entries
				entries = next

				continue outer
			}
		}
	}

	return nil
}

// cleanDirs traverses the dir tree of supplied
// path, removing any folders with zero children.
func cleanDirs(path string) error {
	pb := internal.GetPathBuilder()
	err := cleanDir(pb, path, true)
	internal.PutPathBuilder(pb)
	return err
}

// cleanDir performs the actual dir cleaning logic for the above top-level version.
func cleanDir(pb *fastpath.Builder, path string, top bool) error {

	// Get directory entries at path.
	entries, err := readDir(path)
	if err != nil {
		return err
	}

	// If no entries, delete dir.
	if !top && len(entries) == 0 {
		return rmdir(path)
	}

	var errs []error

	// Iterate all directory entries.
	for _, entry := range entries {

		if entry.IsDir() {
			// Calculate directory path.
			dir := pb.Join(path, entry.Name())

			// Recursively clean sub-dir entries, adding errs.
			if err := cleanDir(pb, dir, false); err != nil {
				err = fmt.Errorf("error(s) cleaning subdir %s: %w", dir, err)
				errs = append(errs, err)
			}
		}
	}

	// Return combined errors.
	return errors.Join(errs...)
}

// readDir will open file at path, read the unsorted list of entries, then close.
func readDir(path string) ([]fs.DirEntry, error) {

	// Open directory at path for read.
	file, err := open(path, readArgs)
	if err != nil {
		return nil, err
	}

	// Read ALL directory entries.
	entries, err := file.ReadDir(-1)

	// Done with file
	_ = file.Close()

	return entries, err
}

// open is a simple wrapper around syscall.Open().
func open(path string, args OpenArgs) (*os.File, error) {
	var fd int
	err := retryOnEINTR(func() (err error) {
		fd, err = syscall.Open(path, args.Flags, args.Perms)
		return
	})
	if err != nil {
		return nil, err
	}
	return os.NewFile(uintptr(fd), path), nil
}

// stat is a simple wrapper around syscall.Stat().
func stat(path string) (*syscall.Stat_t, error) {
	var stat syscall.Stat_t
	err := retryOnEINTR(func() error {
		return syscall.Stat(path, &stat)
	})
	if err != nil {
		if err == syscall.ENOENT {
			// not-found is no error
			err = nil
		}
		return nil, err
	}
	return &stat, nil
}

// lstat is a simple wrapper around syscall.Lstat().
func lstat(path string) (*syscall.Stat_t, error) {
	var stat syscall.Stat_t
	err := retryOnEINTR(func() error {
		return syscall.Lstat(path, &stat)
	})
	if err != nil {
		if err == syscall.ENOENT {
			// not-found is no error
			err = nil
		}
		return nil, err
	}
	return &stat, nil
}

// unlink is a simple wrapper around syscall.Unlink().
func unlink(path string) error {
	return retryOnEINTR(func() error {
		return syscall.Unlink(path)
	})
}

// rmdir is a simple wrapper around syscall.Rmdir().
func rmdir(path string) error {
	return retryOnEINTR(func() error {
		return syscall.Rmdir(path)
	})
}

// symlink is a simple wrapper around syscall.Symlink()
func symlink(oldpath, newpath string) error {
	return retryOnEINTR(func() error {
		return syscall.Symlink(oldpath, newpath)
	})
}

// link is a simple wrapper around syscall.Link()
func link(oldpath, newpath string) error {
	return retryOnEINTR(func() error {
		return syscall.Link(oldpath, newpath)
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
