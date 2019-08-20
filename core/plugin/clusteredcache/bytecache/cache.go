package bytecache

import (
	"hash/fnv"
	"turboengine/common/utils"
)

func sum64(s string) uint64 {
	h := fnv.New64a()
	h.Write([]byte(s))
	return h.Sum64()
}

type Cache struct {
	shards         []*shard
	shardCount     int
	shardItemCount int
	shardMask      uint64
}

func New(shardCount, shardItemCount int) *Cache {
	c := &Cache{}
	c.shardCount = shardCount
	c.shardItemCount = shardItemCount
	if !utils.IsPowerOfTwo(c.shardCount) {
		panic("Shards number must be power of two")
	}

	c.shardMask = uint64(c.shardCount - 1)
	c.shards = make([]*shard, c.shardCount)
	for i := 0; i < c.shardCount; i++ {
		c.shards[i] = newShard(c.shardItemCount)
	}
	return c
}

func (c *Cache) getShard(hashkey uint64) *shard {
	return c.shards[hashkey&c.shardMask]
}

func (c *Cache) Set(key string, value []byte) error {
	hashkey := sum64(key)
	shard := c.getShard(hashkey)
	return shard.Set(key, hashkey, value)
}

func (c *Cache) Get(key string) []byte {
	hashkey := sum64(key)
	shard := c.getShard(hashkey)
	return shard.Get(key, hashkey)
}

func (c *Cache) Del(key string) bool {
	hashkey := sum64(key)
	shard := c.getShard(hashkey)
	return shard.Del(key, hashkey)
}

func (c *Cache) Clear() {
	for _, s := range c.shards {
		s.Clear()
	}
}
