package object

import (
	"fmt"
)

const (
	OBJECT_NONE     = 0
	OBJECT_SAVE     = 1
	OBJECT_PUBLIC   = 1 << 1
	OBJECT_PRIVATE  = 1 << 2
	OBJECT_REALTIME = 1 << 3
	OBJECT_CHANGE   = 1 << 4
	OBJECT_CHANGING = 1 << 5
)

var typeToObject = make(map[int]func(string) Attr)

type ObjectId uint64

// 基础对象，所以游戏内的对象基类
type Object struct {
	id        ObjectId
	attrs     []Attr
	nameToIdx map[string]int
	replicate bool
	dirty     bool
	silent    bool // 静默
	inited    bool
	change    []OnChange
}

func (o *Object) New(cap int) {
	o.attrs = make([]Attr, 0, cap)
	o.nameToIdx = make(map[string]int, cap)
	o.change = make([]OnChange, cap)
}

func (o *Object) Init() {
	if o.inited {
		return
	}
	o.inited = true
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
		o.change[index] = change
		o.attrs[index].SetFlag(OBJECT_CHANGE)
	}
}

func (o *Object) ClearChange(index int) {
	if index >= 0 && index < len(o.attrs) {
		o.attrs[index].Change(nil)
		o.change[index] = nil
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
		if o.change[index] != nil {
			o.change[index](index, val)
		}
		o.attrs[index].ClearFlag(OBJECT_CHANGING)
	}
}
