package redis

import (
	"context"
	"encoding/json"

	"github.com/go-redis/redis/v8"
	"google.golang.org/protobuf/proto"
)

func newSys(options Options) (sys *Redis, err error) {
	sys = &Redis{options: options}
	err = sys.init()
	return
}

type Redis struct {
	options Options
	client  *redis.Client
}

func (this *Redis) init() (err error) {
	this.client = redis.NewClient(&redis.Options{
		Addr:     this.options.RedisUrl,
		Password: this.options.RedisPassword,
		DB:       this.options.RedisDB,
	})
	_, err = this.client.Ping(this.getContext()).Result()
	return
}

func (this *Redis) getContext() (ctx context.Context) {
	ctx, _ = context.WithTimeout(context.Background(), this.options.TimeOut)
	return
}

///数据编码
func (this *Redis) encode(value interface{}) (result []byte, err error) {
	if this.options.RedisStorageType == JsonData {
		result, err = json.Marshal(value)
	} else {
		if _, ok := value.(proto.Message); ok {
			result, err = proto.Marshal(value.(proto.Message))
		} else {
			result, err = json.Marshal(value)
		}
	}
	return
}

func (this *Redis) decode(value []byte, result interface{}) (err error) {
	if this.options.RedisStorageType == JsonData {
		err = json.Unmarshal(value, result)
	} else {
		if _, ok := result.(proto.Message); ok {
			err = proto.Unmarshal(value, result.(proto.Message))
		} else {
			err = json.Unmarshal(value, result)
		}
	}
	return
}
