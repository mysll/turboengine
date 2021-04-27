package level

import (
	"turboengine/common/aoi"
	"turboengine/common/aoi/tower"
	. "turboengine/common/datatype"
	"turboengine/gameplay/object"
)

type EntityContainer map[ObjectId]object.GameObject

type Level struct {
	aoi      aoi.AOIMgr
	entities EntityContainer
	config   *Config
}

func NewLevel(config *Config) *Level {
	if config == nil {
		panic("config is nil")
	}
	l := &Level{
		config:   config,
		entities: EntityContainer{},
	}

	aoi := tower.NewTowerAOI(config.MinX, config.MinY, config.MaxX, config.MaxY,
		config.TileWidth, config.TileHeight,
		config.ViewMaxRange, l)
	l.aoi = aoi
	return l
}

func CreateFromFile(f string) *Level {
	config := &Config{}
	config.LoadFromFile(f)
	return NewLevel(config)
}

func (l *Level) OnEnterAOI(watcher, target ObjectId) {

}

func (l *Level) OnLeaveAOI(watcher, target ObjectId) {

}
