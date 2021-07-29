package redis

import (
	"errors"
	"time"

	"github.com/liwei1dao/lego/core"
)

/*
Redis Scard 命令返回集合中元素的数量
*/
func (this *Redis) NewRedisMutex(key core.Redis_Key, opt ...RMutexOption) (result *RedisMutex, err error) {
	opts := newRMutexOptions(opt...)
	result = &RedisMutex{
		sys:    this,
		key:    key,
		expiry: opts.expiry,
		delay:  opts.delay,
	}
	return
}

func (this *Redis) Lock(key core.Redis_Key, outTime int) (err error) {
	err = this.client.Do(this.getContext(), "set", key, 1, "ex", outTime, "nx").Err()
	return
}
func (this *Redis) UnLock(key core.Redis_Key) (err error) {
	err = this.client.Do(this.getContext(), "del", key).Err()
	return
}

type RedisMutex struct {
	sys    IRedis
	key    core.Redis_Key
	expiry int //过期时间 单位秒
	delay  time.Duration
}

//此接口未阻塞接口
func (this *RedisMutex) Lock() (err error) {
	wait := make(chan error)
	go func() {
		start := time.Now()
		for int(time.Now().Sub(start).Seconds()) <= this.expiry {
			if err := this.sys.Lock(this.key, this.expiry); err == nil {
				wait <- nil
				return
			} else {
				time.Sleep(this.delay)
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
