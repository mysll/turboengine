package service

import "turboengine/common/log"

func (s *service) onServiceChange(event string, id interface{}) {
	switch event {
	case EVENT_ADD:
		s.serviceValid(id.(string))
	case EVENT_DEL:
		log.Info("service del:", id.(string))
	}
}

// async call
func (s *service) notify(event string, id string) {
	s.event.AsyncEmit(event, id)
}

func (s *service) addEvent() {
	s.event.AddListener(EVENT_ADD, s.onServiceChange)
	s.event.AddListener(EVENT_DEL, s.onServiceChange)
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

	log.Info("service ready")
	s.handler.OnDependReady()
}
