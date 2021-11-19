// A driver test program
// to debug: go build -gcflags="all=-N -l" main.go
//
// Profiling:
// $ go build
// $ ./_examples -cpuprofile aaa
// $ go tool pprof aaa
// (pprof) top 10

package main

import (
	. "github.com/kubaroth/Vec3"
	"errors"
	"fmt"
	"image"
	"image/png"
	"os"
	"time"
	"flag"
	"runtime/pprof"
	"sync"
)



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

	world := HittableList{}
	world.Add(Sphere{NewVec3(0,0,-1), 0.5})
	world.Add(Sphere{NewVec3(0,-100.5,-1), 100.0})
	_ = world
	
	bvh := NewBVHSplit(world.Objects,0,len(world.Objects))
	_ = bvh

	
	path := os.Getenv("HOME") + "/storage/downloads/img.png" // termux preview
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		path = "img.png"
	}
	
	fmt.Println("saving into:", path)

	f, err := os.Create(path)
	if err != nil {
		panic(err)
	}

	cam := NewCamera(NewVec3(0,0,0), NewVec3(0,0,-1), 2000)

	upLeft := image.Point{0, 0}
	lowRight := image.Point{cam.Width, cam.Height}
	img := image.NewRGBA(image.Rectangle{upLeft, lowRight})

	start := time.Now()
	samples := 16

	var wg sync.WaitGroup; _ = wg
	wg.Add(cam.Height)

	rand_seeds := make([]float32, samples)
	for i:=0; i<samples; i++ {
		rand_seeds[i] = RandFloat()
	}
	
	for j := 0; j < cam.Height; j++ {

		go func(j int) {
		defer wg.Done()


		for i := 0; i < cam.Width; i++ {

			pixel_color := NewVec3(0,0,0); _ = pixel_color
			for s:=0; s < samples; s++ {
				rr := rand_seeds[s]  // reuse rand and avoids global lock!
				u := (float32(i) + rr) / float32(cam.Width-1)
				v := (float32(j) + rr) / float32(cam.Height-1)
				_, _ = u, v
				ray := cam.GetRay(u,v)
				// pixel_color = pixel_color.Add(RayColorBVH(&ray, bvh)) // BVH scene
				pixel_color = pixel_color.Add(RayColorArray(&ray, world))  // flat list scene
			}
			px_cd := Write_color(pixel_color, samples)
			// _ = px_cd
			img.SetRGBA(i, cam.Height-j, px_cd)

		}

		}(j)
	}

	wg.Wait()
	
	fmt.Println("time", time.Since(start))

	defer f.Close()
	if err = png.Encode(f, img); err != nil {
		fmt.Printf("failed to encode: %v", err)
	}

}
