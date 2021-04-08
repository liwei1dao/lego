package redis

type (
	IRedisFactory interface {
		GetPool() *RedisPool
		CloseAllPool()
	}
)

var (
	defsys IRedisFactory
)

func OnInit(config map[string]interface{}, option ...Option) (err error) {
	defsys, err = newsys(newOptions(config, option...))
	return
}

func NewSys(option ...Option) (sys IRedisFactory, err error) {
	sys, err = newsys(newOptionsByOption(option...))
	return
}

func GetPool() *RedisPool {
	return defsys.GetPool()
}

func CloseAllPool() {
	defsys.CloseAllPool()
}
