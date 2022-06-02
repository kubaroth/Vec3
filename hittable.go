package raytrace

import (
	"fmt"
	"math"
	"sort"
)
type AABB struct {
	min, max Vec3
}

func NewAABB(min, max Vec3) AABB{
	return AABB{min,max}
}

func NewAABBUninit() AABB{
	return NewAABB(NewVec3(-1,-1, -0.0001), NewVec3(1,1, 0.0001))
}

func(aabb AABB) Min() Vec3 {
	return aabb.min
}

func(aabb AABB) Max() Vec3 {
	return aabb.max
}
func (aabb AABB) Hit(r *Ray, t_min, t_max float64) bool{
	for a:=0; a < 3; a++{
		var t0, t1 float64
		ray_dir := float64(r.Direction().At(a))

		// Handling division by 0
		t0_min := float64(aabb.min.At(a) - r.Origin().At(a))
		t0_max := float64(aabb.max.At(a) - r.Origin().At(a))
		if ray_dir == 0.0 {
			if t0_min < 0 {
				t0_min = math.Inf(-1)
			} else {
				t0_min = math.Inf(1)
			}
			if t0_max < 0 {
				t0_max = math.Inf(-1)
			} else {
				t0_max = math.Inf(1)
			}
			t0 = math.Min( t0_min, t0_max)
			t1 = math.Max( t0_min, t0_max)
		} else {
			t0 = math.Min( t0_min / ray_dir, t0_max / ray_dir)
			t1 = math.Max( t0_min / ray_dir, t0_max / ray_dir)
		}

		t_min = math.Max(t0, t_min)
		t_max = math.Min(t1, t_max)
		if (t_max <= t_min){
			return false
		}
	}
	return true
}

func (aabb AABB) HitOptimized(r Ray, t_min, t_max float64) bool{
	for a:=0; a < 3; a++{
		var t0, t1, invD float64
		ray_dir := float64(r.Direction().At(a))

		// Handling division by 0
		t0_min := float64(aabb.min.At(a) - r.Origin().At(a))
		t0_max := float64(aabb.max.At(a) - r.Origin().At(a))
		if ray_dir == 0.0 {
			if t0_min < 0 {
				t0_min = math.Inf(-1)
			} else {
				t0_min = math.Inf(1)
			}
			if t0_max < 0 {
				t0_max = math.Inf(-1)
			} else {
				t0_max = math.Inf(1)
			}
			t0 = math.Min( t0_min, t0_max)
			t1 = math.Max( t0_min, t0_max)
		} else {
			invD = 1.0 / ray_dir
			t0 = math.Min( t0_min / ray_dir, t0_max / ray_dir) * invD
			t1 = math.Max( t0_min / ray_dir, t0_max / ray_dir) * invD
		}
		if invD < 0.0 {
			t0,t1 = t1,t0
		}

		if t0 > t_min {
			t_min = t0
		}
		if t1 < t_max {
			t_max = t1
		}
		if (t_max <= t_min){
			return false
		}
	}
	return true
}

type Hittable interface {
	Hit(r *Ray, t_min, t_max float32, rec *HitRecord) bool
	BBox(aabb *AABB) bool
}

func (rec *HitRecord) set_face_normal(r *Ray, outward_normal * Vec3) {
	if r.Direction().Dot(*outward_normal) < 0 {
		rec.Normal = *outward_normal
		rec.FrontFace = true
	} else {
		rec.Normal = NewVec3(-outward_normal.At(0), -outward_normal.At(1), -outward_normal.At(2))
		rec.FrontFace = false
	}
}


type HittableList struct {
	Objects []Hittable
}

func (hl HittableList) 	Hit(r *Ray, t_min, t_max float32, rec *HitRecord) bool {
	temp_rec := HitRecord{NewVec3(0,0,0), NewVec3(0,0,0), 1.0, true, -1}
	hit_anything:= false;
	closest_so_far := float32(math.Inf(1.0))

	for obj_id, object := range hl.Objects {
		if object.Hit(r, 0.0, float32(closest_so_far), &temp_rec) {
			closest_so_far = temp_rec.T
			hit_anything = true
			temp_rec.ObjectId = obj_id
			*rec = temp_rec
		}
	}

	return hit_anything
}


func(hl HittableList) BBox(output_box *AABB) bool {
	// TODO output_box = &NewAABB(NewVec3(r.X0, r.X1, r.K-0.0001), NewVec3(r.Y0,Y1,r,K+0.0001))

	if len(hl.Objects) == 0{
		return false
	}

	temp_box := NewAABBUninit()
	first_box := true
	for _, object := range hl.Objects {
		if !object.BBox(&temp_box) {
			return false;
		}
		if first_box {
			*output_box = temp_box
		} else {
			*output_box = Surrounding_box(*output_box, temp_box)
		}

	}
	
	return true
}

func Surrounding_box(box0, box1 AABB) AABB {

	x := float32(math.Min(float64(box0.Min().At(0)), float64(box1.Min().At(0))))
	y := float32(math.Min(float64(box0.Min().At(1)), float64(box1.Min().At(1))))
	z := float32(math.Min(float64(box0.Min().At(2)), float64(box1.Min().At(2))))
	small := NewVec3(x,y,z)
	x = float32(math.Max(float64(box0.Max().At(0)), float64(box1.Max().At(0))))
	y = float32(math.Max(float64(box0.Max().At(1)), float64(box1.Max().At(1))))
	z = float32(math.Max(float64(box0.Max().At(2)), float64(box1.Max().At(2))))
	big := NewVec3(x,y,z)

	return NewAABB(small, big)
}

type BVH_node struct {
	Left, Right Hittable // I dont think we can have Interface Hittable as Left, Right
	Box AABB
}

// Returns *BVH_node which, if required, can be casted to Hittable 
func NewBVH() *BVH_node{
	temp := BVH_node{nil, nil, NewAABBUninit()}
	return &temp
}


func NewBVHSplit(objects []Hittable, start, end int) *BVH_node{

	// randomly choose an axis
	// sort the primitives (using std::sort)
	// put half in each subtree

	bvh := NewBVH()
	
	axis := RandInt(0,2)
	fmt.Println("BVH axis", axis)
	

	comparator := box_x_compare
	if axis == 1 {
		comparator = box_y_compare
	} else {
		comparator = box_z_compare
	}
	_ = comparator
	object_span := end - start

	if object_span == 1 {
		bvh.Left = objects[start]
		bvh.Right = objects[start] // the same
	} else if object_span == 2 {
		if comparator(objects[start], objects[start+1]) {
			bvh.Left = objects[start]   // * handles pointer to interface errors
			bvh.Right = objects[start+1]
		} else {
			bvh.Left = objects[start+1]
			bvh.Right = objects[start]
		}
	} else {
		sort.SliceStable(objects, func(i,j int) bool {
			return box_compare(objects[i], objects[j], axis)
		})

		mid := start + object_span/2
		bvh.Left = NewBVHSplit(objects, start, mid)
		bvh.Right = NewBVHSplit(objects, mid, end)
	}

	box_left := NewAABBUninit()
	box_right := NewAABBUninit()
	if !(bvh.Left).BBox(&box_left) || !(bvh.Right).BBox(&box_right){
		fmt.Printf("No boudning box in Bvh_node constructor \n")
	}

	bvh.Box = Surrounding_box(box_left, box_right)
	
	return bvh
	
}

func box_compare(a, b Hittable, axis int) bool {
	box_a := NewAABBUninit()
	box_b := NewAABBUninit()
	if !a.BBox(&box_a) || !b.BBox(&box_b){
		fmt.Printf("No boudning box in Bvh_node constructor \n")
	}
	return box_a.Min().At(axis) < box_b.Min().At(axis)
}

func box_x_compare (a, b Hittable) bool {
    return box_compare(a, b, 0);
}

func box_y_compare (a, b Hittable) bool {
    return box_compare(a, b, 1);
}

func box_z_compare (a, b Hittable) bool {
    return box_compare(a, b, 2);
}
                     
func (bvh BVH_node) Hit(r *Ray, t_min, t_max float32, rec *HitRecord) bool {
	if !bvh.Box.Hit(r, float64(t_min), float64(t_max)){
		return false
	}
	// this is still unclear to me, why do we need to dereference the Left or otherwise we get an error:
	// bvh.Left.Hit undefined (type *Hittable is pointer to interface, not interface)
	//
	// I think this is because bvh.Left is still a pointer 
	hit_left := (bvh.Left).Hit(r, t_min, t_max, rec) 

	var tt float32
	if hit_left{
		tt = rec.T
	} else{
		tt = t_max
	}
	hit_right := (bvh.Right).Hit(r, t_min, tt, rec)
	return hit_left || hit_right
}

func (bvh BVH_node) BBox(output_box *AABB) bool{
	*output_box = bvh.Box
	return true
}


// passing pointer to HittableList as we update 
func (hl *HittableList) Add(object Hittable) {
	hl.Objects = append(hl.Objects, object)
	println("world length", len(hl.Objects))
}


type HitRecord struct {
	P, Normal Vec3 // point and normal
	T float32
	FrontFace bool
	ObjectId int // default -1 : helper to determine which object was hit by a ray
}

type Sphere struct {
	Center Vec3
	Radius float32
}


// Equation of sphere in vector form
// (P - C) dot (P - C) = r**2
// (A + tb - C) dot (A + tb - C) = r**2
// A + tb is a Ray
// quadratic equation x**2 + x + 1 =0 where t is unkown
// t**2 b dot b + 2t b dot (A-C) + (A-C) dot(A-C) - r**2 = 0
//      --a---       -----b-----   ------c--------
func (s Sphere) Hit(r *Ray, t_min, t_max float32, rec *HitRecord) bool {
	oc := r.Origin().Subtr(s.Center)
	a := r.Direction().Dot(r.Direction())  // square of length of the vector
	half_b := oc.Dot(r.Direction())
	c := oc.Dot(oc) - (s.Radius * s.Radius)
	discriminant := float64(half_b * half_b - a*c) // finding roots
	if discriminant < 0{
		return false
	}
    // Find the nearest root that lies in the acceptable range.
	sqrtd := float32(math.Sqrt(discriminant))
	root := (-half_b - sqrtd) / a
	if (root < t_min || t_max < root) {
        root = (-half_b + sqrtd) / a
        if (root < t_min || t_max < root) {
            return false
		}
    }
    rec.T = root;
    rec.P = r.At(rec.T); // hit point at sphere
    outward_normal := (rec.P.Subtr(s.Center)).DivF(s.Radius)
    rec.set_face_normal(r, &outward_normal)
	return true;
}

func (s Sphere) BBox(out_aabb *AABB) bool  {
	aabb := NewAABB(s.Center.Subtr(NewVec3(s.Radius, s.Radius, s.Radius)),
		s.Center.Add(NewVec3(s.Radius, s.Radius, s.Radius)))
	*out_aabb = aabb
	return true
}


type Cylinder struct {
	Center Vec3
	Radius float32
	Height float32
}

func (cyl Cylinder) Hit(r *Ray, t_min, t_max float32, rec *HitRecord) bool {

	cylinder_pos := cyl.Center // NewVec3(0,0,-1)
	cylinder_rot := NewVec3(0,1,0)
	ray := Ray{Orig: r.Orig}
	ray.Dir = r.Direction().Cross(cylinder_rot)
	co := ray.Orig.Subtr(cylinder_pos)
	a := ray.Dir.Dot(ray.Dir)
	b := 2 * ray.Dir.Dot( co.Cross(cylinder_rot) )
	c := co.Cross(cylinder_rot).Dot( co.Cross(cylinder_rot) ) - cyl.Radius * cyl.Radius
	d := b * b - 4 * c * a
	if d < 0{
		return false
	}

	t1 := ( float64(-b) - math.Sqrt(float64(d)) /  (2*float64(a))) 
	t2 := ( float64(-b) + math.Sqrt(float64(d)) /  (2*float64(a))) 

	if t2 < 0{
		return false
	}

	var t float64
	if t1 > 0 {
		t = t1
	} else {
		t = t2
	}

	// at this point we have infinite height - need to cut at hight
	v := r.Origin().At(1) + float32(t) * r.Direction().At(1)
	if ((v > cylinder_pos.At(1) - cyl.Height/2.0) && v <= cylinder_pos.At(1) + cyl.Height/2.0){
		return true
	} else {
		return false
	}
	
}

func (c Cylinder) BBox(out_aabb *AABB) bool {
	return true // TODO: for now always return true
}

type XYRect struct{
	X0, X1, Y0, Y1, K float32
}

func(r XYRect) BBox(output_box *AABB) bool {
	// TODO output_box = &NewAABB(NewVec3(r.X0, r.X1, r.K-0.0001), NewVec3(r.Y0,Y1,r,K+0.0001))
	return true
}


