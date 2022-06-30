package gpmf

// Accel represents acceleration for each axes.
type Accel struct {
	X float64
	Y float64
	Z float64
}

func parseAccel(e *Element) error {
	e.metadata()
	return floatType(e, 3, func(vals []float64) Accel {
		return Accel{
			Z: vals[0],
			X: vals[1],
			Y: vals[2],
		}
	})
}
