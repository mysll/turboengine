package object

import (
	"encoding/binary"
	"io"
)

var Endian = binary.LittleEndian

const (
	TYPE_UNKNOWN = 0
	TYPE_INT     = 1
	TYPE_FLOAT   = 2
)

// 属性接口
type Attr interface {
	Flag() uint32
	SetFlag(f uint32)
	ClearFlag(f uint32)
	Name() string
	Index() uint32
	Type() int
	Write(stream io.Writer) (uint32, error)
	Read(reader io.Reader) (uint32, error)
}

type AttrHolder struct {
	name  string
	index uint32
	flag  uint32
}

func (h *AttrHolder) Type() int {
	return TYPE_UNKNOWN
}

func (h *AttrHolder) Flag() uint32 {
	return h.flag
}

func (h *AttrHolder) SetFlag(f uint32) {
	h.flag |= f
}

func (h *AttrHolder) ClearFlag(f uint32) {
	h.flag &= ^f
}

func (h *AttrHolder) Name() string {
	return h.name
}

func (h *AttrHolder) Index() uint32 {
	return h.index
}

type IntHolder struct {
	AttrHolder
	data int32
}

func NewIntHolder(name string, index uint32, data int32) *IntHolder {
	return &IntHolder{
		AttrHolder: AttrHolder{name, index, 0},
		data:       data,
	}
}

func (i *IntHolder) Type() int {
	return TYPE_INT
}

func (i *IntHolder) SetData(data int32) {
	i.data = data
}

func (i *IntHolder) Data() int32 {
	return i.data
}

func (i *IntHolder) Write(stream io.Writer) (uint32, error) {
	err := binary.Write(stream, Endian, i.data)
	if err != nil {
		return 0, err
	}
	return 4, nil
}

func (i *IntHolder) Read(reader io.Reader) (uint32, error) {
	err := binary.Read(reader, Endian, &i.data)
	if err != nil {
		return 0, err
	}
	return 4, nil
}

type FloatHolder struct {
	AttrHolder
	data float32
}

func NewFloatHolder(name string, index uint32, data float32) *FloatHolder {
	return &FloatHolder{
		AttrHolder: AttrHolder{name, index, 0},
		data:       data,
	}
}

func (f *FloatHolder) Type() int {
	return TYPE_FLOAT
}

func (f *FloatHolder) SetData(data float32) {
	f.data = data
}

func (f *FloatHolder) Data() float32 {
	return f.data
}

func (f *FloatHolder) Write(stream io.Writer) (uint32, error) {
	err := binary.Write(stream, Endian, f.data)
	if err != nil {
		return 0, err
	}
	return 4, nil
}

func (f *FloatHolder) Read(reader io.Reader) (uint32, error) {
	err := binary.Read(reader, Endian, &f.data)
	if err != nil {
		return 0, err
	}
	return 4, nil
}
