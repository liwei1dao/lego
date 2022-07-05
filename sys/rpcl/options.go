package rpcl

import (
	"time"

	"github.com/liwei1dao/lego/sys/log"
	"github.com/liwei1dao/lego/sys/rpcl/core"
	"github.com/liwei1dao/lego/utils/codec"
	"github.com/liwei1dao/lego/utils/mapstructure"
)

type Option func(*Options)
type Options struct {
	BasePath       string              //服务根节点
	ServiceType    string              //服务类型
	ServiceId      string              //服务Id
	CommType       core.NewConnectType //通信类型
	UpdateInterval time.Duration       //更新间隔时间
	Kafka_Addr     []string            //kafka 连接地址
	Kafka_Version  string              //kafka 版本
	Tcp_Addr       string              //Tcp 地址
	Decoder        codec.IDecoder      //解码器
	Encoder        codec.IEncoder      //编码器
	Debug          bool                //日志是否开启
	Log            log.ILog
}

func SetBasePath(v string) Option {
	return func(o *Options) {
		o.BasePath = v
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
func SetCommType(v bool) Option {
	return func(o *Options) {
		o.Debug = v
	}
}

func SetKafka_Addr(v []string) Option {
	return func(o *Options) {
		o.Kafka_Addr = v
	}
}
func SetKafka_Version(v string) Option {
	return func(o *Options) {
		o.Kafka_Version = v
	}
}

func SetDecoder(v codec.IDecoder) Option {
	return func(o *Options) {
		o.Decoder = v
	}
}

func SetEncoder(v codec.IEncoder) Option {
	return func(o *Options) {
		o.Encoder = v
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
		CommType: core.Tcp,
		Debug:    true,
		Log:      log.Clone(log.SetLoglayer(2)),
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
		CommType: core.Tcp,
		Debug:    true,
		Log:      log.Clone(log.SetLoglayer(2)),
	}
	for _, o := range opts {
		o(&options)
	}
	return options
}
