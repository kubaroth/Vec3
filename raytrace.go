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

func RandFloat() float32 {
	return rand.Float32()
}
func RandFloatMinMax(min, max float32) float32 {
	// Returns a random real in [min,max).
    return min + (max-min)*RandFloat()
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
func (aabb AABB) Hit(r Ray, t_min, t_max float64) bool{
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


/////

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


func(r HittableList) BBox(output_box *AABB) bool {
	// TODO output_box = &NewAABB(NewVec3(r.X0, r.X1, r.K-0.0001), NewVec3(r.Y0,Y1,r,K+0.0001))
	return true
}

// passing pointer to HittableList as we update 
func (hl *HittableList) 	Add(object Hittable) {
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
	a := r.Direction().Dot(r.Direction())
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


type XYRect struct{
	X0, X1, Y0, Y1, K float32
}

func(r XYRect) BBox(output_box *AABB) bool {
	// TODO output_box = &NewAABB(NewVec3(r.X0, r.X1, r.K-0.0001), NewVec3(r.Y0,Y1,r,K+0.0001))
	return true
}


