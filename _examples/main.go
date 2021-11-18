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
	"image/png"
	"os"
	"time"
	"flag"
	"runtime/pprof"
_	"sync"
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

	// Enable this to see BVH culling in action. 5sec vs 28sec for []Hittablelist
	// for i:=0; i<500; i++ {
	// 	world.Add(Sphere{NewVec3(0,float32(i)/10., float32(i + 1)), 0.5})
	// }
	
	bvh := NewBVHSplit(world.Objects,0,len(world.Objects))
	_ = bvh


	cam := NewCamera(NewVec3(0,0,0), NewVec3(0,0,-1), 2000)
	samples := 1
	
	start := time.Now()


	done := make(chan int)
	img := Render(cam, samples, &world, nil, done)  // pass bvh instead of nil to use BVH_node container

	// saving png

	path := os.Getenv("HOME") + "/storage/downloads/img.png" // termux preview
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		path = "img.png"
	}
	
	fmt.Println("\nsaving into:", path)

	f, err := os.Create(path)
	if err != nil {
		panic(err)
	}
	
	defer f.Close()
	if err = png.Encode(f, img); err != nil {
		fmt.Printf("failed to encode: %v", err)
	}

	fmt.Println("Waiting...")

	fmt.Println("time", time.Since(start))
}
