//@author 	 	{{.Auth}}
//@create	  	{{.Time.Format "2006-01-02 15:04:05"}}
//@desc

package {{tolower .Pkg}}

import (
	"context"
	"turboengine/common/utils"
	"turboengine/core/api"
	"turboengine/core/module"
)

type {{.Name}} struct{
	module.Module
}

func (m *{{.Name}}) Name() string{
	return "{{.Name}}"
}

func (m *{{.Name}}) OnPrepare(s api.Service) error {
	m.Module.OnPrepare(s)
	// load module resource
	// load module resource end

	return nil
}

func (m *{{.Name}}) OnStart(ctx context.Context) error {
	m.Module.OnStart(ctx)
	// subscribe subject
	// subscribe subject end
	return nil
}

func (m *{{.Name}}) OnUpdate(t *utils.Time) {

}

func (m *{{.Name}}) OnStop() error {
	return nil
}
