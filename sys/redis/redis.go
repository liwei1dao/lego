package redis

import (
	"context"
	"encoding/json"

	"github.com/go-redis/redis/v8"
	"google.golang.org/protobuf/proto"
)

func newSys(options Options) (sys *Redis, err error) {
	sys = &Redis{options: options}
	err = sys.init()
	return
}

type Redis struct {
	options Options
	client  *redis.Client
}

func (this *Redis) init() (err error) {
	this.client = redis.NewClient(&redis.Options{
		Addr:     this.options.RedisUrl,
		Password: this.options.RedisPassword,
		DB:       this.options.RedisDB,
		PoolSize: this.options.PoolSize, // 连接池大小
	})
	_, err = this.client.Ping(this.getContext()).Result()
	return
}

func (this *Redis) getContext() (ctx context.Context) {
	ctx, _ = context.WithTimeout(context.Background(), this.options.TimeOut)
	return
}

///批处理
func (this *Redis) Pipeline(ctx context.Context, fn func(pipe redis.Pipeliner) error) (err error) {
	_, err = this.client.Pipelined(ctx, fn)
	return
}

///事务
func (this *Redis) TxPipelined(ctx context.Context, fn func(pipe redis.Pipeliner) error) (err error) {
	_, err = this.client.TxPipelined(ctx, fn)
	return
}

///监控
func (this *Redis) Watch(ctx context.Context, fn func(*redis.Tx) error, keys ...string) (err error) {
	agrs := make([]string, len(keys))
	for i, v := range keys {
		agrs[i] = string(v)
	}
	err = this.client.Watch(ctx, fn, agrs...)
	return
}

///数据编码
func (this *Redis) Encode(value interface{}) (result []byte, err error) {
	if this.options.RedisStorageType == JsonData {
		result, err = json.Marshal(value)
	} else {
		if _, ok := value.(proto.Message); ok {
			result, err = proto.Marshal(value.(proto.Message))
		} else {
			result, err = json.Marshal(value)
		}
	}
	return
}

func (this *Redis) Decode(value []byte, result interface{}) (err error) {
	if this.options.RedisStorageType == JsonData {
		err = json.Unmarshal(value, result)
	} else {
		if _, ok := result.(proto.Message); ok {
			err = proto.Unmarshal(value, result.(proto.Message))
		} else {
			err = json.Unmarshal(value, result)
		}
	}
	return
}
