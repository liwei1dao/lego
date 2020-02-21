package proto

import (
	"bufio"
	"reflect"

	"github.com/liwei1dao/lego/core"
)

var (
	option               Options
	MsgProtoType         ProtoType
	IsUseBigEndian       bool
	MsgUnMarshal         func(t reflect.Type, d []byte) (interface{}, error)
	StructUnmarshal      func(d []byte, v interface{}) error
	StructMarshal        func(v interface{}) ([]byte, error)
	MsgToString          func(v interface{}) string
	MessageFactory 		 IMessageFactory
	// MsgMarshal           func(comId uint16, msgId uint16, msg interface{}) IMessage
	// MessageDecodeBybufio func(r *bufio.Reader) (IMessage, error)
	// MessageDecodeBybytes func(buffer []byte) (msg IMessage, err error)
)

type (
	ProtoType uint8
	IMessage  interface {
		GetComId() uint16
		GetMsgId() uint16
		GetMsg() []byte
		Serializable() (bytes []byte, err error)
		ToString() string
	}
	IMessageFactory interface {
		MessageDecodeBybufio func(r *bufio.Reader) (IMessage, error)
		MessageDecodeBybytes func(buffer []byte) (msg IMessage, err error)
		MsgMarshal           func(comId uint16, msgId uint16, msg interface{}) IMessage
	}
	IMsgMarshalString interface { //消息输出结构
		ToString() (string, error)
	}
)

var (
	Proto_Json ProtoType = 1
	Proto_Buff ProtoType = 2
)

func OnInit(s core.IService, opt ...Option) (err error) {
	option = newOptions(opt...)
	MsgProtoType = option.MsgProtoType
	if MsgProtoType == Proto_Json {
		MsgUnMarshal = jsonStructUnmarshal
		StructUnmarshal = JsonStructUnmarshal
		StructMarshal = JsonStructMarshal
	} else {
		MsgUnMarshal = protoStructUnmarshal

		StructUnmarshal = ProtoStructUnmarshal
		StructMarshal = ProtoStructMarshal
	}
	IsUseBigEndian = option.IsUseBigEndian
	MsgToString = msgToString
	MessageFactory = option.MessageFactory
	return
}
