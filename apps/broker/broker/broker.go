//@author 	 	sll
//@create	  	2021-10-09 15:43:46
//@desc

package broker

import (
	"turboengine/apps/broker/mod/home"
	coreapi "turboengine/core/api"
	"turboengine/core/module"
	"turboengine/core/plugin/storage"
	"turboengine/core/service"
)

type Broker struct {
	service.Service
	home coreapi.Module
}

func (s *Broker) OnPrepare(srv coreapi.Service, args map[string]string) error {
	s.Service.OnPrepare(srv, args)
	// use plugin
	srv.UsePlugin(storage.Name, "mysql", "root:123456@tcp(127.0.0.1:3306)/turbo?charset=utf8mb4&parseTime=True&loc=Local")
	// use plugin end

	// add module
	s.home = module.New(&home.Home{}, true)
	srv.AddModule(s.home)
	// add module end

	return nil
}

func (s *Broker) OnStart() error {
	return nil
}

func (s *Broker) OnDependReady() {
	s.Ctx.Service().Ready()
}

func (s *Broker) OnShut() bool {
	return true // If you want to close manually return false
}
