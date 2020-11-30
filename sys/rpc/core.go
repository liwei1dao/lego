package rpc

import (
	"reflect"

	"github.com/liwei1dao/lego/core"
)

type (
	IRpc interface {
		RpcId() string
		GetRpcInfo() (rfs []core.Rpc_Key)
		Register(id string, f interface{})
		RegisterGO(id string, f interface{})
		UnRegister(id string, f interface{})
		Done() (err error)
		OnRegisterRpcData(d interface{}, sf func(d interface{}) ([]byte, error), unsf func(dataType reflect.Type, d []byte) (interface{}, error))
		NewRpcClient(sId, rId string) (clent IRpcClient, err error)
	}
	IRpcServer interface {
		RpcId() string
		GetRpcInfo() (rfs []core.Rpc_Key)
		Register(id string, f interface{})
		RegisterGO(id string, f interface{})
		UnRegister(id string, f interface{})
		Done() (err error)
	}
	IRpcClient interface {
		Done() (err error)
		Call(_func string, params ...interface{}) (interface{}, error)
		CallNR(_func string, params ...interface{}) (err error)
	}
)

var (
	defsys IRpc
)

func OnInit(config map[string]interface{}, option ...Option) (err error) {
	defsys, err = newSys(newOptions(config, option...))
	return
}

func NewSys(option ...Option) (sys IRpc, err error) {
	sys, err = newSys(newOptionsByOption(option...))
	return
}

func RpcId() string {
	return defsys.RpcId()
}
func GetRpcInfo() (rfs []core.Rpc_Key) {
	return defsys.GetRpcInfo()
}
func Register(id string, f interface{}) {
	return defsys.Register(id, f)
}
func RegisterGO(id string, f interface{}) {
	return defsys.RegisterGO(id, f)
}
func UnRegister(id string, f interface{}) {
	return defsys.UnRegister(id, f)
}
func Done() (err error) {
	return defsys.Done()
}
func OnRegisterRpcData(d interface{}, sf func(d interface{}) ([]byte, error), unsf func(dataType reflect.Type, d []byte) (interface{}, error)) {
	return defsys.OnRegisterRpcData(d, sf, unsf, d)
}
func NewRpcClient(sId, rId string) (clent IRpcClient, err error) {
	return defsys.NewRpcClient(sId, rId)
}
