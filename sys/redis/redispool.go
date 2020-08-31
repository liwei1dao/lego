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
func (this *RedisPool) ContainsKey(key string) bool {
	pool := this.Pool.Get()
	defer pool.Close()
	iskeep, err := redis.Bool(pool.Do("EXISTS", key))
	if err != nil {
		log.Errorf("检测缓存键失败 key = %s err = %s", key, err.Error())
	}
	return iskeep
}

//查询缓存keys
func (this *RedisPool) QueryKeys(patternkey string) (keys []string) {
	pool := this.Pool.Get()
	defer pool.Close()
	keys, _ = redis.Strings(pool.Do("KEYS", patternkey))
	return
}

//删除Redis 缓存键数据
func (this *RedisPool) Delete(key string) {
	pool := this.Pool.Get()
	defer pool.Close()
	_, err := pool.Do("DEL", key)
	if err != nil {
		log.Errorf("删除缓存数据失败 key = %s", key)
	}
}

//添加键值对
func (this *RedisPool) SetKeyForValue(key string, value interface{}) {
	if b, err := json.Marshal(value); err == nil {
		pool := this.Pool.Get()
		defer pool.Close()
		_, err := pool.Do("SET", key, b)
		if err != nil {
			log.Errorf("SetKey_String 设置缓存数据失败 key = %s", key)
		}
	}
}

//添加过期键值对
func (this *RedisPool) SetExKeyForValue(key string, value interface{}, expire int) {
	if b, err := json.Marshal(value); err == nil {
		pool := this.Pool.Get()
		defer pool.Close()
		_, err := pool.Do("SET", key, b, "ex", expire)
		if err != nil {
			log.Errorf("SetKey_String 设置缓存数据失败 key = %s", key)
		}
	}
}

//添加过期键值对
func (this *RedisPool) GetExKeyForValue(key string, value interface{}, expire int) (err error) {
	pool := this.Pool.Get()
	defer pool.Close()
	v, err := redis.String(pool.Do("GET", key, "ex", expire))
	if err != nil {
		return fmt.Errorf("GetKey_String 获取缓存数据失败 key = %s err:%s", key, err.Error())
	}
	err = json.Unmarshal([]byte(v), value)
	if err != nil {
		return fmt.Errorf("GetKey_String 解析缓存数据失败 key = %s Value=%s err:%s", key, v, err.Error())
	}
	return nil
}

//添加键值对
func (this *RedisPool) GetKeyForValue(key string, value interface{}) (err error) {
	pool := this.Pool.Get()
	defer pool.Close()
	v, err := redis.String(pool.Do("GET", key))
	if err != nil {
		return fmt.Errorf("GetKey_String 获取缓存数据失败 key = %s err:%s", key, err.Error())
	}
	err = json.Unmarshal([]byte(v), value)
	if err != nil {
		return fmt.Errorf("GetKey_String 解析缓存数据失败 key = %s Value=%s err:%s", key, v, err.Error())
	}
	return nil
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
	if err != nil {
		return
	}
	return
}

//List---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------
//读取列表长度
func (this *RedisPool) GetListCount(key string) (count int, err error) {
	pool := this.Pool.Get()
	defer pool.Close()
	count, err = redis.Int(pool.Do("LLEN", key))
	if err != nil {
		return count, fmt.Errorf("GetListCount 执行异常 key:%s err:%s", key, err.Error())
	}
	return count, nil
}

//读取列表
func (this *RedisPool) GetList(key string, valuetype reflect.Type) (values []interface{}) {
	pool := this.Pool.Get()
	defer pool.Close()
	values = []interface{}{}
	movies, err := redis.Values(pool.Do("LRANGE", key, 0, -1))
	if err != nil {
		log.Errorf("SetKey_String 读取缓存列表数据失败 key = %s", key)
		return
	}
	for _, value := range movies {
		v := reflect.New(valuetype.Elem()).Interface()
		err := json.Unmarshal([]byte(value.([]uint8)), &v)
		if err == nil {
			values = append(values, v)
		}
	}
	return values
}

//随机从列表中获取指定数量的元素
func (this *RedisPool) GetListByRandom(key string, num int, valuetype reflect.Type) (values []interface{}) {
	values = make([]interface{}, 0)
	pool := this.Pool.Get()
	defer pool.Close()
	count, err := redis.Int(pool.Do("LLEN", key))
	if err != nil || count == 0 {
		return
	}
	if num > count {
		num = count
	}
	r := rand.New(rand.NewSource(time.Now().Unix()))
	n := 0
	for _, i := range r.Perm(count) {
		n++
		if data, err := redis.String(pool.Do("LINDEX", key, i)); err != nil {
			v := reflect.New(valuetype.Elem()).Interface()
			if err = json.Unmarshal([]byte(data), &v); err != nil {
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
	data, err := redis.String(pool.Do("LINDEX", key, index))
	if err != nil {
		return fmt.Errorf("GetListByLIndex 执行异常 key:%s index:%d err:%s", key, index, err.Error())
	}
	err = json.Unmarshal([]byte(data), value)
	if err != nil {
		return fmt.Errorf("GetListByLIndex 执行异常 key:%s index:%d err:%s", key, index, err.Error())
	}
	return nil
}

//通过索引来设置元素的值
func (this *RedisPool) SetListByLIndex(key string, index int32, value interface{}) (err error) {
	pool := this.Pool.Get()
	defer pool.Close()
	if data, err := json.Marshal(value); err == nil {
		_, err := pool.Do("LSET", key, index, data)
		if err != nil {
			return fmt.Errorf("GetListByLIndex 执行异常 key:%s index:%d err:%s", key, index, err.Error())
		}
	} else {
		return fmt.Errorf("GetListByLIndex 执行异常 key:%s index:%d err:%s", key, index, err.Error())
	}
	return
}

//将一个或多个值插入到列表头部(最左边)
func (this *RedisPool) SetListByLPush(key string, value []interface{}) {
	pool := this.Pool.Get()
	defer pool.Close()
	Values := []interface{}{}
	Values = append(Values, key)
	for _, v := range value {
		if b, err := json.Marshal(v); err == nil {
			Values = append(Values, string(b))
		}
	}
	_, err := pool.Do("LPUSH", Values...)
	if err != nil {
		log.Errorf("SetKey_String 设置缓存列表数据失败 err = %v key = %s", err, key)
	}
}

//将一个或多个值插入到列表的尾部(最右边)
func (this *RedisPool) SetListByRPush(key string, value []interface{}) {
	pool := this.Pool.Get()
	defer pool.Close()
	Values := []interface{}{}
	Values = append(Values, key)
	for _, v := range value {
		if b, err := json.Marshal(v); err == nil {
			Values = append(Values, string(b))
		}
	}
	_, err := pool.Do("RPUSH", Values...)
	if err != nil {
		log.Errorf("SetKey_String 设置缓存列表数据失败 err = %v key = %s", err, key)
	}
}

//读取返回目标区域的数据
func (this *RedisPool) GetListByLrange(key string, start, end int32, valuetype reflect.Type) (values []interface{}) {
	pool := this.Pool.Get()
	defer pool.Close()
	datas, err := redis.Strings(pool.Do("LRANGE", key, start, end))
	values = make([]interface{}, 0)
	if err != nil {
		log.Errorf("GetListByLrange 执行异常 key:%s err:%s", key, err.Error())
		return
	}
	for _, value := range datas {
		v := reflect.New(valuetype.Elem()).Interface()
		err := json.Unmarshal([]byte(value), &v)
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
	data, err := redis.String(pool.Do("RPOP", key))
	if err != nil {
		return fmt.Errorf("GetListByRPop  key:%s err:%s", key, err.Error())
	}
	err = json.Unmarshal([]byte(data), value)
	if err != nil {
		return fmt.Errorf("GetListByRPop  key:%s err:%s", key, err.Error())
	}
	return nil
}

//移除并返回列表的第一个元素
func (this *RedisPool) GetListByLPop(key string, value interface{}) (err error) {
	pool := this.Pool.Get()
	defer pool.Close()
	data, err := redis.String(pool.Do("LPOP", key))
	if err != nil {
		return fmt.Errorf("GetListByLPop  key:%s err:%s", key, err.Error())
	}
	err = json.Unmarshal([]byte(data), value)
	if err != nil {
		return fmt.Errorf("GetListByLPop  key:%s err:%s", key, err.Error())
	}
	return nil
}

//移除列表中与参数 VALUE 相等的元素。
func (this *RedisPool) RemoveListByValue(key string, value interface{}) (err error) {
	pool := this.Pool.Get()
	defer pool.Close()
	valueStr := ""
	if b, err := json.Marshal(value); err != nil {
		log.Errorf("RemoveListByValue  key:%s err:%s ", key, err.Error())
		return err
	} else {
		valueStr = string(b)
	}
	_, err = pool.Do("LREM", key, 0, valueStr)
	if err != nil {
		log.Errorf("RemoveListByValue  key:%s err:%s ", key, err.Error())
		return err
	}
	return
}

//Map--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------
//判断键是否存在 Map中
func (this *RedisPool) ContainsByMap(key string, _FieldKey string) bool {
	pool := this.Pool.Get()
	defer pool.Close()
	iskeep, err := redis.Bool(pool.Do("Hexists", key, _FieldKey))
	if err != nil {
		log.Errorf("检测缓存键失败 key = %s field = %s", key, _FieldKey)
	}
	return iskeep
}

//添加键值 哈希表
func (this *RedisPool) SetKey_Map(key string, value map[string]interface{}) {
	pool := this.Pool.Get()
	defer pool.Close()
	Values := []interface{}{}
	Values = append(Values, key)
	for k, v := range value {
		if b, err := json.Marshal(v); err == nil {
			Values = append(Values, k, string(b))
		}
	}
	_, err := pool.Do("HSET", Values...)
	if err != nil {
		log.Errorf("设置缓存SetKey_Map失败 key = %s err=%s", key, err.Error())
	}
}

//读取键值 哈希表
func (this *RedisPool) GetKey_Map(key string, valuetype reflect.Type) (Value map[string]interface{}) {
	pool := this.Pool.Get()
	defer pool.Close()
	Value = make(map[string]interface{})
	movies, err := redis.StringMap(pool.Do("HGETALL", key))
	if err != nil {
		log.Errorf("SetKey_String 读取缓存哈希表数据失败 key = %s", key)
		return
	}
	for k, value := range movies {
		v := reflect.New(valuetype.Elem()).Interface()
		err := json.Unmarshal([]byte(value), &v)
		if err == nil {
			Value[k] = v
		}
	}
	return Value
}

func (this *RedisPool) GetKey_MapByLen(key string) (count int, err error) {
	pool := this.Pool.Get()
	defer pool.Close()
	count, err = redis.Int(pool.Do("HLEN", key))
	if err != nil {
		return 0, err
	}
	return count, err
}

func (this *RedisPool) GetKey_MapByKeys(key string) (keys []string, err error) {
	pool := this.Pool.Get()
	defer pool.Close()
	keys, err = redis.Strings(pool.Do("HKEYS", key))
	if err != nil {
		log.Errorf("SetKey_String 读取缓存哈希表数据失败 key = %s", key)
		return []string{}, err
	}
	return
}

func (this *RedisPool) GetKey_MapByValues(key string, valuetype reflect.Type) (Values []interface{}) {
	pool := this.Pool.Get()
	defer pool.Close()
	Values = make([]interface{}, 0)
	movies, err := redis.StringMap(pool.Do("HGETALL", key))
	if err != nil {
		log.Errorf("SetKey_String 读取缓存哈希表数据失败 key = %s", key)
		return
	}
	for _, value := range movies {
		v := reflect.New(valuetype.Elem()).Interface()
		err := json.Unmarshal([]byte(value), &v)
		if err == nil {
			Values = append(Values, v)
		}
	}
	return Values
}

//读取键值 哈希表key 值
func (this *RedisPool) GetKey_MapByKey(key string, fieldkey string, value interface{}) (err error) {
	if !this.ContainsByMap(key, fieldkey) {
		return fmt.Errorf("GetKey_MapByKey 读取缓存哈希表数据失败 不存在的 key = %s", key)
	}
	pool := this.Pool.Get()
	defer pool.Close()
	movies, err := redis.String(pool.Do("HGET", key, fieldkey))
	if err != nil {
		return fmt.Errorf("GetKey_MapByKey 读取缓存哈希表数据失败 key = %s", key)
	}
	err = json.Unmarshal([]byte(movies), value)
	if err != nil {
		return fmt.Errorf("读取Redis Map【%s】【%s】错误 %s", key, fieldkey, err)
	}
	return nil
}

//删除键值 哈希表 字段
func (this *RedisPool) DelKey_MapKey(key string, _FieldKey string) {
	pool := this.Pool.Get()
	defer pool.Close()
	_, err := pool.Do("HDEL", key, _FieldKey)
	if err != nil {
		log.Errorf("删除键值哈希表字段失败 key=%s FieldKey=%s", key, _FieldKey)
	}
}

//获取所有符合给定模式 pattern 的 key 的数据
func (this *RedisPool) Get_PatternKeys(_Patternkey string, valuetype reflect.Type) (Values []interface{}) {
	pool := this.Pool.Get()
	defer pool.Close()
	Values = []interface{}{}
	keys, err := redis.Strings(pool.Do("KEYS", _Patternkey))
	if err != nil {
		log.Errorf("SetKey_String 读取缓存列表数据失败 key = %s", _Patternkey)
		return
	}
	for _, key := range keys {
		if v, err := redis.String(pool.Do("GET", key)); err == nil {
			value := reflect.New(valuetype.Elem()).Interface()
			err = json.Unmarshal([]byte(v), &value)
			if err == nil {
				Values = append(Values, value)
			}
		}
	}
	return Values
}
