package level

import (
	. "turboengine/common/datatype"
	"turboengine/common/log"
)

func (l *Level) OnEnterAOI(watcher, target ObjectId) {
	entity := l.GetEntityById(watcher)
	if entity == nil {
		log.Errorf("watcher not found, %d", watcher)
	}
	entity.AOI().AddViewObj(target)
}

func (l *Level) OnLeaveAOI(watcher, target ObjectId) {
	entity := l.GetEntityById(watcher)
	if entity == nil {
		log.Errorf("watcher not found, %d", watcher)
	}
	entity.AOI().RemoveViewObj(target)
}
