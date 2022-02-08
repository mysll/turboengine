package proto

import "turboengine/common/protocol"

type Push struct {
	Ver string `version:"1.0.0"`
	XXX any
	// custom method begin
	PushToUser func(dest protocol.Mailbox, message []byte) error
	// custom method end
}

func init() {
	reg["Push"] = new(Push)
}
