package sql

import (
	"github.com/liwei1dao/lego/sys/log"
	"github.com/liwei1dao/lego/utils/mapstructure"
)

type SqlType string

const (
	SqlServer SqlType = "sqlserver"
	MySql     SqlType = "mysql"
	Oracle    SqlType = "godror"
	DM        SqlType = "dm"
)

type Option func(*Options)
type Options struct {
	SqlType SqlType
	SqlUrl  string
	TimeOut int
}

func SetSqlType(v SqlType) Option {
	return func(o *Options) {
		o.SqlType = v
	}
}

func SetSqlUrl(v string) Option {
	return func(o *Options) {
		o.SqlUrl = v
	}
}

func SetTimeOut(v int) Option {
	return func(o *Options) {
		o.TimeOut = v
	}
}

func newOptions(config map[string]interface{}, opts ...Option) Options {
	options := Options{
		SqlType: MySql,
		TimeOut: 3,
	}
	if config != nil {
		mapstructure.Decode(config, &options)
	}
	for _, o := range opts {
		o(&options)
	}
	if len(options.SqlUrl) == 0 {
		log.Errorf("start sqls Missing necessary configuration : SqlUrl is nul")
	}
	return options
}

func newOptionsByOption(opts ...Option) Options {
	options := Options{
		TimeOut: 3,
	}
	for _, o := range opts {
		o(&options)
	}
	if len(options.SqlUrl) == 0 {
		log.Errorf("start sqls Missing necessary configuration : SqlUrl is nul")
	}
	return options
}
