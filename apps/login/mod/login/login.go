package login

import (
	"context"
	"turboengine/apps/login/api/rpc"
	"turboengine/common/utils"
	"turboengine/core/api"
	"turboengine/core/module"
	"turboengine/core/plugin/lock"
	"turboengine/core/plugin/storage"
)

// Module: 		Login
// Auth: 	 	sll
// Data:	  	2019-08-09 10:50:55
// Desc:
type Login struct {
	module.Module
	dislock *lock.DisLocker
	storege *storage.Storage
}

func (m *Login) Name() string {
	return "Login"
}

func (m *Login) OnPrepare(s api.Service) error {
	m.Module.OnPrepare(s)
	// load module resource
	m.storege = s.Plugin(storage.Name).(*storage.Storage)
	if m.storege == nil {
		panic("storage is nil")
	}
	// load module resource end

	return nil
}

func (m *Login) OnStart(ctx context.Context) error {
	m.Module.OnStart(ctx)
	m.dislock = m.Srv.Plugin(lock.Name).(*lock.DisLocker)
	// subscribe subject
	rpc.SetLoginProvider(m.Srv, "", &LoginServer{storage: m.storege})
	// subscribe subject end
	return nil
}

func (m *Login) OnUpdate(t *utils.Time) {

}

func (m *Login) OnStop() error {
	return nil
}
