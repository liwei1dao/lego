package registry

import "github.com/liwei1dao/lego/core"

type (
	ServiceNode struct {
		Tag          string          `json:"Tag"`
		Type         string          `json:"Type"`
		Category     core.S_Category `json:"Category"`
		Id           string          `json:"Id"`
		Version      int32           `json:"Version"`
		RpcId        string          `json:"RpcId"`
		PreWeight    int32           `json:"PreWeight"`
		RpcSubscribe []core.Rpc_Key  `json:"RpcSubscribe"`
	}

	IRegistry interface {
		Start() error
		Stop() error
		PushServiceInfo() (err error)
		GetServiceById(sId string) (n *ServiceNode, err error)
		GetServiceByType(sType string) (n []*ServiceNode, err error)
		GetServiceByCategory(category core.S_Category) (n []*ServiceNode, err error)
		GetRpcSubById(rId core.Rpc_Key) (n []*ServiceNode)
	}
)
