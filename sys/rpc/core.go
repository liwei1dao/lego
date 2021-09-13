package rpc

import (
	"reflect"

	lgcore "github.com/liwei1dao/lego/core"
	"github.com/liwei1dao/lego/sys/rpc/core"
)

const (
	Nats RPCConnType = iota
	Kafka
)

type (
	RPCConnType int
	ISys        interface {
		IRpcServer
		NewRpcClient(sId, rId string) (clent IRpcClient, err error)
		OnRegisterRpcData(d interface{}, sf func(d interface{}) ([]byte, error), unsf func(dataType reflect.Type, d []byte) (interface{}, error))
		OnRegisterJsonRpcData(d interface{})
		OnRegisterProtoDataData(d interface{})
	}
	IRpcServer interface {
		RpcId() string
		GetRpcInfo() (rfs []lgcore.Rpc_Key)
		Register(id lgcore.Rpc_Key, f interface{})
		RegisterGO(id lgcore.Rpc_Key, f interface{})
		UnRegister(id lgcore.Rpc_Key, f interface{})
		Start() (err error)
		Stop() (err error)
	}
	IRpcClient interface {
		Stop() (err error)
		Call(_func string, params ...interface{}) (interface{}, error)
		CallNR(_func string, params ...interface{}) (err error)
	}
	RPCListener interface {
		NoFoundFunction(fn string) (*core.FunctionInfo, error)
		BeforeHandle(fn string, callInfo *core.CallInfo) error
		OnTimeOut(fn string, Expired int64)
		OnError(fn string, callInfo *core.CallInfo, err error)
		OnComplete(fn string, callInfo *core.CallInfo, result *core.ResultInfo, exec_time int64)
	}
	IRPCConnServer interface {
		Callback(callinfo core.CallInfo) error
		Start(service core.IRpcServer) (err error)
		Stop() (err error)
	}
	IRPCConnClient interface {
		Stop() (err error)
		Delete(key string) (err error)
		Call(callInfo core.CallInfo, callback chan core.ResultInfo) error
		CallNR(callInfo core.CallInfo) error
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

func RpcId() string {
	return defsys.RpcId()
}
func GetRpcInfo() (rfs []lgcore.Rpc_Key) {
	return defsys.GetRpcInfo()
}
func Register(id lgcore.Rpc_Key, f interface{}) {
	defsys.Register(id, f)
}
func RegisterGO(id lgcore.Rpc_Key, f interface{}) {
	defsys.RegisterGO(id, f)
}
func UnRegister(id lgcore.Rpc_Key, f interface{}) {
	defsys.UnRegister(id, f)
}
func NewRpcClient(sId, rId string) (clent IRpcClient, err error) {
	return defsys.NewRpcClient(sId, rId)
}
func OnRegisterRpcData(d interface{}, sf func(d interface{}) ([]byte, error), unsf func(dataType reflect.Type, d []byte) (interface{}, error)) {
	defsys.OnRegisterRpcData(d, sf, unsf)
}
func OnRegisterJsonRpcData(d interface{}) {
	defsys.OnRegisterJsonRpcData(d)
}
func OnRegisterProtoDataData(d interface{}) {
	defsys.OnRegisterProtoDataData(d)
}
