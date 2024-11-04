package redis

import (
	"time"

	"github.com/liwei1dao/lego/sys/log"
	"github.com/liwei1dao/lego/utils/mapstructure"
)

type Option func(*Options)
type Options struct {
	Debug         bool //日志是否开启
	Log           log.ILogger
	RedisAddr     []string
	RedisPassword string
	RedisDB       int
	RedisTLS      bool
	TimeOut       time.Duration
}

func SetRedisDB(v int) Option {
	return func(o *Options) {
		o.RedisDB = v
	}
}

func SetRedisAddr(v []string) Option {
	return func(o *Options) {
		o.RedisAddr = v
	}
}

func SetRedisPassword(v string) Option {
	return func(o *Options) {
		o.RedisPassword = v
	}
}
func SetRedisTLS(v bool) Option {
	return func(o *Options) {
		o.RedisTLS = v
	}
}
func SetTimeOut(v time.Duration) Option {
	return func(o *Options) {
		o.TimeOut = v
	}
}

func newOptions(config map[string]interface{}, opts ...Option) (options *Options, err error) {
	options = &Options{
		RedisAddr:     []string{"127.0.0.1:6379"},
		RedisPassword: "",
		RedisDB:       0,
		TimeOut:       time.Second * 3,
	}
	if config != nil {
		mapstructure.Decode(config, options)
	}
	for _, o := range opts {
		o(options)
	}
	if options.Log == nil {
		options.Log = log.NewTurnlog(options.Debug, log.Clone("sys.redis", 3))
	}
	return
}

func newOptionsByOption(opts ...Option) (options *Options, err error) {
	options = &Options{
		RedisAddr:     []string{"127.0.0.1:6379"},
		RedisPassword: "",
		RedisDB:       0,
		TimeOut:       time.Second * 3,
	}
	for _, o := range opts {
		o(options)
	}
	if options.Log == nil {
		options.Log = log.NewTurnlog(options.Debug, log.Clone("sys.redis", 3))
	}
	return
}

type RMutexOption func(*RMutexOptions)
type RMutexOptions struct {
	expiry time.Duration
	delay  time.Duration
}

func SetExpiry(v int) RMutexOption {
	return func(o *RMutexOptions) {
		o.expiry = time.Second * 5
	}
}
func Setdelay(v time.Duration) RMutexOption {
	return func(o *RMutexOptions) {
		o.delay = v
	}
}

func newRMutexOptions(opts ...RMutexOption) RMutexOptions {
	opt := RMutexOptions{
		expiry: time.Second * 5,
		delay:  time.Millisecond * 50,
	}
	for _, o := range opts {
		o(&opt)
	}
	return opt
}
