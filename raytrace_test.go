// Debuging with Delve:
// dlv test -- -test.run ^TestAABB
// b TestAABB

package raytrace

import (
	"fmt"
	"testing"
	"math"
	"sort"
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
	aabb := NewAABB(NewVec3(-1,-1, -0.0001), NewVec3(1,1, 0.0001)) // a plane infinitly small in Z
	want := true
	result := aabb.HitOptimized(r, math.Inf(-1), math.Inf(1))
	if  result != want {
		t.Errorf(" %v != %v", result, want)
	}
	r = Ray{NewVec3(0, 2, -2), NewVec3(0,0,1)}
	want = false
	result = aabb.HitOptimized(r, math.Inf(-1), math.Inf(1))
	if  result != want {
		t.Errorf(" %v != %v", result, want)
	}

	r = Ray{NewVec3(0.9999, 0.9999, -2), NewVec3(0,0,1)}
	want = true
	result = aabb.HitOptimized(r, math.Inf(-1), math.Inf(1))
	if  result != want {
		t.Errorf(" %v != %v", result, want)
	}

	r = Ray{NewVec3(1.0, 1.0, -2), NewVec3(0,0,1)}
	want = true
	result = aabb.HitOptimized(r, math.Inf(-1), math.Inf(1))
	if  result != want {
		t.Errorf(" %v != %v", result, want)
	}
	
	r = Ray{NewVec3(1.000001, 0.9999, -2), NewVec3(0,0,1)}
	want = false
	result = aabb.HitOptimized(r, math.Inf(-1), math.Inf(1))
	if  result != want {
		t.Errorf(" %v != %v", result, want)
	}
}


func TestSphereBBox(t *testing.T) {

	s := Sphere{NewVec3(0,0,-1), 0.5}

	// Min() test
	want := NewVec3(-0.5, -0.5, -1.5)
	aabb := NewAABBUninit()
	s.BBox(&aabb)
	result := aabb.Min()
	if  result != want {
		t.Errorf(" %v != %v", result, want)
	}
	// Max() test
	want = NewVec3(0.5, 0.5, -0.5)
	aabb = NewAABBUninit()
	s.BBox(&aabb)
	result = aabb.Max()
	if  result != want {
		t.Errorf(" %v != %v", result, want)
	}
}

func TestSortHittables(t *testing.T) {
	var objects []Hittable
	objects = append(objects,
		Sphere{NewVec3(0,0,-1), 1.0},  //
		Sphere{NewVec3(0,0,-3), 1.0},
		Sphere{NewVec3(0,2,-1), 1.0},  // same Z as first sphere - check if it is stable
		Sphere{NewVec3(0,0,-2), 1.0})

	// fmt.Println("Sort Hittables", objects)

	// We just pass custom comparator and dont implement interface for sorting
	sort.SliceStable(objects, func(i,j int) bool {
		return box_compare(objects[i], objects[j], 2) // along Z
	})
	// fmt.Println("Sort Hittables sorted", objects)

	want := []float32{-4,-3,-2,-2} // we expect BBox.Min.Z values sorted from smallest
	result := []float32{}
	for _, o := range objects {
		aabb := NewAABBUninit()
		o.BBox(&aabb)
		result = append(result, aabb.Min().At(2))
	}

	// compare 2 slices
	testEq := func (a, b []float32) bool {
		if len(a) != len(b) {
			return false
		}
		for i := range a {
			if a[i] != b[i] {
				return false
			}
		}
		return true
	}

	if !testEq(result, want) {
		t.Errorf(" %v != %v", result, want)
	}
	
}


func TestBVH(t *testing.T) {
	bvh := NewBVH()
	fmt.Println("Sort Hittables sorted", bvh)

	var objects []Hittable
	objects = append(objects,
		Sphere{NewVec3(0,0,-1), 1.0},
		Sphere{NewVec3(0,0,-2), 1.0})

	
}
