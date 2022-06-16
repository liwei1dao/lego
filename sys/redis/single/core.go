package single

import (
	"context"
	"time"

	"github.com/liwei1dao/lego/utils/codec"

	"github.com/go-redis/redis/v8"
)

func NewSys(RedisUrl, RedisPassword string, RedisDB, PoolSize int, timeOut time.Duration,
	encode codec.IEncoder,
	decode codec.IDecoder,
) (sys *Redis, err error) {
	var (
		client *redis.Client
	)
	client = redis.NewClient(&redis.Options{
		Addr:     RedisUrl,
		Password: RedisPassword,
		DB:       RedisDB,
		PoolSize: PoolSize, // 连接池大小
	})
	sys = &Redis{
		client:  client,
		timeOut: timeOut,
		encode:  encode,
		decode:  decode,
	}
	_, err = sys.Ping()
	return
}

type Redis struct {
	client  *redis.Client
	timeOut time.Duration
	encode  codec.IEncoder
	decode  codec.IDecoder
}

func (this *Redis) getContext() (ctx context.Context) {
	ctx, _ = context.WithTimeout(context.Background(), this.timeOut)
	return
}
func (this *Redis) Close() (err error) {
	return this.client.Close()
}

/// Ping
func (this *Redis) Ping() (string, error) {
	return this.client.Ping(this.getContext()).Result()
}

/// 命令接口
func (this *Redis) Do(ctx context.Context, args ...interface{}) *redis.Cmd {
	return this.client.Do(ctx, args...)
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

//锁
func (this *Redis) Lock(key string, outTime int) (result bool, err error) {
	cmd := redis.NewBoolCmd(this.getContext(), "set", key, 1, "ex", outTime, "nx")
	this.client.Process(this.getContext(), cmd)
	result, err = cmd.Result()
	return
}

//锁
func (this *Redis) UnLock(key string) (err error) {
	err = this.Delete(key)
	return
}
