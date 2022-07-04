package geo

import (
	"math"
)

// distanceHav returns the haversine distance between two points
// specified by their Latitude and Longitude in radians and
// the radius of the globe calculated using [haversine formula].
//
// [haversine formula]: https://en.wikipedia.org/wiki/Haversine_formula
func distanceHav(lat1, lon1, lat2, lon2 float64) float64 {
	return hav(lat1-lat2) + hav(lon1-lon2)*math.Cos(lat1)*math.Cos(lat2)
}

// sinSum returns sin(invHav(x) + invHav(y)).
func sinSum(x, y float64) float64 {
	a := math.Sqrt(x * (1 - x))
	b := math.Sqrt(y * (1 - y))
	return 2 * (a + b - 2*(a*y+b*x))
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

// distanceHaversin returns the distance between two points
// specified by their Latitude and Longitude in radians and
// the radius of the globe calculated using [haversine formula].
//
// [haversine formula]: https://en.wikipedia.org/wiki/Haversine_formula
func distanceHaversin(lat1, lon1, lat2, lon2, radius float64) float64 {
	return invHav(distanceHav(lat1, lon1, lat2, lon2)) * radius
}

// hav returns the haversin of x.
func hav(x float64) float64 {
	s := math.Sin(x * 0.5)
	return s * s
}

// invHav returns the inverse haversin.
func invHav(x float64) float64 {
	return 2 * math.Asin(math.Sqrt(x))
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
