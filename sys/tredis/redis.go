package redis

import (
	"context"

	"github.com/go-redis/redis/v8"
)

func newSys(options Options) (sys *Redis, err error) {
	sys = &Redis{options: options}
	err = sys.init()
	return
}

type Redis struct {
	options Options
	client  *redis.Client
	ctx     context.Context
}

func (this *Redis) init() (err error) {
	var opt *redis.Options
	opt, err = redis.ParseURL(this.options.RedisUrl)
	if err != nil {
		return
	}
	this.client = redis.NewClient(opt)
	this.ctx, _ = context.WithTimeout(context.Background(), this.options.TimeOut)
	return
}

//判断键是否存在
func (this *Redis) ContainsKey(key string) (iskeep bool, err error) {
	
	// pool := this.Pool.Get()
	// defer pool.Close()
	// iskeep, err = redis.Bool(pool.Do("EXISTS", key))
	return
}
