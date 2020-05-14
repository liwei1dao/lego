package rpc

import (
	"fmt"
	"reflect"

	"github.com/liwei1dao/lego/core"
	rcore "github.com/liwei1dao/lego/sys/rpc/rcore"
	rpcserialize "github.com/liwei1dao/lego/sys/rpc/serialize"

	"github.com/nats-io/nats.go"
)

var (
	options *Options
	rnats   *nats.Conn
	srpc    IRpcServer
)

func OnInit(s core.IService, opt ...Option) (err error) {
	options = newOptions(opt...)
	nc, err := nats.Connect(options.NatsAddr)
	if err != nil {
		return err
	}
	rnats = nc
	rpcserialize.SerializeInit()
	srpc, err = rcore.NewRpcServer(
		rcore.SId(options.sId),
		rcore.SNats(rnats),
		rcore.SMaxCoroutine(options.MaxCoroutine))
	return
}

func OnRegisterRpcData(d interface{}, sf func(d interface{}) ([]byte, error), unsf func(dataType reflect.Type, d []byte) (interface{}, error)) {
	rpcserialize.OnRegister(d, sf, unsf)
}

func NewRpcClient(sId, rId string) (clent IRpcClient, err error) {
	if options == nil {
		return nil, fmt.Errorf("rpc 系统未初始化")
	}
	clent, err = rcore.NewRpcClient(
		rcore.CId(sId),
		rcore.CrpcId(rId),
		rcore.CRpcExpired(options.RpcExpired),
		rcore.CNats(rnats))
	return
}

func RpcId() (rpcId string, err error) {
	if srpc == nil {
		return "", fmt.Errorf("rpc 系统未初始化")
	}
	rpcId = srpc.RpcId()
	return
}

func Register(id string, f interface{}) (err error) {
	if srpc == nil {
		return fmt.Errorf("rpc 系统未初始化")
	}
	srpc.Register(id, f)
	return
}

func RegisterGO(id string, f interface{}) (err error) {
	if srpc == nil {
		return fmt.Errorf("rpc 系统未初始化")
	}
	srpc.RegisterGO(id, f)
	return
}

func UnRegister(id string, f interface{}) (err error) {
	if srpc == nil {
		return fmt.Errorf("rpc 系统未初始化")
	}
	srpc.UnRegister(id, f)
	return
}
