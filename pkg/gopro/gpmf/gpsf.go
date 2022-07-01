package gpmf

import (
	"fmt"
)

type gpsFix uint32

const (
	gpsFixNoLock gpsFix = iota
	_
	gpsFix2DLock
	gpsFix3DLock
)

func parseGPSFix(e *Element) error {
	v, ok := e.Data.(uint32)
	if !ok {
		return fmt.Errorf("gps fix: unexpected data type %T (expected uint32)", e.Data)
	}

	switch gpsFix(v) {
	case gpsFixNoLock:
		e.Data = "no lock"
	case gpsFix2DLock:
		e.Data = "2D lock"
	case gpsFix3DLock:
		e.Data = "3D lock"
	default:
		e.Data = fmt.Sprintf("gps fix: unknown lock: %d", v)
	}

	return parseMetadata(e)
}
