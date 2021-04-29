package storage

import (
	"turboengine/common/log"
	"turboengine/common/utils"
	"turboengine/core/api"
	"turboengine/core/plugin"
	"turboengine/core/plugin/storage/driver"
	"turboengine/gameplay/dao"
)

const (
	Name = "Storage"
)

const (
	DB_OP_CREATE = 1 + iota
	DB_OP_SELECT
	DB_OP_INSERT
	DB_OP_UPDATE
	DB_OP_DEL
)

type Driver interface {
	Connect(dsn string) error
	Create(name string, model interface{}) error
	Select(id uint64, data interface{}) error
	Insert(id uint64, data interface{}) (uint64, error)
	Update(id uint64, data interface{}) error
	Del(name string, id uint64) error
}

type ResultCallback func(interface{}, error)

type Request struct {
	op   int
	id   uint64
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
	if dt, ok := args[0].(string); ok {
		if dsn, ok := args[1].(string); ok {
			switch dt {
			case "mysql":
				s.driver = new(driver.MysqlDao)
				if err := s.driver.Connect(dsn); err != nil {
					panic(err)
				}
			}
		}
	}
	for n, m := range dao.GetAllModel() {
		if err := s.driver.Create(n, m); err != nil {
			panic(err)
		}
		log.Info("create model ", n, " ok")
	}
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
