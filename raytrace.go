package raytrace

import (
	"math"
)

// Vec3
type Vec3 struct {
	e [3]float64
}

func NewVec3(x, y, z float64) Vec3 {
	return Vec3{[3]float64{x, y, z}}
}

func (v Vec3) At(x int) float64 {
	return v.e[x]
}

func (v Vec3) Equal(v1 Vec3) bool {
	if v.At(0) == v1.At(0) && v.At(1) == v1.At(1) && v.At(2) == v1.At(2) {
		return true
	} else {
		return false
	}
}

func (v Vec3) Add(v1 Vec3) Vec3 {
	//return NewVec3(v.At(0)+v1.At(0), v.At(1)+v1.At(1), v.At(2)+v1.At(2))
	// optimization
	v.e[0] = v.e[0]+v1.e[0]
	v.e[1] = v.e[1]+v1.e[1]
	v.e[2] = v.e[2]+v1.e[2]
	return v
}
func (v Vec3) AddF(x float64) Vec3 {
	return NewVec3(v.e[0]+x, v.e[1]+x, v.e[2]+x)
}
func (v Vec3) Subtr(v1 Vec3) Vec3 {
	return NewVec3(v.At(0)-v1.At(0), v.At(1)-v1.At(1), v.At(2)-v1.At(2))
}
func (v Vec3) SubtrF(x float64) Vec3 {
	return NewVec3(v.At(0)-x, v.At(1)-x, v.At(2)-x)
}
func (v Vec3) Mult(v1 Vec3) Vec3 {
	return NewVec3(v.At(0)*v1.At(0), v.At(1)*v1.At(1), v.At(2)*v1.At(2))
}
func (v Vec3) MultF(x float64) Vec3 {
	return NewVec3(v.At(0)*x, v.At(1)*x, v.At(2)*x)
}
func (v Vec3) Div(v1 Vec3) Vec3 {
	return NewVec3(v.At(0)/v1.At(0), v.At(1)/v1.At(1), v.At(2)/v1.At(2))
}
func (v Vec3) DivF(x float64) Vec3 {
	return NewVec3(v.At(0)/x, v.At(1)/x, v.At(2)/x)
}

func (v Vec3) Length() float64 {
	return math.Sqrt(v.At(0)*v.At(0) + v.At(1)*v.At(1) + v.At(2)*v.At(2))
}

func (v Vec3) LengthSquared() float64 {
	return v.At(0)*v.At(0) + v.At(1)*v.At(1) + v.At(2)*v.At(2)
}
func (v Vec3) UnitVec() Vec3 {
	return v.DivF(v.Length())
}
func (v Vec3) Dot(v1 Vec3) float64 {
	return v.e[0]*v1.e[0] + v.e[1]*v1.e[1] + v.e[2]*v1.e[2]
}
func (v Vec3) Cross(u Vec3) Vec3 {
	return NewVec3(
		u.At(1)*v.At(2)-u.At(2)*v.At(1),
		u.At(2)*v.At(0)-u.At(0)*v.At(2),
		u.At(0)*v.At(1)-u.At(1)*v.At(0))
}

// Ray
type Ray struct {
	Orig Vec3
	Dir  Vec3
}

func NewRay(origin, dir Vec3) Ray {
	return Ray{origin, dir}
}

func (r Ray) Origin() Vec3 {
	return r.Orig
}
func (r Ray) Direction() Vec3 {
	return r.Dir
}

func (r Ray) At(t float64) Vec3 {
	return r.Orig.Add(r.Dir.MultF(t))
}
