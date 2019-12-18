package gate

import (
	"fmt"
	"lego/core"
	cbase "lego/core/cbase"
	"lego/sys/proto"
	cont "lego/utils/concurrent"
)

type AgentMgrComp struct {
	cbase.ModuleCompBase
	Agents *cont.BeeMap
}

func (this *AgentMgrComp) Init(service core.IService, module core.IModule, comp core.IModuleComp, _Settings map[string]interface{}) (err error) {
	err = this.ModuleCompBase.Init(service, module, comp, _Settings)
	this.Agents = cont.NewBeeMap()
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
	return
}

func (this *AgentMgrComp) Connect(a IAgent) {
	//Log.Infof("有连接进入 %s", a.Id())
	this.Agents.Set(a.Id(), a)
}
func (this *AgentMgrComp) DisConnect(a IAgent) {
	//Log.Infof("有连接关闭 %s", a.Id())
	this.Agents.Delete(a.Id())
}
func (this *AgentMgrComp) SendMsg(aId string, msg proto.IMessage) (result int, err string) {
	//Log.Infof("发送用户消息 aId =%s ComId= %d MsgId= %d",aId, msg.GetComId(), msg.GetMsgId())
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
func (this *AgentMgrComp) Close(aId string) (result string, err string) {
	agent := this.Agents.Get(aId)
	if agent == nil {
		err = fmt.Sprintf("No agent found " + aId)
		return
	}
	agent.(IAgent).OnClose()
	return
}
