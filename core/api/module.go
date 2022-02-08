package api

import (
	"context"
	"turboengine/common/protocol"
	"turboengine/common/utils"
)

type Module interface {
	// Name 模块名
	Name() string
	// Init 模块初始化
	Init(srv Service) error
	// Start 模块启动
	Start(ctx context.Context)
	// Close 模块关闭
	Close()
	// SetInterest 设置关心的事件
	SetInterest(i int)
	// ClearInterest 清除关心的事件
	ClearInterest(i int)
	// Interest 是否关心某个事件
	Interest(i int) bool
	// Handler 模块回调
	Handler() ModuleHandler
}

type ModuleHandler interface {
	// Name 模块名
	Name() string
	// OnPrepare 初始化
	OnPrepare(Service) error
	// OnStart 启动
	OnStart(context.Context) error
	// OnUpdate 主循环更新
	OnUpdate(*utils.Time)
	// OnStop 停止
	OnStop() error
	// OnConnected 新的客户端连接
	OnConnected(session protocol.Mailbox)
	// OnDisconnected 客户端断线
	OnDisconnected(session protocol.Mailbox)
	// OnMessage 收到客户端消息
	OnMessage(*protocol.ProtoMsg)
	// OnServiceAvailable 发现新服务
	OnServiceAvailable(id uint16)
	// OnServiceOffline 服务下线
	OnServiceOffline(id uint16)
	// OnReady 准备就绪
	OnReady()
}
