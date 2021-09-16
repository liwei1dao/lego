/// 爬虫系统
package colly

import "github.com/gocolly/colly/v2"

type (
	ISys interface {
		OnResponse(f colly.ResponseCallback)
		OnScraped(f colly.ScrapedCallback)
		Post(url string, requestData map[string]string) error
		PostRaw(url string, requestData []byte) error
		Visit(url string) error
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

func OnResponse(f colly.ResponseCallback) {
	defsys.OnResponse(f)
}
func OnScraped(f colly.ScrapedCallback) {
	defsys.OnScraped(f)
}

func Post(url string, requestData map[string]string) error {
	return defsys.Post(url, requestData)
}

func PostRaw(url string, requestData []byte) error {
	return defsys.PostRaw(url, requestData)
}
func Visit(url string) error {
	return defsys.Visit(url)
}
