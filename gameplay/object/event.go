package object

import (
	"fmt"
	"reflect"
)

type OnChange func(int, interface{})

type EventCallback struct {
	ptr  uintptr
	c    OnChange
	next *EventCallback
}

func NewEventCallback(cb OnChange) *EventCallback {
	return &EventCallback{
		ptr: reflect.ValueOf(cb).Pointer(),
		c:   cb,
	}
}

type ChangeEvent struct {
	cb []*EventCallback
}

func NewChangeEvent(cap int) *ChangeEvent {
	return &ChangeEvent{
		cb: make([]*EventCallback, cap),
	}
}

func (c *ChangeEvent) add(index int, cb OnChange) error {
	if index < 0 || index > len(c.cb) {
		return fmt.Errorf("index error, get %d", index)
	}
	event := NewEventCallback(cb)
	if c.cb[index] == nil {
		c.cb[index] = event
		return nil
	}
	e := c.cb[index]
	for {
		if e.next == nil {
			e.next = event
			return nil
		}
		e = e.next
	}
}

func (c *ChangeEvent) emit(index int, val interface{}) error {
	if index < 0 || index > len(c.cb) {
		return fmt.Errorf("index error, get %d", index)
	}

	for e := c.cb[index]; e != nil; e = e.next {
		e.c(index, val)
	}
	return nil
}
