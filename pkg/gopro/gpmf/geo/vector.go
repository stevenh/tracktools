package geo

// vector represents a 3D vector.
type vector struct {
	x, y, z float64
}

// newVector2 creates a new vector with x, y and z set to 1.
func newVector2(x, y float64) *vector {
	return &vector{
		x: x,
		y: y,
		z: 1,
	}
}

// cross returns the cross product of the vector and b.
func (v *vector) cross(other *vector) *vector {
	return &vector{
		x: v.y*other.z - v.z*other.y,
		y: v.z*other.x - v.x*other.z,
		z: v.x*other.y - v.y*other.x,
	}
}

// norm normalises the vectore.
func (v *vector) norm() {
	v.x /= v.z
	v.y /= v.z
	v.z = 1
}
