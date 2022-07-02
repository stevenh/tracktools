package geo

const (
	// defaultRadius is the default radius used by Processor
	// which represents the radius of the earth.
	defaultRadius = 6378137
)

// Option represents a Processor Option.
type Option func(*Processor)

// Radius sets the radius of the ellipsoid.
// Default: 6378137
func Radius(val float64) Option {
	return func(p *Processor) {
		p.radius = val
	}
}

// Tolerance sets the tolerance for matches.
// Default: 0.1.
func Tolerance(val float64) Option {
	return func(p *Processor) {
		p.tolerance = val
	}
}

// Processor represents a geographic processor.
type Processor struct {
	// radius is the radius used for spherical calculations.
	radius float64

	// tolerance is the tolerance for matches.
	tolerance float64

	// havTolerance is the haversine tolerance.
	havTolerance float64

	// distFunc is the function used to calculate distances between two points.
	distFunc func(lat0, lon0, lat1, lon1, radius float64) float64
}

// NewProcessor returns a new geographic Processor.
func NewProcessor(options ...Option) *Processor {
	p := &Processor{
		radius:    defaultRadius,
		tolerance: 0.1,
		distFunc:  distanceHaversin,
	}

	for _, f := range options {
		f(p)
	}

	// calculate the haversine tolerance.
	p.havTolerance = hav(p.tolerance / p.radius)

	return p
}

// OnLine returns true if (lat0, lon0) is on the line between
// (lat1, lon2) and (lat2, lon2) within the Processors tolerance,
// false otherwise.
// Latitudes and longitudes are in degrees.
func (p *Processor) OnLine(lat0, lon0, lat1, lon1, lat2, lon2 float64) bool {
	return p.onLineRadians(
		lat0*radians, lon0*radians,
		lat1*radians, lon1*radians,
		lat2*radians, lon2*radians,
	)
}

// onLineRadians returns true if (lat0, lon0) is on the line between
// (lat1, lon2) and (lat2, lon2) within the Processors tolerance,
// false otherwise.
func (p *Processor) onLineRadians(lat0, lon0, lat1, lon1, lat2, lon2 float64) bool {
	dist01 := distanceHav(lat0, lon0, lat1, lon1)
	if dist01 <= p.havTolerance {
		return true
	}

	dist02 := distanceHav(lat0, lon0, lat2, lon2)
	if dist02 <= p.havTolerance {
		return true
	}

	bearing := sinDeltaBearing(lat1, lon1, lat2, lon2, lat0, lon0)
	sinDist1 := sinHav(dist01)
	track := havSin(sinDist1 * bearing)
	if track > p.havTolerance {
		return false
	}

	dist12 := distanceHav(lat1, lon1, lat2, lon2)
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

// DistanceToLine returns the distance in on a sphere between
// a point and the line between start and end specified in degrees.
func (p *Processor) DistanceToLine(pointLat, pointLon, startLat, startLon, endLat, endLon float64) float64 {
	pointLat *= radians
	pointLon *= radians
	startLat *= radians
	startLon *= radians
	endLat *= radians
	endLon *= radians

	if startLat == endLat && startLon == endLon {
		// Single point not a line.
		return distanceHaversin(pointLat, pointLon, endLat, endLon, p.radius)
	}

	diffLat := endLat - startLat
	diffLon := endLon - startLon
	u := ((pointLat-startLat)*diffLat + (pointLon-startLon)*diffLon) /
		(diffLat*diffLat + diffLon*diffLon)
	switch {
	case u <= 0:
		// Start is closer.
		return distanceHaversin(pointLat, pointLon, startLat, startLon, p.radius)
	case u >= 1:
		// End is closer.
		return distanceHaversin(pointLat, pointLon, endLat, endLon, p.radius)
	default:
		aLat := pointLat - startLat
		aLon := pointLon - startLon
		bLat := u * (endLat - startLat)
		bLon := u * (endLon - startLon)
		return distanceHaversin(aLat, aLon, bLat, bLon, p.radius)
	}
}

// Distance returns the distance on a sphere between
// two points expressed as Latitude, Longitude in degrees.
func (p *Processor) Distance(lat1, lon1, lat2, lon2 float64) float64 {
	return p.distFunc(lat1, lon1, lat2, lon2, p.radius)
}
