package gpmf

import (
	"time"
)

// AccelData represents acceleration data.
type AccelData []Accel

// Accel represents acceleration for each axes.
type Accel struct {
	X      float64
	Y      float64
	Z      float64
	Offset time.Duration
}

// offsets implements offseter.
func (d AccelData) offsets(start, end time.Duration) {
	offsets(start, end, d, func(i int, val time.Duration) {
		v := d[i]
		v.Offset = val
		d[i] = v
	})
}

func parseAccel(e *Element) error {
	e.metadata()
	return floatType[AccelData](e, 3, func(vals []float64) Accel {
		return Accel{
			Z: vals[0],
			X: vals[1],
			Y: vals[2],
		}
	})
}
