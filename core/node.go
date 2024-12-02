package core

import (
	"fmt"
	"net/url"
)

const (
	ServiceNode_Tag  = "tag"
	ServiceNode_Type = "type"
	ServiceNode_Id   = "id"
	ServiceNode_Addr = "addr"
)

func NewServiceNode(value string) (node *ServiceNode, err error) {
	node = &ServiceNode{
		value: value,
	}
	node.meta, err = url.ParseQuery(value)
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
	return this.meta.Get(ServiceNode_Tag)
}
func (this *ServiceNode) Type() string {
	return this.meta.Get(ServiceNode_Type)
}

func (this *ServiceNode) Id() string {
	return this.meta.Get(ServiceNode_Id)
}
func (this *ServiceNode) Addr() string {
	return this.meta.Get(ServiceNode_Addr)
}
func (this *ServiceNode) Path() string {
	return fmt.Sprintf("%s/%s/%s", this.meta.Get(ServiceNode_Tag), this.meta.Get(ServiceNode_Type), this.meta.Get(ServiceNode_Id))
}

func (this *ServiceNode) Set(key, value string) {
	this.meta.Set(key, value)
}

// /判断两个节点是否相等
func (this *ServiceNode) Equal(node IServiceNode) bool {
	if node == nil {
		return false
	}
	return true
}
