package service

import (
	"fmt"
	"turboengine/common/log"
	"turboengine/core/api"
	"turboengine/core/plugin"
)

func (s *service) UsePlugin(name string) error {
	p := plugin.NewPlugin(name)
	if p != nil {
		s.usePlugin(name, p)
		return nil
	}

	return fmt.Errorf("plugin %s not found", name)
}

func (s *service) usePlugin(name string, p api.Plugin) {
	if _, ok := s.plugin[name]; ok {
		return
	}
	p.Prepare(s)
	s.plugin[name] = p
	log.Info("use plugin ", name)
}

func (s *service) UnPlugin(name string) {
	if _, ok := s.plugin[name]; ok {
		s.plugin[name].Shut(s)
		delete(s.plugin, name)
	}
}

func (s *service) Plugin(name string) interface{} {
	if p, ok := s.plugin[name]; ok {
		return p
	}
	return nil
}

func (s *service) CallPlugin(plugin string, cmd string, args ...interface{}) (interface{}, error) {
	if p, ok := s.plugin[plugin]; ok {
		return p.Handle(cmd, args...), nil
	}

	return nil, fmt.Errorf("plugin %s not found", plugin)
}
