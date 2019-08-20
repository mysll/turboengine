package cache

import (
	"hash/fnv"
	"turboengine/common/utils"
	"turboengine/core/api"
	"turboengine/core/plugin"
)

var (
	Name = "Cache"
)

func sum64(s string) uint64 {
	h := fnv.New64a()
	h.Write([]byte(s))
	return h.Sum64()
}

type Cache struct {
	srv            api.Service
	shards         []*shard
	shardCount     int
	shardItemCount int
	shardMask      uint64
}

func (c *Cache) Prepare(srv api.Service, args ...interface{}) {
	c.srv = srv
	c.shardCount = 1024
	if len(args) > 0 {
		c.shardCount = args[0].(int)
	}
	if !utils.IsPowerOfTwo(c.shardCount) {
		panic("Shards number must be power of two")
	}

	c.shardMask = uint64(c.shardCount - 1)
	c.shards = make([]*shard, c.shardCount)
	for i := 0; i < c.shardCount; i++ {
		if len(args) > 1 {
			c.shardItemCount = args[1].(int)
		}
		c.shards[i] = newShard(c.shardItemCount)
	}
}

func (c *Cache) Run() {

}

func (c *Cache) Shut(api.Service) {
	c.Clear()
}

func (c *Cache) Handle(cmd string, args ...interface{}) interface{} {
	return nil
}

func (c *Cache) getShard(hashkey uint64) *shard {
	return c.shards[hashkey&c.shardMask]
}

func (c *Cache) Set(key string, value interface{}) {
	hashkey := sum64(key)
	shard := c.getShard(hashkey)
	shard.Set(key, hashkey, value)
}

func (c *Cache) Get(key string) interface{} {
	hashkey := sum64(key)
	shard := c.getShard(hashkey)
	return shard.Get(key, hashkey)
}

func (c *Cache) Del(key string) {
	hashkey := sum64(key)
	shard := c.getShard(hashkey)
	shard.Del(key, hashkey)
}

func (c *Cache) Clear() {
	for _, s := range c.shards {
		s.Clear()
	}
}

func init() {
	plugin.Register(Name, &Cache{})
}
