package tower

import (
	"testing"
	"turboengine/gameplay/object"
)

type aoi struct {
	t *testing.T
}

func (a *aoi) OnEnterAOI(watcher, target object.ObjectId) {
	a.t.Logf("[%d] :%d enter \n", watcher, target)
}

func (a *aoi) OnLeaveAOI(watcher, target object.ObjectId) {
	a.t.Logf("[%d] :%d leave \n", watcher, target)
}

func TestNewTowerAOI(t *testing.T) {
	toi := NewTowerAOI(1000, 1000, 50, 50, 10, &aoi{t})
	toi.Enter(1, object.Vec3{75, 0, 75}, 1)
	toi.Enter(2, object.Vec3{125, 0, 75}, 1)
	toi.Move(1, object.Vec3{75, 0, 75}, object.Vec3{25, 0, 125}, 1)
	toi.Move(2, object.Vec3{125, 0, 75}, object.Vec3{75, 0, 75}, 1)
	toi.Level(1, object.Vec3{25, 0, 125}, 1)
}
