package service

import (
	"errors"
	"fmt"
	"io"
	"strings"
	"sync"
	"turboengine/common/log"
	"turboengine/common/protocol"
	"turboengine/core/api"
)

const (
	PRE_ROUND_IN_MSG_MAX_PROCESS_COUNT = 128
	IN_MSG_LIST_MAX_COUNT              = 128
	NODE_SEND_MSG_COUNT                = 128
)

var (
	EVENT_CONNECTED    = "event_connected"
	EVENT_DISCONNECTED = "event_disconnected"
)

var (
	encoder protocol.ProtoEncoder
	decoder protocol.ProtoDecoder
)

type NetHandle struct {
	svr *service
}

func (h *NetHandle) Handle(conn Conn) {
	n := h.svr.connPool.NewNode(conn)
	if n == nil {
		conn.Close()
		return
	}
	h.svr.event.Emit(EVENT_CONNECTED, n.mailbox)
	go n.send()
	n.input(h.svr.connPool.inMsg)
	h.svr.event.Emit(EVENT_DISCONNECTED, n.mailbox)
	h.svr.connPool.RemoveNode(n.session, true)
}

type Node struct {
	conn      Conn
	session   uint64
	mailbox   protocol.Mailbox
	sendQueue chan *protocol.ProtoMsg
	closed    bool
}

func (n *Node) input(inmsg chan *protocol.ProtoMsg) {
	for {
		data, err := protocol.ReadMsg(n.conn, protocol.MAX_MSG_LEN)
		if err != nil {
			if err != io.EOF && !strings.Contains(err.Error(), "closed by the remote host") {
				log.Error(err)
			}
			break
		}
		if decoder == nil {
			panic("decode is nil")
		}
		msg, err := decoder.Decode(data.Body)
		data.Free()
		if err != nil {
			log.Error("decode msg failed,", err)
			n.Close()
			break
		}
		if n.closed {
			break
		}
		msg.Src = n.mailbox
		inmsg <- msg
	}
}

func (n *Node) write(m *protocol.ProtoMsg) error {
	if encoder == nil {
		panic("encode is nil")
	}
	msg, err := encoder.Encode(m)
	if err != nil {
		log.Error("encode msg failed,", err)
		n.Close()
		return err
	}
	if err = protocol.WriteMsg(n.conn, msg.Body); err != nil {
		if !strings.Contains(err.Error(), "closed by the remote host") {
			log.Error("write msg failed,", err)
		}
		n.Close()
		msg.Free()
		return err
	}
	msg.Free()
	return nil
}

func (n *Node) batchWrite() (err error) {
	for {
		select {
		case m, ok := <-n.sendQueue:
			if !ok {
				return errors.New("closed")
			}
			err = n.write(m)
			if err != nil {
				return err
			}
		default:
			return
		}
	}
}

func (n *Node) send() {
	var err error
L:
	for {
		select {
		case m, ok := <-n.sendQueue:
			if !ok {
				break L
			}
			err = n.write(m)
			if err != nil {
				break L
			}
			err = n.batchWrite()
			if err != nil {
				break L
			}
		}

		if err = n.conn.Flush(); err != nil {
			if !strings.Contains(err.Error(), "closed by the remote host") {
				log.Error("flush msg failed,", err)
			}
			n.Close()
			break
		}
	}
}

func (n *Node) Send(m *protocol.ProtoMsg) error {
	if n.closed {
		return ERR_CLOSED
	}
	select {
	case n.sendQueue <- m:
	default:
		return ERR_MSG_TOO_MANY
	}
	return nil
}

func (n *Node) Close() {
	if !n.closed {
		n.closed = true
		close(n.sendQueue)
		n.conn.Close()
	}
}

type ConnPool struct {
	sync.RWMutex
	svr     *service
	clients map[uint64]*Node
	session uint64
	quit    bool
	inMsg   chan *protocol.ProtoMsg
}

func NewConnPool(s *service) *ConnPool {
	p := &ConnPool{
		svr:     s,
		inMsg:   make(chan *protocol.ProtoMsg, IN_MSG_LIST_MAX_COUNT),
		clients: make(map[uint64]*Node),
	}
	return p
}

func (c *ConnPool) FindNode(session uint64) *Node {
	var n *Node
	c.RLock()
	if node, ok := c.clients[session]; ok {
		n = node
	}
	c.RUnlock()
	return n
}

func (c *ConnPool) NewNode(conn Conn) *Node {
	if c.quit {
		return nil
	}
	c.Lock()
	for {
		c.session++
		if c.session > protocol.ID_MAX {
			c.session = 1
		}
		if _, dup := c.clients[c.session]; dup {
			continue
		}
		break
	}

	n := &Node{
		conn:      conn,
		mailbox:   protocol.NewMailbox(c.svr.id, api.MB_TYPE_CONN, c.session),
		session:   c.session,
		sendQueue: make(chan *protocol.ProtoMsg, NODE_SEND_MSG_COUNT),
	}
	c.clients[n.session] = n
	log.Info("new session:", n.session, ",addr:", conn.Addr())
	c.Unlock()
	return n
}

func (c *ConnPool) RemoveNode(session uint64, close bool) {
	c.Lock()
	if node, ok := c.clients[session]; ok {
		if close {
			node.Close()
		}
		delete(c.clients, session)
		log.Info("remove session:", session, ",addr:", node.conn.Addr())
	}
	c.Unlock()
}

func (c *ConnPool) Close(session uint64) {
	c.RLock()
	if node, ok := c.clients[session]; ok {
		node.Close()
	}
	c.RUnlock()
}

func (c *ConnPool) CloseAll() {
	c.Lock()
	for _, node := range c.clients {
		node.Close()
	}
	c.Unlock()
}

func (c *ConnPool) Inmsg() chan *protocol.ProtoMsg {
	return c.inMsg
}

func (s *service) onConnEvent(event string, args any) {
	mailbox := args.(protocol.Mailbox)
	switch event {
	case EVENT_CONNECTED:
		s.handler.OnConnected(mailbox)
		for _, m := range s.mods {
			if m.Interest(api.INTEREST_CONNECTION_EVENT) {
				m.Handler().OnConnected(mailbox)
			}
		}
	case EVENT_DISCONNECTED:
		s.handler.OnDisconnected(mailbox)
		for _, m := range s.mods {
			if m.Interest(api.INTEREST_CONNECTION_EVENT) {
				m.Handler().OnDisconnected(mailbox)
			}
		}
	}
}

func (s *service) SetProtoEncoder(enc protocol.ProtoEncoder) {
	encoder = enc
}

func (s *service) SetProtoDecoder(dec protocol.ProtoDecoder) {
	decoder = dec
}

func (s *service) SendToClient(dest protocol.Mailbox, msg *protocol.ProtoMsg) error {
	if msg == nil {
		return fmt.Errorf("msg is nil")
	}

	node := s.connPool.FindNode(dest.Id())
	if node == nil {
		return fmt.Errorf("client not found, session:%d", dest.Id())
	}

	log.Infof("send msg to client, msg:%d, to: %s", msg.Id, dest)
	return node.Send(msg)
}

func (s *service) receive() {
	for i := 0; i < PRE_ROUND_IN_MSG_MAX_PROCESS_COUNT; i++ { // max loop PRE_ROUND_IN_MSG_MAX_PROCESS_COUNT
		select {
		case msg := <-s.connPool.inMsg:
			s.handler.OnMessage(msg)
			for _, m := range s.mods {
				if m.Interest(api.INTEREST_CONNECTION_EVENT) {
					m.Handler().OnMessage(msg)
				}
			}
		default:
			return
		}
	}
}
