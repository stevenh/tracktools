package trackaddict

import (
	"fmt"
	"time"
)

// Record represents an individual row of data for a Lap.
type Record struct {
	// Now is the offset since the beginning of the Lap.
	Now time.Duration

	// Time of the data.
	Time time.Time

	// Lap is the lap number.
	Lap int

	// Predicted lap time.
	Predicted time.Duration

	// Offset is the offset against the best Lap.
	Offset time.Duration

	// GPS represents GPS details
	GPS GPS

	// Speed.
	Speed float64

	// Acceleration represents the acceleration in all 3 axis.
	Accel *Acceleration

	// Brake indicates if the break is applied.
	// This may be a calculated value.
	Brake bool

	// BarometricPressure.
	BarometricPressure float64

	// PressureAltitute.
	PressureAltitute float64

	// OBD is the OBD data.
	OBD *OBD
}

// InitOBD initialises if needed and returns the OBD field.
func (r *Record) InitOBD() *OBD {
	if r.OBD == nil {
		r.OBD = &OBD{}
	}

	return r.OBD
}

// InitAccel initialises if needed and returns the Accel field.
func (r *Record) InitAccel() *Acceleration {
	if r.Accel == nil {
		r.Accel = &Acceleration{}
	}

	return r.Accel
}

// parseRecordNow parses and sets Now from value.
func parseRecordNow(r *Record, value string) (err error) {
	r.Now, err = parseDuration("record now", value)
	return err
}

// parseRecordTime parses and sets Now from value.
// It expects value in the form: <unix time in seconds>.<milliseconds>.
func parseRecordTime(r *Record, value string) error {
	var s, ms int64
	n, err := fmt.Sscanf(value, "%d.%d", &s, &ms)
	if err != nil {
		return fmt.Errorf("record parse time %q: %w", value, err)
	} else if n != 2 {
		return fmt.Errorf("record parse time %q", value)
	}

	r.Time = time.Unix(s, ms*int64(time.Millisecond))

	return nil
}

// parseRecordLap parses and sets Lap from value.
func parseRecordLap(r *Record, value string) (err error) {
	r.Lap, err = parseInt("record lap", value)
	return err
}

// parseRecordPredicted parses and sets Predicted from value.
func parseRecordPredicted(r *Record, value string) (err error) {
	r.Predicted, err = parseDuration("record predicted", value)
	return err
}

// parseRecordOffset parses and sets Offset from value.
func parseRecordOffset(r *Record, value string) (err error) {
	r.Offset, err = parseDuration("record offset", value)
	return err
}

// parseRecordBarometricPressure parses and sets BarometricPressure from value.
func parseRecordBarometricPressure(r *Record, value string, converters ...converter) (err error) {
	r.BarometricPressure, err = parseFloat64("record barometric pressure", value, converters...)
	return err
}

// parseRecordPressureAltitude parses and sets PressureAltitude from value.
func parseRecordPressureAltitude(r *Record, value string, converters ...converter) (err error) {
	r.PressureAltitute, err = parseFloat64("record pressure altitude", value, converters...)
	return err
}

// parseRecordSpeed parses and sets Speed from value.
func parseRecordSpeed(r *Record, value string, converters ...converter) (err error) {
	r.Speed, err = parseFloat64("record speed", value, converters...)
	return err
}

// parseRecordBrake parses and sets Brake from value.
func parseRecordBrake(r *Record, value string) (err error) {
	r.Brake, err = parseBool("record brake", value)
	return err
}
