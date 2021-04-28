package terrain

import (
	"math"
	. "turboengine/common/datatype"
	"turboengine/gameplay/level"
)

// 格子地图
type TileMap struct {
}

//TODO 实现Terrain接口
// Load 从文件加载
func (t *TileMap) LoadFromFile(f string) {

}

// CanWalk 某个点是否可以行走
func (t *TileMap) CanWalk(pos Vec3) bool {
	return false
}

// LineCanWalk 两个点之间是否可以行走
func (t *TileMap) LineCanWalk(start, end Vec3) bool {
	return false
}

// Height 获取高度
func (t *TileMap) Height(pos Vec3) float32 {
	return float32(math.NaN())
}

// 某个点的地图类型(MAP_TYPE)
func (t *TileMap) MapType(pos Vec3) int {
	return level.MAP_TYPE_NONE
}

func init() {
	level.RegMap("tile", func() level.Terrain {
		return &TileMap{}
	})
}
