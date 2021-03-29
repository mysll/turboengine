package service

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"
	"turboengine/common/log"
	"turboengine/common/protocol"
	"turboengine/common/utils"
	"turboengine/core/api"
	"turboengine/core/plugin/config"
	"turboengine/core/plugin/event"
	"turboengine/core/plugin/lock"

	_ "net/http/pprof"

	consulapi "github.com/hashicorp/consul/api"
	"github.com/mysll/toolkit"
)

type Dependency struct {
	Name  string
	Count int
}

type attach struct {
	id uint64
	fn api.Update
}

type service struct {
	sync.RWMutex
	wg      toolkit.WaitGroupWrapper
	id      uint16
	name    string
	mailbox protocol.Mailbox
	sid     string

	c            *Config
	handler      api.ServiceHandler
	running      bool
	quit         bool
	time         *utils.Time
	attaches     []attach
	attachId     uint64 // used for attaches
	mods         []api.Module
	exchange     *Exchange
	inMsg        chan *protocol.Message // receive message from message queue
	asyncMsg     chan *protocol.Message // receive async message from message queue
	lockCall     sync.RWMutex           // protect pending
	pending      map[uint64]*api.Call   // pending call
	asyncPending map[uint64]*api.Call   // async pending call
	session      uint64                 // used for pending
	delegate     map[string]api.InvokeFn
	plugin       map[string]api.Plugin
	lookup       *LookupService
	event        *event.Event
	disLocker    *lock.DisLocker
	config       *config.Configuration
	ready        bool
	uuid         uint64
	uuidTs       uint64
	transport    Transporter
	connPool     *ConnPool
	closing      bool
}

func New(h api.ServiceHandler, c *Config) api.Service {
	s := &service{
		c:            c,
		handler:      h,
		attaches:     make([]attach, 0, 8),
		attachId:     1,
		mods:         make([]api.Module, 0, 8),
		inMsg:        make(chan *protocol.Message, 512),
		asyncMsg:     make(chan *protocol.Message, 128),
		pending:      make(map[uint64]*api.Call),
		asyncPending: make(map[uint64]*api.Call),
		session:      0,
		delegate:     make(map[string]api.InvokeFn),
		plugin:       make(map[string]api.Plugin),
		lookup:       NewLookupService(consulapi.DefaultConfig()),
		event:        new(event.Event),
		disLocker:    new(lock.DisLocker),
		config:       new(config.Configuration),
	}

	if s.c != nil {
		s.id = s.c.ID
		s.name = s.c.Name
		s.sid = strconv.Itoa(int(c.ID))
		s.mailbox = protocol.GetServiceMailbox(s.id)
	}
	return s
}

func (s *service) ID() uint16 {
	return s.id
}

func (s *service) Name() string {
	return s.name
}

func (s *service) Mailbox() protocol.Mailbox {
	return s.mailbox
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

func (s *service) Start() error {
	toolkit.RandSeed()
	ctx := context.Background()
	if s.running {
		return fmt.Errorf("service %s already running", s.c.Name)
	}
	s.prepare()
	if err := s.handler.OnPrepare(s, s.c.Args); err != nil {
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
	if s.c.Debug {
		l, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", s.c.DebugPort))
		if err != nil {
			panic("debug error:" + err.Error())
		}
		go http.Serve(l, nil)
		log.Info("debug server start at ", l.Addr())
	}
	err := s.lookup.Register(s.sid, s.name, "127.0.0.1", s.c.DebugPort)
	if err != nil {
		log.Fatal(err)
	}
	s.lookup.Start()
	s.wg.Wrap(func() {
		s.run()
		s.destroy()
	})
	return nil
}

func (s *service) prepare() {
	s.usePlugin(event.Name, s.event)
	s.addEvent()
}

func (s *service) init() {
	err := s.lookup.Init()
	if err != nil {
		log.Fatal(err)
	}
	s.lookup.SetNotify(s.notify)
	s.usePlugin(lock.Name, s.disLocker, s.lookup.Client())
	s.usePlugin(config.Name, s.config, s.lookup.Client())
	p, err := NewExchange(s.inMsg, s.asyncMsg)
	if err != nil {
		log.Fatal(err)
	}

	s.exchange = p
	err = s.exchange.Start(s.c.NatsUrl)
	if err != nil {
		log.Fatal(err)
	}
	s.SubNoInvoke(fmt.Sprintf(DEFAULT_REPLY, s.c.ID)) // inner reply
	s.SubNoInvoke(fmt.Sprintf(ASYNC_REPLY, s.c.ID))   // async reply
	s.SubNoInvoke(SERVICE_SHUT)
	s.SubNoInvoke(SERVICE_SHUT_ALL)
	if s.c.Expose {
		s.connPool = NewConnPool(s)
		s.createTransport(s.c.Addr, s.c.Port)
	}
}

func (s *service) shutInvoke(*api.Call) {

}

func (s *service) Attach(fn api.Update) uint64 {
	id := s.attachId
	s.attachId++
	s.attaches = append(s.attaches, attach{
		id: id,
		fn: fn,
	})
	return id
}

func (s *service) Detach(id uint64) {
	for k, a := range s.attaches {
		if a.id == id {
			s.attaches = append(s.attaches[:k], s.attaches[k+1:]...)
			return
		}
	}
}

func (s *service) run() {
	log.Infof("service %s started", s.c.Name)
	fps := 10
	if s.c.FPS != 0 {
		fps = s.c.FPS
	}
	s.time = utils.NewTime(fps)
	for _, p := range s.plugin {
		p.Run()
	}
	if len(s.c.Depend) == 0 {
		s.handler.OnDependReady()
	}

	s.running = true

	for !s.quit {
		s.time.Update()
		s.input() // message queue a round trip
		if s.c.Expose {
			s.receive() // process client message
		}
		for _, f := range s.attaches {
			f.fn(s.time)
		}
		if s.c.FPS > 0 {
			time.Sleep(s.time.NextFrame())
		} else {
			time.Sleep(time.Millisecond)
		}
	}
}

func (s *service) destroy() {
	if s.c.Expose {
		// close transport
		s.CloseTransport()
		log.Info("kick all connections")
		// close all connections
		if s.connPool != nil {
			s.connPool.quit = true
			s.connPool.CloseAll()
		}
	}

	log.Info("stop modules")
	// close module
	for _, m := range s.mods {
		m.Close()
		log.Infof("mod %s stopped", m.Name())
	}

	// close plugin
	for k, p := range s.plugin {
		log.Infof("unplug %s plugin", k)
		p.Shut(s)
	}
	// close consul
	s.lookup.SetNotify(nil)
	log.Info("unregister service")
	s.lookup.Unregister(s.sid)
	s.lookup.Stop()
	// close message queue
	s.exchange.Close()
	log.Infof("service %s shut", s.c.Name)
}

func (s *service) Close() {
	if s.closing {
		return
	}
	if s.handler.OnShut() {
		s.Shut()
		log.Infof("service %s close", s.c.Name)
	}
	s.closing = true
}

func (s *service) Shut() {
	if s.quit {
		return
	}
	s.quit = true
}

func (s *service) Await() {
	s.wg.Wait()
}

func (s *service) Ready() {
	if s.ready {
		return
	}
	s.ready = true
	for _, m := range s.mods {
		m.Handler().OnReady()
	}
}

func (s *service) addEvent() {
	s.event.AddListener(EVENT_ADD, s.onServiceChange)
	s.event.AddListener(EVENT_DEL, s.onServiceChange)
	s.event.AddListener(EVENT_CONNECTED, s.onConnEvent)
	s.event.AddListener(EVENT_DISCONNECTED, s.onConnEvent)
}

func Capture() {
	f, err := os.OpenFile("./panic.log", os.O_CREATE|os.O_TRUNC|os.O_RDWR, os.ModePerm)
	if err != nil {
		panic(err)
	}
	redirectStderr(f)
}
