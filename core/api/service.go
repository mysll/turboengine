package api

import (
	"time"
	"turboengine/common/protocol"
	"turboengine/common/utils"
)

const (
	MB_TYPE_SERVICE = iota
	MB_TYPE_CONN
)

const (
	LOAD_BALANCE_RAND = iota
	LOAD_BALANCE_ROUND_ROBIN
	LOAD_BALANCE_LEAST_ACTIVE
	LOAD_BALANCE_HASH
)

const (
	INTEREST_CONNECTION_EVENT = iota + 1
	INTEREST_SERVICE_EVENT
)

var MAX_SID = 0x3FF

type Plugin interface {
	Prepare(srv Service, args ...interface{})
	Run()
	Shut(Service)
	Handle(cmd string, args ...interface{}) interface{}
}

type Locker interface {
	Lock()
	Unlock()
}

type Call struct {
	Session  uint64
	DeadLine time.Time
	Callback func(*Call)
	UserData interface{}
	Err      error
	Data     []byte
	Msg      *protocol.Message
	Done     chan *Call
}

type InvokeFn func(uint16, []byte) (*protocol.Message, error)
type Update func(*utils.Time)

type Service interface {
	ID() uint16
	Name() string
	Mailbox() protocol.Mailbox
	AddModule(Module)
	Start() error
	Close()
	Shut()
	Ready()
	Attach(fn Update) uint64
	Detach(id uint64)
	GenGUID() uint64
	Pub(subject string, data []byte) error
	PubWithTimeout(subject string, data []byte, timeout time.Duration) (*Call, error)
	Sub(subject string, invoke InvokeFn) error
	SubNoInvoke(subject string) error
	UnSub(subject string)
	UsePlugin(name string, args ...interface{}) error
	UnPlugin(name string)
	Plugin(name string) interface{}
	CallPlugin(plugin string, cmd string, args ...interface{}) (interface{}, error)
	Await()
	LookupById(id uint16) protocol.Mailbox
	LookupByName(name string) []protocol.Mailbox
	SelectService(name string, balance int, hash string) protocol.Mailbox
	SetProtoEncoder(enc protocol.ProtoEncoder)
	SetProtoDecoder(dec protocol.ProtoDecoder)
	SendToClient(dest protocol.Mailbox, msg *protocol.ProtoMsg) error
	OpenTransport()
	CloseTransport()
}

type ServiceHandler interface {
	OnPrepare(Service, map[string]string) error
	OnStart() error
	OnShut() bool
	OnDependReady()
	OnServiceAvailable(id uint16)
	OnServiceOffline(id uint16)
	OnConnected(session uint64)
	OnDisconnected(session uint64)
	OnMessage(*protocol.ProtoMsg)
}
