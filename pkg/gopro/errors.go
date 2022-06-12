package gopro

import (
	"errors"
	"fmt"
)

var (
	// errNoChapters is returned no chapters exist.
	errNoChapters = errors.New("no chapters")

	// ErrNoMatch is returned by a matcher if no match is found.
	ErrNoMatch = errors.New("no match")

	// ErrNoFiles is returned if no files were found.
	ErrNoFiles = errors.New("no files")
)

type configError string

func (e configError) Error() string {
	return fmt.Sprintf("missing config %q", string(e))
}

// matchError is returned by Matcher if an unexpected match length is found.
type matchError int

// Error implements error.
func (e matchError) Error() string {
	return fmt.Sprintf("unexpected match len: %d", e)
}

// chapterError is returned if an unexpected chapter order is detected.
type chapterError struct {
	chapter string
	idx     int
}

// Error implements error.
func (e chapterError) Error() string {
	return fmt.Sprintf("unexpected chapter %q at index %d", e.chapter, e.idx)
}
