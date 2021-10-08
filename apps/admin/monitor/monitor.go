//@auth 	 	sll
//@create	  	2019-08-07 09:34:16
//@desc

package monitor

import (
	coreapi "turboengine/core/api"
	"turboengine/core/service"
)

type Monitor struct {
	service.Service
}

func (s *Monitor) OnPrepare(srv coreapi.Service, args map[string]string) error {
	s.Service.OnPrepare(srv, args)
	// use plugin
	// use plugin end

	// add module
	// add module end

	return nil
}

func (s *Monitor) OnStart() error {
	go run(s.Ctx.Service())
	return nil
}

func (s *Monitor) OnDependReady() {
	s.Ctx.Service().Ready()
}

func (s *Monitor) OnShut() bool {
	return true // If you want to close manually return false
}
