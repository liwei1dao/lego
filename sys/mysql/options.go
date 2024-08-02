package mysql

import (
	"time"

	"github.com/liwei1dao/lego/sys/log"
	"github.com/liwei1dao/lego/utils/mapstructure"
)

type Option func(*Options)
type Options struct {
	Debug       bool //日志是否开启
	Log         log.ILogger
	MySQLDsn    string
	Database    string
	MaxPoolSize uint64
	TimeOut     time.Duration
}

func SetMySQLDsn(v string) Option {
	return func(o *Options) {
		o.MySQLDsn = v
	}
}
func SetDatabase(v string) Option {
	return func(o *Options) {
		o.Database = v
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
	if options.Log == nil {
		options.Log = log.NewTurnlog(options.Debug, log.Clone("sys.mysql", 3))
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
	if options.Log == nil {
		options.Log = log.NewTurnlog(options.Debug, log.Clone("sys.mysql", 3))
	}
	return options
}
