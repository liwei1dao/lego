package discovery

import (
	"github.com/liwei1dao/lego/core"
)

type (
	ISys interface {
		Start() (err error)
		WatchService() chan []*core.ServiceNode
		GetServices() []*core.ServiceNode
		Close() error
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

func Start() (err error) {
	return defsys.Start()
}
func WatchService() chan []*core.ServiceNode {
	return defsys.WatchService()
}
func GetServices() []*core.ServiceNode {
	return defsys.GetServices()
}
func Close() (err error) {
	return defsys.Close()
}
