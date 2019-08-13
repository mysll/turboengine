package service

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"runtime"
	"strings"
	"sync"
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
	Flush() error
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
	if c.reader != nil {
		return c.reader.Read(p)
	}
	err = fmt.Errorf("reader is nil")
	return
}

func (c *TcpConn) Write(p []byte) (n int, err error) {
	if c.writer != nil {
		return c.writer.Write(p)
	}
	err = fmt.Errorf("writer is nil")
	return
}

func (c *TcpConn) Flush() error {
	if c.writer != nil {
		return c.writer.Flush()
	}

	return fmt.Errorf("writer is nil")
}

func (c *TcpConn) Close() {
	if c.reader != nil {
		putBufioReader(c.reader)
		c.reader = nil
	}
	if c.writer != nil {
		c.writer.Flush()
		putBufioWriter(c.writer)
		c.writer = nil
	}
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

var (
	bufioReaderPool sync.Pool
	bufioWriterPool sync.Pool
)

func newBufioReader(r io.Reader) *bufio.Reader {
	if v := bufioReaderPool.Get(); v != nil {
		br := v.(*bufio.Reader)
		br.Reset(r)
		return br
	}
	return bufio.NewReader(r)
}

func putBufioReader(br *bufio.Reader) {
	br.Reset(nil)
	bufioReaderPool.Put(br)
}

func newBufioWriter(w io.Writer, size int) *bufio.Writer {
	if v := bufioWriterPool.Get(); v != nil {
		bw := v.(*bufio.Writer)
		bw.Reset(w)
		return bw
	}
	return bufio.NewWriterSize(w, size)
}

func putBufioWriter(bw *bufio.Writer) {
	bw.Reset(nil)
	bufioWriterPool.Put(bw)
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
				log.Warnf("temporary Accept() failure - %s", err)
				runtime.Gosched()
				continue
			}
			if !strings.Contains(err.Error(), "use of closed network connection") {
				log.Errorf("listener.Accept() - %s", err)
			}
			break
		}

		go handler.Handle(&TcpConn{
			conn:   conn,
			reader: newBufioReader(conn),
			writer: newBufioWriter(conn, 4<<10),
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
