package object

import (
	"encoding/binary"
	"fmt"
	"io"
	"turboengine/common/utils"
)

var Endian = binary.LittleEndian

const (
	TYPE_UNKNOWN = 0
	TYPE_INT     = 1
	TYPE_INT64   = 2
	TYPE_FLOAT   = 3
	TYPE_FLOAT64 = 4
	TYPE_STRING  = 5
)

type OnChange func(int, interface{})

// 属性接口
type Attr interface {
	Flag() int
	SetFlag(f int)
	ClearFlag(f int)
	// 是否存在标志位
	HasFlag(flag int) bool
	Name() string
	Index() int
	SetIndex(int)
	Type() int
	Write(stream io.Writer) (int, error)
	Read(reader io.Reader) (int, error)
	// 数据变动
	Change(notify OnChange)
}

type AttrHolder struct {
	name   string
	index  int
	flag   int
	change OnChange
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

func (h *AttrHolder) HasFlag(flag int) bool {
	return h.flag&flag != 0
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

func (h *AttrHolder) Change(change OnChange) {
	h.change = change
}

type NoneHolder struct {
	AttrHolder
}

func NewNoneHolder(name string) *NoneHolder {
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

func NewIntHolder(name string) *IntHolder {
	return &IntHolder{
		AttrHolder: AttrHolder{name: name},
		data:       0,
	}
}

func (i *IntHolder) Type() int {
	return TYPE_INT
}

func (i *IntHolder) SetData(data int32) bool {
	if i.data == data {
		return false
	}
	old := i.data
	i.data = data
	if i.change != nil {
		i.change(i.index, old)
	}
	return i.data != old
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

func NewInt64Holder(name string) *Int64Holder {
	return &Int64Holder{
		AttrHolder: AttrHolder{name: name},
		data:       0,
	}
}

func (i *Int64Holder) Type() int {
	return TYPE_INT64
}

func (i *Int64Holder) SetData(data int64) bool {
	if i.data == data {
		return false
	}
	old := i.data
	i.data = data
	if i.change != nil {
		i.change(i.index, old)
	}
	return i.data != old
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

func NewFloatHolder(name string) *FloatHolder {
	return &FloatHolder{
		AttrHolder: AttrHolder{name: name},
		data:       0,
	}
}

func (f *FloatHolder) Type() int {
	return TYPE_FLOAT
}

func (f *FloatHolder) SetData(data float32) bool {
	if utils.IsEqual(float64(f.data), float64(data)) {
		return false
	}
	old := f.data
	f.data = data
	if f.change != nil {
		f.change(f.index, old)
	}
	return !utils.IsEqual(float64(f.data), float64(old))
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

type Float64Holder struct {
	AttrHolder
	data float64
}

func NewFloat64Holder(name string) *Float64Holder {
	return &Float64Holder{
		AttrHolder: AttrHolder{name: name},
		data:       0,
	}
}

func (f *Float64Holder) Type() int {
	return TYPE_FLOAT64
}

func (f *Float64Holder) SetData(data float64) bool {
	if utils.IsEqual(f.data, data) {
		return false
	}
	old := f.data
	f.data = data
	if f.change != nil {
		f.change(f.index, old)
	}
	return !utils.IsEqual(f.data, old)
}

func (f *Float64Holder) Data() float64 {
	return f.data
}

func (f *Float64Holder) Write(stream io.Writer) (int, error) {
	err := binary.Write(stream, Endian, f.data)
	if err != nil {
		return 0, err
	}
	return 4, nil
}

func (f *Float64Holder) Read(reader io.Reader) (int, error) {
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

func NewStringHolder(name string) *StringHolder {
	return &StringHolder{
		AttrHolder: AttrHolder{name: name},
		data:       "",
	}
}

func (s *StringHolder) Type() int {
	return TYPE_STRING
}

func (s *StringHolder) SetData(data string) bool {
	if s.data == data {
		return false
	}
	old := s.data
	s.data = data
	if s.change != nil {
		s.change(s.index, old)
	}
	return s.data != old
}

func (s *StringHolder) Data() string {
	return s.data
}

func (s *StringHolder) Write(stream io.Writer) (int, error) {
	size := uint16(len(s.data))
	binary.Write(stream, Endian, size)
	err := binary.Write(stream, Endian, s.data)
	if err != nil {
		return 0, err
	}
	return int(size) + 2, nil
}

func (s *StringHolder) Read(reader io.Reader) (int, error) {
	var size uint16
	binary.Read(reader, Endian, &size)
	if size == 0 {
		s.data = ""
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
	s.data = string(buf)
	return int(size) + 2, nil
}

func init() {
	typeToObject[TYPE_UNKNOWN] = func(name string) Attr { return NewNoneHolder(name) }
	typeToObject[TYPE_INT] = func(name string) Attr { return NewIntHolder(name) }
	typeToObject[TYPE_FLOAT] = func(name string) Attr { return NewFloatHolder(name) }
	typeToObject[TYPE_FLOAT64] = func(name string) Attr { return NewFloat64Holder(name) }
	typeToObject[TYPE_INT64] = func(name string) Attr { return NewInt64Holder(name) }
	typeToObject[TYPE_STRING] = func(name string) Attr { return NewStringHolder(name) }
}

// Create object with type
func Create(typ int, name string) Attr {
	if t, ok := typeToObject[typ]; ok {
		return t(name)
	}
	return nil
}
