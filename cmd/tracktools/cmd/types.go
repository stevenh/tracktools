package cmd

import (
	"fmt"
	"time"
)

const (
	dateFormat = "2006-01-02"
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
