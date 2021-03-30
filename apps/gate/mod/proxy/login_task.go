package proxy

import (
	"time"
	"turboengine/apps/gate/api/proto"
	"turboengine/apps/login/api/rpc"
	"turboengine/common/log"
	"turboengine/common/protocol"
	"turboengine/core/api"
)

type Login struct {
	l      *proto.Login
	m      *protocol.ProtoMsg
	proxy  *Proxy
	result bool
}

func (l *Login) Select(srv api.Service, service string, args string) protocol.Mailbox {
	dest := srv.SelectService(service, api.LOAD_BALANCE_HASH, args)
	return dest
}

func (l *Login) Run() {
	login := rpc.NewLoginConsumerBySelector(l.proxy.Srv, "", l, time.Second*3)
	res, err := login.Login(l.l.User, l.l.Pass)
	if err != nil {
		log.Error(err)
	}
	l.result = res
}

func (l *Login) Complete() {
	l.SendResult(l.result)
}

func (l *Login) SendResult(result bool) {
	proto := &protocol.ProtoMsg{
		Id:   proto.LOGIN_RESULT,
		Dest: l.m.Src,
		Data: proto.LoginResult{
			Result: result,
		},
	}
	l.proxy.Srv.SendToClient(l.m.Src, proto)
}
