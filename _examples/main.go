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
	. "example.com/raytrace"
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
}

// Polymorphism example
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

func (s Sphere) Hit(r *Ray, t_min, t_max float32, rec *HitRecord) bool {
	oc := r.Origin().Subtr(s.Center)
	a := r.Direction().Dot(r.Direction())
	half_b := oc.Dot(r.Direction())
	c := oc.Dot(oc) - (s.Radius * s.Radius)
	discriminant := float64(half_b * half_b - a*c) // finding roots
	if discriminant < 0{
		return false
	} else {
		// return float64(-half_b - float32(math.Sqrt(discriminant))) / (float64(a))
		// return true
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
	// fmt.Println("rootNew", root)
    rec.T = root;
    rec.P = r.At(rec.T);
	// fmt.Println("P New", rec.P.Subtr(s.Center), r.At(root).Subtr(s.Center), root, s.Center)
    // rec.Normal = (rec.P.Subtr(s.Center)).DivF(s.Radius)
    outward_normal := (rec.P.Subtr(s.Center)).DivF(s.Radius)
	// rec.Normal = (rec.P.Subtr(s.Center))
	// outward_normal := (rec.P.Subtr(s.Center))
    rec.set_face_normal(r, &outward_normal)
	
	return true;
}

func (rec *HitRecord) set_face_normal(r *Ray, outward_normal * Vec3) {
	if r.Direction().Dot(*outward_normal) < 0{ // TODO: this may reqiure flip
		rec.Normal = *outward_normal
		rec.FrontFace = true
	} else {
		rec.Normal = NewVec3(-outward_normal.At(0), -outward_normal.At(1), -outward_normal.At(2))
		rec.FrontFace = false
	}
	//normal = front_face ? outward_normal :-outward_normal;
}


func (b B) Hit(r *Ray, t_min, t_max float32, rec *HitRecord) bool{
	fmt.Println("B")
	return true
}

// Equation of sphere in vector form
// (P - C) dot (P - C) = r**2
// (A + tb - C) dot (A + tb - C) = r**2
// A + tb is a Ray
// quadratic equation x**2 + x + 1 =0 where t is unkown
// t**2 b dot b + 2t b dot (A-C) + (A-C) dot(A-C) - r**2 = 0
//      --a---       -----b-----   ------c--------
func hit_sphere(center *Vec3, radius float32, r *Ray) float64 {
	oc := r.Origin().Subtr(*center)
	a := r.Direction().Dot(r.Direction())
	half_b := oc.Dot(r.Direction())
	c := oc.Dot(oc) - (radius * radius)
	discriminant := float64(half_b * half_b - a*c) // finding roots
	if discriminant < 0{
		return -1.0
	} else {
		root := float64(-half_b - float32(math.Sqrt(discriminant))) / (float64(a)) 
		fmt.Println("root_old", root)
		return root
	}

}

func ray_color(r *Ray, sphere *Sphere) color.RGBA {


	aa := NewVec3(0, 0, -1)
	root := float32(hit_sphere(&aa, 0.5, r))
	if root > 0 {
		aa = r.At(float32(root)).Subtr( NewVec3(0,0,-1) )
		// fmt.Println("P Old", r.At(float32(root)), aa, root, NewVec3(0,0,-1))
		// aa = aa.UnitVec()
	}
	
	rec := HitRecord{NewVec3(0,0,0), NewVec3(0,0,0), 1.0, true}
	hit := sphere.Hit(r, 0.0, 100000, &rec)
	
	if hit { 
		N := (rec.Normal.Add(NewVec3(1,1,1))).MultF(float32(0.5))
		// N = N.UnitVec()
		R := N.At(0)
		G := N.At(1)
		B := N.At(2)

		// N := rec.Normal.UnitVec()
		// R := (N.At(0) + 1) * 0.5
		// G := (N.At(1) + 1) * 0.5
		// B := (N.At(2) + 1) * 0.5
		
	return color.RGBA{uint8(R*255), uint8(G*255), uint8(B*255), 255}
		
	}
	/*

	   hit_record rec;
    if (world.hit(r, 0, infinity, rec)) {
        return 0.5 * (rec.normal + color(1,1,1));
    }
    vec3 unit_direction = unit_vector(r.direction());
    auto t = 0.5*(unit_direction.y() + 1.0);
    return (1.0-t)*color(1.0, 1.0, 1.0) + t*color(0.5, 0.7, 1.0);
	
	aa := NewVec3(0, 0, -1)
	t := float32(hit_sphere(&aa, 0.5, r))
	if t > 0 {
		N := r.At(float32(t)).Subtr( NewVec3(0,0,-1) )
		N = N.UnitVec()
		R := (N.At(0) + 1) * 0.5
		G := (N.At(1) + 1) * 0.5
		B := (N.At(2) + 1) * 0.5
		return color.RGBA{uint8(R*255), uint8(G*255), uint8(B*255), 255}
	}
	*/

	

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

	rec := HitRecord{NewVec3(0,0,0), NewVec3(0,0,0), 1.0, true}
	ray := NewRay(NewVec3(0,0,0), NewVec3(0,0,0))
	a := []Hittable{Sphere{NewVec3(0,0,0), 1.0}, B{2}, Sphere{NewVec3(1,1,1), 0.5}}
	for _, i := range a {
		i.Hit(&ray, 0, 1, &rec)
	}

	sphere := Sphere{NewVec3(0,0,-1), 0.5}
	
	path := os.Getenv("HOME") + "/storage/downloads/img.png"
	fmt.Println("saving into:", path)

	f, err := os.Create(path)
	if err != nil {
		panic(err)
	}

	//Image
	aspect_ratio := 16.0 / 9.0
	width := 200
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
			cd := ray_color(&ray, &sphere)
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
