package protocol_test

import (
	"bytes"
	"testing"
	"turboengine/common/protocol"
)

func TestLoadArchive(t *testing.T) {
	buf := make([]byte, 0, 1024)
	store := protocol.NewStoreArchive(buf)
	v1 := int8(1)
	v2 := int16(2)
	v3 := int32(3)
	v4 := int64(4)
	v5 := float32(5)
	v6 := float64(6)

	v11 := &v1
	v21 := &v2
	v31 := &v3
	v41 := &v4
	v51 := &v5
	v61 := &v6
	v71 := []uint8{1, 2, 3}

	store.Put(v1)
	store.Put(v2)
	store.Put(v3)
	store.Put(v4)
	store.Put(v5)
	store.Put(v6)

	store.Put(v11)
	store.Put(v21)
	store.Put(v31)
	store.Put(v41)
	store.Put(v51)
	store.Put(v61)

	var str = "test"
	store.Put(str)
	store.Put(&str)

	store.Put(v71)

	var x1 int8
	var x2 int16
	var x3 int32
	var x4 int64
	var x5 float32
	var x6 float64

	load := protocol.NewLoadArchive(store.Data())
	if err := load.Get(&x1); err != nil {
		t.Fatalf(err.Error())
	}
	if err := load.Get(&x2); err != nil {
		t.Fatalf(err.Error())
	}
	if err := load.Get(&x3); err != nil {
		t.Fatalf(err.Error())
	}
	if err := load.Get(&x4); err != nil {
		t.Fatalf(err.Error())
	}
	if err := load.Get(&x5); err != nil {
		t.Fatalf(err.Error())
	}
	if err := load.Get(&x6); err != nil {
		t.Fatalf(err.Error())
	}
	if x1 != v1 || x2 != v2 || x3 != v3 || x4 != v4 || x5 != v5 || x6 != v6 {
		t.Fatalf("not match, need %v %v %v %v %v %v get %v %v %v %v %v %v", v1, v2, v3, v4, v5, v6, x1, x2, x3, x4, x5, x6)
	}

	if err := load.Get(&x1); err != nil {
		t.Fatalf(err.Error())
	}
	if err := load.Get(&x2); err != nil {
		t.Fatalf(err.Error())
	}
	if err := load.Get(&x3); err != nil {
		t.Fatalf(err.Error())
	}
	if err := load.Get(&x4); err != nil {
		t.Fatalf(err.Error())
	}
	if err := load.Get(&x5); err != nil {
		t.Fatalf(err.Error())
	}
	if err := load.Get(&x6); err != nil {
		t.Fatalf(err.Error())
	}
	if x1 != v1 || x2 != v2 || x3 != v3 || x4 != v4 || x5 != v5 || x6 != v6 {
		t.Fatalf("not match, need %v %v %v %v %v %v get %v %v %v %v %v %v", v1, v2, v3, v4, v5, v6, x1, x2, x3, x4, x5, x6)
	}

	var xstr string
	load.Get(&xstr)
	if xstr != str {
		t.Fatalf("not match")
	}

	var xstr1 string
	load.Get(&xstr1)
	if xstr1 != str {
		t.Fatalf("not match")
	}

	var xbyte []uint8
	load.Get(&xbyte)
	if !bytes.Equal(xbyte, v71) {
		t.Fatalf("not match")
	}

}

func TestLoadArchive_GetDataNonCopy(t *testing.T) {
	buf := make([]byte, 0, 1024)
	store := protocol.NewStoreArchive(buf)
	store.Put(int32(1))
	store.PutData([]byte{0, 1, 2, 3, 4, 5, 6})
	store.Put(int64(2))

	load := protocol.NewLoadArchive(store.Data())
	var i int32
	var j int64
	load.Get(&i)
	if i != 1 {
		t.Fatal("not equal")
	}
	bytes, err := load.GetDataNonCopy()
	if err != nil {
		t.Fatal(err)
	}

	for k, v := range bytes {
		if k != int(v) {
			t.Fatal("not equal")
		}
	}
	load.Get(&j)
	if j != 2 {
		t.Fatal("not equal")
	}
}

func TestAutoExtendArchive_Put(t *testing.T) {
	a := protocol.NewAutoExtendArchive(64)
	buf := make([]byte, 60)
	a.Put(buf)
	a.Put(int64(1))
	a.Put("hello")
	a.Put(int64(3))
	a.Put(int64(4))

	msg := a.Message()
	load := protocol.NewLoadArchive(msg.Body)
	load.GetData()
	var i, k, l int64
	var j string
	load.Get(&i)
	load.Get(&j)
	load.Get(&k)
	load.Get(&l)
	if i != 1 || j != "hello" || k != 3 || l != 4 && len(msg.Body) != 64+24+7 {
		t.Fatal("not equal")
	}
}
