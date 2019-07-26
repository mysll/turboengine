package rpc

import "errors"

const (
	ERR_REPLY_SUCCEED = iota // 成功
	ERR_SEND_ERR             // 发送失败
	ERR_TIME_OUT             // 超时
	ERR_ARGS_ERROR           // 参数错误
	ERR_SYSTEM_ERROR         // 系统错误
	ERR_RPC_FAILED           // rpc错误
	ERR_REPLY_FAILED         // 失败
	ERR_INNER_MAX
)

var (
	ErrShutdown = errors.New("connection is shut down")
	ErrTimeout  = errors.New("timeout")
)
