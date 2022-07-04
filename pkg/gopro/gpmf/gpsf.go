package gpmf

import (
	"fmt"
)

// GPSFix represents the type of a GPS fix.
type GPSFix uint32

const (
	// GPSNoLock represents a GPS fix which has no satellite lock.
	GPSNoLock GPSFix = iota
	_
	// GPS2DLock represents a GPS fix which only has a 2D lock.
	GPS2DLock

	// GPS3DLock represents a GPS fix which only has a 3D lock.
	GPS3DLock
)

func parseGPSFix(e *Element) error {
	v, ok := e.Data.(uint32)
	if !ok {
		return fmt.Errorf("gps fix: unexpected data type %T (expected uint32)", e.Data)
	}

	f := GPSFix(v)
	e.Data = f
	e.parent.Metadata[e.friendlyName()] = e.Data

	var s string
	switch f {
	case GPSNoLock:
		s = "No lock"
	case GPS2DLock:
		s = "2D lock"
	case GPS3DLock:
		s = "3D lock"
	default:
		s = fmt.Sprintf("unknown lock: %d", v)
	}

	e.parent.Metadata["gps_fix_description"] = s

	return nil
}
