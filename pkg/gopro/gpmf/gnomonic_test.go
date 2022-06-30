package gpmf

import (
	"fmt"
	"testing"

	"github.com/tidwall/geodesic"
)

func TestCrossing(t *testing.T) {
	tests := []struct {
		name string
		lata1, lona1,
		lata2, lona2,
		latb1, lonb1,
		latb2, lonb2 float64
	}{
		{
			name:  "initial",
			lata1: 42,
			lona1: 29,
			lata2: 39,
			lona2: -77,
			latb1: 64,
			lonb1: -22,
			latb2: 6,
			lonb2: 0,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			geod := geodesic.WGS84
			gn := NewGnomonic(geod)

			lat0 := (tc.lata1 + tc.lata2 + tc.latb1 + tc.latb2) / 4
			// Possibly need to deal with longitudes wrapping around
			lon0 := (tc.lona1 + tc.lona2 + tc.lonb1 + tc.lonb2) / 4
			fmt.Printf("Initial guess %.050f %.050f\n", lat0, lon0)
			for i := 0; i < 10; i++ {
				xa1, ya1, _, _ := gn.Forward(lat0, lon0, tc.lata1, tc.lona1)
				xa2, ya2, _, _ := gn.Forward(lat0, lon0, tc.lata2, tc.lona2)
				xb1, yb1, _, _ := gn.Forward(lat0, lon0, tc.latb1, tc.lonb1)
				xb2, yb2, _, _ := gn.Forward(lat0, lon0, tc.latb2, tc.lonb2)
				// See Hartley and Zisserman, Multiple View Geometry, Sec. 2.2.1
				va1 := newVector3(xa1, ya1)
				va2 := newVector3(xa2, ya2)
				vb1 := newVector3(xb1, yb1)
				vb2 := newVector3(xb2, yb2)
				// la is homogeneous representation of line A1,A2
				// lb is homogeneous representation of line B1,B2
				la := va1.cross(va2)
				lb := vb1.cross(vb2)
				// p0 is homogeneous representation of intersection of la and lb
				p0 := la.cross(lb)
				p0.norm()
				lat1, lon1, _, _ := gn.Reverse(lat0, lon0, p0.x, p0.y)
				fmt.Printf("Increment %.050f %.050f\n", lat1-lat0, lon1-lon0)
				lat0 = lat1
				lon0 = lon1
			}
			fmt.Printf("Final result %.050f %.050f\n", lat0, lon0)
			var azi1, azi2 float64
			geod.Inverse(tc.lata1, tc.lona1, lat0, lon0, nil, nil, &azi2)
			geod.Inverse(lat0, lon0, tc.lata2, tc.lona2, nil, &azi1, nil)
			fmt.Printf("Azimuths on line A %.050f %.050f\n", azi2, azi1)
			geod.Inverse(tc.latb1, tc.lonb1, lat0, lon0, nil, nil, &azi2)
			geod.Inverse(lat0, lon0, tc.latb2, tc.lonb2, nil, &azi1, nil)
			fmt.Printf("Azimuths on line B %.050f %.050f\n", azi2, azi1)
		})
	}
}
