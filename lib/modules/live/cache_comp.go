package live

import (
	"github.com/liwei1dao/lego/core"
	"github.com/liwei1dao/lego/core/cbase"
	"github.com/liwei1dao/lego/sys/log"
	"github.com/liwei1dao/lego/sys/redis"
	"github.com/liwei1dao/lego/utils/uid"
)

const (
	Redis_Channel core.Redis_Key = "LiveChannel" //用户数据缓存
	Redis_Key     core.Redis_Key = "LiveKey"     //用户数据缓存
)

//主机信息监控
type CacheComp struct {
	cbase.ModuleCompBase
	options IOptions
	cache   redis.IRedisFactory
}

func (this *CacheComp) Init(service core.IService, module core.IModule, comp core.IModuleComp, options core.IModuleOptions) (err error) {
	err = this.ModuleCompBase.Init(service, module, comp, options)
	this.options = options.(IOptions)
	if this.cache, err = redis.NewSys(redis.SetRedisUrl(this.options.GetCacheAddr())); err == nil {

	}
	return
}

func (this *CacheComp) GetChannel(key string) (channel string, err error) {
	err = this.cache.GetPool().GetKey_MapByKey(string(Redis_Key), key, &channel)
	// if !saveInLocal {
	// 	return r.redisCli.Get(key).Result()
	// }

	// chann, found := r.localCache.Get(key)
	// if found {
	// 	return chann.(string), nil
	// } else {
	// 	return "", fmt.Errorf("%s does not exists", key)
	// }
	return
}

// set/reset a random key for channel
func (this *CacheComp) SetChannelKey(channel string) (key string, err error) {
	key = uid.RandStringRunes(48)
	if err = this.cache.GetPool().SetKey_Map(string(Redis_Channel), map[string]interface{}{channel: key}); err == nil {
		err = this.cache.GetPool().SetKey_Map(string(Redis_Key), map[string]interface{}{key: channel})
	}
	// if _, err = r.redisCli.Get(key).Result(); err == redis.Nil {
	// 	err = r.redisCli.Set(channel, key, 0).Err()
	// 	if err != nil {
	// 		return
	// 	}

	// 	err = r.redisCli.Set(key, channel, 0).Err()
	// 	return
	// } else if err != nil {
	// 	return
	// }
	return
}

func (this *CacheComp) GetChannelKey(channel string) (newKey string, err error) {

	if err = this.cache.GetPool().GetKey_MapByKey(string(Redis_Channel), channel, &newKey); err != nil {
		newKey, err = this.SetChannelKey(channel)
		log.Debugf("[KEY] new channel [%s]: %s", channel, newKey)
	}
	// if newKey, err = r.redisCli.Get(channel).Result(); err == redis.Nil {
	// 	newKey, err = r.SetKey(channel)
	// 	log.Debugf("[KEY] new channel [%s]: %s", channel, newKey)
	// 	return
	// }
	return
}

func (this *CacheComp) DeleteChannel(channel string) bool {
	var (
		key string
		err error
	)
	if err = this.cache.GetPool().GetKey_MapByKey(string(Redis_Channel), channel, &key); err == nil {
		if err = this.cache.GetPool().DelKey_MapKey(string(Redis_Channel), channel); err == nil {
			if err = this.cache.GetPool().DelKey_MapKey(string(Redis_Key), key); err == nil {
				return true
			}
		}
	}

	// if !saveInLocal {
	// 	return r.redisCli.Del(channel).Err() != nil
	// }

	// key, ok := r.localCache.Get(channel)
	// if ok {
	// 	r.localCache.Delete(channel)
	// 	r.localCache.Delete(key.(string))
	// 	return true
	// }
	// return false
	return false
}
