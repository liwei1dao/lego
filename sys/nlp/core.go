package nlp

import "github.com/go-ego/gse"

type (
	//系统接口
	ISys interface {
		//分词接口
		Cut(text string) []string
		//词性分析
		Pos(text string) []gse.SegPos
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

//分词接口
func Cut(text string) []string {
	return defsys.Cut(text)
}

//词性分析
func Pos(text string) []gse.SegPos {
	return defsys.Pos(text)
}
