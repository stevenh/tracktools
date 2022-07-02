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
		expectedLat, expectedLon,
		expectedAzi1a, expectedAzi2a,
		expectedAzi1b, expectedAzi2b float64
	}{
		{
			name:  "initial",
			lata1: 42, lona1: 29,
			lata2: 39, lona2: -77,
			latb1: 64, lonb1: -22,
			latb2: 6, lonb2: 0,
			expectedLat: 54.717030, expectedLon: -14.563856,
			expectedAzi1a: -84.145064, expectedAzi2a: -84.145064,
			expectedAzi1b: 160.917494, expectedAzi2b: 160.917494,
		},
	}

	gn := NewGnomonic(geodesic.WGS84)
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			lat, lon, azi1a, azi2a, azi1b, azi2b := gn.IntersectExt(
				tc.lata1, tc.lona1,
				tc.lata2, tc.lona2,
				tc.latb1, tc.lonb1,
				tc.latb2, tc.lonb2,
			)
			t.Logf("Result %f, %f\n", lat, lon)
			floatEqual(t, tc.expectedLat, lat, 6)
			floatEqual(t, tc.expectedLon, lon, 6)

			t.Logf("Azimuths on line A %f, %f\n", azi2a, azi1a)
			floatEqual(t, tc.expectedAzi1a, azi1a, 6)
			floatEqual(t, tc.expectedAzi2a, azi2a, 6)

			t.Logf("Azimuths on line B %f, %f\n", azi2b, azi1b)
			floatEqual(t, tc.expectedAzi1b, azi1b, 6)
			floatEqual(t, tc.expectedAzi2b, azi2b, 6)
		})
	}
}

func TestLines(t *testing.T) {
	g := geodesic.WGS84
	gn := NewGnomonic(g)
	lat, lon, azi1a, azi2a, azi1b, azi2b := gn.IntersectExt(
		50.857928, -0.752664,
		50.857939, -0.752523,
		50.858006, -0.752614,
		50.857828, -0.752579,
	)
	t.Logf("Result %f, %f\n", lat, lon)
	t.Logf("Azimuths on line A %f, %f\n", azi2a, azi1a)
	t.Logf("Azimuths on line B %f, %f\n", azi2b, azi1b)

	var lat2, lon2, azi2 float64
	for _, d := range []float64{5, 10} {
		g.Direct(
			50.857928, -0.752664, 83, d,
			&lat2, &lon2, &azi2,
		)
		t.Logf("Result %f, %f azi: %f\n", lat2, lon2, azi2)
	}

	g.Direct(
		50.858006, -0.752614, 173, 20,
		&lat2, &lon2, &azi2,
	)
	t.Logf("Result %f, %f azi: %f\n", lat2, lon2, azi2)
}

type testingT interface {
	require.TestingT
	Helper()
}

// floatEqual checks two floats to a given db.
func floatEqual(t testingT, expected, val float64, dp int) {
	t.Helper()
	format := fmt.Sprintf("%%.%df", dp)
	e := fmt.Sprintf(format, expected)
	v := fmt.Sprintf(format, val)
	require.Equal(t, e, v)
}
