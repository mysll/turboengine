package proto

import "turboengine/common/protocol"

type JsonEncoder struct {
}

func (e *JsonEncoder) Encode(*protocol.ProtoMsg) (*protocol.Message, error) {
	return nil, nil
}

type JsonDecoder struct {
}

func (d *JsonDecoder) Decode([]byte) (*protocol.ProtoMsg, error) {
	return nil, nil
}
