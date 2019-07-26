package event

import (
	"container/list"
	"turboengine/common/log"
	"turboengine/common/utils"
	"turboengine/core/api"
	"turboengine/core/plugin"
)

const (
	Name = "event"
)

type EventData struct {
	Name string
	Data interface{}
}

type Event struct {
	srv     api.Service
	id      uint64
	events  chan *EventData
	invokes map[string]*list.List
	serial  uint64
}

type Callback func(event string, data interface{})

type Listener struct {
	id uint64
	fn Callback
}

func (e *Event) Prepare(srv api.Service) {
	e.srv = srv
	e.id = e.srv.Attach(e.AsyncInvoke)
	e.events = make(chan *EventData, 32)
	e.invokes = make(map[string]*list.List)
	e.serial = 1
}

func (e *Event) Shut(api.Service) {
	e.srv.Deatch(e.id)
}

func (e *Event) Handle(cmd string, args ...interface{}) interface{} {
	defer func() {
		if x := recover(); x != nil {
			log.Error(x)
		}
	}()
	switch cmd {
	case "addListener":
		return e.AddListener(args[0].(string), args[1].(Callback))
	case "emit":
		e.Emit(args[0].(string), args[1])
	case "asyncemit":
		e.AsyncEmit(args[0].(string), args[1])
	}

	return nil
}

func (e *Event) AddListener(event string, fn Callback) uint64 {
	if _, ok := e.invokes[event]; !ok {
		e.invokes[event] = list.New()
	}

	serial := e.serial
	e.serial++
	e.invokes[event].PushBack(&Listener{
		id: serial,
		fn: fn,
	})
	return serial
}

func (e *Event) Emit(event string, data interface{}) {
	evt := &EventData{
		Name: event,
		Data: data,
	}
	e.Invoke(evt)
}

func (e *Event) AsyncEmit(event string, data interface{}) {
	evt := &EventData{
		Name: event,
		Data: data,
	}
	select {
	case e.events <- evt:
	default:
	}
}

func (e *Event) AsyncInvoke(t *utils.Time) {
	for {
		select {
		case event := <-e.events:
			e.Invoke(event)
		default:
			return
		}
	}
}

func (e *Event) Invoke(data *EventData) {
	if l, ok := e.invokes[data.Name]; ok {
		for invoke := l.Front(); invoke != nil; invoke = invoke.Next() {
			invoke.Value.(*Listener).fn(data.Name, data.Data)
		}
	}
}

func init() {
	plugin.Register(Name, &Event{})
}
