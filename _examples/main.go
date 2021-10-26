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
	"gonum.org/v1/gonum/floats"
	"image"
	"image/color"
	"image/png"
	"math"
	"os"
	"time"
	"flag"
	"runtime/pprof"
)

// Polymorphism example
type A struct {
	x int
}
type B struct {
	x float64
}

type Mesh interface {
	Hello()
}

func (a A) Hello() {
	fmt.Println(a.x)
}
func (b B) Hello() {
	fmt.Println(b.x)
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
	b := oc.Dot(r.Direction()) * 2.0
	c := oc.Dot(oc) - (radius * radius)
	discriminant := float64(b*b - 4*a*c) // finding roots
	if discriminant < 0{
		return -1.0
	} else {
		return float64(-b - float32(math.Sqrt(discriminant))) / (2.0 * float64(a))
	}

}

func ray_color(r *Ray) color.RGBA {
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

	unit_direction := r.Direction().UnitVec()
	t = float32(0.5 * (unit_direction.At(1) + 1.0))
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

	
	a := []Mesh{A{1}, B{2}, A{3}}
	for _, i := range a {
		i.Hello()
	}

	rr := NewRay(NewVec3(0, 0, 0), NewVec3(1, 2, 3))
	_ = rr

	v1 := []float64{1, 3, -5}
	v2 := []float64{4, -2, -1}
	fmt.Println(floats.Dot(v1, v2))

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
			cd := ray_color(&ray)
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
