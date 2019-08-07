package service

import (
	"strconv"
	"turboengine/core/api"
)

type Context struct {
	service api.Service
	flags   map[string]string
}

func (c *Context) Service() api.Service {
	return c.service
}

func (c *Context) String(key string) string {
	if val, ok := c.flags[key]; ok {
		return val
	}
	return ""
}

func (c *Context) Int(key string) int {
	if val, ok := c.flags[key]; ok {
		i, err := strconv.Atoi(val)
		if err != nil {
			return 0
		}
		return i
	}
	return 0
}

func (c *Context) Float(key string) float64 {
	if val, ok := c.flags[key]; ok {
		f, err := strconv.ParseFloat(val, 64)
		if err != nil {
			return 0
		}
		return f
	}
	return 0
}

type Service struct {
	Ctx Context
}

func (s *Service) OnPrepare(srv api.Service, flags map[string]string) error {
	s.Ctx = Context{
		service: srv,
		flags:   flags,
	}
	return nil
}

func (s *Service) OnStart() error {
	return nil
}

func (s *Service) OnShut() bool {
	return true
}

func (s *Service) OnDependReady() {

}

func (s *Service) OnServiceAvailable(id uint16) {

}

func (s *Service) OnServiceOffline(id uint16) {

}

func (s *Service) OnConnected(session uint64) {

}

func (s *Service) OnDisconnected(session uint64) {

}
