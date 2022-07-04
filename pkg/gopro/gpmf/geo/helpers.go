package geo

import (
	"math"
)

// remquo returns the floating-point remainder of numer/denom and quotient.
// This replicates the C function of the same name.
func remquo(numer, denom float64) (float64, int) {
	return math.Remainder(numer, denom), int(math.Round(numer / denom))
}

// sincosd returns the sine and cosine function with the argument in degrees
// while doing its best to minimize round-off errors.
func sincosd(x float64) (sin, cos float64) {
	// In order to minimize round-off errors, this function exactly reduces
	// the argument to the range [-45, 45] before converting it to radians.
	r, q := remquo(x, quarterDegrees)
	s, c := math.Sincos(r * radians)
	switch uint(q) & 3 {
	case 0:
		sin = s
		cos = c
	case 1:
		sin = c
		cos = -s
	case 2:
		sin = -s
		cos = -c
	default:
		sin = -c
		cos = s
	}

	if sin == 0 {
		sin = math.Copysign(sin, x)
	}

	return sin, cos
}

// atan2d returns atan2 with the result in degrees while doing its best to
// minimize round-off errors.
func atan2d(y float64, x float64) float64 {
	// In order to minimize round-off errors, this function rearranges the
	// arguments so that result of atan2 is in the range [-pi/4, pi/4] before
	// converting it to degrees and mapping the result to the correct
	// quadrant.
	var q int
	if math.Abs(y) > math.Abs(x) {
		x, y = y, x
		q = 2
	}
	if x < 0 {
		x = -x
		q++
	}

	// Here x >= 0 and x >= abs(y), so angle is in [-pi/4, pi/4].
	ang := math.Atan2(y, x) / radians
	switch q {
	case 1:
		return math.Copysign(halfDegrees, y) - ang
	case 2:
		return quarterDegrees - ang
	case 3:
		return -quarterDegrees + ang
	}
	return ang
}
