package redis

import (
	"time"

	"github.com/liwei1dao/lego/utils/mapstructure"
)

type RedisType int8

const (
	Redis_Single RedisType = iota
	Redis_Cluster
)

///redis 存储数据格式化类型
type RedisStorageTyoe int8

const (
	JsonData RedisStorageTyoe = iota
	ProtoData
)

type Option func(*Options)
type Options struct {
	RedisType              RedisType
	Redis_Single_Addr      string
	Redis_Single_Password  string
	Redis_Single_DB        int
	Redis_Single_PoolSize  int
	Redis_Cluster_Addr     []string
	Redis_Cluster_Password string
	RedisStorageType       RedisStorageTyoe
	TimeOut                time.Duration
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

func SetRedis_Single_PoolSize(v int) Option {
	return func(o *Options) {
		o.Redis_Single_PoolSize = v
	}
}
func Redis_Cluster_Addr(v []string) Option {
	return func(o *Options) {
		o.Redis_Cluster_Addr = v
	}
}

func SetRedis_Cluster_Password(v string) Option {
	return func(o *Options) {
		o.Redis_Cluster_Password = v
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
		Redis_Single_Addr:      "127.0.0.1:6379",
		Redis_Single_Password:  "",
		Redis_Single_DB:        1,
		Redis_Cluster_Addr:     []string{"127.0.0.1:6379"},
		Redis_Cluster_Password: "",
		TimeOut:                time.Second * 3,
		Redis_Single_PoolSize:  100,
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
		Redis_Single_Addr:      "127.0.0.1:6379",
		Redis_Single_Password:  "",
		Redis_Single_DB:        1,
		Redis_Cluster_Addr:     []string{"127.0.0.1:6379"},
		Redis_Cluster_Password: "",
		TimeOut:                time.Second * 3,
		Redis_Single_PoolSize:  100,
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
