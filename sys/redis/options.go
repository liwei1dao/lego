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
		RedisUrl: "redis://root:li13451234@127.0.0.1:6379/1",
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
		RedisUrl: "redis://root:li13451234@127.0.0.1:6379/1",
		TimeOut:  time.Second * 3,
		PoolSize: 100,
	}
	for _, o := range opts {
		o(&options)
	}
	return options
}
