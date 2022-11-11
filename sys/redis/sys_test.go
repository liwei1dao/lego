package redis_test

import (
	"context"
	"fmt"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/liwei1dao/lego/utils/codec/json"

	"github.com/go-redis/redis/v8"
	lgredis "github.com/liwei1dao/lego/sys/redis"
)

func TestMain(m *testing.M) {
	if err := lgredis.OnInit(nil,
		lgredis.SetRedisType(lgredis.Redis_Single),
		lgredis.SetRedis_Single_Addr("10.0.0.9:6986"),
		lgredis.SetRedis_Single_Password("li13451234"),
		lgredis.SetRedis_Single_DB(2),
	); err != nil {
		fmt.Println("err:", err)
		return
	}
	defer os.Exit(m.Run())
}

func Test_Redis_ExpireatKey(t *testing.T) {
	var err error
	if err = lgredis.Set("liwei1dao", 123, -1); err != nil {
		fmt.Printf("Redis:err:%v \n", err)
	}
	if err = lgredis.Expire("liwei1dao", 120); err != nil {
		fmt.Printf("Redis:err:%v \n", err)
	}
	fmt.Printf("Redis:end \n")
}

func Test_JsonMarshal(t *testing.T) {
	result, _ := json.Marshal(100)
	fmt.Printf("结果%s \n", string(result))
}

func Test_Redis_SetNX(t *testing.T) {
	wg := new(sync.WaitGroup)
	wg.Add(20)
	for i := 0; i < 20; i++ {
		go func(index int) {
			result, err := lgredis.SetNX("liwei1dao", index)
			fmt.Printf("Redis index:%d result:%d err:%v \n", index, result, err)
			wg.Done()
		}(i)
	}
	wg.Wait()
	fmt.Printf("Redis:end \n")
}
func Test_Redis_Lock(t *testing.T) {
	result, err := lgredis.Lock("liwei2dao", 100000)
	fmt.Printf("Redis result:%v err:%v  \n", result, err)
}

func Test_Redis_Mutex(t *testing.T) {
	wg := new(sync.WaitGroup)
	wg.Add(20)
	for i := 0; i < 20; i++ {
		go func(index int) {
			if lock, err := lgredis.NewRedisMutex("liwei1dao_lock"); err != nil {
				fmt.Printf("NewRedisMutex index:%d err:%v\n", index, err)
			} else {
				fmt.Printf("Lock 0 index:%d time:%s\n", index, time.Now().Format("2006/01/02 15:04:05 000"))
				err = lock.Lock()
				fmt.Printf("Lock 1 index:%d time:%s err:%v \n", index, time.Now().Format("2006/01/02 15:04:05 000"), err)
				value := 0
				lgredis.Get("liwei1dao", &value)
				lgredis.Set("liwei1dao", value+1, -1)
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
	if ty, err := lgredis.Type("test_set"); err != nil {
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
	err := lgredis.Set("test:1001", &TestData{Name: "liwei1dao", Agr: 12}, -1)
	fmt.Printf("err:%v\n", err)
}
func Test_Redis_Encoder_int(t *testing.T) {
	err := lgredis.Set("test:1002", 856, -1)
	fmt.Printf("err:%v \n", err)
	data := 0
	err = lgredis.Get("test:1002", &data)
	fmt.Printf("data:%d err:%v\n", data, err)
}

func Test_Redis_Encoder_Hash(t *testing.T) {
	name := ""
	err := lgredis.HGet("test:103", "name", &name)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(name)
}

//测试redis lua 脚本
func Test_Redis_Lua_HGETALL(t *testing.T) {
	script := lgredis.NewScript(`
	local key    = tostring(KEYS[1])
	local keys = lgredis.call("HGETALL", key)
	local data = {}
	local n = 1
	for i, v in ipairs(keys) do
		if i%2 == 0 then
			data[n] = lgredis.call("HGETALL", v)
			n = n+1
		end
	end 
	return data
`)
	sha, err := script.Result()
	if err != nil {
		fmt.Println(err)
	}
	ret := lgredis.EvalSha(lgredis.Context(), sha, []string{"items:0_62c259916d8cf3e4e06311a8"})
	if result, err := ret.Result(); err != nil {
		fmt.Printf("Execute Redis err: %v", err.Error())
	} else {
		temp1 := result.([]interface{})
		data := make([]map[string]string, len(temp1))
		for i, v := range temp1 {
			temp2 := v.([]interface{})
			data[i] = make(map[string]string)
			for n := 0; n < len(temp2); n += 2 {
				data[i][temp2[n].(string)] = temp2[n+1].(string)
			}
		}
		fmt.Printf("data: %v", data)
	}
}

//测试redis lua 脚本
func Test_Redis_Lua_HSETALL(t *testing.T) {
	script := lgredis.NewScript(`
	local n = 1
	local key = ""
	for i, v in ipairs(KEYS) do
		key = v
		local argv = {}
		for i=n,#ARGV,1 do
			n = n+1
			if ARGV[i] == "#end" then
				lgredis.call("HMSet", key,unpack(argv))
				break
			else
				table.insert(argv, ARGV[i])
			end
		end
	end 
	return "OK"
`)
	sha, err := script.Result()
	if err != nil {
		fmt.Println(err)
	}
	ret := lgredis.EvalSha(lgredis.Context(), sha, []string{"test_HMSet", "test_HMSet_1"}, "a", "1", "b", "2", "#end", "a1", "11", "b", "21", "#end")
	if result, err := ret.Result(); err != nil {
		fmt.Printf("Execute Redis err: %v", err.Error())
	} else {
		fmt.Printf("data: %v", result)
	}
}

func Test_Redis_Lua_Queue(t *testing.T) {
	//Redis 自定义脚本 批量读取队列数据
	script := lgredis.NewScript(`
	local count	= tonumber(ARGV[1])
	local k   = tostring(ARGV[3])
	local keys = {}
	local out = {}
	local n = 1
	for i, v in ipairs(KEYS) do
		if (i == 1) then
			for i=n,#ARGV,1 do
				n = n+1
				if ARGV[i] == "#end" then
					break
				end
			end
		elseif (i == 2) then
			for i=n,#ARGV,1 do
				n = n+1
				if ARGV[i] == "#end" then
					break
				end
			end
		else
			local key = v
			local argv = {}
			table.insert(keys, key)
			for i=n,#ARGV,1 do
				n = n+1
				if ARGV[i] == "#end" then
					lgredis.call("HMSet", key,unpack(argv))
					break
				else
					table.insert(argv, ARGV[i])
				end
			end
		end
	end
	lgredis.call("RPush", k,unpack(keys))
	local c = tonumber(lgredis.call("LLEN", k))
	if (c > count) then
		local off = c-count
		out = lgredis.call("LRANGE", k,0,off-1)
		lgredis.call("LTRIM", k,off,-1)
		for i, v in ipairs(out) do
			lgredis.call("DEL", v)
		end
	end
	return out
	`)
	sha, err := script.Result()
	if err != nil {
		fmt.Println(err)
	}
	ret := lgredis.EvalSha(lgredis.Context(), sha, []string{"count", "key", "Test_Redis_Lua_Queue_1", "Test_Redis_Lua_Queue_2", "Test_Redis_Lua_Queue_3"}, "10", "#end", "Test_Redis_Lua_Queue", "#end", "a", "1", "b", "2", "#end", "a1", "11", "b1", "21", "#end", "a3", "13", "b3", "23", "#end")
	if result, err := ret.Result(); err != nil {
		fmt.Printf("Execute Redis err: %v", err.Error())
	} else {
		fmt.Printf("data: %v", result)
	}
	return
}

//测试数据管道
func Test_RedisPipe(t *testing.T) {
	keys, _ := lgredis.HGetAllToMapString("items:2_62ecff47148b1ef7184ced8f{items}")
	pipe := lgredis.RedisPipe(context.TODO())
	result := make([]*redis.StringStringMapCmd, 0)
	for _, v := range keys {
		cmd := pipe.HGetAllToMapString(v)
		result = append(result, cmd)
	}
	pipe.Exec()
	for _, v := range result {
		data, err := v.Result()
		fmt.Printf("RedisPipe data:%v err: %v \n", data, err)
	}

}
