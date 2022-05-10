package recom

/*
推荐系统 协同过滤算法
*/
type (
	RecomModel uint8
	ISys       interface {
		Fit()
		Wait()
		RecommendItems(uId uint32, howmany int) (itemIds []uint32)
	}
)

const (
	BPRModel RecomModel = iota
	SVD
)

var (
	defsys ISys
)

func OnInit(config map[string]interface{}, option ...Option) (err error) {
	defsys, err = newSys(newOptions(config, option...))
	return
}

func NewSys(option ...Option) (sys ISys, err error) {
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
