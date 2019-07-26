package rpc

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"strings"
	"sync"
	"turboengine/common/log"
	"turboengine/common/protocol"
)

type Fn func(*protocol.Message) (*protocol.Message, error)

type Service interface {
	Register(name string, f Fn)
}

type Exposer interface {
	Expose(Service)
}

type Request struct {
	Session       uint64
	ServiceMethod string            // format: "Service.Method"
	Seq           uint64            // sequence number chosen by client
	Raw           *protocol.Message // raw message
	method        Fn
	service       *service
	err           error    // 出错的调用
	next          *Request // for free list in Server
}

func (r *Request) Invoke() {
	resp := r.service.server.getResponse()
	resp.Error = r.err
	resp.Errcode = 1
	resp.Session = r.Session
	resp.Seq = r.Seq
	if resp.Error == nil {
		reply, err := r.method(r.Raw) // real call function
		resp.Reply = reply
		resp.Error = err
		if err == nil {
			resp.Errcode = 0
		}
	}

	r.service.server.Reply(resp)
}

func (r *Request) Error() error {
	return r.err
}

func (r *Request) Done() {
	r.service.server.freeRequest(r) // 回收
}

type Response struct {
	Session uint64
	Seq     uint64
	Errcode int32
	Reply   *protocol.Message
	Error   error
	next    *Response
}

type service struct {
	server *Server
	rcvr   Exposer
	method map[string]Fn
}

func (s *service) Register(name string, f Fn) {
	if _, dup := s.method[name]; dup {
		log.Fatal("method %s register twice", name)
	}
	s.method[name] = f
}

type Session struct {
	rwc       io.ReadWriteCloser
	codec     ServerCodec
	sendQueue chan *Response
	quit      bool
}

func (s *Session) send() {
	for resp := range s.sendQueue {
		if resp.Seq == 0 {
			log.Error("resp.seq is zero")
			break
		}

		err := s.codec.WriteResponse(resp.Seq, resp.Errcode, resp.Reply)
		if err != nil {
			log.Error(err.Error())
			break
		}
		log.Info("response call, seq: %d", resp.Seq)
		if s.quit {
			break
		}
	}
}

func (s *Session) close() {
	close(s.sendQueue)
	s.quit = true
}

type Server struct {
	sessionLock sync.RWMutex // protects sessions
	serial      uint64
	sessions    map[uint64]*Session

	reqLock  sync.Mutex // protects freeReq
	freeReq  *Request
	respLock sync.Mutex // protects freeResp
	freeResp *Response

	serviceMap map[string]*service
	request    chan *Request
}

func NewServer() *Server {
	s := &Server{}
	s.serial = 1
	s.sessions = make(map[uint64]*Session)
	s.serviceMap = make(map[string]*service)
	return s
}

func (server *Server) getRequest() *Request {
	server.reqLock.Lock()
	req := server.freeReq
	if req == nil {
		req = new(Request)
	} else {
		server.freeReq = req.next
		*req = Request{}
	}
	server.reqLock.Unlock()
	return req
}

func (server *Server) freeRequest(req *Request) {
	server.reqLock.Lock()
	if req.Raw != nil {
		req.Raw.Free() // free message
		req.Raw = nil
	}
	req.next = server.freeReq
	server.freeReq = req
	server.reqLock.Unlock()
}

func (server *Server) getResponse() *Response {
	server.respLock.Lock()
	resp := server.freeResp
	if resp == nil {
		resp = new(Response)
	} else {
		server.freeResp = resp.next
		*resp = Response{}
	}
	server.respLock.Unlock()
	return resp
}

func (server *Server) freeResponse(resp *Response) {
	server.respLock.Lock()
	resp.next = server.freeResp
	server.freeResp = resp
	server.respLock.Unlock()
}

func (server *Server) Reply(resp *Response) {

}

func (s *Server) RegisterName(name string, rcvr Exposer) {
	if _, dup := s.serviceMap[name]; dup {
		log.Fatal("service %s register twice", name)
	}
	s.serviceMap[name] = &service{
		rcvr:   rcvr,
		method: make(map[string]Fn),
	}
}

func (s *Server) ServeConn(conn net.Conn) {
	codec := &ByteServerCodec{
		rwc:    conn,
		encBuf: bufio.NewWriter(conn),
	}
	s.ServeCodec(codec)
}

func (s *Server) ServeCodec(codec ServerCodec) {
	var serial uint64
	s.sessionLock.Lock()
	serial = s.serial
	s.serial++
	session := &Session{rwc: codec.GetConn(), codec: codec, sendQueue: make(chan *Response, 32)}
	s.sessions[serial] = session
	s.sessionLock.Unlock()
	go session.send()

	log.Info("start new rpc server %d", serial)
	for {
		req := s.getRequest()
		req.Session = serial
		err := codec.ReadRequest(req, RPC_MAX_LEN)
		if err != nil {
			if err != io.EOF &&
				!strings.Contains(err.Error(), "An existing connection was forcibly closed by the remote host") &&
				!strings.Contains(err.Error(), "use of closed network connection") {
				log.Error("rpc err:%s", err.Error())
			} else {
				log.Info("service client closed")
			}
			break
		}

		err = s.parseRequest(req)
		if err != nil {
			req.err = err
			log.Error(err.Error())
		}

		s.request <- req
	}

	session.close()
	codec.Close()
	s.sessionLock.Lock()
	delete(s.sessions, serial)
	log.Info("rpc server %d closed", serial)
	s.sessionLock.Unlock()
}

func (s *Server) parseRequest(r *Request) (err error) {

	dot := strings.LastIndex(r.ServiceMethod, ".")
	if dot < 0 {
		err = fmt.Errorf("rpc: service/method request ill-formed: %s", r.ServiceMethod)
		return
	}
	serviceName := r.ServiceMethod[:dot]
	methodName := r.ServiceMethod[dot+1:]

	r.service = s.serviceMap[serviceName]
	if r.service == nil {
		err = fmt.Errorf("rpc: can't find service %s", r.ServiceMethod)
		return
	}

	r.method = r.service.method[methodName]
	if r.method == nil {
		err = fmt.Errorf("rpc: can't find method %s", r.ServiceMethod)
		return
	}
	return
}
