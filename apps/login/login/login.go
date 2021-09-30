//@auth 	 	sll
//@create	  	2019-08-08 19:03:37
//@desc

package login

import (
	"turboengine/apps/login/mod/login"
	"turboengine/core/api"
	coreapi "turboengine/core/api"
	"turboengine/core/module"
	"turboengine/core/plugin/storage"
	"turboengine/core/service"
)

type Login struct {
	service.Service
	login api.Module
}

func (s *Login) OnPrepare(srv coreapi.Service, args map[string]string) error {
	s.Service.OnPrepare(srv, args)
	// use plugin
	srv.UsePlugin(storage.Name, "mysql", "root:123456@tcp(127.0.0.1:3306)/turbo?charset=utf8mb4&parseTime=True&loc=Local")
	// use plugin end

	// add module
	s.login = module.New(&login.Login{}, false)
	srv.AddModule(s.login)
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
