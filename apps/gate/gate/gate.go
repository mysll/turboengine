package gate

import (
	"turboengine/apps/gate/mod/proxy"
	"turboengine/core/api"
	coreapi "turboengine/core/api"
	"turboengine/core/module"
	"turboengine/core/plugin/workqueue"
	"turboengine/core/service"
)

// Service: 	Gate
// Auth: 	 	sll
// Data:	  	2019-08-07 11:21:24
// Desc:
type Gate struct {
	service.Service
	proxy api.Module
}

func (s *Gate) OnPrepare(srv coreapi.Service, args map[string]string) error {
	s.Service.OnPrepare(srv, args)
	// use plugin
	srv.UsePlugin(workqueue.Name)
	// use plugin end

	// add module
	s.proxy = module.New(&proxy.Proxy{}, false)
	s.proxy.SetInterest(coreapi.INTEREST_CONNECTION_EVENT)
	srv.AddModule(s.proxy)
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
