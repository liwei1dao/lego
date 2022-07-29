package protocol

import (
	"encoding/binary"

	lcore "github.com/liwei1dao/lego/sys/rpcl/core"
)

const (
	magicNumber byte = 0x08
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
func (this Header) MessageType() lcore.MessageType {
	return lcore.MessageType(this[2]&0x80) >> 7
}

// 设置 协议类型 请求/回应
func (this *Header) SetMessageType(mt lcore.MessageType) {
	this[2] = this[2] | (byte(mt) << 7)
}

// 是否是握手
func (h Header) IsShakeHands() bool {
	return h[2]&0x10 == 0x10
}

// 设置握手
func (this *Header) SetShakeHands(sh bool) {
	if sh {
		this[2] = this[2] | 0x10
	} else {
		this[2] = this[2] &^ 0x10
	}
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
func (this Header) CompressType() lcore.CompressType {
	return lcore.CompressType((this[2] & 0x1C) >> 2)
}

//设置压缩类型
func (this *Header) SetCompressType(ct lcore.CompressType) {
	this[2] = (this[2] &^ 0x1C) | ((byte(ct) << 2) & 0x1C)
}

//消息状态
func (this Header) MessageStatusType() lcore.MessageStatusType {
	return lcore.MessageStatusType(this[2] & 0x03)
}

// 设置消息状态 正常或者错误
func (this *Header) SetMessageStatusType(mt lcore.MessageStatusType) {
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
func (h Header) SerializeType() lcore.SerializeType {
	return lcore.SerializeType((h[3] & 0xF0) >> 4)
}

// SetSerializeType sets the serialization type.
func (h *Header) SetSerializeType(st lcore.SerializeType) {
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
