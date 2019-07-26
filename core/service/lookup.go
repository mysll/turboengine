package service

import (
	"context"
	"fmt"
	"sync"
	"time"
	"turboengine/common/log"

	consulapi "github.com/hashicorp/consul/api"
)

const (
	SERVICE_TAG = "turbo.service"
	SERVICE_TTL = time.Second * 2 // TTL
	DEREG_TIME  = time.Second * 10
)

const (
	EVENT_ADD = "service_add"
	EVENT_DEL = "service_del"
)

type NotifyFn func(event string, id string)

type ServiceInfo struct {
	ID   string
	Name string
	Addr string
	Port int
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
}

func NewLookupService(config *consulapi.Config) *LookupService {
	sd := &LookupService{}
	sd.config = config
	sd.localServ = make(map[string]*ServiceInfo)
	sd.service = make(map[string]*ServiceInfo)

	return sd
}

func (sd *LookupService) Start() error {
	c, err := consulapi.NewClient(sd.config)
	if err != nil {
		return nil
	}
	sd.client = c

	sd.ctx, sd.cancelFunc = context.WithCancel(context.Background())
	go sd.discover(sd.ctx, true)
	go sd.updateTTL(sd.ctx)
	return nil
}

func (sd *LookupService) Stop() {
	sd.UnregisterAll()
	sd.cancelFunc()
}

func (sd *LookupService) Register(id, name string, addr string, port int) error {
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
		return fmt.Errorf("register server error : %s", err.Error())
	}
	sd.localServ[id] = &ServiceInfo{
		ID:   id,
		Name: name,
		Addr: addr,
		Port: port,
	}
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

			addSet[entry.Service.ID] = &ServiceInfo{ // 新增加
				ID:   entry.Service.ID,
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
			if sd.notify != nil {
				sd.notify(EVENT_DEL, id)
			}
		}

		for id := range addSet {
			if sd.notify != nil {
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
			for k, _ := range sd.localServ {
				err := sd.client.Agent().UpdateTTL("service:"+k, "", consulapi.HealthPassing)
				if err != nil {
					log.Error(err)
				}
			}
		case <-ctx.Done():
			break L
		}
	}

}
