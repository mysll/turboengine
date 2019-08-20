package bytecache

import (
	"bytes"
	"testing"
)

func Test_shard_Set(t *testing.T) {
	s := newShard(128)
	s.Set("hello", sum64("hello"), []byte("world"))
	if bytes.Compare(s.Get("hello", sum64("hello")), []byte("world")) != 0 {
		t.Fail()
	}
	s.Set("hello", sum64("hello"), []byte("dlrow"))
	if bytes.Compare(s.Get("hello", sum64("hello")), []byte("dlrow")) != 0 {
		t.Fail()
	}
}
