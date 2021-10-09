package proto

import "turboengine/common/protocol"

type Station struct {
	Ver string `version:"1.0.0"`
	XXX interface{}
	// custom method begin
	EnterHome      func(userId uint64, mb protocol.Mailbox) error
	LeaveHome      func(userId uint64) error
	BreakAllByGate func(gateId uint64) error
	BreakAll       func() error
	// custom method end
}

func init() {
	reg["Station"] = new(Station)
}
