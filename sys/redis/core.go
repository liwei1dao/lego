package redis

import (
	"context"
	"time"

	v9 "github.com/redis/go-redis/v9"
)

/*
系统描述:redis系统驱动,集成单接单和集群模式，以及之定义lua支持,并配合序列化系统批量处理消息数据
*/
type (
	ISys interface {
		GetClient() v9.UniversalClient
		NewRedisMutex(key string, opt ...RMutexOption) (result *RedisMutex, err error)
		Lock(ctx context.Context, key string, outTime time.Duration) (result bool, err error)
		UnLock(ctx context.Context, key string) (err error)
	}
)

const (
	RedisNil    = v9.Nil //数据为空错误
	TxFailedErr = v9.TxFailedErr
)

var defsys ISys

func OnInit(config map[string]interface{}, opt ...Option) (err error) {
	var option *Options
	if option, err = newOptions(config, opt...); err != nil {
		return
	}
	defsys, err = newSys(option)
	return
}

func NewSys(opt ...Option) (sys ISys, err error) {
	var option *Options
	if option, err = newOptionsByOption(opt...); err != nil {
		return
	}
	sys, err = newSys(option)
	return
}

func GetClient() v9.UniversalClient {
	return defsys.GetClient()
}

/*Lock*/
func NewRedisMutex(key string, opt ...RMutexOption) (result *RedisMutex, err error) {
	return defsys.NewRedisMutex(key, opt...)
}

func Lock(ctx context.Context, key string, outTime time.Duration) (result bool, err error) {
	return defsys.Lock(ctx, key, outTime)
}
func UnLock(ctx context.Context, key string) (err error) {
	return defsys.UnLock(ctx, key)
}
