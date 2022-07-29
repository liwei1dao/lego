package rpcl

import (
	"github.com/liwei1dao/lego/core"
	"github.com/liwei1dao/lego/sys/discovery"
	"github.com/liwei1dao/lego/sys/log"
	lcore "github.com/liwei1dao/lego/sys/rpcl/core"
	"github.com/liwei1dao/lego/utils/mapstructure"
)

type Option func(*Options)
type Options struct {
	ServiceNode        *core.ServiceNode   //服务节点
	SerializeType      lcore.SerializeType //消息序列化方式 0:JSON 1:ProtoBuffer 2:MsgPack 3:Thrift
	CompressType       lcore.CompressType  //消息压缩模式	0:CompressNone 1:CompressGzip
	ConnectType        lcore.ConnectType   //RPC通信类型类型 0:Tcp 1:Kafka 2:Nats
	MessageEndpoints   []string            //消息传输节点 tcp:"ip:prot"
	ConnectionTimeout  int32               //连接超时 单位秒
	ReadTimeout        int32               //读取超时 单位秒
	WriteTimeout       int32               //写入超时 单位秒
	KeepAlivePeriod    int32               //保持活跃时期 单位秒
	DiscoveryStoreType discovery.StoreType //发现服务类型  支持 Consul  Zookeeper Redis
	DiscoveryEndpoints []string            //服务发现节点
	DiscoveryInterval  int32               //发现间隔 单位秒
	Debug              bool                //日志是否开启
	Log                log.ILog
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

func newOptions(config map[string]interface{}, opts ...Option) (options *Options, err error) {
	options = &Options{
		Debug: true,
		Log:   log.Clone(log.SetLoglayer(2)),
	}
	if config != nil {
		mapstructure.Decode(config, options)
	}
	for _, o := range opts {
		o(options)
	}
	return
}

func newOptionsByOption(opts ...Option) (options *Options, err error) {
	options = &Options{
		Debug: true,
		Log:   log.Clone(log.SetLoglayer(2)),
	}
	for _, o := range opts {
		o(options)
	}
	return
}
