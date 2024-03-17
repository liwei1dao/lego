package pipe

import (
	"time"
)

/* Key *******************************************************************************/

///删除redis key
func (this *RedisPipe) Delete(key string) (err error) {
	err = this.client.Do(this.ctx, "DEL", key).Err()
	return
}

///判断是否存在key
func (this *RedisPipe) ExistsKey(key string) (iskeep bool, err error) {
	iskeep, err = this.client.Do(this.ctx, "EXISTS", key).Bool()
	return
}

///设置key的过期时间 单位以秒级
func (this *RedisPipe) Expire(key string, expiration time.Duration) (err error) {
	this.client.Expire(this.ctx, key, expiration)
	return
}

///设置key的过期时间戳 秒级时间戳
func (this *RedisPipe) ExpireAt(key string, tm time.Time) (err error) {
	err = this.client.ExpireAt(this.ctx, key, tm).Err()
	return
}

///设置key的过期时间 单位以毫秒级
func (this *RedisPipe) PExpire(key string, expiration time.Duration) (err error) {
	err = this.client.PExpire(this.ctx, key, expiration).Err()
	return
}

///设置key的过期时间戳 单位以豪秒级
func (this *RedisPipe) PExpireAt(key string, tm time.Time) (err error) {
	err = this.client.PExpireAt(this.ctx, key, tm).Err()
	return
}

///移除Key的过期时间
func (this *RedisPipe) Persist(key string) (err error) {
	err = this.client.Persist(this.ctx, key).Err()
	return
}

///获取key剩余过期时间 单位毫秒
func (this *RedisPipe) PTTL(key string) (leftexpire time.Duration, err error) {
	leftexpire, err = this.client.PTTL(this.ctx, key).Result()
	return
}

///获取key剩余过期时间 单位秒
func (this *RedisPipe) TTL(key string) (leftexpire time.Duration, err error) {
	leftexpire, err = this.client.TTL(this.ctx, key).Result()
	return
}

///重命名Key
func (this *RedisPipe) Rename(oldkey string, newkey string) (err error) {
	err = this.client.Rename(this.ctx, oldkey, newkey).Err()
	return
}

///重命名key 在新的 key 不存在时修改 key 的名称
func (this *RedisPipe) RenameNX(oldkey string, newkey string) (err error) {
	err = this.client.RenameNX(this.ctx, oldkey, newkey).Err()
	return
}

///判断是否存在key pattern:key*
func (this *RedisPipe) Keys(pattern string) (keys []string, err error) {
	keys, err = this.client.Keys(this.ctx, pattern).Result()
	return
}

///获取键类型
func (this *RedisPipe) Type(key string) (ty string, err error) {
	ty, err = this.client.Type(this.ctx, key).Result()
	return
}
