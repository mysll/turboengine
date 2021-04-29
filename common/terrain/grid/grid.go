package grid

import (
	"math"
	. "turboengine/common/datatype"
	. "turboengine/common/terrain"
)

type IdSet map[ObjectId]struct{}

type Callback interface {
	OnEnterGrid(g *Grid, target ObjectId)
	OnLeaveGrid(g *Grid, target ObjectId)
}

type GridPos struct {
	Row, Col int
}

// rectangle
type Grid struct {
	pos    GridPos                // 格子坐标(行,列)
	gtype  int                    // 格子类型
	height float32                // 格子高度
	data   map[string]interface{} // 附加数据
	ids    IdSet                  // 格子中的玩家
}

func NewGrid() *Grid {
	return &Grid{}
}

type GridMap struct {
	pixel  Vec2      // 原点坐标(像素)
	zone   Vec2      // 格子长和宽(像素)
	unit   Vec2      // n坐标单位=1像素单位
	max    GridPos   // 格子最大(行,列)
	origin Vec2      // 原点坐标
	size   Vec2      // 格子大小
	grids  [][]*Grid // 格子数据
	cb     Callback
}

func NewGridMap() *GridMap {
	return &GridMap{}
}

func (g *GridMap) init() {
	g.origin = V2(g.pixel.X()*g.unit.X(), g.pixel.Y()*g.unit.Y())
	g.size = V2(g.zone.X()*g.unit.X(), g.zone.Y()*g.unit.Y())
}

// Load 从文件加载
func (g *GridMap) LoadFromFile(f string) {

}

// CanWalk 某个点是否可以行走
func (g *GridMap) CanWalk(pos Vec3) (b bool) {
	grid := g.getGrid(pos)
	if grid == nil {
		return false
	}
	return grid.gtype|MAP_TYPE_WATER|MAP_TYPE_TREE|MAP_TYPE_BUILD != 0
}

// Walk
func (g *GridMap) Walk(obj ObjectId, old Vec3, new Vec3) {
	/*ngrid := g.grids[new.Row][new.Col]
	if old == new {
		return
	}
	ogrid := g.grids[old.Row][old.Col]
	g.cb.OnLeaveGrid(ogrid, obj)
	delete(ogrid.ids, obj)
	g.cb.OnEnterGrid(ngrid, obj)
	ngrid.ids[obj] = struct{}{}*/
}

// LineCanWalk 两个点之间是否可以行走
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

// Height 获取高度
func (g *GridMap) Height(pos Vec3) float32 {
	grid := g.getGrid(pos)
	if grid == nil {
		return float32(math.NaN())
	}
	return grid.height
}

// 某个点的地图类型(MAP_TYPE)
func (g *GridMap) MapType(pos Vec3) int {
	grid := g.getGrid(pos)
	if grid == nil {
		return MAP_TYPE_NONE
	}
	return grid.gtype
}

func (g *GridMap) getGrid(pos Vec3) *Grid {
	x := (pos.X() - g.origin.X()) / g.size.X()
	y := (pos.Z() - g.origin.Y()) / g.size.Y()
	row := int(math.Ceil(float64(x)))
	col := int(math.Ceil(float64(y)))
	if row > g.max.Row || col > g.max.Col {
		return nil
	}
	return g.grids[row][col]
}
