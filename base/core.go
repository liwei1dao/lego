package base

import (
	"context"

	"github.com/liwei1dao/lego/core"
	"github.com/smallnest/rpcx/client"
)

type Result struct {
	Index  string
	Result interface{}
	Err    error
}

type ISingleService interface {
	core.IService
}

type IClusterServiceBase interface {
	core.IService
	GetTag() string                          //获取集群标签
	GetCategory() core.S_Category            //服务类别 例如游戏服
	GetRpcId() string                        //获取rpc通信id
	GetPreWeight() float64                   //集群服务负载值 暂时可以不用理会
	SetPreWeight(weight int32)               //设置服务器权重
	GetServiceMonitor() core.IServiceMonitor //获取监控模块
}

type IClusterServiceSession interface {
	GetId() string
	GetIp() string
	GetRpcId() string
	GetType() string
	GetVersion() string
	SetVersion(v string)
	GetPreWeight() float64
	SetPreWeight(p float64)
	Done()
	CallNR(_func core.Rpc_Key, params ...interface{}) (err error)
	Call(_func core.Rpc_Key, params ...interface{}) (interface{}, error)
}

type IClusterService interface {
	IClusterServiceBase
	GetSessionsByCategory(category core.S_Category) (ss []IClusterServiceSession)         //按服务类别获取服务列表
	DefauleRpcRouteRules(stype string, sip string) (ss IClusterServiceSession, err error) //默认rpc路由规则
	RpcInvokeById(sId string, rkey core.Rpc_Key, iscall bool, arg ...interface{}) (result interface{}, err error)
	RpcInvokeByIds(sId []string, rkey core.Rpc_Key, iscall bool, arg ...interface{}) (results map[string]*Result, err error)               //执行远程服务Rpc方法
	RpcInvokeByType(sType string, rkey core.Rpc_Key, iscall bool, arg ...interface{}) (result interface{}, err error)                      //根据路由规则执行远程方法
	RpcInvokeByIp(sIp, sType string, rkey core.Rpc_Key, iscall bool, arg ...interface{}) (result interface{}, err error)                   //根据目标IP和类型执行远程方法
	RpcInvokeByIps(sIp []string, sType string, rkey core.Rpc_Key, iscall bool, arg ...interface{}) (results map[string]*Result, err error) //根据目标IP集和类型执行远程方法
	ReleaseRpc(rkey core.Rpc_Key, arg ...interface{})                                                                                      //发布Rpc
	Register(id core.Rpc_Key, f interface{})                                                                                               //注册RPC远程方法
	RegisterGO(id core.Rpc_Key, f interface{})                                                                                             //注册RPC远程方法
	Subscribe(id core.Rpc_Key, f interface{}) (err error)                                                                                  //订阅Rpc
	UnSubscribe(id core.Rpc_Key, f interface{}) (err error)                                                                                //订阅Rpc
}

type IRPCXServiceSession interface {
	GetId() string
	GetIp() string
	GetRpcId() string
	GetType() string
	GetVersion() string
	SetVersion(v string)
	GetPreWeight() float64
	SetPreWeight(p float64)
	Done()
	Call(ctx context.Context, serviceMethod string, args interface{}, reply interface{}) error
	Go(ctx context.Context, serviceMethod string, args interface{}, reply interface{}) (*client.Call, error)
}

type IRPCXService interface {
	IClusterServiceBase
	DefauleRpcRouteRules(stype string, sip string) (ss IRPCXServiceSession, err error) //默认rpc路由规则
	Register(rcvr interface{}) (err error)
	RegisterFunction(fn interface{}) (err error)
	RegisterFunctionName(name string, fn interface{}) (err error)
	RpcCallById(sId string, serviceMethod string, ctx context.Context, args interface{}, reply interface{}) (err error)
	RpcGoById(sId string, serviceMethod string, ctx context.Context, args interface{}, reply interface{}) (call *client.Call, err error)
	RpcCallByType(sType string, serviceMethod string, ctx context.Context, args interface{}, reply interface{}) (err error)
	RpcGoByType(sType string, serviceMethod string, ctx context.Context, args interface{}, reply interface{}) (call *client.Call, err error)
}
