package redis

import "github.com/liwei1dao/lego/core"

/*
Redis Sadd 命令将一个或多个成员元素加入到集合中，已经存在于集合的成员元素将被忽略。
假如集合 key 不存在，则创建一个只包含添加的元素作成员的集合。
当集合 key 不是集合类型时，返回一个错误。
*/
func (this *Redis) SAdd(key core.Redis_Key, values ...interface{}) (err error) {
	agrs := make([]interface{}, 0)
	agrs = append(agrs, "SADD")
	agrs = append(agrs, key)
	for _, v := range values {
		result, _ := this.encode(v)
		agrs = append(agrs, result)
	}
	err = this.client.Do(this.getContext(), agrs...).Err()
	return
}

/*
Redis Scard 命令返回集合中元素的数量
*/
func (this *Redis) Scard(key core.Redis_Key) (result int, err error) {
	result, err = this.client.Do(this.getContext(), "SCARD", key).Int()
	return
}

/*
Redis Sismember 命令判断成员元素是否是集合的成员
*/
func (this *Redis) Sismember(key core.Redis_Key, value interface{}) (iskeep bool, err error) {
	iskeep, err = this.client.Do(this.getContext(), "SISMEMBER", key).Bool()
	return
}
