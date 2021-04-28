package terrain

import (
	"math"
	. "turboengine/common/datatype"
	"turboengine/gameplay/level"
)

// nav mesh
type NavMesh struct {
}

//TODO 实现Terrain接口
// Load 从文件加载
func (n *NavMesh) LoadFromFile(f string) {

}

// CanWalk 某个点是否可以行走
func (n *NavMesh) CanWalk(pos Vec3) bool {
	return false
}

// LineCanWalk 两个点之间是否可以行走
func (n *NavMesh) LineCanWalk(start, end Vec3) bool {
	return false
}

// Height 获取高度
func (n *NavMesh) Height(pos Vec3) float32 {
	return float32(math.NaN())
}

// 某个点的地图类型(MAP_TYPE)
func (n *NavMesh) MapType(pos Vec3) int {
	return level.MAP_TYPE_NONE
}

func init() {
	level.RegMap("navmesh", func() level.Terrain {
		return &NavMesh{}
	})
}
