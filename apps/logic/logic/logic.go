//@auth 	 	sll
//@create	  	2021-04-08 15:35:07
//@desc

package logic

import (
	_ "turboengine/apps/logic/internal/entity"
	coreapi "turboengine/core/api"
	"turboengine/core/plugin/storage"
	"turboengine/core/service"
)

type Logic struct {
	service.Service
}

func (s *Logic) OnPrepare(srv coreapi.Service, args map[string]string) error {
	s.Service.OnPrepare(srv, args)
	// use plugin
	srv.UsePlugin(storage.Name, "mysql", "root:123456@tcp(127.0.0.1:3306)/turbo?charset=utf8mb4&parseTime=True&loc=Local")
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
