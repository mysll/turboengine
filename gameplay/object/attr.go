package object

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	. "turboengine/common/datatype"

	"github.com/mysll/toolkit"
)

var Endian = binary.LittleEndian

const (
	TYPE_UNKNOWN = 0
	TYPE_INT     = 1
	TYPE_INT64   = 2
	TYPE_FLOAT   = 3
	TYPE_FLOAT64 = 4
	TYPE_STRING  = 5
	TYPE_VECTOR2 = 6
	TYPE_VECTOR3 = 7
	TYPE_BYTES   = 8
)

type changeFn func(int, interface{})

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
	Change(change changeFn)
	Equal(Attr) bool
}

type AttrHolder struct {
	name   string
	index  int
	flag   int
	change changeFn
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

func (h *AttrHolder) Change(change changeFn) {
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

func (i *NoneHolder) Equal(other Attr) bool {
	if other.Type() == i.Type() {
		return true
	}
	return false
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

func (i *IntHolder) Equal(other Attr) bool {
	if other.Type() == i.Type() {
		if o, ok := other.(*IntHolder); ok {
			return i.data == o.data
		}
	}
	return false
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

func (i *Int64Holder) Equal(other Attr) bool {
	if other.Type() == i.Type() {
		if o, ok := other.(*Int64Holder); ok {
			return i.data == o.data
		}
	}
	return false
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
	if toolkit.IsEqual32(f.data, data) {
		return false
	}
	old := f.data
	f.data = data
	if f.change != nil {
		f.change(f.index, old)
	}
	return !toolkit.IsEqual32(f.data, old)
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

func (f *FloatHolder) Equal(other Attr) bool {
	if other.Type() == f.Type() {
		if o, ok := other.(*FloatHolder); ok {
			return toolkit.IsEqual32(f.data, o.data)
		}
	}
	return false
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
	if toolkit.IsEqual64(f.data, data) {
		return false
	}
	old := f.data
	f.data = data
	if f.change != nil {
		f.change(f.index, old)
	}
	return !toolkit.IsEqual64(f.data, old)
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

func (f *Float64Holder) Equal(other Attr) bool {
	if other.Type() == f.Type() {
		if o, ok := other.(*Float64Holder); ok {
			return toolkit.IsEqual64(f.data, o.data)
		}
	}
	return false
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

func (s *StringHolder) Equal(other Attr) bool {
	if other.Type() == s.Type() {
		if o, ok := other.(*StringHolder); ok {
			return s.data == o.data
		}
	}
	return false
}

type Vector2Holder struct {
	AttrHolder
	data Vec2
}

func NewVector2Holder(name string) *Vector2Holder {
	return &Vector2Holder{
		AttrHolder: AttrHolder{name: name},
	}
}

func (v *Vector2Holder) Type() int {
	return TYPE_VECTOR2
}

func (v *Vector2Holder) SetData(val Vec2) bool {
	if v.data.Equal(val) {
		return false
	}
	old := v.data
	v.data = val

	if v.change != nil {
		v.change(v.index, old)
	}

	if v.data.Equal(old) {
		return false
	}
	return true
}

func (v *Vector2Holder) Data() Vec2 {
	return v.data
}

func (v *Vector2Holder) Write(stream io.Writer) (int, error) {
	err := binary.Write(stream, Endian, v.data)
	if err != nil {
		return 0, err
	}
	return 4, nil
}

func (v *Vector2Holder) Read(reader io.Reader) (int, error) {
	err := binary.Read(reader, Endian, &v.data)
	if err != nil {
		return 0, err
	}
	return 4, nil
}

func (v *Vector2Holder) Equal(other Attr) bool {
	if other.Type() == v.Type() {
		if o, ok := other.(*Vector2Holder); ok {
			return v.data.Equal(o.data)
		}
	}
	return false
}

type Vector3Holder struct {
	AttrHolder
	data Vec3
}

func NewVector3Holder(name string) *Vector3Holder {
	return &Vector3Holder{
		AttrHolder: AttrHolder{name: name},
	}
}

func (v *Vector3Holder) Type() int {
	return TYPE_VECTOR3
}

func (v *Vector3Holder) SetData(val Vec3) bool {
	if v.data.Equal(val) {
		return false
	}
	old := v.data
	v.data = val

	if v.change != nil {
		v.change(v.index, old)
	}

	if v.data.Equal(old) {
		return false
	}
	return true
}

func (v *Vector3Holder) Data() Vec3 {
	return v.data
}

func (v *Vector3Holder) Write(stream io.Writer) (int, error) {
	err := binary.Write(stream, Endian, v.data)
	if err != nil {
		return 0, err
	}
	return 4, nil
}

func (v *Vector3Holder) Read(reader io.Reader) (int, error) {
	err := binary.Read(reader, Endian, &v.data)
	if err != nil {
		return 0, err
	}
	return 4, nil
}

func (v *Vector3Holder) Equal(other Attr) bool {
	if other.Type() == v.Type() {
		if o, ok := other.(*Vector3Holder); ok {
			return v.data.Equal(o.data)
		}
	}
	return false
}

type BytesHolder struct {
	AttrHolder
	data []byte
}

func NewBytesHolder(name string) *BytesHolder {
	return &BytesHolder{
		AttrHolder: AttrHolder{name: name},
	}
}

func (s *BytesHolder) Type() int {
	return TYPE_BYTES
}

func (s *BytesHolder) SetData(data []byte) bool {
	if bytes.Equal(s.data, data) {
		return false
	}
	old := s.data
	s.data = data
	if s.change != nil {
		s.change(s.index, old)
	}
	return !bytes.Equal(s.data, old)
}

func (s *BytesHolder) Data() []byte {
	return s.data
}

func (s *BytesHolder) Write(stream io.Writer) (int, error) {
	size := uint16(len(s.data))
	binary.Write(stream, Endian, size)
	err := binary.Write(stream, Endian, s.data)
	if err != nil {
		return 0, err
	}
	return int(size) + 2, nil
}

func (s *BytesHolder) Read(reader io.Reader) (int, error) {
	var size uint16
	binary.Read(reader, Endian, &size)
	if size == 0 {
		s.data = nil
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
	s.data = buf
	return int(size) + 2, nil
}

func (s *BytesHolder) Equal(other Attr) bool {
	if other.Type() == s.Type() {
		if o, ok := other.(*BytesHolder); ok {
			return bytes.Equal(s.data, o.data)
		}
	}
	return false
}

func init() {
	typeToObject[TYPE_UNKNOWN] = func(name string) Attr { return NewNoneHolder(name) }
	typeToObject[TYPE_INT] = func(name string) Attr { return NewIntHolder(name) }
	typeToObject[TYPE_FLOAT] = func(name string) Attr { return NewFloatHolder(name) }
	typeToObject[TYPE_FLOAT64] = func(name string) Attr { return NewFloat64Holder(name) }
	typeToObject[TYPE_INT64] = func(name string) Attr { return NewInt64Holder(name) }
	typeToObject[TYPE_STRING] = func(name string) Attr { return NewStringHolder(name) }
	typeToObject[TYPE_VECTOR2] = func(name string) Attr { return NewVector2Holder(name) }
	typeToObject[TYPE_VECTOR3] = func(name string) Attr { return NewVector3Holder(name) }
	typeToObject[TYPE_BYTES] = func(name string) Attr { return NewBytesHolder(name) }
}

// Create object with type
func Create(typ int, name string) Attr {
	if t, ok := typeToObject[typ]; ok {
		return t(name)
	}
	return nil
}
