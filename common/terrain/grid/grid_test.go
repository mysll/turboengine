package grid

import (
	"math"
	"testing"
)

func TestNewGridMap(t *testing.T) {
	f := float32(math.Sin(1))
	x := clampHeight(0, 1000, f, 0xFFFFFF)
	y := getClampHeight(0, 1000, 0xFFFFFF, x)
	t.Log(f, x, y)
}
