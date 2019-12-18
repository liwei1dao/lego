package gate

import (
	"lego/core"
	"lego/core/cbase"
	"lego/sys/proto"
)

type Gate struct {
	cbase.ModuleBase
	//AgentMgrComp IAgentMgrComp
	//LocalRouteMgrComp  *LocalRouteMgrComp
	//RemoteRouteMgrComp *RemoteRouteMgrComp
	//WsServerComp       *WsServerComp
	//TcpServerComp      *TcpServerComp
}

func (this *Gate) RegisterLocalRoute(comId uint16, f func(session core.IUserSession, msg proto.IMessage) (code int, err string)) (err error) {
	//err = this.LocalRouteMgrComp.RegisterRoute(comId, f)
	return
}

func (this *Gate) UnRegisterLocalRoute(comId uint16, f func(session core.IUserSession, msg proto.IMessage)) {
	//this.LocalRouteMgrComp.UnRegisterRoute(comId, f)
}

//需重构处理  内部函数为重构代码
//代理链接
func (this *Gate) Connect(a IAgent) {
	//this.AgentMgrComp.Connect(a)
}

//代理关闭
func (this *Gate) DisConnect(a IAgent) {
	//this.AgentMgrComp.DisConnect(a)
}

//接收代理消息
func (this *Gate) OnRoute(a IAgent, msg proto.IMessage) (code int, err error) {
	return
}

//主动关闭代理
func (this *Gate) CloseAgent(sId string) (result string, err string) {
	return
}

//发送代理消息
func (this *Gate) SendMsg(sId string, msg proto.IMessage) (result int, err string) {
	return
}

//func (this *Gate) OnInstallComp() {
//	this.ModuleBase.OnInstallComp()
//	this.AgentMgrComp = this.RegisterComp(new(AgentMgrComp)).(*AgentMgrComp)
//	this.LocalRouteMgrComp = this.RegisterComp(new(LocalRouteMgrComp)).(*LocalRouteMgrComp)
//	this.RemoteRouteMgrComp = this.RegisterComp(new(LocalRouteMgrComp)).(*RemoteRouteMgrComp)
//	this.TcpServerComp = this.RegisterComp(new(TcpServerComp)).(*TcpServerComp)
//	this.WsServerComp = this.RegisterComp(new(WsServerComp)).(*WsServerComp)
//	this.LocalRouteMgrComp.NewSession = NewLocalSession
//	this.RemoteRouteMgrComp.NewSession = NewRemoteSession
//}
