package rpc

import (
	"bufio"
	"encoding/gob"
	"io"
	"log"
	"net"
	"reflect"
	"sync"
	"time"
)

type CallInfo struct {
	service *service
	mtype   *methodType
	req     *Request
	argv    reflect.Value
	replyv  reflect.Value
	sending *sync.Mutex
	next    *CallInfo
}

type ReplyInfo struct {
	req    *Request
	reply  interface{}
	errmsg string
}

type ChanServer struct {
	Server
	callCh  chan *CallInfo
	replyCh chan *ReplyInfo
	codec   ServerCodec
}

func NewChanServer(size int) *ChanServer {
	s := new(ChanServer)
	s.callCh = make(chan *CallInfo, size)
	s.replyCh = make(chan *ReplyInfo, size)
	return s
}

func (server *ChanServer) Close() {
	if server.codec != nil {
		server.codec.Close()
	}
	close(server.callCh)
	close(server.replyCh)
}

func (server *ChanServer) Accept(lis net.Listener) {
	for {
		conn, err := lis.Accept()
		if err != nil {
			log.Print("rpc.Serve: accept:", err.Error())
			return
		}
		go server.ServeConn(conn)
	}
}

func (server *ChanServer) ServeConn(conn net.Conn) {
	buf := bufio.NewWriter(conn)
	srv := &gobServerCodec{
		rwc:    conn,
		dec:    gob.NewDecoder(conn),
		enc:    gob.NewEncoder(buf),
		encBuf: buf,
	}
	go server.Write(srv)
	server.ServeCodec(srv)
}

func (server *ChanServer) ServeCodec(codec ServerCodec) {
	server.codec = codec
	for {
		service, mtype, req, argv, replyv, keepReading, err := server.readRequest(codec)
		if err != nil {
			if debugLog && err != io.EOF {
				log.Println("rpc:", err)
			}
			if !keepReading {
				break
			}
			// send a response if we actually managed to read a header.
			if req != nil {
				server.response(req, invalidRequest, codec, err.Error())
			}
			continue
		}
		c := &CallInfo{}
		c.service = service
		c.mtype = mtype
		c.req = req
		c.argv = argv
		c.replyv = replyv
		server.callCh <- c

	}
	codec.Close()
}

func (server *ChanServer) response(req *Request, reply interface{}, codec ServerCodec, errmsg string) {
	r := &ReplyInfo{}

	r.req = req
	r.reply = reply
	r.errmsg = errmsg
	server.replyCh <- r
}

// loop write
func (server *ChanServer) Write(codec ServerCodec) {
	for r := range server.replyCh {
		server.send(r.req, r.reply, codec, r.errmsg)
	}
}

// real send response
func (server *ChanServer) send(req *Request, reply interface{}, codec ServerCodec, errmsg string) {
	resp := server.getResponse()
	// Encode the response header
	resp.ServiceMethod = req.ServiceMethod
	if errmsg != "" {
		resp.Error = errmsg
		reply = invalidRequest
	}
	resp.Seq = req.Seq
	err := codec.WriteResponse(resp, reply)
	if debugLog && err != nil {
		log.Println("rpc: writing response:", err)
	}
	server.freeResponse(resp)
	server.freeRequest(req)
}

// one loop
func (server *ChanServer) Exec(timeout time.Duration) {
	st := time.Now()
L:
	for {
		select {
		case c, ok := <-server.callCh:
			if !ok {
				break L
			}
			c.service.chanCall(server, c.mtype, c.req, c.argv, c.replyv, server.codec)
		default:
			break L
		}
		if time.Now().Sub(st) > timeout {
			break L
		}
	}
}

func (s *service) chanCall(server *ChanServer, mtype *methodType, req *Request, argv, replyv reflect.Value, codec ServerCodec) {
	mtype.numCalls++
	function := mtype.method.Func
	// Invoke the method, providing a new value for the reply.
	returnValues := function.Call([]reflect.Value{s.rcvr, argv, replyv})
	// The return value for the method is an error.
	errInter := returnValues[0].Interface()
	errmsg := ""
	if errInter != nil {
		errmsg = errInter.(error).Error()
	}
	server.response(req, replyv.Interface(), codec, errmsg)
}
