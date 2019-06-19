package module

import (
	"context"
	"turboengine/common/utils"
	"turboengine/core/api"
)

type Module struct {
	Srv api.Service
	Ctx context.Context
}

func (m *Module) OnPrepare(s api.Service) error {
	m.Srv = s
	return nil
}

func (m *Module) OnStart(ctx context.Context) error {
	m.Ctx = ctx
	return nil
}

func (m *Module) OnUpdate(*utils.Time) {

}

func (m *Module) OnStop() error {
	return nil
}
