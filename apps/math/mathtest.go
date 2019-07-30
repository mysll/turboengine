package math

import (
	"turboengine/core/api"
	"turboengine/core/service"
)

// Service: 	MathTest
// Auth: 	 	sll
// Data:	  	2019-07-30 11:19:59
// Desc:
type MathTest struct {
	service.Service
}

func (s *MathTest) OnPrepare(srv api.Service, args map[string]string) error {
	s.Service.OnPrepare(srv, args)
	return nil
}

func (s *MathTest) OnStart() error {
	return nil
}

func (s *MathTest) OnDependReady() {
}

func (s *MathTest) OnShut() bool {
	return true // If you want to close manually return false
}
