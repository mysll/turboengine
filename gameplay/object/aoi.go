package object

import (
	. "turboengine/common/datatype"
)

type AOI interface {
	ViewRange() float32
	SetViewRange(f float32)
	AddViewObj(obj ObjectId)
	RemoveViewObj(obj ObjectId)
	Clear()
	CachePosition(pos Vec3)
	Position() Vec3
}

type ViewMap map[ObjectId]struct{}

type ViewChange struct {
	News []GameObject
	Del  []GameObject
}

type View struct {
	around    []GameObject
	owner     GameObject
	neighbor  ViewMap
	viewRange float32
	position  Vec3
}

func NewView(owner GameObject) *View {
	return &View{
		owner:    owner,
		around:   make([]GameObject, 30),
		neighbor: make(ViewMap),
	}
}

func (v *View) Clear() {
	v.neighbor = ViewMap{}
}

func (v *View) AddViewObj(obj ObjectId) {
	if _, ok := v.neighbor[obj]; ok {
		return
	}
	v.neighbor[obj] = struct{}{}
}

func (v *View) RemoveViewObj(obj ObjectId) {
	if _, ok := v.neighbor[obj]; ok {
		delete(v.neighbor, obj)
	}
}

func (v *View) ViewRange() float32 {
	return v.viewRange
}

func (v *View) SetViewRange(f float32) {
	v.viewRange = f
}

func (v *View) CachePosition(pos Vec3) {
	v.position = pos
}

func (v *View) Position() Vec3 {
	return v.position
}
