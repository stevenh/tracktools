package trackaddict

import (
	"fmt"
	"time"
)

// Lap represents a track lap.
// The Records may be partial.
type Lap struct {
	// Duration of the Lap.
	Duration time.Duration

	// Number is the number of this Lap.
	Number int

	// Records are the records for this Lap.
	Records []Record
}

// parseLapNumber parses and sets the lap number.
func parseLapNumber(l *Lap, value string) (err error) {
	l.Number, err = parseInt("lap number", value)
	return err
}

// parseLapDuration parses and sets lap duration.
func parseLapDuration(l *Lap, value string) (err error) {
	// Format: 00:02:03.202
	var h, m, s, ms int
	n, err := fmt.Sscanf(value, "%d:%d:%d.%d", &h, &m, &s, &ms)
	if err != nil {
		return fmt.Errorf("parse lap duration %q: %w", value, err)
	} else if n != 4 {
		return fmt.Errorf("parse lap duration %q", value)
	}

	value = fmt.Sprintf("%dh%dm%ds%dms", h, m, s, ms)
	if l.Duration, err = time.ParseDuration(value); err != nil {
		return fmt.Errorf("parse lap duration %q", value)
	}

	return nil
}
