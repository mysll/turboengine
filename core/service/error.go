package service

import "errors"

var (
	ERR_MSG_TOO_MANY = errors.New("message too many")
	ERR_CLOSED       = errors.New("already closed")
	ERR_TIMEOUT      = errors.New("time out")
)
