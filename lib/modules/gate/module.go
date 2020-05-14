package gate

import (
	"github.com/liwei1dao/lego/core"
	"github.com/liwei1dao/lego/core/cbase"
	"github.com/liwei1dao/lego/sys/log"
	"github.com/liwei1dao/lego/sys/proto"
)

type Gate struct {
	cbase.ModuleBase
	CustomRouteComp    ICustomRouteComp
	LocalRouteMgrComp  ILocalRouteMgrComp
	RemoteRouteMgrComp IRemoteRouteMgrComp
	AgentMgrComp       IAgentMgrComp
	//WsServerComp       *WsServerComp
	//TcpServerComp      *TcpServerComp
}

func (this *Gate) GetLocalRouteMgrComp() ILocalRouteMgrComp {
	return this.LocalRouteMgrComp
}

//需重构处理  内部函数为重构代码d
func (this *Gate) RegisterRemoteRoute(comId uint16, sId string) (result string, err string) {
	result, err = this.RemoteRouteMgrComp.RegisterRoute(comId, sId)
	return
}
func (this *Gate) RegisterLocalRoute(comId uint16, f func(session core.IUserSession, msg proto.IMessage) (code int, err string)) {
	this.LocalRouteMgrComp.RegisterRoute(comId, f)
}
func (this *Gate) UnRegisterRemoteRoute(comId uint16, sType, sId string) {
	this.RemoteRouteMgrComp.UnRegisterRoute(comId, sType, sId)
}
func (this *Gate) UnRegisterLocalRoute(comId uint16, f func(session core.IUserSession, msg proto.IMessage) (code int, err string)) {
	this.LocalRouteMgrComp.UnRegisterRoute(comId, f)
}

//需重构处理  内部函数为重构代码
//代理链接
func (this *Gate) Connect(a IAgent) {
	log.Debugf("有新的用户链接进来IP:[%s] Id:[%s]", a.IP, a.Id)
	this.AgentMgrComp.Connect(a)
}

//代理关闭
func (this *Gate) DisConnect(a IAgent) {
	log.Debugf("有用户链接断开IP:[%s] Id:[%s]", a.IP, a.Id)
	this.AgentMgrComp.DisConnect(a)
}

//接收代理消息
func (this *Gate) OnRoute(a IAgent, msg proto.IMessage) (code core.ErrorCode, err error) {
	if this.CustomRouteComp != nil { //优先自定义网关
		if code, err = this.CustomRouteComp.OnRoute(a, msg); code == core.ErrorCode_Success || err != nil { //优先自定义网关
			return
		}
	}

	if code, err = this.LocalRouteMgrComp.OnRoute(a, msg); code == core.ErrorCode_Success || err != nil { //其次本地网关
		return
	}

	code, err = this.RemoteRouteMgrComp.OnRoute(a, msg)
	return
}

//主动关闭代理
func (this *Gate) CloseAgent(sId string) (result string, err string) {
	return this.AgentMgrComp.Close(sId)
}

//发送代理消息
func (this *Gate) SendMsg(sId string, msg proto.IMessage) (result int, err string) {
	return this.AgentMgrComp.SendMsg(sId, msg)
}

//广播代理消息
func (this *Gate) RadioMsg(sIds []string, msg proto.IMessage) (result int, err string) {
	for _, v := range sIds {
		this.AgentMgrComp.SendMsg(v, msg)
	}
	return
}

func (this *Gate) OnInstallComp() {
	this.ModuleBase.OnInstallComp()
	this.LocalRouteMgrComp = this.RegisterComp(new(LocalRouteMgrComp)).(*LocalRouteMgrComp)
	this.RemoteRouteMgrComp = this.RegisterComp(new(RemoteRouteMgrComp)).(*RemoteRouteMgrComp)
}
