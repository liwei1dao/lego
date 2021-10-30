package redis

import (
	"errors"
	"time"

	"github.com/go-redis/redis/v8"
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

func (this *Redis) Lock(key string, outTime int) (result bool, err error) {
	cmd := redis.NewBoolCmd(this.getContext(), "set", key, 1, "ex", outTime, "nx")
	this.client.Process(this.getContext(), cmd)
	result, err = cmd.Result()
	return
}
func (this *Redis) UnLock(key string) (err error) {
	err = this.Delete(key)
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
