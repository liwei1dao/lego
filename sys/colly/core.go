/// 爬虫系统
package colly

type (
	Isys interface {
	}
)

var (
	defsys Isys
)

func OnInit(config map[string]interface{}, option ...Option) (err error) {
	if defsys, err = newSys(newOptions(config, option...)); err == nil {
	}
	return
}

func NewSys(option ...Option) (sys Isys, err error) {
	if sys, err = newSys(newOptionsByOption(option...)); err == nil {
	}
	return
}
