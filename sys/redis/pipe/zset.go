package pipe

import (
	"github.com/go-redis/redis/v8"
)

/*
Redis ZAdd 向有序集合添加一个或多个成员，或者更新已存在成员的分数
*/
func (this *RedisPipe) ZAdd(key string, members ...*redis.Z) (err error) {
	this.client.ZAdd(this.ctx, key, members...)
	return
}

/*
Redis Zcard 用于计算集合中元素的数量。
*/
func (this *RedisPipe) ZCard(key string) (result int64, err error) {
	result, err = this.client.ZCard(this.ctx, key).Result()
	return
}

/*
Redis ZCount 用于计算集合中指定的范围内的数量
*/
func (this *RedisPipe) ZCount(key string, min string, max string) (result int64, err error) {
	result, err = this.client.ZCount(this.ctx, key, min, max).Result()
	return
}

/*
Redis ZIncrBy 有序集合中对指定成员的分数加上增量 increment
*/
func (this *RedisPipe) ZIncrBy(key string, increment float64, member string) (result float64, err error) {
	result, err = this.client.ZIncrBy(this.ctx, key, increment, member).Result()
	return
}

/*
Redis ZInterStore 计算给定的一个或多个有序集的交集并将结果集存储在新的有序集合 destination 中
*/
func (this *RedisPipe) ZInterStore(destination string, store *redis.ZStore) (result int64, err error) {
	result, err = this.client.ZInterStore(this.ctx, destination, store).Result()
	return
}

/*
Redis ZLexCount 在有序集合中计算指定字典区间内成员数量
*/
func (this *RedisPipe) ZLexCount(key string, min string, max string) (result int64, err error) {
	result, err = this.client.ZLexCount(this.ctx, key, min, max).Result()
	return
}

/*
Redis ZRange 通过索引区间返回有序集合指定区间内的成员
*/
func (this *RedisPipe) ZRange(key string, start int64, stop int64, v interface{}) (err error) {
	var _result []string
	cmd := this.client.ZRange(this.ctx, key, start, stop)
	if _result, err = cmd.Result(); err == nil {
		err = this.codec.UnmarshalSlice(_result, v)
	}
	return
}

/*
Redis ZRangeByLex 通过字典区间返回有序集合的成员
*/
func (this *RedisPipe) ZRangeByLex(key string, opt *redis.ZRangeBy, v interface{}) (err error) {
	var _result []string
	cmd := this.client.ZRangeByLex(this.ctx, key, opt)
	if _result, err = cmd.Result(); err == nil {
		err = this.codec.UnmarshalSlice(_result, v)
	}
	return
}

/*
Redis ZRangeByScore 通过分数返回有序集合指定区间内的成员
*/
func (this *RedisPipe) ZRangeByScore(key string, opt *redis.ZRangeBy, v interface{}) (err error) {
	var _result []string
	cmd := this.client.ZRangeByScore(this.ctx, key, opt)
	if _result, err = cmd.Result(); err == nil {
		err = this.codec.UnmarshalSlice(_result, v)
	}
	return
}

/*
Redis ZRank 返回有序集合中指定成员的索引
*/
func (this *RedisPipe) ZRank(key string, member string) (result int64, err error) {
	result, err = this.client.ZRank(this.ctx, key, member).Result()
	return
}

/*
Redis ZRem 移除有序集合中的一个或多个成员
*/
func (this *RedisPipe) ZRem(key string, members ...interface{}) (result int64, err error) {
	result, err = this.client.ZRem(this.ctx, key, members...).Result()
	return
}

/*
Redis ZRemRangeByLex 移除有序集合中给定的字典区间的所有成员
*/
func (this *RedisPipe) ZRemRangeByLex(key string, min string, max string) (result int64, err error) {
	result, err = this.client.ZRemRangeByLex(this.ctx, key, min, max).Result()
	return
}

/*
Redis ZRemRangeByRank 移除有序集合中给定的排名区间的所有成员
*/
func (this *RedisPipe) ZRemRangeByRank(key string, start int64, stop int64) (result int64, err error) {
	result, err = this.client.ZRemRangeByRank(this.ctx, key, start, stop).Result()
	return
}

/*
Redis ZRemRangeByScore 移除有序集合中给定的分数区间的所有成员
*/
func (this *RedisPipe) ZRemRangeByScore(key string, min string, max string) (result int64, err error) {
	result, err = this.client.ZRemRangeByScore(this.ctx, key, min, max).Result()
	return
}

/*
Redis ZRevRange 返回有序集中指定区间内的成员，通过索引，分数从高到低 ZREVRANGE
*/
func (this *RedisPipe) ZRevRange(key string, start int64, stop int64, v interface{}) (err error) {
	var _result []string
	cmd := this.client.ZRevRange(this.ctx, key, start, stop)
	if _result, err = cmd.Result(); err == nil {
		err = this.codec.UnmarshalSlice(_result, v)
	}
	return
}

/*
Redis ZRevRangeByScore 返回有序集中指定分数区间内的成员，分数从高到低排序
*/
func (this *RedisPipe) ZRevRangeByScore(key string, opt *redis.ZRangeBy, v interface{}) (err error) {
	var _result []string
	cmd := this.client.ZRevRangeByScore(this.ctx, key, opt)
	if _result, err = cmd.Result(); err == nil {
		err = this.codec.UnmarshalSlice(_result, v)
	}
	return
}

/*
Redis ZRevRank 返回有序集中指定分数区间内的成员，分数从高到低排序
*/
func (this *RedisPipe) ZRevRank(key string, member string) (result int64, err error) {
	result, err = this.client.ZRevRank(this.ctx, key, member).Result()
	return
}

/*
Redis ZScore 返回有序集中指定分数区间内的成员，分数从高到低排序
*/
func (this *RedisPipe) ZScore(key string, member string) (result float64, err error) {
	result, err = this.client.ZScore(this.ctx, key, member).Result()
	return
}

/*
Redis ZScore 返回有序集中指定分数区间内的成员，分数从高到低排序 ZUNIONSTORE
*/
func (this *RedisPipe) ZUnionStore(dest string, store *redis.ZStore) (result int64, err error) {
	result, err = this.client.ZUnionStore(this.ctx, dest, store).Result()
	return
}

/*
Redis ZScan 迭代有序集合中的元素（包括元素成员和元素分值）
*/
func (this *RedisPipe) ZScan(key string, _cursor uint64, match string, count int64) (keys []string, cursor uint64, err error) {
	keys, cursor, err = this.client.ZScan(this.ctx, key, _cursor, match, count).Result()
	return
}
