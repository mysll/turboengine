package level

import (
	. "turboengine/common/datatype"
)

const (
	MAP_TYPE_NONE  = iota
	MAP_TYPE_LAND  // 陆地
	MAP_TYPE_WATER // 水
)

type MapCreater func() Terrain

var mapType map[string]MapCreater

func RegMap(typ string, f MapCreater) {
	mapType[typ] = f
}

type Terrain interface {
	// Load 从文件加载
	LoadFromFile(f string)
	// CanWalk 某个点是否可以行走
	CanWalk(pos Vec3) bool
	// LineCanWalk 两个点之间是否可以行走
	LineCanWalk(start, end Vec3) bool
	// Height 获取高度
	Height(pos Vec3) float32
	// 某个点的地图类型(MAP_TYPE)
	MapType(pos Vec3) int
}

func (l *Level) LoadTerrain(path string) {
	f, ok := mapType[l.config.TerrainType]
	if !ok {
		panic("map type not found")
	}

	terrain := f()
	terrain.LoadFromFile(path)
	l.terrain = terrain
}
