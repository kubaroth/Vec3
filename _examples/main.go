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

func ray_color(r *Ray, objects []Hittable) color.RGBA {
	// Iteration over the list of objects can me moved into a separate type
	// class in c++ HittableList (the world) but we leave it here for clarity
	rec := HitRecord{NewVec3(0,0,0), NewVec3(0,0,0), 1.0, true, -1}
	hit:= false
	closest_so_far := float32(math.Inf(1.0))

	s1:= Sphere{NewVec3(0,0,-1), 0.5}
	s2:= Sphere{NewVec3(0,-100.5,-1), 100.0}
	hit = s1.Hit(r, 0.0, float32(closest_so_far), &rec) // half the time comapred
	hit = s2.Hit(r, 0.0, float32(closest_so_far), &rec) // to for range
	
	// for obj_id, object := range objects {
	// 	// NOTE: we don't update hit var in here but inside the if block.
	// 	// With a ray intersecting multiple objects the second object
	// 	// (which can be furhter) will return False as the closest_so_far criteria no longer is met
	// 	hit_object := object.Hit(r, 0.0, float32(closest_so_far), &rec) // this will populate HitRecord
	// 	if hit_object{
	// 		closest_so_far = rec.T
	// 		hit = true
	// 		rec.ObjectId = obj_id
	// 	}
	// }
	if hit {
		N := (rec.Normal.Add(NewVec3(1,1,1))).MultF(float32(0.5))
		// N = N.UnitVec()
		R := N.At(0)
		G := N.At(1)
		B := N.At(2)
		return color.RGBA{uint8(R*255), uint8(G*255), uint8(B*255), 255}
	}
	// Background
	unit_direction := r.Direction().UnitVec()
	t := float32(0.5 * (unit_direction.At(1) + 1.0))
	sky := NewVec3(0.5, 0.7, 1.0).MultF(t)
	sky = sky.Add(NewVec3(1, 1, 1).MultF(1 - t))
	sky = sky.MultF(255) // NOTE: rember to shift range to 0-255
	return color.RGBA{uint8(sky.At(0)), uint8(sky.At(1)), uint8(sky.At(2)), 255}
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

	path := os.Getenv("HOME") + "/storage/downloads/img.png"
	fmt.Println("saving into:", path)

	f, err := os.Create(path)
	if err != nil {
		panic(err)
	}

	//Image
	aspect_ratio := 16.0 / 9.0
	width := 2000
	height := int(float64(width) / aspect_ratio)
	fmt.Printf("image res", width, height)
	
	// Camera
	viewport_height := 2.0
	viewport_width := aspect_ratio * viewport_height
	focal_length := 1.0
	_ = focal_length

	origin := NewVec3(0, 0, 0)
	horizontal := NewVec3(float32(viewport_width), 0, 0)
	vertical := NewVec3(0, float32(viewport_height), 0)

	lower_left_corner := origin.Subtr(horizontal.DivF(2.0))
	llc := lower_left_corner.Subtr(vertical.DivF(2.0))
	llc = llc.Subtr(NewVec3(0, 0, float32(focal_length)))

	upLeft := image.Point{0, 0}
	lowRight := image.Point{width, height}
	img := image.NewRGBA(image.Rectangle{upLeft, lowRight})

	start := time.Now()
	// var r, g uint8
	for j := 0; j < height; j++ {
		for i := 0; i < width; i++ {
			// r = uint8(255 * float64(i) / float64(width-1))
			// g = uint8(255 * float64(j) / float64(height))
			// img.SetRGBA(i, height-j, color.RGBA{r, g, 0, 255})

			u := float32(i) / float32(width-1)
			v := float32(j) / float32(height-1)

			u_horiz := horizontal.MultF(u)
			v_vert := vertical.MultF(v)
			dir := llc.Add(u_horiz)
			dir = dir.Add(v_vert)
			dir = dir.Subtr(origin)
			
			ray := NewRay(origin, dir)
			// fmt.Println("ray", ray)
			cd := ray_color(&ray, world)
			// fmt.Println(i, j)
			img.SetRGBA(i, height-j, cd)
		}
	}

	fmt.Println("time", time.Since(start))

	defer f.Close()
	if err = png.Encode(f, img); err != nil {
		fmt.Printf("failed to encode: %v", err)
	}

}
