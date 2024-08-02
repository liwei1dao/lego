package pools

import (
	"sync"

	"github.com/liwei1dao/lego/sys/log"
)

// 泛型对象池结构体
type ObjectPool[T any] struct {
	pool *sync.Pool
}

// 创建一个新的泛型对象池
func NewObjectPool[T any](newFunc func() T) *ObjectPool[T] {
	return &ObjectPool[T]{
		pool: &sync.Pool{
			New: func() interface{} {
				return newFunc()
			},
		},
	}
}

// 从对象池获取对象
func (p *ObjectPool[T]) Get() T {
	return p.pool.Get().(T)
}

// 将对象放回对象池
func (p *ObjectPool[T]) Put(obj T) {
	p.pool.Put(obj)
}
func NewPoolMessager() *PoolManager {
	return &PoolManager{
		pools: make(map[string]interface{}),
	}
}

type PoolManager struct {
	pools map[string]interface{}
}

// 添加一个对象池到集合的泛型函数
func addPool[T any](pm *PoolManager, name string, newFunc func() T) {
	pm.pools[name] = NewObjectPool(newFunc)
}

// 获取对象
func get[T any](pm *PoolManager, name string) (v T) {
	pool, exists := pm.pools[name]
	if !exists {
		log.Panic("no found pool", log.Field{Key: "name", Value: name})
		return
	}
	typedPool := pool.(*ObjectPool[T])
	return typedPool.Get()
}

// 放回对象
func put[T any](pm *PoolManager, name string, v T) {
	pool, exists := pm.pools[name]
	if !exists {
		log.Panic("no found pool", log.Field{Key: "name", Value: name})
		return
	}
	typedPool := pool.(*ObjectPool[T])
	typedPool.Put(v)
}
