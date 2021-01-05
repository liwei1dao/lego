package pprof

type (
	Ipprof interface {
	}
)

var (
	defsys Ipprof
)

func OnInit(config map[string]interface{}, option ...Option) (err error) {
	if defsys, err = newSys(newOptions(config, option...)); err == nil {
	}
	return
}

func NewSys(option ...Option) (sys Ipprof, err error) {
	if sys, err = newSys(newOptionsByOption(option...)); err == nil {
	}
	return
}
