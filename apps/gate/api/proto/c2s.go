package proto

import "turboengine/common/protocol"

const (
	LOGIN = 1000 + iota
)

type Login struct {
	User string
	Pass string
}

func init() {
	protocol.AddProto(LOGIN, func() interface{} { return &Login{} })
}
