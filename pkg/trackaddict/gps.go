package trackaddict

import (
	"time"
)

// GPS represents GPS data.
type GPS struct {
	Update    bool
	Delay     time.Duration
	Latitude  float64
	Longitude float64
	Altitude  float64
	Accuracy  float64
	Heading   float64
}

func (g *GPS) parseLatitude(value string) (err error) {
	g.Latitude, err = parseFloat64("latitude", value)
	return err
}

func (g *GPS) parseLongitude(value string) (err error) {
	g.Longitude, err = parseFloat64("gps longitude", value)
	return err
}

func (g *GPS) parseHeading(value string) (err error) {
	g.Heading, err = parseFloat64("gps heading", value)
	return err
}

func (g *GPS) parseCoordinate(lat, long, head string) error {
	if err := g.parseLatitude(lat); err != nil {
		return err
	}
	if err := g.parseLongitude(lat); err != nil {
		return err
	}
	return g.parseHeading(lat)
}

func parseGPSUpdate(r *Record, value string) (err error) {
	r.GPS.Update, err = parseBool("gps update", value)
	return err
}

func parseGPSDelay(r *Record, value string) (err error) {
	r.GPS.Delay, err = parseDuration("gps delay", value)
	return err
}

func parseGPSLatitude(r *Record, value string) (err error) {
	return r.GPS.parseLatitude(value)
}

func parseGPSLongitude(r *Record, value string) error {
	return r.GPS.parseLongitude(value)
}

func parseGPSAltitude(r *Record, value string, funcs ...converter) (err error) {
	r.GPS.Altitude, err = parseFloat64("gps altitude", value, funcs...)
	return err
}

func parseGPSHeading(r *Record, value string) error {
	return r.GPS.parseHeading(value)
}

func parseGPSAccuracy(r *Record, value string, funcs ...converter) (err error) {
	r.GPS.Accuracy, err = parseFloat64("gps accuracy", value, funcs...)
	return err
}
