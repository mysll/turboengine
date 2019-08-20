package api

import (
	"context"
	"turboengine/common/protocol"
	"turboengine/common/utils"
)

type Module interface {
	Name() string
	Init(srv Service) error
	Start(ctx context.Context)
	Close()
	SetInterest(i int)
	ClearInterest(i int)
	Interest(i int) bool
	Handler() ModuleHandler
}

type ModuleHandler interface {
	Name() string // 模块名
	OnPrepare(Service) error
	OnStart(context.Context) error
	OnUpdate(*utils.Time)
	OnStop() error
	OnConnected(session uint64)
	OnDisconnected(session uint64)
	OnMessage(*protocol.ProtoMsg)
	OnServiceAvailable(id uint16)
	OnServiceOffline(id uint16)
	OnReady()
}
