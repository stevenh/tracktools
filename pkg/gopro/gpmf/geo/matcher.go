package geo

import (
	"math"

	"github.com/tidwall/geodesic"
)

// MatcherOpt represents a Matcher Option.
type MatcherOpt func(*Matcher)

// Radius sets the radius in meters of the ellipsoid.
func Radius(meters float64) MatcherOpt {
	return func(m *Matcher) {
		m.radius = meters
	}
}

// Tolerance sets the tolerance in meters for matches.
func Tolerance(meters float64) MatcherOpt {
	return func(m *Matcher) {
		m.tolerance = meters
	}
}

// Matcher represents a haversine matcher.
type Matcher struct {
	radius, tolerance float64
}

// NewMatcher returns a new haversine Matcher initialised from e.
func NewMatcher(options ...MatcherOpt) *Matcher {
	m := &Matcher{
		radius:    geodesic.WGS84.Radius(),
		tolerance: 0.1,
	}

	for _, f := range options {
		f(m)
	}

	// calculate the haversine tolerance.
	m.tolerance = hav(m.tolerance / m.radius)

	return m
}

// OnLine returns true if (lat0, lon0) is on the line between
// (lat1, lon2) and (lat2, lon2) within the Matchers tolerance,
// false otherwise.
// Latitudes and longitudes are in degrees.
func (m *Matcher) OnLine(lat0, lon0, lat1, lon1, lat2, lon2 float64) bool {
	return m.onLineRadians(
		lat0*radians, lon0*radians,
		lat1*radians, lon1*radians,
		lat2*radians, lon2*radians,
	)
}

// onLineRadians returns true if (lat0, lon0) is on the line between
// (lat1, lon2) and (lat2, lon2) within the Matchers tolerance,
// false otherwise.
func (m *Matcher) onLineRadians(lat0, lon0, lat1, lon1, lat2, lon2 float64) bool {
	dist01 := distance(lat0, lon0, lat1, lon1)
	if dist01 <= m.tolerance {
		return true
	}

	dist02 := distance(lat0, lon0, lat2, lon2)
	if dist02 <= m.tolerance {
		return true
	}

	bearing := sinDeltaBearing(lat1, lon1, lat2, lon2, lat0, lon0)
	sinDist1 := sinHav(dist01)
	track := havSin(sinDist1 * bearing)
	if track > m.tolerance {
		return false
	}

	dist12 := distance(lat1, lon1, lat2, lon2)
	term := dist12 + track*(1-2*dist12)
	if dist01 > term || dist02 > term {
		return false
	}

	if dist12 < 0.74 {
		return true
	}

	cosTrack := 1 - 2*track
	return sinSum(
		(dist01-track)/cosTrack,
		(dist02-track)/cosTrack,
	) > 0
}

// hav returns the haversine of x.
func hav(x float64) float64 {
	s := math.Sin(x * 0.5)
	return s * s
}

// distance returns the haversine distance between to points.
func distance(lat0, lon0, lat1, lon1 float64) float64 {
	return hav(lat0-lat1) + hav(lon0-lon1)*math.Cos(lat0)*math.Cos(lat1)
}

// sinSum returns sin(invHav(x) + invHav(y)).
func sinSum(x, y float64) float64 {
	a := math.Sqrt(x * (1 - x))
	b := math.Sqrt(y * (1 - y))
	return 2 * (a + b - 2*(a*y+b*x))
}

// sinHav given h == hav(x) returns sin(abs(x)).
func sinHav(h float64) float64 {
	return 2 * math.Sqrt(math.Abs(h))
}

// havSin returns hav(asin(x)).
func havSin(x float64) float64 {
	x2 := x * x
	return x2 / (1 + math.Sqrt(1-x2)) * .5
}

// sinDeltaBearing returns sin of difference between
// bearing from (lat1, lon1) to (lat0, lon0)
// bearing from (lat1, lon1) to (lat2, lon2).
func sinDeltaBearing(lat1, lon1, lat2, lon2, lat0, lon0 float64) float64 {
	sinLat1 := math.Sin(lat1)
	cosLat2 := math.Cos(lat2)
	cosLat3 := math.Cos(lat0)
	lat01 := lat0 - lat1
	lon01 := lon0 - lon1
	lat21 := lat2 - lat1
	lon21 := lon2 - lon1
	a := math.Sin(lon01) * cosLat3
	c := math.Sin(lon21) * cosLat2
	b := math.Sin(lat01) + 2*sinLat1*cosLat3*hav(lon01)
	d := math.Sin(lat21) + 2*sinLat1*cosLat2*hav(lon21)
	denom := (a*a + b*b) * (c*c + d*d)
	if denom <= 0 {
		return 1
	}

	return (a*d - b*c) / math.Sqrt(denom)
}
