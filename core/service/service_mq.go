package service

import (
	"fmt"
	"time"
	"turboengine/common/log"
	"turboengine/common/protocol"
	"turboengine/core/api"
)

const (
	DEFAULT_REPLY = "turbo.service.reply#%s"
)

func makeBody(typ uint8, id string, session uint64, data []byte) *protocol.Message {
	msg := protocol.NewMessage(len(data) + 64)
	sr := protocol.NewStoreArchiver(msg.Body)
	sr.Put(typ)
	sr.Put(id)
	sr.Put(session)
	sr.PutData(data)
	msg.Body = msg.Body[:sr.Len()]
	return msg
}

func parseBody(m *protocol.Message) (typ uint8, id string, session uint64, data []byte, err error) {
	ar := protocol.NewLoadArchiver(m.Body)
	err = ar.Get(&typ)
	if err != nil {
		return
	}
	err = ar.Get(&id)
	if err != nil {
		return
	}
	err = ar.Get(&session)
	if err != nil {
		return
	}
	err = ar.Get(&data)
	return
}

func (s *service) Pub(subject string, data []byte) error {
	msg := makeBody(0, s.c.ID, 0, data)
	return s.exchange.Pub(subject, msg)
}

func (s *service) reply(id string, session uint64, data []byte) error {
	msg := makeBody(1, s.c.ID, session, data)
	return s.exchange.Pub(fmt.Sprintf(DEFAULT_REPLY, id), msg)
}

func (s *service) PubWithTimeout(subject string, data []byte, timeout time.Duration) (*api.Call, error) {
	session := s.session
	s.session++
	msg := makeBody(0, s.c.ID, session, data)
	err := s.exchange.Pub(subject, msg)
	if err != nil {
		msg.Free()
		return nil, err
	}

	call := &api.Call{
		Session:  session,
		DeadLine: time.Now().Add(timeout),
	}

	s.pending[session] = call
	return call, nil
}

func (s *service) SubNoInvoke(subject string) error {
	return s.exchange.Sub(subject)
}

func (s *service) Sub(subject string, invoke api.InvokeFn) error {
	if _, dup := s.delegate[subject]; dup {
		return fmt.Errorf("add subscribe twice, %s ", subject)
	}
	s.delegate[subject] = invoke
	return s.exchange.Sub(subject)
}

func (s *service) UnSub(subject string) {
	if _, ok := s.delegate[subject]; ok {
		delete(s.delegate, subject)
	}
	s.exchange.UnSub(subject)
}

func (s *service) input() { // run on main goroutine
L:
	for {
		select {
		case m := <-s.inMsg:
			subject := string(m.Header)
			s.handle(subject, m)
			m.Free()
		default:
			break L
		}
	}

	// check timeout
	n := time.Now()
	for id, call := range s.pending {
		if call.DeadLine.Sub(n) <= 0 {
			if call.Callback != nil {
				call.Err = ERR_TIMEOUT
				call.Callback(call, nil)
			}
			delete(s.pending, id)
		}
	}
}

func (s *service) handle(subject string, m *protocol.Message) {
	typ, id, session, data, err := parseBody(m)

	if err != nil {
		log.Error("parse msg failed")
		return
	}

	if typ == 0 { // normal message
		// TODO: handle call
		reply := s.invoke(subject, id, data)
		if session != 0 { // need reply
			if reply != nil {
				s.reply(id, session, reply.Body)
				return
			}
			s.reply(id, session, []byte(""))
		}
		return
	}

	if typ == 1 { // reply message
		s.callback(session, data)
	}

}

func (s *service) invoke(subject string, id string, data []byte) *protocol.Message {
	if invoke, ok := s.delegate[subject]; ok {
		reply := invoke(id, data)
		return reply
	}
	return nil
}

func (s *service) callback(session uint64, data []byte) {
	if call, ok := s.pending[session]; ok {
		if call.Callback != nil {
			call.Callback(call, data)
		}
		delete(s.pending, session) // delete
	}
}
