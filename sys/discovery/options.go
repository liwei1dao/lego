package discovery

import (
	"errors"
	"time"

	"github.com/liwei1dao/lego/core"
	dcore "github.com/liwei1dao/lego/sys/discovery/core"
	"github.com/liwei1dao/lego/sys/log"
	"github.com/liwei1dao/lego/utils/mapstructure"
)

type ICodec interface {
	Marshal(v interface{}) ([]byte, error)
	Unmarshal(data []byte, v interface{}) error
}

type StoreType int8

const (
	StoreConsul StoreType = iota
	StoreZookeeper
	StoreRedis
)

type Option func(*Options)
type Options struct {
	BasePath       string            //根节点
	ServiceNode    *core.ServiceNode //服务节点
	UpdateInterval time.Duration     //更新间隔
	StoreType      StoreType         //第三方服务类型 支持 Consul  Zookeeper Redis
	Endpoints      []string          //服务节点集合
	Config         *dcore.Config     //连接配置
	Codec          ICodec            //编解码工具
	Debug          bool              //日志是否开启
	Log            log.ILog
}

func SetBasePath(v string) Option {
	return func(o *Options) {
		o.BasePath = v
	}
}
func SetServiceNode(v *core.ServiceNode) Option {
	return func(o *Options) {
		o.ServiceNode = v
	}
}
func SetUpdateInterval(v time.Duration) Option {
	return func(o *Options) {
		o.UpdateInterval = v
	}
}
func SetStoreType(v StoreType) Option {
	return func(o *Options) {
		o.StoreType = v
	}
}
func SetEndpoints(v []string) Option {
	return func(o *Options) {
		o.Endpoints = v
	}
}
func SetConfig(v *dcore.Config) Option {
	return func(o *Options) {
		o.Config = v
	}
}
func SetCodec(v ICodec) Option {
	return func(o *Options) {
		o.Codec = v
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
		BasePath: "discovery",
	}
	if config != nil {
		mapstructure.Decode(config, &options)
	}
	for _, o := range opts {
		o(options)
	}
	if options.Debug && options.Log == nil {
		if options.Log = log.Clone(log.SetLoglayer(2)); options.Log == nil {
			err = errors.New("log is nil")
		}
	}
	return
}

func newOptionsByOption(opts ...Option) (options *Options, err error) {
	options = &Options{
		BasePath: "discovery",
	}
	for _, o := range opts {
		o(options)
	}
	if options.Debug && options.Log == nil {
		if options.Log = log.Clone(log.SetLoglayer(2)); options.Log == nil {
			err = errors.New("log is nil")
		}
	}
	return
}
