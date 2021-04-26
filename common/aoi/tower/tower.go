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

func NewTowerAOI(w float32, h float32, tw float32, th float32, limit int) *TowerAOI {
	aoi := &TowerAOI{
		width:       w,
		height:      h,
		towerWidth:  tw,
		towerHeight: th,
		rangeLimit:  limit,
	}
	aoi.Init()
	return aoi
}

// Check if the pos is valid;
func (this *TowerAOI) checkPos(pos object.Vec3) bool {
	if pos.X() < 0 || pos.Z() < 0 || pos.X() >= this.width || pos.Z() >= this.height {
		return false
	}
	return true
}

// Trans the absolut pos to tower pos. For example : (210, 110} -> (1, 0), for tower width 200, height 200
func (this *TowerAOI) transPos(pos object.Vec3) TowerPos {
	return TowerPos{
		X: int(math.Floor(float64(pos.X() / this.towerWidth))),
		Y: int(math.Floor(float64(pos.Z() / this.towerHeight))),
	}
}

// getPosLimit Get the postion limit of given range,
func getPosLimit(pos TowerPos, ranges int, max TowerPos) (start TowerPos, end TowerPos) {

	if pos.X-ranges < 0 {
		start.X = 0
		end.X = 2 * ranges
	} else if pos.X+ranges > max.X {
		end.X = max.X
		start.X = max.X - 2*ranges
	} else {
		start.X = pos.X - ranges
		end.X = pos.X + ranges
	}

	if pos.Y-ranges < 0 {
		start.Y = 0
		end.Y = 2 * ranges
	} else if pos.Y+ranges > max.Y {
		end.Y = max.Y
		start.Y = max.Y - 2*ranges
	} else {
		start.Y = pos.Y - ranges
		end.Y = pos.Y + ranges
	}
	if start.X < 0 {
		start.X = 0
	}
	if end.X > max.X {
		end.X = max.X
	}
	if start.Y < 0 {
		start.Y = 0
	}
	if end.Y > max.Y {
		end.Y = max.Y
	}

	return
}

// isInRect  Check if the pos is in the rect
func isInRect(pos TowerPos, start TowerPos, end TowerPos) bool {
	return (pos.X >= start.X && pos.X <= end.X && pos.Y >= start.Y && pos.Y <= end.Y)
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
