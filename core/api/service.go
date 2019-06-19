package api

import (
	"turboengine/common/utils"
)

type AttachFn func(*utils.Time)

type Service interface {
	Register(Module)
	Start() error
	Close()
	Shut()
	Attach(fn AttachFn) int64
	Deatch(id int64)
}

type ServiceHandler interface {
	OnPrepare(Service) error
	OnStart() error
	OnShut() bool
}
