// A driver test program
// to debu: go build -gcflags="all=-N -l" main.go

package main

import (
	"fmt"
	"gonum.org/v1/gonum/floats"
	"image"
	"image/color"
	"image/png"
	"math"
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

// Helper functions - promary advantage compared to gonum floats is the return value
// which helps to write complex expressions
func vec3_div(v []float64, x float64) []float64{
	return []float64{v[0]/x, v[1]/x , v[2]/x}
}
func vec3_mul(v []float64, x float64) []float64{
	return []float64{v[0]*x, v[1]*x , v[2]*x}
}
func vec3_length(v []float64) float64{
 	return math.Sqrt(v[0]*v[0] + v[1]*v[1] + v[2]*v[2]) 
}

func unit_vector(v []float64) []float64{
	length := vec3_length(v)
	return vec3_div(v, length)
} 


func ray_color(r raytrace.Ray) color.RGBA {
	unit_direction := unit_vector(r.Direction())
	t := 0.5 * (unit_direction[1] + 1.0)
	temp := []float64{ 0.5, 0.7, 1.0 }
	floats.Mul(temp, []float64{ t, t, t })
	temp2 := []float64{ 1, 1, 1 }
	floats.Mul(temp2, []float64{ 1-t, 1-t, 1-t })
	floats.Add(temp, temp2)
	temp = vec3_mul(temp, 255) // NOTE: rember to shift range to 0-255
	return color.RGBA{uint8(temp[0]), uint8(temp[1]), uint8(temp[2]), 255}
	
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

	//Image
	aspect_ratio := 16.0/9.0
	width := 200
	height := int(float64(width)/aspect_ratio)

	// Camera
	viewport_height := 2.0
	viewport_width := aspect_ratio *viewport_height
	focal_length := 1.0; _ = focal_length

	origin := []float64{0,0,0}
	horizontal := []float64{float64(viewport_width), 0, 0}
	vertical := []float64{0, float64(viewport_height), 0}
	horizontal_2 := vec3_div(horizontal, 2)
	vertical_2 := vec3_div(vertical, 2)
	lower_left_corner := []float64{0,0,0}
	floats.Add(lower_left_corner, origin)
	floats.Sub(lower_left_corner, horizontal_2)
	floats.Sub(lower_left_corner, vertical_2)
	floats.Add(lower_left_corner, []float64{0,0,focal_length})
	
	
	upLeft := image.Point{0, 0}
	lowRight := image.Point{width, height}
	img := image.NewRGBA(image.Rectangle{upLeft, lowRight})

	// var r, g uint8
	for j := height - 1; j >= 0; j-- {
		for i := 0; i < width; i++ {
			// r = uint8(255 * float64(i) / float64(width-1))
			// g = uint8(255 * float64(j) / float64(height))
			//img.SetRGBA(i, j, color.RGBA{r, g, 0, 255})

			u := float64(i) / float64(width - 1)
			v := float64(j) / float64(height - 1)
			dir := []float64{0,0,0}
			floats.Add(dir, lower_left_corner)
			floats.Add(dir, vec3_mul(horizontal, u))
			floats.Add(dir, vec3_mul(horizontal, v))
			floats.Sub(dir, origin)
			ray := raytrace.Ray{origin, dir} ; _ = ray
			//fmt.Println("ray", ray)
			cd := ray_color(ray)
			// fmt.Println(cd)
			img.SetRGBA(i, j, cd)
			
		}
	}

	defer f.Close()
	if err = png.Encode(f, img); err != nil {
		fmt.Printf("failed to encode: %v", err)
	}

}
