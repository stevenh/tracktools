package gopro

import (
	"io/fs"
	"os"
	"path/filepath"
	"time"
)

// osFS implements fs interfaces using os functions.
type osFS struct{}

// Open implements procFS.
func (osFS) Open(name string) (fs.File, error) {
	return os.Open(filepath.FromSlash(name))
}

// ReadDir implements procFS.
func (osFS) ReadDir(name string) ([]fs.DirEntry, error) {
	return os.ReadDir(filepath.FromSlash(name))
}

// Stat implements procFS.
func (osFS) Stat(name string) (fs.FileInfo, error) {
	return os.Stat(filepath.FromSlash(name))
}

// CreateTemp implements procFS.
func (osFS) CreateTemp(dir, pattern string) (tempFile, error) { // nolint: ireturn
	return os.CreateTemp(filepath.FromSlash(dir), pattern)
}

// Chtimes implements procFS.
func (osFS) Chtimes(name string, atime time.Time, mtime time.Time) error {
	return os.Chtimes(filepath.FromSlash(name), atime, mtime)
}

// Remove implements procFS.
func (osFS) Remove(name string) error {
	return os.Remove(filepath.FromSlash(name))
}
