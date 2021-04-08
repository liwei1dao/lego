package redis

type (
	IRedis interface {
		
	}
)

var defsys IRedis

func OnInit(config map[string]interface{}, option ...Option) (err error) {
	defsys, err = newSys(newOptions(config, option...))
	return
}

func NewSys(option ...Option) (sys IRedis, err error) {
	sys, err = newSys(newOptionsByOption(option...))
	return
}
