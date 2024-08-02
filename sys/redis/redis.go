package redis

import (
	"context"
	"crypto/tls"
	"time"

	"github.com/liwei1dao/lego/sys/log"
	v9 "github.com/redis/go-redis/v9"
)

func newSys(options *Options) (sys *Redis, err error) {
	sys = &Redis{options: options}
	err = sys.init()
	return
}

type Redis struct {
	options *Options
	client  v9.UniversalClient
}

func (this *Redis) init() (err error) {
	rconf := &v9.UniversalOptions{
		Addrs:    this.options.RedisAddr,
		Password: this.options.RedisPassword, // 如果有密码
		DB:       this.options.RedisDB,
	}
	if this.options.RedisTLS {
		rconf.TLSConfig = &tls.Config{}
	}
	// 使用集群模式
	this.client = v9.NewUniversalClient(rconf)

	_, err = this.client.Ping(context.Background()).Result()
	if err != nil {
		this.options.Log.Error(err.Error(), log.Field{Key: "options", Value: this.options})
		return
	}
	return
}

func (this *Redis) GetClient() v9.UniversalClient {
	return this.client
}

func (this *Redis) Lock(ctx context.Context, key string, expiration time.Duration) (result bool, err error) {
	result, err = this.client.SetNX(ctx, key, 1, expiration).Result()
	return
}

func (this *Redis) UnLock(ctx context.Context, key string) (err error) {
	_, err = this.client.Del(ctx, key).Result()
	return
}
