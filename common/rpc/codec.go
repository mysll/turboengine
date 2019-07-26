package rpc

import (
	"bufio"
	"encoding/binary"
	"errors"
	"io"
	"turboengine/common/protocol"

	"github.com/nggenius/ngengine/utils"
)

const RPC_MAX_LEN = 1024 * 64 // 64k

var (
	ErrTooLong = errors.New("message is to long")
)

type ServerCodec interface {
	ReadRequest(req *Request, maxrc uint32) error
	WriteResponse(seq uint64, errcode int32, body *protocol.Message) (err error)
	GetConn() io.ReadWriteCloser
	Close() error
}

type ByteServerCodec struct {
	rwc    io.ReadWriteCloser
	encBuf *bufio.Writer
	closed bool
}

func (c *ByteServerCodec) ReadRequest(req *Request, maxrc uint32) (err error) {
	var msg *protocol.Message

	msg, err = readMessage(c.rwc, maxrc)
	if err != nil {
		return err
	}
	ar := protocol.NewLoadArchiver(msg.Header)

	req.Seq, err = ar.GetUint64()
	if err != nil {
		return err
	}

	req.ServiceMethod, err = ar.GetString()
	if err != nil {
		return err
	}

	req.Raw = msg

	return nil
}

func readMessage(rwc io.ReadWriteCloser, maxrc uint32) (*protocol.Message, error) {
	var sz uint32
	var headlen uint8
	var err error
	var msg *protocol.Message

	if err = binary.Read(rwc, binary.LittleEndian, &sz); err != nil {
		return nil, err
	}

	// Limit messages to the maximum receive value, if not
	// unlimited.  This avoids a potential denaial of service.
	if sz < 0 || (maxrc > 0 && sz > maxrc) {
		return nil, ErrTooLong
	}

	if err = binary.Read(rwc, binary.LittleEndian, &headlen); err != nil {
		return nil, err
	}

	if headlen > protocol.MAX_HEADER_LEN {
		return nil, ErrTooLong
	}

	bodylen := int(sz - uint32(headlen))
	msg = protocol.NewMessage(bodylen)
	msg.Header = msg.Header[0:headlen]
	if _, err = io.ReadFull(rwc, msg.Header); err != nil {
		msg.Free()
		return nil, err
	}

	if bodylen > 0 {
		msg.Body = msg.Body[0:bodylen]

		if _, err = io.ReadFull(rwc, msg.Body); err != nil {
			msg.Free()
			return nil, err
		}
	}

	return msg, nil
}

func (c *ByteServerCodec) WriteResponse(seq uint64, errcode int32, body *protocol.Message) (err error) {
	if body == nil {
		body = protocol.NewMessage(1)
	}

	body.Header = body.Header[:0]
	w := utils.NewStoreArchiver(body.Header)
	w.Put(seq)
	if errcode != 0 {
		w.Put(int8(1))
		w.Put(errcode)
	} else {
		w.Put(int8(0))
	}
	body.Header = body.Header[:w.Len()]
	size := len(body.Header) + len(body.Body)
	if size > RPC_MAX_LEN {
		return ErrTooLong
	}

	binary.Write(c.encBuf, binary.LittleEndian, uint32(size))            //数据大小
	binary.Write(c.encBuf, binary.LittleEndian, uint8(len(body.Header))) //头部大小
	c.encBuf.Write(body.Header)
	if len(body.Body) > 0 {
		c.encBuf.Write(body.Body)
	}
	body.Header = body.Header[:0]
	return c.encBuf.Flush()
}

func (c *ByteServerCodec) Close() error {
	if c.closed {
		// Only call c.rwc.Close once; otherwise the semantics are undefined.
		return nil
	}
	c.closed = true
	return c.rwc.Close()
}

func (c *ByteServerCodec) GetConn() io.ReadWriteCloser {
	return c.rwc
}

type ClientCodec interface {
	// WriteRequest must be safe for concurrent use by multiple goroutines.
	WriteRequest(seq uint64, args *protocol.Message) error
	ReadMessage() (*protocol.Message, error)
	GetAddress() string
	Close() error
}

type ByteClientCodec struct {
	rwc    io.ReadWriteCloser
	encBuf *bufio.Writer
	maxrx  uint32
}

func (b *ByteClientCodec) WriteRequest(seq uint64, args *protocol.Message) error {
	return nil
}

func (b *ByteClientCodec) ReadMessage() (*protocol.Message, error) {
	return nil, nil
}

func (b *ByteClientCodec) GetAddress() string {
	return ""
}

func (b *ByteClientCodec) Close() error {
	return nil
}
