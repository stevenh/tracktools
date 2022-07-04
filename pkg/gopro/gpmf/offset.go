package gpmf

import (
	"time"
)

// setOffset is a function which can set the offset of a dataset.
type setOffset func(idx int, offset time.Duration)

// offsets sets the offsets for s using fn.
func offsets[S ~[]E, E any](start, end time.Duration, s S, fn setOffset) {
	offset := start
	dur := end - start
	inc := dur / time.Duration(len(s))
	for i := range s {
		fn(i, offset)
		offset += inc
	}
}
