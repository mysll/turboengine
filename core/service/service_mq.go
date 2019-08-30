package service

import (
	"errors"
	"fmt"
	"sync/atomic"
	"time"
	"turboengine/common/log"
	"turboengine/common/protocol"
	"turboengine/core/api"
)

const (
	DEFAULT_REPLY               = "%d:#.reply"
	SERVICE_SHUT                = "#.shut"
	SERVICE_SHUT_ALL            = "#.shut_all"
	PRE_ROUND_MAX_PROCESS_COUNT = 100
)

func makeBody(typ uint8, id uint16, session uint64, data []byte) *protocol.Message {
	msg := protocol.NewMessage(len(data) + 64)
	sr := protocol.NewStoreArchiver(msg.Body)
	sr.Put(typ)
	sr.Put(id)
	sr.Put(session)
	sr.Put(int8(0)) // error
	sr.PutData(data)
	msg.Body = msg.Body[:sr.Len()]
	return msg
}

func makeErrorBody(typ uint8, id uint16, session uint64, err error) *protocol.Message {
	msg := protocol.NewMessage(len(err.Error()) + 64)
	sr := protocol.NewStoreArchiver(msg.Body)
	sr.Put(typ)
	sr.Put(id)
	sr.Put(session)
	sr.Put(int8(1)) // error
	sr.Put(err.Error())
	msg.Body = msg.Body[:sr.Len()]
	return msg
}

func parseBody(m *protocol.Message) (typ uint8, id uint16, session uint64, data []byte, err error) {
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
	var code int8
	err = ar.Get(&code)
	if err != nil {
		return
	}

	if code != 0 { // check error
		var e string
		ar.Get(&e)
		err = errors.New(e)
		return
	}

	data, err = ar.GetDataNonCopy()
	return
}

func (s *service) Pub(subject string, data []byte) error {
	msg := makeBody(0, s.c.ID, 0, data)
	return s.exchange.Pub(subject, msg)
}

func (s *service) reply(id uint16, session uint64, data []byte) error {
	msg := makeBody(1, s.c.ID, session, data)
	return s.exchange.Pub(fmt.Sprintf(DEFAULT_REPLY, id), msg)
}

func (s *service) replyError(id uint16, session uint64, err error) error {
	msg := makeErrorBody(1, s.c.ID, session, err)
	return s.exchange.Pub(fmt.Sprintf(DEFAULT_REPLY, id), msg)
}

func (s *service) PubWithTimeout(subject string, data []byte, timeout time.Duration) (*api.Call, error) {
	session := atomic.AddUint64(&s.session, 1)
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

	s.lockCall.Lock()
	s.pending[session] = call
	s.lockCall.Unlock()
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
	log.Info("subscribe subject ", subject)
	return s.exchange.Sub(subject)
}

func (s *service) UnSub(subject string) {
	if _, ok := s.delegate[subject]; ok {
		delete(s.delegate, subject)
		s.exchange.UnSub(subject)
	}
}

func (s *service) innerHandle(subject string, m *protocol.Message) bool {
	res := true
	switch subject {
	case SERVICE_SHUT:
		_, _, _, data, err := parseBody(m)
		if err != nil {
			panic(err)
		}
		if s.sid == string(data) {
			s.Close()
		}
	case SERVICE_SHUT_ALL:
		if s.id != 0 {
			s.Close()
		}
	default:
		res = false
	}

	if res {
		m.Free()
	}
	return res
}

func (s *service) input() { // run on main goroutine
L:
	for i := 0; i < PRE_ROUND_MAX_PROCESS_COUNT; i++ {
		select {
		case m := <-s.inMsg:
			subject := string(m.Header)
			if !s.innerHandle(subject, m) {
				s.handle(subject, m)
			}
		default:
			break L
		}
	}

	// check timeout
	n := time.Now()
	s.lockCall.Lock()
	var timeout []*api.Call
	for id, call := range s.pending {
		if call.DeadLine.Sub(n) <= 0 {
			timeout = append(timeout, call)
			delete(s.pending, id)
		}
	}
	s.lockCall.Unlock()
	if len(timeout) > 0 {
		for _, call := range timeout {
			if call.DeadLine.Sub(n) <= 0 {
				call.Err = ERR_TIMEOUT
				call.Data = nil

				if call.Done != nil {
					call.Done <- call
				} else if call.Callback != nil {
					call.Callback(call)
				}
			}
		}
	}

}

func (s *service) handle(subject string, m *protocol.Message) {
	typ, id, session, data, err := parseBody(m)

	if err != nil {
		m.Free()
		if typ == 1 && session != 0 {
			s.callbackError(session, err)
			return
		}
		log.Error("parse msg failed")
		return
	}

	if typ == 0 { // normal message
		//  sync invoke call
		reply, err := s.invoke(subject, id, data)
		m.Free()
		if session != 0 { // need reply
			if err != nil {
				s.replyError(id, session, err)
				return
			}
			if reply != nil {
				s.reply(id, session, reply.Body)
				reply.Free()
				return
			}
			s.reply(id, session, []byte(""))
		}
		return
	}

	if typ == 1 { // reply message
		s.callback(session, m, data)
	}

}

func (s *service) invoke(subject string, id uint16, data []byte) (*protocol.Message, error) {
	if invoke, ok := s.delegate[subject]; ok {
		reply, err := invoke(id, data)
		return reply, err
	}
	return nil, fmt.Errorf("subject %s not handle", subject)
}

func (s *service) callback(session uint64, msg *protocol.Message, data []byte) {
	s.lockCall.RLock()
	call, ok := s.pending[session]
	s.lockCall.RUnlock()
	if ok {
		call.Data = data

		if call.Done != nil {
			call.Msg = msg // if msg call free, call.Data will be gc.
			call.Done <- call
		} else if call.Callback != nil {
			call.Callback(call)
			msg.Free()
		}

		s.lockCall.Lock()
		delete(s.pending, session) // delete
		s.lockCall.Unlock()
	}
}

func (s *service) callbackError(session uint64, err error) {
	s.lockCall.RLock()
	call, ok := s.pending[session]
	s.lockCall.RUnlock()

	if ok {
		call.Err = err
		call.Data = nil

		if call.Done != nil {
			call.Done <- call
		} else if call.Callback != nil {
			call.Callback(call)
		}
		s.lockCall.Lock()
		delete(s.pending, session) // delete
		s.lockCall.Unlock()
	}
}
