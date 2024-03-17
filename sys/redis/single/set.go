package single

/*
Redis Sadd 命令将一个或多个成员元素加入到集合中，已经存在于集合的成员元素将被忽略。
假如集合 key 不存在，则创建一个只包含添加的元素作成员的集合。
当集合 key 不是集合类型时，返回一个错误。
*/
func (this *Redis) SAdd(key string, values ...interface{}) (err error) {
	agrs := make([]interface{}, 0)
	agrs = append(agrs, "SADD")
	agrs = append(agrs, key)
	for _, v := range values {
		result, _ := this.codec.Marshal(v)
		agrs = append(agrs, result)
	}
	err = this.client.Do(this.client.Context(), agrs...).Err()
	return
}

/*
Redis Scard 命令返回集合中元素的数量
*/
func (this *Redis) SCard(key string) (result int64, err error) {
	result, err = this.client.SCard(this.client.Context(), key).Result()
	return
}

/*
Redis Sdiff 命令返回第一个集合与其他集合之间的差异，也可以认为说第一个集合中独有的元素。不存在的集合 key 将视为空集。
差集的结果来自前面的 FIRST_KEY ,而不是后面的 OTHER_KEY1，也不是整个 FIRST_KEY OTHER_KEY1..OTHER_KEYN 的差集。
实例:
*/
func (this *Redis) SDiff(v interface{}, keys ...string) (err error) {
	var _result []string
	cmd := this.client.SDiff(this.client.Context(), keys...)
	if _result, err = cmd.Result(); err == nil {
		err = this.codec.UnmarshalSlice(_result, v)
	}
	return
}

/*
Redis Sdiffstore 命令将给定集合之间的差集合存储在指定的集合中。
*/
func (this *Redis) SDiffStore(destination string, keys ...string) (result int64, err error) {
	result, err = this.client.SDiffStore(this.client.Context(), destination, keys...).Result()
	return
}

/*
Redis Sismember 命令返回给定所有给定集合的交集。 不存在的集合 key 被视为空集。 当给定集合当中有一个空集时，结果也为空集(根据集合运算定律)。
*/
func (this *Redis) SInter(v interface{}, keys ...string) (err error) {
	var _result []string
	cmd := this.client.SInter(this.client.Context(), keys...)
	if _result, err = cmd.Result(); err == nil {
		err = this.codec.UnmarshalSlice(_result, v)
	}
	return
}

/*
Redis Sinterstore 决定将给定集合之间的交集在指定的集合中。如果指定的集合已经存在，则将其覆盖
*/
func (this *Redis) SInterStore(destination string, keys ...string) (result int64, err error) {
	result, err = this.client.SInterStore(this.client.Context(), destination, keys...).Result()
	return
}

/*
Redis Sismember 命令判断成员元素是否是集合的成员
*/
func (this *Redis) Sismember(key string, value interface{}) (iskeep bool, err error) {
	iskeep, err = this.client.SIsMember(this.client.Context(), key, value).Result()
	return
}

/*
Redis Smembers 号召返回集合中的所有成员。
*/
func (this *Redis) SMembers(v interface{}, key string) (err error) {
	var _result []string
	cmd := this.client.SMembers(this.client.Context(), key)
	if _result, err = cmd.Result(); err == nil {
		err = this.codec.UnmarshalSlice(_result, v)
	}
	return
}

/*
Redis Smove 命令将指定成员 member 元素从 source 集合移动到 destination 集合。
SMOVE 是原子性操作。
如果 source 集合不存在或不包含指定的 member 元素，则 SMOVE 命令不执行任何操作，仅返回 0 。否则， member 元素从 source 集合中被移除，并添加到 destination 集合中去。
当 destination 集合已经包含 member 元素时， SMOVE 命令只是简单地将 source 集合中的 member 元素删除。
当 source 或 destination 不是集合类型时，返回一个错误。
*/
func (this *Redis) SMove(source string, destination string, member interface{}) (result bool, err error) {
	result, err = this.client.SMove(this.client.Context(), source, destination, member).Result()
	return
}

/*
Redis Spop命令用于移除集合中的指定键的一个或多个随机元素，移除后会返回移除的元素。
该命令类似于Srandmember命令，但SPOP将随机元素从集合中移除并返回，而Srandmember则返回元素，而不是对集合进行任何改动。
*/
func (this *Redis) Spop(key string) (result string, err error) {
	result, err = this.client.SPop(this.client.Context(), key).Result()
	return
}

/*
Redis Srandmember 命令用于返回集合中的一个随机元素。
从 Redis 2.6 版本开始， Srandmember 命令接受可选的 count 参数：
如果 count 为正数，且小于集合基数，那么命令返回一个包含 count 个元素的数组，数组中的元素各不相同。如果 count 大于等于集合基数，那么返回整个集合。
如果 count 为负数，那么命令返回一个数组，数组中的元素可能会重复出现多次，而数组的长度为 count 的绝对值。
该操作和 SPOP 相似，但 SPOP 将随机元素从集合中移除并返回，而 Srandmember 则仅仅返回随机元素，而不对集合进行任何改动。
*/
func (this *Redis) Srandmember(key string) (result string, err error) {
	result, err = this.client.SRandMember(this.client.Context(), key).Result()
	return
}

/*
Redis Srem 呼吁用于移除集合中的一个或多个元素元素，不存在的元素元素会被忽略。
当键不是集合类型，返回一个错误。
在 Redis 2.4 版本以前，SREM 只接受个别成员值。
*/
func (this *Redis) SRem(key string, members ...interface{}) (result int64, err error) {
	result, err = this.client.SRem(this.client.Context(), key, members...).Result()
	return
}

/*
Redis Sunion 命令返回给定集合的并集。
*/
func (this *Redis) SUnion(v interface{}, keys ...string) (err error) {
	var _result []string
	cmd := this.client.SUnion(this.client.Context(), keys...)
	if _result, err = cmd.Result(); err == nil {
		err = this.codec.UnmarshalSlice(_result, v)
	}
	return
}

/*
Redis Sunionstore 命令将给定集合的并集存储在指定的集合 destination 中。如果 destination 已经存在，则将其覆盖。
*/
func (this *Redis) Sunionstore(destination string, keys ...string) (result int64, err error) {
	result, err = this.client.SUnionStore(this.client.Context(), destination, keys...).Result()
	return
}

/*
Redis Sscan 用于继承集合中键的元素，Sscan 继承自Scan。
*/
func (this *Redis) Sscan(key string, _cursor uint64, match string, count int64) (keys []string, cursor uint64, err error) {
	keys, cursor, err = this.client.SScan(this.client.Context(), key, _cursor, match, count).Result()
	return
}
