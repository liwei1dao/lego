package consul

import (
	"crypto/tls"
	"time"

	"github.com/liwei1dao/lego/sys/log"
	"github.com/liwei1dao/lego/utils/mapstructure"
)

type Option func(*Options)
type Options struct {
	Debug             bool //日志是否开启
	Log               log.ILogger
	ConsulServers     []string
	TLS               *tls.Config
	ConnectionTimeout time.Duration
	Username          string
	Password          string
}

func newOptions(config map[string]interface{}, opts ...Option) (options Options) {
	options = Options{}
	if config != nil {
		mapstructure.Decode(config, &options)
	}
	for _, o := range opts {
		o(&options)
	}
	if options.Log == nil {
		options.Log = log.NewTurnlog(options.Debug, log.Clone("sys.discovery_etcd3", 3))
	}
	return
}

func newOptionsByOption(opts ...Option) (options Options) {
	options = Options{}
	for _, o := range opts {
		o(&options)
	}
	if options.Log == nil {
		options.Log = log.NewTurnlog(options.Debug, log.Clone("sys.discovery_etcd3", 3))
	}
	return
}
