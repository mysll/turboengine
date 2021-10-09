package home

import "turboengine/common/protocol"

type StationServer struct {
}

func (s *StationServer) EnterHome(userId uint64, mb protocol.Mailbox) error {
	return nil
}

func (s *StationServer) LeaveHome(userId uint64) error {
	return nil
}

func (s *StationServer) BreakAllByGate(gateId uint64) error {
	return nil
}

func (s *StationServer) BreakAll() error {
	return nil
}
