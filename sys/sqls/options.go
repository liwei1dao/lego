package sqls

import (
	"time"
)

type Option func(*Options)
type Options struct {
	SqlUrl  string
	TimeOut time.Duration
}

func SetSqlUrl(v string) Option {
	return func(o *Options) {
		o.SqlUrl = v
	}
}

func SetTimeOut(v time.Duration) Option {
	return func(o *Options) {
		o.TimeOut = v
	}
}

func newOptions(opts ...Option) Options {
	opt := Options{
		TimeOut: time.Second * 3,
	}
	for _, o := range opts {
		o(&opt)
	}
	return opt
}
