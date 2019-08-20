package lock

import (
	"errors"
	"time"
	"turboengine/common/log"
	"turboengine/common/utils"
	"turboengine/core/api"
	"turboengine/core/plugin"

	consulapi "github.com/hashicorp/consul/api"
)

var (
	Name = "DisLocker"
)

type Locker interface {
	Unlock()
}

type Handler interface {
	Ok(<-chan struct{}, Locker)
	Fail(error)
}

type lockEntry struct {
	key      string
	lock     *consulapi.Lock
	h        Handler
	locked   bool
	ch       chan struct{}
	resp     <-chan struct{}
	waitTime time.Duration
	err      error
}

func (l *lockEntry) tryLock() {
	l.ch = make(chan struct{}, 1)
	resp, err := l.lock.Lock(l.ch)
	if err != nil {
		l.err = err
		return
	}
	l.locked = true
	l.resp = resp
	log.Info("locked ", l.key)
}

func (l *lockEntry) Unlock() {
	if l.locked {
		if err := l.lock.Unlock(); err != nil {
			log.Error(err)
		}
		log.Info("unlock ", l.key)
	}
}

type DisLocker struct {
	srv      api.Service
	client   *consulapi.Client
	queue    chan *lockEntry
	complete chan *lockEntry
	attachid uint64
	shut     bool
}

func (l *DisLocker) Prepare(srv api.Service, args ...interface{}) {
	l.srv = srv
	l.attachid = srv.Attach(l.process)
	l.client = args[0].(*consulapi.Client)
	l.queue = make(chan *lockEntry, 64)
	l.complete = make(chan *lockEntry, 64)
}

func (l *DisLocker) Run() {
	go l.exec()
}

func (l *DisLocker) Shut(api.Service) {
	l.shut = true
	close(l.queue)
	l.srv.Deatch(l.attachid)
L:
	for { //drain
		select {
		case e := <-l.complete:
			if e.locked {
				e.Unlock()
			}
		default:
			break L
		}
	}
	close(l.complete)
}

func (l *DisLocker) exec() {
	for entry := range l.queue {
		entry.tryLock()
		if l.shut {
			if entry.locked {
				entry.Unlock()
			}
			return
		}
		l.complete <- entry
	}
}

func (l *DisLocker) process(*utils.Time) {
	for i := 0; i < 8; i++ {
		select {
		case entry := <-l.complete:
			if entry.locked {
				entry.h.Ok(entry.resp, entry)
				break
			}
			entry.h.Fail(entry.err)
		default:
			return
		}
	}
}

func (l *DisLocker) Handle(cmd string, args ...interface{}) interface{} {
	return nil
}

func (l *DisLocker) AcquireLock(key string, h Handler, waitTime time.Duration) error {
	lock, err := l.client.LockOpts(
		&consulapi.LockOptions{
			Key: key,
			SessionOpts: &consulapi.SessionEntry{
				// Checks:   checks,
				Behavior: consulapi.SessionBehaviorDelete,
				// after release lock, other get lock wating lockDelay time.
				LockDelay: time.Millisecond,
			},
			// block wait to acquire, consul defualt 15s
			LockWaitTime: waitTime,
		},
	)

	if err != nil {
		return err
	}

	entry := &lockEntry{
		key:      key,
		lock:     lock,
		h:        h,
		waitTime: waitTime,
	}

	select {
	case l.queue <- entry:
	default:
		return errors.New("too many lock")
	}

	return nil
}

func init() {
	plugin.Register(Name, &DisLocker{})
}
