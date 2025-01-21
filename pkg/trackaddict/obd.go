package trackaddict

import "reflect"

// OBD represents On-Board Diagnostic information.
type OBD struct {
	// Update indicates if this value was updated.
	Update bool

	// Speed is vehicle speed.
	Speed *float64

	// RPM is the engine RPM.
	EngineSpeed *float64

	// Throttle is the throttle position percentage.
	Throttle *float64

	// ColantTemp is the coolant temperature.
	CoolantTemp *float64

	// IntakeTemp is the air intake temperature.
	IntakeTemp *float64

	// ManifoldPressure is the Intake Manifold Pressure.
	ManifoldPressure *float64

	// TODO(steve): Other values such as mass air flow.
}

func (o OBD) appendValues(slices [][]float64) [][]float64 {
	v := reflect.ValueOf(o)
	t := v.Type()
	if len(slices) == 0 {
		var j int
		for i := range t.NumField() {
			f := v.Field(i)
			if t.Field(i).Type.Kind() == reflect.Pointer {
				f = f.Elem()
			}

			if f.CanFloat() {
				j++
			}
		}
		slices = make([][]float64, j)
	}

	var j int
	for i := range t.NumField() {
		f := v.Field(i)
		if t.Field(i).Type.Kind() == reflect.Pointer {
			f = f.Elem()
		}

		if f.CanFloat() {
			slices[j] = append(slices[j], f.Float())
			j++
		}
	}

	return slices
}

// set sets the float fields from values.
// panics if len(values) is less than the number of float fields.
func (o *OBD) set(values []float64) {
	v := reflect.ValueOf(o).Elem()
	t := v.Type()

	var j int
	for i := range t.NumField() {
		f := v.Field(i)
		if f.CanFloat() {
			f.SetFloat(values[j])
			j++
		}
	}
}

func parseOBDUpdate(o *OBD, value string) (err error) {
	o.Update, err = parseBool("obd updated", value)
	return err
}

func parseOBDSpeed(o *OBD, value string, converters ...converter) (err error) {
	o.Speed, err = parseFloat64p("obd speed", value, converters...)
	return err
}

func parseOBDEngineSpeed(o *OBD, value string) (err error) {
	o.EngineSpeed, err = parseFloat64p("obd engine rpm", value)
	return err
}

func parseOBDThrottle(o *OBD, value string) (err error) {
	o.Throttle, err = parseFloat64p("obd throttle rpm", value)
	return err
}

func parseOBDCoolantTemp(o *OBD, value string, converters ...converter) (err error) {
	o.CoolantTemp, err = parseFloat64p("obd coolant temp", value, converters...)
	return err
}

func parseOBDIntakeTemp(o *OBD, value string, converters ...converter) (err error) {
	o.IntakeTemp, err = parseFloat64p("obd intake temp", value, converters...)
	return err
}

func parseOBDManifoldPressure(o *OBD, value string, converters ...converter) (err error) {
	o.ManifoldPressure, err = parseFloat64p("obd manifold pressure", value, converters...)
	return err
}
