//@author 	 	sll
//@create	  	2021-10-09 15:51:00
//@desc

package home

import (
	"context"
	"turboengine/common/utils"
	"turboengine/core/api"
	"turboengine/core/module"
)

type Home struct {
	module.Module
}

func (m *Home) Name() string {
	return "Home"
}

func (m *Home) OnPrepare(s api.Service) error {
	m.Module.OnPrepare(s)
	// load module resource
	// load module resource end

	return nil
}

func (m *Home) OnStart(ctx context.Context) error {
	m.Module.OnStart(ctx)
	// subscribe subject
	// subscribe subject end
	return nil
}

func (m *Home) OnUpdate(t *utils.Time) {

}

func (m *Home) OnStop() error {
	return nil
}
