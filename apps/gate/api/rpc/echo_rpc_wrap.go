package rpc

import (
	"fmt"
	"time"
	"turboengine/apps/gate/api/proto"
	"turboengine/common/protocol"
	coreapi "turboengine/core/api"
)

type IEcho_RPC_Go_V1_0_0 interface {
	Print(string) error
	Echo(string) (string, error)
}

type Echo_RPC_Go_V1_0_0 struct {
	handler IEcho_RPC_Go_V1_0_0
}

func (p *Echo_RPC_Go_V1_0_0) Print(id uint16, data []byte) (ret *protocol.Message, err error) {
	ar := protocol.NewLoadArchive(data)

	var arg0 string
	err = ar.Get(&arg0)
	if err != nil {
		return
	}
	err1 := p.handler.Print(arg0)
	if err1 != nil {
		err = err1
		return
	}

	return
}

func (p *Echo_RPC_Go_V1_0_0) Echo(id uint16, data []byte) (ret *protocol.Message, err error) {
	ar := protocol.NewLoadArchive(data)

	var arg0 string
	err = ar.Get(&arg0)
	if err != nil {
		return
	}
	reply0, err1 := p.handler.Echo(arg0)
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

func SetEchoProvider(svr coreapi.Service, prefix string, provider IEcho_RPC_Go_V1_0_0) error {
	m := new(Echo_RPC_Go_V1_0_0)
	m.handler = provider

	if err := svr.Sub(fmt.Sprintf("%s%d:Echo.Print", prefix, svr.ID()), m.Print); err != nil {
		return err
	}
	if err := svr.Sub(fmt.Sprintf("%s%d:Echo.Echo", prefix, svr.ID()), m.Echo); err != nil {
		return err
	}
	return nil
}

// client
type Echo_RPC_Go_V1_0_0_Client struct {
	svr      coreapi.Service
	prefix   string
	dest     protocol.Mailbox
	timeout  time.Duration
	selector coreapi.Selector
}

func (m *Echo_RPC_Go_V1_0_0_Client) Redirect(dest protocol.Mailbox) {
	m.dest = dest
}

func (m *Echo_RPC_Go_V1_0_0_Client) SetSelector(selector coreapi.Selector) {
	m.selector = selector
}

// Print
func (m *Echo_RPC_Go_V1_0_0_Client) Print(arg0 string) (err error) {
	sr := protocol.NewAutoExtendArchive(128)
	err = sr.Put(arg0)
	if err != nil {
		return
	}

	msg := sr.Message()
	remote := m.dest
	if remote.IsNil() {
		remote = m.selector.Select(m.svr, "Gate", arg0)
	}
	if remote.IsNil() {
		err = fmt.Errorf("service Gate not found")
		return
	}
	call, err := m.svr.AsyncPubWithTimeout(fmt.Sprintf("%s%d:Echo.Print", m.prefix, remote.ServiceId()), msg.Body, m.timeout)
	msg.Free()
	if err != nil {
		return
	}
	call = <-call.Done
	if call.Err != nil {
		err = call.Err
		return
	}

	if call.Msg != nil {
		call.Msg.Free()
		call.Msg = nil
	}
	return
}

// Echo
func (m *Echo_RPC_Go_V1_0_0_Client) Echo(arg0 string) (reply0 string, err error) {
	sr := protocol.NewAutoExtendArchive(128)
	err = sr.Put(arg0)
	if err != nil {
		return
	}

	msg := sr.Message()
	remote := m.dest
	if remote.IsNil() {
		remote = m.selector.Select(m.svr, "Gate", arg0)
	}
	if remote.IsNil() {
		err = fmt.Errorf("service Gate not found")
		return
	}
	call, err := m.svr.AsyncPubWithTimeout(fmt.Sprintf("%s%d:Echo.Echo", m.prefix, remote.ServiceId()), msg.Body, m.timeout)
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

func NewEchoConsumer(svr coreapi.Service, prefix string, dest protocol.Mailbox, selector coreapi.Selector, timeout time.Duration) *proto.Echo {
	m := new(proto.Echo)
	mc := new(Echo_RPC_Go_V1_0_0_Client)
	mc.svr = svr
	mc.dest = dest
	mc.prefix = prefix
	mc.timeout = timeout
	mc.selector = selector
	m.XXX = mc
	m.Print = mc.Print
	m.Echo = mc.Echo
	return m
}

func NewEchoConsumerBySelector(svr coreapi.Service, prefix string, selector coreapi.Selector, timeout time.Duration) *proto.Echo {
	return NewEchoConsumer(svr, prefix, 0, selector, timeout)
}

func NewEchoConsumerByMailbox(svr coreapi.Service, prefix string, remote protocol.Mailbox, timeout time.Duration) *proto.Echo {
	return NewEchoConsumer(svr, prefix, remote, nil, timeout)
}

type Echo_RPC_Go_V1_0_0_Print_Reply struct {
}

type Echo_RPC_Go_V1_0_0_Echo_Reply struct {
	Arg0 string
}

type IEcho_RPC_Go_V1_0_0_Handler interface {
	OnPrint(error)
	OnEcho(string, error)
}

type Echo_RPC_Go_V1_0_0_Client_Handle struct {
	svr      coreapi.Service
	prefix   string
	dest     protocol.Mailbox
	timeout  time.Duration
	handler  IEcho_RPC_Go_V1_0_0_Handler
	selector coreapi.Selector
}

func (m *Echo_RPC_Go_V1_0_0_Client_Handle) Redirect(dest protocol.Mailbox) {
	m.dest = dest
}

func (m *Echo_RPC_Go_V1_0_0_Client_Handle) SetSelector(selector coreapi.Selector) {
	m.selector = selector
}

func (m *Echo_RPC_Go_V1_0_0_Client_Handle) Print(arg0 string) (err error) {
	sr := protocol.NewAutoExtendArchive(128)
	err = sr.Put(arg0)
	if err != nil {
		return
	}

	msg := sr.Message()
	remote := m.dest
	if remote.IsNil() {
		remote = m.selector.Select(m.svr, "Gate", arg0)
	}
	if remote.IsNil() {
		err = fmt.Errorf("service Gate not found")
		return
	}
	call, err := m.svr.PubWithTimeout(fmt.Sprintf("%s%d:Echo.Print", m.prefix, remote.ServiceId()), msg.Body, m.timeout)
	msg.Free()
	if err != nil {
		return
	}
	call.Callback = m.OnPrint
	return
}

func (m *Echo_RPC_Go_V1_0_0_Client_Handle) OnPrint(call *coreapi.Call) {

	var err error
	err = call.Err
	if err != nil {
		m.handler.OnPrint(err)
		return
	}

	m.handler.OnPrint(err)
}

func (m *Echo_RPC_Go_V1_0_0_Client_Handle) Echo(arg0 string) (reply0 string, err error) {
	sr := protocol.NewAutoExtendArchive(128)
	err = sr.Put(arg0)
	if err != nil {
		return
	}

	msg := sr.Message()
	remote := m.dest
	if remote.IsNil() {
		remote = m.selector.Select(m.svr, "Gate", arg0)
	}
	if remote.IsNil() {
		err = fmt.Errorf("service Gate not found")
		return
	}
	call, err := m.svr.PubWithTimeout(fmt.Sprintf("%s%d:Echo.Echo", m.prefix, remote.ServiceId()), msg.Body, m.timeout)
	msg.Free()
	if err != nil {
		return
	}
	call.Callback = m.OnEcho
	return
}

func (m *Echo_RPC_Go_V1_0_0_Client_Handle) OnEcho(call *coreapi.Call) {
	var reply Echo_RPC_Go_V1_0_0_Echo_Reply
	var err error
	err = call.Err
	if err != nil {
		m.handler.OnEcho(reply.Arg0, err)
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
	m.handler.OnEcho(reply.Arg0, err)
}

func NewEchoConsumerWithHandle(svr coreapi.Service, prefix string, dest protocol.Mailbox, selector coreapi.Selector, timeout time.Duration, handler IEcho_RPC_Go_V1_0_0_Handler) *proto.Echo {
	m := new(proto.Echo)
	mc := new(Echo_RPC_Go_V1_0_0_Client_Handle)
	mc.svr = svr
	mc.dest = dest
	mc.prefix = prefix
	mc.timeout = timeout
	mc.handler = handler
	mc.selector = selector
	m.XXX = mc
	m.Print = mc.Print
	m.Echo = mc.Echo
	return m
}

func NewEchoConsumerWithHandleBySelector(svr coreapi.Service, prefix string, selector coreapi.Selector, timeout time.Duration, handler IEcho_RPC_Go_V1_0_0_Handler) *proto.Echo {
	return NewEchoConsumerWithHandle(svr, prefix, 0, selector, timeout, handler)
}

func NewEchoConsumerWithHandleByMailbox(svr coreapi.Service, prefix string, remote protocol.Mailbox, timeout time.Duration, handler IEcho_RPC_Go_V1_0_0_Handler) *proto.Echo {
	return NewEchoConsumerWithHandle(svr, prefix, remote, nil, timeout, handler)
}
