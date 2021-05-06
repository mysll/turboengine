package grid

import (
	"math"
	. "turboengine/common/datatype"
	. "turboengine/common/terrain"
)

const limit float32 = 0xFFFFFF

type Grid struct {
	flag uint32      // [8bit类型][24bit高度]
	data interface{} // 附加数据
}

func NewGrid() *Grid {
	return &Grid{}
}

func (g *Grid) SetHeight(h uint32) {
	g.flag = (g.flag & 0xFF000000) | (h & 0xFFFFFF)
}

func (g *Grid) SetType(t uint32) {
	g.flag = (t << 24) | (g.flag & 0xFFFFFF)
}

func (g *Grid) Height() uint32 {
	return g.flag & 0xFFFFFF
}

func (g *Grid) Type() uint32 {
	return g.flag >> 24
}

// 映射到[0,limit]
func (g *Grid) clampHeight(min, max, height float32) uint32 {
	return uint32((height - min) / (max - min) * limit)
}

// 映射回[min, max]
func (g *Grid) getClampHeight(min, max float32) float32 {
	return float32(g.Height())*(max-min)/limit + min
}

type GridMap struct {
	pixel     Vec2     // 原点坐标(像素)
	zone      Vec2     // 格子长和宽(像素)
	unit      Vec2     // n坐标单位=1像素单位
	row       uint32   // 格子最大行
	col       uint32   // 格子最大列
	origin    Vec2     // 原点坐标
	size      Vec2     // 格子大小
	grids     [][]Grid // 格子数据
	maxHeight float32
	minHeight float32
}

func (g *GridMap) init() {
	g.origin = V2(g.pixel.X()*g.unit.X(), g.pixel.Y()*g.unit.Y())
	g.size = V2(g.zone.X()*g.unit.X(), g.zone.Y()*g.unit.Y())
}

func (g *GridMap) LoadFromFile(f string) {
	// var g grid
	// grid.setHeight(clampHeight(g.minHeight, g.maxHeight,height, 0xFFFFFF)
}

func (g *GridMap) CanWalk(pos Vec3) (b bool) {
	grid := g.getGrid(pos)
	if grid == nil {
		return false
	}
	return grid.Type()|MAP_TYPE_WATER|MAP_TYPE_TREE|MAP_TYPE_BUILD != 0
}

func (g *GridMap) LineCanWalk(step float32, start, end Vec3) bool {
	sx := end.X() - start.X()
	sz := end.Z() - start.Z()
	l := float32(math.Sqrt(float64(sx*sx + sz*sz)))
	if l < step {
		return g.CanWalk(start) && g.CanWalk(end)
	}
	xstep := step * (sx / l)
	zstep := step * (sz / l)
	sx = start.X()
	sz = start.Z()
	loop := int(l / step)
	for i := 0; i < loop; i++ {
		sx += xstep
		sz += zstep
		if !g.CanWalk(V3(sx, 0, sz)) {
			return false
		}
	}
	return false
}

func (g *GridMap) Height(pos Vec3) float32 {
	grid := g.getGrid(pos)
	if grid == nil {
		return float32(math.NaN())
	}
	return grid.getClampHeight(g.minHeight, g.maxHeight)
}

func (g *GridMap) MapType(pos Vec3) int {
	grid := g.getGrid(pos)
	if grid == nil {
		return MAP_TYPE_NONE
	}
	return int(grid.Type())
}

func (g *GridMap) getGrid(pos Vec3) *Grid {
	x := (pos.X() - g.origin.X()) / g.size.X()
	y := (pos.Z() - g.origin.Y()) / g.size.Y()
	row := uint32(math.Ceil(float64(x)))
	col := uint32(math.Ceil(float64(y)))
	if row > g.row || col > g.col {
		return nil
	}
	return &g.grids[row][col]
}
