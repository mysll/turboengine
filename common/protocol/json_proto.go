package protocol

import (
	"encoding/json"
	"fmt"
)

var (
	proto = make(map[int32]Builder)
)

type JsonEncoder struct {
}

func NewJsonEncoder() *JsonEncoder {
	j := &JsonEncoder{}
	return j
}

func (e *JsonEncoder) Encode(p *ProtoMsg) (*Message, error) {
	msg := NewAutoExtendArchive(128)
	msg.Put(p.Id)
	data, err := json.Marshal(p.Data)
	if err != nil {
		msg.Free()
		return nil, err
	}

	msg.Append(data)
	return msg.Message(), nil
}

type JsonDecoder struct {
}

type Builder func() interface{}

func NewJsonDecoder() *JsonDecoder {
	j := &JsonDecoder{}
	return j
}

func AddProto(id int32, builder Builder) {
	proto[id] = builder
}

func (d *JsonDecoder) Decode(data []byte) (*ProtoMsg, error) {
	ar := NewLoadArchive(data)
	id, err := ar.GetInt32()
	if err != nil {
		return nil, err
	}

	if fn, ok := proto[id]; ok {
		pt := fn()
		err = json.Unmarshal(data[4:], pt)
		if err != nil {
			return nil, err
		}
		p := &ProtoMsg{
			Id:   id,
			Data: pt,
		}

		return p, nil
	}
	return nil, fmt.Errorf("proto not found, %d", id)
}
