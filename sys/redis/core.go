package redis

import (
	"reflect"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/liwei1dao/lego/core"
)

type (
	IRedis interface {
		Pipeline(fn func(pipe redis.Pipeliner) error) (err error)
		Watch(fn func(*redis.Tx) error, keys ...core.Redis_Key) (err error)
		/*Key*/
		Delete(key core.Redis_Key) (err error)
		ExistsKey(key core.Redis_Key) (iskeep bool, err error)
		ExpireKey(key core.Redis_Key, expire int) (err error)
		ExpireatKey(key core.Redis_Key, expire_unix int64) (err error)
		Pexpirekey(key core.Redis_Key, expire int) (err error)
		PexpireatKey(key core.Redis_Key, expire_unix int64) (err error)
		PersistKey(key core.Redis_Key) (err error)
		PttlKey(key core.Redis_Key) (leftexpire int64, err error)
		TtlKey(key core.Redis_Key) (leftexpire int64, err error)
		RenameKye(oldkey core.Redis_Key, newkey string) (err error)
		RenamenxKey(oldkey core.Redis_Key, newkey string) (err error)
		Keys(pattern core.Redis_Key) (keys []string, err error)
		/*String*/
		Set(key core.Redis_Key, value interface{}, expiration time.Duration) (err error)
		SetNX(key core.Redis_Key, value interface{}) (err error)
		MSet(keyvalues map[core.Redis_Key]interface{}) (err error)
		MSetNX(keyvalues map[core.Redis_Key]interface{}) (err error)
		Incr(key core.Redis_Key) (err error)
		IncrBY(key core.Redis_Key, value int) (err error)
		Incrbyfloat(key core.Redis_Key, value float32) (err error)
		Decr(key core.Redis_Key, value int) (err error)
		DecrBy(key core.Redis_Key, value int) (err error)
		Append(key core.Redis_Key, value interface{}) (err error)
		Get(key core.Redis_Key, value interface{}) (err error)
		GetSet(key core.Redis_Key, value interface{}, result interface{}) (err error)
		MGet(keys ...core.Redis_Key) (result []string, err error)
		/*List*/
		Lindex(key core.Redis_Key, value interface{}) (err error)
		Linsert(key core.Redis_Key, isbefore bool, tager interface{}, value interface{}) (err error)
		Llen(key core.Redis_Key) (result int, err error)
		LPop(key core.Redis_Key, value interface{}) (err error)
		LPush(key core.Redis_Key, values ...interface{}) (err error)
		LPushX(key core.Redis_Key, values ...interface{}) (err error)
		LRange(key core.Redis_Key, start, end int, valuetype reflect.Type) (result []interface{}, err error)
		LRem(key core.Redis_Key, count int, target interface{}) (err error)
		LSet(key core.Redis_Key, index int, value interface{}) (err error)
		Ltrim(key core.Redis_Key, start, stop int) (err error)
		Rpop(key core.Redis_Key, value interface{}) (err error)
		RPopLPush(oldkey core.Redis_Key, newkey core.Redis_Key, value interface{}) (err error)
		RPush(key core.Redis_Key, values ...interface{}) (err error)
		RPushX(key core.Redis_Key, values ...interface{}) (err error)
		/*Hash*/
		HDel(key core.Redis_Key, fields ...string) (err error)
		HExists(key core.Redis_Key, field string) (result bool, err error)
		HGet(key core.Redis_Key, field string, value interface{}) (err error)
		HGetAll(key core.Redis_Key, valuetype reflect.Type) (result []interface{}, err error)
		HIncrBy(key core.Redis_Key, field string, value int) (err error)
		HIncrByFloat(key core.Redis_Key, field string, value float32) (err error)
		Hkeys(key core.Redis_Key) (result []string, err error)
		Hlen(key core.Redis_Key) (result int, err error)
		HMGet(key core.Redis_Key, valuetype reflect.Type, fields ...string) (result []interface{}, err error)
		HMSet(key core.Redis_Key, value map[string]interface{}) (err error)
		HSet(key core.Redis_Key, field string, value interface{}) (err error)
		HSetNX(key core.Redis_Key, field string, value interface{}) (err error)
		/*Set*/
		SAdd(key core.Redis_Key, values ...interface{}) (err error)
		Scard(key core.Redis_Key) (result int, err error)
		Sismember(key core.Redis_Key, value interface{}) (iskeep bool, err error)
	}
)

const (
	RedisNil = redis.Nil //数据为空错误
)

var defsys IRedis

func OnInit(config map[string]interface{}, option ...Option) (err error) {
	defsys, err = newSys(newOptions(config, option...))
	return
}

func NewSys(option ...Option) (sys IRedis, err error) {
	sys, err = newSys(newOptionsByOption(option...))
	return
}

func Delete(key core.Redis_Key) (err error) {
	return defsys.Delete(key)
}
func ExistsKey(key core.Redis_Key) (iskeep bool, err error) {
	return defsys.ExistsKey(key)

}
func ExpireKey(key core.Redis_Key, expire int) (err error) {
	return defsys.ExpireKey(key, expire)
}
func ExpireatKey(key core.Redis_Key, expire_unix int64) (err error) {
	return defsys.ExpireatKey(key, expire_unix)
}
func Pexpirekey(key core.Redis_Key, expire int) (err error) {
	return defsys.Pexpirekey(key, expire)
}
func PexpireatKey(key core.Redis_Key, expire_unix int64) (err error) {
	return defsys.PexpireatKey(key, expire_unix)
}
func PersistKey(key core.Redis_Key) (err error) {
	return defsys.PersistKey(key)
}
func PttlKey(key core.Redis_Key) (leftexpire int64, err error) {
	return defsys.PttlKey(key)
}
func TtlKey(key core.Redis_Key) (leftexpire int64, err error) {
	return defsys.TtlKey(key)
}
func RenameKye(oldkey core.Redis_Key, newkey string) (err error) {
	return defsys.RenameKye(oldkey, newkey)
}
func RenamenxKey(oldkey core.Redis_Key, newkey string) (err error) {
	return defsys.RenamenxKey(oldkey, newkey)
}
func Keys(pattern core.Redis_Key) (keys []string, err error) {
	return defsys.Keys(pattern)
}

/*String*/
func Set(key core.Redis_Key, value interface{}, expiration time.Duration) (err error) {
	return defsys.Set(key, value, expiration)
}
func SetNX(key core.Redis_Key, value interface{}) (err error) {
	return defsys.SetNX(key, value)
}
func MSet(keyvalues map[core.Redis_Key]interface{}) (err error) {
	return defsys.MSet(keyvalues)
}
func MSetNX(keyvalues map[core.Redis_Key]interface{}) (err error) {
	return defsys.MSetNX(keyvalues)
}
func Incr(key core.Redis_Key) (err error) {
	return defsys.Incr(key)
}
func IncrBY(key core.Redis_Key, value int) (err error) {
	return defsys.IncrBY(key, value)
}
func Incrbyfloat(key core.Redis_Key, value float32) (err error) {
	return defsys.Incrbyfloat(key, value)
}
func Decr(key core.Redis_Key, value int) (err error) {
	return defsys.Decr(key, value)
}
func DecrBy(key core.Redis_Key, value int) (err error) {
	return defsys.DecrBy(key, value)
}
func Append(key core.Redis_Key, value interface{}) (err error) {
	return defsys.Append(key, value)
}
func Get(key core.Redis_Key, value interface{}) (err error) {
	return defsys.Get(key, value)
}
func GetSet(key core.Redis_Key, value interface{}, result interface{}) (err error) {
	return defsys.GetSet(key, value, result)
}
func MGet(keys ...core.Redis_Key) (result []string, err error) {
	return defsys.MGet(keys...)
}

/*List*/
func Lindex(key core.Redis_Key, value interface{}) (err error) {
	return defsys.Lindex(key, value)
}
func Linsert(key core.Redis_Key, isbefore bool, tager interface{}, value interface{}) (err error) {
	return defsys.Linsert(key, isbefore, tager, value)
}
func Llen(key core.Redis_Key) (result int, err error) {
	return defsys.Llen(key)
}
func LPop(key core.Redis_Key, value interface{}) (err error) {
	return defsys.LPop(key, value)
}
func LPush(key core.Redis_Key, values ...interface{}) (err error) {
	return defsys.LPush(key, values...)
}
func LPushX(key core.Redis_Key, values ...interface{}) (err error) {
	return defsys.LPushX(key, values...)
}
func LRange(key core.Redis_Key, start, end int, valuetype reflect.Type) (result []interface{}, err error) {
	return defsys.LRange(key, start, end, valuetype)
}
func LRem(key core.Redis_Key, count int, target interface{}) (err error) {
	return defsys.LRem(key, count, target)
}
func LSet(key core.Redis_Key, index int, value interface{}) (err error) {
	return defsys.LSet(key, index, value)
}
func Ltrim(key core.Redis_Key, start, stop int) (err error) {
	return defsys.Ltrim(key, start, stop)
}
func Rpop(key core.Redis_Key, value interface{}) (err error) {
	return defsys.Rpop(key, value)
}
func RPopLPush(oldkey core.Redis_Key, newkey core.Redis_Key, value interface{}) (err error) {
	return defsys.RPopLPush(oldkey, newkey, value)
}
func RPush(key core.Redis_Key, values ...interface{}) (err error) {
	return defsys.RPush(key, values...)
}
func RPushX(key core.Redis_Key, values ...interface{}) (err error) {
	return defsys.RPushX(key, values...)
}

/*Hash*/
func HDel(key core.Redis_Key, fields ...string) (err error) {
	return defsys.HDel(key, fields...)
}
func HExists(key core.Redis_Key, field string) (result bool, err error) {
	return defsys.HExists(key, field)
}
func HGet(key core.Redis_Key, field string, value interface{}) (err error) {
	return defsys.HGet(key, field, value)
}
func HGetAll(key core.Redis_Key, valuetype reflect.Type) (result []interface{}, err error) {
	return defsys.HGetAll(key, valuetype)
}
func HIncrBy(key core.Redis_Key, field string, value int) (err error) {
	return defsys.HIncrBy(key, field, value)
}
func HIncrByFloat(key core.Redis_Key, field string, value float32) (err error) {
	return defsys.HIncrByFloat(key, field, value)
}
func Hkeys(key core.Redis_Key) (result []string, err error) {
	return defsys.Hkeys(key)
}
func Hlen(key core.Redis_Key) (result int, err error) {
	return defsys.Hlen(key)
}
func HMGet(key core.Redis_Key, valuetype reflect.Type, fields ...string) (result []interface{}, err error) {
	return defsys.HMGet(key, valuetype, fields...)
}
func HMSet(key core.Redis_Key, value map[string]interface{}) (err error) {
	return defsys.HMSet(key, value)
}
func HSet(key core.Redis_Key, field string, value interface{}) (err error) {
	return defsys.HSet(key, field, value)
}
func HSetNX(key core.Redis_Key, field string, value interface{}) (err error) {
	return defsys.HSetNX(key, field, value)
}

/*Set*/
func SAdd(key core.Redis_Key, values ...interface{}) (err error) {
	return defsys.SAdd(key)
}
func Scard(key core.Redis_Key) (result int, err error) {
	return defsys.Scard(key)
}
func Sismember(key core.Redis_Key, value interface{}) (iskeep bool, err error) {
	return defsys.Sismember(key, value)
}
