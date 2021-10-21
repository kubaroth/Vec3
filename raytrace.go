package raytrace

import (
	"fmt"
	"gonum.org/v1/gonum/floats"
)

type Ray struct {
	// It would be nice to use [3]floats64 instead of slices here as we know the
	// size of origin and direction but  gonum/floats works with slices
	Orig []float64
	Dir  []float64
}

func NewRay(){
	fmt.Println("NewRay")
}

func (r Ray) Origin() []float64 {
	return r.Orig
}
func (r Ray) Direction() []float64 {
	return r.Dir
}

func (r Ray) At(t float64) []float64 {
	ret := make([]float64, 3)
    dir := []float64{r.Dir[0], r.Dir[1], r.Dir[2]}
	floats.Mul(dir, []float64{t,t,t})
	floats.AddTo(ret, r.Orig, dir) 
	return ret
}
