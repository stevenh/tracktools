package trackaddict

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// converters returns the required converters eliminating
// nil functions.
func converters(funcs ...converter) []converter {
	var ret []converter
	for _, fn := range funcs {
		if fn != nil {
			ret = append(ret, fn)
		}
	}
	return ret
}

// parseDuration parses a duration in the form:
// <seconds>.<milliseconds>
// and returns the result.
func parseDuration(name, value string) (time.Duration, error) {
	value = strings.Replace(value, ".", "s", 1)
	d, err := time.ParseDuration(value + "ms")
	if err != nil {
		return 0, fmt.Errorf("parse %s %q: %w", name, value, err)
	}

	return d, nil
}

// parseFloat64p parses a float64 and returns a pointer to the result.
func parseFloat64p(name, value string, converters ...converter) (*float64, error) {
	v, err := parseFloat64(name, value, converters...)
	if err != nil {
		return nil, err
	}

	return &v, nil
}

// parseFloat64 parses a float64 and returns the result.
func parseFloat64(name, value string, converters ...converter) (float64, error) {
	f, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return 0, fmt.Errorf("parse %s %q: %w", name, value, err)
	}

	for _, fn := range converters {
		f = fn(f)
	}

	return f, nil
}

// parseBool parses a bool and returns the result.
func parseBool(name, value string) (bool, error) {
	b, err := strconv.ParseBool(value)
	if err != nil {
		return false, fmt.Errorf("parse %s %q: %w", name, value, err)
	}

	return b, nil
}

// parseInt parses a int and returns the result.
func parseInt(name, value string) (int, error) {
	i, err := strconv.Atoi(value)
	if err != nil {
		return 0, fmt.Errorf("parse %s %q: %w", name, value, err)
	}

	return i, nil
}
