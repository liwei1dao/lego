package redis

import (
	"github.com/liwei1dao/lego/core"
)

var (
	deffactory IRedisFactory
)

type (
	IRedisFactory interface {
		GetPool() *RedisPool
		CloseAllPool()
	}
)

func OnInit(s core.IService, opt ...Option) (err error) {
	deffactory = newRedisFactory()
	return
}

func NewRedisSys(opt ...Option) (factory IRedisFactory, err error) {
	factory = newRedisFactory()
	return
}

func GetPool() *RedisPool {
	return deffactory.GetPool()
}

func CloseAllPool() {
	deffactory.CloseAllPool()
}
