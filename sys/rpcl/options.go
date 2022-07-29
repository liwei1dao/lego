package rpcl

import (
	"github.com/liwei1dao/lego/core"
	"github.com/liwei1dao/lego/sys/log"
	lcore "github.com/liwei1dao/lego/sys/rpcl/core"
	"github.com/liwei1dao/lego/utils/mapstructure"
)

type Option func(*Options)
type Options struct {
	ServiceNode   *core.ServiceNode   //服务节点
	SerializeType lcore.SerializeType //消息序列化方式
	CompressType  lcore.CompressType  //消息压缩模式
	Auth          string              //认证
	Config        *lcore.Config       //连接配置
	Debug         bool                //日志是否开启
	Log           log.ILog
}

func SetServiceNode(v *core.ServiceNode) Option {
	return func(o *Options) {
		o.ServiceNode = v
	}
}
func SetCommType(v bool) Option {
	return func(o *Options) {
		o.Debug = v
	}
}
func SetEndpoints(v bool) Option {
	return func(o *Options) {
		o.Debug = v
	}
}

func SetDebug(v bool) Option {
	return func(o *Options) {
		o.Debug = v
	}
}

func SetLog(v log.ILog) Option {
	return func(o *Options) {
		o.Log = v
	}
}

func newOptions(config map[string]interface{}, opts ...Option) Options {
	options := Options{
		Debug: true,
		Log:   log.Clone(log.SetLoglayer(2)),
	}
	if config != nil {
		mapstructure.Decode(config, &options)
	}
	for _, o := range opts {
		o(&options)
	}
	return options
}

func newOptionsByOption(opts ...Option) Options {
	options := Options{
		Debug: true,
		Log:   log.Clone(log.SetLoglayer(2)),
	}
	for _, o := range opts {
		o(&options)
	}
	return options
}
