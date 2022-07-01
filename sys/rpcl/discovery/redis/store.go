package redis

import (
	"github.com/liwei1dao/lego/sys/redis"
)

type RedisStore struct {
	client redis.ISys
}
