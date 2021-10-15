package home

import "turboengine/common/protocol"

type BrokerServer struct {
}

func (s *BrokerServer) EnterHome(userId uint64, mb protocol.Mailbox) error {
	return nil
}

func (s *BrokerServer) LeaveHome(userId uint64) error {
	return nil
}

func (s *BrokerServer) BreakAllByGate(gateId uint64) error {
	return nil
}

func (s *BrokerServer) BreakAll() error {
	return nil
}
