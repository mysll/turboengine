package proxy

import (
	"turboengine/apps/gate/api/proto"
	"turboengine/common/log"
	"turboengine/common/protocol"
)

type User struct {
	proxy       *Proxy
	mailbox     protocol.Mailbox
	userAccount string
	logging     bool // 正在登录
	logged      bool // 是否已经登录
	userId      uint64
}

func NewUser(mailbox protocol.Mailbox, proxy *Proxy) *User {
	return &User{
		proxy:   proxy,
		mailbox: mailbox,
	}
}

func (user *User) OnMessage(msg *protocol.ProtoMsg) {
	switch msg.Id {
	case proto.LOGIN:
		loginMsg := msg.Data.(*proto.Login)
		user.login(loginMsg.User, loginMsg.Pass)
	}
}

func (user *User) login(account, pass string) {
	res, err := login.login(account, pass)
	if err != nil {
		log.Error(err)
	}
	outMsg := &protocol.ProtoMsg{
		Id:   proto.LOGIN_RESULT,
		Dest: user.mailbox,
		Data: proto.LoginResult{
			Result: res,
		},
	}
	user.proxy.Srv.SendToClient(user.mailbox, outMsg)
}
