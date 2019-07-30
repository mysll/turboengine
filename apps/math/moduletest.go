package math

import (
	"context"
	"turboengine/common/utils"
	"turboengine/core/api"
	"turboengine/core/module"
)

// Module: 		ModuleTest
// Auth: 	 	sll
// Data:	  	2019-07-30 11:31:31
// Desc:
type ModuleTest struct {
	module.Module
}

func (m *ModuleTest) OnPrepare(s api.Service) error {
	m.Module.OnPrepare(s)
	return nil
}

func (m *ModuleTest) OnStart(ctx context.Context) error {
	m.Module.OnStart(ctx)
	return nil
}

func (m *ModuleTest) OnUpdate(t *utils.Time) {

}

func (m *ModuleTest) OnStop() error {
	return nil
}
