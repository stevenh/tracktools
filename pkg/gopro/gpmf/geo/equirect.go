package geo

import (
	"math"
)

// distanceEquirect returns the distance between two points
// specified by their Latitude and Longitude in radians and
// the radius of the globe calculated using [Equirectangular]
// projection.
//
// [Equirectangular]: https://en.wikipedia.org/wiki/Equirectangular_projection
func distanceEquirect(lat1, lon1, lat2, lon2, radius float64) float64 {
	lat1 *= radians
	lon1 *= radians
	lat2 *= radians
	lon2 *= radians

	x := (lon2 - lon1) * math.Cos((lat2+lat1)*0.5)
	y := lat2 - lat1
	return math.Sqrt((x*x)+(y*y)) * radius
}
