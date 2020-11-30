package rpc

import (
	"reflect"

	"github.com/liwei1dao/lego/core"
	"github.com/liwei1dao/lego/sys/rpc/core"
	rcore "github.com/liwei1dao/lego/sys/rpc/rcore"
	"github.com/nats-io/nats.go"
)

func newSys(options Options) (sys *Rpc, err error) {
	natsconn, err := nats.Connect(options.NatsAddr)
	if err != nil {
		return err
	}
	service, err = rcore.NewRpcServer(
		rcore.SId(options.sId),
		rcore.SNats(natsconn),
		rcore.SMaxCoroutine(options.MaxCoroutine))
	if err != nil {
		return err
	}
	sys = &Rpc{
		options:  options,
		natsconn: natsconn,
		service:  service,
	}
	return
}

type Rpc struct {
	options  Options
	natsconn *nats.Conn
	service  IRpcServer
}

func (this *Rpc) RpcId() (rpcId string) {
	rpcId = this.service.RpcId()
	return
}

func GetRpcInfo() (rfs []core.Rpc_Key) {
	return this.service.GetRpcInfo()
}

func Register(id string, f interface{}) {
	this.service.Register(id, f)
}

func RegisterGO(id string, f interface{}) {
	this.service.RegisterGO(id, f)
}

func UnRegister(id string, f interface{}) {
	this.service.UnRegister(id, f)
}

func OnRegisterRpcData(d interface{}, sf func(d interface{}) ([]byte, error), unsf func(dataType reflect.Type, d []byte) (interface{}, error)) {
	onRegister(d, sf, unsf)
}

func NewRpcClient(sId, rId string) (clent IRpcClient, err error) {
	clent, err = core.NewRpcClient(
		core.CId(sId),
		core.CrpcId(rId),
		core.CRpcExpired(this.options.RpcExpired),
		core.CNats(this.natsconn))
	return
}
