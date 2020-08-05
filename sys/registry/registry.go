package registry

import (
	"github.com/liwei1dao/lego/base"
	"github.com/liwei1dao/lego/core"
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
	return defregistry.PushServiceInfo()
}

func GetServiceById(sId string) (n *ServiceNode, err error) {
	return defregistry.GetServiceById(sId)
}
func GetServiceByType(sType string) (n []*ServiceNode) {
	return defregistry.GetServiceByType(sType)
}

func GetAllServices() (n []*ServiceNode) {
	return defregistry.GetAllServices()
}

func GetServiceByCategory(category core.S_Category) (n []*ServiceNode) {
	return defregistry.GetServiceByCategory(category)
}

func GetRpcSubById(rId core.Rpc_Key) (n []*ServiceNode) {
	return defregistry.GetRpcSubById(rId)
}
