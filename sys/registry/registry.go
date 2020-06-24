package registry

import (
	"github.com/liwei1dao/lego/base"
	"github.com/liwei1dao/lego/core"
	"github.com/liwei1dao/lego/sys/log"
)

var (
	service     base.IClusterService
	defregistry IRegistry
)

func OnInit(s core.IService, opt ...Option) (err error) {
	service = s.(base.IClusterService)
	defregistry, err = newConsulregistry(opt...)
	return
}

func Start() error {
	return defregistry.Start()
}
func Stop() error {
	return defregistry.Stop()
}

func PushServiceInfo() (err error) {
	log.Debugf("主动推送服务信息")
	return defregistry.PushServiceInfo()
}

func GetServiceById(sId string) (n *ServiceNode, err error) {
	return defregistry.GetServiceById(sId)
}
func GetServiceByType(sType string) (n []*ServiceNode, err error) {
	return defregistry.GetServiceByType(sType)
}
func GetServiceByCategory(category core.S_Category) (n []*ServiceNode, err error) {
	return defregistry.GetServiceByCategory(category)
}

func GetRpcSubById(rId core.Rpc_Key) (n []*ServiceNode) {
	return defregistry.GetRpcSubById(rId)
}
