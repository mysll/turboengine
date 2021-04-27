package object

import (
	. "turboengine/common/datatype"
)

type AOI interface {
	AddObj(obj ObjectId)
	RemoveObj(obj ObjectId)
	Clear()
}

type ViewMap map[ObjectId]struct{}

type ViewChange struct {
	News []GameObject
	Del  []GameObject
}

type View struct {
	around   []GameObject
	owner    GameObject
	interest ViewMap
}

func NewView(owner GameObject) *View {
	return &View{
		owner:    owner,
		around:   make([]GameObject, 30),
		interest: make(ViewMap),
	}
}

func (v *View) Clear() {
	v.interest = ViewMap{}
}

func (v *View) AddObj(obj ObjectId) {
	if _, ok := v.interest[obj]; ok {
		return
	}
	v.interest[obj] = struct{}{}
}

func (v *View) RemoveObj(obj ObjectId) {
	if _, ok := v.interest[obj]; ok {
		delete(v.interest, obj)
	}
}
