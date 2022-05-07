package blockcache

/*
堵塞缓存系统 缓存区写满后自动堵塞 读取是无数据自动堵塞
*/

type (
	ISys interface {
		In() chan<- interface{}
		Out() <-chan interface{}
	}
)

var (
	defsys ISys
)

func OnInit(config map[string]interface{}, option ...Option) (err error) {
	if defsys, err = newSys(newOptions(config, option...)); err == nil {

	}
	return
}

func NewSys(option ...Option) (sys ISys, err error) {
	if sys, err = newSys(newOptionsByOption(option...)); err == nil {

	}
	return
}

func In() chan<- interface{} {
	return defsys.In()
}

func Out() <-chan interface{} {
	return defsys.Out()
}


