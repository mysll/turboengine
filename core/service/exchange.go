package service

import (
	"errors"
	"fmt"
	"time"
	"turboengine/common/log"
	"turboengine/common/protocol"

	nats "github.com/nats-io/nats.go"
)

type Exchange struct {
	conn    *nats.Conn
	recvCh  chan *nats.Msg
	sendCh  chan *protocol.Message
	msgCh   chan *protocol.Message
	close   chan struct{}
	closing bool
	cancel  func()
	subs    map[string]*nats.Subscription
}

func NewExchange(recv chan *protocol.Message) (*Exchange, error) {
	p := &Exchange{
		recvCh: make(chan *nats.Msg, 128),
		msgCh:  recv,
		sendCh: make(chan *protocol.Message, 128),
		close:  make(chan struct{}),
		subs:   make(map[string]*nats.Subscription),
	}
	return p, nil
}

func (p *Exchange) Start(url string) error {
	conn, err := nats.Connect(url)
	if err != nil {
		return err
	}
	p.conn = conn

	go p.send()
	go p.input()
	return nil
}

func (p *Exchange) Close() {
	if p.closing {
		return
	}
	p.closing = true
	p.conn.Drain()
	p.conn.Close()
	close(p.close)
}

func (p *Exchange) Sub(subject string) error {
	if p.closing {
		return ERR_CLOSED
	}
	if _, dup := p.subs[subject]; dup {
		return errors.New("subject subscribe twice")
	}
	sub, err := p.conn.ChanSubscribe(subject, p.recvCh)
	if err != nil {
		return err
	}
	p.subs[subject] = sub
	return nil
}

func (p *Exchange) UnSub(subject string) {
	if sub, ok := p.subs[subject]; ok {
		sub.Unsubscribe()
		delete(p.subs, subject)
	}
}

func (p *Exchange) Pub(subject string, msg *protocol.Message) error {
	if p.closing {
		return ERR_CLOSED
	}
	msg.Header = msg.Header[:0]
	msg.Header = append(msg.Header, []byte(subject)...)
	select {
	case p.sendCh <- msg:
		return nil
	default:
		return ERR_MSG_TOO_MANY
	}
}

func (p *Exchange) input() {
L:
	for {
		select {
		case m := <-p.recvCh:
			msg := protocol.NewMessage(len(m.Data))
			msg.Header = msg.Header[:0]
			msg.Header = append(msg.Header, []byte(m.Subject)...)
			msg.Body = append(msg.Body, m.Data...)
			p.msgCh <- msg
			if p.closing {
				break L
			}
		case <-p.close:
			break L
		}
	}

	//drain
	for {
		select { // drain
		case m := <-p.recvCh:
			msg := protocol.NewMessage(len(m.Data))
			msg.Header = msg.Header[:0]
			msg.Header = append(msg.Header, []byte(m.Subject)...)
			msg.Body = append(msg.Body, m.Data...)
			p.msgCh <- msg
		default:
			return
		}
	}

}

func (p *Exchange) send() {
L:
	for {
		select {
		case m := <-p.sendCh:
			err := p.conn.Publish(string(m.Header), m.Body)
			log.Info("send:", string(m.Header))
			m.Free()
			if err != nil && !p.closing {
				for i := 0; i < 3; i++ {
					time.Sleep(time.Second)
					err = p.conn.Publish(string(m.Header), m.Body)
					if err == nil {
						break
					}
					log.Errorf("send message failed %s, retry after 1s", err.Error())
				}
			}
			if p.closing {
				break L
			}
		case <-p.close:
			break L
		}
	}

	// drain
	for {
		select {
		case m := <-p.sendCh:
			err := p.conn.Publish(string(m.Header), m.Body)
			fmt.Println("send:", string(m.Header), "body:", string(m.Body))
			m.Free()
			if err != nil {
				log.Error("send message error", err)
				continue
			}
		default:
			close(p.sendCh)
			return
		}
	}

}
