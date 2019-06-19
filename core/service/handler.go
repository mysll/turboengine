package service

import "turboengine/core/api"

type Service struct {
	Srv api.Service
}

func (s *Service) OnPrepare(srv api.Service) error {
	s.Srv = srv
	return nil
}

func (s *Service) OnStart() error {
	return nil
}

func (s *Service) OnShut() bool {
	return true
}
