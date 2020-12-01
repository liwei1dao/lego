package proto

import (
	"github.com/liwei1dao/lego/utils/mapstructure"
)

type ProtoType uint8

const (
	Proto_Json ProtoType = 1
	Proto_Buff ProtoType = 2
)

type Option func(*Options)
type Options struct {
	MsgProtoType   ProtoType
	IsUseBigEndian bool
}

func SetMsgProtoType(v ProtoType) Option {
	return func(o *Options) {
		o.MsgProtoType = v
	}
}

func SetIsUseBigEndian(v bool) Option {
	return func(o *Options) {
		o.IsUseBigEndian = v
	}
}

func newOptions(config map[string]interface{}, opts ...Option) Options {
	options := Options{
		MsgProtoType:   Proto_Buff,
		IsUseBigEndian: false,
	}
	if config != nil {
		mapstructure.Decode(config, &options)
	}
	for _, o := range opts {
		o(&options)
	}
	return options
}

func newOptionsByOption(opts ...Option) Options {
	options := Options{
		MsgProtoType:   Proto_Buff,
		IsUseBigEndian: false,
	}
	for _, o := range opts {
		o(&options)
	}
	return options
}
