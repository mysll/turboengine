package api

import (
	"time"
	"turboengine/common/protocol"
	"turboengine/common/utils"
)

type Plugin interface {
	Prepare(Service)
	Shut(Service)
	Handle(cmd string, args ...interface{}) interface{}
}

type Call struct {
	Session  uint64
	DeadLine time.Time
	Callback func(*Call, []byte)
	UserData interface{}
	Err      error
}

type InvokeFn func(string, []byte) *protocol.Message
type Update func(*utils.Time)

type Service interface {
	AddModule(Module)
	Start() error
	Close()
	Shut()
	Attach(fn Update) uint64
	Deatch(id uint64)
	Pub(subject string, data []byte) error
	PubWithTimeout(subject string, data []byte, timeout time.Duration) (*Call, error)
	Sub(subject string, invoke InvokeFn) error
	SubNoInvoke(subject string) error
	UsePlugin(name string) error
	UnPlugin(name string)
	Plugin(name string) interface{}
	CallPlugin(plugin string, cmd string, args ...interface{}) (interface{}, error)
}

type ServiceHandler interface {
	OnPrepare(Service) error
	OnStart() error
	OnShut() bool
}
