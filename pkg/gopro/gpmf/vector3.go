package gpmf

type vector3 struct {
	x, y, z float64
}

func newVector3(x, y float64) *vector3 {
	return &vector3{
		x: x,
		y: y,
		z: 1,
	}
}

func (v *vector3) cross(b *vector3) *vector3 {
	return &vector3{
		x: v.y*b.z - v.z*b.y,
		y: v.z*b.x - v.x*b.z,
		z: v.x*b.y - v.y*b.x,
	}
}

func (v *vector3) norm() {
	v.x /= v.z
	v.y /= v.z
	v.z = 1
}
