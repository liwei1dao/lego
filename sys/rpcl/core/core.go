package core

import (
	"context"
	"errors"
	"time"

	"github.com/liwei1dao/lego/core"
	"github.com/liwei1dao/lego/sys/log"
)

var (
	ErrServerClosed          = errors.New("http: Server closed")                                   //服务关闭
	ErrMetaKVMissing         = errors.New("wrong metadata lines. some keys or values are missing") //解析Meta对象错误
	ErrUnsupportedCompressor = errors.New("unsupported compressor")                                //解压缩错误
	ErrXClientNoServer       = errors.New("can not found any server")
	ErrUnsupportedCodec      = errors.New("unsupported codec")
)

const (
	ServiceError   = "__rpcx_error__"   //服务错误信息字段
	ReqMetaDataKey = "__req_metadata"   //请求元数据字段
	ResMetaDataKey = "__res_metadata"   //返回元数据字段
	ServiceAddrKey = "__service_addr__" //服务端地址
)

type ConnectType int //通信类型
const (
	Http  ConnectType = iota //http  连接对象
	Kafka                    //Kafka 连接
	Nats                     //Nats  连接
	Tcp                      //Tcp   连接
)

type SelectMode int //选择器类型
const (
	RandomSelect       SelectMode = iota //随机选择器
	RoundRobin                           //轮询选择器
	WeightedRoundRobin                   //权重轮询选择器
	WeightedICMP                         //网络质量选择器
	RuleRobin                            //规则选择器 默认
)

//消息类型
type MessageType byte

const (
	Request  MessageType = iota //请求
	Response                    //回应
)

// 消息序列化方式
type SerializeType byte

const (
	// JSON for payload.
	JSON SerializeType = iota
	// ProtoBuffer for payload.
	ProtoBuffer
	// MsgPack for payload
	MsgPack
	// Thrift
	// Thrift for payload
	Thrift
)

//消息压缩类型
type CompressType byte

const (
	CompressNone CompressType = iota //无压缩
	CompressGzip                     //gzip压缩
)

//消息状态
type MessageStatusType byte

const (
	Normal MessageStatusType = iota //正常消息
	Error                           //错误消息
)

//系统对象
type ISys interface {
	log.Ilogf
	ServiceNode() *core.ServiceNode                                        //服务节点路径
	Heartbeat() []byte                                                     //心跳包数据 可以复用
	Handle(client IConnClient, message IMessage)                           //接收到远程消息
	ShakehandsRequest(ctx context.Context, client IConnClient) (err error) //握手请求
}

//消息对象
type IMessage interface {
	Clone() IMessage
	CheckMagicNumber() bool
	Version() byte
	SetVersion(v byte)
	MessageType() MessageType
	SetMessageType(mt MessageType)
	IsShakeHands() bool
	SetShakeHands(sh bool)
	IsHeartbeat() bool
	SetHeartbeat(hb bool)
	CompressType() CompressType
	SetCompressType(ct CompressType)
	MessageStatusType() MessageStatusType
	SetMessageStatusType(mt MessageStatusType)
	IsOneway() bool
	SetOneway(oneway bool)
	SerializeType() SerializeType
	SetSerializeType(st SerializeType)
	Seq() uint64
	SetSeq(seq uint64)
	EncodeSlicePointer() *[]byte
	ServiceMethod() string
	SetServiceMethod(v string)
	From() *core.ServiceNode
	SetFrom(v *core.ServiceNode)
	Metadata() map[string]string
	SetMetadata(map[string]string)
	Payload() []byte
	SetPayload(b []byte)
}

//路由
type ICodec interface {
	Marshal(v interface{}) ([]byte, error)
	Unmarshal(data []byte, v interface{}) error
}

//选择器
type ISelector interface {
	Select(ctx context.Context, servicePath string) []*core.ServiceNode
	UpdateServer(servers []*core.ServiceNode)
}

//连接对象池
type IConnPool interface {
	Start() error
	GetClient(node *core.ServiceNode) (client IConnClient, err error)
	Close() error
}

type IConnClient interface {
	ServiceNode() *core.ServiceNode
	Start()
	ResetHbeat()
	Write(msg []byte) (err error)
	Close() (err error)
}

//连接配置信息
type Config struct {
	ConnectType       ConnectType   //通信类型
	Endpoints         []string      //节点信息
	ConnectionTimeout time.Duration //连接超时
	ReadTimeout       time.Duration //读取超时
	WriteTimeout      time.Duration //写入超时
	KeepAlivePeriod   time.Duration //保持活跃时期
	Username          string        //用户名
	Password          string        //密码
	Vsersion          string        //版本
}
