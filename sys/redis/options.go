package redis

import (
	"time"

	"github.com/liwei1dao/lego/utils/mapstructure"
)

///redis 存储数据格式化类型
type RedisStorageTyoe int8

const (
	JsonData RedisStorageTyoe = iota
	ProtoData
)

type Option func(*Options)
type Options struct {
	RedisUrl         string
	RedisPassword    string
	RedisDB          int
	PoolSize         int
	RedisStorageType RedisStorageTyoe
	TimeOut          time.Duration
}

///RedisUrl = "127.0.0.1:6379"
func SetRedisUrl(v string) Option {
	return func(o *Options) {
		o.RedisUrl = v
	}
}

func SetRedisPassword(v string) Option {
	return func(o *Options) {
		o.RedisPassword = v
	}
}
func SetRedisDB(v int) Option {
	return func(o *Options) {
		o.RedisDB = v
	}
}

func SetPoolSize(v int) Option {
	return func(o *Options) {
		o.PoolSize = v
	}
}

func SetRedisStorageType(v RedisStorageTyoe) Option {
	return func(o *Options) {
		o.RedisStorageType = v
	}
}

func SetTimeOut(v time.Duration) Option {
	return func(o *Options) {
		o.TimeOut = v
	}
}

func newOptions(config map[string]interface{}, opts ...Option) Options {
	options := Options{
		RedisUrl: "127.0.0.1:6379",
		RedisDB:  1,
		TimeOut:  time.Second * 3,
		PoolSize: 100,
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
		RedisUrl: "127.0.0.1:6379",
		RedisDB:  1,
		TimeOut:  time.Second * 3,
		PoolSize: 100,
	}
	for _, o := range opts {
		o(&options)
	}
	return options
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
