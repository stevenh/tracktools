package convert

import (
	"fmt"
	"time"

	"github.com/stevenh/tracktools/pkg/laptimer"
	"github.com/stevenh/tracktools/pkg/trackaddict"
	"github.com/tidwall/geodesic"
	"gonum.org/v1/gonum/interp"
)

// TrackAddict converts between trackaddict format
// and LapTimer format.
type TrackAddict struct {
	track      string
	vehicle    string
	tags       laptimer.Tags
	note       string
	hdop       laptimer.Float2dp
	satellites int
	predictor  interp.FittablePredictor
	diffStatus laptimer.DifferentialStatus
	posFixing  laptimer.PositionFixing
	startDate  time.Time
	dateAdjust time.Duration
}

// Option represents a TrackAddict option.
type Option func(*TrackAddict) error

// TrackOpt sets the Track for a TrackAddict.
// Default is a blank string.
func TrackOpt(name string) Option {
	return func(ta *TrackAddict) error {
		ta.track = name

		return nil
	}
}

// VehicleOpt sets the Vehicle output of a TrackAddict.
// Default is empty, using the value from the source.
func VehicleOpt(name string) Option {
	return func(ta *TrackAddict) error {
		ta.vehicle = name

		return nil
	}
}

// TagsOpt sets tags used in the output of a TrackAddict.
// Default is empty.
func TagsOpt(tags ...string) Option {
	return func(ta *TrackAddict) error {
		ta.tags = laptimer.Tags(tags)

		return nil
	}
}

// PredictorOpt sets a predictor used to interpolated between OBD data
// for Records with have updated GPS data but not updated OBD data
// in the output of a TrackAddict.
// Default is a interp.PiecewiseLinear.
// If set to nil no interpolation will be done.
func PredictorOpt(predictor interp.FittablePredictor) Option {
	return func(ta *TrackAddict) error {
		ta.predictor = predictor

		return nil
	}
}

// NoteOpt sets a note used in the output of a TrackAddict.
// Default is blank string.
func NoteOpt(value string) Option {
	return func(ta *TrackAddict) error {
		ta.note = value

		return nil
	}
}

// DifferentialOpt sets a differential status used for fixes in the output of a TrackAddict.
// Default is laptimer.DifferentialStatusUnknown.
func DifferentialOpt(value laptimer.DifferentialStatus) Option {
	return func(ta *TrackAddict) error {
		ta.diffStatus = value

		return nil
	}
}

// PositionOpt sets a position fixing used for fixes in the output of a TrackAddict.
// Default is PositionFixing3D.
func PositionOpt(value laptimer.PositionFixing) Option {
	return func(ta *TrackAddict) error {
		ta.posFixing = value

		return nil
	}
}

// StartDateOpt overrides the start date for times in the output of a TrackAddict.
func StartDateOpt(date time.Time) Option {
	return func(ta *TrackAddict) error {
		ta.startDate = date

		return nil
	}
}

// NewTrackAddict creates a new TrackAddict with a given set of options.
func NewTrackAddict(options ...Option) (*TrackAddict, error) {
	c := &TrackAddict{
		predictor:  &interp.PiecewiseLinear{},
		hdop:       1,
		diffStatus: laptimer.DifferentialStatusUnknown,
		posFixing:  laptimer.PositionFixing3D,
	}
	for _, f := range options {
		if err := f(c); err != nil {
			return nil, err
		}
	}

	return c, nil
}

// LapTimer returns s converted to a laptimer.DB.
func (ta *TrackAddict) LapTimer(s *trackaddict.Session) (*laptimer.DB, error) {
	db := laptimer.NewDB()
	vehicle := s.Vehicle
	if ta.vehicle != "" {
		vehicle = ta.vehicle
	}

	if ta.predictor != nil {
		if err := s.PredictOBD(ta.predictor); err != nil {
			return nil, fmt.Errorf("predict: %w", err)
		}
	}

	// First lap is the outlap and last is the in lap.
	// TODO(steve): should we filter it better?
	if len(s.Laps) < 3 {
		return db, nil
	}

	fixID := 1
	for _, l := range s.Laps[1 : len(s.Laps)-1] {
		lap := ta.lapTimerLap(l, vehicle, fixID)
		db.Laps = append(db.Laps, lap)
		fixID += len(lap.Recording.Fixes)
	}

	return db, nil
}

// lapTimerLap returns laptimer.Lap representation of l.
func (ta *TrackAddict) lapTimerLap(l *trackaddict.Lap, vehicle string, id int) laptimer.Lap {
	lap := laptimer.Lap{
		ID:               l.Number,
		LapTime:          laptimer.Duration(l.Duration),
		Vehicle:          vehicle,
		Track:            ta.track,
		Tags:             ta.tags,
		Note:             ta.note,
		LapRecordingType: laptimer.LapRecordingTriggered,
	}

	if len(l.Records) == 0 {
		return lap
	}

	var dist, d float64
	r := l.Records[0]
	lastGPS := r.GPS
	firstNow := r.Now

	if !ta.startDate.IsZero() && ta.dateAdjust == 0 {
		y, m, d := r.Time.UTC().Date()
		t := time.Date(y, m, d, 0, 0, 0, 0, time.UTC)
		ta.dateAdjust = ta.startDate.Sub(t)
	}

	lap.Date = laptimer.LapDate(r.Time.Add(ta.dateAdjust))

	for j, r := range l.Records {
		if j != 0 && !r.GPS.Update {
			continue
		}

		if j > 0 && r.GPS.Update {
			geodesic.WGS84.Inverse(
				lastGPS.Latitude,
				lastGPS.Longitude,
				r.GPS.Latitude,
				r.GPS.Longitude,
				&d, nil, nil,
			)
			dist += d
			lastGPS = r.GPS
		}

		lap.Recording.Fixes = append(lap.Recording.Fixes,
			ta.lapTimerFix(id, dist, r, firstNow),
		)
		id++
	}

	lap.OverallDistance = round1dp(dist)

	return lap
}

// lapTimerFix returns a laptimer.Fix representation of the data in r.
func (ta *TrackAddict) lapTimerFix(id int, dist float64, r trackaddict.Record, firstNow time.Duration) laptimer.Fix {
	f := laptimer.Fix{
		ID:   id,
		Date: laptimer.FixDate(r.Time.Add(ta.dateAdjust)),
		Coordinate: laptimer.AltitudeCoordinate{
			Coordinate: laptimer.Coordinate{
				Longitude: r.GPS.Longitude,
				Latitude:  r.GPS.Latitude,
			},
			Altitude: r.GPS.Altitude,
		},
		Speed: round1dp(r.Speed),
		Positioning: laptimer.Positioning{
			DifferentialStatus: ta.diffStatus,
			PositionFixing:     ta.posFixing,
			Interpolated:       false,
		},
		Satellites: ta.satellites,
		Direction:  round1dp(r.GPS.Heading),
		Hdop:       ta.hdop,
		Accuracy:   round1dp(r.GPS.Accuracy),
		RelativeToStart: laptimer.RelativeToStart{
			Distance: dist,
			Offset:   laptimer.Duration(r.Now - firstNow),
		},
	}

	ta.lapTimerAccel(&f, r)
	ta.lapTimerOBD(&f, r)

	return f
}

// lapTimerAccell adds acceleration data from r to f if available.
func (ta *TrackAddict) lapTimerAccel(f *laptimer.Fix, r trackaddict.Record) {
	if r.Accel == nil {
		return
	}

	f.Acceleration = &laptimer.Acceleration{
		// TODO(steve): Calculate correctly according to:
		// http://www.gps-laptimer.com/LapTimerDocumentation%20-%20Acceleration%20Chapter.pdf
		Source:  0, // TODO(steve): set correctly.
		Lateral: round2dp(r.Accel.X),
		Lineal:  round2dp(r.Accel.Y),
		Coordinate: laptimer.Coordinate{
			Longitude: r.GPS.Longitude,
			Latitude:  r.GPS.Latitude,
		},
	}
}

// lapTimerOBD adds OBD data from r to f if available.
func (ta *TrackAddict) lapTimerOBD(f *laptimer.Fix, r trackaddict.Record) {
	if r.OBD == nil {
		return
	}

	f.OBD = &laptimer.OBD{}
	if r.OBD.EngineSpeed != nil {
		v := roundint(*r.OBD.EngineSpeed)
		f.OBD.EngineRPM = &v
	}

	if r.OBD.ManifoldPressure != nil {
		v := round2dp(*r.OBD.ManifoldPressure)
		f.OBD.ManifoldAbsolutePressure = &v
	}

	if r.OBD.Speed != nil {
		v := round1dp(*r.OBD.Speed)
		f.OBD.VehicleSpeed = &v
	}

	if r.OBD.Throttle != nil {
		v := round2dp(*r.OBD.Throttle)
		f.OBD.Throttle = &v
	}

	if r.OBD.CoolantTemp != nil {
		v := round1dp(*r.OBD.CoolantTemp)
		f.OBD.CoolantTemp = &v
	}

	if r.OBD.IntakeTemp != nil {
		v := round0dp(*r.OBD.IntakeTemp)
		f.OBD.IntakeAirTemperature = &v
	}
}
