package model

import (
	"time"

	"github.com/liwei1dao/lego/sys/redis"
)

//文档模型
type Document struct {
	Model     *Model
	TableName string
	Expired   time.Duration //过期时间
}

//创建锁对象
func (this *Document) NewRedisMutex(key string) (result *redis.RedisMutex, err error) {
	result, err = this.Model.Redis.NewRedisMutex(key)
	return
}

//创建锁对象 自定义过期时间
func (this *Document) NewRedisMutexForExpiry(key string, outtime int) (result *redis.RedisMutex, err error) {
	result, err = this.Model.Redis.NewRedisMutex(key)
	return
}
