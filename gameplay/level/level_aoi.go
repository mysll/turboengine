package level

import (
	. "turboengine/common/datatype"
	"turboengine/common/log"
)

func (l *Level) OnEnterAOI(self, target ObjectId) {
	entity := l.GetEntityById(self)
	if entity == nil {
		log.Errorf("entity not found, %d", self)
		return
	}
	entity.AOI().AddViewObj(target)
}

func (l *Level) OnLeaveAOI(self, target ObjectId) {
	entity := l.GetEntityById(self)
	if entity == nil {
		log.Errorf("entity not found, %d", self)
		return
	}
	entity.AOI().RemoveViewObj(target)
}

// Move 移动一段距离
func (l *Level) Move(obj ObjectId, step Vec3) {
	entity := l.GetEntityById(obj)
	if entity == nil {
		log.Errorf("entity not found, %d", obj)
		return
	}

	old := entity.Movement().Position()
	pos := entity.Movement().Translate(step)
	if old.Equal(pos) {
		return
	}
	// TODO: 检查是否能够移动，地形和碰撞检查，不通过则回退到上一步的位置

	if entity.HasView() {
		if l.aoi.Move(obj, entity.AOI().Position(), pos, entity.AOI().ViewRange()) {
			entity.AOI().CachePosition(pos)
		}
	}
}

// Locate 移动到pos位置
func (l *Level) Locate(obj ObjectId, pos Vec3) {
	entity := l.GetEntityById(obj)
	if entity == nil {
		log.Errorf("entity not found, %d", obj)
		return
	}

	old := entity.Movement().Position()
	if old.Equal(pos) {
		return
	}
	entity.Movement().MoveTo(pos)
	// TODO: 检查是否能够移动，地形和碰撞检查，不通过则回退到上一步的位置

	if entity.HasView() {
		if l.aoi.Move(obj, entity.AOI().Position(), pos, entity.AOI().ViewRange()) {
			entity.AOI().CachePosition(pos)
		}
	}
}
