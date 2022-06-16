package mgo

import (
	"time"

	"github.com/liwei1dao/lego/utils/mapstructure"
)

type Option func(*Options)
type Options struct {
	MongodbUrl      string
	MongodbDatabase string
	MaxPoolSize     uint64
	TimeOut         time.Duration
}

func SetMongodbUrl(v string) Option {
	return func(o *Options) {
		o.MongodbUrl = v
	}
}

func SetMongodbDatabase(v string) Option {
	return func(o *Options) {
		o.MongodbDatabase = v
	}
}

func SetMaxPoolSize(v uint64) Option {
	return func(o *Options) {
		o.MaxPoolSize = v
	}
}

func SetTimeOut(v time.Duration) Option {
	return func(o *Options) {
		o.TimeOut = v
	}
}

func newOptions(config map[string]interface{}, opts ...Option) Options {
	options := Options{
		MaxPoolSize: 100,
		TimeOut:     time.Second * 3,
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
		MaxPoolSize: 100,
		TimeOut:     time.Second * 3,
	}
	for _, o := range opts {
		o(&options)
	}
	return options
}
