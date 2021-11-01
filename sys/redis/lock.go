package redis

import (
	"errors"
	"time"
)

/*
Redis Scard 命令返回集合中元素的数量
*/
func (this *Redis) NewRedisMutex(key string, opt ...RMutexOption) (result *RedisMutex, err error) {
	opts := newRMutexOptions(opt...)
	result = &RedisMutex{
		sys:    this.client,
		key:    key,
		expiry: opts.expiry,
		delay:  opts.delay,
	}
	return
}

type RedisMutex struct {
	sys    IRedis
	key    string
	expiry int //过期时间 单位秒
	delay  time.Duration
}

//此接口未阻塞接口
func (this *RedisMutex) Lock() (err error) {
	wait := make(chan error)
	go func() {
		start := time.Now()
		for int(time.Now().Sub(start).Seconds()) <= this.expiry {
			if result, err := this.sys.Lock(this.key, this.expiry); err == nil && result {
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

func (this *RedisMutex) Unlock() {
	this.sys.UnLock(this.key)
}
