package protocol

import (
	"fmt"
	"testing"
	"turboengine/common/protocol"
)

type Test struct {
	X, Y int
}

func Test_Decode(t *testing.T) {
	enc := NewEncoder()
	msg, err := enc.Encode(&protocol.ProtoMsg{
		Id:   1,
		Data: &Test{1, 1},
	})
	if err != nil {
		t.Fatal(err)
	}
	AddProto(1, func() interface{} { return &Test{} })
	dec := NewDecoder()
	p, err := dec.Decode(msg.Body)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(p, p.Data)
	if p.Id != 1 || p.Data.(*Test).X != 1 || p.Data.(*Test).Y != 1 {
		t.Fatal("fail")
	}

}
