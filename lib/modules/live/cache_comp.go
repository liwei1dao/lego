package live

import (
	"github.com/liwei1dao/lego/core"
	"github.com/liwei1dao/lego/core/cbase"
	"github.com/liwei1dao/lego/sys/redis"
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
	this.cache, err = redis.NewSys(redis.SetRedisUrl(this.options.GetCacheAddr()))
	return
}

func (this *CacheComp) GetChannel(key string) (channel string, err error) {
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
	// key = uid.RandStringRunes(48)
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
	// if newKey, err = r.redisCli.Get(channel).Result(); err == redis.Nil {
	// 	newKey, err = r.SetKey(channel)
	// 	log.Debugf("[KEY] new channel [%s]: %s", channel, newKey)
	// 	return
	// }
	return
}

func (this *CacheComp) DeleteChannel(channel string) bool {
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
