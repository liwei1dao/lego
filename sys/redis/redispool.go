package redis

import (
	"encoding/json"
	"fmt"
	"lego/sys/log"
	"reflect"

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
func (this *RedisPool) ContainsKey(_Key string) bool {
	pool := this.Pool.Get()
	defer pool.Close()
	iskeep, err := redis.Bool(pool.Do("EXISTS", _Key))
	if err != nil {
		log.Errorf("检测缓存键失败 key = %s err = %s", _Key, err.Error())
	}
	return iskeep
}

//删除Redis 缓存键数据
func (this *RedisPool) Delete(_Key string) {
	pool := this.Pool.Get()
	defer pool.Close()
	_, err := pool.Do("DEL", _Key)
	if err != nil {
		log.Errorf("删除缓存数据失败 key = %s", _Key)
	}
}

//添加键值对
func (this *RedisPool) SetKey_Value(_Key string, _Value interface{}) {
	if b, err := json.Marshal(_Value); err == nil {
		pool := this.Pool.Get()
		defer pool.Close()
		_, err := pool.Do("SET", _Key, b)
		if err != nil {
			log.Errorf("SetKey_String 设置缓存数据失败 key = %s", _Key)
		}
	}
}

//添加过期键值对
func (this *RedisPool) SetExKey_Value(_Key string, _Value interface{}, expire int) {
	if b, err := json.Marshal(_Value); err == nil {
		pool := this.Pool.Get()
		defer pool.Close()
		_, err := pool.Do("SET", _Key, b, "ex", expire)
		if err != nil {
			log.Errorf("SetKey_String 设置缓存数据失败 key = %s", _Key)
		}
	}
}

//添加键值对
func (this *RedisPool) GetKey_Value(_Key string, _Value interface{}) (err error) {
	pool := this.Pool.Get()
	defer pool.Close()
	v, err := redis.String(pool.Do("GET", _Key))
	if err != nil {
		return fmt.Errorf("GetKey_String 获取缓存数据失败 key = %s err:%s", _Key, err.Error())
	}
	err = json.Unmarshal([]byte(v), _Value)
	if err != nil {
		return fmt.Errorf("GetKey_String 解析缓存数据失败 key = %s Value=%s err:%s", _Key, v, err.Error())
	}
	return nil
}

func (this *RedisPool) Lock(_Key string, outTime int) (err error) {
	//从连接池中娶一个con链接，pool可以自己定义。
	pool := this.Pool.Get()
	defer pool.Close()
	//这里需要redis.String包一下，才能返回redis.ErrNil
	_, err = redis.String(pool.Do("set", _Key, 1, "ex", outTime, "nx"))
	return
}
func (this *RedisPool) UnLock(_Key string) (err error) {
	pool := this.Pool.Get()
	defer pool.Close()
	_, err = pool.Do("del", _Key)
	if err != nil {
		return
	}
	return
}

//添加键值列表
func (this *RedisPool) SetKey_List(_Key string, _Value []interface{}) {
	pool := this.Pool.Get()
	defer pool.Close()
	Values := []interface{}{}
	Values = append(Values, _Key)
	for _, v := range _Value {
		if b, err := json.Marshal(v); err == nil {
			Values = append(Values, string(b))
		}
	}
	_, err := pool.Do("LPUSH", Values...)
	if err != nil {
		log.Errorf("SetKey_String 设置缓存列表数据失败 err = %v key = %s", err, _Key)
	}
}

//读取键值列表
func (this *RedisPool) GetKey_List(_Key string, _Vakue_type reflect.Type) (Value []interface{}) {
	pool := this.Pool.Get()
	defer pool.Close()
	Value = []interface{}{}
	movies, err := redis.Values(pool.Do("LRANGE", _Key, 0, -1))
	if err != nil {
		log.Errorf("SetKey_String 读取缓存列表数据失败 key = %s", _Key)
		return
	}
	for _, value := range movies {
		v := reflect.New(_Vakue_type.Elem()).Interface()
		err := json.Unmarshal([]byte(value.([]uint8)), &v)
		if err == nil {
			Value = append(Value, v)
		}
	}
	return Value
}

//判断键是否存在 Map中
func (this *RedisPool) ContainsKey_Map(_Key string, _FieldKey string) bool {
	pool := this.Pool.Get()
	defer pool.Close()
	iskeep, err := redis.Bool(pool.Do("Hexists", _Key, _FieldKey))
	if err != nil {
		log.Errorf("检测缓存键失败 key = %s field = %s", _Key, _FieldKey)
	}
	return iskeep
}

//添加键值 哈希表
func (this *RedisPool) SetKey_Map(_Key string, _Value map[string]interface{}) {
	pool := this.Pool.Get()
	defer pool.Close()
	Values := []interface{}{}
	Values = append(Values, _Key)
	for k, v := range _Value {
		if b, err := json.Marshal(v); err == nil {
			Values = append(Values, k, string(b))
		}
	}
	_, err := pool.Do("HSET", Values...)
	if err != nil {
		log.Errorf("设置缓存SetKey_Map失败 key = %s err=%s", _Key, err.Error())
	}
}

//添加键值 哈希表
func (this *RedisPool) SetExKey_Map(_Key string, _Value map[string]interface{}) {
	pool := this.Pool.Get()
	defer pool.Close()
	Values := []interface{}{}
	Values = append(Values, _Key)
	for k, v := range _Value {
		if b, err := json.Marshal(v); err == nil {
			Values = append(Values, k, string(b))
		}
	}
	_, err := pool.Do("HSET", Values...)
	if err != nil {
		log.Errorf("设置缓存SetKey_Map失败 key = %s err=%s", _Key, err.Error())
	}
}

//读取键值 哈希表
func (this *RedisPool) GetKey_Map(_Key string, _Vakue_type reflect.Type) (Value map[string]interface{}) {
	pool := this.Pool.Get()
	defer pool.Close()
	Value = make(map[string]interface{})
	movies, err := redis.StringMap(pool.Do("HGETALL", _Key))
	if err != nil {
		log.Errorf("SetKey_String 读取缓存哈希表数据失败 key = %s", _Key)
		return
	}
	for k, value := range movies {
		v := reflect.New(_Vakue_type.Elem()).Interface()
		err := json.Unmarshal([]byte(value), &v)
		if err == nil {
			Value[k] = v
		}
	}
	return Value
}

func (this *RedisPool) GetKey_MapByLen(_Key string) (count int, err error) {
	pool := this.Pool.Get()
	defer pool.Close()
	count, err = redis.Int(pool.Do("HLEN", _Key))
	if err != nil {
		return 0, err
	}
	return count, err
}

func (this *RedisPool) GetKey_MapByKeys(_Key string) (keys []string, err error) {
	pool := this.Pool.Get()
	defer pool.Close()
	keys, err = redis.Strings(pool.Do("HKEYS", _Key))
	if err != nil {
		log.Errorf("SetKey_String 读取缓存哈希表数据失败 key = %s", _Key)
		return []string{}, err
	}
	return
}

func (this *RedisPool) GetKey_MapByValues(_Key string, _Vakue_type reflect.Type) (Values []interface{}) {
	pool := this.Pool.Get()
	defer pool.Close()
	Values = make([]interface{}, 0)
	movies, err := redis.StringMap(pool.Do("HGETALL", _Key))
	if err != nil {
		log.Errorf("SetKey_String 读取缓存哈希表数据失败 key = %s", _Key)
		return
	}
	for _, value := range movies {
		v := reflect.New(_Vakue_type.Elem()).Interface()
		err := json.Unmarshal([]byte(value), &v)
		if err == nil {
			Values = append(Values, v)
		}
	}
	return Values
}

//读取键值 哈希表_key 值
func (this *RedisPool) GetKey_MapByKey(key string, fieldkey string, value interface{}) (err error) {
	if !this.ContainsKey_Map(key, fieldkey) {
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
func (this *RedisPool) DelKey_MapKey(_Key string, _FieldKey string) {
	pool := this.Pool.Get()
	defer pool.Close()
	_, err := pool.Do("HDEL", _Key, _FieldKey)
	if err != nil {
		log.Errorf("删除键值哈希表字段失败 key=%s FieldKey=%s", _Key, _FieldKey)
	}
}
