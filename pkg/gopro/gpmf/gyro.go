package gpmf

// Gyro represents gyroscope metric for each axes.
type Gyro struct {
	X float64
	Y float64
	Z float64
}

func parseGyro(e *Element) error {
	e.metadata()
	return floatType(e, 3, func(vals []float64) Gyro {
		return Gyro{
			Z: vals[0],
			X: vals[1],
			Y: vals[2],
		}
	})
}
