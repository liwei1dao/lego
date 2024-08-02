package redis

import (
	"context"
	"errors"
	"time"
)

/*
Redis Scard 命令返回集合中元素的数量
*/
func (this *Redis) NewRedisMutex(key string, opt ...RMutexOption) (result *RedisMutex, err error) {
	opts := newRMutexOptions(opt...)
	result = &RedisMutex{
		sys:    this,
		key:    key,
		expiry: opts.expiry,
		delay:  opts.delay,
	}
	return
}

type RedisMutex struct {
	sys    ISys
	key    string
	expiry time.Duration //过期时间 单位秒
	delay  time.Duration
}

// 此接口未阻塞接口
func (this *RedisMutex) Lock(ctx context.Context) (err error) {
	wait := make(chan error)
	go func() {
		start := time.Now()
		for time.Now().Sub(start) <= this.expiry {
			if result, err := this.sys.Lock(ctx, this.key, this.expiry); err == nil && result {
				wait <- nil
				return
			} else if err == nil && !result {
				time.Sleep(this.delay)
			} else {
				wait <- err
				return
			}
		}
		wait <- errors.New("time out")
	}()
	err = <-wait
	return
}

func (this *RedisMutex) Unlock(ctx context.Context) {
	this.sys.UnLock(ctx, this.key)
}
