package gpmf

import (
	"time"
)

// MagnetometerData represents Magnetometer data.
type MagnetometerData []Magnetometer

// offsets implements offseter.
func (d MagnetometerData) offsets(start, end time.Duration) {
	offsets(start, end, d, func(i int, val time.Duration) {
		v := d[i]
		v.Offset = val
		d[i] = v
	})
}

// Magnetometer represents camera pointing direction.
type Magnetometer struct {
	X      float64
	Y      float64
	Z      float64
	Offset time.Duration
}

func parseMagnetometer(e *Element) error {
	return floatType[MagnetometerData](e, 3, func(vals []float64) Magnetometer {
		return Magnetometer{
			Z: vals[0],
			X: vals[1],
			Y: vals[2],
		}
	})
}
