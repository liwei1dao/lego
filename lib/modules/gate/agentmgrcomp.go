package gate

import (
	"fmt"

	"github.com/liwei1dao/lego/core"
	"github.com/liwei1dao/lego/core/cbase"
	"github.com/liwei1dao/lego/sys/proto"
	"github.com/liwei1dao/lego/utils/container"
)

type AgentMgrComp struct {
	cbase.ModuleCompBase
	Agents *container.BeeMap
}

func (this *AgentMgrComp) Init(service core.IService, module core.IModule, comp core.IModuleComp, options core.IModuleOptions) (err error) {
	err = this.ModuleCompBase.Init(service, module, comp, options)
	this.Agents = container.NewBeeMap()
	return
}
func (this *AgentMgrComp) Destroy() (err error) {
	if err = this.ModuleCompBase.Destroy(); err != nil {
		return
	}
	for _, v := range this.Agents.Items() {
		agent := v.(IAgent)
		agent.OnClose()
	}
	this.Agents.DeleteAll()
	return
}
func (this *AgentMgrComp) Connect(a IAgent) {
	this.Agents.Set(a.Id(), a)
}
func (this *AgentMgrComp) DisConnect(a IAgent) {
	this.Agents.Delete(a.Id())
}
func (this *AgentMgrComp) SendMsg(aId string, msg proto.IMessage) (result int, err string) {
	agent := this.Agents.Get(aId)
	if agent == nil {
		err = fmt.Sprintf("No agent found " + aId)
		return
	}
	e := agent.(IAgent).WriteMsg(msg)
	if e != nil {
		err = e.Error()
	} else {
		result = 0
	}
	return
}
func (this *AgentMgrComp) SendMsgByGroup(aIds []string, msg proto.IMessage) (result []string, err string) {
	result = make([]string, 0)
	for _, a := range aIds {
		agent := this.Agents.Get(a)
		if agent == nil {
			result = append(result, a)
			continue
		}
		e := agent.(IAgent).WriteMsg(msg)
		if e != nil {
			result = append(result, a)
		}
	}
	return
}
func (this *AgentMgrComp) SendMsgByBroadcast(msg proto.IMessage) (result int, err string) {
	for _, v := range this.Agents.Items() {
		v.(IAgent).WriteMsg(msg)
	}
	return
}
func (this *AgentMgrComp) Close(aId string) (result string, err string) {
	agent := this.Agents.Get(aId)
	if agent == nil {
		err = fmt.Sprintf("No agent found " + aId)
		return
	}
	agent.(IAgent).OnClose()
	return
}
