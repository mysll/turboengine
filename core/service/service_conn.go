package service

import (
	"fmt"
	"sync"
	"turboengine/common/log"
	"turboengine/common/protocol"
	"turboengine/core/api"
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
	h.svr.event.AsyncEmit(EVENT_CONNECTED, n.session)
	if n == nil {
		conn.Close()
		return
	}
	go n.send()
	n.input(h.svr.connPool.inmsg)
	h.svr.event.AsyncEmit(EVENT_DISCONNECTED, n.session)
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
	buff := make([]byte, 0, protocol.MAX_BUF_LEN)
	for {
		data, err := protocol.ReadMsg(n.conn, buff[:0])
		if err != nil {
			log.Error(err)
			break
		}
		if decoder == nil {
			panic("decode is nil")
		}
		msg, err := decoder.Decode(data)
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

func (n *Node) send() {
	for m := range n.sendQueue {
		if encoder == nil {
			panic("encode is nil")
		}
		msg, err := encoder.Encode(m)
		if err != nil {
			log.Error("encode msg failed,", err)
			n.Close()
			break
		}
		if err := protocol.WriteMsg(n.conn, msg.Body); err != nil {
			log.Error("encode msg failed,", err)
			n.Close()
			msg.Free()
			break
		}
		msg.Free()
	}
}

func (n *Node) Send(m *protocol.ProtoMsg) error {
	select {
	case n.sendQueue <- m:
	default:
		return ERR_MSG_TOO_MANY
	}
	return nil
}

func (n *Node) Close() {
	if !n.closed {
		close(n.sendQueue)
		n.conn.Close()
		n.closed = true
	}
}

type ConnPool struct {
	sync.RWMutex
	svr     *service
	clients map[uint64]*Node
	session uint64
	quit    bool
	inmsg   chan *protocol.ProtoMsg
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
		sendQueue: make(chan *protocol.ProtoMsg, 128),
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
	return c.inmsg
}

func (s *service) onConnEvent(event string, args interface{}) {
	session := args.(uint64)
	switch event {
	case EVENT_CONNECTED:
		s.handler.OnConnected(session)
	case EVENT_DISCONNECTED:
		s.handler.OnDisconnected(session)
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

	return node.Send(msg)
}

func (s *service) receive() {
	for {
		select {
		case m := <-s.connPool.inmsg:
			s.handler.OnMessage(m)
		default:
			return
		}
	}
}
