package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/liwei1dao/lego/sys/redis/cluster"
	"github.com/liwei1dao/lego/sys/redis/pipe"
	"github.com/liwei1dao/lego/sys/redis/single"
	"github.com/liwei1dao/lego/utils/codec/json"

	"github.com/go-redis/redis/v8"
)

func newSys(options *Options) (sys *Redis, err error) {
	sys = &Redis{options: options}
	err = sys.init()
	return
}

type Redis struct {
	options *Options
	client  IRedis
}

func (this *Redis) init() (err error) {
	if this.options.RedisType == Redis_Single {
		this.client, err = single.NewSys(
			this.options.Redis_Single_Addr,
			this.options.Redis_Single_Password,
			this.options.Redis_Single_DB,
			this.options.TimeOut,
			this,
		)
	} else if this.options.RedisType == Redis_Cluster {
		this.client, err = cluster.NewSys(
			this.options.Redis_Cluster_Addr,
			this.options.Redis_Cluster_Password,
			this.options.TimeOut,
			this,
		)
	} else {
		err = fmt.Errorf("init Redis err:RedisType - %d", this.options.RedisType)
	}
	return
}
func (this *Redis) Close() (err error) {
	return this.client.Close()
}
func (this *Redis) GetClient() IRedis {
	return this.client
}

func (this *Redis) Context() context.Context {
	return this.client.Context()
}

func (this *Redis) Do(ctx context.Context, args ...interface{}) *redis.Cmd {
	return this.client.Do(ctx, args...)
}
func (this *Redis) RedisPipe(ctx context.Context) *pipe.RedisPipe {
	return pipe.NewPipe(ctx, this.client.Pipeline(), this)
}
func (this *Redis) Pipeline() redis.Pipeliner {
	return this.client.Pipeline()
}
func (this *Redis) Pipelined(ctx context.Context, fn func(pipe redis.Pipeliner) error) (err error) {
	return this.client.Pipelined(ctx, fn)
}
func (this *Redis) TxPipelined(ctx context.Context, fn func(pipe redis.Pipeliner) error) (err error) {
	return this.client.TxPipelined(ctx, fn)
}
func (this *Redis) Watch(ctx context.Context, fn func(*redis.Tx) error, keys ...string) (err error) {
	return this.client.Watch(ctx, fn)
}
func (this *Redis) Lock(key string, outTime int) (result bool, err error) {
	return this.client.Lock(key, outTime)
}
func (this *Redis) UnLock(key string) (err error) {
	return this.client.UnLock(key)
}
func (this *Redis) Delete(key string) (err error) {
	return this.client.Delete(key)
}
func (this *Redis) ExistsKey(key string) (iskeep bool, err error) {
	return this.client.ExistsKey(key)
}
func (this *Redis) Expire(key string, expire time.Duration) (err error) {
	return this.client.Expire(key, expire)
}
func (this *Redis) ExpireAt(key string, tm time.Time) (err error) {
	return this.client.ExpireAt(key, tm)
}
func (this *Redis) PExpire(key string, expire time.Duration) (err error) {
	return this.client.PExpire(key, expire)
}
func (this *Redis) PExpireAt(key string, tm time.Time) (err error) {
	return this.client.PExpireAt(key, tm)
}
func (this *Redis) Persist(key string) (err error) {
	return this.client.Persist(key)
}
func (this *Redis) PTTL(key string) (leftexpire time.Duration, err error) {
	return this.client.PTTL(key)
}
func (this *Redis) TTL(key string) (leftexpire time.Duration, err error) {
	return this.client.TTL(key)
}
func (this *Redis) Rename(oldkey string, newkey string) (err error) {
	return this.client.Rename(oldkey, newkey)
}
func (this *Redis) RenameNX(oldkey string, newkey string) (err error) {
	return this.client.RenameNX(oldkey, newkey)
}
func (this *Redis) Keys(pattern string) (keys []string, err error) {
	return this.client.Keys(pattern)
}
func (this *Redis) Type(key string) (ty string, err error) {
	return this.client.Type(key)
}

/*String*/
func (this *Redis) Set(key string, value interface{}, expiration time.Duration) (err error) {
	return this.client.Set(key, value, expiration)
}
func (this *Redis) SetNX(key string, value interface{}) (result int64, err error) {
	return this.client.SetNX(key, value)
}
func (this *Redis) MSet(keyvalues map[string]interface{}) (err error) {
	return this.client.MSet(keyvalues)
}
func (this *Redis) MSetNX(keyvalues map[string]interface{}) (err error) {
	return this.client.MSetNX(keyvalues)
}
func (this *Redis) Incr(key string) (err error) {
	return this.client.Incr(key)
}
func (this *Redis) IncrBY(key string, value int) (err error) {
	return this.client.IncrBY(key, value)
}
func (this *Redis) Incrbyfloat(key string, value float32) (err error) {
	return this.client.Incrbyfloat(key, value)
}
func (this *Redis) Decr(key string, value int) (err error) {
	return this.client.Decr(key, value)
}
func (this *Redis) DecrBy(key string, value int) (err error) {
	return this.client.DecrBy(key, value)
}
func (this *Redis) Append(key string, value interface{}) (err error) {
	return this.client.Append(key, value)
}
func (this *Redis) Get(key string, value interface{}) (err error) {
	return this.client.Get(key, value)
}
func (this *Redis) GetSet(key string, value interface{}, result interface{}) (err error) {
	return this.client.GetSet(key, value, result)
}
func (this *Redis) MGet(v interface{}, keys ...string) (err error) {
	return this.client.MGet(v, keys...)
}
func (this *Redis) INCRBY(key string, amount int64) (result int64, err error) {
	return this.client.INCRBY(key, amount)
}

/*List*/
func (this *Redis) Lindex(key string, value interface{}) (err error) {
	return this.client.Lindex(key, value)
}
func (this *Redis) Linsert(key string, isbefore bool, tager interface{}, value interface{}) (err error) {
	return this.client.Linsert(key, isbefore, tager, value)
}
func (this *Redis) Llen(key string) (result int, err error) {
	return this.client.Llen(key)
}
func (this *Redis) LPop(key string, value interface{}) (err error) {
	return this.client.LPop(key, value)
}
func (this *Redis) LPush(key string, values ...interface{}) (err error) {
	return this.client.LPush(key, values...)
}
func (this *Redis) LPushX(key string, values ...interface{}) (err error) {
	return this.client.LPushX(key, values...)
}
func (this *Redis) LRange(key string, start, end int, v interface{}) (err error) {
	return this.client.LRange(key, start, end, v)
}
func (this *Redis) LRem(key string, count int, target interface{}) (err error) {
	return this.client.LRem(key, count, target)
}
func (this *Redis) LSet(key string, index int, value interface{}) (err error) {
	return this.client.LSet(key, index, value)
}
func (this *Redis) Ltrim(key string, start, stop int) (err error) {
	return this.client.Ltrim(key, start, stop)
}
func (this *Redis) Rpop(key string, value interface{}) (err error) {
	return this.client.Rpop(key, value)
}
func (this *Redis) RPopLPush(oldkey string, newkey string, value interface{}) (err error) {
	return this.client.RPopLPush(oldkey, newkey, value)
}
func (this *Redis) RPush(key string, values ...interface{}) (err error) {
	return this.client.RPush(key, values...)
}
func (this *Redis) RPushX(key string, values ...interface{}) (err error) {
	return this.client.RPushX(key, values...)
}

/*Hash*/
func (this *Redis) HDel(key string, fields ...string) (err error) {
	return this.client.HDel(key, fields...)
}
func (this *Redis) HExists(key string, field string) (result bool, err error) {
	return this.client.HExists(key, field)
}
func (this *Redis) HGet(key string, field string, value interface{}) (err error) {
	return this.client.HGet(key, field, value)
}
func (this *Redis) HGetAll(key string, v interface{}) (err error) {
	return this.client.HGetAll(key, v)
}
func (this *Redis) HGetAllToMapString(key string) (result map[string]string, err error) {
	return this.client.HGetAllToMapString(key)
}
func (this *Redis) HIncrBy(key string, field string, value int) (err error) {
	return this.client.HIncrBy(key, field, value)
}
func (this *Redis) HIncrByFloat(key string, field string, value float32) (err error) {
	return this.client.HIncrByFloat(key, field, value)
}
func (this *Redis) Hkeys(key string) (result []string, err error) {
	return this.client.Hkeys(key)
}
func (this *Redis) Hlen(key string) (result int, err error) {
	return this.client.Hlen(key)
}
func (this *Redis) HMGet(key string, v interface{}, fields ...string) (err error) {
	return this.client.HMGet(key, v, fields...)
}
func (this *Redis) HMSet(key string, v interface{}) (err error) {
	return this.client.HMSet(key, v)
}

func (this *Redis) HMSetForMap(key string, v map[string]string) (err error) {
	return this.client.HMSetForMap(key, v)
}

func (this *Redis) HSet(key string, field string, value interface{}) (err error) {
	return this.client.HSet(key, field, value)
}
func (this *Redis) HSetNX(key string, field string, value interface{}) (err error) {
	return this.client.HSetNX(key, field, value)
}

/*Set*/
func (this *Redis) SAdd(key string, values ...interface{}) (err error) {
	return this.client.SAdd(key, values...)
}
func (this *Redis) SCard(key string) (result int64, err error) {
	return this.client.SCard(key)
}
func (this *Redis) SDiff(v interface{}, keys ...string) (err error) {
	return this.client.SDiff(v, keys...)
}
func (this *Redis) SDiffStore(destination string, keys ...string) (result int64, err error) {
	return this.client.SDiffStore(destination, keys...)
}
func (this *Redis) SInter(v interface{}, keys ...string) (err error) {
	return this.client.SInter(v, keys...)
}
func (this *Redis) SInterStore(destination string, keys ...string) (result int64, err error) {
	return this.client.SInterStore(destination, keys...)
}
func (this *Redis) Sismember(key string, value interface{}) (iskeep bool, err error) {
	return this.client.Sismember(key, value)
}
func (this *Redis) SMembers(v interface{}, key string) (err error) {
	return this.client.SMembers(v, key)
}
func (this *Redis) SMove(source string, destination string, member interface{}) (result bool, err error) {
	return this.client.SMove(source, destination, member)
}
func (this *Redis) Spop(key string) (result string, err error) {
	return this.client.Spop(key)
}
func (this *Redis) Srandmember(key string) (result string, err error) {
	return this.client.Srandmember(key)
}
func (this *Redis) SRem(key string, members ...interface{}) (result int64, err error) {
	return this.client.SRem(key, members...)
}
func (this *Redis) SUnion(v interface{}, keys ...string) (err error) {
	return this.client.SUnion(v, keys...)
}
func (this *Redis) Sunionstore(destination string, keys ...string) (result int64, err error) {
	return this.client.Sunionstore(destination, keys...)
}
func (this *Redis) Sscan(key string, _cursor uint64, match string, count int64) (keys []string, cursor uint64, err error) {
	return this.client.Sscan(key, _cursor, match, count)
}

/*ZSet*/
func (this *Redis) ZAdd(key string, members ...*redis.Z) (err error) {
	return this.client.ZAdd(key, members...)
}
func (this *Redis) ZCard(key string) (result int64, err error) {
	return this.client.ZCard(key)
}
func (this *Redis) ZCount(key string, min string, max string) (result int64, err error) {
	return this.client.ZCount(key, min, max)
}
func (this *Redis) ZIncrBy(key string, increment float64, member string) (result float64, err error) {
	return this.client.ZIncrBy(key, increment, member)
}
func (this *Redis) ZInterStore(destination string, store *redis.ZStore) (result int64, err error) {
	return this.client.ZInterStore(destination, store)
}
func (this *Redis) ZLexCount(key string, min string, max string) (result int64, err error) {
	return this.client.ZLexCount(key, min, max)
}
func (this *Redis) ZRange(key string, start int64, stop int64, v interface{}) (err error) {
	return this.client.ZRange(key, start, stop, v)
}
func (this *Redis) ZRangeByLex(key string, opt *redis.ZRangeBy, v interface{}) (err error) {
	return this.client.ZRangeByLex(key, opt, v)
}
func (this *Redis) ZRangeByScore(key string, opt *redis.ZRangeBy, v interface{}) (err error) {
	return this.client.ZRangeByScore(key, opt, v)
}
func (this *Redis) ZRank(key string, member string) (result int64, err error) {
	return this.client.ZRank(key, member)
}
func (this *Redis) ZRem(key string, members ...interface{}) (result int64, err error) {
	return this.client.ZRem(key, members...)
}
func (this *Redis) ZRemRangeByLex(key string, min string, max string) (result int64, err error) {
	return this.client.ZRemRangeByLex(key, min, max)
}
func (this *Redis) ZRemRangeByRank(key string, start int64, stop int64) (result int64, err error) {
	return this.client.ZRemRangeByRank(key, start, stop)
}
func (this *Redis) ZRemRangeByScore(key string, min string, max string) (result int64, err error) {
	return this.client.ZRemRangeByScore(key, min, max)
}
func (this *Redis) ZRevRange(key string, start int64, stop int64, v interface{}) (err error) {
	return this.client.ZRevRange(key, start, stop, v)
}
func (this *Redis) ZRevRangeByScore(key string, opt *redis.ZRangeBy, v interface{}) (err error) {
	return this.client.ZRevRangeByScore(key, opt, v)
}
func (this *Redis) ZRevRank(key string, member string) (result int64, err error) {
	return this.client.ZRevRank(key, member)
}
func (this *Redis) ZScore(key string, member string) (result float64, err error) {
	return this.client.ZScore(key, member)
}
func (this *Redis) ZUnionStore(dest string, store *redis.ZStore) (result int64, err error) {
	return this.client.ZUnionStore(dest, store)
}
func (this *Redis) ZScan(key string, _cursor uint64, match string, count int64) (keys []string, cursor uint64, err error) {
	return this.client.ZScan(key, _cursor, match, count)
}

//lua Script
func (this *Redis) NewScript(src string) *redis.StringCmd {
	return this.client.NewScript(src)
}
func (this *Redis) Eval(ctx context.Context, script string, keys []string, args ...interface{}) *redis.Cmd {
	return this.client.Eval(ctx, script, keys, args...)
}
func (this *Redis) EvalSha(ctx context.Context, sha1 string, keys []string, args ...interface{}) *redis.Cmd {
	return this.client.EvalSha(ctx, sha1, keys, args...)
}
func (this *Redis) ScriptExists(ctx context.Context, hashes ...string) *redis.BoolSliceCmd {
	return this.client.ScriptExists(ctx, hashes...)
}

//Codec---------------------------------------------------------------------------------------------------------------------------------------
func (this *Redis) Marshal(v interface{}) ([]byte, error) {
	if this.options.Codec != nil {
		return this.options.Codec.Marshal(v)
	} else {
		return json.Marshal(v)
	}
}
func (this *Redis) Unmarshal(data []byte, v interface{}) error {
	if this.options.Codec != nil {
		return this.options.Codec.Unmarshal(data, v)
	} else {
		return json.Unmarshal(data, v)
	}
}
func (this *Redis) MarshalMap(val interface{}) (ret map[string]string, err error) {
	if this.options.Codec != nil {
		return this.options.Codec.MarshalMap(val)
	} else {
		return json.MarshalMap(val)
	}
}
func (this *Redis) UnmarshalMap(data map[string]string, val interface{}) (err error) {
	if this.options.Codec != nil {
		return this.options.Codec.UnmarshalMap(data, val)
	} else {
		return json.UnmarshalMap(data, val)
	}
}
func (this *Redis) MarshalSlice(val interface{}) (ret []string, err error) {
	if this.options.Codec != nil {
		return this.options.Codec.MarshalSlice(val)
	} else {
		return json.MarshalSlice(val)
	}
}
func (this *Redis) UnmarshalSlice(data []string, val interface{}) (err error) {
	if this.options.Codec != nil {
		return this.options.Codec.UnmarshalSlice(data, val)
	} else {
		return json.UnmarshalSlice(data, val)
	}
}
