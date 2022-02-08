package election

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"
	"turboengine/common/log"
	"turboengine/common/protocol"
	"turboengine/common/utils"
	"turboengine/core/api"
	"turboengine/core/plugin"
	"turboengine/core/plugin/config"
	"turboengine/core/plugin/event"
	"turboengine/core/plugin/lock"
	"turboengine/core/service"
)

var (
	Name = "Election"
)

const (
	EVENT_ELECTED = "elected"
	EVENT_FOLLOW  = "follow"
)

type LeaderInfo struct {
	Id   uint16
	Ping string
	Job  string
}

type Election struct {
	srv        api.Service
	disLock    *lock.DisLocker
	leaderInfo LeaderInfo
	cfg        *config.Configuration
	event      *event.Event
	leader     bool
	job        string
	leaderKey  string
	attachId   uint64
	lastChk    time.Duration
	shut       bool
}

func (e *Election) Prepare(srv api.Service, args ...any) {
	e.srv = srv
	e.attachId = srv.Attach(e.update)
}

func (e *Election) Run() {
	e.disLock = e.srv.Plugin(lock.Name).(*lock.DisLocker)
	e.cfg = e.srv.Plugin(config.Name).(*config.Configuration)
	e.event = e.srv.Plugin(event.Name).(*event.Event)
}

func (e *Election) Shut(api.Service) {
	e.shut = true
	e.srv.Detach(e.attachId)
	if e.leader {
		e.cfg.DelKey(e.leaderKey)
	}
}

func (e *Election) Handle(cmd string, args ...any) any {
	return nil
}

func (e *Election) update(t *utils.Time) {
	if e.shut {
		return
	}
	diff := t.Time() - e.lastChk
	if diff > time.Millisecond*500 {
		e.check()
		e.lastChk = t.Time()
	}
}

func (e *Election) Announce(job string) {
	e.job = job
	e.leaderKey = fmt.Sprintf("%s/info", job)
	e.leader = false
	e.leaderInfo = LeaderInfo{}
	e.disLock.AcquireLock(job, e, time.Second*3)
}

func (e *Election) check() {
	if !e.leader {
		if e.leaderInfo.Id != 0 {
			if e.srv.LookupById(e.leaderInfo.Id).IsNil() { // 不存在
				e.Announce(e.job)
			}
		}
	}
}

func (e *Election) electLeader() {
	if e.shut {
		return
	}
	sub := fmt.Sprintf("%d:%s", e.srv.ID(), e.job)
	var info LeaderInfo
	info.Id = e.srv.ID()
	info.Ping = sub
	info.Job = e.job

	data, err := json.Marshal(info)
	if err != nil {
		panic(err)
	}
	err = e.cfg.StoreKV(e.leaderKey, data)
	if err != nil {
		log.Error(err)
	}
	e.srv.Sub(sub, e.ping)
	e.leader = true
	log.Info("elected leader for ", e.job)
	e.event.AsyncEmit(EVENT_ELECTED, e.job)
}

func (e *Election) ping(id uint16, data []byte) (msg *protocol.Message, err error) {
	if string(data) == "ping" {
		msg = protocol.NewMessage(32)
		msg.Body = append(msg.Body, []byte("pong")...)
		return
	}
	return nil, errors.New("unknown cmd")
}

func (e *Election) follow(leader LeaderInfo) {
	if e.shut {
		return
	}
	log.Info("follow ", leader)
	e.leaderInfo = leader
	e.leader = false
	e.event.AsyncEmit(EVENT_FOLLOW, leader)
}

// lock interface
func (e *Election) Ok(ch <-chan struct{}, locker lock.Locker) {
	if e.shut {
		locker.Unlock()
		return
	}
	l, err := e.cfg.GetKey(e.leaderKey)
	if err != nil {
		locker.Unlock()
		log.Error(err)
		return
	}

	if l == nil {
		e.electLeader()
		locker.Unlock()
		return
	}

	var info LeaderInfo
	err = json.Unmarshal(l, &info)
	if err != nil {
		locker.Unlock()
		log.Error(err)
		return
	}
	call, err := e.srv.PubWithTimeout(info.Ping, []byte("ping"), time.Second*2)
	if err != nil {
		log.Error(err)
		return
	}
	call.Done = make(chan *api.Call, 1)
	go func(call *api.Call, ch <-chan struct{}, leader string, locker lock.Locker) {
		select {
		case <-ch:
			// lock released
			return
		case c := <-call.Done:
			if c.Err != nil {
				if c.Err == service.ERR_TIMEOUT {
					e.electLeader()
					break
				}
				log.Error(c.Err)
				break
			}

			e.follow(info)
			if call.Msg != nil {
				call.Msg.Free()
				call.Msg = nil
			}
		}
		locker.Unlock()
	}(call, ch, string(l), locker)
}

func (e *Election) Fail(err error) {

}

func init() {
	plugin.Register(Name, &Election{})
}
