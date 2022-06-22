package trackaddict

import (
	"fmt"
	"reflect"
	"time"

	"gonum.org/v1/gonum/interp"
)

// Session represents a track session.
type Session struct {
	Laps     []*Lap
	Metadata map[string]string
	Vehicle  string
	Endpoint GPS
}

// NewSession returns a new initialised session.
func NewSession() *Session {
	return &Session{
		Metadata: make(map[string]string, 9),
		Laps:     make([]*Lap, 0, 1),
	}
}

type obdNeeded struct {
	obd *OBD
	x   float64
}

// PredictOBD interpolates the OBD values for each GPS record using
// given predictor type.
func (s *Session) PredictOBD(predictor interp.FittablePredictor) error {
	var (
		start  time.Time
		ys     [][]float64
		xs     []float64
		needed []obdNeeded
	)

	for _, l := range s.Laps {
		for _, r := range l.Records {
			switch {
			case r.OBD != nil && r.OBD.Update:
				// OBD value was updated, capture its values for prediction.
				if start.IsZero() {
					start = r.Time
					xs = append(xs, 0)
				} else {
					xs = append(xs, r.Time.Sub(start).Seconds())
				}

				ys = r.OBD.appendValues(ys)
			case r.GPS.Update:
				// GPS value was update, capture it for update.
				needed = append(needed, obdNeeded{
					obd: r.OBD,
					x:   r.Time.Sub(start).Seconds(),
				})
			}
		}
	}

	if len(needed) == 0 {
		// No needed values so return early.
		return nil
	}

	// Create a predictor for each OBD field so we don't have
	// to keep reinitialising the predictor.
	predictors := make([]interp.FittablePredictor, len(ys))
	for i := range predictors {
		switch i {
		case 0:
			predictors[i] = predictor
		default:
			// nolint: forcetypeassert
			predictors[i] = reflect.New(reflect.ValueOf(predictor).Elem().Type()).
				Interface().(interp.FittablePredictor)
		}

		if err := predictors[i].Fit(xs, ys[i]); err != nil {
			return fmt.Errorf("fit %d: %w", i, err)
		}
	}

	// Predict the needed values and update.
	vals := make([]float64, len(predictors))
	for _, n := range needed {
		for i, p := range predictors {
			vals[i] = p.Predict(n.x)
		}
		n.obd.set(vals)
	}

	return nil
}
