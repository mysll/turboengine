package rpc

import (
	"bufio"
	"io"
	"strings"
	"sync"
	"time"
	"turboengine/common/log"
	"turboengine/common/protocol"

	"github.com/nggenius/ngengine/utils"
)

type Call struct {
	ServiceMethod string // The name of the service and method to call.
	Args          *protocol.Message
	Reply         *protocol.Message
	Error         error      // After completion, the error status.
	Done          chan *Call // Strobes when call is complete.
	Tag           interface{}
	deadline      time.Time
}

func (call *Call) done() {

}

type Client struct {
	mutex    sync.Mutex
	seq      uint64
	codec    ClientCodec
	pending  map[uint64]*Call
	sending  chan *Call
	response chan *Call

	closing  bool // user has called Close
	shutdown bool // server has told us to stop
}

func NewClient(conn io.ReadWriteCloser) *Client {
	encBuf := bufio.NewWriter(conn)
	client := &ByteClientCodec{conn, encBuf, RPC_MAX_LEN}
	return NewClientWithCodec(client)
}

// NewClientWithCodec is like NewClient but uses the specified
// codec to encode requests and decode responses.
func NewClientWithCodec(codec ClientCodec) *Client {
	client := &Client{
		codec:   codec,
		pending: make(map[uint64]*Call),
	}
	go client.input()
	go client.write()
	return client
}

func (client *Client) send(call *Call) (err error) {
	return
}

func (client *Client) write() {
	var err error
	t := time.NewTicker(time.Second)
L:
	for err == nil {
		select {
		case call := <-client.sending:
			err := client.send(call)
			if err != nil {
				log.Error("send error: %s", err.Error())
				break L
			}
		case <-t.C:
			client.mutex.Lock()
			now := time.Now()
			for k, v := range client.pending {
				if now.Sub(v.deadline) > 0 { // 超时删除
					delete(client.pending, k)
					if v.Args == nil {
						v.Args = protocol.NewMessage(1)
					}
					sr := utils.NewStoreArchiver(v.Args.Header)
					sr.Put(k)
					sr.Put(int8(1))
					sr.Put(int32(ERR_TIME_OUT))
					v.Args.Header = v.Args.Header[:sr.Len()]
					v.Error = ErrTimeout
					client.response <- v
					log.Info("response timeout, seq: %d", k)
				}
			}
			client.mutex.Unlock()
		default:
			if client.shutdown || client.closing {
				log.Info("quit sending loop")
				return
			}
			time.Sleep(time.Millisecond)
		}
	}
	t.Stop()
}

func (client *Client) input() {
	var err error
	for err == nil {
		message, err := client.codec.ReadMessage()
		if err != nil {
			if err != io.EOF &&
				!strings.Contains(err.Error(), "An existing connection was forcibly closed by the remote host") &&
				!strings.Contains(err.Error(), "use of closed network connection") {
				log.Error(client.codec.GetAddress(), err)
			}
			break
		}

		ar := utils.NewLoadArchiver(message.Header)
		seq, err := ar.GetUint64()
		client.mutex.Lock()
		call := client.pending[seq]
		delete(client.pending, seq)
		client.mutex.Unlock()

		switch {
		case call == nil:
		default:
			call.Reply = message
			client.response <- call
			log.Info("response replyed, seq: %d", seq)
		}
	}
	// Terminate pending calls.
	client.mutex.Lock()
	client.shutdown = true
	closing := client.closing
	if err == io.EOF {
		if closing {
			err = ErrShutdown
		} else {
			err = io.ErrUnexpectedEOF
		}
	}
	for _, call := range client.pending {
		call.Error = err
		call.done()
	}
	client.mutex.Unlock()

	log.Info("quit read loop")

}

func (client *Client) Process() {
	start_time := time.Now()
	for {
		select {
		case call := <-client.response:
			call.done()
			if time.Now().Sub(start_time) > time.Millisecond*500 { //消息太多,先返回,等下一帧再处理
				return
			}
		default:
			return
		}
	}
}

func (client *Client) Call(serviceMethod string, msg *protocol.Message) error {
	call := new(Call)
	call.ServiceMethod = serviceMethod
	call.Args = msg

	return nil
}
