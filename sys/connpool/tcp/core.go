package tcp

import "github.com/liwei1dao/lego/sys/connpool"

var (
	defsys connpool.IConnPool
)

func OnInit(config map[string]interface{}, opts ...Option) (err error) {
	var option *Options
	if option, err = newOptions(config, opts...); err != nil {
		return
	}
	if defsys, err = newSys(option); err == nil {
		Start()
	}
	return
}

func NewSys(opts ...Option) (sys connpool.IConnPool, err error) {
	var option *Options
	if option, err = newOptionsByOption(opts...); err != nil {
		return
	}
	if sys, err = newSys(option); err == nil {
		sys.Start()
	}
	return
}

func Start() {
	defsys.Start()
}

func Close() {
	defsys.Close()
}
