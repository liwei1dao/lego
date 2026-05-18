package postgresql

import (
	"time"

	"github.com/liwei1dao/lego/sys/log"
	"github.com/liwei1dao/lego/utils/mapstructure"
)

// Option 函数式选项,用于在初始化阶段修改 Options
type Option func(*Options)

// Options postgresql 系统的可配置项
// 字段标签上的中文说明同时作为 mapstructure 解码后的人类可读注释
type Options struct {
	Debug             bool          //是否开启调试日志(影响 log.NewTurnlog 行为)
	Log               log.ILogger   //日志组件,未配置时由系统基于 sys.postgresql 自动派生
	PostgresqlUrl     string        //连接字符串,格式 postgres://user:password@host:port/dbname?sslmode=disable
	MaxConns          int32         //连接池最大连接数,达到上限后获取连接将阻塞或超时
	MinConns          int32         //连接池保留的最小空闲连接数,有助于减少冷启动开销
	MaxConnLifetime   time.Duration //单个连接的最长生命周期,到期后池会关闭并重建,规避中间件踢链
	MaxConnIdleTime   time.Duration //单个连接最大空闲时长,长期空闲的连接会被回收
	HealthCheckPeriod time.Duration //后台健康检查周期,池会定期 Ping 空闲连接剔除异常连接
	ConnectTimeout    time.Duration //建立单条连接(以及启动阶段首次 Ping)的超时
}

// SetDebug 设置是否开启调试日志
func SetDebug(v bool) Option {
	return func(o *Options) {
		o.Debug = v
	}
}

// SetLog 注入外部日志组件,未调用时使用默认 sys.postgresql 日志
func SetLog(v log.ILogger) Option {
	return func(o *Options) {
		o.Log = v
	}
}

// SetPostgresqlUrl 设置 DSN,必填项
func SetPostgresqlUrl(v string) Option {
	return func(o *Options) {
		o.PostgresqlUrl = v
	}
}

// SetMaxConns 设置连接池最大连接数
func SetMaxConns(v int32) Option {
	return func(o *Options) {
		o.MaxConns = v
	}
}

// SetMinConns 设置连接池最小空闲连接数
func SetMinConns(v int32) Option {
	return func(o *Options) {
		o.MinConns = v
	}
}

// SetMaxConnLifetime 设置单连接最长存活时间
func SetMaxConnLifetime(v time.Duration) Option {
	return func(o *Options) {
		o.MaxConnLifetime = v
	}
}

// SetMaxConnIdleTime 设置单连接最长空闲时间
func SetMaxConnIdleTime(v time.Duration) Option {
	return func(o *Options) {
		o.MaxConnIdleTime = v
	}
}

// SetHealthCheckPeriod 设置后台健康检查周期
func SetHealthCheckPeriod(v time.Duration) Option {
	return func(o *Options) {
		o.HealthCheckPeriod = v
	}
}

// SetConnectTimeout 设置建连超时
func SetConnectTimeout(v time.Duration) Option {
	return func(o *Options) {
		o.ConnectTimeout = v
	}
}

// newOptions 由 OnInit 使用:先以默认值为底,再用 config(map)解码覆盖,
// 最后应用函数式选项,确保选项优先级高于配置文件
func newOptions(config map[string]interface{}, opts ...Option) Options {
	options := defaultOptions()
	if config != nil {
		mapstructure.Decode(config, &options)
	}
	for _, o := range opts {
		o(&options)
	}
	if options.Log == nil {
		options.Log = log.NewTurnlog(options.Debug, log.Clone("sys.postgresql", 3))
	}
	return options
}

// newOptionsByOption 由 NewSys 使用:仅通过函数式选项配置,适合代码内直接构造实例
func newOptionsByOption(opts ...Option) Options {
	options := defaultOptions()
	for _, o := range opts {
		o(&options)
	}
	if options.Log == nil {
		options.Log = log.NewTurnlog(options.Debug, log.Clone("sys.postgresql", 3))
	}
	return options
}

// defaultOptions 返回内置默认值,作为所有初始化路径的基线
// 默认值倾向于服务端常见场景:中等并发、不轻易踢链、5s 建连超时
func defaultOptions() Options {
	return Options{
		MaxConns:          16,
		MinConns:          2,
		MaxConnLifetime:   time.Hour,
		MaxConnIdleTime:   time.Minute * 30,
		HealthCheckPeriod: time.Minute,
		ConnectTimeout:    time.Second * 5,
	}
}
