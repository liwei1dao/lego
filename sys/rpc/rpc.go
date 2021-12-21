package rpc

import (
	"fmt"
	"reflect"
	"time"

	lgcore "github.com/liwei1dao/lego/core"
	"github.com/liwei1dao/lego/sys/rpc/conn"
	"github.com/liwei1dao/lego/sys/rpc/core"
	"github.com/liwei1dao/lego/sys/rpc/serialize"
)

func newSys(options Options) (sys *RPC, err error) {
	var (
		rpcId   string
		service *RPCService
		c       IRPCConnServer
	)
	rpcId = fmt.Sprintf("%s-%s", options.ClusterTag, options.ServiceId)
	if options.RPCConnType == Nats {
		if c, err = conn.NewNatsService(options.Nats_Addr, rpcId); err != nil {
			return
		}
	} else if options.RPCConnType == Kafka {
		if c, err = conn.NewKafkaService(options.ServiceId, options.Kafka_Version, options.Kafka_Host, rpcId); err != nil {
			return
		}
	}
	service = &RPCService{
		rpcId:        rpcId,
		listener:     options.Listener,
		maxCoroutine: options.MaxCoroutine,
		conn:         c,
		functions:    make(map[lgcore.Rpc_Key][]*core.FunctionInfo),
		runch:        make(chan int, options.MaxCoroutine),
	}
	sys = &RPC{
		options: options,
		service: service,
	}
	return
}

type RPC struct {
	options Options
	service IRpcServer
}

func (this *RPC) NewRpcClient(sId, rId string) (clent IRpcClient, err error) {
	var (
		c         IRPCConnClient
		receiveId string
	)
	receiveId = fmt.Sprintf("%s-%s-%s", this.options.ClusterTag, this.options.ServiceId, sId)
	if this.options.RPCConnType == Nats {
		if c, err = conn.NewNatsClient(this.options.Nats_Addr, rId, receiveId); err != nil {
			return
		}
	} else if this.options.RPCConnType == Kafka {
		if c, err = conn.NewKafkaClient(this.options.ServiceId, this.options.Kafka_Version, this.options.Kafka_Host, rId, receiveId); err != nil {
			return
		}
	}
	clent = &RPCClient{
		ServiceId:  sId,
		rpcExpired: time.Second * time.Duration(this.options.RpcExpired),
		conn:       c,
	}
	return
}

func (this *RPC) Start() (err error) {
	return this.service.Start()
}
func (this *RPC) Stop() (err error) {
	return this.service.Stop()
}

func (this *RPC) RpcId() (rpcId string) {
	rpcId = this.service.RpcId()
	return
}

func (this *RPC) GetRpcInfo() (rfs []lgcore.Rpc_Key) {
	return this.service.GetRpcInfo()
}

func (this *RPC) Register(id lgcore.Rpc_Key, f interface{}) {
	this.service.Register(id, f)
}

func (this *RPC) RegisterGO(id lgcore.Rpc_Key, f interface{}) {
	this.service.RegisterGO(id, f)
}

func (this *RPC) UnRegister(id lgcore.Rpc_Key, f interface{}) {
	this.service.UnRegister(id, f)
}

func (this *RPC) OnRegisterRpcData(d interface{}, sf func(d interface{}) ([]byte, error), unsf func(dataType reflect.Type, d []byte) (interface{}, error)) {
	serialize.OnRegisterRpcData(d, sf, unsf)
}
func (this *RPC) OnRegisterJsonRpcData(d interface{}) {
	serialize.OnRegisterJsonRpc(d)
}
func (this *RPC) OnRegisterProtoDataData(d interface{}) {
	serialize.OnRegisterProtoData(d)
}
