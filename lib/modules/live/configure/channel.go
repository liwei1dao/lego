package configure

import (
	"fmt"

	"github.com/go-redis/redis/v7"
	"github.com/liwei1dao/lego/sys/log"
	"github.com/liwei1dao/lego/utils/uid"
	"github.com/patrickmn/go-cache"
)

var saveInLocal = true

var RoomKeys = &RoomKeysType{
	localCache: cache.New(cache.NoExpiration, 0),
}

type RoomKeysType struct {
	redisCli   *redis.Client
	localCache *cache.Cache
}

// set/reset a random key for channel
func (r *RoomKeysType) SetKey(channel string) (key string, err error) {
	if !saveInLocal {
		for {
			key = uid.RandStringRunes(48)
			if _, err = r.redisCli.Get(key).Result(); err == redis.Nil {
				err = r.redisCli.Set(channel, key, 0).Err()
				if err != nil {
					return
				}

				err = r.redisCli.Set(key, channel, 0).Err()
				return
			} else if err != nil {
				return
			}
		}
	}

	for {
		key = uid.RandStringRunes(48)
		if _, found := r.localCache.Get(key); !found {
			r.localCache.SetDefault(channel, key)
			r.localCache.SetDefault(key, channel)
			break
		}
	}
	return
}

func (r *RoomKeysType) GetKey(channel string) (newKey string, err error) {
	if !saveInLocal {
		if newKey, err = r.redisCli.Get(channel).Result(); err == redis.Nil {
			newKey, err = r.SetKey(channel)
			log.Debugf("[KEY] new channel [%s]: %s", channel, newKey)
			return
		}

		return
	}

	var key interface{}
	var found bool
	if key, found = r.localCache.Get(channel); found {
		return key.(string), nil
	}
	newKey, err = r.SetKey(channel)
	log.Debugf("[KEY] new channel [%s]: %s", channel, newKey)
	return
}

func (r *RoomKeysType) GetChannel(key string) (channel string, err error) {
	if !saveInLocal {
		return r.redisCli.Get(key).Result()
	}

	chann, found := r.localCache.Get(key)
	if found {
		return chann.(string), nil
	} else {
		return "", fmt.Errorf("%s does not exists", key)
	}
}
