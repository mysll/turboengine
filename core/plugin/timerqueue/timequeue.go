package timerqueue

import (
	"container/list"
	"sync"
	"sync/atomic"
	"time"
	"turboengine/common/utils"
	"turboengine/core/api"
	"turboengine/core/plugin"
)

const (
	TVR_BITS          = 8
	TVN_BITS          = 6
	TVR_SIZE          = 1 << TVR_BITS
	TVN_SIZE          = 1 << TVN_BITS
	TVR_MASK          = TVR_SIZE - 1
	TVN_MASK          = TVN_SIZE - 1
	MIN_TICK_INTERVAL = 1e6 // nanoseconds, 1ms
	MAXN_LEVEL        = 5
	FPS               = 50
)

const (
	Name = "timequeue"
)

type TimerInvoke func(id int64) bool
type timer struct {
	id     int64
	expire int64
	node   *list.Element
	root   *list.List
	fn     TimerInvoke
	delay  int
	count  int
}

type TimerQueue struct {
	srv           api.Service
	tickTime      int64
	ticks         int64
	nextTimerId   int64
	tvec          [MAXN_LEVEL][]*list.List
	pendingTimers *list.List
	mutex         sync.Mutex
	attachId      uint64 // attach id
	running       bool
	deleteLock    sync.Mutex
	deleting      map[int64]struct{}
}

func (tq *TimerQueue) Prepare(srv api.Service, args ...interface{}) {
	tq.srv = srv
	tq.pendingTimers = list.New()
	tq.deleting = make(map[int64]struct{})
	for i := 0; i < MAXN_LEVEL; i++ {
		if i == 0 {
			tq.tvec[i] = make([]*list.List, TVR_SIZE)
		} else {
			tq.tvec[i] = make([]*list.List, TVN_SIZE)
		}
		for j := 0; j < len(tq.tvec[i]); j++ {
			tq.tvec[i][j] = list.New()
		}
	}

	tq.attachId = srv.Attach(tq.step)
	tq.running = false
}

func (tq *TimerQueue) Shut(srv api.Service) {
	tq.srv.Deatch(tq.attachId)
	tq.running = false
}

func (tq *TimerQueue) Run() {
	tq.tickTime = now()
	tq.running = true
}

func (tq *TimerQueue) Handle(cmd string, args ...interface{}) interface{} {
	switch cmd {
	case "Schedule":
		return tq.Schedule(args[0].(int), args[1].(TimerInvoke), args[2].(int))
	case "Delete":
		tq.Delete(args[0].(int64))
	}
	return nil
}

func (tq *TimerQueue) step(t *utils.Time) {
	tq.tick(t.DeltaTime().Nanoseconds())
}

func (tq *TimerQueue) Delete(id int64) {
	tq.deleteLock.Lock()
	tq.deleting[id] = struct{}{}
	tq.deleteLock.Unlock()
}

func (tq *TimerQueue) Schedule(delay int, cb TimerInvoke, count int) int64 {
	if count < -1 {
		return -1
	}

	delayns := int64(0)
	if tq.running {
		delayns = int64(delay * 1e6)
		if delayns < MIN_TICK_INTERVAL {
			delayns = MIN_TICK_INTERVAL
		}
		delayns += atomic.LoadInt64(&(tq.tickTime))
	}

	ev := &timer{
		id:     tq.genID(),
		expire: delayns,
		node:   nil,
		root:   nil,
		fn:     cb,
		delay:  delay,
		count:  count,
	}

	tq.mutex.Lock()
	tq.pendingTimers.PushBack(ev)
	tq.mutex.Unlock()
	return ev.id
}

func (tq *TimerQueue) repeat(id int64, delay int, cb TimerInvoke, count int) {

	delayns := int64(delay * 1e6)
	if delayns < MIN_TICK_INTERVAL {
		delayns = MIN_TICK_INTERVAL
	}
	delayns += atomic.LoadInt64(&(tq.tickTime))

	ev := &timer{
		id:     id,
		expire: delayns,
		node:   nil,
		root:   nil,
		fn:     cb,
		delay:  delay,
		count:  count,
	}

	tq.mutex.Lock()
	tq.pendingTimers.PushBack(ev)
	tq.mutex.Unlock()
}

func (tq *TimerQueue) genID() int64 {
	tq.nextTimerId++
	return tq.nextTimerId
}

func now() int64 {
	return time.Now().UnixNano()
}

func (tq *TimerQueue) addTimer(t *timer) int64 {
	var vec *list.List

	ticks := (t.expire - tq.tickTime) / MIN_TICK_INTERVAL
	if ticks < 0 {
		ticks = 0
	}
	idx := tq.ticks + ticks
	level := 0

	if ticks < TVR_SIZE {
		idx = idx & TVR_MASK
		level = 0
	} else if ticks < 1<<(TVR_BITS+TVN_BITS) {
		idx = (idx >> (TVR_BITS)) & TVN_MASK
		level = 1
	} else if ticks < 1<<(TVR_BITS+2*TVN_BITS) {
		idx = (idx >> (TVR_BITS + TVN_BITS)) & TVN_MASK
		level = 2
	} else if ticks < 1<<(TVR_BITS+3*TVN_BITS) {
		idx = (idx >> (TVR_BITS + 2*TVN_BITS)) & TVN_MASK
		level = 3
	} else {
		idx = (idx >> (TVR_BITS + 3*TVN_BITS)) & TVN_MASK
		level = 4
	}
	vec = tq.tvec[level][idx]
	t.node = vec.PushBack(t)
	t.root = vec
	return t.id
}

func (tq *TimerQueue) cascade(n uint32) uint32 {
	idx := uint32(tq.ticks>>(TVR_BITS+(n-1)*TVN_BITS)) & TVN_MASK
	vec := tq.tvec[n][idx]
	tq.tvec[n][idx] = list.New()

	for e := vec.Front(); e != nil; e = e.Next() {
		t := e.Value.(*timer)
		tq.addTimer(t)
	}
	return idx
}

func (tq *TimerQueue) tick(dt int64) {
	if !tq.running {
		return
	}
	// schedule pending timers
	tq.mutex.Lock()
	pendingTimers := tq.pendingTimers
	tq.pendingTimers = list.New()
	tq.mutex.Unlock()
	for e := pendingTimers.Front(); e != nil; e = e.Next() {
		t := e.Value.(*timer)
		if t.expire == 0 {
			delayns := int64(t.delay * 1e6)
			if delayns < MIN_TICK_INTERVAL {
				delayns = MIN_TICK_INTERVAL
			}
			t.expire = atomic.LoadInt64(&(tq.tickTime)) + delayns
		}
		tq.addTimer(t)
	}

	// tick
	for ticks := dt / MIN_TICK_INTERVAL; ticks > 0; ticks-- {
		idx := tq.ticks & TVR_MASK
		if idx == 0 &&
			tq.cascade(1) == 0 &&
			tq.cascade(2) == 0 {
			tq.cascade(3)
		}

		root := tq.tvec[0][idx]
		tq.tvec[0][idx] = list.New()
		for e := root.Front(); e != nil; e = e.Next() {
			t := e.Value.(*timer)
			t.node = nil
			t.root = nil
			if t.count > 0 {
				t.count--
			}

			tq.deleteLock.Lock()
			_, del := tq.deleting[t.id]
			delete(tq.deleting, t.id)
			tq.deleteLock.Unlock()
			if !del {
				keep := t.fn(t.id)
				tq.deleteLock.Lock()
				_, del = tq.deleting[t.id]
				delete(tq.deleting, t.id)
				tq.deleteLock.Unlock()
				if !del && keep && t.count != 0 {
					tq.repeat(t.id, t.delay, t.fn, t.count)
				}

			}
		}
		tq.ticks++
		atomic.AddInt64(&(tq.tickTime), MIN_TICK_INTERVAL)
	}
}

func init() {
	plugin.Register(Name, &TimerQueue{})
}
