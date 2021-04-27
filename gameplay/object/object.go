package object

import (
	"fmt"
	. "turboengine/common/datatype"
)

const (
	OBJECT_NONE          = 0
	OBJECT_SAVE          = 1
	OBJECT_PUBLIC        = 1 << 1
	OBJECT_PRIVATE       = 1 << 2
	OBJECT_REALTIME      = 1 << 3
	OBJECT_CHANGE        = 1 << 4
	OBJECT_CHANGING      = 1 << 5
	OBJECT_DIRTY         = 1 << 6
	OBJECT_PUBLIC_DIRTY  = 1 << 7
	OBJECT_PRIVATE_DIRTY = 1 << 8
)

const (
	FEATURES_NONE      = 0
	FEATURES_MOVEMENT  = 1
	FEATURES_REPLICATE = 1 << 1
	FEATURES_AOI       = 1<<2 | FEATURES_MOVEMENT
	FEATURES_COLLIDER  = 1<<3 | FEATURES_MOVEMENT
	FEATURES_ALL       = -1
)

var typeToObject = make(map[int]func(string) Attr)

type GameObject interface {
	Id() ObjectId
	Dirty() bool
	SetDirty()
	ClearDirty()
	PublicDirty() bool
	SetPublicDirty()
	ClearPublicDirty()
	PrivateDirty() bool
	SetPrivateDirty()
	ClearPrivateDirty()
	Silent() bool
	SetSilent(s bool)
	AttrCount() int
	GetAttr(index int) Attr
	GetAttrByName(name string) Attr
	IsMovement() bool
	IsReplicate() bool
	HasView() bool
	Movement() Movement
	Collider() Collider
	AOI() AOI
}

// 基础对象，所以游戏内的对象基类
type Object struct {
	*Collision
	*Transform
	*Replication
	*View
	id        ObjectId
	attrs     []Attr
	nameToIdx map[string]int
	dirty     bool
	silent    bool // 静默
	inited    bool
	change    *changeEvent
	pubDirty  bool
	priDirty  bool
	holder    GameObject
	features  int
}

func (o *Object) SetFeature(features int) {
	o.features = features
	if o.features|FEATURES_MOVEMENT != 0 {
		o.Transform = NewTransform(o.holder)
	} else {
		o.Transform = nil
	}

	if o.features|FEATURES_REPLICATE != 0 {
		o.Replication = NewReplication(o.holder)
	} else {
		o.Replication = nil
	}

	if o.features|FEATURES_AOI != 0 {
		o.View = NewView(o.holder)
	} else {
		o.View = nil
	}

	if o.features|FEATURES_COLLIDER != 0 {
		o.Collision = NewCollision(o.holder)
	} else {
		o.Collision = nil
	}
}

func (o *Object) IsMovement() bool {
	return o.features|FEATURES_MOVEMENT != 0
}

func (o *Object) IsReplicate() bool {
	return o.features|FEATURES_REPLICATE != 0
}

func (o *Object) HasView() bool {
	return o.features|FEATURES_AOI != 0
}

func (o *Object) hasCollider() bool {
	return o.features|FEATURES_COLLIDER != 0
}

func (o *Object) Movement() Movement {
	return o.Transform
}

func (o *Object) Collider() Collider {
	return o.Collision
}

func (o *Object) AOI() AOI {
	return o.View
}

func (o *Object) new(cap int) {
	o.attrs = make([]Attr, 0, cap)
	o.nameToIdx = make(map[string]int, cap)
	o.change = newChangeEvent(cap)
}

func (o *Object) InitOnce(self GameObject, cap int) {
	if o.inited {
		return
	}
	o.new(cap)
	o.holder = self
	o.inited = true
}

func (o *Object) SetOwner(self GameObject) {
	o.holder = self
}

func (o *Object) Id() ObjectId {
	return o.id
}

func (o *Object) Dirty() bool {
	return o.dirty
}

func (o *Object) SetDirty() {
	o.dirty = true
}

func (o *Object) ClearDirty() {
	o.dirty = false
}

func (o *Object) PublicDirty() bool {
	return o.pubDirty
}

func (o *Object) SetPublicDirty() {
	o.pubDirty = true
}

func (o *Object) ClearPublicDirty() {
	o.pubDirty = false
}

func (o *Object) PrivateDirty() bool {
	return o.priDirty
}

func (o *Object) SetPrivateDirty() {
	o.priDirty = true
}

func (o *Object) ClearPrivateDirty() {
	o.priDirty = false
}

func (o *Object) Silent() bool {
	return o.silent
}

func (o *Object) SetSilent(s bool) {
	o.silent = s
}

func (o *Object) AttrCount() int {
	return len(o.attrs)
}

func (o *Object) AddAttr(attr Attr) (int, error) {
	if _, dup := o.nameToIdx[attr.Name()]; dup {
		return -1, fmt.Errorf("attr %s already exist", attr.Name())
	}
	idx := len(o.attrs)
	attr.SetIndex(idx)
	o.attrs = append(o.attrs, attr)
	o.nameToIdx[attr.Name()] = idx
	return idx, nil
}

func (o *Object) AddAttrByType(name string, typ int) (int, error) {
	fn, ok := typeToObject[typ]
	if !ok {
		return -1, fmt.Errorf("attr type not found, %s %d", name, typ)
	}
	attr := fn(name)
	return o.AddAttr(attr)
}

func (o *Object) GetAttr(index int) Attr {
	if index >= 0 && index < len(o.attrs) {
		return o.attrs[index]
	}
	return nil
}

func (o *Object) GetAttrByName(name string) Attr {
	if index, ok := o.nameToIdx[name]; ok {
		return o.attrs[index]
	}
	return nil
}

func (o *Object) Change(index int, change OnChange) {
	if index >= 0 && index < len(o.attrs) {
		o.attrs[index].Change(o.onChange)
		o.change.add(index, change)
		o.attrs[index].SetFlag(OBJECT_CHANGE)
	}
}

func (o *Object) RemoveChange(index int, change OnChange) {
	if index >= 0 && index < len(o.attrs) {
		o.change.remove(index, change)
		if !o.change.hasEvent(index) {
			o.attrs[index].ClearFlag(OBJECT_CHANGE)
		}
	}
}

func (o *Object) ClearChange(index int) {
	if index >= 0 && index < len(o.attrs) {
		o.attrs[index].Change(nil)
		o.change.clear(index)
		o.attrs[index].ClearFlag(OBJECT_CHANGE)
	}
}

func (o *Object) onChange(index int, val interface{}) {
	if o.silent {
		return
	}
	if index >= 0 && index < len(o.attrs) {
		if o.attrs[index].HasFlag(OBJECT_CHANGING) {
			return
		}
		o.attrs[index].SetFlag(OBJECT_CHANGING)
		o.change.emit(o.holder, index, val)
		o.attrs[index].ClearFlag(OBJECT_CHANGING)
		if o.IsReplicate() {
			o.Replication.change(index, o.attrs[index])
		}
	}
}
