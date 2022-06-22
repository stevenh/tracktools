package trackaddict

// Acceleration represents acceleration in all axis.
type Acceleration struct {
	// X is the acceleration in X.
	X float64

	// Y is the acceleration in Y.
	Y float64

	// Z is the acceleration in Z.
	Z float64
}

func parseAccelX(r *Record, value string) (err error) {
	r.InitAccel().X, err = parseFloat64("accel x", value)
	return err
}

func parseAccelY(r *Record, value string) (err error) {
	r.InitAccel().Y, err = parseFloat64("accel y", value)
	return err
}

func parseAccelZ(r *Record, value string) (err error) {
	r.InitAccel().Z, err = parseFloat64("accel z", value)
	return err
}
