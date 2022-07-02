package geo

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tidwall/geodesic"
)

func TestProcessorOnLine(t *testing.T) {
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
			m := NewProcessor(Tolerance(tc.tol))
			on := m.OnLine(
				tc.lat0, tc.lon0,
				tc.lat1, tc.lon1,
				tc.lat2, tc.lon2,
			)
			require.Equal(t, tc.expected, on)
		})
	}
}

func TestProcessorDistanceToLine(t *testing.T) {
	tests := []struct {
		name string
		pLat, pLon,
		sLat, sLon,
		eLat, eLon,
		expected float64
	}{
		/*
			{
				name: "within-10m",
				pLat: 50.858006, pLon: -0.752614,
				sLat: 50.857928, sLon: -0.752664,
				eLat: 50.857939, eLon: -0.752523,
				expected: 6.515,
			},
		*/
		{
			name: "within-5cm",
			pLat: 50.857933, pLon: -0.752594,
			sLat: 50.857928, sLon: -0.752664,
			eLat: 50.857939, eLon: -0.752523,
			expected: 0.051,
		},
	}

	m := NewProcessor()
	radius := geodesic.WGS84.Radius()
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			e := distanceEquirect(tc.pLat, tc.pLon, tc.sLat, tc.sLon, radius)
			fmt.Printf("est1: %.9f\n", e)

			e = distanceHaversin(tc.pLat*radians, tc.pLon*radians, tc.sLat*radians, tc.sLon*radians, radius)
			fmt.Printf("est2: %.9f\n", e)

			d := m.DistanceToLine(
				tc.pLat, tc.pLon,
				tc.sLat, tc.sLon,
				tc.eLat, tc.eLon,
			)
			floatEqual(t, tc.expected, d, 3)
		})
	}
}

var result float64

func BenchmarkDistance(b *testing.B) {
	tests := []struct {
		name string
		f    func(lat1, lon1, lat2, lon2, radius float64) float64
	}{
		{
			name: "distance-equirectangular",
			f:    distanceEquirect,
		},
		{
			name: "distance-haversine",
			f: func(lat1, lon1, lat2, lon2, radius float64) float64 {
				return distanceHaversin(lat1*radians, lon1*radians, lat2*radians, lon2*radians, radius)
			},
		},
	}

	var r float64
	radius := geodesic.WGS84.Radius()
	for _, tc := range tests {
		b.Run(tc.name, func(b *testing.B) {
			for n := 0; n < b.N; n++ {
				r = tc.f(
					50.857933, -0.752594,
					50.857928, -0.752664,
					radius,
				)
			}
			result = r
			floatEqual(b, 4.950285274, r, 9)
		})
	}
}
