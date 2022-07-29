package protocol

import (
	"encoding/binary"
)

const (
	magicNumber byte = 0x08
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

//消息头
type Header [12]byte

//协议头校验
func (this Header) CheckMagicNumber() bool {
	return this[0] == magicNumber
}

func (this Header) Version() byte {
	return this[1]
}

//设置协议版本
func (this *Header) SetVersion(v byte) {
	this[1] = v
}

// 协议类型 请求/回应
func (this Header) MessageType() MessageType {
	return MessageType(this[2]&0x80) >> 7
}

// 设置 协议类型 请求/回应
func (this *Header) SetMessageType(mt MessageType) {
	this[2] = this[2] | (byte(mt) << 7)
}

// 是否是心跳消息
func (h Header) IsHeartbeat() bool {
	return h[2]&0x40 == 0x40
}

// 设置心跳消息
func (h *Header) SetHeartbeat(hb bool) {
	if hb {
		h[2] = h[2] | 0x40
	} else {
		h[2] = h[2] &^ 0x40
	}
}

//读取压缩方式
func (this Header) CompressType() CompressType {
	return CompressType((this[2] & 0x1C) >> 2)
}

//设置压缩类型
func (this *Header) SetCompressType(ct CompressType) {
	this[2] = (this[2] &^ 0x1C) | ((byte(ct) << 2) & 0x1C)
}

//消息状态
func (this Header) MessageStatusType() MessageStatusType {
	return MessageStatusType(this[2] & 0x03)
}

// 设置消息状态 正常或者错误
func (this *Header) SetMessageStatusType(mt MessageStatusType) {
	this[2] = (this[2] &^ 0x03) | (byte(mt) & 0x03)
}

//是否是单向消息 是 true 则服务器不会发出回应消息
func (this Header) IsOneway() bool {
	return this[2]&0x20 == 0x20
}

//设置是否有回应
func (this *Header) SetOneway(oneway bool) {
	if oneway {
		this[2] = this[2] | 0x20
	} else {
		this[2] = this[2] &^ 0x20
	}
}

// SerializeType returns serialization type of payload.
func (h Header) SerializeType() SerializeType {
	return SerializeType((h[3] & 0xF0) >> 4)
}

// SetSerializeType sets the serialization type.
func (h *Header) SetSerializeType(st SerializeType) {
	h[3] = (h[3] &^ 0xF0) | (byte(st) << 4)
}

//读取回应id
func (this Header) Seq() uint64 {
	return binary.BigEndian.Uint64(this[4:])
}

//设置回应id
func (this *Header) SetSeq(seq uint64) {
	binary.BigEndian.PutUint64(this[4:], seq)
}
