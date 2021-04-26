package tower

import (
	"math"
	"turboengine/gameplay/object"
)

type IdSet map[object.ObjectId]struct{}

type Callback interface {
	OnEnterAOI(watcher, target object.ObjectId)
	OnLeaveAOI(watcher, target object.ObjectId)
}

type tower struct {
	Ids      IdSet
	Watchers IdSet
}

func (t *tower) clear() {
	t.Ids = IdSet{}
	t.Watchers = IdSet{}
}

func (t *tower) add(obj object.ObjectId) bool {
	if _, ok := t.Ids[obj]; ok {
		return false
	}
	t.Ids[obj] = struct{}{}
	return true
}

func (t *tower) remove(obj object.ObjectId) bool {
	if _, ok := t.Ids[obj]; !ok {
		return false
	}
	delete(t.Ids, obj)
	return true
}

func (t *tower) getIds() []object.ObjectId {
	if len(t.Ids) == 0 {
		return nil
	}

	objs := make([]object.ObjectId, 0, len(t.Ids))
	for o := range t.Ids {
		objs = append(objs, o)
	}
	return objs
}

func (t *tower) addWatcher(watcher object.ObjectId) bool {

	if _, ok := t.Watchers[watcher]; ok {
		return false
	}
	t.Watchers[watcher] = struct{}{}
	return true
}

func (t *tower) removeWatcher(watcher object.ObjectId) bool {
	if _, ok := t.Watchers[watcher]; !ok {
		return false
	}

	delete(t.Watchers, watcher)
	return true
}

func (t *tower) getAllWatchers() []object.ObjectId {
	if len(t.Watchers) == 0 {
		return nil
	}

	result := make([]object.ObjectId, 0, len(t.Watchers))

	for o := range t.Watchers {
		result = append(result, o)
	}
	return result
}

func NewTower() *tower {
	t := &tower{}
	t.Ids = IdSet{}
	t.Watchers = IdSet{}
	return t
}

type towerPos struct {
	X, Y int
}

func (tp towerPos) equal(rhs towerPos) bool {
	return tp.X == rhs.X && tp.Y == rhs.Y
}

type TowerAOI struct {
	width       float32
	height      float32
	towerWidth  float32
	towerHeight float32
	rangeLimit  int
	max         towerPos
	towers      [][]*tower
	callback    Callback
}

func NewTowerAOI(w float32, h float32, tw float32, th float32, limit int, cb Callback) *TowerAOI {
	aoi := &TowerAOI{
		width:       w,
		height:      h,
		towerWidth:  tw,
		towerHeight: th,
		rangeLimit:  limit,
		callback:    cb,
	}
	aoi.init()
	return aoi
}

func (aoi *TowerAOI) Enter(obj object.ObjectId, pos object.Vec3, ranges int) {
	if !aoi.checkPos(pos) {
		panic("pos invalid")
	}
	aoi.addWatcher(obj, pos, ranges)
	aoi.addObject(pos, obj)
}

func (aoi *TowerAOI) Level(obj object.ObjectId, pos object.Vec3, ranges int) {
	if !aoi.checkPos(pos) {
		panic("pos invalid")
	}
	if aoi.removeObject(pos, obj) {
		return
	}
	aoi.removeWatcher(obj, pos, ranges)
}

func (aoi *TowerAOI) Move(obj object.ObjectId, oldpos object.Vec3, dest object.Vec3, ranges int) bool {
	if !aoi.checkPos(oldpos) || !aoi.checkPos(dest) {
		return false
	}
	p1 := aoi.transPos(oldpos)
	p2 := aoi.transPos(dest)

	if p1.equal(p2) {
		return true
	}

	t1 := aoi.towers[p1.X][p1.Y]
	t2 := aoi.towers[p2.X][p2.Y]

	aoi.innerRemoveObject(t1, obj)
	aoi.innerAddObj(t2, obj)
	addTowers, removeTowers := aoi.getChangedTowers(p1, p2, ranges, ranges)
	for _, t := range removeTowers {
		aoi.innerRemoveWatcher(t, obj)
	}
	for _, t := range addTowers {
		aoi.innerAddWatch(t, obj)
	}
	return true
}

func (aoi *TowerAOI) Clear() {
	for i := 0; i <= aoi.max.X; i++ {
		for j := 0; j <= aoi.max.Y; j++ {
			aoi.towers[i][j].clear()
		}
	}
}

func (aoi *TowerAOI) GetIdsByRange(pos object.Vec3, ranges int) []object.ObjectId {
	if !aoi.checkPos(pos) || ranges < 0 {
		return nil
	}

	result := make([]object.ObjectId, 0, 100)
	if ranges > aoi.rangeLimit {
		ranges = aoi.rangeLimit
	}
	p := aoi.transPos(pos)
	start, end := getPosLimit(p, ranges, aoi.max)

	for i := start.X; i <= end.X; i++ {
		for j := start.Y; j <= end.Y; j++ {
			result = append(result, aoi.towers[i][j].getIds()...)
		}
	}
	return result
}

// Check if the pos is valid;
func (aoi *TowerAOI) checkPos(pos object.Vec3) bool {
	if pos.X() < 0 || pos.Z() < 0 || pos.X() >= aoi.width || pos.Z() >= aoi.height {
		return false
	}
	return true
}

// Trans the absolut pos to tower pos. For example : (210, 110} -> (1, 0), for tower width 200, height 200
func (aoi *TowerAOI) transPos(pos object.Vec3) towerPos {
	return towerPos{
		X: int(math.Floor(float64(pos.X() / aoi.towerWidth))),
		Y: int(math.Floor(float64(pos.Z() / aoi.towerHeight))),
	}
}

// getPosLimit Get the postion limit of given range,
func getPosLimit(pos towerPos, ranges int, max towerPos) (start towerPos, end towerPos) {

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
func isInRect(pos towerPos, start towerPos, end towerPos) bool {
	return (pos.X >= start.X && pos.X <= end.X && pos.Y >= start.Y && pos.Y <= end.Y)
}

func (aoi *TowerAOI) init() {
	iloop := int(math.Ceil(float64(aoi.width / aoi.towerWidth)))
	jloop := int(math.Ceil(float64(aoi.height / aoi.towerHeight)))
	aoi.max.X = iloop - 1
	aoi.max.Y = jloop - 1
	aoi.towers = make([][]*tower, iloop)
	for i := 0; i < iloop; i++ {
		aoi.towers[i] = make([]*tower, jloop)
		for j := 0; j < jloop; j++ {
			aoi.towers[i][j] = NewTower()
		}
	}
}

func (aoi *TowerAOI) addObject(pos object.Vec3, obj object.ObjectId) bool {
	p := aoi.transPos(pos)
	return aoi.innerAddObj(aoi.towers[p.X][p.Y], obj)
}

func (aoi *TowerAOI) innerAddObj(t *tower, obj object.ObjectId) bool {
	if t.add(obj) {
		for _, watcher := range t.getAllWatchers() {
			if watcher == obj {
				continue
			}
			aoi.callback.OnEnterAOI(watcher, obj)
		}
		return true
	}
	return false
}

func (aoi *TowerAOI) removeObject(pos object.Vec3, obj object.ObjectId) bool {
	p := aoi.transPos(pos)
	return aoi.innerRemoveObject(aoi.towers[p.X][p.Y], obj)
}

func (aoi *TowerAOI) innerRemoveObject(t *tower, obj object.ObjectId) bool {
	if t.remove(obj) {
		for _, watcher := range t.getAllWatchers() {
			if watcher == obj {
				continue
			}
			aoi.callback.OnLeaveAOI(watcher, obj)
		}
		return true
	}
	return false
}

func (aoi *TowerAOI) getWatchers(pos object.Vec3) []object.ObjectId {
	if aoi.checkPos(pos) {
		p := aoi.transPos(pos)
		return aoi.towers[p.X][p.Y].getAllWatchers()
	}
	return nil
}

func (aoi *TowerAOI) addWatcher(watcher object.ObjectId, pos object.Vec3, ranges int) {
	if ranges <= 0 {
		panic("ranges <= 0")
	}
	if ranges > aoi.rangeLimit {
		ranges = aoi.rangeLimit
	}
	p := aoi.transPos(pos)
	start, end := getPosLimit(p, ranges, aoi.max)
	for i := start.X; i <= end.X; i++ {
		for j := start.Y; j <= end.Y; j++ {
			aoi.innerAddWatch(aoi.towers[i][j], watcher)
		}
	}
}

func (aoi *TowerAOI) innerAddWatch(t *tower, watcher object.ObjectId) {
	if t.addWatcher(watcher) {
		for neighbor := range t.Ids {
			if neighbor != watcher {
				aoi.callback.OnEnterAOI(watcher, neighbor)
			}
		}
	}
}

func (aoi *TowerAOI) removeWatcher(watcher object.ObjectId, pos object.Vec3, ranges int) {
	if ranges <= 0 {
		panic("ranges <= 0")
	}

	if ranges > aoi.rangeLimit {
		ranges = aoi.rangeLimit
	}

	p := aoi.transPos(pos)

	start, end := getPosLimit(p, ranges, aoi.max)

	for i := start.X; i <= end.X; i++ {
		for j := start.Y; j <= end.Y; j++ {
			aoi.innerRemoveWatcher(aoi.towers[i][j], watcher)
		}
	}
}

func (aoi *TowerAOI) innerRemoveWatcher(t *tower, watcher object.ObjectId) {
	if t.removeWatcher(watcher) {
		for neighbor := range t.Ids {
			if neighbor != watcher {
				aoi.callback.OnLeaveAOI(watcher, neighbor)
			}
		}
	}
}

func (aoi *TowerAOI) getChangedTowers(p1 towerPos, p2 towerPos, r1 int, r2 int) ([]*tower, []*tower) {
	var start1, end1 = getPosLimit(p1, r1, aoi.max)
	var start2, end2 = getPosLimit(p2, r2, aoi.max)

	removeTowers := make([]*tower, 0, 10)
	addTowers := make([]*tower, 0, 10)

	for i := start1.X; i <= end1.X; i++ {
		for j := start1.Y; j <= end1.Y; j++ {
			if !isInRect(towerPos{i, j}, start2, end2) {
				removeTowers = append(removeTowers, aoi.towers[i][j])
			}
		}
	}

	for i := start2.X; i <= end2.X; i++ {
		for j := start2.Y; j <= end2.Y; j++ {
			if !isInRect(towerPos{i, j}, start1, end1) {
				addTowers = append(addTowers, aoi.towers[i][j])
			}
		}
	}

	return addTowers, removeTowers
}
