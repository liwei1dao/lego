package redis

import (
	"encoding/json"
	"fmt"
	"sync"
	"testing"
	"time"
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

func Test_Redis_ExpireatKey(t *testing.T) {
	err := OnInit(map[string]interface{}{
		"RedisUrl":      "172.20.27.145:10001",
		"RedisDB":       0,
		"RedisPassword": "li13451234",
	})
	if err != nil {
		fmt.Printf("Redis:err:%v \n", err)
		return
	}
	fmt.Printf("Redis:succ \n")
	if err = Set("liwei1dao", 123, -1); err != nil {
		fmt.Printf("Redis:err:%v \n", err)
	}
	if err = ExpireKey("liwei1dao", 120); err != nil {
		fmt.Printf("Redis:err:%v \n", err)
	}
	fmt.Printf("Redis:end \n")
}

func Test_JsonMarshal(t *testing.T) {
	result, _ := json.Marshal(100)
	fmt.Printf("结果%s \n", string(result))
}

func Test_Redis_SetNX(t *testing.T) {
	err := OnInit(map[string]interface{}{
		"RedisUrl":      "172.20.27.145:10001",
		"RedisDB":       0,
		"RedisPassword": "li13451234",
	})
	if err != nil {
		fmt.Printf("Redis:err:%v \n", err)
		return
	}

	wg := new(sync.WaitGroup)
	wg.Add(20)
	for i := 0; i < 20; i++ {
		go func(index int) {
			result, err := SetNX("liwei1dao", index)
			fmt.Printf("Redis index:%d result:%d err:%v \n", index, result, err)
			wg.Done()
		}(i)
	}
	wg.Wait()
	fmt.Printf("Redis:end \n")
}
func Test_Redis_Lock(t *testing.T) {
	err := OnInit(map[string]interface{}{
		"RedisUrl":      "172.20.27.145:10001",
		"RedisDB":       0,
		"RedisPassword": "li13451234",
	})
	if err != nil {
		fmt.Printf("Redis:err:%v \n", err)
		return
	}
	result, err := Lock("liwei2dao", 100000)
	fmt.Printf("Redis result:%v err:%v  \n", result, err)
}

func Test_Redis_Mutex(t *testing.T) {
	err := OnInit(map[string]interface{}{
		"RedisUrl":      "172.20.27.145:10001",
		"RedisDB":       0,
		"RedisPassword": "li13451234",
	})
	if err != nil {
		fmt.Printf("Redis:err:%v \n", err)
		return
	}

	wg := new(sync.WaitGroup)
	wg.Add(20)
	for i := 0; i < 20; i++ {
		go func(index int) {
			if lock, err := NewRedisMutex("liwei1dao_lock"); err != nil {
				fmt.Printf("NewRedisMutex index:%d err:%v\n", index, err)
			} else {
				fmt.Printf("Lock 0 index:%d time:%s\n", index, time.Now().Format("2006/01/02 15:04:05 000"))
				err = lock.Lock()
				fmt.Printf("Lock 1 index:%d time:%s err:%v \n", index, time.Now().Format("2006/01/02 15:04:05 000"), err)
				value := 0
				Get("liwei1dao", &value)
				Set("liwei1dao", value+1, -1)
				lock.Unlock()
				fmt.Printf("Lock 2 index:%d time:%s\n", index, time.Now().Format("2006/01/02 15:04:05 000"))
			}
			wg.Done()
		}(i)
	}
	wg.Wait()
	fmt.Printf("Redis:end \n")
}
