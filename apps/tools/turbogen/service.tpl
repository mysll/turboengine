package {{tolower .Pkg}} 

import (
	coreapi "turboengine/core/api"
	"turboengine/core/service"
)
 
// Service: 	{{.Name}}
// Auth: 	 	{{.Auth}}
// Data:	  	{{.Time.Format "2006-01-02 15:04:05"}}
// Desc:
type {{.Name}} struct{
	service.Service
}

func (s *{{.Name}}) OnPrepare(srv coreapi.Service, args map[string]string) error {
	s.Service.OnPrepare(srv, args)
	// use plugin
	// use plugin end

	// add module
	// add module end 
	
	return nil
}

func (s *{{.Name}}) OnStart() error {
	return nil
}

func (s *{{.Name}}) OnDependReady() {
	s.Ctx.Service().Ready()
}

func (s *{{.Name}}) OnShut() bool {
	return true // If you want to close manually return false
}