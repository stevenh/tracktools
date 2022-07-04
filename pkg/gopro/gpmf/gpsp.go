package gpmf

import (
	"fmt"
)

// GPSDoP represents GPS Dilution of precision (DoP).
type GPSDoP float64

// Ideal returns true if positional measurements are suitable for
// applications demanding the highest possible precision at all
// times, false otherwise.
func (v GPSDoP) Ideal() bool {
	return v < 1
}

// Excellent returns true if positional measurements are considered
// accurate enough to meet all but the most sensitive applications,
// false otherwise.
func (v GPSDoP) Excellent() bool {
	return v >= 1 && v <= 2
}

// Good returns true if postional measurements meet the minimum
// appropriate for making accurate decisions. Positional measurements
// could be used to make reliable in-route navigation suggestions to
// the user, false otherwise.
func (v GPSDoP) Good() bool {
	return v >= 2 && v <= 5
}

// Moderate returns true if positional measurements could be used for
// calculations, but the fix quality could still be improved. A more
// open view of the sky is recommended, false otherwise.
func (v GPSDoP) Moderate() bool {
	return v >= 5 && v <= 10
}

// Fair returns true if positional measurements should be discarded
// or used only to indicate a very rough estimate of the current
// location, false otherwise.
func (v GPSDoP) Fair() bool {
	return v >= 10 && v <= 20
}

// Poor returns true if measurements are inaccurate by as much as 300
// meters with a 6-meter accurate device (50 DOP × 6 meters) and
// should be discarded, false otherwise.
func (v GPSDoP) Poor() bool {
	return v > 20
}

func parseGPSDoP(e *Element) error {
	v, ok := e.Data.(uint16)
	if !ok {
		return fmt.Errorf("gps dop: unexpected data type %T (expected uint16)", e.Data)
	}

	e.Data = GPSDoP(v) / 100

	return parseMetadata(e)
}
