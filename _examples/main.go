// A driver test program
package main

import (
	"fmt"
	"gonum.org/v1/gonum/floats"
	"image"
	"image/color"
	"image/png"
	"os"
	"example.com/raytrace"
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


func main() {
	a := []Mesh{A{1}, B{2}, A{3}}
	for _, i := range a {
		i.Hello()
	}

	rr := raytrace.Ray{ []float64{0,0,0}, []float64{1,2,3} }
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

	
	
	width := 200
	height := 100
	upLeft := image.Point{0, 0}
	lowRight := image.Point{width, height}
	img := image.NewRGBA(image.Rectangle{upLeft, lowRight})

	var r, g uint8
	for j := height - 1; j >= 0; j-- {
		for i := 0; i < width; i++ {
			r = uint8(255 * float64(i) / float64(width-1))
			g = uint8(255 * float64(j) / float64(height))
			img.SetRGBA(i, j, color.RGBA{r, g, 0, 255})
		}
	}

	defer f.Close()
	if err = png.Encode(f, img); err != nil {
		fmt.Printf("failed to encode: %v", err)
	}

}
