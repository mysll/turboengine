package rpc

import (
	"io"
	"turboengine/common/protocol"
)

type Codec interface {
	ReadRequest(maxrc uint16) (*protocol.Message, error)
	WriteResponse(seq uint64, errcode int32, body *protocol.Message) (err error)
	GetConn() io.ReadWriteCloser
	Close() error
}

type ByteCodec struct {
	rwc io.ReadWriteCloser
}

func (codec *ByteCodec) ReadRequest(maxrc uint16) (*protocol.Message, error) {
	return nil, nil
}

func (codec *ByteCodec) WriteResponse(seq uint64, errcode int32, body *protocol.Message) (err error) {
	return nil
}

func (codec *ByteCodec) GetConn() io.ReadWriteCloser {
	return codec.rwc
}

func (codec *ByteCodec) Close() error {
	return nil
}
