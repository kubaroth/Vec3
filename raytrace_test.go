// Debuging with Delve:
// dlv test -- -test.run ^TestAABB
// b TestAABB

package raytrace

import (
	_ "fmt"
	"testing"
)

func TestRay1(t *testing.T) {
	r := Ray{NewVec3(0, 0, 0), NewVec3(1, 2, 3)}
	want := float32(2.0)
	if r.Dir.At(1) != want {
		t.Errorf(" %v != %v", r.Dir.At(2), want)
	}

}

func TestVecIndexAccesso(t *testing.T) {
	vec := NewVec3(1, 2, 3)
	want := float32(1.0)
	if vec.At(0) != (want) {
		t.Errorf(" %v != %v", vec.At(0), want)
	}
	want = float32(2.0)
	if vec.At(1) != (want) {
		t.Errorf(" %v != %v", vec.At(1), want)
	}
	want = float32(3.0)
	if vec.At(2) != (want) {
		t.Errorf(" %v != %v", vec.At(2), want)
	}
}

func TestVecEqual(t *testing.T) {
	vec := NewVec3(1, 2, 3)
	want := NewVec3(1, 2, 3)
	if !vec.Equal(want) {
		t.Errorf(" %v != %v", vec, want)
	}
	vec = NewVec3(1, 2, 3)
	want = NewVec3(0, 0, 0)
	if vec.Equal(want) {
		t.Errorf(" %v != %v", vec, want)
	}
}

func TestVecAdd(t *testing.T) {
	// Add
	vec := NewVec3(1, 2, 3)
	want := NewVec3(2, 4, 6)
	res := vec.Add(vec)
	if !res.Equal(want) {
		t.Errorf(" %v != %v", res, want)
	}

	want = NewVec3(2, 3, 4)
	res = vec.AddF(1)
	if !res.Equal(want) {
		t.Errorf(" %v != %v", res, want)
	}
	// Subtract
	want = NewVec3(0, 0, 0)
	res = vec.Subtr(vec)
	if !res.Equal(want) {
		t.Errorf(" %v != %v", res, want)
	}

	want = NewVec3(0, 1, 2)
	res = vec.SubtrF(1)
	if !res.Equal(want) {
		t.Errorf(" %v != %v", res, want)
	}

	// Multiply
	want = NewVec3(1, 4, 9)
	res = vec.Mult(vec)
	if !res.Equal(want) {
		t.Errorf(" %v != %v", res, want)
	}

	want = NewVec3(2, 4, 6)
	res = vec.MultF(2)
	if !res.Equal(want) {
		t.Errorf(" %v != %v", res, want)
	}

	// Divide
	vec = NewVec3(4, 6, 8)
	want = NewVec3(1, 1, 1)
	res = vec.Div(vec)
	if !res.Equal(want) {
		t.Errorf(" %v != %v", res, want)
	}

	want = NewVec3(2, 3, 4)
	res = vec.DivF(float32(2))
	if !res.Equal(want) {
		t.Errorf(" %v != %v", res, want)
	}

}

func TestVecOther(t *testing.T) {
	vec := NewVec3(0, 3, 4)
	want := float32(5.0)
	res := vec.Length()
	if res != want {
		t.Errorf(" %v != %v", res, want)
	}

	vec = NewVec3(2, 3, 4)
	want = float32(29.0)
	res = vec.LengthSquared()
	if res != want {
		t.Errorf(" %v != %v", res, want)
	}

	// Dot product
	vec1 := NewVec3(2, 3, 4)
	vec2 := NewVec3(2, 3, 4)
	if vec1.Dot(vec2) != vec2.LengthSquared(){
		t.Errorf(" %v != %v", vec1, vec2)
	}
	
	// Cross product

}


func TestAABB(t *testing.T) {
	r := Ray{NewVec3(0, 0, -2), NewVec3(0,0,1)}
	aabb := NewAABB(NewVec3(-1,-1,0), NewVec3(1,1,0))
	want := true
	result := aabb.Hit(r,1000, 0)
	if  result != want {
		t.Errorf(" %v != %v", result, want)
	}
	// chexk 0,2 still hits?
	r = Ray{NewVec3(0, 2, -2), NewVec3(0,0,1)}
	want = false
	result = aabb.Hit(r, 1000, 0)
	if  result != want {
		t.Errorf(" %v != %v", result, want)
	}

}
