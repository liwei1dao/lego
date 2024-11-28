package core

import (
	"fmt"
	"net/url"
)

const (
	ServiceNode_Tag     = "tag"
	ServiceNode_Type    = "type"
	ServiceNode_Id      = "id"
	ServiceNode_Version = "ver"
	ServiceNode_Addr    = "addr"
	ServiceNode_State   = "state"
)

func NewServiceNode(value string) (node *ServiceNode, err error) {
	var (
		values url.Values
	)
	if values, err = url.ParseQuery(value); err == nil {
		return
	}
	node = &ServiceNode{
		value: value,
		meta:  values,
	}
	return
}

type ServiceNode struct {
	value string
	meta  url.Values
}

func (this *ServiceNode) Value() string {
	return this.value
}
func (this *ServiceNode) Tag() string {
	return this.GetMate(ServiceNode_Tag)
}

func (this *ServiceNode) Type() string {
	return this.GetMate(ServiceNode_Type)
}
func (this *ServiceNode) Id() string {
	return this.GetMate(ServiceNode_Id)
}
func (this *ServiceNode) Version() string {
	return this.GetMate(ServiceNode_Version)
}
func (this *ServiceNode) Addr() string {
	return this.GetMate(ServiceNode_Addr)
}
func (this *ServiceNode) State() string {
	return this.GetMate(ServiceNode_State)
}

func (this *ServiceNode) GetNodePath() string {
	return fmt.Sprintf("%s/%s/%s", this.meta.Get(ServiceNode_Tag), this.meta.Get(ServiceNode_Type), this.meta.Get(ServiceNode_Id))
}

// 写入元数据
func (this *ServiceNode) SetMate(name, value string) {
	this.meta.Set(name, value)
}

// 写入元数据
func (this *ServiceNode) GetMate(name string) (value string) {
	value = this.meta.Get(name)
	return
}

// /判断两个节点是否相等
func (this *ServiceNode) Equal(node IServiceNode) bool {
	return true
}
