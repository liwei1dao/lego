package core

import "sync"

const (
	magicNumber byte = 0x08
)

//消息压缩类型
type CompressType byte

const (
	// None does not compress.
	None CompressType = iota
	// Gzip uses gzip compression.
	Gzip
)

type MessageType byte

const (
	// Request is message type of request
	Request MessageType = iota
	// Response is message type of response
	Response
)

// MessageStatusType is status of messages.
type MessageStatusType byte

const (
	// Normal is normal requests and responses.
	Normal MessageStatusType = iota
	// Error indicates some errors occur.
	Error
)

var (
	zeroHeaderArray Header
	zeroHeader      = zeroHeaderArray[1:]
)

var msgPool = sync.Pool{
	New: func() interface{} {
		header := Header([12]byte{})
		header[0] = magicNumber

		return &Message{
			Header: &header,
		}
	},
}

func GetPooledMsg() *Message {
	return msgPool.Get().(*Message)
}

func FreeMsg(msg *Message) {
	msg.Reset()
	msgPool.Put(msg)
}

var poolUint32Data = sync.Pool{
	New: func() interface{} {
		data := make([]byte, 4)
		return &data
	},
}

type Header [12]byte

func (this Header) CheckMagicNumber() bool {
	return this[0] == magicNumber
}
func (this Header) Version() byte {
	return this[1]
}
func (this *Header) SetVersion(v byte) {
	this[1] = v
}

//设置压缩类型
func (this *Header) SetCompressType(ct CompressType) {
	this[2] = (this[2] &^ 0x1C) | ((byte(ct) << 2) & 0x1C)
}

// 设置消息类型 请求或是回应
func (this *Header) SetMessageType(mt MessageType) {
	this[2] = this[2] | (byte(mt) << 7)
}

// 设置消息状态 正常或者错误
func (this *Header) SetMessageStatusType(mt MessageStatusType) {
	this[2] = (this[2] &^ 0x03) | (byte(mt) & 0x03)
}

//是否是单向消息 是 则服务器不会发出回应消息
func (this Header) IsOneway() bool {
	return this[2]&0x20 == 0x20
}

type Message struct {
	*Header
	ServiceMethod string
	Metadata      map[string]string
	Payload       []byte
}

func (this *Message) Reset() {
	resetHeader(this.Header)
	this.Metadata = nil
	this.Payload = []byte{}
	this.ServiceMethod = ""
}
func (this Message) Clone() *Message {
	header := *this.Header
	c := GetPooledMsg()
	header.SetCompressType(None)
	c.Header = &header
	c.ServiceMethod = this.ServiceMethod
	return c
}

func resetHeader(h *Header) {
	copy(h[1:], zeroHeader)
}
