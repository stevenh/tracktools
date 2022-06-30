package geo

import (
	"math"

	"github.com/tidwall/geodesic"
)

const (
	// quarterDegrees the number of degrees in a quarter of a turn.
	quarterDegrees = 90.0

	// halfDegrees the number of degrees in a half of a turn.
	halfDegrees = 2 * quarterDegrees

	// numIterations is the number of iterations used for solves.
	numIterations = 20
)

var (
	// epsilon is the difference between 1 and the least value greater
	// than 1 that is representable.
	epsilon = math.Nextafter(1, 2) - 1
	eps     = 0.01 * math.Sqrt(epsilon)

	// radians is the number of radians in a degree.
	radians = math.Pi / halfDegrees
)

// Gnomonic projection centred at an arbitrary position C on the ellipsoid.
type Gnomonic struct {
	earth  *geodesic.Ellipsoid
	radius float64
}

// NewGnomonic returns a new Gnomonic initialised from e.
func NewGnomonic(e *geodesic.Ellipsoid) *Gnomonic {
	return &Gnomonic{
		earth:  e,
		radius: e.Radius(),
	}
}

// Forward projection, from geographic to gnomonic.
//
// lat0	latitude of center point of projection in the range [−90°, 90°].
// lon0	longitude of center point of projection in degrees.
// lat	latitude of point in the range [−90°, 90°].
// lon	longitude of point in degrees.
//
// It returns:
// x	easting of point in meters.
// y	northing of point in meters.
// azi	azimuth of geodesic at point in degrees.
// rk	reciprocal of azimuthal scale at point.
//
// The scale of the projection is 1/rk² in the "radial" direction, azi clockwise
// from true north, and is 1/rk in the direction perpendicular to this. If the
// point lies "over the horizon", i.e., if rk ≤ 0, then NaNs are returned for
// x and y (the correct values are returned for azi and rk).
//
// A call to Forward followed by a call to Reverse will return the original (lat, lon)
// (to within roundoff) provided the point in not over the horizon.
func (g *Gnomonic) Forward(lat0, lon0, lat, lon float64) (x, y, azi, rk float64) {
	var m12 float64
	g.earth.GenInverse(
		lat0, lon0, lat, lon,
		nil, &azi, nil, &m12, &rk, nil, nil,
	)

	if rk <= 0 {
		nan := math.NaN()
		return nan, nan, nan, nan
	}

	rho := m12 / rk
	x, y = sincosd(azi)
	x *= rho
	y *= rho

	return x, y, azi, rk
}

// Reverse projection, from gnomonic to geographic.
// lat0	latitude of center point of projection in the range [−90°, 90°] (degrees).
// lon0	longitude of center point of projection (degrees).
// x	easting of point (meters).
// y	northing of point (meters).
//
// Returns:
// lat	latitude of point in the range [−90°, 90°] (degrees).
// lon	longitude of point in the range [−180°, 180°] (degrees).
// azi	azimuth of geodesic at point (degrees).
// rk	reciprocal of azimuthal scale at point.
//
// The scale of the projection is 1/rk² in the "radial" direction, azi clockwise
// from true north, and is 1/rk in the direction perpendicular to this. Even though
// all inputs should return a valid lat and lon, it's possible that the procedure
// fails to converge for very large x or y; in this case NaNs are returned for all
// the output arguments. A call to Reverse followed by a call to Forward will
// return the original (x, y) (to roundoff).
func (g *Gnomonic) Reverse(lat0, lon0, x, y float64) (lat, lon, azi, rk float64) {
	azi0 := atan2d(x, y)
	rho := math.Hypot(x, y)
	s := g.radius * math.Atan(rho/g.radius)
	little := rho <= g.radius

	if !little {
		rho = 1 / rho
	}

	line := g.earth.LineInit(lat0, lon0, azi0,
		geodesic.Latitude|geodesic.Longitude|
			geodesic.Azimuth|geodesic.DistanceIn|
			geodesic.ReducedLength|geodesic.GeodesicScale,
	)

	var trip int
	var s12, m12 float64
	for count := numIterations - 1; count > 0; count-- {
		line.GenPosition(geodesic.NoFlags,
			s, &lat, &lon, &azi, &s12, &m12, &rk, nil, nil,
		)

		if trip > 0 {
			break
		}

		var ds float64
		if little {
			ds = (m12 - rho*rk) * rk
		} else {
			ds = (rho*m12 - rk) * m12
		}
		s -= ds

		if math.Abs(ds) < eps*g.radius {
			trip++
		}
	}

	if trip == 0 {
		nan := math.NaN()
		return nan, nan, nan, nan
	}

	return lat, lon, azi, rk
}
