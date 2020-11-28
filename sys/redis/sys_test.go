package redis

import (
	"reflect"
	"sync"
	"testing"
	"time"
)

func TestRedisMutex_Lock(t *testing.T) {
	if err := OnInit(map[string]interface{}{
		"RedisUrl":"redis://127.0.0.1:6379/1"
	}); err != nil {
		t.Errorf("初始化 redis 失败 err:%s", err.Error())
		return
	}
	wg := sync.WaitGroup{}
	wg.Add(100)
	for i := 0; i < 100; i++ {
		go func() {
			rlock := NewRedisMutex("roomIds")
			start := time.Now()
			if err := rlock.Lock(); err == nil {
				t.Log("通过 ", time.Now().Sub(start))
			} else {
				t.Error("超时通过 ", time.Now().Sub(start))
			}
			time.Sleep(time.Millisecond * 5)
			rlock.Unlock()
			wg.Done()
		}()
	}
	wg.Wait()
	t.Log("结束测试")
}

func TestRedisList_Lock(t *testing.T) {
	if err := OnInit(map[string]interface{}{
		"RedisUrl":"redis://127.0.0.1:6379/1"
	}); err != nil {
		t.Errorf("初始化 redis 失败 err:%s", err.Error())
		return
	}
	var item string
	GetPool().SetListByLPush("TestList", []interface{}{"liwei1dao"})
	GetPool().SetListByLPush("TestList", []interface{}{"liwei2dao"})
	GetPool().SetListByLPush("TestList", []interface{}{"liwei3dao"})
	GetPool().SetListByLPush("TestList", []interface{}{"liwei4dao"})
	GetPool().SetListByLPush("TestList", []interface{}{"liwei5dao"})
	data := GetPool().GetListByLrange("TestList", 0, 100, reflect.TypeOf(&item))
	// GetPool().GetListByLPop("TestList", &item)
	t.Logf("结束测试 data:%v", data)
}
