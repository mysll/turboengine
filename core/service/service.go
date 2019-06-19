package service

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
	"turboengine/common/log"
	"turboengine/common/utils"
	"turboengine/core/api"

	"github.com/mysll/toolkit"
)

type Config struct {
	Name string
}

type service struct {
	sync.RWMutex
	toolkit.WaitGroupWrapper
	c        *Config
	handler  api.ServiceHandler
	running  bool
	quit     bool
	time     *utils.Time
	attachs  map[int64]api.AttachFn
	attachId int64
	mods     map[string]api.Module
}

func New(h api.ServiceHandler, c *Config) api.Service {
	s := &service{}
	s.c = c
	s.handler = h
	s.attachs = make(map[int64]api.AttachFn)
	s.mods = make(map[string]api.Module)
	return s
}

// call before Start
func (s *service) Register(mod api.Module) {
	if _, ok := s.mods[mod.Name()]; ok {
		panic(fmt.Errorf("register %s mod twice", mod.Name()))
	}
	s.mods[mod.Name()] = mod
}

func (s *service) Start() error {
	ctx := context.Background()
	if s.running {
		return fmt.Errorf("service %s already running", s.c.Name)
	}
	s.prepare()
	if err := s.handler.OnPrepare(s); err != nil {
		log.Error("prepare %s failed, %s", s.c.Name, err.Error())
		return err
	}

	for _, m := range s.mods {
		if err := m.Init(s); err != nil {
			log.Error("init mod %s failed, %s", m.Name(), err.Error())
			return err
		}
		log.Info("mod %s init", m.Name())
	}

	s.init()

	for _, m := range s.mods {
		m.Start(ctx)
		log.Info("mod %s start", m.Name())
	}

	if err := s.handler.OnStart(); err != nil {
		log.Error("start %s failed, %s", s.c.Name, err.Error())
		return err
	}
	s.run()
	s.destroy()
	return nil
}

func (s *service) prepare() {

}

func (s *service) init() {

}

func (s *service) Attach(fn api.AttachFn) int64 {
	id := atomic.AddInt64(&s.attachId, 1)
	s.attachs[id] = fn
	return id
}

func (s *service) Deatch(id int64) {
	if _, ok := s.attachs[id]; ok {
		delete(s.attachs, id)
	}
}

func (s *service) run() {
	s.running = true
	log.Info("service %s start", s.c.Name)
	s.time = utils.NewTime(10)
	for !s.quit {
		s.time.Update()
		for _, f := range s.attachs {
			f(s.time)
		}
		time.Sleep(s.time.NextFrame())
	}
}

func (s *service) destroy() {
	for _, m := range s.mods {
		m.Close()
		log.Info("mod %s stop", m.Name())
	}
	log.Info("service %s stop", s.c.Name)
}

func (s *service) Close() {
	if s.handler.OnShut() {
		s.Shut()
		log.Info("service %s close", s.c.Name)
	}
}

func (s *service) Shut() {
	if s.quit {
		return
	}
	s.quit = true
}
