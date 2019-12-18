package timer

type ITimer interface {
	Add(inteval uint32, handler func(string, ...interface{}), args ...interface{}) string
	Remove(key string)
}
