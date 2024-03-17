package sitemap

/*
系统描述:网站的seo优化系统，自动同步sitemap 文件
*/
type (
	ISys interface {
		AppendUrl(url *Url)
		GetUrls() (urls []*Url)
		SetUrls(urls []*Url)
		ToXml() ([]byte, error)
		Storage() (filename string, err error)
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

func AppendUrl(url *Url) {
	defsys.AppendUrl(url)
}

func GetUrls() (urls []*Url) {
	return defsys.GetUrls()
}

func SetUrls(urls []*Url) {
	defsys.SetUrls(urls)
}

func ToXml() ([]byte, error) {
	return defsys.ToXml()
}

func Storage() (filename string, err error) {
	return defsys.Storage()
}
