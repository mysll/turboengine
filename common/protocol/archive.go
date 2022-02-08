package protocol

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"errors"
	"fmt"
	"io"
	"reflect"
)

// 输入字节流
type StoreArchive struct {
	buf []byte
	pos int
}

type Encoder interface {
	MarshalArchive(io.Writer) error
}

type Decoder interface {
	UnmarshalArchive(io.Reader) error
}

// NewStoreArchive 创建一个新的输出流
func NewStoreArchive(buf []byte) *StoreArchive {
	if buf == nil || cap(buf) == 0 {
		return nil
	}
	ar := &StoreArchive{}
	ar.buf = buf[:0]
	ar.pos = 0
	return ar
}

// Write 写入字节数组
func (ar *StoreArchive) Write(p []byte) (n int, err error) {
	l := len(p)
	if l == 0 {
		return l, nil
	}
	if ar.pos+l > cap(ar.buf) {
		return 0, io.EOF
	}
	ar.buf = ar.buf[:ar.pos+l]
	copy(ar.buf[ar.pos:], p)
	ar.pos += l
	return l, nil
}

// Data 获取已经写入的字节数组
func (ar *StoreArchive) Data() []byte {
	return ar.buf[:ar.pos]
}

// Len 写入的字节数组长度
func (ar *StoreArchive) Len() int {
	return ar.pos
}

// WriteAt 在指定位置定义数据，覆盖写入，不修改原始长度
func (ar *StoreArchive) WriteAt(offset int, val any) error {
	if offset >= cap(ar.buf) {
		return fmt.Errorf("offset out of range")
	}

	old := ar.pos
	ar.pos = offset
	var err error
	switch val.(type) {
	case int8, int16, int32, int64, uint8, uint16, uint32, uint64, float32, float64:
		err = binary.Write(ar, binary.LittleEndian, val)
	case int:
		err = binary.Write(ar, binary.LittleEndian, int32(val.(int)))
	default:
		err = fmt.Errorf("unsupport type")
	}

	ar.pos = old
	return err
}

// Put 写入任意类型数据
func (ar *StoreArchive) Put(val any) error {
	if m, ok := val.(Encoder); ok {
		return m.MarshalArchive(ar)
	}

	switch t := val.(type) {
	case int8, int16, int32, int64, uint8, uint16, uint32, uint64, float32, float64:
		return binary.Write(ar, binary.LittleEndian, t)
	case *int8, *int16, *int32, *int64, *uint8, *uint16, *uint32, *uint64, *float32, *float64:
		v := reflect.ValueOf(val).Elem()
		return binary.Write(ar, binary.LittleEndian, v.Interface())
	case int:
		return binary.Write(ar, binary.LittleEndian, int32(t))
	case uint:
		return binary.Write(ar, binary.LittleEndian, uint32(t))
	case *int:
		return binary.Write(ar, binary.LittleEndian, int32(*t))
	case *uint:
		return binary.Write(ar, binary.LittleEndian, uint32(*t))
	case bool:
		return binary.Write(ar, binary.LittleEndian, t)
	case *bool:
		return binary.Write(ar, binary.LittleEndian, *t)
	case string:
		return ar.PutString(t)
	case *string:
		return ar.PutString(*t)
	case []byte:
		return ar.PutData(t)
	default:
		return ar.PutObject(val)
	}
}

// PutString 写入字符串，格式：uint16长度+字符串
func (ar *StoreArchive) PutString(val string) error {
	if len(val) > 0xFFFF {
		return errors.New("string size too big")
	}
	data := []byte(val)
	size := len(data)
	err := ar.Put(uint16(size))
	if err != nil {
		return err
	}
	_, err = ar.Write(data)
	return err
}

// PutObject 写入go对象
func (ar *StoreArchive) PutObject(obj any) error {
	enc := gob.NewEncoder(ar)
	return enc.Encode(obj)
}

// PutData 写入字节数据，格式为：4字节长度+数据,最大0xFFFFFFFF
func (ar *StoreArchive) PutData(data []byte) error {
	if len(data) > 0xFFFFFFFF {
		return errors.New("data size too big")
	}
	err := ar.Put(uint32(len(data)))
	if err != nil {
		return err
	}
	_, err = ar.Write(data)
	return err
}

// 输出字节流
type LoadArchive struct {
	bytes  []byte
	reader *bytes.Reader
}

func NewLoadArchive(data []byte) *LoadArchive {
	ar := &LoadArchive{}
	ar.reader = bytes.NewReader(data)
	ar.bytes = data
	return ar
}

// Position 获取当前位置
func (ar *LoadArchive) Position() int {
	return int(ar.reader.Size()) - ar.reader.Len()
}

// AvailableBytes 剩余字节长度
func (ar *LoadArchive) AvailableBytes() int {
	return ar.reader.Len()
}

// Size 总容量
func (ar *LoadArchive) Size() int {
	return int(ar.reader.Size())
}

// Seek 移动读指针
func (ar *LoadArchive) Seek(offset int, whence int) (int, error) {
	ret, err := ar.reader.Seek(int64(offset), whence)
	return int(ret), err
}

// io.Reader
func (ar *LoadArchive) Read(p []byte) (n int, err error) {
	return ar.reader.Read(p)
}

// Get 读取任意类型的数据
func (ar *LoadArchive) Get(val any) (err error) {
	// dpv := reflect.ValueOf(val)
	// if dpv.Kind() != reflect.Ptr {
	// 	return errors.New("destination not a pointer")
	// }
	// if dpv.IsNil() {
	// 	return errors.New("destination pointer is nil")
	// }

	if m, ok := val.(Decoder); ok {
		return m.UnmarshalArchive(ar.reader)
	}

	switch val.(type) {
	case *int8, *int16, *int32, *int64, *uint8, *uint16, *uint32, *uint64, *float32, *float64:
		return binary.Read(ar.reader, binary.LittleEndian, val)
	case *int:
		var out int32
		err = binary.Read(ar.reader, binary.LittleEndian, &out)
		if err != nil {
			return err
		}
		*(val.(*int)) = int(out)
		return nil
	case *uint:
		var out uint32
		err = binary.Read(ar.reader, binary.LittleEndian, &out)
		if err != nil {
			return err
		}
		*(val.(*uint)) = uint(out)
		return nil
	case *bool:
		var out bool
		err = binary.Read(ar.reader, binary.LittleEndian, &out)
		if err != nil {
			return err
		}
		*(val.(*bool)) = out
		return nil
	case *string:
		inst := val.(*string)
		*inst, err = ar.GetString()
		return err
	case *[]byte:
		inst := val.(*[]byte)
		*inst, err = ar.GetData()
		return err
	default:
		return ar.GetObject(val)
	}
}

func (ar *LoadArchive) GetInt8() (val int8, err error) {
	err = ar.Get(&val)
	return
}

func (ar *LoadArchive) GetUInt8() (val uint8, err error) {
	err = ar.Get(&val)
	return
}

func (ar *LoadArchive) GetInt16() (val int16, err error) {
	err = ar.Get(&val)
	return
}

func (ar *LoadArchive) GetUint16() (val uint16, err error) {
	err = ar.Get(&val)
	return
}

func (ar *LoadArchive) GetInt32() (val int32, err error) {
	err = ar.Get(&val)
	return
}

func (ar *LoadArchive) GetUint32() (val uint32, err error) {
	err = ar.Get(&val)
	return
}

func (ar *LoadArchive) GetInt64() (val int64, err error) {
	err = ar.Get(&val)
	return
}

func (ar *LoadArchive) GetUint64() (val uint64, err error) {
	err = ar.Get(&val)
	return
}

func (ar *LoadArchive) GetFloat32() (val float32, err error) {
	err = ar.Get(&val)
	return
}

func (ar *LoadArchive) GetFloat64() (val float64, err error) {
	err = ar.Get(&val)
	return
}

// GetString 读取带前缀长度的字符串
func (ar *LoadArchive) GetString() (val string, err error) {
	l, err := ar.GetUint16()
	if err != nil {
		return "", err
	}
	if l == 0 {
		val = ""
		return
	}
	data := make([]byte, l)
	_, err = ar.reader.Read(data)
	if err != nil {
		return
	}
	val = string(data)
	return
}

// GetObject 读取go对象
func (ar *LoadArchive) GetObject(val any) error {
	dec := gob.NewDecoder(ar.reader)
	return dec.Decode(val)
}

// GetData 读带前缀长度的字节流
func (ar *LoadArchive) GetData() (data []byte, err error) {
	var l uint32
	l, err = ar.GetUint32()
	if err != nil {
		return nil, err
	}
	if l == 0 {
		return []byte{}, nil
	}
	data = make([]byte, int(l))
	_, err = ar.reader.Read(data)
	return data, err
}

func (ar *LoadArchive) GetDataNonCopy() (data []byte, err error) {
	var l uint32
	l, err = ar.GetUint32()
	if err != nil {
		return nil, err
	}
	if l == 0 {
		return []byte{}, nil
	}
	curpos, _ := ar.reader.Seek(0, io.SeekCurrent)
	data = ar.bytes[int(curpos) : int(curpos)+int(l)]
	_, err = ar.reader.Seek(int64(l), io.SeekCurrent)
	if err != nil {
		return nil, err
	}
	return data, nil
}
