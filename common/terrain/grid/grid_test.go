package grid

import (
	"fmt"
	"testing"
)

func TestNewGridMap(t *testing.T) {
	f := float32(100.12332211)
	grids := make([][]*Grid, 10)
	for row := 0; row < 10; row++ {
		grids[row] = make([]*Grid, 10)
		for col := 0; col < 10; col++ {
			g := NewGrid()
			grids[row][col] = g
			g.setType(col)
			f += 0.01
			g.setHeight(f)
			fmt.Println(fmt.Sprintf("type: %d, height: %f", g.getType(), g.getHeight()))
		}
	}
}
