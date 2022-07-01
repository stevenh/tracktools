package geo

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMatcher(t *testing.T) {
	tests := []struct {
		name string
		lat0, lon0,
		lat1, lon1,
		lat2, lon2,
		tol float64
		expected bool
	}{
		{
			name: "not-within-1m",
			lat0: 50.858006, lon0: -0.752614,
			lat1: 50.857928, lon1: -0.752664,
			lat2: 50.857939, lon2: -0.752523,
			tol: 1,
		},
		{
			name: "within-10m",
			lat0: 50.858006, lon0: -0.752614,
			lat1: 50.857928, lon1: -0.752664,
			lat2: 50.857939, lon2: -0.752523,
			tol:      10,
			expected: true,
		},
		{
			name: "within-1mm",
			lat0: 50.857933, lon0: -0.752600,
			lat1: 50.857928, lon1: -0.752664,
			lat2: 50.857939, lon2: -0.752523,
			tol:      0.001,
			expected: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			m := NewMatcher(Tolerance(tc.tol))
			on := m.OnLine(
				tc.lat0, tc.lon0,
				tc.lat1, tc.lon1,
				tc.lat2, tc.lon2,
			)
			require.Equal(t, tc.expected, on)
		})
	}
}
