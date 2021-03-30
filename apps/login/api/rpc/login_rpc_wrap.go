package rpc

import (
	"fmt"
	"time"
	"turboengine/apps/login/api/proto"
	"turboengine/common/protocol"
	coreapi "turboengine/core/api"
)

type ILogin_RPC_Go_V1_0_0 interface {
	Login(string, string) (bool, error)
}

type Login_RPC_Go_V1_0_0 struct {
	handler ILogin_RPC_Go_V1_0_0
}

func (p *Login_RPC_Go_V1_0_0) Login(id uint16, data []byte) (ret *protocol.Message, err error) {
	ar := protocol.NewLoadArchive(data)

	var arg0 string
	err = ar.Get(&arg0)
	if err != nil {
		return
	}
	var arg1 string
	err = ar.Get(&arg1)
	if err != nil {
		return
	}
	reply0, err1 := p.handler.Login(arg0, arg1)
	if err1 != nil {
		err = err1
		return
	}
	//reply
	sr := protocol.NewAutoExtendArchive(128)
	err = sr.Put(reply0)
	if err != nil {
		return
	}

	ret = sr.Message()

	return
}

func SetLoginProvider(svr coreapi.Service, prefix string, provider ILogin_RPC_Go_V1_0_0) error {
	m := new(Login_RPC_Go_V1_0_0)
	m.handler = provider

	if err := svr.Sub(fmt.Sprintf("%s%d:Login.Login", prefix, svr.ID()), m.Login); err != nil {
		return err
	}
	return nil
}

// client
type Login_RPC_Go_V1_0_0_Client struct {
	svr      coreapi.Service
	prefix   string
	dest     protocol.Mailbox
	timeout  time.Duration
	selector coreapi.Selector
}

func (m *Login_RPC_Go_V1_0_0_Client) Redirect(dest protocol.Mailbox) {
	m.dest = dest
}

func (m *Login_RPC_Go_V1_0_0_Client) SetSelector(selector coreapi.Selector) {
	m.selector = selector
}

// Login
func (m *Login_RPC_Go_V1_0_0_Client) Login(arg0 string, arg1 string) (reply0 bool, err error) {
	sr := protocol.NewAutoExtendArchive(128)
	err = sr.Put(arg0)
	if err != nil {
		return
	}
	err = sr.Put(arg1)
	if err != nil {
		return
	}

	msg := sr.Message()
	remote := m.dest
	if remote.IsNil() {
		remote = m.selector.Select(m.svr, "Login", arg0)
	}
	if remote.IsNil() {
		err = fmt.Errorf("service Login not found")
		return
	}
	call, err := m.svr.AsyncPubWithTimeout(fmt.Sprintf("%s%d:Login.Login", m.prefix, remote.ServiceId()), msg.Body, m.timeout)
	msg.Free()
	if err != nil {
		return
	}
	call = <-call.Done
	if call.Err != nil {
		err = call.Err
		return
	}

	for {
		ar := protocol.NewLoadArchive(call.Data)

		err = ar.Get(&reply0)
		if err != nil {
			break
		}
		break
	}

	if call.Msg != nil {
		call.Msg.Free()
		call.Msg = nil
	}
	return
}

func NewLoginConsumer(svr coreapi.Service, prefix string, dest protocol.Mailbox, selector coreapi.Selector, timeout time.Duration) *proto.Login {
	m := new(proto.Login)
	mc := new(Login_RPC_Go_V1_0_0_Client)
	mc.svr = svr
	mc.dest = dest
	mc.prefix = prefix
	mc.timeout = timeout
	mc.selector = selector
	m.XXX = mc
	m.Login = mc.Login
	return m
}

func NewLoginConsumerBySelector(svr coreapi.Service, prefix string, selector coreapi.Selector, timeout time.Duration) *proto.Login {
	return NewLoginConsumer(svr, prefix, 0, selector, timeout)
}

func NewLoginConsumerByMailbox(svr coreapi.Service, prefix string, remote protocol.Mailbox, timeout time.Duration) *proto.Login {
	return NewLoginConsumer(svr, prefix, remote, nil, timeout)
}

type Login_RPC_Go_V1_0_0_Login_Reply struct {
	Arg0 bool
}

type ILogin_RPC_Go_V1_0_0_Handler interface {
	OnLogin(bool, error)
}

type Login_RPC_Go_V1_0_0_Client_Handle struct {
	svr      coreapi.Service
	prefix   string
	dest     protocol.Mailbox
	timeout  time.Duration
	handler  ILogin_RPC_Go_V1_0_0_Handler
	selector coreapi.Selector
}

func (m *Login_RPC_Go_V1_0_0_Client_Handle) Redirect(dest protocol.Mailbox) {
	m.dest = dest
}

func (m *Login_RPC_Go_V1_0_0_Client_Handle) SetSelector(selector coreapi.Selector) {
	m.selector = selector
}

func (m *Login_RPC_Go_V1_0_0_Client_Handle) Login(arg0 string, arg1 string) (reply0 bool, err error) {
	sr := protocol.NewAutoExtendArchive(128)
	err = sr.Put(arg0)
	if err != nil {
		return
	}
	err = sr.Put(arg1)
	if err != nil {
		return
	}

	msg := sr.Message()
	remote := m.dest
	if remote.IsNil() {
		remote = m.selector.Select(m.svr, "Login", arg0)
	}
	if remote.IsNil() {
		err = fmt.Errorf("service Login not found")
		return
	}
	call, err := m.svr.PubWithTimeout(fmt.Sprintf("%s%d:Login.Login", m.prefix, remote.ServiceId()), msg.Body, m.timeout)
	msg.Free()
	if err != nil {
		return
	}
	call.Callback = m.OnLogin
	return
}

func (m *Login_RPC_Go_V1_0_0_Client_Handle) OnLogin(call *coreapi.Call) {
	var reply Login_RPC_Go_V1_0_0_Login_Reply
	var err error
	err = call.Err
	if err != nil {
		m.handler.OnLogin(reply.Arg0, err)
		return
	}

	for {
		ar := protocol.NewLoadArchive(call.Data)

		err = ar.Get(&reply.Arg0)
		if err != nil {
			break
		}
		break
	}
	m.handler.OnLogin(reply.Arg0, err)
}

func NewLoginConsumerWithHandle(svr coreapi.Service, prefix string, dest protocol.Mailbox, selector coreapi.Selector, timeout time.Duration, handler ILogin_RPC_Go_V1_0_0_Handler) *proto.Login {
	m := new(proto.Login)
	mc := new(Login_RPC_Go_V1_0_0_Client_Handle)
	mc.svr = svr
	mc.dest = dest
	mc.prefix = prefix
	mc.timeout = timeout
	mc.handler = handler
	mc.selector = selector
	m.XXX = mc
	m.Login = mc.Login
	return m
}

func NewLoginConsumerWithHandleBySelector(svr coreapi.Service, prefix string, selector coreapi.Selector, timeout time.Duration, handler ILogin_RPC_Go_V1_0_0_Handler) *proto.Login {
	return NewLoginConsumerWithHandle(svr, prefix, 0, selector, timeout, handler)
}

func NewLoginConsumerWithHandleByMailbox(svr coreapi.Service, prefix string, remote protocol.Mailbox, timeout time.Duration, handler ILogin_RPC_Go_V1_0_0_Handler) *proto.Login {
	return NewLoginConsumerWithHandle(svr, prefix, remote, nil, timeout, handler)
}
