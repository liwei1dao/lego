package sqls

import "time"

type Option func(*Options)
type Options struct {
	ServerAdd string
	Port      int
	User      string
	Password  string
	Database  string
	TimeOut   time.Duration
}

func SetServerAdd(v string) Option {
	return func(o *Options) {
		o.ServerAdd = v
	}
}
func SetPort(v int) Option {
	return func(o *Options) {
		o.Port = v
	}
}

func SetUser(v string) Option {
	return func(o *Options) {
		o.User = v
	}
}

func SetPassword(v string) Option {
	return func(o *Options) {
		o.Password = v
	}
}
func SetDatabase(v string) Option {
	return func(o *Options) {
		o.Database = v
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
