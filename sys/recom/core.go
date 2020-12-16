package recom

/*
推荐系统 协同过滤算法
*/
type (
	IRecom interface {
		Fit()
		Wait()
		RecommendItems(uId uint32, howmany int) (itemIds []uint32)
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

func Fit() {
	defsys.Fit()
}

func Wait() {
	defsys.Wait()
}

func RecommendItems(uId uint32, howmany int) (itemIds []uint32) {
	return defsys.RecommendItems(uId, howmany)
}
