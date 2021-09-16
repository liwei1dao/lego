package redis

import (
	"github.com/go-redis/redis/v8"
	"github.com/liwei1dao/lego/core"
)

/* Key *******************************************************************************/

///删除redis key
func (this *Redis) Delete(key core.Redis_Key) (err error) {
	err = this.client.Do(this.getContext(), "DEL", string(key)).Err()
	return
}

///判断是否存在key
func (this *Redis) ExistsKey(key core.Redis_Key) (iskeep bool, err error) {
	iskeep, err = this.client.Do(this.getContext(), "EXISTS", string(key)).Bool()
	return
}

///设置key的过期时间 单位以秒级
func (this *Redis) ExpireKey(key core.Redis_Key, expire int) (err error) {
	err = this.client.Do(this.getContext(), "EXPIRE", key, expire).Err()
	return
}

///设置key的过期时间戳 秒级时间戳
func (this *Redis) ExpireatKey(key core.Redis_Key, expire_unix int64) (err error) {
	err = this.client.Do(this.getContext(), "EXPIREAT", key, expire_unix).Err()
	return
}

///设置key的过期时间 单位以毫秒级
func (this *Redis) Pexpirekey(key core.Redis_Key, expire int) (err error) {
	err = this.client.Do(this.getContext(), "PEXPIRE", key, expire).Err()
	return
}

///设置key的过期时间戳 单位以豪秒级
func (this *Redis) PexpireatKey(key core.Redis_Key, expire_unix int64) (err error) {
	err = this.client.Do(this.getContext(), "PEXPIREAT", key, expire_unix).Err()
	return
}

///移除Key的过期时间
func (this *Redis) PersistKey(key core.Redis_Key) (err error) {
	err = this.client.Do(this.getContext(), "PERSIST", key).Err()
	return
}

///获取key剩余过期时间 单位毫秒
func (this *Redis) PttlKey(key core.Redis_Key) (leftexpire int64, err error) {
	leftexpire, err = this.client.Do(this.getContext(), "PTTL", key).Int64()
	return
}

///获取key剩余过期时间 单位秒
func (this *Redis) TtlKey(key core.Redis_Key) (leftexpire int64, err error) {
	leftexpire, err = this.client.Do(this.getContext(), "TTL", key).Int64()
	return
}

///重命名Key
func (this *Redis) RenameKye(oldkey core.Redis_Key, newkey string) (err error) {
	err = this.client.Do(this.getContext(), "RENAME", oldkey, newkey).Err()
	return
}

///重命名key 在新的 key 不存在时修改 key 的名称
func (this *Redis) RenamenxKey(oldkey core.Redis_Key, newkey string) (err error) {
	err = this.client.Do(this.getContext(), "RENAMENX", oldkey, newkey).Err()
	return
}

///判断是否存在key pattern:key*
func (this *Redis) Keys(pattern core.Redis_Key) (keys []string, err error) {
	cmd := redis.NewStringSliceCmd(this.getContext(), "KEYS", string(pattern))
	this.client.Process(this.getContext(), cmd)
	keys, err = cmd.Result()
	return
}
