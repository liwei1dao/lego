package registry

import (
	"github.com/liwei1dao/lego/core"
)

type (
	ServiceNode struct {
		Tag          string          `json:"Tag"`          //服务集群标签
		Type         string          `json:"Type"`         //服务类型
		Category     core.S_Category `json:"Category"`     //服务列别
		Id           string          `json:"Id"`           //服务Id
		Version      int32           `json:"Version"`      //服务版本
		RpcId        string          `json:"RpcId"`        //服务通信Id
		PreWeight    int32           `json:"PreWeight"`    //服务负载权重
		RpcSubscribe []core.Rpc_Key  `json:"RpcSubscribe"` //服务开放接口
	}

	IRegistry interface {
		Start() error
		Stop() error
		PushServiceInfo() (err error)
		GetServiceById(sId string) (n *ServiceNode, err error)
		GetServiceByType(sType string) (n []*ServiceNode)
		GetAllServices() (n []*ServiceNode)
		GetServiceByCategory(category core.S_Category) (n []*ServiceNode)
		GetRpcSubById(rId core.Rpc_Key) (n []*ServiceNode)
	}
)

var (
	defsys IRegistry
)

func OnInit(config map[string]interface{}, option ...Option) (err error) {
	defsys, err = newSys(newOptions(config, option...))
	return
}

func NewSys(option ...Option) (sys IRegistry, err error) {
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

func GetServiceById(sId string) (n *ServiceNode, err error) {
	return defsys.GetServiceById(sId)
}
func GetServiceByType(sType string) (n []*ServiceNode) {
	return defsys.GetServiceByType(sType)
}

func GetAllServices() (n []*ServiceNode) {
	return defsys.GetAllServices()
}

func GetServiceByCategory(category core.S_Category) (n []*ServiceNode) {
	return defsys.GetServiceByCategory(category)
}

func GetRpcSubById(rId core.Rpc_Key) (n []*ServiceNode) {
	return defsys.GetRpcSubById(rId)
}
