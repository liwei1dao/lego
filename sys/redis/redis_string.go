package redis

import (
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/liwei1dao/lego/core"
)

/* String *******************************************************************************/
/*
命令用于设置给定 key 的值。如果 key 已经存储其他值， SET 就覆写旧值，且无视类型。
*/
func (this *Redis) Set(key core.Redis_Key, value interface{}, expiration time.Duration) (err error) {
	var result []byte
	if result, err = this.encode(value); err == nil {
		err = this.client.Set(this.getContext(), string(key), result, expiration).Err()
	}
	return
}

/*
指定的 key 不存在时，为 key 设置指定的值
*/
func (this *Redis) SetNX(key core.Redis_Key, value interface{}) (err error) {
	var result []byte
	if result, err = this.encode(value); err == nil {
		err = this.client.Do(this.getContext(), "SETNX", key, result).Err()
	}
	return
}

/*
同时设置一个或多个 key-value 对
*/
func (this *Redis) MSet(keyvalues map[core.Redis_Key]interface{}) (err error) {
	agrs := make([]interface{}, 0)
	agrs = append(agrs, "MSET")
	for k, v := range keyvalues {
		result, _ := this.encode(v)
		agrs = append(agrs, k, result)
	}
	err = this.client.Do(this.getContext(), agrs...).Err()
	return
}

/*
命令用于所有给定 key 都不存在时，同时设置一个或多个 key-value 对
*/
func (this *Redis) MSetNX(keyvalues map[core.Redis_Key]interface{}) (err error) {
	agrs := make([]interface{}, 0)
	agrs = append(agrs, "MSETNX")
	for k, v := range keyvalues {
		result, _ := this.encode(v)
		agrs = append(agrs, k, result)
	}
	err = this.client.Do(this.getContext(), agrs...).Err()
	return
}

/*
Redis Incr 命令将 key 中储存的数字值增一。
如果 key 不存在，那么 key 的值会先被初始化为 0 ，然后再执行 INCR 操作。
如果值包含错误的类型，或字符串类型的值不能表示为数字，那么返回一个错误。
本操作的值限制在 64 位(bit)有符号数字表示之内。
*/
func (this *Redis) Incr(key core.Redis_Key) (err error) {
	err = this.client.Do(this.getContext(), "INCR", key).Err()
	return
}

/*
Redis Incrby 命令将 key 中储存的数字加上指定的增量值。
如果 key 不存在，那么 key 的值会先被初始化为 0 ，然后再执行 INCRBY 命令。
如果值包含错误的类型，或字符串类型的值不能表示为数字，那么返回一个错误。
本操作的值限制在 64 位(bit)有符号数字表示之内
*/
func (this *Redis) IncrBY(key core.Redis_Key, value int) (err error) {
	err = this.client.Do(this.getContext(), "INCRBY", key, value).Err()
	return
}

/*
Redis Incrbyfloat 命令为 key 中所储存的值加上指定的浮点数增量值。
如果 key 不存在，那么 INCRBYFLOAT 会先将 key 的值设为 0 ，再执行加法操作
*/
func (this *Redis) Incrbyfloat(key core.Redis_Key, value float32) (err error) {
	err = this.client.Do(this.getContext(), "INCRBYFLOAT", key, value).Err()
	return
}

/*
Redis Decr 命令将 key 中储存的数字值减一。
如果 key 不存在，那么 key 的值会先被初始化为 0 ，然后再执行 DECR 操作。
如果值包含错误的类型，或字符串类型的值不能表示为数字，那么返回一个错误。
本操作的值限制在 64 位(bit)有符号数字表示之内
*/
func (this *Redis) Decr(key core.Redis_Key, value int) (err error) {
	err = this.client.Do(this.getContext(), "DECR", key, value).Err()
	return
}

/*
Redis Decrby 命令将 key 所储存的值减去指定的减量值。
如果 key 不存在，那么 key 的值会先被初始化为 0 ，然后再执行 DECRBY 操作。
如果值包含错误的类型，或字符串类型的值不能表示为数字，那么返回一个错误。
本操作的值限制在 64 位(bit)有符号数字表示之内
*/
func (this *Redis) DecrBy(key core.Redis_Key, value int) (err error) {
	err = this.client.Do(this.getContext(), "DECRBY", key, value).Err()
	return
}

/*
Redis Append 命令用于为指定的 key 追加值。
如果 key 已经存在并且是一个字符串， APPEND 命令将 value 追加到 key 原来的值的末尾。
如果 key 不存在， APPEND 就简单地将给定 key 设为 value ，就像执行 SET key value 一样。
*/
func (this *Redis) Append(key core.Redis_Key, value interface{}) (err error) {
	var result []byte
	if result, err = this.encode(value); err == nil {
		err = this.client.Do(this.getContext(), "APPEND", key, result).Err()
	}
	return
}

/*
命令用于设置给定 key 的值。如果 key 已经存储其他值， SET 就覆写旧值，且无视类型
*/
func (this *Redis) Get(key core.Redis_Key, value interface{}) (err error) {
	var result []byte
	if result, err = this.client.Get(this.getContext(), string(key)).Bytes(); err == nil {
		err = this.decode(result, value)
	}
	return
}

/*
设置指定 key 的值，并返回 key 的旧值
*/
func (this *Redis) GetSet(key core.Redis_Key, value interface{}, result interface{}) (err error) {
	var (
		data   string
		_value []byte
	)
	if _value, err = this.encode(value); err == nil {
		if data = this.client.Do(this.getContext(), "GETSET", key, _value).String(); data != string(redis.Nil) {
			err = this.decode([]byte(data), result)
		} else {
			err = fmt.Errorf(string(redis.Nil))
		}
	}
	return
}

/*
返回所有(一个或多个)给定 key 的值。 如果给定的 key 里面，有某个 key 不存在，那么这个 key 返回特殊值 nil
*/
func (this *Redis) MGet(keys ...core.Redis_Key) (result []string, err error) {
	agrs := make([]interface{}, 0)
	agrs = append(agrs, "MGET")
	for _, v := range keys {
		agrs = append(agrs, v)
	}
	cmd := redis.NewStringSliceCmd(this.getContext(), agrs...)
	this.client.Process(this.getContext(), cmd)
	result, err = cmd.Result()
	return
}
