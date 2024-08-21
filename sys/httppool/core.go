package httppool

import "net/http"

/*
系统描述:进程级别的事件系统
*/
type (
	ISys interface {
		GetClient() *http.Client
	}
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

func GetClient() *http.Client {
	return defsys.GetClient()
}
