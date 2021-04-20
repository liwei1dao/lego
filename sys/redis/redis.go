package redis

import (
	"context"
	"fmt"

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
}

func (this *Redis) init() (err error) {
	this.client = redis.NewClient(&redis.Options{
		Addr:     this.options.RedisUrl,
		Password: this.options.RedisPassword,
		DB:       this.options.RedisDB,
	})
	_, err = this.client.Ping(this.getContext()).Result()
	return
}

func (this *Redis) getContext() (ctx context.Context) {
	ctx, _ = context.WithTimeout(context.Background(), this.options.TimeOut)
	return
}

///判断是否存在key
func (this *Redis) ContainsKey(key string) (iskeep bool, err error) {
	iskeep, err = this.client.Do(this.getContext(), "EXISTS", key).Bool()
	return
}

///判断是否存在key
func (this *Redis) QueryPatternKeys(key string) (keys []string, err error) {
	cmd := redis.NewStringSliceCmd(this.getContext(), "KEYS", fmt.Sprintf("%s*", key))
	this.client.Process(this.getContext(), cmd)
	keys, err = cmd.Result()
	return
}

//删除Redis 缓存键数据
func (this *Redis) Delete(key string) (err error) {
	err = this.client.Del(this.getContext(), key).Err()
	return
}

//添加键值对
func (this *Redis) SetKeyForValue(key string, value interface{}) (err error) {
	// this.client.Set(this.getContext(), key, value)
	// if b, err := json.Marshal(value); err == nil {
	// 	pool := this.Pool.Get()
	// 	defer pool.Close()
	// 	_, err = pool.Do("SET", key, b)
	// }
	return
}
