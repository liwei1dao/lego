package rpc

import (
	"reflect"

	lgcore "github.com/liwei1dao/lego/core"
	"github.com/liwei1dao/lego/sys/rpc/core"
	"github.com/nats-io/nats.go"
)

func newSys(options Options) (sys *Rpc, err error) {
	var (
		natsconn *nats.Conn
		service  IRpcServer
	)
	natsconn, err = nats.Connect(options.NatsAddr)
	if err != nil {
		return
	}
	service, err = core.NewRpcServer(
		core.SId(options.sId),
		core.SNats(natsconn),
		core.SMaxCoroutine(options.MaxCoroutine))
	if err != nil {
		return
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

func (this *Rpc) GetRpcInfo() (rfs []lgcore.Rpc_Key) {
	return this.service.GetRpcInfo()
}

func (this *Rpc) Register(id lgcore.Rpc_Key, f interface{}) {
	this.service.Register(string(id), f)
}

func (this *Rpc) RegisterGO(id lgcore.Rpc_Key, f interface{}) {
	this.service.RegisterGO(string(id), f)
}

func (this *Rpc) UnRegister(id lgcore.Rpc_Key, f interface{}) {
	this.service.UnRegister(string(id), f)
}

func (this *Rpc) Done() (err error) {
	return this.service.Done()
}

func (this *Rpc) OnRegisterRpcData(d interface{}, sf func(d interface{}) ([]byte, error), unsf func(dataType reflect.Type, d []byte) (interface{}, error)) {
	core.OnRegisterRpcData(d, sf, unsf)
}
func (this *Rpc) OnRegisterJsonRpcData(d interface{}) {
	core.OnRegisterJsonRpc(d)
}
func (this *Rpc) OnRegisterProtoDataData(d interface{}) {
	core.OnRegisterProtoData(d)
}
func (this *Rpc) NewRpcClient(sId, rId string) (clent IRpcClient, err error) {
	clent, err = core.NewRpcClient(
		core.CId(sId),
		core.CrpcId(rId),
		core.CRpcExpired(this.options.RpcExpired),
		core.CNats(this.natsconn))
	return
}
