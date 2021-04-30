package grid

import (
	"encoding/binary"
	"math"
	. "turboengine/common/datatype"
	. "turboengine/common/terrain"
)

type GridPos struct {
	Row, Col uint16
}

// rectangle
type Grid struct {
	pos    GridPos                // 格子坐标(行,列)
	gtype  byte                   // 格子类型
	height [3]byte                // 格子高度  3byte
	data   map[string]interface{} // 附加数据
}

func NewGrid() *Grid {
	return &Grid{}
}

func (g *Grid) setType(t int) {
	g.gtype = byte(t)
}

func (g *Grid) Type() uint32 {
	return uint32(g.gtype)
}

func (g *Grid) setHeight(f float32) {
	bits := math.Float32bits(f)
	b := make([]byte, 4)
	binary.LittleEndian.PutUint32(b, bits)
	g.height[0] = b[1]
	g.height[1] = b[2]
	g.height[2] = b[3]
}

func (g *Grid) Height() float32 {
	b := make([]byte, 4)
	b[0] = 0
	b[1] = g.height[0]
	b[2] = g.height[1]
	b[3] = g.height[2]
	bits := binary.LittleEndian.Uint32(b)
	float := math.Float32frombits(bits)
	return float
}

type grid uint32

func (g *grid) SetHeight(h uint32) {
	*g = grid(uint32(*g)&0xFF000000 | (h & 0xFFFFFF))
}

func (g *grid) SetType(t uint32) {
	*g = grid((t << 24) | (uint32(*g) & 0xFFFFFF))
}

func (g grid) Height() uint32 {
	return uint32(g) & 0xFFFFFF
}

func (g grid) Type() uint32 {
	return uint32(g) >> 24
}

// 映射到[0,0xFFFFFF]
func clampHeight(min, max, height, limit float32) uint32 {
	return uint32((height - min) / (max - min) * limit)
}

// 映射回[min, max]
func getClampHeight(min, max, limit float32, height uint32) float32 {
	return float32(height)*(max-min)/limit + min
}

type GridMap struct {
	pixel     Vec2     // 原点坐标(像素)
	zone      Vec2     // 格子长和宽(像素)
	unit      Vec2     // n坐标单位=1像素单位
	max       GridPos  // 格子最大(行,列)
	origin    Vec2     // 原点坐标
	size      Vec2     // 格子大小
	grids     [][]grid // 格子数据
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
	return getClampHeight(g.minHeight, g.maxHeight, 0xFFFFFF, grid.Height())
}

func (g *GridMap) MapType(pos Vec3) int {
	grid := g.getGrid(pos)
	if grid == nil {
		return MAP_TYPE_NONE
	}
	return int(grid.Type())
}

func (g *GridMap) getGrid(pos Vec3) *grid {
	x := (pos.X() - g.origin.X()) / g.size.X()
	y := (pos.Z() - g.origin.Y()) / g.size.Y()
	row := uint16(math.Ceil(float64(x)))
	col := uint16(math.Ceil(float64(y)))
	if row > g.max.Row || col > g.max.Col {
		return nil
	}
	return &g.grids[row][col]
}
