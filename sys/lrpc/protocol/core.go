package protocol

import "github.com/liwei1dao/lego/core"

// 消息类型
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
	ServiceMethod() string
	SetServiceMethod(v string)
	From() core.IServiceNode
	SetFrom(v core.IServiceNode)
	Metadata() map[string]string
	SetMetadata(map[string]string)
	Payload() []byte
	SetPayload(b []byte)
	PrintHeader() string
}
