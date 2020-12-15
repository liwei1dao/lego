package recom

/*
推荐系统 协同过滤算法
*/
type (
	IRecom interface {
		
	}
)

var (
	defsys IRecom
)

func OnInit(config map[string]interface{}, option ...Option) (err error) {
	defsys, err = newSys(newOptions(config, option...))
	return
}

func NewSys(option ...Option) (sys IRecom, err error) {
	sys, err = newSys(newOptionsByOption(option...))
	return
}
