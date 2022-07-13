package blockcache

/*
系统描述:限容堵塞缓冲池系统,设置最大内存大小 当缓存区为存满是直接写入 写满后进入堵塞状态 等待缓存区释放
*/

type (
	ISys interface {
		In() chan<- interface{}
		Out() <-chan interface{}
		Close()
	}
)

var (
	defsys ISys
)

func OnInit(config map[string]interface{}, opt ...Option) (err error) {
	var option *Options
	if option, err = newOptions(config, opt...); err != nil {
		return
	}
	defsys, err = newSys(option)
	return
}

func NewSys(opt ...Option) (sys ISys, err error) {
	var option *Options
	if option, err = newOptionsByOption(opt...); err != nil {
		return
	}
	sys, err = newSys(option)
	return
}

func In() chan<- interface{} {
	return defsys.In()
}

func Out() <-chan interface{} {
	return defsys.Out()
}

func Close() {
	defsys.Close()
}
