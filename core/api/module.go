package api

import (
	"context"
	"turboengine/common/utils"
)

type Module interface {
	Name() string
	Init(srv Service) error
	Start(ctx context.Context)
	Close()
}

type ModuleHandler interface {
	Name() string // 模块名
	OnPrepare(Service) error
	OnStart(context.Context) error
	OnUpdate(*utils.Time)
	OnStop() error
}
