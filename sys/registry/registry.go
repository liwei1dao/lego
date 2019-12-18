package registry

import (
	"fmt"
	"lego/base"
	"lego/core"
	"lego/sys/log"
	"time"
)

var (
	service  base.IClusterService
	registry IRegistry
	exit     chan bool
	isstart  bool
)

func OnInit(s base.IClusterService, opt ...Option) (err error) {
	service = s
	registry, err = newConsulregistry(opt...)
	return
}

//警告 禁止反复调用
func Registry() (err error) {
	if !isstart {
		if err := registry.RegisterSNode(&ServiceNode{
			Tag:       service.GetTag(),
			Id:        service.GetId(),
			Type:      service.GetType(),
			Category:  service.GetCategory(),
			Version:   service.GetVersion(),
			RpcId:     service.GetRpcId(),
			PreWeight: service.GetPreWeight(),
		}); err != nil {
			return err
		}
		exit = make(chan bool)
		go run(exit)
		return
	} else {
		return fmt.Errorf("重复 Registry")
	}
}

//注销
func Deregister() (err error) {
	exit <- true
	isstart = false
	if err = registry.DeregisterSNode(service.GetId()); err != nil {
		log.Errorf("registry 注销服务失败 err:%s", err.Error())
	} else {
		log.Infof("registry 成功注销服务 %s", service.GetId())
	}
	return
}

func GetServiceById(sId string) (node *ServiceNode, err error) {
	if registry == nil {
		return nil, fmt.Errorf("registry 系统未初始化")
	}
	node, err = registry.GetServiceById(sId)
	return
}
func GetServiceByType(sId string) (nodes []*ServiceNode, err error) {
	if registry == nil {
		return nil, fmt.Errorf("registry 系统未初始化")
	}
	nodes, err = registry.GetServiceByType(sId)
	return
}

func GetServiceByCategory(category core.S_Category) (nodes []*ServiceNode, err error) {
	if registry == nil {
		return nil, fmt.Errorf("registry 系统未初始化")
	}
	nodes, err = registry.GetServiceByCategory(category)
	return
}

//注册Rpc订阅
func RegisterRpcFunc(rId core.Rpc_Key, sId string) (err error) {
	if registry == nil {
		return fmt.Errorf("registry 系统未初始化")
	}
	err = registry.RegisterRpcSub(rId, sId)
	return
}
func UnRegisterRpcFunc(rId core.Rpc_Key, sId string) (err error) {
	if registry == nil {
		return fmt.Errorf("registry 系统未初始化")
	}
	err = registry.UnRegisterRpcSub(rId, sId)
	return
}

func GetRpcFunc(rId core.Rpc_Key) (data *RpcFuncInfo, err error) {
	if registry == nil {
		return nil, fmt.Errorf("registry 系统未初始化")
	}
	data, err = registry.GetRpcSubById(rId)
	return
}
func run(exit chan bool) {
	if registry == nil {
		log.Errorf("registry 系统未初始化")
		return
	}
	if registry.Options().RegisterInterval <= time.Duration(0) {
		return
	}
	t := time.NewTicker(registry.Options().RegisterInterval)
	for {
		select {
		case <-t.C:
			err := registry.RegisterSNode(&ServiceNode{
				Tag:       service.GetTag(),
				Id:        service.GetId(),
				Type:      service.GetType(),
				Category:  service.GetCategory(),
				Version:   service.GetVersion(),
				RpcId:     service.GetRpcId(),
				PreWeight: service.GetPreWeight(),
			})
			if err != nil {
				log.Warnf("service run Server.Register error: ", err)
			}
		case <-exit:
			t.Stop()
			return
		}
	}
}
