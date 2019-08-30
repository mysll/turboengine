package module

import (
	"context"
	"time"
	"turboengine/common/utils"
	"turboengine/core/api"
)

type Config struct {
	Name  string
	Async bool
	FPS   int
}

type module struct {
	srv      api.Service
	mod      api.ModuleHandler
	ctx      context.Context
	async    bool
	attachId uint64
	close    chan struct{}
	name     string
	fps      int
	interest int
}

func New(h api.ModuleHandler, async bool) api.Module {
	m := new(module)
	m.mod = h
	m.name = h.Name()
	m.close = make(chan struct{})
	m.async = async
	m.fps = 60
	return m
}

func NewWithConfig(h api.ModuleHandler, c Config) api.Module {
	m := new(module)
	m.mod = h
	m.name = h.Name()
	if m.name == "" {
		m.name = c.Name
	}
	m.close = make(chan struct{})
	m.async = c.Async
	m.fps = 60 // default
	if c.FPS != 0 {
		m.fps = c.FPS
	}
	return m
}

func (m *module) Handler() api.ModuleHandler {
	return m.mod
}

func (m *module) SetInterest(i int) {
	m.interest |= i
}

func (m *module) ClearInterest(i int) {
	m.interest &= ^i
}

func (m *module) Interest(i int) bool {
	return (m.interest & i) != 0
}

func (m *module) Name() string {
	return m.name
}

func (m *module) Init(srv api.Service) error {
	m.srv = srv
	m.mod.OnPrepare(srv)
	return nil
}

func (m *module) Start(ctx context.Context) {
	if m.async {
		go m.loop(ctx)
		return
	}
	m.run(ctx)
}

func (m *module) run(ctx context.Context) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	if err := m.mod.OnStart(ctx); err != nil {
		panic(err)
	}
	m.attachId = m.srv.Attach(m.Update)
}

func (m *module) loop(ctx context.Context) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	m.mod.OnStart(ctx)
	st := utils.NewTime(m.fps)
L:
	for {
		st.Update()
		m.Update(st)
		select {
		case <-m.close:
			break L
		case <-ctx.Done():
			break L
		default:
		}
		time.Sleep(st.NextFrame())
	}
}

func (m *module) Update(t *utils.Time) {
	m.mod.OnUpdate(t)
}

func (m *module) Close() {
	m.mod.OnStop()
	if m.attachId > 0 {
		m.srv.Detach(m.attachId)
	}
	close(m.close)
}
