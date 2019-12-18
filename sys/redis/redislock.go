package redis

import (
	"github.com/pkg/errors"
	"time"
)

func NewRedisMutex(key string, opt ...RMutexOption) *RedisMutex {
	opts := newRMutexOptions(opt...)
	rm := &RedisMutex{
		key:    key,
		expiry: opts.expiry,
		delay:  opts.delay,
	}
	return rm
}

type RedisMutex struct {
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
			pool := GetPool()
			if err := pool.Lock(this.key, this.expiry); err == nil {
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
	pool := GetPool()
	pool.UnLock(this.key)
}
