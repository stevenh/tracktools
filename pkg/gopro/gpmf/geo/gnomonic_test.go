package geo

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tidwall/geodesic"
)

func TestIntersect(t *testing.T) {
	tests := []struct {
		name string
		lata1, lona1,
		lata2, lona2,
		latb1, lonb1,
		latb2, lonb2,
		expected_lat, expected_lon,
		expected_azi1a, expected_azi2a,
		expected_azi1b, expected_azi2b float64
	}{
		{
			name:  "initial",
			lata1: 42, lona1: 29,
			lata2: 39, lona2: -77,
			latb1: 64, lonb1: -22,
			latb2: 6, lonb2: 0,
			expected_lat: 54.717030, expected_lon: -14.563856,
			expected_azi1a: -84.145064, expected_azi2a: -84.145064,
			expected_azi1b: 160.917494, expected_azi2b: 160.917494,
		},
	}

	gn := NewGnomonic(geodesic.WGS84)
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			lat, lon, azi1a, azi2a, azi1b, azi2b := gn.Intersect(
				tc.lata1, tc.lona1,
				tc.lata2, tc.lona2,
				tc.latb1, tc.lonb1,
				tc.latb2, tc.lonb2,
			)
			t.Logf("Result %f, %f\n", lat, lon)
			floatEqual(t, tc.expected_lat, lat)
			floatEqual(t, tc.expected_lon, lon)

			t.Logf("Azimuths on line A %f, %f\n", azi2a, azi1a)
			floatEqual(t, tc.expected_azi1a, azi1a)
			floatEqual(t, tc.expected_azi2a, azi2a)

			t.Logf("Azimuths on line B %f, %f\n", azi2b, azi1b)
			floatEqual(t, tc.expected_azi1b, azi1b)
			floatEqual(t, tc.expected_azi2b, azi2b)
		})
	}
}

func TestIntersect2(t *testing.T) {
	g := geodesic.WGS84
	gn := NewGnomonic(g)
	lat, lon, azi1a, azi2a, azi1b, azi2b := gn.Intersect(
		50.857907, -0.752633,
		50.857911, -0.752567,
		50.857941, -0.752600,
		50.857930, -0.752601,
		//50.857833, -0.752584,
	)
	t.Logf("Result %f, %f\n", lat, lon)
	t.Logf("Azimuths on line A %f, %f\n", azi2a, azi1a)
	t.Logf("Azimuths on line B %f, %f\n", azi2b, azi1b)

	var lat2, lon2, azi2 float64
	g.Direct(
		50.857907, -0.752633, 83, 10,
		&lat2, &lon2, &azi2,
	)
	t.Logf("Result %f, %f azi: %f\n", lat2, lon2, azi2)
}

// floatEqual checks two floats to 6dp.
func floatEqual(t *testing.T, expected, val float64) {
	t.Helper()
	e := fmt.Sprintf("%.6f", expected)
	v := fmt.Sprintf("%.6f", val)
	require.Equal(t, e, v)
}
