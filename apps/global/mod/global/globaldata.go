//@auth 	 	sll
//@create	  	2019-08-14 17:08:18
//@desc

package global

import (
	"context"
	"turboengine/common/log"
	"turboengine/common/utils"
	"turboengine/core/api"
	"turboengine/core/module"
	"turboengine/core/plugin/election"
	"turboengine/core/plugin/event"
)

type GlobalData struct {
	module.Module
	election *election.Election
	event    *event.Event
}

func (m *GlobalData) Name() string {
	return "GlobalData"
}

func (m *GlobalData) OnPrepare(s api.Service) error {
	m.Module.OnPrepare(s)
	// load module resource
	// load module resource end

	return nil
}

func (m *GlobalData) OnStart(ctx context.Context) error {
	m.Module.OnStart(ctx)
	m.election = m.Srv.Plugin(election.Name).(*election.Election)
	m.event = m.Srv.Plugin(event.Name).(*event.Event)

	// subscribe subject
	// subscribe subject end
	return nil
}

func (m *GlobalData) OnUpdate(t *utils.Time) {

}

func (m *GlobalData) OnStop() error {
	return nil
}

func (m *GlobalData) OnReady() {
	m.reqLeader()
}

func (m *GlobalData) reqLeader() {
	m.event.AddListener(election.EVENT_ELECTED, m.elected)
	m.event.AddListener(election.EVENT_FOLLOW, m.follow)
	m.election.Announce("shared/leader")
}

func (m *GlobalData) elected(event string, data interface{}) {
	log.Info("leader ", data.(string))
}

func (m *GlobalData) follow(event string, data interface{}) {
	log.Info("follow ", data.(election.LeaderInfo).Job)

}
