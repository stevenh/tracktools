package gopro

import (
	"io"
	"io/fs"
	"time"
)

// Log is implemented by types which can act as Processor logger.
type Log interface {
	Println(v ...interface{})
}

// procFS represents a filesystem implementation as needed by Processor.
type procFS interface {
	fs.StatFS

	// CreateTemp creates a temporary file in dir with pattern.
	CreateTemp(dir, pattern string) (tempFile, error)

	// Chtimes changes the access and modification times of the named file,
	// similar to the Unix utime() or utimes() functions.
	//
	// The underlying filesystem may truncate or round the values to a less
	// precise time unit. If there is an error, it will be of type *fs.PathError.
	Chtimes(name string, atime, mtime time.Time) error

	// Remove removes the named file or (empty) directory.
	// If there is an error, it will be of type *fs.PathError.
	Remove(name string) error
}

// tempFile represents a temporary file.
type tempFile interface {
	io.WriteCloser

	// Name returns the name of the file.
	Name() string
}
