package login

import (
	coreapi "turboengine/core/api"
	"turboengine/core/service"
)

// Service: 	Login
// Auth: 	 	sll
// Data:	  	2019-08-08 19:03:37
// Desc:
type Login struct {
	service.Service
}

func (s *Login) OnPrepare(srv coreapi.Service, args map[string]string) error {
	s.Service.OnPrepare(srv, args)
	// use plugin
	// use plugin end

	// add module
	// add module end

	return nil
}

func (s *Login) OnStart() error {
	return nil
}

func (s *Login) OnDependReady() {
	s.Ctx.Service().Ready()
}

func (s *Login) OnShut() bool {
	return true // If you want to close manually return false
}
