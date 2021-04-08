package logic

import (
	coreapi "turboengine/core/api"
	"turboengine/core/service"
)

// Service: 	Logic
// Auth: 	 	sll
// Data:	  	2021-04-08 14:39:58
// Desc:
type Logic struct {
	service.Service
}

func (s *Logic) OnPrepare(srv coreapi.Service, args map[string]string) error {
	s.Service.OnPrepare(srv, args)
	// use plugin
	// use plugin end

	// add module
	// add module end

	return nil
}

func (s *Logic) OnStart() error {
	return nil
}

func (s *Logic) OnDependReady() {
	s.Ctx.Service().Ready()
}

func (s *Logic) OnShut() bool {
	return true // If you want to close manually return false
}
