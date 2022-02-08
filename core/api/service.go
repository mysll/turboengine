package api

import (
	"time"
	"turboengine/common/protocol"
	"turboengine/common/utils"
)

const (
	MB_TYPE_SERVICE = iota
	MB_TYPE_CONN
)

const (
	LOAD_BALANCE_RAND         = iota // 随机
	LOAD_BALANCE_ROUND_ROBIN         // 轮询
	LOAD_BALANCE_LEAST_ACTIVE        // 最小负载
	LOAD_BALANCE_HASH                // 通过hash选取
)

const (
	INTEREST_CONNECTION_EVENT = iota + 1 // 连接事件
	INTEREST_SERVICE_EVENT               // 服务事件
)

var MAX_SID = 0x3FF

type Plugin interface {
	Prepare(srv Service, args ...any)
	Run()
	Shut(Service)
	Handle(cmd string, args ...any) any
}

type Locker interface {
	Lock()
	Unlock()
}

type Call struct {
	Session  uint64
	DeadLine time.Time
	Callback func(*Call)
	UserData any
	Err      error
	Data     []byte
	Msg      *protocol.Message
	Done     chan *Call
}

type InvokeFn func(uint16, []byte) (*protocol.Message, error)
type Update func(*utils.Time)

type Service interface {
	// ID 服务ID
	ID() uint16
	// Name 服务名
	Name() string
	// Mailbox 服务地址
	Mailbox() protocol.Mailbox
	// AddModule 增加module
	AddModule(Module)
	// Start 启动服务
	Start() error
	// Close 关闭服务
	Close()
	// Shut 关闭服务(如果选择手动关闭,则调用Shut)
	Shut()
	// Ready 服务已经就绪
	Ready()
	// Attach 将fn挂载到主循环
	Attach(fn Update) uint64
	// Detach 分离挂载
	Detach(id uint64)
	// GenGUID 生成guid
	GenGUID() uint64
	// Pub 发布消息
	Pub(subject string, data []byte) error
	// PubWithTimeout 发布消息并设置超时
	PubWithTimeout(subject string, data []byte, timeout time.Duration) (*Call, error)
	// AsyncPubWithTimeout 异步发布消息并设置超时(不依赖主循环)
	AsyncPubWithTimeout(subject string, data []byte, timeout time.Duration) (*Call, error)
	// Sub 订阅消息,invoke为收到消息时的回调函数
	Sub(subject string, invoke InvokeFn) error
	// SubNoInvoke 订阅消息
	SubNoInvoke(subject string) error
	// UnSub 取消订阅
	UnSub(subject string)
	// UsePlugin 加载插件
	UsePlugin(name string, args ...any) error
	// UnPlugin 卸载插件
	UnPlugin(name string)
	// Plugin 通过插件名获取插件
	Plugin(name string) any
	// CallPlugin 调用插件
	CallPlugin(plugin string, cmd string, args ...any) (any, error)
	// Await 阻塞直到服务结束
	Await()
	// LookupById 通过服务ID获取服务地址
	LookupById(id uint16) protocol.Mailbox
	// LookupByName 通过服务名获取服务列表
	LookupByName(name string) []protocol.Mailbox
	// SelectService 通过服务名选择一个服务,balance负载均衡策略,如果是LOAD_BALANCE_HASH,则通过hash参数进行散列处理
	SelectService(name string, balance int, hash string) protocol.Mailbox
	// SetProtoEncoder 设置协议编码器
	SetProtoEncoder(enc protocol.ProtoEncoder)
	// SetProtoDecoder 设置协议解码器
	SetProtoDecoder(dec protocol.ProtoDecoder)
	// SendToClient 发送消息到客户端
	SendToClient(dest protocol.Mailbox, msg *protocol.ProtoMsg) error
	// OpenTransport 开启socket连接(服务expose打开的情况下,才有用)
	OpenTransport()
	// CloseTransport 关闭socket连接
	CloseTransport()
}

type ServiceHandler interface {
	// OnPrepare 初始化回调
	OnPrepare(Service, map[string]string) error
	// OnStart 启动回调
	OnStart() error
	// OnShut 关闭回调
	OnShut() bool
	// OnDependReady 依赖的服务都已启动
	OnDependReady()
	// OnServiceAvailable 发现新服务
	OnServiceAvailable(id uint16)
	// OnServiceOffline 服务离线
	OnServiceOffline(id uint16)
	// OnConnected 新客户端连接
	OnConnected(session protocol.Mailbox)
	// OnDisconnected 客户端断线
	OnDisconnected(session protocol.Mailbox)
	// OnMessage 收到客户端消息
	OnMessage(*protocol.ProtoMsg)
}

// Selector 服务选择器
type Selector interface {
	Select(srv Service, service string, args string) protocol.Mailbox
}
