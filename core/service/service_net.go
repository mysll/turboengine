package service

import (
	"fmt"
	"net"
	"runtime"
	"strings"
	"turboengine/common/log"
)

var transport = make(map[string]Transporter)

func RegisterTransport(typ string, tr Transporter) {
	if _, dup := transport[typ]; dup {
		panic("register transport twice" + typ)
	}

	transport[typ] = tr
}

type Conn interface {
	Addr() string
	Read(p []byte) (n int, err error)
	Write(p []byte) (n int, err error)
	Close()
}

type ConnHandler interface {
	Handle(Conn)
}

type Transporter interface {
	ListenAndServ(addr string, port int, handler ConnHandler)
	Close()
}

type TcpConn struct {
	conn net.Conn
}

func (c *TcpConn) Addr() string {
	return c.conn.RemoteAddr().String()
}

func (c *TcpConn) Read(p []byte) (n int, err error) {
	return c.conn.Read(p)
}

func (c *TcpConn) Write(p []byte) (n int, err error) {
	return c.conn.Write(p)
}

func (c *TcpConn) Close() {
	c.conn.Close()
}

type TcpTransport struct {
	l net.Listener
}

func (t *TcpTransport) ListenAndServ(addr string, port int, handler ConnHandler) {
	l, err := net.Listen("tcp", fmt.Sprintf("%s:%d", addr, port))
	if err != nil {
		panic(err)
	}

	t.l = l
	log.Info("open transport at ", l.Addr())

	for {
		conn, err := l.Accept()
		if err != nil {
			if nerr, ok := err.(net.Error); ok && nerr.Temporary() {
				log.Warnf("NOTICE: temporary Accept() failure - %s", err)
				runtime.Gosched()
				continue
			}
			if !strings.Contains(err.Error(), "use of closed network connection") {
				log.Errorf("ERROR: listener.Accept() - %s", err)
			}
			break
		}
		go handler.Handle(&TcpConn{conn: conn})
	}
	log.Info("close transport")
}

func (t *TcpTransport) Close() {
	if t.l != nil {
		t.l.Close()
	}
}

func (s *service) UseTransport(typ string) {
	if tr, ok := transport[typ]; ok {
		s.tr = tr
		return
	}

	panic("transport not found " + typ)
}

func (s *service) OpenTransport(addr string, port int) {
	if s.tr == nil {
		s.tr = &TcpTransport{}
	}

	s.wg.Wrap(func() {
		s.tr.ListenAndServ(addr, port, &NetHandle{svr: s})
	})
}
