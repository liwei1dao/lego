package base

import (
	"github.com/liwei1dao/lego/core"
)

type ISingleService interface {
	core.IService
}

type IClusterService interface {
	core.IService
	GetTag() string                                                                                                   //获取集群标签
	GetCategory() core.S_Category                                                                                     //服务类别 例如游戏服
	GetRpcId() string                                                                                                 //获取rpc通信id
	GetPreWeight() int32                                                                                              //集群服务负载值 暂时可以不用理会
	SetPreWeight(weight int32)                                                                                        //设置服务器权重
	GetServiceMonitor() core.IServiceMonitor                                                                          //获取监控模块
	GetSessionsByCategory(category core.S_Category) (ss []core.IServiceSession)                                       //按服务类别获取服务列表
	DefauleRpcRouteRules(stype string) (ss core.IServiceSession, err error)                                           //默认rpc路由规则
	RpcInvokeById(sId string, rkey core.Rpc_Key, iscall bool, arg ...interface{}) (result interface{}, err error)     //执行远程服务Rpc方法
	RpcInvokeByType(sType string, rkey core.Rpc_Key, iscall bool, arg ...interface{}) (result interface{}, err error) //根据路由规则执行远程方法
	ReleaseRpc(rkey core.Rpc_Key, arg ...interface{})                                                                 //发布Rpc
	Register(id core.Rpc_Key, f interface{})                                                                          //注册RPC远程方法
	RegisterGO(id core.Rpc_Key, f interface{})                                                                        //注册RPC远程方法
	Subscribe(id core.Rpc_Key, f interface{}) (err error)                                                             //订阅Rpc
	UnSubscribe(id core.Rpc_Key, f interface{}) (err error)                                                           //订阅Rpc
}
