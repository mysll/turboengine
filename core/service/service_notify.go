package service

import (
	"strconv"
	"turboengine/common/log"
)

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

	s.ready = true
	s.handler.OnDependReady()
}
