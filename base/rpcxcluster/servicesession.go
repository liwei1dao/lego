package rpcxcluster

import (
	"context"
	"fmt"

	"github.com/liwei1dao/lego/base"
	"github.com/liwei1dao/lego/core"
	"github.com/liwei1dao/lego/sys/registry"
	"github.com/liwei1dao/lego/sys/rpcx"
	"github.com/smallnest/rpcx/client"
)

func NewServiceSession(node *registry.ServiceNode) (ss base.IRPCXServiceSession, err error) {
	session := new(ServiceSession)
	session.node = node
	session.client, err = rpcx.NewRpcClient(fmt.Sprintf("%s:%d", node.IP, node.PreWeight))
	ss = session
	return
}

type ServiceSession struct {
	node   *registry.ServiceNode
	client rpcx.IRPCXClient
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
func (this *ServiceSession) Call(ctx context.Context, serviceMethod core.Rpc_Key, args interface{}, reply interface{}) error {
	return this.client.Call(ctx, serviceMethod, args, reply)
}
func (this *ServiceSession) CallNR(ctx context.Context, serviceMethod core.Rpc_Key, args interface{}, reply interface{}) (*client.Call, error) {
	return this.client.Go(ctx, serviceMethod, args, reply, nil)
}
