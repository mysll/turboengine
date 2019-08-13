package protocol

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
)

func WriteMsg(w io.Writer, data []byte) (err error) {
	if w == nil {
		err = errors.New("writer is nil")
		return
	}

	if err = binary.Write(w, binary.LittleEndian, uint32(len(data))); err != nil {
		return
	}

	if len(data) > 0 {
		if _, err = w.Write(data); err != nil {
			return
		}
	}

	return
}

func ReadMsg(r io.Reader, max uint32) (msg *Message, err error) {
	if r == nil {
		err = errors.New("reader is nil")
		return
	}
	var size uint32
	if err = binary.Read(r, binary.LittleEndian, &size); err != nil {
		return
	}

	if size >= max {
		err = errors.New(fmt.Sprintf("message size exceed, package size: %d, max: %d", size, max))
		return
	}

	msg = NewMessage(int(size))
	msg.Body = msg.Body[0:size]
	if _, err = io.ReadFull(r, msg.Body); err != nil {
		msg.Free()
		msg = nil
	}
	return
}
