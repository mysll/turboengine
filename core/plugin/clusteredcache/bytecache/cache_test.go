package bytecache

import (
	"bytes"
	"testing"
)

func TestCache_Set(t *testing.T) {
	cache := New(1024, 128)
	cache.Set("hello", []byte("world"))
	if bytes.Compare(cache.Get("hello"), []byte("world")) != 0 {
		t.Fail()
	}
}
