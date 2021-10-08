package proxy

import (
	"time"
	"turboengine/apps/login/api/rpc"
	"turboengine/common/protocol"
	"turboengine/core/api"
)

var login *Login

type Login struct {
	service api.Service
}

func (l *Login) Select(srv api.Service, service string, args string) protocol.Mailbox {
	dest := srv.SelectService(service, api.LOAD_BALANCE_HASH, args)
	return dest
}

func (l *Login) login(user, pass string) (bool, error) {
	login := rpc.NewLoginConsumerBySelector(l.service, "", l, time.Second*3)
	return login.Login(user, pass)
}

func InitLogin(service api.Service) {
	login = &Login{service: service}
}
