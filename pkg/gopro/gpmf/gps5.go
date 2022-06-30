package gpmf

// GPS represents GPS5 data.
type GPS struct {
	Latitude  float64
	Longitude float64
	Altitude  float64
	Speed     float64
	Speed3D   float64
}

func parseGPS(e *Element) error {
	e.metadata()
	return floatType(e, 5, func(vals []float64) GPS {
		return GPS{
			Latitude:  vals[0],
			Longitude: vals[1],
			Altitude:  vals[2],
			Speed:     vals[3],
			Speed3D:   vals[4],
		}
	})
}
