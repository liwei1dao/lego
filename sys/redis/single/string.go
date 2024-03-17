package single

import (
	"time"

	"github.com/go-redis/redis/v8"
)

/* String *******************************************************************************/
/*
命令用于设置给定 key 的值。如果 key 已经存储其他值， SET 就覆写旧值，且无视类型。
*/
func (this *Redis) Set(key string, value interface{}, expiration time.Duration) (err error) {
	var result []byte
	if result, err = this.codec.Marshal(value); err != nil {
		return
	}
	err = this.client.Set(this.client.Context(), key, result, expiration).Err()
	return
}

/*
指定的 key 不存在时，为 key 设置指定的值
*/
func (this *Redis) SetNX(key string, value interface{}) (result int64, err error) {
	cmd := redis.NewIntCmd(this.client.Context(), "SETNX", key, value)
	this.client.Process(this.client.Context(), cmd)
	result, err = cmd.Result()
	// }
	return
}

/*
同时设置一个或多个 key-value 对
*/
func (this *Redis) MSet(v map[string]interface{}) (err error) {
	agrs := make([]interface{}, 0)
	agrs = append(agrs, "MSET")
	for k, v := range v {
		result, _ := this.codec.Marshal(v)
		agrs = append(agrs, k, result)
	}
	err = this.client.Do(this.client.Context(), agrs...).Err()
	return
}

/*
命令用于所有给定 key 都不存在时，同时设置一个或多个 key-value 对
*/
func (this *Redis) MSetNX(v map[string]interface{}) (err error) {
	agrs := make([]interface{}, 0)
	agrs = append(agrs, "MSETNX")
	for k, v := range v {
		result, _ := this.codec.Marshal(v)
		agrs = append(agrs, k, result)
	}
	err = this.client.Do(this.client.Context(), agrs...).Err()
	return
}

/*
Redis Incr 命令将 key 中储存的数字值增一。
如果 key 不存在，那么 key 的值会先被初始化为 0 ，然后再执行 INCR 操作。
如果值包含错误的类型，或字符串类型的值不能表示为数字，那么返回一个错误。
本操作的值限制在 64 位(bit)有符号数字表示之内。
*/
func (this *Redis) Incr(key string) (err error) {
	err = this.client.Do(this.client.Context(), "INCR", key).Err()
	return
}

/*
Redis Incrby 命令将 key 中储存的数字加上指定的增量值。
如果 key 不存在，那么 key 的值会先被初始化为 0 ，然后再执行 INCRBY 命令。
如果值包含错误的类型，或字符串类型的值不能表示为数字，那么返回一个错误。
本操作的值限制在 64 位(bit)有符号数字表示之内
*/
func (this *Redis) IncrBY(key string, value int) (err error) {
	err = this.client.Do(this.client.Context(), "INCRBY", key, value).Err()
	return
}

/*
Redis Incrbyfloat 命令为 key 中所储存的值加上指定的浮点数增量值。
如果 key 不存在，那么 INCRBYFLOAT 会先将 key 的值设为 0 ，再执行加法操作
*/
func (this *Redis) Incrbyfloat(key string, value float32) (err error) {
	err = this.client.Do(this.client.Context(), "INCRBYFLOAT", key, value).Err()
	return
}

/*
Redis Decr 命令将 key 中储存的数字值减一。
如果 key 不存在，那么 key 的值会先被初始化为 0 ，然后再执行 DECR 操作。
如果值包含错误的类型，或字符串类型的值不能表示为数字，那么返回一个错误。
本操作的值限制在 64 位(bit)有符号数字表示之内
*/
func (this *Redis) Decr(key string, value int) (err error) {
	err = this.client.Do(this.client.Context(), "DECR", key, value).Err()
	return
}

/*
Redis Decrby 命令将 key 所储存的值减去指定的减量值。
如果 key 不存在，那么 key 的值会先被初始化为 0 ，然后再执行 DECRBY 操作。
如果值包含错误的类型，或字符串类型的值不能表示为数字，那么返回一个错误。
本操作的值限制在 64 位(bit)有符号数字表示之内
*/
func (this *Redis) DecrBy(key string, value int) (err error) {
	err = this.client.Do(this.client.Context(), "DECRBY", key, value).Err()
	return
}

/*
Redis Append 命令用于为指定的 key 追加值。
如果 key 已经存在并且是一个字符串， APPEND 命令将 value 追加到 key 原来的值的末尾。
如果 key 不存在， APPEND 就简单地将给定 key 设为 value ，就像执行 SET key value 一样。
*/
func (this *Redis) Append(key string, value interface{}) (err error) {
	var result []byte
	if result, err = this.codec.Marshal(value); err != nil {
		return
	}
	err = this.client.Do(this.client.Context(), "APPEND", key, result).Err()
	return
}

/*
命令用于设置给定 key 的值。如果 key 已经存储其他值， SET 就覆写旧值，且无视类型
*/
func (this *Redis) Get(key string, value interface{}) (err error) {
	var result []byte
	if result, err = this.client.Get(this.client.Context(), key).Bytes(); err == nil {
		err = this.codec.Unmarshal(result, value)
	}
	return
}

/*
设置指定 key 的值，并返回 key 的旧值
*/
func (this *Redis) GetSet(key string, value interface{}, result interface{}) (err error) {
	var (
		_value []byte
	)
	if _value, err = this.codec.Marshal(value); err == nil {
		cmd := redis.NewStringCmd(this.client.Context(), "GETSET", key, _value)
		this.client.Process(this.client.Context(), cmd)
		var _result []byte
		if _result, err = cmd.Bytes(); err == nil {
			err = this.codec.Unmarshal(_result, result)
		}
	}
	return
}

/*
返回所有(一个或多个)给定 key 的值。 如果给定的 key 里面，有某个 key 不存在，那么这个 key 返回特殊值 nil
*/
func (this *Redis) MGet(v interface{}, keys ...string) (err error) {
	agrs := make([]interface{}, 0)
	agrs = append(agrs, "MGET")
	for _, v := range keys {
		agrs = append(agrs, v)
	}
	cmd := redis.NewStringSliceCmd(this.client.Context(), agrs...)
	this.client.Process(this.client.Context(), cmd)
	var result []string
	if result, err = cmd.Result(); err != nil {
		return
	}
	err = this.codec.UnmarshalSlice(result, v)
	return
}

///判断是否存在key pattern:key*
func (this *Redis) INCRBY(key string, amount int64) (result int64, err error) {
	cmd := redis.NewIntCmd(this.client.Context(), "INCRBY", key, amount)
	this.client.Process(this.client.Context(), cmd)
	result, err = cmd.Result()
	return
}
