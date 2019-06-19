package rpc

import (
	"fmt"
	"io"
)

type RpcFn0 func(src, dest Mailbox, args interface{}) error
type RpcFn1 func(src, dest Mailbox, args interface{}) (interface{}, error)

type RpcService interface {
	Export(string, interface{})
}

type Expose interface {
	Expose(RpcService)
}

type Call struct {
	session uint64
	srv     *service
}

type service struct {
	svr    *Server
	rcvr   interface{}
	method map[string]interface{}
}

// cb 原型：RpcFn0/RpcFn1
func (srv *service) Export(method string, cb interface{}) {
	switch cb.(type) {
	case RpcFn0:
	case RpcFn1:
	default:
		panic("register rpc method type error")
	}
}

type Server struct {
	codec  Codec
	srvMap map[string]*service
}

func NewServer(codec Codec) *Server {
	s := new(Server)
	s.codec = codec
	return s
}

// Register 注册service
func (svr *Server) Register(name string, rcvr Expose) error {
	if _, ok := svr.srvMap[name]; ok {
		return fmt.Errorf("service %s register twice", name)
	}

	s := new(service)
	s.svr = svr
	s.rcvr = rcvr
	rcvr.Expose(s)
	svr.srvMap[name] = s
	return nil
}

func (svr *Server) ServConn(conn io.ReadWriteCloser) {
	codec := &ByteCodec{
		rwc: conn,
	}
	svr.ServeCodec(codec)
}

func (svr *Server) ServeCodec(codec Codec) {
	svr.codec = codec

}
