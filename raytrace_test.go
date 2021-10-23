package raytrace

import (
	_ "fmt"
	"gonum.org/v1/gonum/floats"
	"testing"
)

func TestRay1(t *testing.T) {
	r := Ray{[]float64{0, 0, 0}, []float64{1, 2, 3}}
	vec := r.At(2)
	want := []float64{2, 4, 6}
	if !floats.Same(vec, want) {
		t.Errorf(" %v != %v", vec, want)
	}

}

func TestVecIndexAccesso(t *testing.T) {
	vec := NewVec3(1, 2, 3)
	want := 1.0
	if vec.At(0) != (want) {
		t.Errorf(" %v != %v", vec.At(0), want)
	}
	want = 2.0
	if vec.At(1) != (want) {
		t.Errorf(" %v != %v", vec.At(1), want)
	}
	want = 3.0
	if vec.At(2) != (want) {
		t.Errorf(" %v != %v", vec.At(2), want)
	}
}

func TestVecEqual(t *testing.T) {
	vec := NewVec3(1, 2, 3)
	want := Vec3{[3]float64{1, 2, 3}}
	if !vec.Equal(want) {
		t.Errorf(" %v != %v", vec, want)
	}
	vec = NewVec3(1, 2, 3)
	want = Vec3{[3]float64{0, 0, 0}}
	if vec.Equal(want) {
		t.Errorf(" %v != %v", vec, want)
	}
}

func TestVecAdd(t *testing.T) {
	// Add
	vec := NewVec3(1, 2, 3)
	want := Vec3{[3]float64{2, 4, 6}}
	res := vec.Add(vec)
	if !res.Equal(want) {
		t.Errorf(" %v != %v", res, want)
	}

	want = Vec3{[3]float64{2, 3, 4}}
	res = vec.AddF(1)
	if !res.Equal(want) {
		t.Errorf(" %v != %v", res, want)
	}
	// Subtract
	want = Vec3{[3]float64{0, 0, 0}}
	res = vec.Subtr(vec)
	if !res.Equal(want) {
		t.Errorf(" %v != %v", res, want)
	}

	want = Vec3{[3]float64{0, 1, 2}}
	res = vec.SubtrF(1)
	if !res.Equal(want) {
		t.Errorf(" %v != %v", res, want)
	}

	// Multiply
	want = Vec3{[3]float64{1, 4, 9}}
	res = vec.Mult(vec)
	if !res.Equal(want) {
		t.Errorf(" %v != %v", res, want)
	}

	want = Vec3{[3]float64{2, 4, 6}}
	res = vec.MultF(2)
	if !res.Equal(want) {
		t.Errorf(" %v != %v", res, want)
	}

	// Divide
	vec = Vec3{[3]float64{4, 6, 8}}
	want = Vec3{[3]float64{1, 1, 1}}
	res = vec.Div(vec)
	if !res.Equal(want) {
		t.Errorf(" %v != %v", res, want)
	}

	want = Vec3{[3]float64{2, 3, 4}}
	res = vec.DivF(2)
	if !res.Equal(want) {
		t.Errorf(" %v != %v", res, want)
	}

}

func TestVecOther(t *testing.T) {
	vec := NewVec3(0, 3, 4)
	want := 5.0
	res := vec.Length()
	if res != want {
		t.Errorf(" %v != %v", res, want)
	}

	vec = NewVec3(2, 3, 4)
	want = 29.0
	res = vec.LengthSquared()
	if res != want {
		t.Errorf(" %v != %v", res, want)
	}

	vec = Vec3{[3]float64{3, 3, 3}}
	want = Vec3{[3]float64{1, 1, 1}}
	res = vec.Div(vec)
	if !res.Equal(want) {
		t.Errorf(" %v != %v", res, want)
	}

	// Dot product

	// Cross product

}
