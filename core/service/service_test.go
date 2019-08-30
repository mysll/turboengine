package service

import (
	"fmt"
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
	s := New(&Echo{}, &Config{
		ID:      100,
		Name:    "echo",
		NatsUrl: "nats://0.0.0.0:4222",
		Depend: []Dependency{
			{
				Name:  "echo",
				Count: 1,
			},
		},
		Expose: true,
		Addr:   "0.0.0.0",
		Port:   0,
	})
	go s.Start()
	time.Sleep(time.Second)
	s.Sub("test", func(id uint16, data []byte) (*protocol.Message, error) {
		log.Info("call test from:", id, ", content:", string(data))
		m := protocol.NewMessage(len(data))
		m.Body = append(m.Body, data...)
		return m, fmt.Errorf("error")
	})
	time.Sleep(time.Second * 10)
	s.Close()
	s.Await()
}

func Srv2() {
	s := New(&Echo{}, &Config{ID: 101,
		Name:    "echo",
		NatsUrl: "nats://0.0.0.0:4222",
		Depend: []Dependency{
			{
				Name:  "echo",
				Count: 1,
			},
		}})
	go s.Start()
	time.Sleep(time.Second * 2)
	c, err := s.PubWithTimeout("test", []byte("hello world"), time.Second)
	if err != nil {
		panic(err)
	}

	c.Callback = func(call *api.Call) {
		log.Info("reply:", call.Err, string(call.Data))
	}

	time.Sleep(time.Second * 10)

	s.Close()
	s.Await()
}

func TestNewService(t *testing.T) {
	log.Init(nil)
	defer log.Close()

	go Srv1()
	go Srv2()

	time.Sleep(time.Second * 15)

}
