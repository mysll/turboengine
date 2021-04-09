package object

import (
	"encoding/binary"
	"fmt"
	"io"
)

var Endian = binary.LittleEndian

const (
	TYPE_UNKNOWN = 0
	TYPE_INT     = 1
	TYPE_FLOAT   = 2
	TYPE_INT64   = 3
	TYPE_STRING  = 4
)

// 属性接口
type Attr interface {
	Flag() int
	SetFlag(f int)
	ClearFlag(f int)
	Name() string
	Index() int
	SetIndex(int)
	Type() int
	Write(stream io.Writer) (int, error)
	Read(reader io.Reader) (int, error)
}

type AttrHolder struct {
	name  string
	index int
	flag  int
}

func (h *AttrHolder) Type() int {
	return TYPE_UNKNOWN
}

func (h *AttrHolder) Flag() int {
	return h.flag
}

func (h *AttrHolder) SetFlag(f int) {
	h.flag |= f
}

func (h *AttrHolder) ClearFlag(f int) {
	h.flag &= ^f
}

func (h *AttrHolder) Name() string {
	return h.name
}

func (h *AttrHolder) Index() int {
	return h.index
}

func (h *AttrHolder) SetIndex(idx int) {
	h.index = idx
}

type NoneHolder struct {
	AttrHolder
}

func NewNoneHolder(name string) Attr {
	return &NoneHolder{
		AttrHolder: AttrHolder{name: name},
	}
}

func (i *NoneHolder) Write(stream io.Writer) (int, error) {
	return 0, nil
}

func (i *NoneHolder) Read(reader io.Reader) (int, error) {
	return 0, nil
}

type IntHolder struct {
	AttrHolder
	data int32
}

func NewIntHolder(name string) Attr {
	return &IntHolder{
		AttrHolder: AttrHolder{name: name},
		data:       0,
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

func (i *IntHolder) Write(stream io.Writer) (int, error) {
	err := binary.Write(stream, Endian, i.data)
	if err != nil {
		return 0, err
	}
	return 4, nil
}

func (i *IntHolder) Read(reader io.Reader) (int, error) {
	err := binary.Read(reader, Endian, &i.data)
	if err != nil {
		return 0, err
	}
	return 4, nil
}

type Int64Holder struct {
	AttrHolder
	data int64
}

func NewInt64Holder(name string) Attr {
	return &Int64Holder{
		AttrHolder: AttrHolder{name: name},
		data:       0,
	}
}

func (i *Int64Holder) Type() int {
	return TYPE_INT64
}

func (i *Int64Holder) SetData(data int64) {
	i.data = data
}

func (i *Int64Holder) Data() int64 {
	return i.data
}

func (i *Int64Holder) Write(stream io.Writer) (int, error) {
	err := binary.Write(stream, Endian, i.data)
	if err != nil {
		return 0, err
	}
	return 4, nil
}

func (i *Int64Holder) Read(reader io.Reader) (int, error) {
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

func NewFloatHolder(name string) Attr {
	return &FloatHolder{
		AttrHolder: AttrHolder{name: name},
		data:       0,
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

func (f *FloatHolder) Write(stream io.Writer) (int, error) {
	err := binary.Write(stream, Endian, f.data)
	if err != nil {
		return 0, err
	}
	return 4, nil
}

func (f *FloatHolder) Read(reader io.Reader) (int, error) {
	err := binary.Read(reader, Endian, &f.data)
	if err != nil {
		return 0, err
	}
	return 4, nil
}

type StringHolder struct {
	AttrHolder
	data string
}

func NewStringHolder(name string) Attr {
	return &StringHolder{
		AttrHolder: AttrHolder{name: name},
		data:       "",
	}
}

func (f *StringHolder) Type() int {
	return TYPE_STRING
}

func (f *StringHolder) SetData(data string) {
	f.data = data
}

func (f *StringHolder) Data() string {
	return f.data
}

func (f *StringHolder) Write(stream io.Writer) (int, error) {
	size := uint16(len(f.data))
	binary.Write(stream, Endian, size)
	err := binary.Write(stream, Endian, f.data)
	if err != nil {
		return 0, err
	}
	return int(size) + 2, nil
}

func (f *StringHolder) Read(reader io.Reader) (int, error) {
	var size uint16
	binary.Read(reader, Endian, &size)
	if size == 0 {
		f.data = ""
		return 0, nil
	}
	buf := make([]byte, size)
	n, err := reader.Read(buf)
	if err != nil {
		return 0, err
	}
	if size != uint16(n) {
		return 0, fmt.Errorf("size not match")
	}
	f.data = string(buf)
	return int(size) + 2, nil
}

func init() {
	typeToObject[TYPE_UNKNOWN] = NewNoneHolder
	typeToObject[TYPE_INT] = NewIntHolder
	typeToObject[TYPE_FLOAT] = NewFloatHolder
	typeToObject[TYPE_INT64] = NewInt64Holder
	typeToObject[TYPE_STRING] = NewStringHolder
}

// Create object with type
func Create(typ int, name string) Attr {
	if t, ok := typeToObject[typ]; ok {
		return t(name)
	}
	return nil
}
