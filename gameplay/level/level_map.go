package level

import (
	. "turboengine/common/terrain"
	. "turboengine/common/terrain/grid"
	. "turboengine/common/terrain/nav"
)

type MapCreater func() Terrain

var mapType map[string]MapCreater

func RegMap(typ string, f MapCreater) {
	mapType[typ] = f
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

func init() {
	RegMap("grid", func() Terrain {
		return &GridMap{}
	})
	RegMap("nav", func() Terrain {
		return &NavMesh{}
	})
}
