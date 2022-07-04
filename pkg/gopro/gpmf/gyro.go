package gpmf // nolint: dupl

import (
	"time"
)

// GyroData represents Gyro data.
type GyroData []Gyro

func (d GyroData) offsets(start, end time.Duration) {
	offsets(start, end, d, func(i int, val time.Duration) {
		v := d[i]
		v.Offset = val
		d[i] = v
	})
}

// Gyro represents gyroscope metric for each axes.
type Gyro struct {
	X      float64
	Y      float64
	Z      float64
	Offset time.Duration
}

func parseGyro(e *Element) error {
	e.initMetadata()
	return floatType[GyroData](e, 3, func(vals []float64) Gyro {
		return Gyro{
			Z: vals[0],
			X: vals[1],
			Y: vals[2],
		}
	})
}
