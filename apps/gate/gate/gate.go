package gate

import (
	"turboengine/apps/gate/mod/proxy"
	"turboengine/common/protocol"
	coreapi "turboengine/core/api"
	"turboengine/core/module"
	"turboengine/core/service"
)

// Service: 	Gate
// Auth: 	 	sll
// Data:	  	2019-08-07 11:21:24
// Desc:
type Gate struct {
	service.Service
	proxy *proxy.Proxy
}

func (s *Gate) OnPrepare(srv coreapi.Service, args map[string]string) error {
	s.Service.OnPrepare(srv, args)
	// use plugin
	// use plugin end

	// add module
	s.proxy = &proxy.Proxy{}
	srv.AddModule(module.New(s.proxy, false))
	// add module end

	return nil
}

func (s *Gate) OnStart() error {
	return nil
}

func (s *Gate) OnDependReady() {
	s.Ctx.Service().Ready()
	s.Ctx.Service().OpenTransport()
}

func (s *Gate) OnShut() bool {
	s.Ctx.Service().CloseTransport()
	return true // If you want to close manually return false
}

func (s *Gate) OnMessage(msg *protocol.ProtoMsg) {
	s.proxy.OnMessage(msg)
}
