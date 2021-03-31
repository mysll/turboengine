package object

import (
	"fmt"
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
}

func (o *Object) Id() ObjectId {
	return o.id
}

func (o *Object) Dirty() bool {
	return o.dirty
}

func (o *Object) ClearDirty() {
	o.dirty = false
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
	if index > 0 && index < len(o.attrs) {
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
