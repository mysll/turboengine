package proxy

import (
	"context"
	"turboengine/apps/gate/internal/proto"
	"turboengine/common/protocol"
	"turboengine/common/utils"
	"turboengine/core/api"
	"turboengine/core/module"
)

// Module: 		Proxy
// Auth: 	 	sll
// Data:	  	2019-08-08 15:16:00
// Desc:
type Proxy struct {
	module.Module
}

func (m *Proxy) Name() string {
	return "Proxy"
}

func (m *Proxy) OnPrepare(s api.Service) error {
	m.Module.OnPrepare(s)
	// load module resource
	s.SetProtoEncoder(&proto.JsonEncoder{})
	s.SetProtoDecoder(&proto.JsonDecoder{})
	// load module resource end

	return nil
}

func (m *Proxy) OnStart(ctx context.Context) error {
	m.Module.OnStart(ctx)
	// subscribe subject
	// subscribe subject end
	return nil
}

func (m *Proxy) OnUpdate(t *utils.Time) {

}

func (m *Proxy) OnStop() error {
	return nil
}

func (m *Proxy) OnMessage(msg *protocol.ProtoMsg) {

}
