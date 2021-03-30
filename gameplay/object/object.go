package object

import (
	"fmt"
)

var typeToObject = make(map[int]func(string) Attr)

// 基础对象，所以游戏内的对象基类
type Object struct {
	attrs     []Attr
	nameToIdx map[string]int
}

func (o *Object) AddAttr(attr Attr) error {
	if _, dup := o.nameToIdx[attr.Name()]; dup {
		return fmt.Errorf("attr %s already exist", attr.Name())
	}
	idx := len(o.attrs)
	attr.SetIndex(idx)
	o.attrs = append(o.attrs, attr)
	o.nameToIdx[attr.Name()] = idx
	return nil
}

func (o *Object) AddAttrByType(name string, typ int) error {
	fn, ok := typeToObject[typ]
	if !ok {
		return fmt.Errorf("attr type not found, %s %d", name, typ)
	}
	attr := fn(name)
	return o.AddAttr(attr)
}
