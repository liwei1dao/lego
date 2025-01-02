package rpc

import (
	"context"
	"errors"
	"io"

	"github.com/liwei1dao/lego/core"
)

var (
	ErrServerClosed          = errors.New("http: Server closed")                                   //服务关闭
	ErrMetaKVMissing         = errors.New("wrong metadata lines. some keys or values are missing") //解析Meta对象错误
	ErrUnsupportedCompressor = errors.New("unsupported compressor")                                //解压缩错误
	ErrXClientNoServer       = errors.New("can not found any server")
	ErrUnsupportedCodec      = errors.New("unsupported codec")
)

// 消息类型
type MessageType byte

const (
	Request  MessageType = iota //请求
	Response                    //回应
)

// 消息压缩类型
type CompressType byte

const (
	CompressNone CompressType = iota //无压缩
	CompressGzip                     //gzip压缩
)

// 消息状态
type MessageStatusType byte

const (
	Normal MessageStatusType = iota //正常消息
	Error                           //错误消息
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

// 消息对象
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
	GetService() string
	SetService(v string)
	From() core.IServiceNode
	SetFrom(v core.IServiceNode)
	Metadata() map[string]string
	SetMetadata(map[string]string)
	Payload() []byte
	SetPayload(b []byte)
	PrintHeader() string
}

type KV struct {
	Key   string
	Value string
}

type ClientState int32

const (
	ClientClose      ClientState = iota //关闭状态
	ClientShakeHands                    //握手状态
	ClientRuning                        //运行中
	ClientCloseing                      //关闭中
)

type IConnClient interface {
	ServiceNode() core.IServiceNode
	SetServiceNode(node core.IServiceNode)
	State() ClientState
	Start()
	ResetHbeat()
	Write(msg []byte) (err error)
	Close() (err error)
}

// 连接对象池
type IConnPool interface {
	Start(host IBodyHost) error
	GetClient(node core.IServiceNode) (client IConnClient, err error)
	AddClient(client IConnClient, node core.IServiceNode) (err error)
	Close() error
}

// 主机对象
type IBodyHost interface {
	ReadProtocol(client IConnClient, r io.Reader) (err error)
	Heartbeat() (buff []byte)
	ShakehandsRequest(ctx context.Context, client IConnClient) (err error)
}
