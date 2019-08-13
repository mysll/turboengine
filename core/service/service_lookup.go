package service

import (
	"context"
	"fmt"
	"sort"
	"strconv"
	"sync"
	"time"
	"turboengine/common/log"
	"turboengine/common/protocol"
	"turboengine/core/api"

	consulapi "github.com/hashicorp/consul/api"
	"github.com/mysll/toolkit"

	"github.com/stathat/consistent"
)

const (
	SERVICE_TAG = "turbo.service"
	SERVICE_TTL = time.Second * 5 // TTL
	DEREG_TIME  = time.Second * 10
)

const (
	EVENT_ADD = "service_add"
	EVENT_DEL = "service_del"
)

type services []*ServiceInfo

func (s services) Len() int {
	return len(s)
}

func (s services) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s services) Less(i, j int) bool { return s[i].NID < s[j].NID }

type NotifyFn func(event string, id string)

type ServiceInfo struct {
	update string
	ID     string
	NID    uint16
	Name   string
	Addr   string
	Port   int
	Load   int // load balance
}

func (s *ServiceInfo) dirty(service *consulapi.AgentService) bool {
	if s.ID != service.ID ||
		s.Name != service.Service ||
		s.Addr != service.Address ||
		s.Port != service.Port {
		return true
	}

	return false
}

type LookupService struct {
	sync.RWMutex
	config     *consulapi.Config
	client     *consulapi.Client
	localServ  map[string]*ServiceInfo
	service    map[string]*ServiceInfo
	cancelFunc context.CancelFunc
	ctx        context.Context
	notify     NotifyFn
	lastSel    int
}

func NewLookupService(config *consulapi.Config) *LookupService {
	sd := &LookupService{}
	sd.config = config
	sd.localServ = make(map[string]*ServiceInfo)
	sd.service = make(map[string]*ServiceInfo)

	return sd
}

func (sd *LookupService) Init() error {
	c, err := consulapi.NewClient(sd.config)
	if err != nil {
		return nil
	}
	sd.client = c

	sd.ctx, sd.cancelFunc = context.WithCancel(context.Background())
	return nil
}

func (sd *LookupService) Start() {
	go sd.discover(sd.ctx, true)
	go sd.updateTTL(sd.ctx)
}

func (sd *LookupService) Stop() {
	sd.UnregisterAll()
	sd.cancelFunc()
}

func (sd *LookupService) Register(id string, name string, addr string, port int) error {
	registration := new(consulapi.AgentServiceRegistration)
	registration.ID = id
	registration.Name = name
	registration.Port = port
	registration.Tags = []string{SERVICE_TAG}
	registration.Address = addr
	registration.Check = &consulapi.AgentServiceCheck{
		TTL:                            SERVICE_TTL.String(),
		DeregisterCriticalServiceAfter: DEREG_TIME.String(),
	}

	err := sd.client.Agent().ServiceRegister(registration)

	if err != nil {
		return fmt.Errorf("register service error : %s", err.Error())
	}
	sd.localServ[id] = &ServiceInfo{
		update: "service:" + id,
		ID:     id,
		Name:   name,
		Addr:   addr,
		Port:   port,
	}

	log.Info("register succeed, id:", id)
	return nil
}

func (sd *LookupService) UnregisterAll() {
	for k, _ := range sd.localServ {
		sd.client.Agent().ServiceDeregister(k)
	}
	sd.localServ = make(map[string]*ServiceInfo)
}

func (sd *LookupService) Unregister(serviceId string) {
	if s, ok := sd.localServ[serviceId]; ok {
		sd.client.Agent().ServiceDeregister(s.ID)
		delete(sd.localServ, serviceId)
	}
}

func (sd *LookupService) StoreKV(key string, value []byte) error {
	kv := &consulapi.KVPair{
		Key:   key,
		Flags: 0,
		Value: value,
	}
	_, err := sd.client.KV().Put(kv, nil)
	return err
}

func (sd *LookupService) GetKey(key string) ([]byte, error) {
	kv, _, err := sd.client.KV().Get(key, nil)
	if err != nil {
		return nil, err
	}
	return kv.Value, err
}

func (sd *LookupService) DelKey(key string) error {
	_, err := sd.client.KV().Delete(key, nil)
	return err
}

func (sd *LookupService) Lookup(id string) *ServiceInfo {
	sd.RLock()
	if s, ok := sd.service[id]; ok {
		sd.RUnlock()
		return s
	}
	sd.RUnlock()
	return nil
}

func (sd *LookupService) AmountByName(name string) int {
	i := 0
	sd.RLock()
	for _, s := range sd.service {
		if s.Name == name {
			i++
		}
	}
	sd.RUnlock()
	return i
}

func (sd *LookupService) LookupByName(name string) []*ServiceInfo {
	var s services
	sd.RLock()
	for _, svr := range sd.service {
		if svr.Name == name {
			s = append(s, svr)
		}
	}
	sd.RUnlock()
	sort.Sort(s)
	return s
}

func (sd *LookupService) UpdateLoad(id string, load int) {
	sd.Lock()
	if s, ok := sd.service[id]; ok {
		s.Load = load
	}
	sd.Unlock()
}

func (sd *LookupService) Exist(id string) bool {
	sd.RLock()
	_, ok := sd.service[id]
	sd.RUnlock()
	return ok
}

func (sd *LookupService) SetNotify(fn NotifyFn) {
	sd.notify = fn
}

func (sd *LookupService) discover(ctx context.Context, healthyOnly bool) {
	t := time.NewTicker(time.Second)
	sd.discoverServer(healthyOnly)
L:
	for {
		select {
		case <-t.C:
			sd.discoverServer(healthyOnly)
		case <-ctx.Done():
			break L
		}
	}
	t.Stop()
}

func (sd *LookupService) discoverServer(healthyOnly bool) {
	option := &consulapi.QueryOptions{}
	services, _, err := sd.client.Catalog().Services(option)
	if err != nil {
		log.Error(err)
		return
	}

	oldSet := make(map[string]struct{})
	for k, _ := range sd.service {
		oldSet[k] = struct{}{}
	}

	addSet := make(map[string]*ServiceInfo)
	for service, _ := range services {
		serviceData, _, err := sd.client.Health().Service(service, SERVICE_TAG, healthyOnly, option)
		if err != nil {
			log.Error(err)
			return
		}
		for _, entry := range serviceData {
			if service, ok := sd.service[entry.Service.ID]; ok { // 已经存在
				if !service.dirty(entry.Service) {
					delete(oldSet, entry.Service.ID) // 没有变动，不需要处理。
					continue                         // 如果有变动，则说明服务已经重启了。则会将旧的服务移除，然后再新增加
				}
			}

			nid, err := strconv.Atoi(entry.Service.ID)
			if err != nil {
				panic(err)
			}
			addSet[entry.Service.ID] = &ServiceInfo{ // 新增加
				ID:   entry.Service.ID,
				NID:  uint16(nid),
				Name: entry.Service.Service,
				Addr: entry.Service.Address,
				Port: entry.Service.Port,
			}
		}
	}

	if len(oldSet) > 0 || len(addSet) > 0 {
		sd.Lock()
		for id := range oldSet {
			delete(sd.service, id)
		}
		for id, s := range addSet {
			sd.service[id] = s
		}
		sd.Unlock()

		// notify
		for id := range oldSet {
			if id != "0" && sd.notify != nil {
				sd.notify(EVENT_DEL, id)
			}
		}

		for id := range addSet {
			if id != "0" && sd.notify != nil {
				sd.notify(EVENT_ADD, id)
			}
		}
	}
}

func (sd *LookupService) updateTTL(ctx context.Context) {
	t := time.NewTicker(time.Second)
L:
	for {
		select {
		case <-t.C:
			for _, s := range sd.localServ {
				err := sd.client.Agent().UpdateTTL(s.update, "", consulapi.HealthPassing)
				if err != nil {
					log.Error(err)
				}
			}
		case <-ctx.Done():
			break L
		}
	}

}

func (s *service) onServiceChange(event string, id interface{}) {
	nid, _ := strconv.Atoi(id.(string))
	switch event {
	case EVENT_ADD:
		if id != s.sid {
			log.Infof("service %s avaliable", id)
			s.serviceValid(id.(string))
			s.handler.OnServiceAvailable(uint16(nid))
		}
	case EVENT_DEL:
		if id != s.sid {
			log.Infof("service %s offline", id)
		}
		s.handler.OnServiceOffline(uint16(nid))
	}
}

// async call
func (s *service) notify(event string, id string) {
	s.event.AsyncEmit(event, id)
}

func (s *service) serviceValid(id string) {
	if s.ready {
		return
	}

	for _, dep := range s.c.Depend {
		svrs := s.lookup.LookupByName(dep.Name)
		count := 0
		for _, svr := range svrs {
			if svr.ID != s.sid {
				count++
			}
		}
		if dep.Count != count {
			return
		}
	}

	s.handler.OnDependReady()
}

func (s *service) LookupById(id uint16) protocol.Mailbox {
	si := s.lookup.Lookup(strconv.Itoa(int(id)))
	if si == nil {
		return 0
	}

	return protocol.NewMailbox(id, api.MB_TYPE_SERVICE, 0)
}

func (s *service) LookupByName(name string) []protocol.Mailbox {
	var ret []protocol.Mailbox
	ss := s.lookup.LookupByName(name)
	for _, s := range ss {
		ret = append(ret, protocol.NewMailbox(s.NID, api.MB_TYPE_SERVICE, 0))
	}
	return ret
}

func (s *service) SelectService(name string, balance int, hash string) protocol.Mailbox {
	ss := s.lookup.LookupByName(name)
	count := len(ss)
	if count == 0 {
		return 0
	}
	if count == 1 {
		return protocol.NewMailbox(ss[0].NID, api.MB_TYPE_SERVICE, 0)
	}
	var id uint16
	switch balance {
	case api.LOAD_BALANCE_RAND:
		id = ss[toolkit.RandRange(0, len(ss))].NID
	case api.LOAD_BALANCE_ROUNDROBIN:
		s.lookup.lastSel++
		id = ss[s.lookup.lastSel%len(ss)].NID
	case api.LOAD_BALANCE_LEASTACTIVE:
		l := ss[0].Load
		sel := 0
		for i := 1; i < count; i++ {
			if l > ss[i].Load {
				l = ss[i].Load
				sel = i
			}
		}
		id = ss[sel].NID
	case api.LOAD_BALANCE_HASH:
		c := consistent.New()
		for _, s := range ss {
			c.Add(s.ID)
		}
		sid, err := c.Get(hash)
		if err != nil {
			panic(err)
		}

		nid, err := strconv.Atoi(sid)
		if err != nil {
			panic(err)
		}
		id = uint16(nid)
	}

	return protocol.NewMailbox(id, api.MB_TYPE_SERVICE, 0)
}
