package clusteredcache

import (
	"errors"
	"time"
	"turboengine/common/protocol"
	"turboengine/core/api"
	"turboengine/core/plugin/clusteredcache/bytecache"
	"turboengine/core/plugin/election"
	"turboengine/core/plugin/event"
	"turboengine/core/plugin/lock"
)

var (
	Name = "ClusteredCache"
)

type ClusteredCache struct {
	srv        api.Service
	election   *election.Election
	event      *event.Event
	dislocker  *lock.DisLocker
	leader     bool
	domain     string
	version    int
	leaderinfo election.LeaderInfo
	canLeader  bool
	cache      *bytecache.Cache
	beginTx    bool
}

func (c *ClusteredCache) Prepare(srv api.Service, args ...interface{}) {
	c.srv = srv
	if srv.Plugin(election.Name) == nil {
		srv.UsePlugin(election.Name)
	}
	c.election = srv.Plugin(election.Name).(*election.Election)
	c.event = srv.Plugin(event.Name).(*event.Event)
	c.dislocker = srv.Plugin(lock.Name).(*lock.DisLocker)
	if len(args) == 0 {
		panic("args is nil")
	}
	c.domain = args[0].(string)
	if len(args) == 2 {
		c.canLeader = args[1].(bool)
	}
	c.cache = bytecache.New(1024, 128)
}

func (c *ClusteredCache) Run() {
	if c.canLeader {
		c.event.AddListener(election.EVENT_ELECTED, c.elected)
		c.event.AddListener(election.EVENT_FOLLOW, c.follow)
		c.election.Announce("domain:" + c.domain)
		return
	}
}

func (c *ClusteredCache) Shut(srv api.Service) {

}

func (c *ClusteredCache) Handle(cmd string, args ...interface{}) interface{} {
	return nil
}

func (c *ClusteredCache) elected(event string, data interface{}) {
	c.leader = true
}

func (c *ClusteredCache) follow(event string, data interface{}) {
	c.leader = false
	c.leaderinfo = data.(election.LeaderInfo)
}

func (c *ClusteredCache) Set(key string, value interface{}) error {
	if c.leader {
		c.version++
		ar := protocol.NewAutoExtendArchive(32)
		if err := ar.Put(value); err != nil {
			return err
		}
		err := c.cache.Set(key, ar.Message().Body)
		ar.Free()
		return err
	}
	return nil
}

type Tx struct {
	f func(bool)
	c *ClusteredCache
}

func (t *Tx) Ok(ch <-chan struct{}, l lock.Locker) {
	t.c.beginTx = true
	t.f(true)
	t.c.beginTx = false
}

func (t *Tx) Fail(error) {
	t.f(false)
}

func (c *ClusteredCache) Tx(f func(bool)) {
	c.dislocker.AcquireLock("lock/"+c.domain, &Tx{c: c, f: f}, time.Second*2)
}

func (c *ClusteredCache) Get(key string) []byte {
	return c.cache.Get(key)
}

func (c *ClusteredCache) GetWithValue(key string, value interface{}) error {
	data := c.cache.Get(key)
	if data == nil {
		return errors.New("not found")
	}

	ar := protocol.NewLoadArchiver(data)
	return ar.Get(value)
}

func (c *ClusteredCache) Del(key string) {
	if c.leader {
		c.version++
		if c.cache.Del(key) {

		}
	}
}
