package fsm

import (
	"sync"
	"turboengine/common/utils"
	"turboengine/core/api"
	"turboengine/core/plugin"
)

const (
	Name = "FSM"
)

type State struct {
}

type FSM struct {
	sync.Mutex
	srv api.Service
	id  uint64
}

func (f *FSM) Prepare(srv api.Service, args ...any) {
	f.srv = srv
	f.id = f.srv.Attach(f.timer)
}

func (f *FSM) Shut(api.Service) {
	f.srv.Detach(f.id)
}

func (f *FSM) Run() {

}

func (f *FSM) Handle(cmd string, args ...any) any {
	return nil
}

func (f *FSM) timer(t *utils.Time) {

}

func init() {
	plugin.Register(Name, &FSM{})
}
