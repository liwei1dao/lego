package single

import (
	"github.com/liwei1dao/lego/sys/redis/core"

	"github.com/go-redis/redis/v8"
)

/*
Redis Hdel 命令用于删除哈希表 key 中的一个或多个指定字段，不存在的字段将被忽略
*/
func (this *Redis) HDel(key string, fields ...string) (err error) {
	agrs := make([]interface{}, 0)
	agrs = append(agrs, "HDEL")
	agrs = append(agrs, key)
	for _, v := range fields {
		agrs = append(agrs, v)
	}
	err = this.client.Do(this.client.Context(), agrs...).Err()
	return
}

/*
Redis Hexists 命令用于查看哈希表的指定字段是否存在
*/
func (this *Redis) HExists(key string, field string) (result bool, err error) {
	result, err = this.client.Do(this.client.Context(), "HEXISTS", key, field).Bool()
	return
}

/*
Redis Hmset 命令用于同时将多个 field-value (字段-值)对设置到哈希表中。
此命令会覆盖哈希表中已存在的字段。
如果哈希表不存在，会创建一个空哈希表，并执行 HMSET 操作
*/
func (this *Redis) HMSet(key string, v interface{}) (err error) {
	agrs := make([]interface{}, 0)
	agrs = append(agrs, "HMSET")
	agrs = append(agrs, key)
	var data map[string]string
	if data, err = this.codec.MarshalMap(v); err != nil {
		return
	}
	for k, v := range data {
		agrs = append(agrs, k, v)
	}
	err = this.client.Do(this.client.Context(), agrs...).Err()
	return
}
func (this *Redis) HMSetForMap(key string, v map[string]string) (err error) {
	agrs := make([]interface{}, 0)
	agrs = append(agrs, "HMSET")
	agrs = append(agrs, key)
	for k, v := range v {
		agrs = append(agrs, k, v)
	}
	err = this.client.Do(this.client.Context(), agrs...).Err()
	return
}

/*
Redis Hget 命令用于返回哈希表中指定字段的值
*/
func (this *Redis) HGet(key string, field string, v interface{}) (err error) {
	cmd := redis.NewStringCmd(this.client.Context(), "HGET", key, field)
	this.client.Process(this.client.Context(), cmd)
	var _result []byte
	if _result, err = cmd.Bytes(); err == nil {
		err = this.codec.Unmarshal(_result, v)
	}
	return
}

/*
Redis Hgetall 命令用于返回哈希表中，所有的字段和值。
在返回值里，紧跟每个字段名(field name)之后是字段的值(value)，所以返回值的长度是哈希表大小的两倍
*/
func (this *Redis) HGetAll(key string, v interface{}) (err error) {
	cmd := redis.NewStringStringMapCmd(this.client.Context(), "HGETALL", key)
	this.client.Process(this.client.Context(), cmd)
	var _result map[string]string
	if _result, err = cmd.Result(); err == nil {
		if len(_result) == 0 {
			return redis.Nil
		}
		err = this.codec.UnmarshalMap(_result, v)
	}
	return
}

/*
	读取全部hash集合数据到map中
*/
func (this *Redis) HGetAllToMapString(key string) (result map[string]string, err error) {
	cmd := redis.NewStringStringMapCmd(this.client.Context(), "HGETALL", key)
	this.client.Process(this.client.Context(), cmd)
	if result, err = cmd.Result(); err == nil {
		if len(result) == 0 {
			err = redis.Nil
			return
		}
	}
	return
}

/*
Redis Hincrby 命令用于为哈希表中的字段值加上指定增量值。
增量也可以为负数，相当于对指定字段进行减法操作。
如果哈希表的 key 不存在，一个新的哈希表被创建并执行 HINCRBY 命令。
如果指定的字段不存在，那么在执行命令前，字段的值被初始化为 0 。
对一个储存字符串值的字段执行 HINCRBY 命令将造成一个错误。
本操作的值被限制在 64 位(bit)有符号数字表示之内
*/
func (this *Redis) HIncrBy(key string, field string, value int) (err error) {
	err = this.client.Do(this.client.Context(), "HINCRBY", key, field, value).Err()
	return
}

/*
Redis Hincrbyfloat 命令用于为哈希表中的字段值加上指定浮点数增量值。
如果指定的字段不存在，那么在执行命令前，字段的值被初始化为 0
*/
func (this *Redis) HIncrByFloat(key string, field string, value float32) (err error) {
	err = this.client.Do(this.client.Context(), "HINCRBYFLOAT", key, field, value).Err()
	return
}

/*
Redis Hkeys 命令用于获取哈希表中的所有域(field)
*/
func (this *Redis) Hkeys(key string) (result []string, err error) {
	cmd := redis.NewStringSliceCmd(this.client.Context(), "HKEYS", key)
	this.client.Process(this.client.Context(), cmd)
	result, err = cmd.Result()
	return
}

/*
Redis Hlen 命令用于获取哈希表中字段的数量
*/
func (this *Redis) Hlen(key string) (result int, err error) {
	result, err = this.client.Do(this.client.Context(), "HLEN", key).Int()
	return
}

/*
Redis Hmget 命令用于返回哈希表中，一个或多个给定字段的值。
如果指定的字段不存在于哈希表，那么返回一个 nil 值
*/
func (this *Redis) HMGet(key string, v interface{}, fields ...string) (err error) {
	agrs := make([]interface{}, 0)
	agrs = append(agrs, "HMGET")
	agrs = append(agrs, key)
	for _, v := range fields {
		agrs = append(agrs, v)
	}
	cmd := redis.NewStringStringMapCmd(this.client.Context(), agrs...)
	this.client.Process(this.client.Context(), cmd)
	var _result map[string]string
	if _result, err = cmd.Result(); err == nil {
		err = this.codec.UnmarshalMap(_result, v)
	}
	return
}

/*
Redis Hset 命令用于为哈希表中的字段赋值
如果哈希表不存在，一个新的哈希表被创建并进行 HSET 操作
如果字段已经存在于哈希表中，旧值将被覆盖
*/
func (this *Redis) HSet(key string, field string, value interface{}) (err error) {
	var resultvalue []byte
	if !core.IsBaseType(value) {
		if resultvalue, err = this.codec.Marshal(value); err == nil {
			err = this.client.Do(this.client.Context(), "HSET", key, field, resultvalue).Err()
		}
	} else {
		err = this.client.Do(this.client.Context(), "HSET", key, field, value).Err()
	}
	return
}

/*
Redis Hsetnx 命令用于为哈希表中不存在的的字段赋值
如果哈希表不存在，一个新的哈希表被创建并进行 HSET 操作
如果字段已经存在于哈希表中，操作无效
如果 key 不存在，一个新哈希表被创建并执行 HSETNX 命令
*/
func (this *Redis) HSetNX(key string, field string, value interface{}) (err error) {
	var resultvalue []byte
	if resultvalue, err = this.codec.Marshal(value); err == nil {
		err = this.client.Do(this.client.Context(), "HSETNX", key, field, resultvalue).Err()
	}
	return
}
