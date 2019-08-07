package gate

import (
	coreapi "turboengine/core/api"
	"turboengine/core/service"
)

// Service: 	Gate
// Auth: 	 	sll
// Data:	  	2019-08-07 11:21:24
// Desc:
type Gate struct {
	service.Service
}

func (s *Gate) OnPrepare(srv coreapi.Service, args map[string]string) error {
	s.Service.OnPrepare(srv, args)
	// use plugin
	// use plugin end

	// add module
	// add module end

	return nil
}

func (s *Gate) OnStart() error {
	return nil
}

func (s *Gate) OnDependReady() {
}

func (s *Gate) OnShut() bool {
	return true // If you want to close manually return false
}
