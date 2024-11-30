package rpc

import (
	"github.com/liwei1dao/lego/core"
	"github.com/liwei1dao/lego/sys/log"
	"github.com/liwei1dao/lego/sys/rpc/rpccore"
	"github.com/liwei1dao/lego/utils/mapstructure"
)

type Option func(*Options)
type Options struct {
	ServiceNode         core.IServiceNode     //服务节点
	ProtoVersion        byte                  //协议版本
	SerializeType       rpccore.SerializeType //消息序列化方式 0:JSON 1:ProtoBuffer 2:MsgPack 3:Thrift
	CompressType        rpccore.CompressType  //消息压缩模式	0:CompressNone 1:CompressGzip
	ConnectType         rpccore.ConnectType   //RPC通信类型类型 0:Tcp 1:Kafka 2:Nats
	CommAddrs           []string              //消息传输节点 tcp:["ip:prot"] kafka:["topic_up","topic_back"]
	ConnectionTimeout   int32                 //连接超时 单位秒
	ReadTimeout         int32                 //读取超时 单位秒
	WriteTimeout        int32                 //写入超时 单位秒
	KeepAlivePeriod     int32                 //保持活跃时期 单位秒
	Heartbeat           bool                  //是否开启心跳
	HeartbeatInterval   int32                 //心跳间隔 单位秒
	MaxWaitForHeartbeat int32                 //最大等待心跳时间 单位秒
	Debug               bool                  //日志是否开启
	Log                 log.ILogger
}

func SetServiceNode(v core.IServiceNode) Option {
	return func(o *Options) {
		o.ServiceNode = v
	}
}
func SetProtoVersion(v byte) Option {
	return func(o *Options) {
		o.ProtoVersion = v
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

func SetLog(v log.ILogger) Option {
	return func(o *Options) {
		o.Log = v
	}
}

func newOptions(config map[string]interface{}, opts ...Option) (options *Options, err error) {
	options = &Options{
		ProtoVersion:  1,
		SerializeType: rpccore.MsgPack,
	}
	if config != nil {
		mapstructure.Decode(config, options)
	}
	for _, o := range opts {
		o(options)
	}
	if options.Log == nil {
		options.Log = log.NewTurnlog(options.Debug, log.Clone("sys.rpc", 3))
	}
	return
}

func newOptionsByOption(opts ...Option) (options *Options, err error) {
	options = &Options{
		ProtoVersion:  1,
		SerializeType: rpccore.MsgPack,
	}
	for _, o := range opts {
		o(options)
	}
	if options.Log == nil {
		options.Log = log.NewTurnlog(options.Debug, log.Clone("sys.rpc", 3))
	}
	return
}
