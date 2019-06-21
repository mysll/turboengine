package math

import (
	"reflect"
	"turboengine/common/net/rpc"
	"turboengine/common/protocol"
)

type IMath_RPC_Go_V1_0 interface {
	Do(protocol.Mailbox, protocol.Mailbox, int, int) (int, error)
	Print(protocol.Mailbox, protocol.Mailbox, string) error
}

type Math_RPC_Go_V1_0 struct {
	handler IMath_RPC_Go_V1_0
}

type Math_RPC_Go_V1_0_Do struct {
	Arg0 protocol.Mailbox
	Arg1 protocol.Mailbox
	Arg2 int
	Arg3 int
}

type Math_RPC_Go_V1_0_Do_Reply struct {
	Arg0 int
}

type Math_RPC_Go_V1_0_Print struct {
	Arg0 protocol.Mailbox
	Arg1 protocol.Mailbox
	Arg2 string
}

type Math_RPC_Go_V1_0_Print_Reply struct {
}

func (math *Math_RPC_Go_V1_0) Do(arg *Math_RPC_Go_V1_0_Do, reply *Math_RPC_Go_V1_0_Do_Reply) (err error) {
	reply.Arg0, err = math.handler.Do(arg.Arg0, arg.Arg1, arg.Arg2, arg.Arg3)
	return
}

func (math *Math_RPC_Go_V1_0) Print(arg *Math_RPC_Go_V1_0_Print, reply *Math_RPC_Go_V1_0_Print_Reply) (err error) {
	err = math.handler.Print(arg.Arg0, arg.Arg1, arg.Arg2)
	return
}

func SetMathProvider(svr *rpc.Server, name string, provider IMath_RPC_Go_V1_0) {
	m := new(Math_RPC_Go_V1_0)
	m.handler = provider
	regName := "Math"
	if name != "" {
		regName = name
	}
	svr.RegisterName(regName+"_V1_0", m)
}

// client
type Math_RPC_Go_V1_0_Client struct {
	c   *rpc.Client
	srv string
}

func (m *Math_RPC_Go_V1_0_Client) Redirect(c *rpc.Client) {
	m.c = c
}

func (m *Math_RPC_Go_V1_0_Client) Do(arg0 protocol.Mailbox, arg1 protocol.Mailbox, arg2 int, arg3 int) (int, error) {
	_arg := &Math_RPC_Go_V1_0_Do{}
	_arg.Arg0 = arg0
	_arg.Arg1 = arg1
	_arg.Arg2 = arg2
	_arg.Arg3 = arg3

	_reply := &Math_RPC_Go_V1_0_Do_Reply{}
	err := m.c.Call(m.srv+"_V1_0.Do", _arg, _reply)
	return _reply.Arg0, err
}

func (m *Math_RPC_Go_V1_0_Client) Print(arg0 protocol.Mailbox, arg1 protocol.Mailbox, arg2 string) error {
	_arg := &Math_RPC_Go_V1_0_Print{}
	_arg.Arg0 = arg0
	_arg.Arg1 = arg1
	_arg.Arg2 = arg2

	_reply := &Math_RPC_Go_V1_0_Print_Reply{}
	err := m.c.Call(m.srv+"_V1_0.Print", _arg, _reply)
	return err
}

func NewMathConsumer(client *rpc.Client, srv string) *Math {
	m := new(Math)
	mc := new(Math_RPC_Go_V1_0_Client)
	mc.c = client
	mc.srv = "Math"
	if srv != "" {
		mc.srv = srv
	}
	m.XXX = mc
	value := reflect.ValueOf(m)

	value.Elem().FieldByName("Do").Set(reflect.ValueOf(mc.Do))
	value.Elem().FieldByName("Print").Set(reflect.ValueOf(mc.Print))
	return m
}
