package gpmf

// Magnetometer represents camera pointing direction.
type Magnetometer struct {
	X float64
	Y float64
	Z float64
}

func parseMagnetometer(e *Element) error {
	return floatType(e, 3, func(vals []float64) Magnetometer {
		return Magnetometer{
			Z: vals[0],
			X: vals[1],
			Y: vals[2],
		}
	})
}
