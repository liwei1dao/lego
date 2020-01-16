package redis

import (
	"github.com/liwei1dao/lego/core"
	cont "github.com/liwei1dao/lego/utils/concurrent"
)

var (
	opts    *Options
	service core.IService
	factory *RedisFactory
)

func OnInit(s core.IService, opt ...Option) (err error) {
	service = s
	opts = newOptions(opt...)
	factory = &RedisFactory{
		pools: cont.NewBeeMap(),
	}
	return
}

func GetService() core.IService {
	return service
}

func GetPool() *RedisPool {
	return factory.GetPool(opts.RedisUrl)
}
