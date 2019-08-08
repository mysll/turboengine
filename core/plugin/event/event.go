package event

import (
	"container/list"
	"sync"
	"turboengine/common/log"
	"turboengine/common/utils"
	"turboengine/core/api"
	"turboengine/core/plugin"
)

const (
	Name = "Event"
)

type EventData struct {
	Name string
	Data interface{}
}

type Event struct {
	sync.Mutex
	srv     api.Service
	id      uint64
	pending *list.List
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
	e.id = e.srv.Attach(e.roundInvoke)
	e.pending = list.New()
	e.invokes = make(map[string]*list.List)
	e.serial = 1
}

func (e *Event) Shut(api.Service) {
	e.srv.Deatch(e.id)
}

func (e *Event) Run() {

}

func (e *Event) Handle(cmd string, args ...interface{}) interface{} {
	defer func() {
		if x := recover(); x != nil {
			log.Error(x)
		}
	}()
	switch cmd {
	case "AddListener":
		return e.AddListener(args[0].(string), args[1].(Callback))
	case "RemoveListener":
		e.RemoveListener(args[0].(string), args[1].(uint64))
	case "Emit":
		e.Emit(args[0].(string), args[1])
	case "AsyncEmit":
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

func (e *Event) RemoveListener(event string, id uint64) {
	if l, ok := e.invokes[event]; ok {
		for ele := l.Front(); ele != nil; ele = ele.Next() {
			if ele.Value.(*Listener).id == id {
				l.Remove(ele)
				return
			}
		}
	}
}

// sync invoke
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
	e.Lock()
	e.pending.PushBack(evt)
	e.Unlock()
}

// service update
func (e *Event) roundInvoke(t *utils.Time) {
	if e.pending.Len() == 0 {
		return
	}
	e.Lock()
	ele := e.pending.Front()
	for ele != nil {
		cur := ele
		ele = ele.Next()
		e.Invoke(cur.Value.(*EventData))
		e.pending.Remove(cur)
	}
	e.Unlock()
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
