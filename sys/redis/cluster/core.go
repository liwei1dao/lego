package cluster

import (
	"context"
	"time"

	"github.com/liwei1dao/lego/sys/redis/core"

	"github.com/go-redis/redis/v8"
)

func NewSys(RedisUrl []string, RedisPassword string, timeOut time.Duration,
	codec core.ICodec,
) (sys *Redis, err error) {
	var (
		client *redis.ClusterClient
	)
	client = redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:        RedisUrl,
		Password:     RedisPassword,
		DialTimeout:  timeOut,
		ReadTimeout:  timeOut,
		WriteTimeout: timeOut,
	})
	sys = &Redis{
		client:  client,
		timeOut: timeOut,
		codec:   codec,
	}
	_, err = sys.Ping()
	return
}

type Redis struct {
	client  *redis.ClusterClient
	timeOut time.Duration
	codec   core.ICodec
}

///事务
func (this *Redis) Close() (err error) {
	err = this.client.Close()
	return
}

/// Context
func (this *Redis) Context() context.Context {
	return this.client.Context()
}

/// Ping
func (this *Redis) Ping() (string, error) {
	return this.client.Ping(this.client.Context()).Result()
}

/// 命令接口
func (this *Redis) Do(ctx context.Context, args ...interface{}) *redis.Cmd {
	return this.client.Do(ctx, args...)
}

///批处理
func (this *Redis) Pipeline() redis.Pipeliner {
	return this.client.Pipeline()
}

///批处理
func (this *Redis) Pipelined(ctx context.Context, fn func(pipe redis.Pipeliner) error) (err error) {
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
	cmd := redis.NewBoolCmd(this.client.Context(), "set", key, 1, "ex", outTime, "nx")
	this.client.Process(this.client.Context(), cmd)
	result, err = cmd.Result()
	return
}

//锁
func (this *Redis) UnLock(key string) (err error) {
	err = this.Delete(key)
	return
}

//lua Script
func (this *Redis) NewScript(src string) *redis.StringCmd {
	script := redis.NewScript(src)
	return script.Load(this.Context(), this.client)
}
func (this *Redis) Eval(ctx context.Context, script string, keys []string, args ...interface{}) *redis.Cmd {
	return this.client.Eval(ctx, script, keys, args...)
}
func (this *Redis) EvalSha(ctx context.Context, sha1 string, keys []string, args ...interface{}) *redis.Cmd {
	return this.client.EvalSha(ctx, sha1, keys, args...)
}
func (this *Redis) ScriptExists(ctx context.Context, hashes ...string) *redis.BoolSliceCmd {
	return this.client.ScriptExists(ctx, hashes...)
}

// func (this *Redis) ScriptLoad(ctx context.Context, script string) *redis.StringCmd {
// 	return this.client.ScriptLoad(ctx, script)
// }
