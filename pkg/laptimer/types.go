package laptimer

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"strconv"
	"strings"
	"time"
)

const (
	fixDateFormat = "02-Jan-06,15:04:05.00"
	lapDateFormat = "02-Jan-06,15:04:05"
)

// DB represents a Harry's LapTimer export including
// laps and vehicles.
type DB struct {
	Name string `xml:"name"`
	Laps []Lap  `xml:"lap"`
	// Vehicles being a slice with one vehcile per looks like a bug in LapTimer
	// but is currently how the data is exported by v24.6.
	Vehicles []Vehicles `xml:"vehicles,omitempty"`
	XMLName  struct{}   `xml:"LapTimerDB"`
}

// NewDB returns a new DB with the default name.
func NewDB() *DB {
	return &DB{Name: "LapTimer Database"}
}

// Lap represents a LapTimer lap.
type Lap struct {
	ID               int              `xml:"index,attr"`
	Date             LapDate          `xml:"date"`
	LapTime          Duration         `xml:"lapTime"`
	Vehicle          string           `xml:"vehicle,omitempty"`
	Track            string           `xml:"track"`
	Intermediates    Intermediates    `xml:"intermediates,omitempty"`
	LapRecordingType LapRecordingType `xml:"lapRecordingType"`
	// OverallDistance should be same as distance in the last
	// Recording -> Fix -> RelativeToStart and is 2D only.
	OverallDistance  Float1dp  `xml:"overallDistance"`
	WeatherCondition int       `xml:"weatherCondition,omitempty"` // TODO(steve): Use enum.
	AmbientTemp      Float1dp  `xml:"ambientTemp,omitempty"`
	AmbientPressure  Float0dp  `xml:"ambientPressure,omitempty"`  // TODO(steve): Is this in hPa?
	RelativeHumidity Float2dp  `xml:"relativeHumidity,omitempty"` // TODO(steve): Percentage?
	Videos           []Video   `xml:"video,omitempty"`
	Note             string    `xml:"note,omitempty"`
	Tags             Tags      `xml:"tags,omitempty"`
	Recording        Recording `xml:"recording"`
}

// Intermediates is a slice of Intermediate readings.
type Intermediates []Intermediate

// MarshalXML implements xml.Marshaler.
func (i Intermediates) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	if len(i) == 0 {
		return e.EncodeElement("", start)
	}

	// LapTimer doesn't encode values so we manually format with newlines and tabs.
	var buf bytes.Buffer
	for _, v := range i {
		fmt.Fprintf(&buf, "\n\t\t\t%s,%.1f", v.Time, v.Distance)
	}

	s := buf.String()
	return e.EncodeElement(s+"\n\t\t", start)
}

// UnmarshalXML implements xml.Unmarshaler.
func (i *Intermediates) UnmarshalXML(dec *xml.Decoder, start xml.StartElement) error {
	var v string
	if err := dec.DecodeElement(&v, &start); err != nil {
		return fmt.Errorf("unmarshal intermediates element: %w", err)
	}

	for _, s := range strings.Split(v, "\n") {
		s = strings.TrimSpace(s)
		if s == "" {
			continue
		}
		parts := strings.SplitN(s, ",", 2)
		if len(parts) != 2 {
			return fmt.Errorf("unmarshal intermediates split %q: %d != 2", s, len(parts))
		}

		var val Intermediate
		if err := val.Time.Parse(parts[0]); err != nil {
			return fmt.Errorf("unmarshal intermediates time parse %q: %w", parts[0], err)
		}

		var err error
		if val.Distance, err = strconv.ParseFloat(parts[1], 64); err != nil {
			return fmt.Errorf("unmarshal intermediates distance parse %q: %w", parts[1], err)
		}

		*i = append(*i, val)
	}

	return nil
}

// Intermediate represents a measured sub section of a Lap.
type Intermediate struct {
	// Time is the time taken for the intermediate.
	Time Duration

	// Distance is the distance in meters for the intermediate.
	Distance float64
}

// Vehicles represents the vehicles configured in Harry's LapTimer.
type Vehicles struct {
	Vehicles []Vehicle `xml:"vehicledef"`
}

// Vehicle represents a configured vehicle.
type Vehicle struct {
	// Index is the index of the vehicle.
	Index int `xml:"index,attr"`

	// Name is the name of the Vehicle.
	Name string `xml:"name"`

	// PowerLoss is the total power loss in percentage.
	PowerLoss Float `xml:"powerloss,omitempty"`

	// DragCoefficient is the vehicles aerodynamic resistance.
	DragCoefficient Float `xml:"cW,omitempty"`

	// UnladenWeight is the curb weight in Kilograms including
	// fluids and fuel but not loaded with passengers or cargo.
	UnladenWeight Float `xml:"unladenWeight,omitempty"`

	// Payload is the carrying weight on top of curb weight made
	// up of passengers and cargo.
	Payload Float `xml:"payload,omitempty"`

	// GrossVehicelMass (GVM) is the maximum operating weight as
	// specified by the manufacturer.
	GrossVehicleMass Float `xml:"grossVehicleMass,omitempty"`

	// MaxTorqueRPM is the engine RPM at which maximum torque is
	// developed.
	MaxTorqueRPM int `xml:"maxTorqueRPM,omitempty"`

	// MaxTorque is maximum torque developed by the engine.
	MaxTorque int `xml:"maxTorque,omitempty"`

	// MaxPowerRPM is the engine RPM at which maximum power is
	// developed.
	MaxPowerRPM int `xml:"maxPowerRPM,omitempty"`

	// MaxPower is maximum power developed by the engine.
	MaxPower int `xml:"maxPower,omitempty"`

	// Vin is the Vehicle Identification Number.
	Vin string `xml:"vin,omitempty"`

	// Gears specifies the gears of the vehicle.
	Gears *Gears `xml:"gears,omitempty"`

	// DriveRatio is the relationshup between the gear box's and
	// the axle's RPM.
	DriveRatio Float `xml:"driveRatio,omitempty"`

	// FrontWheel specifies the type of tyres used on the front wheels.
	FrontWheels *Tyre `xml:"frontWheels,omitempty"`
	// RearWheels specifies the type of tyres used on the rear wheels.
	RearWheels *Tyre `xml:"rearWheels,omitempty"`

	// DriveWheels specifies the wheels which drive the vehicle
	DriveWheels DriveWheels `xml:"driveWheels,omitempty"`

	// IntakeType specifies engine intake type.
	IntakeType IntakeType `xml:"intakeType,omitempty"`

	// EngineType specifies the type of the engine powering the vehicle.
	EngineType EngineType `xml:"engineType,omitempty"`

	// MaxRPM is the maximum RPM of the engine.
	MaxRPM int `xml:"maxRPM,omitempty"`

	// ShiftGearThreshold is the percentage of the engine max RPM at which
	// the Shift Gear Flash gets activated.
	ShiftGearThreshold Threshold `xml:"shiftGearThreshold,omitempty"`

	// VolumetricEfficiency is the efficiency in percentage with which
	// the engine can move charge into and out of its cylinders.
	VolumetricEfficiency int `xml:"volumetricEfficiency,omitempty"`

	// Displacement is the engine displacement in cubic centimetres.
	Displacement int `xml:"displacement,omitempty"`

	// TankVolume is the fuel take capacity including reserve.
	TankVolume Float `xml:"tankVolume,omitempty"`

	// Axles is the number of axles of the vehicle.
	Axles int `xml:"axles,omitempty"`

	// BodyWidth is the width of the vehicle excluding mirrors.
	BodyWidth int `xml:"bodyWidth,omitempty"`

	// BodyLength is the full length of the vehicle.
	BodyLength int `xml:"bodyLength,omitempty"`

	// BodyHeight is the height of the vehicle with the doors closed.
	BodyHeight int `xml:"bodyHeight,omitempty"`

	// Wheelbase is the distance between the cents of the front and
	// rear wheels.
	Wheelbase int `xml:"wheelbase,omitempty"`

	// Style is the specific model style of the vehicle.
	// e.g. 'GT3 RS' for a Porsche 991 GT3 RS or a 'V6 Coupe' for
	// a Ford Mustang V6 Coupe.
	Style string `xml:"style,omitempty"`

	// Make is the name of manufacture.
	Make string `xml:"make,omitempty"`

	// Model is the model of the vehicle, e.g. '911' for a Porsche 991
	// or 'Mustang' for a Ford Mustang.
	Model string `xml:"model,omitempty"`

	// Year is the full year the vehicle was manufactured.
	Year int `xml:"year,omitempty"`

	// ContryCode is the two digit country code for the vehicle.
	CountryCode string `xml:"countryCode,omitempty"`

	// TurningCircle is the minimum turning circle in meters.
	TurningCircle Float1dp `xml:"turningCircle,omitempty"`

	// OverhangFront is the length in meters which extends beyond the
	// wheelbase at the front.
	OverhangFront int `xml:"overhangFront,omitempty"`

	// VehicleID is the ID of a vehicle submitted to LapTimer.
	VehicleID int `xml:"vehicleID,omitempty"`

	// OriginalContributor is true if this user submitted the vehicle,
	// false otherwise.
	OriginalContributor bool `xml:"originalContributor,omitempty"`
}

// Threshold represents shift gear threshold.
type Threshold int

// MarshalXML implements xml.Marshaler.
func (t Threshold) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	return e.EncodeElement(fmt.Sprintf("%d%%", int(t)), start)
}

// UnmarshalXML implements xml.Unmarshaler.
func (t *Threshold) UnmarshalXML(dec *xml.Decoder, start xml.StartElement) error {
	var v string
	if err := dec.DecodeElement(&v, &start); err != nil {
		return fmt.Errorf("unmarshal threshold element: %w", err)
	}

	d, err := strconv.Atoi(strings.TrimSuffix(v, "%"))
	if err != nil {
		return fmt.Errorf("unmarshal threshold atoi: %w", err)
	}

	*t = Threshold(d)

	return nil
}

// Gears represents the set of gears for a vehicle.
type Gears struct {
	Gears []Gear `xml:"gear"`
}

// Gear represents an individual gear of a vehicle.
type Gear struct {
	// Number is the number of the gear.
	Number int

	// Ratio is the relationship between the engine's and the
	// gear box's RPM.
	Ratio float64
}

// MarshalXML implements xml.Marshaler.
func (g Gear) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	return e.EncodeElement(fmt.Sprintf("%d,%f", g.Number, g.Ratio), start)
}

// UnmarshalXML implements xml.Unmarshaler.
func (g *Gear) UnmarshalXML(dec *xml.Decoder, start xml.StartElement) error {
	var v string
	if err := dec.DecodeElement(&v, &start); err != nil {
		return fmt.Errorf("unmarshal gear element: %w", err)
	}

	n, err := fmt.Sscanf(v, "%d,%f", &g.Number, &g.Ratio)
	if err != nil {
		return fmt.Errorf("unmarshal gear scan %q: %w", v, err)
	} else if n != 2 {
		return fmt.Errorf("unmarshal gear scan %q: %d != 2", v, n)
	}

	return nil
}

// Tyre represents a tyre specification.
type Tyre struct {
	// Width is the width of the tyre.
	Width int

	// Profile is the height of the sidewall as a percentage
	// of is width.
	Profile int

	// Size is the diameter of the tyre in inches.
	Size int

	// SpeedRating is the speed rating for the tyre.
	SpeedRating string
}

// MarshalXML implements xml.Marshaler.
func (t Tyre) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	return e.EncodeElement(
		fmt.Sprintf("%d / %d %s %d",
			t.Width,
			t.Profile,
			t.SpeedRating,
			t.Size,
		),
		start,
	)
}

// UnmarshalXML implements xml.Unmarshaler.
func (t *Tyre) UnmarshalXML(dec *xml.Decoder, start xml.StartElement) error {
	var v string
	if err := dec.DecodeElement(&v, &start); err != nil {
		return fmt.Errorf("unmarshal tyre element: %w", err)
	}

	n, err := fmt.Sscanf(v, "%d / %d %s %d", &t.Width, &t.Profile, &t.SpeedRating, &t.Size)
	if err != nil {
		return fmt.Errorf("unmarshal tyre scan %q: %w", v, err)
	} else if n != 4 {
		return fmt.Errorf("unmarshal tyre scan %q: %d != 4", v, n)
	}

	return nil
}

// Tags represents the tags applied to a Lap.
type Tags []string

// MarshalXML implements xml.Marshaler.
func (t Tags) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	return e.EncodeElement(strings.Join([]string(t), ","), start)
}

// UnmarshalXML implements xml.Unmarshaler.
func (t *Tags) UnmarshalXML(dec *xml.Decoder, start xml.StartElement) error {
	var v string
	if err := dec.DecodeElement(&v, &start); err != nil {
		return fmt.Errorf("unmarshal tags element: %w", err)
	}

	*t = strings.Split(v, ",")

	return nil
}

// Video represents a video for a Lap.
type Video struct {
	// Overlaid is true if this video has been overlaid.
	Overlaid bool `xml:"overlaid"`

	// URL is the internal URL of the video.
	URL string `xml:"url"`

	// SyncPoint represents the offset in seconds from
	// the start of the video which the Lap starts.
	SyncPoint SyncPoint `xml:"syncPoint"`
}

// Recording represents the data recorded for a Lap.
type Recording struct {
	Fixes []Fix `xml:"fix"`
}

// Fix represents a single data entry for a Lap.
type Fix struct {
	// ID of the fix, this value isn't imported.
	ID int `xml:"index,attr"`

	// Date is the date time of the Fix.
	Date FixDate `xml:"date"`

	// Coordinate is the coordinate including altitude.
	Coordinate AltitudeCoordinate `xml:"coordinate"`

	// Speed is the vehicle speed based on GPS.
	Speed Float1dp `xml:"speed"`

	// Position indicates the state and method of the
	// GSP coordinates.
	Positioning Positioning `xml:"positioning"`

	// Satellites indicates how many satellites were visible.
	Satellites int `xml:"satellites"`

	// Direction indicates the direction of travel in degrees.
	Direction Float1dp `xml:"direction"`

	// Hdop is the GPS Horizontal Dilution of Precision.
	// 1: Ideal
	// 1-2: Excellent
	// 2-5: Good
	// 5-10: Moderate
	// 10-20: Fair
	// >20: Poor
	// Default to 1 as advised by Harry here:
	// http://forum.gps-laptimer.de/viewtopic.php?f=9&t=6134
	Hdop Float2dp `xml:"hdop"`

	// Accuracy is the GPS accuracy in meters.
	// Typically either Hdop or Accuracy is specified but not both.
	Accuracy Float1dp `xml:"accuracy"`

	// RelativeToStart provides values relative to the start of the Lap.
	RelativeToStart RelativeToStart `xml:"relativeToStart"`

	// Acceleration provides acceleration relative to the
	// direction of travel.
	Acceleration *Acceleration `xml:"acceleration,omitempty"`

	// OBD provides OBD data if available and is
	// interpolated to GPS storage rate.
	OBD *OBD `xml:"obd,omitempty"`

	// TPMS provides tyre pressure monitoring data.
	TPMS *TPMS `xml:"tpms,omitempty"`
}

// Coordinate represents a GPS coordinate without altitude.
type Coordinate struct {
	Latitude  float64
	Longitude float64
}

// MarshalXML implements xml.Marshaler.
func (c Coordinate) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	return e.EncodeElement(
		fmt.Sprintf("%.08f,%.08f",
			c.Latitude,
			c.Longitude,
		),
		start,
	)
}

// UnmarshalXML implements xml.Unmarshaler.
func (c *Coordinate) UnmarshalXML(dec *xml.Decoder, start xml.StartElement) error {
	var v string
	if err := dec.DecodeElement(&v, &start); err != nil {
		return fmt.Errorf("unmarshal coordinate element: %w", err)
	}

	n, err := fmt.Sscanf(v, "%f,%f", &c.Latitude, &c.Longitude)
	if err != nil {
		return fmt.Errorf("unmarshal coordinate scan %q: %w", v, err)
	} else if n != 2 {
		return fmt.Errorf("unmarshal coordinate scan %q: %d != 2", v, n)
	}

	return nil
}

// AltitudeCoordinate represents a GPS coordinate including altitude.
type AltitudeCoordinate struct {
	Coordinate

	// Altitude is the distance above sea level in meters.
	Altitude float64
}

// MarshalXML implements xml.Marshaler.
func (c AltitudeCoordinate) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	return e.EncodeElement(
		fmt.Sprintf("%.08f,%.08f,%.1f",
			c.Latitude,
			c.Longitude,
			c.Altitude,
		),
		start,
	)
}

// UnmarshalXML implements xml.Unmarshaler.
func (c *AltitudeCoordinate) UnmarshalXML(dec *xml.Decoder, start xml.StartElement) error {
	var v string
	if err := dec.DecodeElement(&v, &start); err != nil {
		return fmt.Errorf("unmarshal altitude coordinate element: %w", err)
	}

	n, err := fmt.Sscanf(v, "%f,%f,%f", &c.Latitude, &c.Longitude, &c.Altitude)
	if err != nil {
		return fmt.Errorf("unmarshal altitude coordinate scan %q: %w", v, err)
	} else if n != 3 {
		return fmt.Errorf("unmarshal altitude coordinate scan %q: %d != 3", v, n)
	}

	return nil
}

// Positioning represents the state of positioning information.
type Positioning struct {
	// DifferentialStatus identifies the differential GPS status of this Fix.
	DifferentialStatus DifferentialStatus

	// PositionFixing identifies how the Fix was obtained.
	PositionFixing PositionFixing

	// Interpolated is true for any Fix calculated by LapTimer using
	// interpolation to approximate a trigger line has passed, false
	// if read directly from the GPS.
	Interpolated bool
}

// MarshalXML implements xml.Marshaler.
func (p Positioning) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	var inter int
	if p.Interpolated {
		inter = 1
	}
	return e.EncodeElement(
		fmt.Sprintf("%d,%d,%d",
			p.DifferentialStatus,
			p.PositionFixing,
			inter,
		),
		start,
	)
}

// UnmarshalXML implements xml.Unmarshaler.
func (p *Positioning) UnmarshalXML(dec *xml.Decoder, start xml.StartElement) error {
	var v string
	if err := dec.DecodeElement(&v, &start); err != nil {
		return fmt.Errorf("unmarshal positioning element: %w", err)
	}

	var di, po int
	n, err := fmt.Sscanf(v, "%d,%d,%t", &di, &po, &p.Interpolated)
	if err != nil {
		return fmt.Errorf("unmarshal positioning scan %q: %w", v, err)
	} else if n != 3 {
		return fmt.Errorf("unmarshal positioning scan %q: %d != 3", v, n)
	}

	p.DifferentialStatus = DifferentialStatus(di)
	p.PositionFixing = PositionFixing(po)

	return nil
}

// RelativeToStart represents the time a distance relative to the start of the Lap.
type RelativeToStart struct {
	// Distance in meters in 2D only.
	Distance float64

	// Offset is the time since the start of the Lap.
	Offset Duration
}

// MarshalXML implements xml.Marshaler.
func (r RelativeToStart) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	return e.EncodeElement(
		fmt.Sprintf("%.1f,%s", r.Distance, r.Offset.String()),
		start,
	)
}

// UnmarshalXML implements xml.Unmarshaler.
func (r *RelativeToStart) UnmarshalXML(dec *xml.Decoder, start xml.StartElement) error {
	var v string
	if err := dec.DecodeElement(&v, &start); err != nil {
		return fmt.Errorf("unmarshal relative to start element: %w", err)
	}

	parts := strings.SplitN(v, ",", 2)
	if len(parts) != 2 {
		return fmt.Errorf("unmarshal relative to start split %q: %d != 2", v, len(parts))
	}

	f, err := strconv.ParseFloat(parts[0], 64)
	if err != nil {
		return fmt.Errorf("unmarshal relative to start parse %q: %w", parts[0], err)
	}

	r.Distance = f

	return r.Offset.Parse(parts[1])
}

// Acceleration represents the acceleration G values relative to driving direction.
// See the following for details:
// http://www.gps-laptimer.com/LapTimerDocumentation%20-%20Acceleration%20Chapter.pdf
type Acceleration struct {
	// Source is the source of the acceleration measurement.
	Source int `xml:"source"` // TODO(steve): Use an enum.

	// Lateral is the lateral acceleration in G as a result of turning.
	Lateral Float2dp `xml:"lateral"`

	// Lineral is the lineral acceleration in G as a result of accelerating or braking.
	Lineral Float2dp `xml:"lineal"`

	// Coordinate GPS position for the red/yellow/green display of
	// lateral acceleration.
	Coordinate Coordinate `xml:"coordinate"`
}

// LapDate represents the data and time of a Fix.
type LapDate time.Time

func (d LapDate) String() string {
	return strings.ToUpper(time.Time(d).UTC().Format(lapDateFormat))
}

// MarshalXML implements xml.Marshaler.
func (d LapDate) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	return e.EncodeElement(d.String(), start)
}

// UnmarshalXML implements xml.Unmarshaler.
func (d *LapDate) UnmarshalXML(dec *xml.Decoder, start xml.StartElement) error {
	var v string
	if err := dec.DecodeElement(&v, &start); err != nil {
		return fmt.Errorf("unmarshal lap date element: %w", err)
	}

	t, err := time.Parse(lapDateFormat, v)
	if err != nil {
		return fmt.Errorf("unmarshal lap date parse: %w", err)
	}

	*d = LapDate(t)

	return nil
}

// FixDate represents the data and time of a Fix.
type FixDate time.Time

func (d FixDate) String() string {
	return strings.ToUpper(time.Time(d).UTC().Format(fixDateFormat))
}

// MarshalXML implements xml.Marshaler.
func (d FixDate) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	return e.EncodeElement(d.String(), start)
}

// UnmarshalXML implements xml.Unmarshaler.
func (d *FixDate) UnmarshalXML(dec *xml.Decoder, start xml.StartElement) error {
	var v string
	if err := dec.DecodeElement(&v, &start); err != nil {
		return fmt.Errorf("unmarshal fix date element: %w", err)
	}

	t, err := time.Parse(fixDateFormat, v)
	if err != nil {
		return fmt.Errorf("unmarshal fix date parse: %w", err)
	}

	*d = FixDate(t)

	return nil
}

// Duration represents a time duration.
type Duration time.Duration

// String implements Stringer.
func (d Duration) String() string {
	dur := time.Duration(d)
	m := dur / time.Minute
	mv := m * time.Minute // nolint: durationcheck
	s := (dur - mv) / time.Second
	sv := s * time.Second // nolint: durationcheck
	cs := (dur - mv - sv) / time.Millisecond / 10

	return fmt.Sprintf("%02d:%02d.%02d",
		m,
		s,
		cs,
	)
}

// MarshalXML implements xml.Marshaler.
func (d Duration) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	return e.EncodeElement(d.String(), start)
}

// UnmarshalXML implements xml.Unmarshaler.
func (d *Duration) UnmarshalXML(dec *xml.Decoder, start xml.StartElement) error {
	var v string
	if err := dec.DecodeElement(&v, &start); err != nil {
		return fmt.Errorf("unmarshal duration element: %w", err)
	}

	return d.Parse(v)
}

// Parse parses value into d.
func (d *Duration) Parse(value string) error {
	var m, s, cs int
	n, err := fmt.Sscanf(value, "%d:%d.%d", &m, &s, &cs)
	if err != nil {
		return fmt.Errorf("unmarshal duration scan %q: %w", value, err)
	} else if n != 3 {
		return fmt.Errorf("unmarshal duration scan %q: %d != 3", value, n)
	}

	*d = Duration(
		(time.Duration(m) * time.Minute) +
			(time.Duration(s) * time.Second) +
			(time.Duration(cs*10) * time.Millisecond),
	)

	return nil
}

// OBD represents On-Board Diagnostics from a vehicle.
// Not all values are provided by all cars / sensors.
// Info from:
// http://forum.gps-laptimer.de/viewtopic.php?f=9&t=6134
// All temperatures are in Celsius.
type OBD struct {
	// EngineRPM is the RPM of the engine.
	EngineRPM *int `xml:"rpm,omitempty"`

	// MassAirFlow is mass air flow in grams per second.
	// If not defined map is required.
	MassAirFlow *Float2dp `xml:"maf,omitempty"`

	// VehicleSpeed in Kmh.
	// VSS: Vehicle Speed Sensor, WheelSpeed.
	VehicleSpeed *Float1dp `xml:"speed"`

	// Throttle is the throttle position in percentage as reported
	// by the Throttle Position Sensor (TPS).
	Throttle *Float2dp `xml:"throttle"`

	// FuelLevel in percentage.
	FuelLevel *Float2dp `xml:"fuelLevel,omitempty"`

	// CoolantTemp is the engine coolant temperature in Celsius as
	// reported by Engine Coolant Temperature (ECT).
	CoolantTemp *Float1dp `xml:"coolant,omitempty"`

	// OilTemp is the engine oil temperature in Celsius as
	// reported by Engine Oil Temperature (EOT).
	OilTemp *Float1dp `xml:"oil,omitempty"`

	// ManifoldAbsolutePressure is a pressure of the manifold as
	// reported by Manifold Absolute Pressure (MAP).
	// For turbo / super charged engines.
	ManifoldAbsolutePressure *Float2dp `xml:"map,omitempty"`

	// IntakeAirTemperature is the temperature in Celsius of
	// the air at the intake as reported by IAT.
	IntakeAirTemperature *Float0dp `xml:"iat,omitempty"`

	// TODO(steve): correct all the tags below here which are a guess.

	// Gear is the selected gear which may be calculated.
	Gear *int `xml:"gear,omitempty"`

	// OilPressure is the oil pressure in Bar not Kpa.
	OilPressure *float64 `xml:"oilp,omitempty"`

	// BrakePressure is the applied break pressure in Bar not Kpa.
	BrakePressure *float64 `xml:"brake,omitempty"`

	// WheelSpeedRearLeft is the speed in RPM of the rear left wheel.
	WheelSpeedRearLeft *float64 `xml:"speedrl,omitempty"`

	// WheelSpeedReadRight is the speed in RPM of the rear right wheel.
	WheelSpeedRearRight *float64 `xml:"speedrr,omitempty"`

	// WheelSpeedFrontLeft is the speed in RPM of the front left wheel.
	WheelSpeedFrontLeft *float64 `xml:"speedfl,omitempty"`

	// WheelSpeedFrontLeft is the speed in RPM of the front right wheel.
	WheelSpeedFrontRight *float64 `xml:"speedfr,omitempty"`

	// YawRate is the angular velocity of rotation.
	YawRate *float64 `xml:"yaw,omitempty"`

	// Odometer measures the distance moved in meters.
	Odometer *float64 `xml:"odometer,omitempty"`

	// SteerAngle is the angle of steering wheel in degrees.
	SteerAngle *float64 `xml:"steerangle,omitempty"`

	// SteerRate is the rate of steering wheel angle change.
	SteerRate *float64 `xml:"steerrate,omitempty"`

	// FixType TODO(steve): describe properly and check type.
	FixType *FixType `xml:"fixtype,omitempty"`

	// SupportLevel TODO(steve): describe properly and check type.
	SupportLevel *int `xml:"supportlevel,omitempty"`

	// DriverPower TODO(steve): describe properly.
	DriverPower *float64 `xml:"driverpower,omitempty"`

	// BatteryTemp is the temperature in Celsius of the battery.
	BatteryTemp *float64 `xml:"batterytemp,omitempty"`

	// BatteryVoltage is the voltage of battery.
	BatteryVoltage *float64 `xml:"voltage,omitempty"`

	// ConsumedPower TODO(steve): describe properly.
	ConsumedPower *float64 `xml:"consumedpower,omitempty"`
}

// Float2dp represents a float that is output with 2 decimal places of precision.
type Float2dp float64

// MarshalXML implements xml.Marshaler.
func (f Float2dp) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	return e.EncodeElement(fmt.Sprintf("%.02f", f), start)
}

// Float represents a float that is output with 6 decimal places of precision.
type Float float64

// MarshalXML implements xml.Marshaler.
func (f Float) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	return e.EncodeElement(fmt.Sprintf("%f", f), start)
}

// Float0dp represents a float that is output with 0 decimal places of precision.
type Float0dp float64

// MarshalXML implements xml.Marshaler.
func (f Float0dp) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	return e.EncodeElement(fmt.Sprintf("%.0f", f), start)
}

// Float1dp represents a float that is output with 1 decimal places of precision.
type Float1dp float64

// MarshalXML implements xml.Marshaler.
func (f Float1dp) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	return e.EncodeElement(fmt.Sprintf("%.01f", f), start)
}

// SyncPoint represents the time from the start of a video to sync.
type SyncPoint time.Duration

// MarshalXML implements xml.Marshaler.
func (s SyncPoint) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	return e.EncodeElement(fmt.Sprintf("%.02f", time.Duration(s).Seconds()), start)
}

// UnmarshalXML implements xml.Unmarshaler.
func (s *SyncPoint) UnmarshalXML(dec *xml.Decoder, start xml.StartElement) error {
	var v string
	if err := dec.DecodeElement(&v, &start); err != nil {
		return fmt.Errorf("unmarshal duration element: %w", err)
	}

	var sec, cs int64
	n, err := fmt.Sscanf(v, "%d.%d", &sec, &cs)
	if err != nil {
		return fmt.Errorf("unmarshal sync point scan %q: %w", v, err)
	} else if n != 2 {
		return fmt.Errorf("unmarshal sync point scan %q: %d != 2", v, n)
	}

	*s = SyncPoint(sec*int64(time.Second) + cs*10*int64(time.Millisecond))

	return nil
}

// TPMS represents Tyre Pressure Monitoring System data.
type TPMS struct {
	Tyres []TPMSTyre `xml:"tire"`
}

// TPMSTyre represents Tyre Pressure Monitoring System tyre.
type TPMSTyre struct {
	Position    TyrePosition `xml:"tirePosition"`
	Temperature Float1dp     `xml:"temperatures"`
}
