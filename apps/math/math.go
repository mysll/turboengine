package math

import (
	"turboengine/core/api"
	"turboengine/core/service"
)

// Service: 	Math
// Auth: 	 	sll
// Data:	  	2019-07-30 10:48:40
// Desc:
type Math struct {
	service.Service
}

func (s *Math) OnPrepare(srv api.Service, args map[string]string) error {
	s.Service.OnPrepare(srv, args)
	return nil
}

func (s *Math) OnStart() error {
	return nil
}

func (s *Math) OnDependReady() {
}

func (s *Math) OnShut() bool {
	return true // If you want to close manually return false
}
