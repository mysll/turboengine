package cache

import "testing"

func Test_shard(t *testing.T) {
	shard := newShard(3)
	shard.Set("test1", 1, 1)
	shard.Set("test2", 2, 2)
	shard.Set("test3", 3, 3)
	x1 := shard.Get("test1", 1).(int)
	x2 := shard.Get("test2", 2).(int)
	x3 := shard.Get("test3", 3).(int)
	if x1 != 1 || x2 != 2 || x3 != 3 {
		t.Fatal("error")
	}

	shard.Set("test4", 1, 4)
	x1 = shard.Get("test1", 1).(int)
	x2 = shard.Get("test2", 2).(int)
	x3 = shard.Get("test3", 3).(int)
	x4 := shard.Get("test4", 1).(int)
	if x1 != 1 || x2 != 2 || x3 != 3 || x4 != 4 {
		t.Fatal("error")
	}
	shard.Del("test1", 1)
	x4 = shard.Get("test4", 1).(int)
	if x4 != 4 {
		t.Fatal("error")
	}
	if shard.Get("test1", 1) != nil {
		t.Fatal("error")
	}

	shard.Del("test2", 2)

	shard.Set("test2", 2, 2)
	shard.Set("test5", 5, 5)
	shard.Set("test6", 6, 6)
	shard.Set("test7", 7, 7)
	shard.output()
	shard.Del("test4", 1)
	shard.Del("test2", 2)
	shard.Del("test6", 6)
	//shard.Del("test7", 4)
	shard.output()
	shard.Set("test9", 9, 9)
	shard.output()

	x3 = shard.Get("test3", 3).(int)
	x5 := shard.Get("test5", 5).(int)
	x7 := shard.Get("test7", 7).(int)
	x9 := shard.Get("test9", 9).(int)
	if x3 != 3 || x5 != 5 || x7 != 7 || x9 != 9 {
		t.Fatal("error")
	}
}
