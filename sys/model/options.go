package model

import (
	"errors"

	"github.com/liwei1dao/lego/sys/log"
	"github.com/liwei1dao/lego/utils/mapstructure"
)

type Option func(*Options)
type Options struct {
	Tag             string   //标签
	RedisIsCluster  bool     //是否是集群
	RedisAddr       []string //redis 的集群地址
	RedisPassword   string   //redis的密码
	RedisDB         int      //数据库位置
	MongodbUrl      string   //数据库连接地址
	MongodbDatabase string   //数据库名
	Debug           bool     //日志是否开启
	Log             log.ILogger
}

func SetTag(v string) Option {
	return func(o *Options) {
		o.Tag = v
	}
}
func SetDebug(v bool) Option {
	return func(o *Options) {
		o.Debug = v
	}
}

func SetLog(v log.ILogger) Option {
	return func(o *Options) {
		o.Log = v
	}
}
func newOptions(config map[string]interface{}, opts ...Option) (options *Options, err error) {
	options = &Options{}
	if config != nil {
		if err = mapstructure.Decode(config, options); err != nil {
			return
		}
	}
	for _, o := range opts {
		o(options)
	}
	if options.Log = log.NewTurnlog(options.Debug, log.Clone("sys.db", 2)); options.Log == nil {
		err = errors.New("log is nil")
	}
	return
}

func newOptionsByOption(opts ...Option) (options *Options, err error) {
	options = &Options{}
	for _, o := range opts {
		o(options)
	}
	if options.Log = log.NewTurnlog(options.Debug, log.Clone("sys.db", 2)); options.Log == nil {
		err = errors.New("log is nil")
	}
	return
}

type DBOption func(*DBOptions)
type DBOptions struct {
	IsMgoLog bool //是否写mgolog
}

//设置是否写mgor日志
func SetDBMgoLog(v bool) DBOption {
	return func(o *DBOptions) {
		o.IsMgoLog = v
	}
}

//更具 Option 序列化 系统参数对象
func newDBOption(opts ...DBOption) DBOptions {
	options := DBOptions{
		IsMgoLog: true,
	}
	for _, o := range opts {
		o(&options)
	}
	return options
}
