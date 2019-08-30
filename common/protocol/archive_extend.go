package protocol

import "io"

type AutoExtendArchive struct {
	msg *Message
	sr  *StoreArchive
}

func NewAutoExtendArchive(initCap int) *AutoExtendArchive {
	a := &AutoExtendArchive{}
	a.msg = NewMessage(initCap)
	a.sr = NewStoreArchiver(a.msg.Body)
	return a
}

func (a *AutoExtendArchive) Put(val interface{}) error {
	err := a.sr.Put(val)
	for err == io.EOF {
		msg := NewMessage(cap(a.msg.Body) * 2)
		sr := NewStoreArchiver(msg.Body)
		sr.Write(a.sr.Data())
		a.sr = sr
		a.msg = msg
		err = a.sr.Put(val)
		if err == nil {
			break
		}
	}
	return err
}

func (a *AutoExtendArchive) Append(data []byte) error {
	_, err := a.sr.Write(data)
	if err == io.EOF {
		msg := NewMessage(cap(a.msg.Body) * 2)
		sr := NewStoreArchiver(msg.Body)
		sr.Write(a.sr.Data())
		a.sr = sr
		a.msg = msg
		_, err = a.sr.Write(data)
	}
	return err
}

func (a *AutoExtendArchive) Message() *Message {
	a.msg.Body = a.msg.Body[:a.sr.Len()]
	return a.msg
}

func (a *AutoExtendArchive) Free() {
	a.msg.Free()
}
