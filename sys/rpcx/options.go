package rpcx

import (
	"errors"

	"github.com/liwei1dao/lego/core"
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
	ServiceNode    *core.ServiceNode      //服务节点
	ETCDServers    []string               //ETCD集群服务地址
	UpdateInterval int32                  //更新间隔
	RpcxStartType  RpcxStartType          //Rpcx启动类型
	AutoConnect    bool                   //自动连接 客户端启动模式下 主动连接发现的节点服务器
	SerializeType  protocol.SerializeType //序列化方式
	OutTime        int32                  //超时配置 单位秒 0 无超时限制
	Debug          bool                   //日志是否开启
	Log            log.ILogger
}

func SetServiceNode(v *core.ServiceNode) Option {
	return func(o *Options) {
		o.ServiceNode = v
	}
}

func SetETCDServers(v []string) Option {
	return func(o *Options) {
		o.ETCDServers = v
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
	if len(options.ServiceNode.Tag) == 0 || len(options.ServiceNode.Type) == 0 || len(options.ServiceNode.Id) == 0 || len(options.ETCDServers) == 0 {
		return options, errors.New("[Sys.RPCX] newOptions err: 启动参数异常")
	}

	if options.Log = log.NewTurnlog(options.Debug, log.Clone("sys.rpc", 3)); options.Log == nil {
		err = errors.New("log is nil")
		return
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
	if len(options.ServiceNode.Tag) == 0 || len(options.ServiceNode.Type) == 0 || len(options.ServiceNode.Id) == 0 || len(options.ETCDServers) == 0 {
		return options, errors.New("[Sys.RPCX] newOptions err: 启动参数异常")
	}
	if options.Log = log.NewTurnlog(options.Debug, log.Clone("sys.rpc", 3)); options.Log == nil {
		err = errors.New("log is nil")
	}
	return options, nil
}
