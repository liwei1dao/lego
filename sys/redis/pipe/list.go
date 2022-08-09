package pipe

import (
	"github.com/go-redis/redis/v8"
)

/*
Redis Lindex 命令用于通过索引获取列表中的元素。你也可以使用负数下标，以 -1 表示列表的最后一个元素， -2 表示列表的倒数第二个元素，以此类推
*/
func (this *RedisPipe) Lindex(key string, v interface{}) (err error) {
	cmd := redis.NewStringCmd(this.ctx, "LINDEX", key)
	this.client.Process(this.ctx, cmd)
	var _result []byte
	if _result, err = cmd.Bytes(); err == nil {
		if len(_result) == 0 {
			err = redis.Nil
			return
		}
		err = this.codec.Unmarshal(_result, v)
	}
	return
}

/*
Redis Linsert 命令用于在列表的元素前或者后插入元素。当指定元素不存在于列表中时，不执行任何操作。
当列表不存在时，被视为空列表，不执行任何操作。
如果 key 不是列表类型，返回一个错误
*/
func (this *RedisPipe) Linsert(key string, isbefore bool, tager interface{}, value interface{}) (err error) {
	var (
		tagervalue  []byte
		resultvalue []byte
	)
	if tagervalue, err = this.codec.Marshal(tager); err != nil {
		return
	}
	if resultvalue, err = this.codec.Marshal(value); err != nil {
		return
	}
	if isbefore {
		err = this.client.Do(this.ctx, "LINSERT", key, "BEFORE", tagervalue, resultvalue).Err()
	} else {
		err = this.client.Do(this.ctx, "LINSERT", key, "AFTER", tagervalue, resultvalue).Err()
	}
	return
}

/*
Redis Llen 命令用于返回列表的长度。 如果列表 key 不存在，则 key 被解释为一个空列表，返回 0 。 如果 key 不是列表类型，返回一个错误
*/
func (this *RedisPipe) Llen(key string) (result int, err error) {
	result, err = this.client.Do(this.ctx, "LLEN", key).Int()
	return
}

/*
Redis Lpop 命令用于移除并返回列表的第一个元素
*/
func (this *RedisPipe) LPop(key string, v interface{}) (err error) {
	cmd := redis.NewStringCmd(this.ctx, "LPOP", key)
	this.client.Process(this.ctx, cmd)
	var _result []byte
	if _result, err = cmd.Bytes(); err == nil {
		err = this.codec.Unmarshal(_result, v)
	}
	return
}

/*
Redis Lpush 命令将一个或多个值插入到列表头部。 如果 key 不存在，一个空列表会被创建并执行 LPUSH 操作。 当 key 存在但不是列表类型时，返回一个错误
*/
func (this *RedisPipe) LPush(key string, values ...interface{}) (err error) {
	agrs := make([]interface{}, 0)
	agrs = append(agrs, "LPUSH")
	for _, v := range values {
		result, _ := this.codec.Marshal(v)
		agrs = append(agrs, result)
	}
	err = this.client.Do(this.ctx, agrs...).Err()
	return
}

/*
Redis Lpushx 将一个值插入到已存在的列表头部，列表不存在时操作无效
*/
func (this *RedisPipe) LPushX(key string, values ...interface{}) (err error) {
	agrs := make([]interface{}, 0)
	agrs = append(agrs, "LPUSHX")
	for _, v := range values {
		result, _ := this.codec.Marshal(v)
		agrs = append(agrs, result)
	}
	err = this.client.Do(this.ctx, agrs...).Err()
	return
}

/*
Redis Lrange 返回列表中指定区间内的元素，区间以偏移量 START 和 END 指定。 其中 0 表示列表的第一个元素， 1 表示列表的第二个元素，
以此类推。 你也可以使用负数下标，以 -1 表示列表的最后一个元素， -2 表示列表的倒数第二个元素，以此类推
*/
func (this *RedisPipe) LRange(key string, start, end int, v interface{}) (err error) {
	var _result []string
	cmd := redis.NewStringSliceCmd(this.ctx, "LRANGE", key, start, end)
	this.client.Process(this.ctx, cmd)
	if _result, err = cmd.Result(); err == nil {
		err = this.codec.UnmarshalSlice(_result, v)
	}
	return
}

/*
Redis Lrem 根据参数 COUNT 的值，移除列表中与参数 VALUE 相等的元素。
COUNT 的值可以是以下几种：
count > 0 : 从表头开始向表尾搜索，移除与 VALUE 相等的元素，数量为 COUNT 。
count < 0 : 从表尾开始向表头搜索，移除与 VALUE 相等的元素，数量为 COUNT 的绝对值。
count = 0 : 移除表中所有与 VALUE 相等的值
*/
func (this *RedisPipe) LRem(key string, count int, target interface{}) (err error) {
	var resultvalue []byte
	if resultvalue, err = this.codec.Marshal(target); err != nil {
		return
	}
	err = this.client.Do(this.ctx, "LREM", key, count, resultvalue).Err()
	return
}

/*
Redis Lset 通过索引来设置元素的值。
当索引参数超出范围，或对一个空列表进行 LSET 时，返回一个错误
*/
func (this *RedisPipe) LSet(key string, index int, value interface{}) (err error) {
	var resultvalue []byte
	if resultvalue, err = this.codec.Marshal(value); err == nil {
		return
	}
	err = this.client.Do(this.ctx, "LSET", key, index, resultvalue).Err()
	return
}

/*
Redis Ltrim 对一个列表进行修剪(trim)，就是说，让列表只保留指定区间内的元素，不在指定区间之内的元素都将被删除。
下标 0 表示列表的第一个元素，以 1 表示列表的第二个元素，以此类推。 你也可以使用负数下标，
以 -1 表示列表的最后一个元素， -2 表示列表的倒数第二个元素，以此类推
*/
func (this *RedisPipe) Ltrim(key string, start, stop int) (err error) {
	err = this.client.Do(this.ctx, "LTRIM", key, start, stop).Err()
	return
}

/*
Redis Rpop 命令用于移除列表的最后一个元素，返回值为移除的元素
*/
func (this *RedisPipe) Rpop(key string, v interface{}) (err error) {
	cmd := redis.NewStringCmd(this.ctx, "RPOP", key)
	this.client.Process(this.ctx, cmd)
	var _result []byte
	if _result, err = cmd.Bytes(); err == nil {
		err = this.codec.Unmarshal(_result, v)
	}
	return
}

/*
Redis Rpoplpush 命令用于移除列表的最后一个元素，并将该元素添加到另一个列表并返回
*/
func (this *RedisPipe) RPopLPush(oldkey string, newkey string, v interface{}) (err error) {
	cmd := redis.NewStringCmd(this.ctx, "RPOPLPUSH", oldkey, newkey)
	this.client.Process(this.ctx, cmd)
	var _result []byte
	if _result, err = cmd.Bytes(); err == nil {
		err = this.codec.Unmarshal(_result, v)
	}
	return
}

/*
Redis Rpush 命令用于将一个或多个值插入到列表的尾部(最右边)。
如果列表不存在，一个空列表会被创建并执行 RPUSH 操作。 当列表存在但不是列表类型时，返回一个错误。
注意：在 Redis 2.4 版本以前的 RPUSH 命令，都只接受单个 value 值
*/
func (this *RedisPipe) RPush(key string, values ...interface{}) (err error) {
	agrs := make([]interface{}, 0)
	agrs = append(agrs, "RPUSH")
	for _, v := range values {
		result, _ := this.codec.Marshal(v)
		agrs = append(agrs, result)
	}
	err = this.client.Do(this.ctx, agrs...).Err()
	return
}

/*
Redis Rpushx 命令用于将一个值插入到已存在的列表尾部(最右边)。如果列表不存在，操作无效
*/
func (this *RedisPipe) RPushX(key string, values ...interface{}) (err error) {
	agrs := make([]interface{}, 0)
	agrs = append(agrs, "RPUSHX")
	for _, v := range values {
		result, _ := this.codec.Marshal(v)
		agrs = append(agrs, result)
	}
	err = this.client.Do(this.ctx, agrs...).Err()
	return
}
