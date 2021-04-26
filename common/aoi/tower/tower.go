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

func (tp TowerPos) Equal(rhs TowerPos) bool {
	return tp.X == rhs.X && tp.Y == rhs.Y
}

type TowerAOI struct {
	width       float32
	height      float32
	towerWidth  float32
	towerHeight float32
	rangeLimit  int
	max         TowerPos
	towers      [][]*Tower
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

func (aoi *TowerAOI) Enter(obj object.ObjectId, pos object.Vec3, ranges int) bool {
	if !aoi.addWatcher(obj, pos, ranges) {
		return false
	}
	return aoi.addObject(pos, obj)
}

func (aoi *TowerAOI) Level(obj object.ObjectId, pos object.Vec3, ranges int) bool {
	if !aoi.removeObject(pos, obj) {
		return false
	}
	return aoi.removeWatcher(obj, pos, ranges)
}

func (aoi *TowerAOI) Move(obj object.ObjectId, oldpos object.Vec3, dest object.Vec3, ranges int) bool {
	if !aoi.checkPos(oldpos) || !aoi.checkPos(dest) {
		return false
	}
	p1 := aoi.transPos(oldpos)
	p2 := aoi.transPos(dest)

	if p1.Equal(p2) {
		return true
	}

	return false
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

func (this *TowerAOI) init() {
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

func (this *TowerAOI) GetIdsByPos(pos object.Vec3, ranges int) []object.ObjectId {
	if !this.checkPos(pos) || ranges < 0 {
		return nil
	}

	result := make([]object.ObjectId, 0, 100)
	if ranges > this.rangeLimit {
		ranges = this.rangeLimit
	}
	p := this.transPos(pos)
	start, end := getPosLimit(p, ranges, this.max)

	for i := start.X; i <= end.X; i++ {
		for j := start.Y; j <= end.Y; j++ {
			result = append(result, this.towers[i][j].getIds()...)
		}
	}
	return result
}

func (this *TowerAOI) addObject(pos object.Vec3, obj object.ObjectId) bool {
	if this.checkPos(pos) {
		p := this.transPos(pos)
		return this.innerAddObj(this.towers[p.X][p.Y], obj)
	}
	return false
}

func (this *TowerAOI) innerAddObj(t *Tower, obj object.ObjectId) bool {
	if t.add(obj) {
		for _, watcher := range t.getAllWatchers() {
			if watcher == obj {
				continue
			}
			this.callback.OnEnterAOI(watcher, obj)
		}
		return true
	}
	return false
}

func (this *TowerAOI) removeObject(pos object.Vec3, obj object.ObjectId) bool {
	if this.checkPos(pos) {
		p := this.transPos(pos)
		return this.innerRemoveObject(this.towers[p.X][p.Y], obj)
	}
	return false
}

func (this *TowerAOI) innerRemoveObject(t *Tower, obj object.ObjectId) bool {
	if t.remove(obj) {
		for _, watcher := range t.getAllWatchers() {
			if watcher == obj {
				continue
			}
			this.callback.OnLeaveAOI(watcher, obj)
		}
		return true
	}
	return false
}

func (this *TowerAOI) getWatchers(pos object.Vec3) []object.ObjectId {
	if this.checkPos(pos) {
		p := this.transPos(pos)
		return this.towers[p.X][p.Y].getAllWatchers()
	}
	return nil
}

func (this *TowerAOI) addWatcher(watcher object.ObjectId, pos object.Vec3, ranges int) bool {
	if ranges < 0 {
		return false
	}
	if ranges > this.rangeLimit {
		ranges = this.rangeLimit
	}
	p := this.transPos(pos)
	start, end := getPosLimit(p, ranges, this.max)
	for i := start.X; i <= end.X; i++ {
		for j := start.Y; j <= end.Y; j++ {
			this.towers[i][j].addWatcher(watcher)
			for neighbor := range this.towers[i][j].Ids {
				if neighbor != watcher {
					this.callback.OnEnterAOI(watcher, neighbor)
				}
			}
		}
	}

	return true
}

func (this *TowerAOI) removeWatcher(watcher object.ObjectId, pos object.Vec3, ranges int) bool {
	if ranges < 0 {
		return false
	}

	if ranges > this.rangeLimit {
		ranges = this.rangeLimit
	}

	p := this.transPos(pos)

	start, end := getPosLimit(p, ranges, this.max)

	for i := start.X; i <= end.X; i++ {
		for j := start.Y; j <= end.Y; j++ {
			this.towers[i][j].removeWatcher(watcher)
			for neighbor := range this.towers[i][j].Ids {
				if neighbor != watcher {
					this.callback.OnLeaveAOI(watcher, neighbor)
				}
			}
		}
	}

	return true
}

func (this *TowerAOI) updateWatcher(watcher object.ObjectId, oldPos object.Vec3, newPos object.Vec3, oldRange int, newRange int) bool {
	if !this.checkPos(oldPos) || !this.checkPos(newPos) {
		return false
	}
	p1 := this.transPos(oldPos)
	p2 := this.transPos(newPos)

	if p1.X == p2.X && p1.Y == p2.Y {
		return true
	} else {
		if oldRange < 0 || newRange < 0 {
			return false
		}
		if oldRange > this.rangeLimit {
			oldRange = this.rangeLimit
		}
		if newRange > this.rangeLimit {
			newRange = this.rangeLimit
		}
		addTowers, removeTowers := this.getChangedTowers(p1, p2, oldRange, newRange)
		for _, t := range removeTowers {
			t.removeWatcher(watcher)
			t.remove(watcher)
			for _, neighbor := range t.getAllWatchers() {
				if neighbor != watcher {
					this.callback.OnLeaveAOI(neighbor, watcher)
				}
			}
		}

		for _, t := range addTowers {
			t.addWatcher(watcher)
			t.add(watcher)
			for _, neighbor := range t.getAllWatchers() {
				if neighbor != watcher {
					this.callback.OnEnterAOI(neighbor, watcher)
				}
			}
		}

	}
	return true
}

func (this *TowerAOI) getChangedTowers(p1 TowerPos, p2 TowerPos, r1 int, r2 int) ([]*Tower, []*Tower) {
	var start1, end1 = getPosLimit(p1, r1, this.max)
	var start2, end2 = getPosLimit(p2, r2, this.max)

	removeTowers := make([]*Tower, 0, 10)
	addTowers := make([]*Tower, 0, 10)

	for i := start1.X; i <= end1.X; i++ {
		for j := start1.Y; j <= end1.Y; j++ {
			if !isInRect(TowerPos{i, j}, start2, end2) {
				removeTowers = append(removeTowers, this.towers[i][j])
			}
		}
	}

	for i := start2.X; i <= end2.X; i++ {
		for j := start2.Y; j <= end2.Y; j++ {
			if !isInRect(TowerPos{i, j}, start1, end1) {
				addTowers = append(addTowers, this.towers[i][j])
			}
		}
	}

	return addTowers, removeTowers
}
