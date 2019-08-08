package service

import (
	"bufio"
	"fmt"
	"net"
	"runtime"
	"strings"
	"turboengine/common/log"
	"turboengine/common/protocol"
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
	ListenAndServe(addr string, port int, handler ConnHandler)
	Open()
	Close()
}

type TcpConn struct {
	conn   net.Conn
	reader *bufio.Reader
	writer *bufio.Writer
}

func (c *TcpConn) Addr() string {
	return c.conn.RemoteAddr().String()
}

func (c *TcpConn) Read(p []byte) (n int, err error) {
	return c.reader.Read(p)
}

func (c *TcpConn) Write(p []byte) (n int, err error) {
	return c.writer.Write(p)
}

func (c *TcpConn) Close() {
	c.conn.Close()
}

type TcpTransport struct {
	l     net.Listener
	open  bool
	ready chan bool
}

func (t *TcpTransport) Open() {
	if !t.open {
		t.open = true
		t.ready <- true
	}
}

func (t *TcpTransport) ListenAndServe(addr string, port int, handler ConnHandler) {
	l, err := net.Listen("tcp", fmt.Sprintf("%s:%d", addr, port))
	if err != nil {
		panic(err)
	}

	t.l = l
	log.Info("transport listen on ", l.Addr())

	if t.open {
		log.Info("transport waiting for connection")
	}

	for {
		if !t.open {
			log.Info("waiting for open")
			_, ok := <-t.ready
			if !ok { // close
				break
			}
			log.Info("transport waiting for connection")
		}

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
		go handler.Handle(&TcpConn{
			conn:   conn,
			reader: bufio.NewReaderSize(conn, protocol.MAX_BUF_LEN),
			writer: bufio.NewWriterSize(conn, protocol.MAX_BUF_LEN),
		})
	}
	log.Info("transport closed")
}

func (t *TcpTransport) Close() {
	if t.l != nil {
		t.l.Close()
		close(t.ready)
	}
}

func (s *service) UseTransport(typ string) {
	if tr, ok := transport[typ]; ok {
		s.tr = tr
		return
	}

	panic("transport not found " + typ)
}

func (s *service) createTransport(addr string, port int) {
	if s.tr == nil {
		s.tr = &TcpTransport{
			ready: make(chan bool, 1),
		}
	}

	s.wg.Wrap(func() {
		s.tr.ListenAndServe(addr, port, &NetHandle{svr: s})
	})
}

func (s *service) OpenTransport() {
	if s.tr != nil {
		s.tr.Open()
	}
}

func (s *service) CloseTransport() {
	if s.tr != nil {
		s.tr.Close()
		s.tr = nil
	}
}
