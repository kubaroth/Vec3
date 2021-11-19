package raytrace

import (
	"math"
	"math/rand"
)

func Deg_to_Rad(deg float64) float64 {
	return deg * math.Pi / 180.0
}

// Vec3
type Vec3 struct {
	// e [3]float32
	x,y,z float32
}

func NewVec3(x, y, z float32) Vec3 {
	return Vec3{x, y, z}
}

func (v Vec3) At(x int) float32 {
	if x == 0{
		return v.x
	} else if x == 1 {
		return v.y
	} else {
		return v.z
	}
}

func (v Vec3) Equal(v1 Vec3) bool {
	if v.x == v1.x && v.y == v1.y && v.z == v1.z {
		return true
	} else {
		return false
	}
}

func (v Vec3) Add(v1 Vec3) Vec3 {
	return NewVec3(v.At(0)+v1.At(0), v.At(1)+v1.At(1), v.At(2)+v1.At(2))
}
// Mutates
func (v *Vec3) Add_(v1 *Vec3) *Vec3 {
	v.x = v.x + v1.x
	v.y = v.y + v1.y
	v.z = v.z + v1.z
	return v
}
func (v Vec3) AddF(x float32) Vec3 {
	return NewVec3(v.x+x, v.y+x, v.z+x)
}
func (v Vec3) Subtr(v1 Vec3) Vec3 {
	return NewVec3(v.At(0)-v1.At(0), v.At(1)-v1.At(1), v.At(2)-v1.At(2))
}
func (v Vec3) SubtrF(x float32) Vec3 {
	return NewVec3(v.At(0)-x, v.At(1)-x, v.At(2)-x)
}
func (v Vec3) Mult(v1 Vec3) Vec3 {
	return NewVec3(v.At(0)*v1.At(0), v.At(1)*v1.At(1), v.At(2)*v1.At(2))
}
func (v Vec3) MultF(x float32) Vec3 {
	return NewVec3(v.At(0)*x, v.At(1)*x, v.At(2)*x)
}
func (v Vec3) Div(v1 Vec3) Vec3 {
	return NewVec3(v.At(0)/v1.At(0), v.At(1)/v1.At(1), v.At(2)/v1.At(2))
}
func (v Vec3) DivF(x float32) Vec3 {
	return NewVec3(v.At(0)/x, v.At(1)/x, v.At(2)/x)
}

func (v Vec3) Length() float32 {
	return float32(math.Sqrt(float64(v.At(0)*v.At(0) + v.At(1)*v.At(1) + v.At(2)*v.At(2))))
}

func (v Vec3) LengthSquared() float32 {
	return v.At(0)*v.At(0) + v.At(1)*v.At(1) + v.At(2)*v.At(2)
}
func (v Vec3) UnitVec() Vec3 {
	return v.DivF(v.Length())
}
func (v Vec3) Dot(v1 Vec3) float32 {
	return v.x * v1.x + v.y * v1.y + v.z * v1.z
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

func (r Ray) At(t float32) Vec3 {
	return r.Orig.Add(r.Dir.MultF(t))
}

// Not the global rand.Rand instance has a mutex and needs to be initialize outside
// the short-lived goroutine
// https://stackoverflow.com/questions/14298523/why-does-adding-concurrency-slow-down-this-golang-code
func RandFloat() float32 {
	return rand.Float32()
}
func RandFloatMinMax(min, max float32) float32 {
	// Returns a random real in [min,max).
    return min + (max-min)*RandFloat()
}

func RandInt(min, max float32) int {
	// Returns a random int in [min,max].
    return int(RandFloatMinMax(min, max+1))
}

func Clamp(x, min, max float32 ) float32{
	if x < min {
		return min
	}
	if x > max {
		return max
	}
	return x

}
