package object

import (
	"fmt"
	"reflect"
)

type OnChange func(self interface{}, index int, val interface{})

type ChangeEventCallback struct {
	ptr  uintptr
	c    OnChange
	next *ChangeEventCallback
}

func NewEventCallback(cb OnChange) *ChangeEventCallback {
	return &ChangeEventCallback{
		ptr: reflect.ValueOf(cb).Pointer(),
		c:   cb,
	}
}

func (e *ChangeEventCallback) Equal(cb OnChange) bool {
	return e.ptr == reflect.ValueOf(cb).Pointer()
}

type ChangeEvent struct {
	cb []*ChangeEventCallback
}

func NewChangeEvent(cap int) *ChangeEvent {
	return &ChangeEvent{
		cb: make([]*ChangeEventCallback, cap),
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

func (c *ChangeEvent) remove(index int, cb OnChange) error {
	if index < 0 || index > len(c.cb) {
		return fmt.Errorf("index error, get %d", index)
	}

	e := c.cb[index]
	if e.Equal(cb) {
		c.cb[index] = e.next
		return nil
	}
	pre := e
	e = e.next
	for ; e != nil; e = e.next {
		if e.Equal(cb) {
			pre.next = e.next
			break
		}
		pre = e
	}
	return nil
}

func (c *ChangeEvent) hasEvent(index int) bool {
	if index < 0 || index > len(c.cb) {
		return false
	}
	return c.cb[index] != nil
}

func (c *ChangeEvent) clear(index int) {
	c.cb[index] = nil
}

func (c *ChangeEvent) emit(self interface{}, index int, val interface{}) error {
	if index < 0 || index > len(c.cb) {
		return fmt.Errorf("index error, get %d", index)
	}

	for e := c.cb[index]; e != nil; e = e.next {
		e.c(self, index, val)
	}
	return nil
}
