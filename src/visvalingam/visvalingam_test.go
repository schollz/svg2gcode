package visvalingam

import "testing"

func TestSimplify(t *testing.T) {
	newcoords := Simplify([]float64{[]float64{1, 2}, []float64{3, 4}})
	if len(newcoords) == 0 {
		t.Errorf("should not be 0")
	}
}
