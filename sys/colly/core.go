/// 爬虫系统
package colly

import (
	"github.com/gocolly/colly/v2"
)

type (
	ISys interface {
		OnResponse(f colly.ResponseCallback)
		OnScraped(f colly.ScrapedCallback)
		OnError(f colly.ErrorCallback)
		Post(url string, requestData map[string]string) error
		PostRaw(url string, requestData []byte) error
		Wait()
		Visit(url string) error
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

func OnResponse(f colly.ResponseCallback) {
	defsys.OnResponse(f)
}
func OnScraped(f colly.ScrapedCallback) {
	defsys.OnScraped(f)
}

func OnError(f colly.ErrorCallback) {
	defsys.OnError(f)
}

func Post(url string, requestData map[string]string) error {
	return defsys.Post(url, requestData)
}

func PostRaw(url string, requestData []byte) error {
	return defsys.PostRaw(url, requestData)
}

func Wait() {
	defsys.Wait()
}
func Visit(url string) error {
	return defsys.Visit(url)
}
