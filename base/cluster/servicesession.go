package cluster

import (
	"github.com/liwei1dao/lego/base"
	"github.com/liwei1dao/lego/core"
	"github.com/liwei1dao/lego/sys/rpc"
)

func NewServiceSession(node *core.ServiceNode) (ss base.IClusterServiceSession, err error) {
	session := new(ServiceSession)
	session.node = node
	session.client, err = rpc.NewRpcClient(node.Id, node.RpcId)
	ss = session
	return
}

type ServiceSession struct {
	node   *core.ServiceNode
	client rpc.IRpcClient
}

func (this *ServiceSession) GetId() string {
	return this.node.Id
}
func (this *ServiceSession) GetIp() string {
	return this.node.IP
}
func (this *ServiceSession) GetRpcId() string {
	return this.node.RpcId
}

func (this *ServiceSession) GetType() string {
	return this.node.Type
}
func (this *ServiceSession) GetVersion() string {
	return this.node.Version
}
func (this *ServiceSession) SetVersion(v string) {
	this.node.Version = v
}

func (this *ServiceSession) GetPreWeight() float64 {
	return this.node.PreWeight
}
func (this *ServiceSession) SetPreWeight(p float64) {
	this.node.PreWeight = p
}
func (this *ServiceSession) Done() {
	this.client.Stop()
}
func (this *ServiceSession) Call(f core.Rpc_Key, params ...interface{}) (interface{}, error) {
	return this.client.Call(string(f), params...)
}
func (this *ServiceSession) CallNR(f core.Rpc_Key, params ...interface{}) (err error) {
	return this.client.CallNR(string(f), params...)
}
