// A driver test program
// to debu: go build -gcflags="all=-N -l" main.go
//
// Profiling:
// $ go build
// $ ./_examples -cpuprofile aaa
// $ go tool pprof aaa
// (pprof) top 10

package main

import (
	. "github.com/kubaroth/Vec3"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math"
	"os"
	"time"
	"flag"
	"runtime/pprof"
)

type Camera struct{
	Height, Width int
	Origin Vec3
	Lower_left_corner Vec3
	Horizontal Vec3
	Vertical Vec3
}

func NewCamera(lookfrom, lookat Vec3, width int) Camera {

	//Image
	vfov := 90.0
	aspect_ratio := 16.0 / 9.0
	height := int(float64(width) / aspect_ratio)
	fmt.Printf("image res", width, height)
	
	// Camera
	cam := Camera{}
	cam.Width = width
	cam.Height = height
	theta := float64(Deg_to_Rad(vfov))
	h := math.Tan(theta/2)
	_ = h
	viewport_height := 2.0 *h
	viewport_width := aspect_ratio * viewport_height

	
	vup := NewVec3(0,1,0)
	w := (lookfrom.Subtr(lookat)).UnitVec()
	u := (vup.Cross(w)).UnitVec()
	v := w.Cross(u)
	_ = v
	// fmt.Println("w,u,v", w, u, v, viewport_width, viewport_height)	
	cam.Origin = lookfrom

	cam.Horizontal = NewVec3(float32(viewport_width), 0, 0)
	cam.Vertical = NewVec3(0, float32(viewport_height), 0)

	// Disable scalling by u,v if using UnitVec() when calculatingg
	// basis vector / cross product - otherwise image gets stretched vertically
	// cam.Horizontal = cam.Horizontal.Mult(u) //
	// cam.Vertical = cam.Vertical.Mult(v)     

	fmt.Println("hor/ver", cam.Horizontal, cam.Vertical)
	// cam.Horizontal = NewVec3(4,0,0)
	// cam.Vertical = NewVec3(0,2,0)
	
	o := cam.Origin.Subtr(cam.Horizontal.DivF(2.0))
	o = o.Subtr(cam.Vertical.DivF(2.0))
	cam.Lower_left_corner = o.Subtr(w)
	// fmt.Println("llc:", cam.Lower_left_corner)
	return cam
}

func (c Camera) GetRay(u,v float32) Ray {
	u_horiz := c.Horizontal.MultF(u)
	v_vert := c.Vertical.MultF(v)
	dir := c.Lower_left_corner.Add(u_horiz)
	dir = dir.Add(v_vert)
	dir = dir.Subtr(c.Origin)
	return NewRay(c.Origin, dir)
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
type B struct {
	x float64
}

type Hittable interface {
	Hit(r *Ray, t_min, t_max float32, rec *HitRecord) bool
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

func (rec *HitRecord) set_face_normal(r *Ray, outward_normal * Vec3) {
	if r.Direction().Dot(*outward_normal) < 0 {
		rec.Normal = *outward_normal
		rec.FrontFace = true
	} else {
		rec.Normal = NewVec3(-outward_normal.At(0), -outward_normal.At(1), -outward_normal.At(2))
		rec.FrontFace = false
	}
}

func (b B) Hit(r *Ray, t_min, t_max float32, rec *HitRecord) bool{
	fmt.Println("B")
	return true
}

func write_color(cd Vec3, samples int) color.RGBA {
	scale := float32(1.0) / float32(samples)
	R := cd.At(0) * scale
	G := cd.At(1) * scale
	B := cd.At(2) * scale
	R = Clamp(R, float32(0.0), float32(0.9999))
	G = Clamp(G, float32(0.0), float32(0.9999))
	B = Clamp(B, float32(0.0), float32(0.9999))
	return color.RGBA{uint8(R*255), uint8(G*255), uint8(B*255), 255}	
}

func ray_color(r *Ray, objects []Hittable) Vec3 {
	// Iteration over the list of objects can me moved into a separate type
	// class in c++ HittableList (the world) but we leave it here for clarity
	
	rec := HitRecord{NewVec3(0,0,0), NewVec3(0,0,0), 1.0, true, -1}
	hit:= false
	closest_so_far := float32(math.Inf(1.0))

	// Testing Hit() without dynamic dispatch
	// s1:= Sphere{NewVec3(0,0,-1), 0.5}
	// s2:= Sphere{NewVec3(0,-100.5,-1), 100.0}	
	// hit_object := s1.Hit(r, 0.0, float32(closest_so_far), &rec) // half the time comapred
	// if hit_object{
	// 	closest_so_far = rec.T
	// 	hit = true
	// }
	// hit_object = s2.Hit(r, 0.0, float32(closest_so_far), &rec) // to for range
	// if hit_object{
	// 	closest_so_far = rec.T
	// 	hit = true
	// }

	// Option 2 - using dynamic dispatch
	for obj_id, object := range objects {
		// NOTE: we don't update hit var in here but inside the if block.
		// With a ray intersecting multiple objects the second object
		// (which can be furhter) will return False as the closest_so_far criteria no longer is met
		hit_object := object.Hit(r, 0.0, float32(closest_so_far), &rec) // this will populate HitRecord
		if hit_object{
			closest_so_far = rec.T
			hit = true
			rec.ObjectId = obj_id
		}
	}
	if hit {
		N := (rec.Normal.Add(NewVec3(1,1,1))).MultF(float32(0.5))
		// N = N.UnitVec()
		return N
	}
	// Background
	unit_direction := r.Direction().UnitVec()
	t := float32(0.5 * (unit_direction.At(1) + 1.0))
	sky := NewVec3(0.5, 0.7, 1.0).MultF(t)
	sky = sky.Add(NewVec3(1, 1, 1).MultF(1 - t))
	return sky
}

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")

func main() {
	flag.Parse()
    if *cpuprofile != "" {
        f, err := os.Create(*cpuprofile)
        if err != nil {
            fmt.Println(err)
        }
        pprof.StartCPUProfile(f)
        defer pprof.StopCPUProfile()
    }

	
	// The slice of types implements Hittable interface (HittableList in c++)n
	world := []Hittable{
		Sphere{NewVec3(0,0,-1), 0.5},
		Sphere{NewVec3(0,-100.5,-1), 100.0}}

	path := "img.png"
	path = os.Getenv("HOME") + "/storage/downloads/img.png" // termux preview
	
	fmt.Println("saving into:", path)

	f, err := os.Create(path)
	if err != nil {
		panic(err)
	}

	cam := NewCamera(NewVec3(0,0,0), NewVec3(0,0,-1), 200)

	upLeft := image.Point{0, 0}
	lowRight := image.Point{cam.Width, cam.Height}
	img := image.NewRGBA(image.Rectangle{upLeft, lowRight})

	start := time.Now()
	samples := 1

	for j := 0; j < cam.Height; j++ {
		for i := 0; i < cam.Width; i++ {
			pixel_color := NewVec3(0,0,0); _ = pixel_color
			for s:=0; s < samples; s++ {
				u := (float32(i) + RandFloat()) / float32(cam.Width-1)
				v := (float32(j) + RandFloat()) / float32(cam.Height-1)
				ray := cam.GetRay(u,v)
				pixel_color = pixel_color.Add(ray_color(&ray, world))
			}
			px_cd := write_color(pixel_color, samples)
			img.SetRGBA(i, cam.Height-j, px_cd)
		}
	}

	fmt.Println("time", time.Since(start))

	defer f.Close()
	if err = png.Encode(f, img); err != nil {
		fmt.Printf("failed to encode: %v", err)
	}

}
