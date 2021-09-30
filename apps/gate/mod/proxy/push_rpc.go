package proxy

import (
	"turboengine/common/protocol"
	"turboengine/core/api"
)

type PushServer struct {
	srv api.Service
}

func (p *PushServer) PushToUser(dest protocol.Mailbox, message []byte) error {
	return nil
}
