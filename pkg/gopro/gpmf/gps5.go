package gpmf

import (
	"fmt"
	"time"

	"github.com/rs/zerolog"
)

// GPSData represents GPS data.
type GPSData []GPS

// offsets implements offseter.
func (d GPSData) offsets(start, end time.Duration) {
	offsets(start, end, d, func(i int, val time.Duration) {
		v := d[i]
		v.Offset = val
		d[i] = v
	})
}

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

// MarshalZerologObject implements zerolog.LogObjectMarshaler.
func (g GPS) MarshalZerologObject(e *zerolog.Event) {
	e.Float64("latitude", g.Latitude).
		Float64("longitude", g.Longitude).
		Float64("altitude", g.Altitude).
		Float64("speed", g.Speed).
		Float64("speed3d", g.Speed3D).
		Str("offset", g.Offset.String())
}

func parseGPS(e *Element) error {
	e.initMetadata()
	return floatType[GPSData](e, 5, func(vals []float64) GPS {
		return GPS{
			Latitude:  vals[0],
			Longitude: vals[1],
			Altitude:  vals[2],
			Speed:     vals[3],
			Speed3D:   vals[4],
		}
	})
}
