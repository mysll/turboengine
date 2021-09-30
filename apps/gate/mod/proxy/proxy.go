//@auth 	 	sll
//@create	  	2019-08-08 15:16:00
//@desc

package proxy

import (
	"context"
	"turboengine/apps/gate/api/rpc"
	"turboengine/common/log"
	"turboengine/common/protocol"
	"turboengine/common/utils"
	"turboengine/core/api"
	"turboengine/core/module"
	"turboengine/core/plugin/workqueue"
)

type Proxy struct {
	module.Module
	workQueue *workqueue.WorkQueue
	users     map[protocol.Mailbox]*User
}

func (m *Proxy) Name() string {
	return "Proxy"
}

func (m *Proxy) OnPrepare(s api.Service) error {
	m.users = make(map[protocol.Mailbox]*User)
	m.Module.OnPrepare(s)
	// load module resource
	s.SetProtoEncoder(protocol.NewJsonEncoder())
	s.SetProtoDecoder(protocol.NewJsonDecoder())
	// load module resource end
	InitLogin(s)
	return nil
}

func (m *Proxy) OnStart(ctx context.Context) error {
	m.Module.OnStart(ctx)
	m.workQueue = m.Srv.Plugin(workqueue.Name).(*workqueue.WorkQueue)
	// subscribe subject
	rpc.SetPushProvider(m.Srv, "", &PushServer{m.Srv})
	// subscribe subject end
	return nil
}

func (m *Proxy) OnUpdate(t *utils.Time) {

}

func (m *Proxy) OnStop() error {
	return nil
}

func (m *Proxy) OnConnected(mailbox protocol.Mailbox) {
	m.users[mailbox] = NewUser(mailbox, m)
	log.Info("new client")
}

func (m *Proxy) OnDisconnected(mailbox protocol.Mailbox) {
	if user, ok := m.users[mailbox]; ok {
		user.onDisconnect()
	}
	delete(m.users, mailbox)
	log.Info("remove client")
}

func (m *Proxy) OnMessage(msg *protocol.ProtoMsg) {
	log.Infof("recv msg %d, from %s", msg.Id, msg.Src)
	if user, ok := m.users[msg.Src]; ok {
		user.OnMessage(msg)
	}
}
