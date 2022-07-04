package gpmf

import (
	"fmt"
	"time"
)

// floatType returns a slice of T created using f.
func floatType[T ~[]E, E any](e *Element, size int, f func([]float64) E) error {
	vals, err := floatSlice(e.Data)
	if err != nil {
		return err
	}

	var elem E
	if len(vals)%size != 0 {
		return fmt.Errorf("parse %T: invalid number of elements %d (not a multiple of %d)", elem, len(vals), size)
	}

	d := make(T, len(vals)/size)
	for i := range d {
		d[i] = f(vals[0:size])
		vals = vals[size:]
	}

	e.Data = d

	return nil
}

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
