package geo

import (
	"fmt"
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
// x	easting of point.
// y	northing of point.
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
// x	easting of point.
// y	northing of point.
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

// IntersectExt returns the intersection of lines A and B defined
// by their end points as latitudes in the range [−90°, 90°] and
// longitudes in degrees if they are extended.
//   Line A: point 1 (lat1a, lon1a) -> point 2 (lat2a, lon2a)
//   Line B: point 1 (lat1b, lon1b) -> point 2 (lat2b, lon2b)
//
// It returns the latitude, longitide of intersection point and
// the azimuths:
// azia1: azimuth from line A point 1 to intersection point.
// azia2: azimuth from line A point 2 to intersection point.
// azib1: azimuth from line B point 1 to intersection point.
// azib1: azimuth from line B point 2 to intersection point.
func (g *Gnomonic) IntersectExt(
	lat1a, lon1a,
	lat2a, lon2a,
	lat1b, lon1b,
	lat2b, lon2b float64,
) (lat, lon, azia1, azia2, azib1, azib2 float64) {
	lat = (lat1a + lat2a + lat1b + lat2b) / 4
	// Possibly need to deal with longitudes wrapping around.
	lon = (lon1a + lon2a + lon1b + lon2b) / 4
	for i := 0; i < numIterations; i++ {
		xa1, ya1, _, _ := g.Forward(lat, lon, lat1a, lon1a)
		xa2, ya2, _, _ := g.Forward(lat, lon, lat2a, lon2a)
		xb1, yb1, _, _ := g.Forward(lat, lon, lat1b, lon1b)
		xb2, yb2, _, _ := g.Forward(lat, lon, lat2b, lon2b)
		// See Hartley and Zisserman, Multiple View Geometry, Sec. 2.2.1
		va1 := newVector2(xa1, ya1)
		va2 := newVector2(xa2, ya2)
		vb1 := newVector2(xb1, yb1)
		vb2 := newVector2(xb2, yb2)
		// la is homogeneous representation of line A1,A2.
		// lb is homogeneous representation of line B1,B2.
		la := va1.cross(va2)
		lb := vb1.cross(vb2)
		// p0 is homogeneous representation of intersection of la and lb.
		p0 := la.cross(lb)
		p0.norm()

		lat0, lon0, _, _ := g.Reverse(lat, lon, p0.x, p0.y)
		if lat0 == lat && lon0 == lon {
			// No change so stop early.
			break
		}
		lat, lon = lat0, lon0
	}

	g.earth.Inverse(lat1a, lon1a, lat, lon, nil, nil, &azia2)
	g.earth.Inverse(lat, lon, lat2a, lon2a, nil, &azia1, nil)

	g.earth.Inverse(lat1b, lon1b, lat, lon, nil, nil, &azib2)
	g.earth.Inverse(lat, lon, lat2b, lon2b, nil, &azib1, nil)

	return lat, lon, azia1, azia2, azib1, azib2
}

// Intersect returns the intersection latitude and longitude
// of lines A and B defined by their end points as latitudes in
// the range [−90°, 90°] and longitudes in degrees.
//   Line A: point 1 (lat1a, lon1a) -> point 2 (lat2a, lon2a)
//   Line B: point 1 (lat1b, lon1b) -> point 2 (lat2b, lon2b)
func (g *Gnomonic) Intersect(
	lat1a, lon1a,
	lat2a, lon2a,
	lat1b, lon1b,
	lat2b, lon2b float64,
) (lat, lon float64, err error) {
	lat, lon, azia1, azia2, azib1, azib2 := g.IntersectExt(
		lat1a, lon1a,
		lat2a, lon2a,
		lat1b, lon1b,
		lat2b, lon2b,
	)

	if math.Signbit(azia1) == math.Signbit(azia2) ||
		math.Signbit(azib1) == math.Signbit(azib2) {
		return 0, 0, fmt.Errorf("line a %f, %f -> %f, %f doesn't intersect with line b %f, %f -> %f, %f",
			lat1a, lon1a,
			lat2a, lon2a,
			lat1b, lon1b,
			lat2b, lon2b,
		)
	}

	return lat, lon, nil
}
