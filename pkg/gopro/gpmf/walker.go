package gpmf

import (
	"errors"
)

var (
	// ErrSkip is used as a return value from WalkFuncs to
	// indicate that the element in the call is to be skipped.
	// It is not returned as an error by any function.
	ErrSkip = errors.New("skip element")
)

// WalkFunc is the type of the function called by Walk to visit
// each Element.
type WalkFunc func(e *Element) error

// Walk walks the Element tree rooted at elems, calling fn for each
// Element in the tree.
func Walk(elems []*Element, fn WalkFunc) error {
	for _, e := range elems {
		err := fn(e)
		switch {
		case err == nil:
		case errors.Is(err, ErrSkip):
			continue
		default:
			return err
		}

		if err := Walk(e.Nested, fn); err != nil {
			return err
		}
	}

	return nil
}
