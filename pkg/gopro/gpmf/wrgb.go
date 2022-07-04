package gpmf // nolint: dupl

import (
	"time"
)

// WhiteBalanceRGBData represents WhiteBalanceRGB data.
type WhiteBalanceRGBData []WhiteBalanceRGB

// offsets implements offseter.
func (d WhiteBalanceRGBData) offsets(start, end time.Duration) {
	offsets(start, end, d, func(i int, val time.Duration) {
		v := d[i]
		v.Offset = val
		d[i] = v
	})
}

// WhiteBalanceRGB represents a white balance with RGB channels.
type WhiteBalanceRGB struct {
	Red    float64
	Green  float64
	Blue   float64
	Offset time.Duration
}

func parseWhiteBalanceRGB(e *Element) error {
	e.initMetadata()
	return floatType[WhiteBalanceRGBData](e, 3, func(vals []float64) WhiteBalanceRGB {
		return WhiteBalanceRGB{
			Red:   vals[0],
			Green: vals[1],
			Blue:  vals[2],
		}
	})
}
