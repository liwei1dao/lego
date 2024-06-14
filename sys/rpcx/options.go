package rpcx

import (
	"errors"

	"github.com/liwei1dao/lego/sys/log"
	"github.com/liwei1dao/lego/utils/mapstructure"

	"github.com/smallnest/rpcx/protocol"
)

type RpcxStartType int8

const (
	RpcxStartByService RpcxStartType = iota //启动服务端
	RpcxStartByClient                       //启动客户端
	RpcxStartByAll                          //服务端客户端都启动
)

type Option func(*Options)
type Options struct {
	ServiceTag     string                 //集群标签
	ServiceType    string                 //服务类型
	ServiceId      string                 //服务id
	ServiceVersion string                 //服务版本
	ServiceAddr    string                 //服务地址
	ConsulServers  []string               //Consul集群服务地址
	RpcxStartType  RpcxStartType          //Rpcx启动类型
	AutoConnect    bool                   //自动连接 客户端启动模式下 主动连接发现的节点服务器
	SerializeType  protocol.SerializeType //序列化方式
	OutTime        int32                  //超时配置 单位秒 0 无超时限制
	Debug          bool                   //日志是否开启
	Log            log.ILogger
}

func SetServiceTag(v string) Option {
	return func(o *Options) {
		o.ServiceTag = v
	}
}

func SetServiceType(v string) Option {
	return func(o *Options) {
		o.ServiceType = v
	}
}

func SetServiceId(v string) Option {
	return func(o *Options) {
		o.ServiceId = v
	}
}

func SetServiceVersion(v string) Option {
	return func(o *Options) {
		o.ServiceVersion = v
	}
}

func SetServiceAddr(v string) Option {
	return func(o *Options) {
		o.ServiceAddr = v
	}
}

func SetConsulServers(v []string) Option {
	return func(o *Options) {
		o.ConsulServers = v
	}
}

// 设置启动类型
func SetRpcxStartType(v RpcxStartType) Option {
	return func(o *Options) {
		o.RpcxStartType = v
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
	options = &Options{
		AutoConnect:   true,
		SerializeType: protocol.MsgPack,
		OutTime:       5,
	}
	if config != nil {
		mapstructure.Decode(config, options)
	}
	for _, o := range opts {
		o(options)
	}
	if len(options.ServiceTag) == 0 || len(options.ServiceType) == 0 || len(options.ServiceId) == 0 || len(options.ConsulServers) == 0 {
		return options, errors.New("[Sys.RPCX] newOptions err: 启动参数异常")
	}

	if options.Log = log.NewTurnlog(options.Debug, log.Clone("sys.rpc", 3)); options.Log == nil {
		err = errors.New("log is nil")
	}

	return options, nil
}

func newOptionsByOption(opts ...Option) (options *Options, err error) {
	options = &Options{
		AutoConnect:   true,
		SerializeType: protocol.MsgPack,
		OutTime:       5,
	}
	for _, o := range opts {
		o(options)
	}
	if len(options.ServiceTag) == 0 || len(options.ServiceType) == 0 || len(options.ServiceId) == 0 || len(options.ConsulServers) == 0 {
		return options, errors.New("[Sys.RPCX] newOptions err: 启动参数异常")
	}
	if options.Log = log.NewTurnlog(options.Debug, log.Clone("sys.rpc", 3)); options.Log == nil {
		err = errors.New("log is nil")
	}
	return options, nil
}
