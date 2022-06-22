package trackaddict

const (
	// To metric.
	m2km    = 1.60934
	ft2m    = 0.3048
	psi2kpa = 6.89476

	// To imperial.
	km2m    = 0.621371
	m2ft    = 0.3048
	kpa2psi = 0.145038
)

var (
	// Imperialunits configures all units to imperial.
	Imperialunits = units{
		Altitude:    meters2Feet,
		Accuracy:    meters2Feet,
		Speed:       kilometers2Miles,
		Pressure:    kpa2Psi,
		Temperature: celsius2Fahrenheit,
	}

	// Ukunits configures all units to metric except
	// Speed which is configured to imperial.
	Ukunits = units{
		Speed: kilometers2Miles,
	}
)

// converter represents a function which can convert from one unit to another.
type converter func(float64) float64

// units represents converters for known measurement types.
type units struct {
	Altitude    converter
	Accuracy    converter
	Speed       converter
	Pressure    converter
	Temperature converter
}

// celsius2Fahrenheit converts a Celsius temperature to Fahrenheit.
func celsius2Fahrenheit(c float64) float64 {
	return (c * 9 / 5) + 32
}

// fahrenheit2Celsius converts a Fahrenheit temperature to Celsius.
func fahrenheit2Celsius(f float64) float64 {
	return (f - 32) * 5 / 9
}

// feet2Meters converts a measurement in feet to meters.
func feet2Meters(v float64) float64 {
	return v * ft2m
}

// meters2Feet converts a measurement in meters to feet.
func meters2Feet(v float64) float64 {
	return v * m2ft
}

// kilometers2Miles converts a measurement in kilometers to miles.
func kilometers2Miles(v float64) float64 {
	return v * km2m
}

// miles2Kilometers converts a measurement in miles to kilometers.
func miles2Kilometers(v float64) float64 {
	return v * m2km
}

// kpa2Psi converts a pressure in KPA to PSI.
func kpa2Psi(v float64) float64 {
	return v * kpa2psi
}

// psi2Kpa converts a pressure in PSI to KPA.
func psi2Kpa(v float64) float64 {
	return v * psi2kpa
}
