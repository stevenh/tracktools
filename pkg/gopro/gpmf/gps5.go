package gpmf

import (
	"fmt"
	"time"
)

// GPS represents GPS5 data.
type GPS struct {
	Latitude  float64
	Longitude float64
	Altitude  float64
	Speed     float64
	Speed3D   float64
	Offset    time.Duration
}

func (g GPS) String() string {
	return fmt.Sprintf("pos: %.7f,%.7f, alt: %.2f, speed: %.2f, speed3d: %.2f off: %s",
		g.Latitude,
		g.Longitude,
		g.Altitude,
		g.Speed,
		g.Speed3D,
		g.Offset,
	)
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
