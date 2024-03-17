package protocol

import (
	"encoding/binary"

	"github.com/liwei1dao/lego/sys/rpc/rpccore"
)

const (
	magicNumber byte = 0x08
)

//消息头
type Header [12]byte

//协议头校验
func (this *Header) CheckMagicNumber() bool {
	return this[0] == magicNumber
}

func (this *Header) Version() byte {
	return this[1]
}

//设置协议版本
func (this *Header) SetVersion(v byte) {
	this[1] = v
}

// 协议类型 请求/回应  10000000
func (this *Header) MessageType() rpccore.MessageType {
	return rpccore.MessageType(this[2]&0x80) >> 7
}

// 设置 协议类型 请求/回应
func (this *Header) SetMessageType(mt rpccore.MessageType) {
	this[2] = this[2] | (byte(mt) << 7)
}

// 是否是心跳消息 01000000
func (this *Header) IsHeartbeat() bool {
	return this[2]&0x40 == 0x40
}

// 设置心跳消息
func (this *Header) SetHeartbeat(hb bool) {
	if hb {
		this[2] = this[2] | 0x40
	} else {
		this[2] = this[2] &^ 0x40
	}
}

//是否是单向消息 是 true 则服务器不会发出回应消息
///00100000
func (this *Header) IsOneway() bool {
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

// 是否是握手
// 00010000
func (this *Header) IsShakeHands() bool {
	return this[2]&0x10 == 0x10
}

// 设置握手
func (this *Header) SetShakeHands(sh bool) {
	if sh {
		this[2] = this[2] | 0x10
	} else {
		this[2] = this[2] &^ 0x10
	}
}

//读取压缩方式 00001110
func (this *Header) CompressType() rpccore.CompressType {
	return rpccore.CompressType((this[2] & 0x0E) >> 1)
}

//设置压缩类型
func (this *Header) SetCompressType(ct rpccore.CompressType) {
	this[2] = (this[2] &^ 0x0E) | ((byte(ct) << 1) & 0x0E)
}

// 消息状态
// 00000001
func (this *Header) MessageStatusType() rpccore.MessageStatusType {
	return rpccore.MessageStatusType(this[2] & 0x01)
}

// 设置消息状态 正常或者错误
func (this *Header) SetMessageStatusType(mt rpccore.MessageStatusType) {
	this[2] = (this[2] &^ 0x01) | (byte(mt) & 0x01)
}

// SerializeType returns serialization type of payload.
func (this *Header) SerializeType() rpccore.SerializeType {
	return rpccore.SerializeType((this[3] & 0xF0) >> 4)
}

// SetSerializeType sets the serialization type.
func (this *Header) SetSerializeType(st rpccore.SerializeType) {
	this[3] = (this[3] &^ 0xF0) | (byte(st) << 4)
}

//读取回应id
func (this *Header) Seq() uint64 {
	return binary.BigEndian.Uint64(this[4:])
}

//设置回应id
func (this *Header) SetSeq(seq uint64) {
	binary.BigEndian.PutUint64(this[4:], seq)
}
