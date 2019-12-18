package proto

import "bufio"

type Option func(*Options)
type Options struct {
	MsgProtoType         ProtoType
	IsUseBigEndian       bool
	MessageDecodeBybufio func(r *bufio.Reader) (IMessage, error)
	MessageDecodeBybytes func(buffer []byte) (msg IMessage, err error)
	MsgMarshal           func(comId uint16, msgId uint16, msg interface{}) IMessage
}

func SetProtoType(v ProtoType) Option {
	return func(o *Options) {
		o.MsgProtoType = v
	}
}

func SetMessageDecodeBybufio(v func(r *bufio.Reader) (IMessage, error)) Option {
	return func(o *Options) {
		o.MessageDecodeBybufio = v
	}
}

func SetMessageDecodeBybytes(v func(buffer []byte) (msg IMessage, err error)) Option {
	return func(o *Options) {
		o.MessageDecodeBybytes = v
	}
}

func SetMsgMarshal(v func(comId uint16, msgId uint16, msg interface{}) IMessage) Option {
	return func(o *Options) {
		o.MsgMarshal = v
	}
}

func SetIsUseBigEndian(v bool) Option {
	return func(o *Options) {
		o.IsUseBigEndian = v
	}
}

func newOptions(opts ...Option) Options {
	opt := Options{
		IsUseBigEndian:       false,
		MsgProtoType:         Proto_Json,
		MessageDecodeBybufio: DefMessageDecodeBybufio,
		MessageDecodeBybytes: DefMessageDecodeBybytes,
		MsgMarshal:           DefMessageMarshal,
	}
	for _, o := range opts {
		o(&opt)
	}
	return opt
}
