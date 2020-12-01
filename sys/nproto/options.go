package proto




type Option func(*Options)
type Options struct {
	MsgProtoType   ProtoType
	IsUseBigEndian bool
	MessageFactory IMessageFactory
}

func SetProtoType(v ProtoType) Option {
	return func(o *Options) {
		o.MsgProtoType = v
	}
}

func SetIsUseBigEndian(v bool) Option {
	return func(o *Options) {
		o.IsUseBigEndian = v
	}
}

func SetMessageFactory(v IMessageFactory) Option {
	return func(o *Options) {
		o.MessageFactory = v
	}
}

func newOptions(opts ...Option) Options {
	opt := Options{
		IsUseBigEndian: false,
		MsgProtoType:   Proto_Json,
		MessageFactory: &DefMessageFactory{},
	}
	for _, o := range opts {
		o(&opt)
	}
	return opt
}
