package global

import (
	"turboengine/apps/global/mod/global"
	"turboengine/core/api"
	coreapi "turboengine/core/api"
	"turboengine/core/module"
	"turboengine/core/plugin/election"
	"turboengine/core/service"
)

// Service: 	Global
// Auth: 	 	sll
// Data:	  	2019-08-14 17:05:45
// Desc:
type Global struct {
	service.Service
	globalData api.Module
}

func (s *Global) OnPrepare(srv coreapi.Service, args map[string]string) error {
	s.Service.OnPrepare(srv, args)
	// use plugin
	srv.UsePlugin(election.Name)
	// use plugin end

	// add module
	s.globalData = module.New(&global.GlobalData{}, false)
	s.globalData.SetInterest(api.INTEREST_SERVICE_EVENT)
	srv.AddModule(s.globalData)
	// add module end

	return nil
}

func (s *Global) OnStart() error {
	return nil
}

func (s *Global) OnDependReady() {
	s.Ctx.Service().Ready()
}

func (s *Global) OnShut() bool {
	return true // If you want to close manually return false
}
