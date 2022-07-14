package redis_test

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/liwei1dao/lego/sys/redis"
)

func TestMain(m *testing.M) {
	if err := redis.OnInit(nil,
		redis.SetRedisType(redis.Redis_Cluster),
		redis.SetRedis_Cluster_Addr([]string{"10.0.0.9:9001", "10.0.0.9:9002", "10.0.0.9:9003", "10.0.1.45:9004", "10.0.1.45:9005", "10.0.1.45:9006"}),
		redis.SetRedis_Cluster_Password(""),
	); err != nil {
		fmt.Println("err:", err)
		return
	}
	defer os.Exit(m.Run())
	// if err := cache.OnInit(nil, cache.Set_Redis_Addr([]string{"10.0.0.9:9001", "10.0.0.9:9002", "10.0.0.9:9003", "10.0.1.45:9004", "10.0.1.45:9005", "10.0.1.45:9006"}), cache.Set_Redis_Password("")); err != nil {
	// 	fmt.Printf("err:%v\n", err)
	// 	return
	// }
}

func Test_SysIPV6(t *testing.T) {
	err := redis.OnInit(map[string]interface{}{
		"Redis_Single_Addr":     "172.27.100.143:6382",
		"Redis_Single_DB":       0,
		"Redis_Single_Password": "idss@sjzt",
	})
	if err != nil {
		fmt.Printf("Redis:err:%v \n", err)
		return
	}
	fmt.Printf("Redis:succ \n")
	if err = redis.Set("liwei1dao", 123, -1); err != nil {
		fmt.Printf("Redis:err:%v \n", err)
	}
}

func Test_Redis_ExpireatKey(t *testing.T) {
	err := redis.OnInit(map[string]interface{}{
		"Redis_Single_Addr":     "172.20.27.145:10001",
		"Redis_Single_DB":       0,
		"Redis_Single_Password": "li13451234",
	})
	if err != nil {
		fmt.Printf("Redis:err:%v \n", err)
		return
	}
	fmt.Printf("Redis:succ \n")
	if err = redis.Set("liwei1dao", 123, -1); err != nil {
		fmt.Printf("Redis:err:%v \n", err)
	}
	if err = redis.Expire("liwei1dao", 120); err != nil {
		fmt.Printf("Redis:err:%v \n", err)
	}
	fmt.Printf("Redis:end \n")
}

func Test_JsonMarshal(t *testing.T) {
	result, _ := json.Marshal(100)
	fmt.Printf("结果%s \n", string(result))
}

func Test_Redis_SetNX(t *testing.T) {
	err := redis.OnInit(map[string]interface{}{
		"Redis_Single_Addr":     "172.20.27.145:10001",
		"RedisDB":               0,
		"Redis_Single_Password": "li13451234",
	})
	if err != nil {
		fmt.Printf("Redis:err:%v \n", err)
		return
	}

	wg := new(sync.WaitGroup)
	wg.Add(20)
	for i := 0; i < 20; i++ {
		go func(index int) {
			result, err := redis.SetNX("liwei1dao", index)
			fmt.Printf("Redis index:%d result:%d err:%v \n", index, result, err)
			wg.Done()
		}(i)
	}
	wg.Wait()
	fmt.Printf("Redis:end \n")
}
func Test_Redis_Lock(t *testing.T) {
	err := redis.OnInit(map[string]interface{}{
		"Redis_Single_Addr":     "172.20.27.145:10001",
		"Redis_Single_DB":       0,
		"Redis_Single_Password": "li13451234",
	})
	if err != nil {
		fmt.Printf("Redis:err:%v \n", err)
		return
	}
	result, err := redis.Lock("liwei2dao", 100000)
	fmt.Printf("Redis result:%v err:%v  \n", result, err)
}

func Test_Redis_Mutex(t *testing.T) {
	err := redis.OnInit(map[string]interface{}{
		"Redis_Single_Addr":     "172.20.27.145:10001",
		"Redis_Single_DB":       0,
		"Redis_Single_Password": "li13451234",
	})
	if err != nil {
		fmt.Printf("Redis:err:%v \n", err)
		return
	}

	wg := new(sync.WaitGroup)
	wg.Add(20)
	for i := 0; i < 20; i++ {
		go func(index int) {
			if lock, err := redis.NewRedisMutex("liwei1dao_lock"); err != nil {
				fmt.Printf("NewRedisMutex index:%d err:%v\n", index, err)
			} else {
				fmt.Printf("Lock 0 index:%d time:%s\n", index, time.Now().Format("2006/01/02 15:04:05 000"))
				err = lock.Lock()
				fmt.Printf("Lock 1 index:%d time:%s err:%v \n", index, time.Now().Format("2006/01/02 15:04:05 000"), err)
				value := 0
				redis.Get("liwei1dao", &value)
				redis.Set("liwei1dao", value+1, -1)
				lock.Unlock()
				fmt.Printf("Lock 2 index:%d time:%s\n", index, time.Now().Format("2006/01/02 15:04:05 000"))
			}
			wg.Done()
		}(i)
	}
	wg.Wait()
	fmt.Printf("Redis:end \n")
}

func Test_Redis_Type(t *testing.T) {
	err := redis.OnInit(map[string]interface{}{
		"Redis_Single_Addr":     "172.20.27.145:10001",
		"Redis_Single_DB":       1,
		"Redis_Single_Password": "li13451234",
	})
	if err != nil {
		fmt.Printf("Redis:err:%v \n", err)
		return
	}
	fmt.Printf("Redis:succ \n")

	if ty, err := redis.Type("test_set"); err != nil {
		fmt.Printf("Test_Redis_Type:err:%v \n", err)
	} else {
		fmt.Printf("Test_Redis_Type:%s \n", ty)
	}
}

type TestAny struct {
	SubName string `json:"subname"`
	Age     int32  `json:"age"`
}
type TestData struct {
	Name string   `json:"name"`
	Agr  int      `json:"agr"`
	Sub  *TestAny `json:"sub"`
}

func Test_Redis_Encoder_Struct(t *testing.T) {
	err := redis.Set("test:1001", &TestData{Name: "liwei1dao", Agr: 12}, -1)
	fmt.Printf("err:%v\n", err)
}
func Test_Redis_Encoder_int(t *testing.T) {
	err := redis.Set("test:1002", 856, -1)
	fmt.Printf("err:%v \n", err)
	data := 0
	err = redis.Get("test:1002", &data)
	fmt.Printf("data:%d err:%v\n", data, err)
}

func Test_Redis_Encoder_Hash(t *testing.T) {
	name := ""
	err := redis.HGet("test:103", "name", &name)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(name)
}

//测试redis lua 脚本
func Test_Redis_Lua(t *testing.T) {
	script := redis.NewScript(`
		local goodsSurplus
		local flag
		local existUserIds    = tostring(KEYS[1])
		local memberUid       = tonumber(ARGV[1])
		local goodsSurplusKey = tostring(KEYS[2])
		local hasBuy = redis.call("sIsMember", existUserIds, memberUid)

		if hasBuy ~= 0 then
		return 0
		end
		

		goodsSurplus =  redis.call("GET", goodsSurplusKey)
		if goodsSurplus == false then
		return 0
		end
		
		-- 没有剩余可抢购物品
		goodsSurplus = tonumber(goodsSurplus)
		if goodsSurplus <= 0 then
		return 0
		end
		
		flag = redis.call("SADD", existUserIds, memberUid)
		flag = redis.call("DECR", goodsSurplusKey)
		return 1
	`)
	sha, err := script.Result()
	if err != nil {
		log.Fatalln(err)
	}
	ret := redis.EvalSha(sha, []string{
		"hadBuyUids",
		"goodsSurplus",
	}, "userId")
	if result, err := ret.Result(); err != nil {
		log.Fatalf("Execute Redis fail: %v", err.Error())
	} else {
		fmt.Println("")
		fmt.Printf("userid: %s, result: %d", "userId", result)
	}
}
