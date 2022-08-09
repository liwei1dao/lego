package redis

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/liwei1dao/lego/sys/redis/pipe"
)

type (
	IRedis interface {
		Close() (err error)
		Context() context.Context
		Do(ctx context.Context, args ...interface{}) *redis.Cmd
		Lock(key string, outTime int) (result bool, err error)
		UnLock(key string) (err error)
		Pipeline() redis.Pipeliner
		Pipelined(ctx context.Context, fn func(pipe redis.Pipeliner) error) (err error)
		TxPipelined(ctx context.Context, fn func(pipe redis.Pipeliner) error) (err error)
		Watch(ctx context.Context, fn func(*redis.Tx) error, keys ...string) (err error)

		/*Key*/
		Delete(key string) (err error)
		ExistsKey(key string) (iskeep bool, err error)
		Expire(key string, expiration time.Duration) (err error)
		ExpireAt(key string, tm time.Time) (err error)
		PExpire(key string, expiration time.Duration) (err error)
		PExpireAt(key string, tm time.Time) (err error)
		Persist(key string) (err error)
		PTTL(key string) (leftexpire time.Duration, err error)
		TTL(key string) (leftexpire time.Duration, err error)
		Rename(oldkey string, newkey string) (err error)
		RenameNX(oldkey string, newkey string) (err error)
		Keys(pattern string) (keys []string, err error)
		Type(key string) (ty string, err error)
		/*String*/
		Set(key string, value interface{}, expiration time.Duration) (err error)
		SetNX(key string, value interface{}) (result int64, err error)
		MSet(keyvalues map[string]interface{}) (err error)
		MSetNX(keyvalues map[string]interface{}) (err error)
		Incr(key string) (err error)
		IncrBY(key string, value int) (err error)
		Incrbyfloat(key string, value float32) (err error)
		Decr(key string, value int) (err error)
		DecrBy(key string, value int) (err error)
		Append(key string, value interface{}) (err error)
		Get(key string, value interface{}) (err error)
		GetSet(key string, value interface{}, result interface{}) (err error)
		MGet(v interface{}, keys ...string) (err error)
		INCRBY(key string, amount int64) (result int64, err error)
		/*List*/
		Lindex(key string, value interface{}) (err error)
		Linsert(key string, isbefore bool, tager interface{}, value interface{}) (err error)
		Llen(key string) (result int, err error)
		LPop(key string, value interface{}) (err error)
		LPush(key string, values ...interface{}) (err error)
		LPushX(key string, values ...interface{}) (err error)
		LRange(key string, start, end int, v interface{}) (err error)
		LRem(key string, count int, target interface{}) (err error)
		LSet(key string, index int, value interface{}) (err error)
		Ltrim(key string, start, stop int) (err error)
		Rpop(key string, value interface{}) (err error)
		RPopLPush(oldkey string, newkey string, value interface{}) (err error)
		RPush(key string, values ...interface{}) (err error)
		RPushX(key string, values ...interface{}) (err error)
		/*Hash*/
		HDel(key string, fields ...string) (err error)
		HExists(key string, field string) (result bool, err error)
		HGet(key string, field string, value interface{}) (err error)
		HGetAll(key string, v interface{}) (err error)
		HGetAllToMapString(key string) (result map[string]string, err error)
		HIncrBy(key string, field string, value int) (err error)
		HIncrByFloat(key string, field string, value float32) (err error)
		Hkeys(key string) (result []string, err error)
		Hlen(key string) (result int, err error)
		HMGet(key string, v interface{}, fields ...string) (err error)
		HMSet(key string, v interface{}) (err error)
		HMSetForMap(key string, v map[string]string) (err error)
		HSet(key string, field string, value interface{}) (err error)
		HSetNX(key string, field string, value interface{}) (err error)
		/*Set*/
		SAdd(key string, values ...interface{}) (err error)
		SCard(key string) (result int64, err error)
		SDiff(v interface{}, keys ...string) (err error)
		SDiffStore(destination string, keys ...string) (result int64, err error)
		SInter(v interface{}, keys ...string) (err error)
		SInterStore(destination string, keys ...string) (result int64, err error)
		Sismember(key string, value interface{}) (iskeep bool, err error)
		SMembers(v interface{}, key string) (err error)
		SMove(source string, destination string, member interface{}) (result bool, err error)
		Spop(key string) (result string, err error)
		Srandmember(key string) (result string, err error)
		SRem(key string, members ...interface{}) (result int64, err error)
		SUnion(v interface{}, keys ...string) (err error)
		Sunionstore(destination string, keys ...string) (result int64, err error)
		Sscan(key string, _cursor uint64, match string, count int64) (keys []string, cursor uint64, err error)
		/*ZSet*/
		ZAdd(key string, members ...*redis.Z) (err error)
		ZCard(key string) (result int64, err error)
		ZCount(key string, min string, max string) (result int64, err error)
		ZIncrBy(key string, increment float64, member string) (result float64, err error)
		ZInterStore(destination string, store *redis.ZStore) (result int64, err error)
		ZLexCount(key string, min string, max string) (result int64, err error)
		ZRange(key string, start int64, stop int64, v interface{}) (err error)
		ZRangeByLex(key string, opt *redis.ZRangeBy, v interface{}) (err error)
		ZRangeByScore(key string, opt *redis.ZRangeBy, v interface{}) (err error)
		ZRank(key string, member string) (result int64, err error)
		ZRem(key string, members ...interface{}) (result int64, err error)
		ZRemRangeByLex(key string, min string, max string) (result int64, err error)
		ZRemRangeByRank(key string, start int64, stop int64) (result int64, err error)
		ZRemRangeByScore(key string, min string, max string) (result int64, err error)
		ZRevRange(key string, start int64, stop int64, v interface{}) (err error)
		ZRevRangeByScore(key string, opt *redis.ZRangeBy, v interface{}) (err error)
		ZRevRank(key string, member string) (result int64, err error)
		ZScore(key string, member string) (result float64, err error)
		ZUnionStore(dest string, store *redis.ZStore) (result int64, err error)
		ZScan(key string, _cursor uint64, match string, count int64) (keys []string, cursor uint64, err error)

		/*Lua Script*/
		NewScript(src string) *redis.StringCmd
		Eval(ctx context.Context, script string, keys []string, args ...interface{}) *redis.Cmd
		EvalSha(ctx context.Context, sha1 string, keys []string, args ...interface{}) *redis.Cmd
		ScriptExists(ctx context.Context, hashes ...string) *redis.BoolSliceCmd
		// ScriptLoad(ctx context.Context, script string) *redis.StringCmd
	}

	ISys interface {
		IRedis
		GetClient() IRedis
		RedisPipe(ctx context.Context) *pipe.RedisPipe
		/*Lock*/
		NewRedisMutex(key string, opt ...RMutexOption) (result *RedisMutex, err error)
	}
)

const (
	RedisNil    = redis.Nil //数据为空错误
	TxFailedErr = redis.TxFailedErr
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

func Close() (err error) {
	return defsys.Close()
}

func GetClient() IRedis {
	return defsys.GetClient()
}
func RedisPipe(ctx context.Context) *pipe.RedisPipe {
	return defsys.RedisPipe(ctx)
}
func Context() context.Context {
	return defsys.Context()
}

func Do(ctx context.Context, args ...interface{}) *redis.Cmd {
	return defsys.Do(ctx, args...)
}
func Pipeline() redis.Pipeliner {
	return defsys.Pipeline()
}
func Pipelined(ctx context.Context, fn func(pipe redis.Pipeliner) error) (err error) {
	return defsys.Pipelined(ctx, fn)
}
func TxPipelined(ctx context.Context, fn func(pipe redis.Pipeliner) error) (err error) {
	return defsys.TxPipelined(ctx, fn)
}
func Watch(ctx context.Context, fn func(*redis.Tx) error, keys ...string) (err error) {
	return defsys.Watch(ctx, fn)
}
func Delete(key string) (err error) {
	return defsys.Delete(key)
}
func ExistsKey(key string) (iskeep bool, err error) {
	return defsys.ExistsKey(key)
}
func Expire(key string, expiration time.Duration) (err error) {
	return defsys.Expire(key, expiration)
}
func ExpireAt(key string, tm time.Time) (err error) {
	return defsys.ExpireAt(key, tm)
}
func PExpire(key string, expiration time.Duration) (err error) {
	return defsys.PExpire(key, expiration)
}
func PExpireAt(key string, tm time.Time) (err error) {
	return defsys.PExpireAt(key, tm)
}
func Persist(key string) (err error) {
	return defsys.Persist(key)
}
func PTTL(key string) (leftexpire time.Duration, err error) {
	return defsys.PTTL(key)
}
func TTL(key string) (leftexpire time.Duration, err error) {
	return defsys.TTL(key)
}
func Rename(oldkey string, newkey string) (err error) {
	return defsys.Rename(oldkey, newkey)
}
func RenameNX(oldkey string, newkey string) (err error) {
	return defsys.RenameNX(oldkey, newkey)
}
func Keys(pattern string) (keys []string, err error) {
	return defsys.Keys(pattern)
}

///获取键类型
func Type(key string) (ty string, err error) {
	return defsys.Type(key)
}

/*String*/
func Set(key string, value interface{}, expiration time.Duration) (err error) {
	return defsys.Set(key, value, expiration)
}
func SetNX(key string, value interface{}) (result int64, err error) {
	return defsys.SetNX(key, value)
}
func MSet(keyvalues map[string]interface{}) (err error) {
	return defsys.MSet(keyvalues)
}
func MSetNX(keyvalues map[string]interface{}) (err error) {
	return defsys.MSetNX(keyvalues)
}
func Incr(key string) (err error) {
	return defsys.Incr(key)
}
func IncrBY(key string, value int) (err error) {
	return defsys.IncrBY(key, value)
}
func Incrbyfloat(key string, value float32) (err error) {
	return defsys.Incrbyfloat(key, value)
}
func Decr(key string, value int) (err error) {
	return defsys.Decr(key, value)
}
func DecrBy(key string, value int) (err error) {
	return defsys.DecrBy(key, value)
}
func Append(key string, value interface{}) (err error) {
	return defsys.Append(key, value)
}
func Get(key string, value interface{}) (err error) {
	return defsys.Get(key, value)
}
func GetSet(key string, value interface{}, result interface{}) (err error) {
	return defsys.GetSet(key, value, result)
}
func MGet(v interface{}, keys ...string) (err error) {
	return defsys.MGet(v, keys...)
}
func INCRBY(key string, amount int64) (result int64, err error) {
	return defsys.INCRBY(key, amount)
}

/*Lock*/
func NewRedisMutex(key string, opt ...RMutexOption) (result *RedisMutex, err error) {
	return defsys.NewRedisMutex(key, opt...)
}

func Lock(key string, outTime int) (result bool, err error) {
	return defsys.Lock(key, outTime)
}
func UnLock(key string) (err error) {
	return defsys.UnLock(key)
}

/*List*/
func Lindex(key string, value interface{}) (err error) {
	return defsys.Lindex(key, value)
}
func Linsert(key string, isbefore bool, tager interface{}, value interface{}) (err error) {
	return defsys.Linsert(key, isbefore, tager, value)
}
func Llen(key string) (result int, err error) {
	return defsys.Llen(key)
}
func LPop(key string, value interface{}) (err error) {
	return defsys.LPop(key, value)
}
func LPush(key string, values ...interface{}) (err error) {
	return defsys.LPush(key, values...)
}
func LPushX(key string, values ...interface{}) (err error) {
	return defsys.LPushX(key, values...)
}
func LRange(key string, start, end int, v interface{}) (err error) {
	return defsys.LRange(key, start, end, v)
}
func LRem(key string, count int, target interface{}) (err error) {
	return defsys.LRem(key, count, target)
}
func LSet(key string, index int, value interface{}) (err error) {
	return defsys.LSet(key, index, value)
}
func Ltrim(key string, start, stop int) (err error) {
	return defsys.Ltrim(key, start, stop)
}
func Rpop(key string, value interface{}) (err error) {
	return defsys.Rpop(key, value)
}
func RPopLPush(oldkey string, newkey string, value interface{}) (err error) {
	return defsys.RPopLPush(oldkey, newkey, value)
}
func RPush(key string, values ...interface{}) (err error) {
	return defsys.RPush(key, values...)
}
func RPushX(key string, values ...interface{}) (err error) {
	return defsys.RPushX(key, values...)
}

/*Hash*/
func HDel(key string, fields ...string) (err error) {
	return defsys.HDel(key, fields...)
}
func HExists(key string, field string) (result bool, err error) {
	return defsys.HExists(key, field)
}
func HGet(key string, field string, value interface{}) (err error) {
	return defsys.HGet(key, field, value)
}
func HGetAll(key string, v interface{}) (err error) {
	return defsys.HGetAll(key, v)
}
func HGetAllToMapString(key string) (result map[string]string, err error) {
	return defsys.HGetAllToMapString(key)
}
func HIncrBy(key string, field string, value int) (err error) {
	return defsys.HIncrBy(key, field, value)
}
func HIncrByFloat(key string, field string, value float32) (err error) {
	return defsys.HIncrByFloat(key, field, value)
}
func Hkeys(key string) (result []string, err error) {
	return defsys.Hkeys(key)
}
func Hlen(key string) (result int, err error) {
	return defsys.Hlen(key)
}
func HMGet(key string, v interface{}, fields ...string) (err error) {
	return defsys.HMGet(key, v, fields...)
}
func HMSet(key string, v interface{}) (err error) {
	return defsys.HMSet(key, v)
}
func HMSetForMap(key string, v map[string]string) (err error) {
	return defsys.HMSetForMap(key, v)
}
func HSet(key string, field string, value interface{}) (err error) {
	return defsys.HSet(key, field, value)
}
func HSetNX(key string, field string, value interface{}) (err error) {
	return defsys.HSetNX(key, field, value)
}

/*Set*/
func SAdd(key string, values ...interface{}) (err error) {
	return defsys.SAdd(key, values...)
}
func SCard(key string) (result int64, err error) {
	return defsys.SCard(key)
}
func SDiff(v interface{}, keys ...string) (err error) {
	return defsys.SDiff(v, keys...)
}
func SDiffStore(destination string, keys ...string) (result int64, err error) {
	return defsys.SDiffStore(destination, keys...)
}
func SInter(v interface{}, keys ...string) (err error) {
	return defsys.SInter(v, keys...)
}
func SInterStore(destination string, keys ...string) (result int64, err error) {
	return defsys.SInterStore(destination, keys...)
}
func Sismember(key string, value interface{}) (iskeep bool, err error) {
	return defsys.Sismember(key, value)
}
func SMembers(v interface{}, key string) (err error) {
	return defsys.SMembers(v, key)
}
func SMove(source string, destination string, member interface{}) (result bool, err error) {
	return defsys.SMove(source, destination, member)
}
func Spop(key string) (result string, err error) {
	return defsys.Spop(key)
}
func Srandmember(key string) (result string, err error) {
	return defsys.Srandmember(key)
}
func SRem(key string, members ...interface{}) (result int64, err error) {
	return defsys.SRem(key, members...)
}
func SUnion(v interface{}, keys ...string) (err error) {
	return defsys.SUnion(v, keys...)
}
func Sunionstore(destination string, keys ...string) (result int64, err error) {
	return defsys.Sunionstore(destination, keys...)
}
func Sscan(key string, _cursor uint64, match string, count int64) (keys []string, cursor uint64, err error) {
	return defsys.Sscan(key, _cursor, match, count)
}

/*ZSet*/
func ZAdd(key string, members ...*redis.Z) (err error) {
	return defsys.ZAdd(key, members...)
}
func ZCard(key string) (result int64, err error) {
	return defsys.ZCard(key)
}
func ZCount(key string, min string, max string) (result int64, err error) {
	return defsys.ZCount(key, min, max)
}
func ZIncrBy(key string, increment float64, member string) (result float64, err error) {
	return defsys.ZIncrBy(key, increment, member)
}
func ZInterStore(destination string, store *redis.ZStore) (result int64, err error) {
	return defsys.ZInterStore(destination, store)
}
func ZLexCount(key string, min string, max string) (result int64, err error) {
	return defsys.ZLexCount(key, min, max)
}
func ZRange(key string, start int64, stop int64, v interface{}) (err error) {
	return defsys.ZRange(key, start, stop, v)
}
func ZRangeByLex(key string, opt *redis.ZRangeBy, v interface{}) (err error) {
	return defsys.ZRangeByLex(key, opt, v)
}
func ZRangeByScore(key string, opt *redis.ZRangeBy, v interface{}) (err error) {
	return defsys.ZRangeByScore(key, opt, v)
}
func ZRank(key string, member string) (result int64, err error) {
	return defsys.ZRank(key, member)
}
func ZRem(key string, members ...interface{}) (result int64, err error) {
	return defsys.ZRem(key, members...)
}
func ZRemRangeByLex(key string, min string, max string) (result int64, err error) {
	return defsys.ZRemRangeByLex(key, min, max)
}
func ZRemRangeByRank(key string, start int64, stop int64) (result int64, err error) {
	return defsys.ZRemRangeByRank(key, start, stop)
}
func ZRemRangeByScore(key string, min string, max string) (result int64, err error) {
	return defsys.ZRemRangeByScore(key, min, max)
}
func ZRevRange(key string, start int64, stop int64, v interface{}) (err error) {
	return defsys.ZRevRange(key, start, stop, v)
}
func ZRevRangeByScore(key string, opt *redis.ZRangeBy, v interface{}) (err error) {
	return defsys.ZRevRangeByScore(key, opt, v)
}
func ZRevRank(key string, member string) (result int64, err error) {
	return defsys.ZRevRank(key, member)
}
func ZScore(key string, member string) (result float64, err error) {
	return defsys.ZScore(key, member)
}
func ZUnionStore(dest string, store *redis.ZStore) (result int64, err error) {
	return defsys.ZUnionStore(dest, store)
}
func ZScan(key string, _cursor uint64, match string, count int64) (keys []string, cursor uint64, err error) {
	return defsys.ZScan(key, _cursor, match, count)
}

/*Lua Script*/
func NewScript(src string) *redis.StringCmd {
	return defsys.NewScript(src)
}
func Eval(ctx context.Context, script string, keys []string, args ...interface{}) *redis.Cmd {
	return defsys.Eval(ctx, script, keys, args...)
}
func EvalSha(ctx context.Context, sha1 string, keys []string, args ...interface{}) *redis.Cmd {
	return defsys.EvalSha(ctx, sha1, keys, args...)
}
func ScriptExists(ctx context.Context, hashes ...string) *redis.BoolSliceCmd {
	return defsys.ScriptExists(ctx, hashes...)
}

// func ScriptLoad(ctx context.Context, script string) *redis.StringCmd {
// 	return defsys.ScriptLoad(ctx, script)
// }
