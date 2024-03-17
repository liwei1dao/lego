package pipe

import (
	"context"

	"github.com/liwei1dao/lego/sys/redis/core"

	"github.com/go-redis/redis/v8"
)

func NewPipe(ctx context.Context, pipe redis.Pipeliner, codec core.ICodec) *RedisPipe {
	return &RedisPipe{
		ctx:    ctx,
		client: pipe,
		codec:  codec,
	}
}

type RedisPipe struct {
	ctx    context.Context
	client redis.Pipeliner
	codec  core.ICodec
}

func (this *RedisPipe) Exec() ([]redis.Cmder, error) {
	return this.client.Exec(this.ctx)
}

/// 命令接口
func (this *RedisPipe) Do(args ...interface{}) *redis.Cmd {
	return this.client.Do(this.ctx, args...)
}

///批处理
func (this *RedisPipe) Pipeline() redis.Pipeliner {
	return this.client.Pipeline()
}

///批处理
func (this *RedisPipe) Pipelined(fn func(pipe redis.Pipeliner) error) (err error) {
	_, err = this.client.Pipelined(this.ctx, fn)
	return
}

///事务
func (this *RedisPipe) TxPipelined(fn func(pipe redis.Pipeliner) error) (err error) {
	_, err = this.client.TxPipelined(this.ctx, fn)
	return
}

//锁
func (this *RedisPipe) Lock(key string, outTime int) (result bool, err error) {
	cmd := redis.NewBoolCmd(this.ctx, "set", key, 1, "ex", outTime, "nx")
	this.client.Process(this.ctx, cmd)
	result, err = cmd.Result()
	return
}

//锁
func (this *RedisPipe) UnLock(key string) (err error) {
	err = this.Delete(key)
	return
}

//lua Script
func (this *RedisPipe) NewScript(src string) *redis.StringCmd {
	script := redis.NewScript(src)
	return script.Load(this.ctx, this.client)
}
func (this *RedisPipe) Eval(ctx context.Context, script string, keys []string, args ...interface{}) *redis.Cmd {
	return this.client.Eval(ctx, script, keys, args...)
}
func (this *RedisPipe) EvalSha(ctx context.Context, sha1 string, keys []string, args ...interface{}) *redis.Cmd {
	return this.client.EvalSha(ctx, sha1, keys, args...)
}
func (this *RedisPipe) ScriptExists(ctx context.Context, hashes ...string) *redis.BoolSliceCmd {
	return this.client.ScriptExists(ctx, hashes...)
}
