package pools

import (
	"reflect"
)

/*
系统描述:泛型对象池
*/

var (
	defobjpools  *PoolManager = NewPoolMessager()
	deftypepools *typePools   = NewTypePool()
)

func Add[T any](name string, newFunc func() T) {
	defobjpools.pools[name] = NewObjectPool(newFunc)
}

func Get[T any](name string) (v T) {
	return get[T](defobjpools, name)
}

func Put[T any](name string, v T) {
	put[T](defobjpools, name, v)
}

// 注册类型
func InitTypes(ts ...reflect.Type) {
	for _, t := range ts {
		deftypepools.Init(t)
	}
}

func GetForType(t reflect.Type) interface{} {
	return deftypepools.Get(t)
}

func PutForType(t reflect.Type, v interface{}) {
	deftypepools.Put(t, v)
}
