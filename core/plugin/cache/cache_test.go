package cache

import (
	"fmt"
	"testing"
)

func TestCache_Set(t *testing.T) {
	cache := &Cache{}
	cache.Prepare(nil, 1024)
	cache.Set("test1", 1)
	cache.Set("test2", 2)
	cache.Set("test3", 3)
	x1 := cache.Get("test1")
	x2 := cache.Get("test2")
	x3 := cache.Get("test3")
	if x1.(int) != 1 || x2.(int) != 2 || x3.(int) != 3 {
		t.Fail()
	}
	cache.Del("test1")
	x1 = cache.Get("test1")
	if x1 != nil {
		t.Fail()
	}

	cache.Set("test2", 200)
	x2 = cache.Get("test2")
	if x2.(int) != 200 {
		t.Fail()
	}

	for i := 100; i < 200; i++ {
		cache.Set(fmt.Sprintf("test%d", i), i)
	}

	for i := 100; i < 200; i++ {
		if cache.Get(fmt.Sprintf("test%d", i)).(int) != i {
			t.Fail()
		}
	}

	for i := 100; i < 200; i++ {
		cache.Set(fmt.Sprintf("test%d", i), i+1)
	}

	for i := 100; i < 200; i++ {
		if cache.Get(fmt.Sprintf("test%d", i)).(int) != i+1 {
			t.Fail()
		}
	}

	for i := 100; i < 200; i++ {
		cache.Del(fmt.Sprintf("test%d", i))
	}

	for i := 100; i < 200; i++ {
		if cache.Get(fmt.Sprintf("test%d", i)) != nil {
			t.Fail()
		}
	}
}

func BenchmarkCacheSet(b *testing.B) {
	cache := &Cache{}
	cache.Prepare(nil, 1024)
	for i := 0; i < b.N; i++ {
		cache.Set(fmt.Sprintf("test%d", i), i)
	}
}

func BenchmarkCacheGet(b *testing.B) {
	cache := &Cache{}
	cache.Prepare(nil, 1024)
	for i := 0; i < 128; i++ {
		cache.Set(fmt.Sprintf("test%d", i), i)
	}
	mask := 127
	for i := 0; i < b.N; i++ {
		cache.Get(fmt.Sprintf("test%d", i&mask))
	}
}
