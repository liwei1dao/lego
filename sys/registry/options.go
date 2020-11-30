package registry

import (
	"time"

	"github.com/liwei1dao/lego/utils/mapstructure"
)

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

func newOptions(config map[string]interface{}, opts ...Option) Options {
	options := Options{
		RegisterInterval: time.Second * 10,
		RegisterTTL:      time.Second * 30,
		Tag:              "lego",
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
		RegisterInterval: time.Second * 10,
		RegisterTTL:      time.Second * 30,
		Tag:              "lego",
	}
	for _, o := range opts {
		o(&options)
	}
	return options
}
