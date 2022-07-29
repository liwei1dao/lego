package redis

import (
	"time"

	"github.com/liwei1dao/lego/sys/redis/core"
	"github.com/liwei1dao/lego/utils/mapstructure"
)

type RedisType int8

const (
	Redis_Single RedisType = iota
	Redis_Cluster
)

type Option func(*Options)
type Options struct {
	RedisType              RedisType
	Redis_Single_Addr      string
	Redis_Single_Password  string
	Redis_Single_DB        int
	Redis_Cluster_Addr     []string
	Redis_Cluster_Password string
	TimeOut                time.Duration
	Codec                  core.ICodec
}

func SetRedisType(v RedisType) Option {
	return func(o *Options) {
		o.RedisType = v
	}
}

///RedisUrl = "127.0.0.1:6379"
func SetRedis_Single_Addr(v string) Option {
	return func(o *Options) {
		o.Redis_Single_Addr = v
	}
}

func SetRedis_Single_Password(v string) Option {
	return func(o *Options) {
		o.Redis_Single_Password = v
	}
}
func SetRedis_Single_DB(v int) Option {
	return func(o *Options) {
		o.Redis_Single_DB = v
	}
}

func SetRedis_Cluster_Addr(v []string) Option {
	return func(o *Options) {
		o.Redis_Cluster_Addr = v
	}
}

func SetRedis_Cluster_Password(v string) Option {
	return func(o *Options) {
		o.Redis_Cluster_Password = v
	}
}

func SetTimeOut(v time.Duration) Option {
	return func(o *Options) {
		o.TimeOut = v
	}
}

func SetCodec(v core.ICodec) Option {
	return func(o *Options) {
		o.Codec = v
	}
}

func newOptions(config map[string]interface{}, opts ...Option) (options *Options, err error) {
	options = &Options{
		Redis_Single_Addr:      "127.0.0.1:6379",
		Redis_Single_Password:  "",
		Redis_Single_DB:        1,
		Redis_Cluster_Addr:     []string{"127.0.0.1:6379"},
		Redis_Cluster_Password: "",
		TimeOut:                time.Second * 3,
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
		Redis_Single_Addr:      "127.0.0.1:6379",
		Redis_Single_Password:  "",
		Redis_Single_DB:        1,
		Redis_Cluster_Addr:     []string{"127.0.0.1:6379"},
		Redis_Cluster_Password: "",
		TimeOut:                time.Second * 3,
	}
	for _, o := range opts {
		o(options)
	}
	return
}

type RMutexOption func(*RMutexOptions)
type RMutexOptions struct {
	expiry int
	delay  time.Duration
}

func SetExpiry(v int) RMutexOption {
	return func(o *RMutexOptions) {
		o.expiry = v
	}
}
func Setdelay(v time.Duration) RMutexOption {
	return func(o *RMutexOptions) {
		o.delay = v
	}
}

func newRMutexOptions(opts ...RMutexOption) RMutexOptions {
	opt := RMutexOptions{
		expiry: 5,
		delay:  time.Millisecond * 50,
	}
	for _, o := range opts {
		o(&opt)
	}
	return opt
}
