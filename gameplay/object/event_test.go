package object

import (
	"testing"
)

func TestNewChangeEvent(t *testing.T) {
	event := newChangeEvent(1)
	fc := func(object any, index int, val any) {
		t.Log(val)
	}
	fc1 := func(object any, index int, val any) {
		t.Log(val, ",2")
	}
	event.add(0, fc)
	event.add(0, fc1)
	event.emit(nil, 0, 123)
	event.remove(0, fc)
	event.emit(nil, 0, 456)
}
