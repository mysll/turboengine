package storage

import (
	"turboengine/common/utils"
	"turboengine/core/api"
	"turboengine/core/plugin"
)

const (
	Name = "Storage"
)

type Driver interface {
}

type ResultCallback func(interface{}, error)

type Request struct {
	op   int
	data interface{}
	cb   ResultCallback
}

type Response struct {
	data interface{}
	err  error
	cb   ResultCallback
}

type Storage struct {
	srv     api.Service
	id      uint64
	pending chan *Request
	done    chan *Response
	driver  Driver
}

func (s *Storage) Prepare(srv api.Service, args ...interface{}) {
	s.srv = srv
	s.id = s.srv.Attach(s.roundInvoke)
	s.pending = make(chan *Request, 1024)
	s.done = make(chan *Response, 1024)
}

func (s *Storage) Shut(srv api.Service) {
	s.srv.Detach(s.id)
}

func (s *Storage) Run() {
}

func (s *Storage) Handle(cmd string, args ...interface{}) interface{} {
	return nil
}

func (s *Storage) roundInvoke(t *utils.Time) {
	for {
		select {
		case r := <-s.done:
			r.cb(r.data, r.err)
		default:
			return
		}
	}
}

func init() {
	plugin.Register(Name, &Storage{})
}
