// Package geo implements gnomonic projection centered at an arbitrary position C on the ellipsoid.
//
// It is a [Go] port of the gnomonic routines from [GeographicLib] and supporting helpers.
//
// The projection is derived in Section 8 of
// - C. F. F. Karney, [Algorithms for geodesics], J. Geodesy 87, 43–55 (2013); DOI: [10.1007/s00190-012-0578-z]; [addenda].
//
// The projection of P is defined as follows: compute the geodesic line from C to P;
// compute the reduced length m12, geodesic scale M12, and ρ = m12/M12; finally x = ρ sin azi1;
// y = ρ cos azi1, where azi1 is the azimuth of the geodesic at C.
// The [Gnomonic.Forward] and [Gnomonic.Reverse] methods also return the azimuth azi of the geodesic
// at P and reciprocal scale rk in the azimuthal direction. The scale in the radial direction if 1/rk².
//
// For a sphere, ρ is reduces to a tan(s12/a), where s12 is the length of the geodesic from C to P,
// and the gnomonic projection has the property that all geodesics appear as straight lines. For an
// ellipsoid, this property holds only for geodesics interesting the centers. However geodesic segments
// close to the center are approximately straight.
//
// Consider a geodesic segment of length l. Let T be the point on the geodesic (extended if necessary)
// closest to C the center of the projection and t be the distance CT. To lowest order, the maximum
// deviation (as a true distance) of the corresponding gnomonic line segment (i.e., with the same end
// points) from the geodesic is the following where K is the Gaussian curvature.
//  (K(T) - K(C)) l² t / 32.
//
// This result applies for any surface. For an ellipsoid of revolution, consider all geodesics whose
// end points are within a distance r of C. For a given r, the deviation is maximum when the latitude
// of C is 45°, when endpoints are a distance r away, and when their azimuths from the center are
// ± 45° or ± 135°. To lowest order in r and the flattening f, the deviation is f (r/2a)³ r.
//
//
// [Algorithms for geodesics]: https://doi.org/10.1007/s00190-012-0578-z
// [10.1007/s00190-012-0578-z]: https://doi.org/10.1007/s00190-012-0578-z
// [addenda]: https://geographiclib.sourceforge.io/geod-addenda.html
// [Go]: https://go.dev/
// [GeographicLib]: https://geographiclib.sourceforge.io/
package geo
