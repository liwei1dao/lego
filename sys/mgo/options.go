package mgo

import (
	"time"
)

type Option func(*Options)
type Options struct {
	MongodbUrl      string
	MongodbDatabase string
	MaxPoolSize     uint64
	TimeOut         time.Duration
}

func MongodbUrl(v string) Option {
	return func(o *Options) {
		o.MongodbUrl = v
	}
}

func MongodbDatabase(v string) Option {
	return func(o *Options) {
		o.MongodbDatabase = v
	}
}

func MaxPoolSize(v uint64) Option {
	return func(o *Options) {
		o.MaxPoolSize = v
	}
}

func TimeOut(v time.Duration) Option {
	return func(o *Options) {
		o.TimeOut = v
	}
}

func newOptions(opts ...Option) Options {
	opt := Options{
		MaxPoolSize: 1000,
		TimeOut:     time.Second * 3,
	}
	for _, o := range opts {
		o(&opt)
	}
	return opt
}
