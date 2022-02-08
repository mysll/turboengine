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
	DB_OP_FIND
	DB_OP_SAVE
	DB_OP_UPDATE
	DB_OP_DEL
)

type Driver interface {
	Connect(dsn string) error
	Create(name string, model dao.Persistent) error
	Find(id uint64, data dao.Persistent) error
	FindBy(data dao.Persistent, where string, args ...any) error
	FindAll(data any, where string, args ...any) error
	Save(data dao.Persistent) (uint64, error)
	Update(data dao.Persistent) error
	Del(data dao.Persistent) error
	DelBy(data dao.Persistent, where string, args ...any) error
}

type ResultCallback func(any, error)

type Request struct {
	op    int
	id    uint64
	multi bool
	data  any
	where string
	args  []any
	cb    ResultCallback
}

type Response struct {
	data any
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

func (s *Storage) Prepare(srv api.Service, args ...any) {
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

func (s *Storage) Handle(cmd string, args ...any) any {
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

func (s *Storage) Find(id uint64, data dao.Persistent) error {
	return s.driver.Find(id, data)
}

func (s *Storage) FindWithCallback(id uint64, data dao.Persistent, cb ResultCallback) error {
	if cb == nil {
		return s.driver.Find(id, data)
	}
	req := &Request{
		op:   DB_OP_FIND,
		id:   id,
		data: data,
		cb:   cb,
	}
	s.pending <- req
	return nil
}

func (s *Storage) FindBy(data dao.Persistent, where string, args ...any) error {
	return s.driver.FindBy(data, where, args...)
}

func (s *Storage) FindByWithCallback(data dao.Persistent, where string, args []any, cb ResultCallback) error {
	if cb == nil {
		return s.driver.FindBy(data, where, args...)
	}
	req := &Request{
		op:    DB_OP_FIND,
		data:  data,
		where: where,
		args:  args,
		cb:    cb,
	}
	s.pending <- req
	return nil
}

func (s *Storage) FindAll(data any, where string, args ...any) error {
	return s.driver.FindAll(data, where, args...)
}

func (s *Storage) FindAllWithCallback(data any, where string, args []any, cb ResultCallback) error {
	if cb == nil {
		return s.driver.FindAll(data, where, args...)
	}
	req := &Request{
		op:    DB_OP_FIND,
		data:  data,
		multi: true,
		where: where,
		args:  args,
		cb:    cb,
	}
	s.pending <- req
	return nil
}

func (s *Storage) Save(data dao.Persistent) (uint64, error) {
	return s.driver.Save(data)
}

func (s *Storage) SaveWithCallback(data dao.Persistent, cb ResultCallback) (uint64, error) {
	if cb == nil {
		return s.driver.Save(data)
	}
	req := &Request{
		op:   DB_OP_SAVE,
		data: data,
		cb:   cb,
	}
	s.pending <- req
	return 0, nil
}

func (s *Storage) Update(data dao.Persistent) error {
	return s.driver.Update(data)
}

func (s *Storage) UpdateWithCallback(data dao.Persistent, cb ResultCallback) error {
	if cb == nil {
		return s.driver.Update(data)
	}
	req := &Request{
		op:   DB_OP_UPDATE,
		data: data,
		cb:   cb,
	}
	s.pending <- req
	return nil
}

func (s *Storage) Del(data dao.Persistent) error {
	return s.driver.Del(data)
}

func (s *Storage) DelWithCallback(data dao.Persistent, cb ResultCallback) error {
	if cb == nil {
		return s.driver.Del(data)
	}
	req := &Request{
		op:   DB_OP_DEL,
		data: data,
		cb:   cb,
	}
	s.pending <- req
	return nil
}

func (s *Storage) DelBy(data dao.Persistent, where string, args ...any) error {
	return s.driver.DelBy(data, where, args...)
}

func (s *Storage) DelByWithCallback(data dao.Persistent, where string, args []any, cb ResultCallback) error {
	if cb == nil {
		return s.driver.DelBy(data, where, args...)
	}
	req := &Request{
		op:    DB_OP_DEL,
		data:  data,
		where: where,
		args:  args,
		cb:    cb,
	}
	s.pending <- req
	return nil
}

func init() {
	plugin.Register(Name, &Storage{})
}
