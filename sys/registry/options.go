package registry

import "time"

type IListener interface {
	FindServiceHandlefunc(snode ServiceNode)
	UpDataServiceHandlefunc(snode ServiceNode)
	LoseServiceHandlefunc(sId string)
}

type Option func(*Options)
type Options struct {
	Address          string
	Timeout          time.Duration
	RegisterInterval time.Duration
	RegisterTTL      time.Duration
	Tag              string //集群标签
	Listener         IListener
}

func SetAddress(v string) Option {
	return func(o *Options) {
		o.Address = v
	}
}

func SetTag(v string) Option {
	return func(o *Options) {
		o.Tag = v
	}
}

func SetListener(v IListener) Option {
	return func(o *Options) {
		o.Listener = v
	}
}

func newOptions(opts ...Option) Options {
	opt := Options{
		RegisterInterval: time.Second * time.Duration(10),
		RegisterTTL:      time.Second * time.Duration(30),
		Tag:              "lego",
	}
	for _, o := range opts {
		o(&opt)
	}
	return opt
}
