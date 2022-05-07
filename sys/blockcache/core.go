package blockcache

/*
限容堵塞缓冲池系统
设置最大内存大小 当缓存区为存满是直接写入 写满后进入堵塞状态 等待缓存区释放
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
