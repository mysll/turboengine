package tower

import (
	"math"
	"turboengine/gameplay/object"
)

type IdSet map[object.ObjectId]struct{}

type Tower struct {
	Ids      IdSet
	Watchers IdSet
}

func (t *Tower) clear() {
	t.Ids = IdSet{}
	t.Watchers = IdSet{}
}

func (t *Tower) add(obj object.ObjectId) bool {
	if _, ok := t.Ids[obj]; ok {
		return false
	}
	t.Ids[obj] = struct{}{}
	return true
}

func (t *Tower) remove(obj object.ObjectId) bool {
	if _, ok := t.Ids[obj]; !ok {
		return false
	}
	delete(t.Ids, obj)
	return true
}

func (t *Tower) getIds() []object.ObjectId {
	if len(t.Ids) == 0 {
		return nil
	}

	objs := make([]object.ObjectId, 0, len(t.Ids))
	for o := range t.Ids {
		objs = append(objs, o)
	}
	return objs
}

func (t *Tower) addWatcher(watcher object.ObjectId) bool {

	if _, ok := t.Watchers[watcher]; ok {
		return false
	}
	t.Watchers[watcher] = struct{}{}
	return true
}

func (t *Tower) removeWatcher(watcher object.ObjectId) {
	if _, ok := t.Watchers[watcher]; !ok {
		return
	}

	delete(t.Watchers, watcher)
}

func (t *Tower) getAllWatchers() []object.ObjectId {
	if len(t.Watchers) == 0 {
		return nil
	}

	result := make([]object.ObjectId, 0, len(t.Watchers))

	for o := range t.Watchers {
		result = append(result, o)
	}
	return result
}

func NewTower() *Tower {
	t := &Tower{}
	t.Ids = IdSet{}
	t.Watchers = IdSet{}
	return t
}

type TowerPos struct {
	X, Y int
}

type TowerAOI struct {
	width       float32
	height      float32
	towerWidth  float32
	towerHeight float32
	rangeLimit  int
	max         TowerPos
	towers      [][]*Tower
}

func (this *TowerAOI) Init() {
	iloop := int(math.Ceil(float64(this.width / this.towerWidth)))
	jloop := int(math.Ceil(float64(this.height / this.towerHeight)))
	this.max.X = iloop - 1
	this.max.Y = jloop - 1
	this.towers = make([][]*Tower, iloop)
	for i := 0; i < iloop; i++ {
		this.towers[i] = make([]*Tower, jloop)
		for j := 0; j < jloop; j++ {
			this.towers[i][j] = NewTower()
		}
	}
}

func (this *TowerAOI) Clear() {
	for i := 0; i <= this.max.X; i++ {
		for j := 0; j <= this.max.Y; j++ {
			this.towers[i][j].clear()
		}
	}
}
