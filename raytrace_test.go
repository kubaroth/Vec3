package raytrace
import (
	"testing"
	_ "fmt"
	"gonum.org/v1/gonum/floats"
)

func TestRay1(t *testing.T){
	r := Ray{ []float64{0,0,0}, []float64{1,2,3} }
	vec := r.At(2)
	want := []float64{2,4,6}
	if !floats.Same(vec, want) {
		t.Errorf(" %v != %v", vec, want)
	}

}
