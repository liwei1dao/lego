package cbase

import (
	"github.com/liwei1dao/lego/core"
	"github.com/liwei1dao/lego/sys/registry"
	"github.com/liwei1dao/lego/sys/rpc"
)

func NewServiceSession(node *registry.ServiceNode) (ss core.IServiceSession, err error) {
	session := new(ServiceSession)
	session.node = node
	session.rpc, err = rpc.NewRpcClient(node.Id, node.RpcId)
	ss = session
	return
}

type ServiceSession struct {
	node *registry.ServiceNode
	rpc  rpc.IRpcClient
}

func (this *ServiceSession) GetId() string {
	return this.node.Id
}

func (this *ServiceSession) GetRpcId() string {
	return this.node.RpcId
}

func (this *ServiceSession) GetType() string {
	return this.node.Type
}
func (this *ServiceSession) GetVersion() int32 {
	return this.node.Version
}
func (this *ServiceSession) SetVersion(v int32) {
	this.node.Version = v
}

func (this *ServiceSession) GetPreWeight() int32 {
	return this.node.PreWeight
}
func (this *ServiceSession) SetPreWeight(p int32) {
	this.node.PreWeight = p
}
func (this *ServiceSession) Done() {
	this.rpc.Done()
}
func (this *ServiceSession) Call(f core.Rpc_Key, params ...interface{}) (interface{}, error) {
	return this.rpc.Call(string(f), params...)
}
func (this *ServiceSession) CallNR(f core.Rpc_Key, params ...interface{}) (err error) {
	return this.rpc.CallNR(string(f), params...)
}
