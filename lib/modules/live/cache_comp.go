package live

import (
	"github.com/liwei1dao/lego/core"
	"github.com/liwei1dao/lego/core/cbase"
	"github.com/liwei1dao/lego/lib/modules/live/utils/uid"
	"github.com/liwei1dao/lego/sys/log"
	"github.com/liwei1dao/lego/sys/redis"
)

const (
	Redis_Channel string = "LiveChannel" //用户数据缓存
	Redis_Key     string = "LiveKey"     //用户数据缓存
)

//主机信息监控
type CacheComp struct {
	cbase.ModuleCompBase
	options IOptions
	cache   redis.IRedis
}

func (this *CacheComp) Init(service core.IService, module core.IModule, comp core.IModuleComp, options core.IModuleOptions) (err error) {
	err = this.ModuleCompBase.Init(service, module, comp, options)
	this.options = options.(IOptions)
	if this.cache, err = redis.NewSys(
		redis.SetRedisType(redis.Redis_Single),
		redis.SetRedis_Single_Addr(this.options.GetCacheAddr()),
		redis.SetRedis_Single_Password(this.options.GetCachePassword()),
		redis.SetRedis_Single_DB(this.options.GetCacheDB())); err == nil {
		this.cache.Delete(Redis_Channel)
		this.cache.Delete(Redis_Key)
	}
	return
}

func (this *CacheComp) GetChannel(key string) (channel string, err error) {
	err = this.cache.HGet(Redis_Key, key, &channel)
	return
}

// set/reset a random key for channel
func (this *CacheComp) SetChannelKey(channel string) (key string, err error) {
	key = uid.RandStringRunes(48)
	if err = this.cache.HSet(Redis_Channel, channel, key); err == nil {
		err = this.cache.HSet(Redis_Key, key, channel)
	}
	return
}

func (this *CacheComp) GetChannelKey(channel string) (newKey string, err error) {

	if err = this.cache.HGet(Redis_Channel, channel, &newKey); err != nil {
		newKey, err = this.SetChannelKey(channel)
		log.Debugf("[KEY] new channel [%s]: %s", channel, newKey)
	}
	return
}

func (this *CacheComp) DeleteChannel(channel string) bool {
	var (
		key string
		err error
	)
	if err = this.cache.HGet(Redis_Channel, channel, &key); err == nil {
		if err = this.cache.HDel(Redis_Channel, channel); err == nil {
			if err = this.cache.HDel(Redis_Key, key); err == nil {
				return true
			}
		}
	}
	return false
}
