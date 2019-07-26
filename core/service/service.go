package service

import (
	"context"
	"fmt"
	"sync"
	"time"
	"turboengine/common/log"
	"turboengine/common/protocol"
	"turboengine/common/utils"
	"turboengine/core/api"
	"turboengine/core/plugin/event"

	consulapi "github.com/hashicorp/consul/api"
	"github.com/mysll/toolkit"
)

type Config struct {
	ID      string
	Name    string
	NatsUrl string
}

type attach struct {
	id uint64
	fn api.Update
}

type service struct {
	sync.RWMutex
	toolkit.WaitGroupWrapper
	c        *Config
	handler  api.ServiceHandler
	running  bool
	quit     bool
	time     *utils.Time
	attachs  []attach
	attachId uint64 // used for attachs
	mods     []api.Module
	exchange *Exchange
	inMsg    chan *protocol.Message // receive message from message queue
	pending  map[uint64]*api.Call   // pending call
	session  uint64                 // used for pending
	delegate map[string]api.InvokeFn
	plugin   map[string]api.Plugin
	lookup   *LookupService
	event    *event.Event
}

func New(h api.ServiceHandler, c *Config) api.Service {
	s := &service{
		c:        c,
		handler:  h,
		attachs:  make([]attach, 0, 8),
		attachId: 1,
		mods:     make([]api.Module, 0, 8),
		inMsg:    make(chan *protocol.Message, 256),
		pending:  make(map[uint64]*api.Call),
		session:  1,
		delegate: make(map[string]api.InvokeFn),
		plugin:   make(map[string]api.Plugin),
		lookup:   NewLookupService(consulapi.DefaultConfig()),
		event:    new(event.Event),
	}

	return s
}

// call before Start
func (s *service) AddModule(mod api.Module) {
	for _, m := range s.mods {
		if m.Name() == mod.Name() {
			panic(fmt.Errorf("register %s mod twice", mod.Name()))
		}
	}
	s.mods = append(s.mods, mod)
}

// async call
func (s *service) notify(event string, id string) {
	s.event.AsyncEmit(event, id)
}

func (s *service) OnServiceChange(event string, id interface{}) {
	switch event {
	case EVENT_ADD:
		log.Info("service add:", id.(string))
	case EVENT_DEL:
		log.Info("service del:", id.(string))
	}
}

func (s *service) Start() error {
	ctx := context.Background()
	if s.running {
		return fmt.Errorf("service %s already running", s.c.Name)
	}
	s.prepare()
	if err := s.handler.OnPrepare(s); err != nil {
		log.Errorf("prepare %s failed, %s", s.c.Name, err.Error())
		return err
	}

	for _, m := range s.mods {
		if err := m.Init(s); err != nil {
			log.Errorf("init mod %s failed, %s", m.Name(), err.Error())
			return err
		}
		log.Infof("mod %s init", m.Name())
	}

	s.init()

	for _, m := range s.mods {
		m.Start(ctx)
		log.Infof("mod %s start", m.Name())
	}

	if err := s.handler.OnStart(); err != nil {
		log.Errorf("start %s failed, %s", s.c.Name, err.Error())
		return err
	}
	s.run()
	s.destroy()
	return nil
}

func (s *service) prepare() {
	s.usePlugin(event.Name, s.event)
	s.event.AddListener(EVENT_ADD, s.OnServiceChange)
	s.event.AddListener(EVENT_DEL, s.OnServiceChange)
}

func (s *service) init() {
	err := s.lookup.Start()
	if err != nil {
		log.Fatal(err)
	}
	err = s.lookup.Register(s.c.ID, s.c.Name, "127.0.0.1", toolkit.RandRange(1, 65535))
	if err != nil {
		log.Fatal(err)
	}
	s.lookup.SetNotify(s.notify)
	p, err := NewExchange(s.inMsg)
	if err != nil {
		log.Fatal(err)
	}

	s.exchange = p
	err = s.exchange.Start(s.c.NatsUrl)
	if err != nil {
		log.Fatal(err)
	}
	s.SubNoInvoke(fmt.Sprintf(DEFAULT_REPLY, s.c.ID)) // inner reply
}

func (s *service) Attach(fn api.Update) uint64 {
	id := s.attachId
	s.attachId++
	s.attachs = append(s.attachs, attach{
		id: id,
		fn: fn,
	})
	return id
}

func (s *service) Deatch(id uint64) {
	for k, a := range s.attachs {
		if a.id == id {
			s.attachs = append(s.attachs[:k], s.attachs[k+1:]...)
			return
		}
	}
}

func (s *service) run() {
	s.running = true
	log.Infof("service %s start", s.c.Name)
	s.time = utils.NewTime(10)
	for !s.quit {
		s.time.Update()
		s.input() // message queue a round trip
		for _, f := range s.attachs {
			f.fn(s.time)
		}
		time.Sleep(s.time.NextFrame())
	}
}

func (s *service) destroy() {
	for _, m := range s.mods {
		m.Close()
		log.Infof("mod %s stop", m.Name())
	}
	for _, p := range s.plugin {
		p.Shut(s)
	}
	s.lookup.SetNotify(nil)
	s.lookup.Unregister(s.c.ID)
	s.lookup.Stop()
	s.exchange.Close()
	log.Infof("service %s shut", s.c.Name)
}

func (s *service) Close() {
	if s.handler.OnShut() {
		s.Shut()
		log.Infof("service %s close", s.c.Name)
	}
}

func (s *service) Shut() {
	if s.quit {
		return
	}
	s.quit = true
}
