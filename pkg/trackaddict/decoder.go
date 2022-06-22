package trackaddict

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"regexp"
	"strings"
)

var (
	endpointRe = regexp.MustCompile(`([0-9\.\-]+), +([0-9\.\-]+) +@ +([0-9\.\-]+)`)
)

// Decoder reads and decodes CVS data from an input stream.
type Decoder struct {
	r       io.Reader
	parsers []func(r *Record, value string) error
	units   units
}

// Option represents a Decoder option.
type Option func(*Decoder) error

// NewDecoder returns a fully initialised Decoder which reads from r.
func NewDecoder(r io.Reader, options ...Option) (*Decoder, error) {
	d := &Decoder{r: r}
	for _, f := range options {
		if err := f(d); err != nil {
			return nil, err
		}
	}

	return d, nil
}

// setup sets up column parsers.
func (d *Decoder) columns(s *Session, cols []string) error { // nolint: gocyclo,cyclop
	for _, col := range cols {
		switch col {
		case "Time":
			d.parsers = append(d.parsers, parseRecordNow)
		case "UTC Time":
			d.parsers = append(d.parsers, parseRecordTime)
		case "Lap":
			d.parsers = append(d.parsers, parseRecordLap)
		case "Predicted Lap Time":
			d.parsers = append(d.parsers, parseRecordPredicted)
		case "Predicted vs Best Lap":
			d.parsers = append(d.parsers, parseRecordOffset)
		case "GPS_Update":
			d.parsers = append(d.parsers, parseGPSUpdate)
		case "GPS_Delay":
			d.parsers = append(d.parsers, parseGPSDelay)
		case "Latitude":
			d.parsers = append(d.parsers, parseGPSLatitude)
		case "Longitude":
			d.parsers = append(d.parsers, parseGPSLongitude)
		case "Altitude (m)":
			funcs := converters(d.units.Altitude)
			d.parsers = append(d.parsers, func(r *Record, value string) error {
				return parseGPSAltitude(r, value, funcs...)
			})
		case "Altitude (ft)":
			funcs := converters(feet2Meters, d.units.Altitude)
			d.parsers = append(d.parsers, func(r *Record, value string) error {
				return parseGPSAltitude(r, value, funcs...)
			})
		case "Speed (MPH)":
			funcs := converters(miles2Kilometers, d.units.Speed)
			d.parsers = append(d.parsers, func(r *Record, value string) error {
				return parseRecordSpeed(r, value, funcs...)
			})
		case "Speed (Km/h)":
			funcs := converters(d.units.Speed)
			d.parsers = append(d.parsers, func(r *Record, value string) error {
				return parseRecordSpeed(r, value, funcs...)
			})
		case "Heading":
			d.parsers = append(d.parsers, parseGPSHeading)
		case "Accuracy (m)":
			funcs := converters(d.units.Accuracy)
			d.parsers = append(d.parsers, func(r *Record, value string) error {
				return parseGPSAccuracy(r, value, funcs...)
			})
		case "Accuracy (ft)":
			funcs := converters(d.units.Accuracy)
			d.parsers = append(d.parsers, func(r *Record, value string) error {
				return parseGPSAccuracy(r, value, funcs...)
			})
		case "Accel X":
			d.parsers = append(d.parsers, parseAccelX)
		case "Accel Y":
			d.parsers = append(d.parsers, parseAccelY)
		case "Accel Z":
			d.parsers = append(d.parsers, parseAccelZ)
		case "Brake (calculated)":
			d.parsers = append(d.parsers, parseRecordBrake)
		case "Barometric Pressure (PSI)":
			funcs := converters(psi2Kpa, d.units.Pressure)
			d.parsers = append(d.parsers, func(r *Record, value string) error {
				return parseRecordBarometricPressure(r, value, funcs...)
			})
		case "Barometric Pressure (kPa)":
			funcs := converters(d.units.Pressure)
			d.parsers = append(d.parsers, func(r *Record, value string) error {
				return parseRecordBarometricPressure(r, value, funcs...)
			})
		case "Pressure Altitude (ft)":
			funcs := converters(feet2Meters, d.units.Altitude)
			d.parsers = append(d.parsers, func(r *Record, value string) error {
				return parseRecordPressureAltitude(r, value, funcs...)
			})
		case "Pressure Altitude (m)":
			funcs := converters(d.units.Altitude)
			d.parsers = append(d.parsers, func(r *Record, value string) error {
				return parseRecordPressureAltitude(r, value, funcs...)
			})
		case "OBD_Update":
			d.parsers = append(d.parsers, func(r *Record, value string) error {
				return parseOBDUpdate(r.InitOBD(), value)
			})
		case "Engine Speed (RPM) *OBD":
			d.parsers = append(d.parsers, func(r *Record, value string) error {
				return parseOBDEngineSpeed(r.InitOBD(), value)
			})
		case "Vehicle Speed (mph) *OBD":
			funcs := converters(miles2Kilometers, d.units.Speed)
			d.parsers = append(d.parsers, func(r *Record, value string) error {
				return parseOBDSpeed(r.InitOBD(), value, funcs...)
			})
		case "Vehicle Speed (km/h) *OBD":
			funcs := converters(d.units.Speed)
			d.parsers = append(d.parsers, func(r *Record, value string) error {
				return parseOBDSpeed(r.InitOBD(), value, funcs...)
			})
		case "Throttle Position (%) *OBD":
			d.parsers = append(d.parsers, func(r *Record, value string) error {
				return parseOBDThrottle(r.InitOBD(), value)
			})
		case "Engine Coolant Temp (F) *OBD":
			funcs := converters(fahrenheit2Celsius, d.units.Temperature)
			d.parsers = append(d.parsers, func(r *Record, value string) error {
				return parseOBDCoolantTemp(r.InitOBD(), value, funcs...)
			})
		case "Engine Coolant Temp (C) *OBD":
			funcs := converters(d.units.Temperature)
			d.parsers = append(d.parsers, func(r *Record, value string) error {
				return parseOBDCoolantTemp(r.InitOBD(), value, funcs...)
			})
		case "Intake Air Temp (F) *OBD":
			funcs := converters(fahrenheit2Celsius, d.units.Temperature)
			d.parsers = append(d.parsers, func(r *Record, value string) error {
				return parseOBDIntakeTemp(r.InitOBD(), value, funcs...)
			})
		case "Intake Air Temp (C) *OBD":
			funcs := converters(d.units.Temperature)
			d.parsers = append(d.parsers, func(r *Record, value string) error {
				return parseOBDIntakeTemp(r.InitOBD(), value, funcs...)
			})
		case "Intake Manifold Pressure (PSI) *OBD":
			funcs := converters(psi2Kpa, d.units.Pressure)
			d.parsers = append(d.parsers, func(r *Record, value string) error {
				return parseOBDManifoldPressure(r.InitOBD(), value, funcs...)
			})
		case "Intake Manifold Pressure (kPa) *OBD":
			funcs := converters(d.units.Pressure)
			d.parsers = append(d.parsers, func(r *Record, value string) error {
				return parseOBDManifoldPressure(r.InitOBD(), value, funcs...)
			})
		default:
			return fmt.Errorf("unknown metric %q", col)
		}
	}

	return nil
}

// process creates a new lap record adding it to lap.
func (d *Decoder) process(line int, lap *Lap, data []string) error {
	r := Record{}
	for i, f := range d.parsers {
		// No need to check out of bounds as encoding/csv does that for us.
		if err := f(&r, data[i]); err != nil {
			return fmt.Errorf("line %d: %w", line, err)
		}
	}

	lap.Records = append(lap.Records, r)

	return nil
}

// Decode reads CSV encoded data from its input and stores it in s.
func (d *Decoder) Decode() (*Session, error) {
	s := NewSession()
	lap := &Lap{}
	s.Laps = append(s.Laps, lap)

	sc := bufio.NewScanner(d.r)

	var buf bytes.Buffer
	csvReader := csv.NewReader(&buf)
	n := 1
	for sc.Scan() {
		line := sc.Text()
		switch {
		case strings.HasPrefix(line, "# "):
			if err := d.parseMetadata(s, lap, line); err != nil {
				return nil, err
			}

			// Update lap which might have changed.
			lap = s.Laps[len(s.Laps)-1]
		default:
			// CSV Header / data.
			if _, err := buf.WriteString(line + "\n"); err != nil {
				return nil, fmt.Errorf("write line: %w", err)
			}

			rec, err := csvReader.Read()
			if err != nil {
				return nil, fmt.Errorf("csv read: %w", err)
			}

			if d.parsers == nil {
				// First row must be the column names.
				if err := d.columns(s, rec); err != nil {
					return nil, err
				}
			} else if err = d.process(n, lap, rec); err != nil {
				return nil, err
			}
		}
		n++
	}

	if err := sc.Err(); err != nil {
		return nil, fmt.Errorf("file scan: %w", err)
	}

	return s, nil
}

func (d *Decoder) parseMetadata(s *Session, lap *Lap, line string) error {
	// Metadata for example:
	// # RaceRender Data: TrackAddict 4.8.0 on iOS 15.5 [iPhone12,1] (Mode: 0)
	// # Vehicle: 2019 McLaren 720S
	// # Vehicle Tune: OSID's: 14MA891CP.01., CVN's: C7C38B43
	// # End Point: 50.857952, -0.752617  @ -1.00 deg
	// # GPS: iOS; Type: 1
	// # OBD Mode: BLE; ID: "OBDII  v2.2"
	// # OBD Settings: AP1;AF1;RPR0
	// # User Settings: U0;AS1;LT0/1;EC1;VC1;VQ3;VS0
	// # Device Free Space: 23784 MB
	// # Lap 0: 00:02:03.202
	// # Sector 1: 00:02:03.202
	// # Session End
	parts := strings.SplitN(line[2:], ":", 2)
	switch {
	case parts[0] == "End Point":
		// End Point.
		m := endpointRe.FindStringSubmatch(parts[1])
		if len(m) != 4 {
			return fmt.Errorf("unexpected end point: %q", parts[1])
		}

		if err := s.Endpoint.parseCoordinate(m[1], m[2], m[3]); err != nil {
			return err
		}

	case parts[0] == "Vehicle":
		// Vehicle.
		s.Vehicle = strings.TrimSpace(parts[1])
	case strings.HasPrefix(parts[0], "Lap "):
		// End of lap.
		if err := parseLapNumber(lap, parts[0][4:]); err != nil {
			return err
		}

		if len(s.Laps)-1 > lap.Number {
			return fmt.Errorf("unexpected lap %d have %d", lap.Number, len(s.Laps))
		}

		if err := parseLapDuration(lap, strings.TrimSpace(parts[1])); err != nil {
			return err
		}

		lap = &Lap{}
		s.Laps = append(s.Laps, lap)
	case len(parts) == 2:
		// Other metadata.
		s.Metadata[parts[0]] = strings.TrimSpace(parts[1])
	}

	return nil
}
