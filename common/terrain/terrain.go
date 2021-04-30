package terrain

import . "turboengine/common/datatype"

/* 地图的坐标系 */
//                (0, z)   (x, z)
//                ___________
//   y  z         |         |
//   | /          |         |
//   |/           |         |
//   +-------x    |_________|
//
//               (0, 0)    (x, 0)

const (
	MAP_TYPE_NONE  = iota
	MAP_TYPE_LAND  = 1      // 陆地
	MAP_TYPE_WATER = 1 << 1 // 水
	MAP_TYPE_TREE  = 1 << 2 // 树
	MAP_TYPE_BUILD = 1 << 3 // 建筑
)

type Terrain interface {
	// Load 从文件加载
	LoadFromFile(f string)
	// CanWalk 某个点是否可以行走
	CanWalk(pos Vec3) bool
	// LineCanWalk 两个点之间是否可以行走
	LineCanWalk(step float32, start, end Vec3) bool
	// Height 获取高度
	Height(pos Vec3) float32
	// 某个点的地图类型(MAP_TYPE)
	MapType(pos Vec3) int
}
