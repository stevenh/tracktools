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

// offsetWalker traverses Elements and sets their time offset.
type offsetWalker struct {
	start, end time.Duration
}

// newOffsetWalker creates a new offsetWalker.
func newOffsetWalker(start, end uint64, units time.Duration) *offsetWalker {
	return &offsetWalker{
		start: time.Duration(start) * units,
		end:   time.Duration(end) * units,
	}
}

// walk is WalkFunc which sets offsets.
func (o *offsetWalker) walk(e *Element) error {
	if v, ok := e.Data.(offseter); ok {
		v.offsets(o.start, o.end)
	}

	return nil
}
