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
	minX, minY, maxX, maxY float32
	width, height          float32
	towerWidth             float32
	towerHeight            float32
	rangeLimit             float32
	max                    towerPos
	towers                 [][]*tower
	callback               Callback
}

func NewTowerAOI(minx, miny, maxx, maxy float32, tw float32, th float32, limit float32, cb Callback) *TowerAOI {
	aoi := &TowerAOI{
		minX:        minx,
		minY:        miny,
		maxX:        maxx,
		maxY:        maxy,
		towerWidth:  tw,
		towerHeight: th,
		rangeLimit:  limit,
		callback:    cb,
	}
	aoi.init()
	return aoi
}

func (aoi *TowerAOI) Enter(obj object.ObjectId, pos object.Vec3, ranges float32) {
	if !aoi.checkPos(pos) {
		panic("pos invalid")
	}
	aoi.addWatcher(obj, pos, ranges)
	aoi.addObject(pos, obj)
}

func (aoi *TowerAOI) Level(obj object.ObjectId, pos object.Vec3, ranges float32) {
	if !aoi.checkPos(pos) {
		panic("pos invalid")
	}
	if aoi.removeObject(pos, obj) {
		return
	}
	aoi.removeWatcher(obj, pos, ranges)
}

func (aoi *TowerAOI) Move(obj object.ObjectId, oldpos object.Vec3, dest object.Vec3, ranges float32) bool {
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
	addTowers, removeTowers := aoi.getChangedTowers(oldpos, dest, ranges, ranges)
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

func (aoi *TowerAOI) GetIdsByRange(pos object.Vec3, ranges float32) []object.ObjectId {
	if !aoi.checkPos(pos) || ranges < 0 {
		return nil
	}

	result := make([]object.ObjectId, 0, 100)
	if ranges > aoi.rangeLimit {
		ranges = aoi.rangeLimit
	}
	min := aoi.transPos(object.Vec3{pos.X() - ranges, 0, pos.Z() - ranges})
	max := aoi.transPos(object.Vec3{pos.X() + ranges, 0, pos.Z() + ranges})
	for i := min.X; i <= max.X; i++ {
		for j := min.Y; j <= max.Y; j++ {
			result = append(result, aoi.towers[i][j].getIds()...)
		}
	}
	return result
}

// Check if the pos is valid;
func (aoi *TowerAOI) checkPos(pos object.Vec3) bool {
	if pos.X() < aoi.minX || pos.Z() < aoi.minY || pos.X() > aoi.maxX || pos.Z() > aoi.maxY {
		return false
	}
	return true
}

// Trans the absolut pos to tower pos. For example : (210, 110} -> (1, 0), for tower width 200, height 200
func (aoi *TowerAOI) transPos(pos object.Vec3) towerPos {
	tx, ty := int(math.Floor(float64((pos.X()-aoi.minX)/aoi.towerWidth))),
		int(math.Floor(float64((pos.Z()-aoi.minY)/aoi.towerHeight)))
	if tx < 0 {
		tx = 0
	} else if tx > aoi.max.X {
		tx = aoi.max.X
	}

	if ty < 0 {
		ty = 0
	} else if ty > aoi.max.Y {
		ty = aoi.max.Y
	}
	return towerPos{
		X: tx,
		Y: ty,
	}
}

func (aoi *TowerAOI) init() {
	aoi.width = aoi.maxX - aoi.minX + 1
	aoi.height = aoi.maxY - aoi.minY + 1
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

func (aoi *TowerAOI) addWatcher(watcher object.ObjectId, pos object.Vec3, ranges float32) {
	if ranges <= 0 {
		panic("ranges <= 0")
	}
	if ranges > aoi.rangeLimit {
		ranges = aoi.rangeLimit
	}
	min := aoi.transPos(object.Vec3{pos.X() - ranges, 0, pos.Z() - ranges})
	max := aoi.transPos(object.Vec3{pos.X() + ranges, 0, pos.Z() + ranges})
	for i := min.X; i <= max.X; i++ {
		for j := min.Y; j <= max.Y; j++ {
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

func (aoi *TowerAOI) removeWatcher(watcher object.ObjectId, pos object.Vec3, ranges float32) {
	if ranges <= 0 {
		panic("ranges <= 0")
	}

	if ranges > aoi.rangeLimit {
		ranges = aoi.rangeLimit
	}

	min := aoi.transPos(object.Vec3{pos.X() - ranges, 0, pos.Z() - ranges})
	max := aoi.transPos(object.Vec3{pos.X() + ranges, 0, pos.Z() + ranges})
	for i := min.X; i <= max.X; i++ {
		for j := min.Y; j <= max.Y; j++ {
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

func (aoi *TowerAOI) getChangedTowers(p1, p2 object.Vec3, r1 float32, r2 float32) ([]*tower, []*tower) {
	oldmin := aoi.transPos(object.Vec3{p1.X() - r1, 0, p1.Z() - r1})
	oldmax := aoi.transPos(object.Vec3{p1.X() + r1, 0, p1.Z() + r1})
	destmin := aoi.transPos(object.Vec3{p2.X() - r2, 0, p2.Z() - r2})
	destmax := aoi.transPos(object.Vec3{p2.X() + r2, 0, p2.Z() + r2})
	removeTowers := make([]*tower, 0, 10)
	addTowers := make([]*tower, 0, 10)

	for x := oldmin.X; x <= oldmax.X; x++ {
		for y := oldmin.Y; y <= oldmax.Y; y++ {
			if x >= destmin.X && x <= destmax.X && y >= destmin.Y && y <= destmax.Y {
				continue
			}
			removeTowers = append(removeTowers, aoi.towers[x][y])
		}
	}

	for x := destmin.X; x <= destmax.X; x++ {
		for y := destmin.Y; y <= destmax.Y; y++ {
			if x >= oldmin.X && x <= oldmax.X && y >= oldmin.Y && y <= oldmax.Y {
				continue
			}
			addTowers = append(addTowers, aoi.towers[x][y])
		}
	}

	return addTowers, removeTowers
}
