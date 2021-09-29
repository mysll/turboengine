package api

import (
	"context"
	"turboengine/common/protocol"
	"turboengine/common/utils"
)

type Module interface {
	// 模块名
	Name() string
	// 模块初始化
	Init(srv Service) error
	// 模块启动
	Start(ctx context.Context)
	// 模块关闭
	Close()
	// 设置关心的事件
	SetInterest(i int)
	// 清除关心的事件
	ClearInterest(i int)
	// 是否关心某个事件
	Interest(i int) bool
	// 模块回调
	Handler() ModuleHandler
}

type ModuleHandler interface {
	// 模块名
	Name() string
	// 初始化
	OnPrepare(Service) error
	// 启动
	OnStart(context.Context) error
	// 主循环更新
	OnUpdate(*utils.Time)
	// 停止
	OnStop() error
	// 新的客户端连接
	OnConnected(session protocol.Mailbox)
	// 客户端断线
	OnDisconnected(session protocol.Mailbox)
	// 收到客户端消息
	OnMessage(*protocol.ProtoMsg)
	// 发现新服务
	OnServiceAvailable(id uint16)
	// 服务下线
	OnServiceOffline(id uint16)
	// 准备就绪
	OnReady()
}
