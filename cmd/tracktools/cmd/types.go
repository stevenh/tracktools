package cmd

import (
	"fmt"
	"time"

	"github.com/tidwall/geodesic"
)

const (
	// dateFormat is the format used for date output.
	dateFormat = "2006-01-02"
)

var (
	// gd is the geodesic used for calculations.
	gd = geodesic.WGS84
)

type date time.Time

// MarshalText implements encoding.TextMarshaler.
func (d date) MarshalText() ([]byte, error) {
	return []byte(d.String()), nil
}

// UnmarshalText implements encoding.TextUnmarshaler.
func (d *date) UnmarshalText(text []byte) error {
	return d.Set(string(text))
}

// String implements pflags.Value.
func (d date) String() string {
	return time.Time(d).Format(dateFormat)
}

// Set implements pflags.Value.
func (d *date) Set(val string) error {
	fmt.Println("date:", val)
	if val == "" {
		// Ignore empty values.
		return nil
	}

	t, err := time.Parse(dateFormat, val)
	if err != nil {
		return err
	}

	*d = date(t)

	return nil
}

// Type implements pflags.Value.
func (d date) Type() string {
	return "date"
}

// Start represents a track start point.
type Start struct {
	Latitude  float64
	Longitude float64
	Bearing   float64
	Distance  float64

	lat1, lon1 float64
	lat2, lon2 float64
}

// calculates the start and end latitudes and longitudes of the start line.
func (s *Start) calculate() {
	gd.Direct(s.Latitude, s.Longitude, s.Bearing+90, s.Distance, &s.lat1, &s.lon1, nil)
	gd.Direct(s.Latitude, s.Longitude, s.Bearing-90, s.Distance, &s.lat2, &s.lon2, nil)
}
