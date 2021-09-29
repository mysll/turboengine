package module

import (
	"context"
	"turboengine/common/protocol"
	"turboengine/common/utils"
	"turboengine/core/api"
)

type Module struct {
	Srv api.Service
	Ctx context.Context
}

func (m *Module) OnPrepare(s api.Service) error {
	m.Srv = s
	return nil
}

func (m *Module) OnStart(ctx context.Context) error {
	m.Ctx = ctx
	return nil
}

func (m *Module) OnUpdate(*utils.Time) {

}

func (m *Module) OnStop() error {
	return nil
}

func (m *Module) OnConnected(session protocol.Mailbox) {

}

func (m *Module) OnDisconnected(session protocol.Mailbox) {

}

func (m *Module) OnMessage(msg *protocol.ProtoMsg) {

}

func (m *Module) OnServiceAvailable(id uint16) {

}

func (m *Module) OnServiceOffline(id uint16) {

}

func (m *Module) OnReady() {

}
