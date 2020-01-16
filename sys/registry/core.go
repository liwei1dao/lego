package registry

import (
	"github.com/liwei1dao/lego/core"
)

type ServiceNode struct {
	Tag       string          `json:"Tag"`
	Type      string          `json:"Type"`
	Category  core.S_Category `json:"Category"`
	Id        string          `json:"Id"`
	Version   int32           `json:"Version"`
	RpcId     string          `json:"RpcId"`
	PreWeight int32           `json:"PreWeight"`
}

type RpcFuncInfo struct {
	Id            string                    //Rpc 方法Id
	SubServiceIds map[string]*ServiceRpcSub //订阅数据
}

type ServiceRpcSub struct {
	Id     string
	SubNum int32
}

type IRegistry interface {
	Options() Options
	GetServiceById(sId string) (n *ServiceNode, err error)
	GetServiceByType(sType string) (n []*ServiceNode, err error)
	GetServiceByCategory(category core.S_Category) (n []*ServiceNode, err error)
	RegisterSNode(s *ServiceNode) (err error)
	DeregisterSNode(sId string) (err error)
	RegisterRpcSub(rId core.Rpc_Key, sId string) (err error)
	UnRegisterRpcSub(rId core.Rpc_Key, sId string) (err error)
	GetRpcSubById(rId core.Rpc_Key) (rf *RpcFuncInfo, err error)
}
