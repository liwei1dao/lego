package sqlserver

import (
	"time"

	"github.com/liwei1dao/lego/sys/log"
	"github.com/liwei1dao/lego/utils/mapstructure"
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

func newOptions(config map[string]interface{}, opts ...Option) Options {
	options := Options{
		TimeOut: 3 * time.Second,
	}
	if config != nil {
		mapstructure.Decode(config, &options)
	}
	for _, o := range opts {
		o(&options)
	}
	if len(options.SqlUrl) == 0 {
		log.Panicf("start sqls Missing necessary configuration : SqlUrl is nul")
	}
	return options
}

func newOptionsByOption(opts ...Option) Options {
	options := Options{
		TimeOut: 3 * time.Second,
	}
	for _, o := range opts {
		o(&options)
	}
	if len(options.SqlUrl) == 0 {
		log.Panicf("start sqls Missing necessary configuration : SqlUrl is nul")
	}
	return options
}
