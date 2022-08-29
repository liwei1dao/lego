package wordfilter

import "io"

/*
系统:单词守卫
描述:DFA 算法实践敏感词过滤 敏感词 过滤 验证 替换
*/

type (
	ISys interface {
		//加载文件
		LoadWordDict(path string) error
		//加载网络数据
		LoadNetWordDict(url string) error
		///加载
		Load(rd io.Reader) error
		///添加敏感词
		AddWord(words ...string)
		///删除敏感词
		DelWord(words ...string)
		///过滤敏感词
		Filter(text string) string
		///和谐敏感词
		Replace(text string, repl rune) string
		///检测敏感词
		FindIn(text string) (bool, string)
		///找到所有匹配词
		FindAll(text string) []string
		///检测字符串是否合法
		Validate(text string) (bool, string)
		///去除空格等噪音
		RemoveNoise(text string) string
		///更新去噪模式
		UpdateNoisePattern(pattern string)
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

//加载文件
func LoadWordDict(path string) error {
	return defsys.LoadWordDict(path)
}

//加载网络数据
func LoadNetWordDict(url string) error {
	return defsys.LoadNetWordDict(url)
}

///加载
func Load(rd io.Reader) error {
	return defsys.Load(rd)
}

///添加敏感词
func AddWord(words ...string) {
	defsys.AddWord(words...)
}

///删除敏感词
func DelWord(words ...string) {
	defsys.DelWord(words...)
}

///过滤敏感词
func Filter(text string) string {
	return defsys.Filter(text)
}

///和谐敏感词
func Replace(text string, repl rune) string {
	return defsys.Replace(text, repl)
}

///检测敏感词
func FindIn(text string) (bool, string) {
	return defsys.FindIn(text)
}

///找到所有匹配词
func FindAll(text string) []string {
	return defsys.FindAll(text)
}

///检测字符串是否合法
func Validate(text string) (bool, string) {
	return defsys.Validate(text)
}

///去除空格等噪音
func RemoveNoise(text string) string {
	return defsys.RemoveNoise(text)
}

///更新去噪模式
func UpdateNoisePattern(pattern string) {
	defsys.UpdateNoisePattern(pattern)
}
