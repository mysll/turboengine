package proxy

import (
	"context"
	"turboengine/apps/gate/api/proto"
	"turboengine/common/log"
	"turboengine/common/protocol"
	"turboengine/common/utils"
	"turboengine/core/api"
	"turboengine/core/module"
	"turboengine/core/plugin/lock"
	"turboengine/core/plugin/workqueue"
)

// Module: 		Proxy
// Auth: 	 	sll
// Data:	  	2019-08-08 15:16:00
// Desc:
type Proxy struct {
	module.Module
	workqueue *workqueue.WorkQueue
	dislock   *lock.DisLocker
}

func (m *Proxy) Name() string {
	return "Proxy"
}

func (m *Proxy) OnPrepare(s api.Service) error {
	m.Module.OnPrepare(s)
	// load module resource
	s.SetProtoEncoder(protocol.NewJsonEncoder())
	s.SetProtoDecoder(protocol.NewJsonDecoder())
	// load module resource end

	return nil
}

func (m *Proxy) OnStart(ctx context.Context) error {
	m.Module.OnStart(ctx)
	m.workqueue = m.Srv.Plugin(workqueue.Name).(*workqueue.WorkQueue)
	m.dislock = m.Srv.Plugin(lock.Name).(*lock.DisLocker)

	// subscribe subject
	// subscribe subject end
	return nil
}

func (m *Proxy) OnUpdate(t *utils.Time) {

}

func (m *Proxy) OnStop() error {
	return nil
}

func (m *Proxy) OnConnected(session uint64) {
	log.Info("new client")
}

func (m *Proxy) OnDisconnected(session uint64) {
	log.Info("remove client")
}

func (m *Proxy) OnMessage(msg *protocol.ProtoMsg) {
	log.Infof("recv msg %d, from %s", msg.Id, msg.Src)
	switch msg.Id {
	case proto.LOGIN:
		login := msg.Data.(*proto.Login)
		task := &Login{l: login, proxy: m, m: msg}
		if !m.Schedule(login.User, task) {
			task.SendResult(false)
		}
	}
}

func (m *Proxy) Schedule(key string, task workqueue.Task) bool {
	return m.workqueue.Schedule(utils.Hash64(key), task)
}
