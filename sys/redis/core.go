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

func OnInit(config map[string]interface{}) (err error) {
	defsys, err = newsys(newOptionsByConfig(config))
	return
}

func NewRedisSys(opt ...Option) (sys IRedisFactory, err error) {
	sys, err = newsys(newOptionsByOption(opt...))
	return
}

func GetPool() *RedisPool {
	return defsys.GetPool()
}

func CloseAllPool() {
	defsys.CloseAllPool()
}
