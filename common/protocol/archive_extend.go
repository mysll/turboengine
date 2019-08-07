package protocol

import "io"

type AutoExtendArchive struct {
	msg *Message
	sr  *StoreArchive
}

func NewAutoExtendArchive(init_cap int) *AutoExtendArchive {
	a := &AutoExtendArchive{}
	a.msg = NewMessage(init_cap)
	a.sr = NewStoreArchiver(a.msg.Body)
	return a
}

func (a *AutoExtendArchive) Put(val interface{}) error {
	err := a.sr.Put(val)
	if err == io.EOF {
		msg := NewMessage(cap(a.msg.Body) * 2)
		sr := NewStoreArchiver(msg.Body)
		sr.Write(a.sr.Data())
		a.sr = sr
		a.msg = msg
		return a.sr.Put(val)
	}
	return err
}

func (a *AutoExtendArchive) Message() *Message {
	a.msg.Body = a.msg.Body[:a.sr.Len()]
	return a.msg
}
