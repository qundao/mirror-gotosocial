package disk

import (
	"bytes"
	"context"
	"io"
	"io/fs"
	"os"
	"path"
	"strings"
	"syscall"

	"codeberg.org/gruf/go-fastpath/v2"
	"codeberg.org/gruf/go-storage"
	"codeberg.org/gruf/go-storage/internal"
)

// ensure DiskStorage conforms to storage.Storage.
var _ storage.Storage = (*DiskStorage)(nil)

// DefaultConfig returns the default DiskStorage configuration.
func DefaultConfig() Config {
	return defaultConfig
}

// immutable default configuration.
var defaultConfig = Config{
	Create:     OpenArgs{syscall.O_CREAT | syscall.O_WRONLY, 0o644},
	MkdirPerms: 0o755,
	CopyFn:     io.Copy,
}

// OpenArgs defines args passed
// in a syscall.Open() operation.
type OpenArgs struct {
	Flags int
	Perms uint32
}

// Config defines options to be
// used when opening a DiskStorage.
type Config struct {

	// Create are the arguments passed
	// to syscall.Open() when creating
	// a file for write operations.
	Create OpenArgs

	// MkdirPerms are the permissions used
	// when creating necessary sub-dirs in
	// a storage key with slashes.
	MkdirPerms uint32

	// CopyFn ...
	CopyFn func(io.Writer, io.Reader) (int64, error)
}

// getDiskConfig returns valid (and owned!) Config for given ptr.
func getDiskConfig(cfg *Config) Config {
	if cfg == nil {
		// use defaults.
		return defaultConfig
	}

	// Ensure non-zero syscall args.
	if cfg.Create.Flags == 0 {
		cfg.Create.Flags = defaultConfig.Create.Flags
	}
	if cfg.Create.Perms == 0 {
		cfg.Create.Perms = defaultConfig.Create.Perms
	}
	if cfg.MkdirPerms == 0 {
		cfg.MkdirPerms = defaultConfig.MkdirPerms
	}

	return Config{
		Create:     cfg.Create,
		MkdirPerms: cfg.MkdirPerms,
		CopyFn:     cfg.CopyFn,
	}
}

// DiskStorage is a Storage implementation
// that stores directly to a filesystem.
type DiskStorage struct {
	cfg  Config // cfg is the supplied configuration for this store
	path string // path is the root path of this store
}

// Open opens a DiskStorage instance for given folder path and configuration.
func Open(path string, cfg *Config) (*DiskStorage, error) {

	// Check + set config defaults.
	config := getDiskConfig(cfg)

	// Clean provided storage path, ensure
	// final '/' to help with path trimming.
	pb := internal.GetPathBuilder()
	path = pb.Clean(path) + "/"
	internal.PutPathBuilder(pb)

	// Ensure directories up-to path exist.
	perms := fs.FileMode(config.MkdirPerms)
	err := os.MkdirAll(path, perms)
	if err != nil {
		return nil, err
	}

	// Prepare DiskStorage.
	st := &DiskStorage{
		path: path,
		cfg:  config,
	}

	return st, nil
}

// Clean: implements Storage.Clean().
func (st *DiskStorage) Clean(_ context.Context) error {
	return cleanDirs(st.path)
}

// ReadBytes: implements Storage.ReadBytes().
func (st *DiskStorage) ReadBytes(ctx context.Context, key string) ([]byte, error) {

	// Get stream reader for key.
	rc, err := st.ReadStream(ctx, key)
	if err != nil {
		return nil, err
	}

	// Read all data to memory.
	data, err := io.ReadAll(rc)

	// Close the reader.
	err2 := rc.Close()

	if err != nil {
		return nil, err
	} else if err2 != nil {
		return nil, err2
	}

	return data, nil
}

// ReadStream: implements Storage.ReadStream().
func (st *DiskStorage) ReadStream(_ context.Context, key string) (io.ReadCloser, error) {

	// Generate file path for key.
	kpath, err := st.Filepath(key)
	if err != nil {
		return nil, err
	}

	// Attempt to open file for read.
	file, err := open(kpath, readArgs)
	if err != nil {

		if err == syscall.ENOENT {
			// Translate not-found errors and wrap with key.
			err = internal.ErrWithKey(storage.ErrNotFound, key)
		}

		return nil, err
	}

	return file, nil
}

// WriteBytes: implements Storage.WriteBytes().
func (st *DiskStorage) WriteBytes(ctx context.Context, key string, value []byte) (int, error) {
	n, err := st.WriteStream(ctx, key, bytes.NewReader(value))
	return int(n), err
}

// WriteStream: implements Storage.WriteStream().
func (st *DiskStorage) WriteStream(_ context.Context, key string, stream io.Reader) (int64, error) {

	// Acquire path builder buffer.
	pb := internal.GetPathBuilder()

	// Generate the file path for key.
	kpath, err := st.filepath(pb, key)

	// Done with path buffer.
	internal.PutPathBuilder(pb)

	if err != nil {
		return 0, err
	}

	// Fast check for whether this may be a
	// sub-directory. This is not a definitive
	// check, but it indicates to try MkdirAll.
	if strings.ContainsRune(key, '/') {

		// Get dir of key path.
		dir := path.Dir(kpath)

		// Ensure required key path dirs exist.
		perms := fs.FileMode(st.cfg.MkdirPerms)
		err = os.MkdirAll(dir, perms)
		if err != nil {
			return 0, err
		}
	}

	// Attempt to open file with create args.
	file, err := open(kpath, st.cfg.Create)
	if err != nil {

		if st.cfg.Create.Flags&syscall.O_EXCL != 0 &&
			err == syscall.EEXIST {
			// Translate already exists errors and wrap with key.
			err = internal.ErrWithKey(storage.ErrAlreadyExists, key)
		}

		return 0, err
	}

	var n int64

	if st.cfg.CopyFn != nil {
		// Use provided io copy function.
		n, err = st.cfg.CopyFn(file, stream)
	} else {
		// Use default io.Copy func.
		n, err = io.Copy(file, stream)
	}

	if err != nil {
		_ = file.Close()
		return n, err
	}

	// Finally, close file.
	return n, file.Close()
}

// Stat implements Storage.Stat().
func (st *DiskStorage) Stat(_ context.Context, key string) (*storage.Entry, error) {
	stat, err := st.Stat_(key)
	if stat == nil {
		return nil, err
	}
	return &storage.Entry{
		Key:      key,
		Size:     stat.Size,
		Modified: modtime(stat),
	}, nil
}

// Remove implements Storage.Remove().
func (st *DiskStorage) Remove(_ context.Context, key string) error {
	return st.Unlink(key)
}

// WalkKeys implements Storage.WalkKeys().
func (st *DiskStorage) WalkKeys(_ context.Context, opts storage.WalkKeysOpts) error {
	if opts.Step == nil {
		panic("nil step fn")
	}

	// Acquire path builder for walk.
	pb := internal.GetPathBuilder()
	defer internal.PutPathBuilder(pb)

	// Dir to walk.
	dir := st.path

	if opts.Prefix != "" {
		// Convert key prefix to one of our filepaths.
		pathprefix, err := st.filepath(pb, opts.Prefix)
		if err != nil {
			return internal.ErrWithMsg(err, "prefix error")
		}

		// Fast check for whether this may be a
		// sub-directory. This is not a definitive
		// check, but it allows us to update the
		// directory we walk to narrow search params.
		if strings.ContainsRune(opts.Prefix, '/') {
			dir = path.Dir(pathprefix)
		}

		// Set updated storage
		// path prefix in opts.
		opts.Prefix = pathprefix
	}

	return walkDir(pb, dir, func(kpath string, fsentry fs.DirEntry) error {
		if !fsentry.Type().IsRegular() {
			// Ignore anything but
			// regular file types.
			return nil
		}

		// Get full item path (without root).
		kpath = pb.Join(kpath, fsentry.Name())

		// Perform a fast filter check against storage path prefix (if set).
		if opts.Prefix != "" && !strings.HasPrefix(kpath, opts.Prefix) {
			return nil // ignore
		}

		// Storage key without base.
		key := kpath[len(st.path):]

		// Ignore filtered keys.
		if opts.Filter != nil &&
			!opts.Filter(key) {
			return nil // ignore
		}

		// Load file info. This should already
		// be loaded due to the underlying call
		// to os.File{}.ReadDir() populating them.
		info, err := fsentry.Info()
		if err != nil {
			return err
		}

		// Perform provided walk function
		return opts.Step(storage.Entry{
			Modified: info.ModTime(),
			Size:     info.Size(),
			Key:      key,
		})
	})
}

// Filepath checks and returns a formatted Filepath for given key.
func (st *DiskStorage) Filepath(key string) (path string, err error) {
	pb := internal.GetPathBuilder()
	path, err = st.filepath(pb, key)
	internal.PutPathBuilder(pb)
	return
}

// Open performs syscall.Open() on the file in storage at key, with given OpenArgs{}.
//
// NOTE: this does not perform much of the wrapping that os.OpenFile() does, it may
// not set appropriate arguments for opening files other than regular / directories!
func (st *DiskStorage) Open(key string, args OpenArgs) (*os.File, error) {

	// Generate file path for key.
	kpath, err := st.Filepath(key)
	if err != nil {
		return nil, err
	}

	// Open file path with args.
	file, err := open(kpath, args)
	switch err {

	case syscall.ENOENT:
		// Translate not-found errors and wrap with key.
		err = internal.ErrWithKey(storage.ErrNotFound, key)

	case syscall.EEXIST:
		if args.Flags&syscall.O_EXCL != 0 {
			// Translate already exists errors and wrap with key.
			err = internal.ErrWithKey(storage.ErrAlreadyExists, key)
		}
	}

	return file, err
}

// OpenRead calls .Open() with syscall.O_RDONLY OpenArgs{}.
func (st *DiskStorage) OpenRead(key string) (*os.File, error) {
	return st.Open(key, readArgs)
}

// OpenWrite calls .Open() with configured 'Create' OpenArgs{}.
func (st *DiskStorage) OpenWrite(key string) (*os.File, error) {
	return st.Open(key, st.cfg.Create)
}

// ReadDir performs syscall.ReadDir() on the file in storage at key.
func (st *DiskStorage) ReadDir(key string) ([]fs.DirEntry, error) {

	// Generate file path for key.
	kpath, err := st.Filepath(key)
	if err != nil {
		return nil, err
	}

	// Read entries in directory.
	entries, err := readDir(kpath)
	switch err {

	case syscall.ENOENT:
		// Translate not-found errors and wrap with key.
		err = internal.ErrWithKey(storage.ErrNotFound, key)
	}

	return entries, err
}

// Stat_ performs syscall.Stat() on the file in storage at key.
func (st *DiskStorage) Stat_(key string) (*syscall.Stat_t, error) {

	// Generate file path for key.
	kpath, err := st.Filepath(key)
	if err != nil {
		return nil, err
	}

	// Stat file on disk.
	return stat(kpath)
}

// Lstat performs syscall.Lstat() on the file in storage at key.
func (st *DiskStorage) Lstat(key string) (*syscall.Stat_t, error) {

	// Generate file path for key.
	kpath, err := st.Filepath(key)
	if err != nil {
		return nil, err
	}

	// Stat file on disk.
	return lstat(kpath)
}

// Unlink performs syscall.Unlink() on the file in storage at key.
func (st *DiskStorage) Unlink(key string) error {

	// Generate file path for key.
	kpath, err := st.Filepath(key)
	if err != nil {
		return err
	}

	// Remove at path (must be a file).
	if err := unlink(kpath); err != nil {

		if err == syscall.ENOENT {
			// Translate not-found errors and wrap with key.
			err = internal.ErrWithKey(storage.ErrNotFound, key)
		}

		return err
	}

	return nil
}

// Rmdir performs syscall.Rmdir() on the dir in storage at key.
func (st *DiskStorage) Rmdir(key string) error {

	// Generate file path for key.
	kpath, err := st.Filepath(key)
	if err != nil {
		return err
	}

	// Remove at path (must be a dir).
	if err := rmdir(kpath); err != nil {

		if err == syscall.ENOENT {
			// Translate not-found errors and wrap with key.
			err = internal.ErrWithKey(storage.ErrNotFound, key)
		}

		return err
	}

	return nil
}

// Symlink performs syscall.Symlink() on the source and destination keys in storage.
func (st *DiskStorage) Symlink(srcKey, dstKey string) error {

	// Acquire path builder buffer.
	pb := internal.GetPathBuilder()

	// Generate file path for source.
	src, err1 := st.filepath(pb, srcKey)

	// Generate file path for destination.
	dst, err2 := st.filepath(pb, dstKey)

	// Done with path buffer.
	internal.PutPathBuilder(pb)

	if err1 != nil {
		return err1
	} else if err2 != nil {
		return err2
	}

	// Create disk symlink.
	return symlink(src, dst)
}

// Link performs syscall.Link() on the source and destination keys in storage.
func (st *DiskStorage) Link(srcKey, dstKey string) error {

	// Acquire path builder buffer.
	pb := internal.GetPathBuilder()

	// Generate file path for source.
	src, err1 := st.filepath(pb, srcKey)

	// Generate file path for destination.
	dst, err2 := st.filepath(pb, dstKey)

	// Done with path buffer.
	internal.PutPathBuilder(pb)

	if err1 != nil {
		return err1
	} else if err2 != nil {
		return err2
	}

	// Create disk hardlink.
	return link(src, dst)
}

// filepath performs the "meat" of Filepath(), working with an existing fastpath.Builder{}.
func (st *DiskStorage) filepath(pb *fastpath.Builder, key string) (path string, err error) {

	// Build from base.
	pb.Append(st.path)
	pb.Append(key)

	// Take COPY of bytes.
	path = string(pb.B)

	// Check for dir traversal outside base.
	if isDirTraversal(st.path, path) {
		err = internal.ErrWithKey(storage.ErrInvalidKey, key)
	}

	return
}

// isDirTraversal will check if rootPlusPath is a dir traversal outside of root,
// assuming that both are cleaned and that rootPlusPath is path.Join(root, somePath).
func isDirTraversal(root, rootPlusPath string) bool {
	switch root {

	// Root is $PWD, check
	// for traversal out of
	case ".":
		return strings.HasPrefix(rootPlusPath, "../")

	// Root is *root*, ensure
	// it's not trying escape
	case "/":
		switch l := len(rootPlusPath); {
		case l == 3: // i.e. root=/ plusPath=/..
			return rootPlusPath[:3] == "/.."
		case l >= 4: // i.e. root=/ plusPath=/../[etc]
			return rootPlusPath[:4] == "/../"
		}
		return false
	}
	switch {

	// The path MUST be prefixed by storage root
	case !strings.HasPrefix(rootPlusPath, root):
		return true

	// In all other cases,
	// check not equal
	default:
		return len(root) == len(rootPlusPath)
	}
}
