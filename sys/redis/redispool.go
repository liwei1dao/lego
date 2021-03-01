package redis

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"reflect"
	"time"

	"github.com/liwei1dao/lego/sys/log"

	"github.com/garyburd/redigo/redis"
)

func NewRedisPool(pool *redis.Pool) *RedisPool {
	p := &RedisPool{
		pool,
	}
	return p
}

type RedisPool struct {
	Pool *redis.Pool
}

//自定义处理方法
func (this *RedisPool) CustomFunc(f func(pool redis.Conn)) {
	if f == nil {
		log.Errorf("Redis CustomFunc 异常调用")
		return
	}
	pool := this.Pool.Get()
	defer pool.Close()
	f(pool)
}

//判断键是否存在
func (this *RedisPool) ContainsKey(key string) (iskeep bool, err error) {
	pool := this.Pool.Get()
	defer pool.Close()
	iskeep, err = redis.Bool(pool.Do("EXISTS", key))
	return
}

//查询缓存keys
func (this *RedisPool) QueryKeys(patternkey string) (keys []string, err error) {
	pool := this.Pool.Get()
	defer pool.Close()
	keys, err = redis.Strings(pool.Do("KEYS", patternkey))
	return
}

//删除Redis 缓存键数据
func (this *RedisPool) Delete(key string) (err error) {
	pool := this.Pool.Get()
	defer pool.Close()
	_, err = pool.Do("DEL", key)
	return
}

//添加键值对
func (this *RedisPool) SetKeyForValue(key string, value interface{}) (err error) {
	if b, err := json.Marshal(value); err == nil {
		pool := this.Pool.Get()
		defer pool.Close()
		_, err = pool.Do("SET", key, b)
	}
	return
}

//添加过期键值对
func (this *RedisPool) SetExKeyForValue(key string, value interface{}, expire int) (err error) {
	if b, err := json.Marshal(value); err == nil {
		pool := this.Pool.Get()
		defer pool.Close()
		_, err = pool.Do("SET", key, b, "ex", expire)
	}
	return
}

//添加过期键值对
func (this *RedisPool) GetExKeyForValue(key string, value interface{}, expire int) (err error) {
	pool := this.Pool.Get()
	defer pool.Close()
	var valueStr string
	if valueStr, err = redis.String(pool.Do("GET", key)); err != nil {
		err = fmt.Errorf("GetExKeyForValue key:%s err:%v", key, err)
		return
	}
	err = json.Unmarshal([]byte(valueStr), value)
	if err != nil {
		err = fmt.Errorf("GetExKeyForValue  key:%s valueStr:%s err:%v", key, valueStr, err)
		return
	}
	_, err = pool.Do("EXPIRE", key, expire)
	return nil
}

//添加键值对
func (this *RedisPool) GetKeyForValue(key string, value interface{}) (err error) {
	pool := this.Pool.Get()
	defer pool.Close()
	var valueStr string
	if valueStr, err = redis.String(pool.Do("GET", key)); err != nil {
		return fmt.Errorf("GetKeyForValue  key:%s err:%v", key, err)
	}
	if err = json.Unmarshal([]byte(valueStr), value); err != nil {
		return fmt.Errorf("GetKeyForValue key:%s valueStr:%s err:%v", key, valueStr, err)
	}
	return
}

func (this *RedisPool) Lock(key string, outTime int) (err error) {
	//从连接池中娶一个con链接，pool可以自己定义。
	pool := this.Pool.Get()
	defer pool.Close()
	//这里需要redis.String包一下，才能返回redis.ErrNil
	_, err = redis.String(pool.Do("set", key, 1, "ex", outTime, "nx"))
	return
}
func (this *RedisPool) UnLock(key string) (err error) {
	pool := this.Pool.Get()
	defer pool.Close()
	_, err = pool.Do("del", key)
	return
}

//List---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------
//读取列表长度
func (this *RedisPool) GetListCount(key string) (count int, err error) {
	pool := this.Pool.Get()
	defer pool.Close()
	if count, err = redis.Int(pool.Do("LLEN", key)); err != nil {
		err = fmt.Errorf("GetListCount key:%s err:%v", key, err)
	}
	return
}

//读取列表
func (this *RedisPool) GetList(key string, valuetype reflect.Type) (values []interface{}, err error) {
	pool := this.Pool.Get()
	defer pool.Close()
	values = []interface{}{}
	var _values []interface{}
	if _values, err = redis.Values(pool.Do("LRANGE", key, 0, -1)); err != nil {
		err = fmt.Errorf("GetList key:%s err:%v", key, err)
		return
	}
	for _, value := range _values {
		v := reflect.New(valuetype.Elem()).Interface()
		if err = json.Unmarshal([]byte(value.([]uint8)), &v); err == nil {
			values = append(values, v)
		}
	}
	return
}

//随机从列表中获取指定数量的元素
func (this *RedisPool) GetListByRandom(key string, num int, valuetype reflect.Type) (values []interface{}, err error) {
	values = make([]interface{}, 0)
	pool := this.Pool.Get()
	defer pool.Close()
	var (
		count int
		value string
	)
	if count, err = redis.Int(pool.Do("LLEN", key)); err != nil || count == 0 {
		return
	}
	if num > count {
		num = count
	}
	r := rand.New(rand.NewSource(time.Now().Unix()))
	n := 0
	for _, i := range r.Perm(count) {
		n++
		if value, err = redis.String(pool.Do("LINDEX", key, i)); err != nil {
			v := reflect.New(valuetype.Elem()).Interface()
			if err = json.Unmarshal([]byte(value), &v); err != nil {
				values = append(values, v)
				if n >= num {
					break
				}
			}
		}
	}
	return
}

//通过索引获取列表中的元素
func (this *RedisPool) GetListByLIndex(key string, index int32, value interface{}) (err error) {
	pool := this.Pool.Get()
	defer pool.Close()
	var valueStr string
	if valueStr, err = redis.String(pool.Do("LINDEX", key, index)); err != nil {
		err = fmt.Errorf("GetListByLIndex key:%s index:%d err:%v", key, index, err)
		return
	}
	if err = json.Unmarshal([]byte(valueStr), value); err != nil {
		err = fmt.Errorf("GetListByLIndex key:%s index:%d err:%v", key, index, err)
	}
	return
}

//通过索引来设置元素的值
func (this *RedisPool) SetListByLIndex(key string, index int32, value interface{}) (err error) {
	pool := this.Pool.Get()
	defer pool.Close()
	var _value []byte
	if _value, err = json.Marshal(value); err == nil {
		if _, err = pool.Do("LSET", key, index, _value); err != nil {
			err = fmt.Errorf("SetListByLIndex key:%s index:%d err:%v", key, index, err)
		}
	} else {
		err = fmt.Errorf("SetListByLIndex key:%s index:%d err:%v", key, index, err)
	}
	return
}

//将一个或多个值插入到列表头部(最左边)
func (this *RedisPool) SetListByLPush(key string, value []interface{}) (err error) {
	pool := this.Pool.Get()
	defer pool.Close()
	var _value []byte
	Values := []interface{}{}
	Values = append(Values, key)
	for _, v := range value {
		if _value, err = json.Marshal(v); err == nil {
			Values = append(Values, string(_value))
		}
	}
	if _, err = pool.Do("LPUSH", Values...); err != nil {
		err = fmt.Errorf("SetListByLPush key:%s err:%v", key, err)
	}
	return
}

//将一个或多个值插入到列表的尾部(最右边)
func (this *RedisPool) SetListByRPush(key string, value []interface{}) (err error) {
	pool := this.Pool.Get()
	defer pool.Close()
	var _value []byte
	Values := []interface{}{}
	Values = append(Values, key)
	for _, v := range value {
		if _value, err = json.Marshal(v); err == nil {
			Values = append(Values, string(_value))
		}
	}
	if _, err = pool.Do("RPUSH", Values...); err != nil {
		err = fmt.Errorf("SetListByRPush   key:%s err:%v", key, err)
	}
	return
}

//读取返回目标区域的数据
func (this *RedisPool) GetListByLrange(key string, start, end int32, valuetype reflect.Type) (values []interface{}, err error) {
	pool := this.Pool.Get()
	defer pool.Close()
	var valuesStr []string
	if valuesStr, err = redis.Strings(pool.Do("LRANGE", key, start, end)); err != nil {
		err = fmt.Errorf("GetListByLrange  key:%s err:%v", key, err)
		return
	}
	values = make([]interface{}, 0)
	for _, value := range valuesStr {
		v := reflect.New(valuetype.Elem()).Interface()
		err = json.Unmarshal([]byte(value), &v)
		if err == nil {
			values = append(values, v)
		}
	}
	return
}

//移除列表的最后一个元素，返回值为移除的元素
func (this *RedisPool) GetListByRPop(key string, value interface{}) (err error) {
	pool := this.Pool.Get()
	defer pool.Close()
	var valueStr string
	if valueStr, err = redis.String(pool.Do("RPOP", key)); err != nil {
		err = fmt.Errorf("GetListByRPop  key:%s err:%v", key, err)
		return
	}
	if err = json.Unmarshal([]byte(valueStr), value); err != nil {
		err = fmt.Errorf("GetListByRPop  key:%s err:%v", key, err)
	}
	return
}

//移除并返回列表的第一个元素
func (this *RedisPool) GetListByLPop(key string, value interface{}) (err error) {
	pool := this.Pool.Get()
	defer pool.Close()
	var valueStr string
	if valueStr, err = redis.String(pool.Do("LPOP", key)); err != nil {
		err = fmt.Errorf("GetListByLPop  key:%s err:%v", key, err)
		return
	}
	if err = json.Unmarshal([]byte(valueStr), value); err != nil {
		err = fmt.Errorf("GetListByLPop  key:%s err:%v", key, err)
	}
	return
}

//移除列表中与参数 VALUE 相等的元素。
func (this *RedisPool) RemoveListByValue(key string, value interface{}) (err error) {
	pool := this.Pool.Get()
	defer pool.Close()
	var (
		valueStr string
		_value   []byte
	)
	if _value, err = json.Marshal(value); err != nil {
		err = fmt.Errorf("RemoveListByValue  key:%s err:%v ", key, err)
		return
	}
	valueStr = string(_value)
	if _, err = pool.Do("LREM", key, 0, valueStr); err != nil {
		err = fmt.Errorf("RemoveListByValue  key:%s err:%v ", key, err)
	}
	return
}

//Map--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------
//判断键是否存在 Map中
func (this *RedisPool) ContainsByMap(key string, field string) (iskeep bool, err error) {
	pool := this.Pool.Get()
	defer pool.Close()
	if iskeep, err = redis.Bool(pool.Do("Hexists", key, field)); err != nil {
		err = fmt.Errorf("ContainsByMap key:%s field:%s err:%v", key, field, err)
	}
	return
}

//添加键值 哈希表
func (this *RedisPool) SetKey_Map(key string, value map[string]interface{}) (err error) {
	pool := this.Pool.Get()
	defer pool.Close()
	var (
		_value []byte
	)
	Values := []interface{}{}
	Values = append(Values, key)
	for k, v := range value {
		if _value, err = json.Marshal(v); err == nil {
			Values = append(Values, k, string(_value))
		}
	}
	if _, err = pool.Do("HMSET", Values...); err != nil {
		err = fmt.Errorf("SetKey_Map key:%s err:%v", key, err)
	}
	return
}

//读取键值 哈希表
func (this *RedisPool) GetKey_Map(key string, valuetype reflect.Type) (value map[string]interface{}, err error) {
	pool := this.Pool.Get()
	defer pool.Close()
	var (
		_value map[string]string
		_v     interface{}
	)
	value = make(map[string]interface{})
	if _value, err = redis.StringMap(pool.Do("HGETALL", key)); err != nil {
		err = fmt.Errorf("GetKey_Map key:%s err:%v", key, err)
		return
	}
	for k, v := range _value {
		_v = reflect.New(valuetype.Elem()).Interface()
		if err = json.Unmarshal([]byte(v), &v); err == nil {
			value[k] = _v
		}
	}
	return
}

func (this *RedisPool) GetKey_MapByLen(key string) (count int, err error) {
	pool := this.Pool.Get()
	defer pool.Close()
	count, err = redis.Int(pool.Do("HLEN", key))
	return
}

func (this *RedisPool) GetKey_MapByKeys(key string) (keys []string, err error) {
	pool := this.Pool.Get()
	defer pool.Close()
	if keys, err = redis.Strings(pool.Do("HKEYS", key)); err != nil {
		err = fmt.Errorf("GetKey_MapByLen  key:%s err:%v", key, err)
	}
	return
}

func (this *RedisPool) GetKey_MapByValues(key string, valuetype reflect.Type) (Values []interface{}, err error) {
	pool := this.Pool.Get()
	defer pool.Close()
	var (
		_value map[string]string
	)
	Values = make([]interface{}, 0)
	if _value, err = redis.StringMap(pool.Do("HGETALL", key)); err != nil {
		err = fmt.Errorf("GetKey_MapByValues  key:%s err:%v", key, err)
		return
	}
	for _, value := range _value {
		v := reflect.New(valuetype.Elem()).Interface()
		if err = json.Unmarshal([]byte(value), &v); err == nil {
			Values = append(Values, v)
		}
	}
	return
}

//读取键值 哈希表key 值
func (this *RedisPool) GetKey_MapByKey(key string, field string, value interface{}) (err error) {
	pool := this.Pool.Get()
	defer pool.Close()
	var (
		_value string
	)
	if _value, err = redis.String(pool.Do("HGET", key, field)); err != nil {
		err = fmt.Errorf("GetKey_MapByKey  key:%s field:%s err:%v", key, field, err)
		return
	}
	if err = json.Unmarshal([]byte(_value), value); err != nil {
		err = fmt.Errorf("GetKey_MapByKey  key:%s field:%s value:%s err:%v", key, field, _value, err)
	}
	return
}

//删除键值 哈希表 字段
func (this *RedisPool) DelKey_MapKey(key string, field string) (err error) {
	pool := this.Pool.Get()
	defer pool.Close()
	if _, err = pool.Do("HDEL", key, field); err != nil {
		err = fmt.Errorf("DelKey_MapKey key:%s field:%s err:%v", key, field, err)
	}
	return
}

//获取所有符合给定模式 pattern 的 key 的数据
func (this *RedisPool) Get_PatternKeys(patternkey string, valuetype reflect.Type) (Values []interface{}, err error) {
	pool := this.Pool.Get()
	defer pool.Close()
	var (
		keys   []string
		_value string
	)
	Values = []interface{}{}
	if keys, err = redis.Strings(pool.Do("KEYS", patternkey)); err != nil {
		err = fmt.Errorf("Get_PatternKeys patternkey:%s err:%v", patternkey, err)
		return
	}
	for _, key := range keys {
		if _value, err = redis.String(pool.Do("GET", key)); err == nil {
			value := reflect.New(valuetype.Elem()).Interface()
			if err = json.Unmarshal([]byte(_value), &value); err == nil {
				Values = append(Values, value)
			}
		}
	}
	return
}
