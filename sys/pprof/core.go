package pprof

/*
系统描述:pprof 性能检测驱动，排查内存以及cpu等指标异常问题，
*/

type (
	ISys interface {
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
