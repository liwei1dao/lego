package redis

import (
	"encoding/json"
	"fmt"
	"testing"
)

func Test_SysIPV6(t *testing.T) {
	// this.client = redis.NewClient(&redis.Options{
	// 	Addr:     this.options.RedisUrl,
	// 	Password: this.options.RedisPassword,
	// 	DB:       this.options.RedisDB,
	// 	PoolSize: this.options.PoolSize, // 连接池大小
	// })
	// _, err = this.client.Ping(this.getContext()).Result()
}

func Test_JsonMarshal(t *testing.T) {
	result, _ := json.Marshal(100)
	fmt.Printf("结果%s \n", string(result))
}
