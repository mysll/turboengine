package grid

import (
	"math"
	"testing"
)

func TestNewGridMap(t *testing.T) {
	f := float32(math.Sin(1))
	g := &Grid{}
	x := g.clampHeight(0, 1000, f)
	g.SetHeight(x)
	y := g.getClampHeight(0, 1000)
	t.Log(f, x, y)
	// output: 0.84147096 14117 0.84143883
}
