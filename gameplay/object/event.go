package object

import (
	"fmt"
	"reflect"
)

type OnChange func(self any, index int, val any)

type changeEventCallback struct {
	ptr  uintptr
	c    OnChange
	next *changeEventCallback
}

func (e *changeEventCallback) Equal(cb OnChange) bool {
	return e.ptr == reflect.ValueOf(cb).Pointer()
}

type changeEvent struct {
	cb []*changeEventCallback
}

func newChangeEvent(cap int) *changeEvent {
	return &changeEvent{
		cb: make([]*changeEventCallback, cap),
	}
}

func (c *changeEvent) add(index int, cb OnChange) error {
	if index < 0 || index > len(c.cb) {
		return fmt.Errorf("index error, get %d", index)
	}
	event := &changeEventCallback{
		ptr: reflect.ValueOf(cb).Pointer(),
		c:   cb,
	}
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

func (c *changeEvent) remove(index int, cb OnChange) error {
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

func (c *changeEvent) hasEvent(index int) bool {
	if index < 0 || index > len(c.cb) {
		return false
	}
	return c.cb[index] != nil
}

func (c *changeEvent) clear(index int) {
	c.cb[index] = nil
}

func (c *changeEvent) emit(self any, index int, val any) error {
	if index < 0 || index > len(c.cb) {
		return fmt.Errorf("index error, get %d", index)
	}

	for e := c.cb[index]; e != nil; e = e.next {
		e.c(self, index, val)
	}
	return nil
}
