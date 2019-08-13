package client

import (
	"fmt"
	"net"
	"strings"
	"time"
	"turboengine/apps/gate/api/proto"
	"turboengine/common/protocol"
)

type Client struct {
	conn  net.Conn
	addr  string
	port  int
	login chan bool
}

func NewClient() *Client {
	c := &Client{
		login: make(chan bool),
	}
	return c
}

func (c *Client) Connect(addr string, port int) bool {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", addr, port))
	if err != nil {
		//fmt.Println(err)
		return false
	}

	c.conn = conn
	c.addr = addr
	c.port = port
	go c.read()
	return true
}

func (c *Client) Close() {
	c.conn.Close()
	c.conn = nil
}

func (c *Client) Send(msg *protocol.Message) bool {
	if err := protocol.WriteMsg(c.conn, msg.Body); err != nil {
		msg.Free()
		fmt.Println(err)
		return false
	}
	msg.Free()
	return true
}

func (c *Client) read() {
	for {
		body, err := protocol.ReadMsg(c.conn, protocol.MAX_MSG_LEN)
		if err != nil {
			if !strings.Contains(err.Error(), "use of closed network connection") {
				fmt.Println(err)
			}

			break
		}
		c.onMessage(body.Body)
		body.Free()
	}
}

func (c *Client) onMessage(body []byte) {
	dec := protocol.NewJsonDecoder()
	pt, err := dec.Decode(body)
	if err != nil {
		panic(err)
	}

	switch pt.Id {
	case proto.LOGIN_RESULT:
		result := pt.Data.(*proto.LoginResult).Result
		c.login <- result
	}
}

func (c *Client) Login(user, pass string) bool {
	login := &proto.Login{User: user, Pass: pass}
	enc := protocol.NewJsonEncoder()
	pm, err := enc.Encode(&protocol.ProtoMsg{
		Id:   proto.LOGIN,
		Data: login,
	})
	if err != nil {
		fmt.Println(err)
		return false
	}
	return c.Send(pm)
}

func (c *Client) WaitLogin() bool {
	t := time.NewTicker(time.Second * 30)
	select {
	case res := <-c.login:
		return res
	case <-t.C:
		return false
	}
}
