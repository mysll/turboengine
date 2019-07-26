package service

import (
	"testing"
	"time"
	"turboengine/common/log"
	"turboengine/common/protocol"
	"turboengine/core/api"
)

type Echo struct {
	Service
}

func (e *Echo) OnStart() error {
	return nil
}

func Srv1() {
	s := New(&Echo{}, &Config{ID: "100", Name: "echo", NatsUrl: "nats://0.0.0.0:4222"})
	go s.Start()
	time.Sleep(time.Second)
	s.Sub("test", func(id string, data []byte) *protocol.Message {
		log.Info("call test from:", id, ", content:", string(data))
		m := protocol.NewMessage(len(data))
		m.Body = append(m.Body, data...)
		return m
	})
	time.Sleep(time.Second * 5)
	s.Close()
}

func Srv2() {
	s := New(&Echo{}, &Config{ID: "101", Name: "echo", NatsUrl: "nats://0.0.0.0:4222"})
	go s.Start()
	time.Sleep(time.Second * 2)
	c, err := s.PubWithTimeout("test", []byte("hello world"), time.Second)
	if err != nil {
		panic(err)
	}

	c.Callback = func(cb *api.Call, data []byte) {
		log.Info("reply", cb.Err, string(data))
	}

	time.Sleep(time.Second * 5)

	s.Close()
}

func TestNewService(t *testing.T) {
	log.Init(nil)
	defer log.Close()

	go Srv1()
	go Srv2()

	time.Sleep(time.Second * 15)

}
