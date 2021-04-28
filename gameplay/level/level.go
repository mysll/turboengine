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

func CreateFromData(data string) *Level {
	config := &Config{}
	config.LoadFromData(data)
	return NewLevel(config)
}

func (l *Level) AddEntity(obj object.GameObject) {
	if _, ok := l.entities[obj.Id()]; ok {
		return
	}

	l.entities[obj.Id()] = obj
	if obj.HasView() {
		obj.AOI().Clear()
		pos := obj.Movement().Position()
		l.aoi.Enter(obj.Id(), pos, obj.AOI().ViewRange())
		obj.AOI().CachePosition(pos)
	}
}

func (l *Level) RemoveEntity(obj object.GameObject) {
	if _, ok := l.entities[obj.Id()]; !ok {
		return
	}

	if obj.HasView() {
		l.aoi.Leave(obj.Id(), obj.AOI().Position(), obj.AOI().ViewRange())
		obj.AOI().Clear()
	}

	delete(l.entities, obj.Id())
}

func (l *Level) GetEntityById(id ObjectId) object.GameObject {
	return l.entities[id]
}
