package registry

import "github.com/liwei1dao/lego/core"

type (
	IListener interface {
		FindServiceHandlefunc(snode core.ServiceNode)
		UpDataServiceHandlefunc(snode core.ServiceNode)
		LoseServiceHandlefunc(sId string)
	}
	ISys interface {
		Start() error
		Stop() error
		PushServiceInfo() (err error)
		GetServiceById(sId string) (n *core.ServiceNode, err error)
		GetServiceByType(sType string) (n []*core.ServiceNode)
		GetAllServices() (n []*core.ServiceNode)
		GetServiceByCategory(category core.S_Category) (n []*core.ServiceNode)
		GetRpcSubById(rId core.Rpc_Key) (n []*core.ServiceNode)
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

func Start() error {
	return defsys.Start()
}
func Stop() error {
	return defsys.Stop()
}

func PushServiceInfo() (err error) {
	return defsys.PushServiceInfo()
}

func GetServiceById(sId string) (n *core.ServiceNode, err error) {
	return defsys.GetServiceById(sId)
}
func GetServiceByType(sType string) (n []*core.ServiceNode) {
	return defsys.GetServiceByType(sType)
}

func GetAllServices() (n []*core.ServiceNode) {
	return defsys.GetAllServices()
}

func GetServiceByCategory(category core.S_Category) (n []*core.ServiceNode) {
	return defsys.GetServiceByCategory(category)
}

func GetRpcSubById(rId core.Rpc_Key) (n []*core.ServiceNode) {
	return defsys.GetRpcSubById(rId)
}

func SubscribeRpc(rpc core.Rpc_Key) {

}
