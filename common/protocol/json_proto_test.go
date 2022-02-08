package protocol

import (
	"fmt"
	"testing"
)

type Test struct {
	X, Y int
}

func Test_Decode(t *testing.T) {
	enc := NewJsonEncoder()
	msg, err := enc.Encode(&ProtoMsg{
		Id:   1,
		Data: &Test{1, 1},
	})
	if err != nil {
		t.Fatal(err)
	}
	AddProto(1, func() any { return &Test{} })
	dec := NewJsonDecoder()
	p, err := dec.Decode(msg.Body)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(p, p.Data)
	if p.Id != 1 || p.Data.(*Test).X != 1 || p.Data.(*Test).Y != 1 {
		t.Fatal("fail")
	}

}
