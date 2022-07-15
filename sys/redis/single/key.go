package single

import "time"

/* Key *******************************************************************************/

///删除redis key
func (this *Redis) Delete(key string) (err error) {
	err = this.client.Do(this.client.Context(), "DEL", key).Err()
	return
}

///判断是否存在key
func (this *Redis) ExistsKey(key string) (iskeep bool, err error) {
	iskeep, err = this.client.Do(this.client.Context(), "EXISTS", key).Bool()
	return
}

///设置key的过期时间 单位以秒级
func (this *Redis) Expire(key string, expiration time.Duration) (err error) {
	this.client.Expire(this.client.Context(), key, expiration)
	return
}

///设置key的过期时间戳 秒级时间戳
func (this *Redis) ExpireAt(key string, tm time.Time) (err error) {
	err = this.client.ExpireAt(this.client.Context(), key, tm).Err()
	return
}

///设置key的过期时间 单位以毫秒级
func (this *Redis) PExpire(key string, expiration time.Duration) (err error) {
	err = this.client.PExpire(this.client.Context(), key, expiration).Err()
	return
}

///设置key的过期时间戳 单位以豪秒级
func (this *Redis) PExpireAt(key string, tm time.Time) (err error) {
	err = this.client.PExpireAt(this.client.Context(), key, tm).Err()
	return
}

///移除Key的过期时间
func (this *Redis) Persist(key string) (err error) {
	err = this.client.Persist(this.client.Context(), key).Err()
	return
}

///获取key剩余过期时间 单位毫秒
func (this *Redis) PTTL(key string) (leftexpire time.Duration, err error) {
	leftexpire, err = this.client.PTTL(this.client.Context(), key).Result()
	return
}

///获取key剩余过期时间 单位秒
func (this *Redis) TTL(key string) (leftexpire time.Duration, err error) {
	leftexpire, err = this.client.TTL(this.client.Context(), key).Result()
	return
}

///重命名Key
func (this *Redis) Rename(oldkey string, newkey string) (err error) {
	err = this.client.Rename(this.client.Context(), oldkey, newkey).Err()
	return
}

///重命名key 在新的 key 不存在时修改 key 的名称
func (this *Redis) RenameNX(oldkey string, newkey string) (err error) {
	err = this.client.RenameNX(this.client.Context(), oldkey, newkey).Err()
	return
}

///判断是否存在key pattern:key*
func (this *Redis) Keys(pattern string) (keys []string, err error) {
	keys, err = this.client.Keys(this.client.Context(), pattern).Result()
	return
}

///获取键类型
func (this *Redis) Type(key string) (ty string, err error) {
	ty, err = this.client.Type(this.client.Context(), key).Result()
	return
}
