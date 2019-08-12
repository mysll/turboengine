package proto

import "turboengine/common/protocol"

const (
	LOGIN_RESULT = 2000 + iota
)

type LoginResult struct {
	Result bool
}

func init() {
	protocol.AddProto(LOGIN_RESULT, func() interface{} { return &LoginResult{} })
}
