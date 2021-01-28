package gate

import (
	"net"

	"github.com/liwei1dao/lego/base"
	"github.com/liwei1dao/lego/core"
	"github.com/liwei1dao/lego/sys/proto"
)

const ( //Rpc
	Rpc_GateRouteRegister       core.Rpc_Key = "Rpc_GateRouteRegister"       //网关路由注册
	Rpc_GateCustomRouteRegister core.Rpc_Key = "Rpc_GateCustomRouteRegister" //自定义网关路由注册
	Rpc_GateRoute               core.Rpc_Key = "Rpc_GateRoute"               //网关路由
	Rpc_GateAgentsIsKeep        core.Rpc_Key = "Rpc_GateAgentsIsKeep"        //校验代理是否还在
	RPC_GateAgentBuild          core.Rpc_Key = "RPC_GateAgentBuild"          //代理绑定
	RPC_GateAgentUnBuild        core.Rpc_Key = "RPC_GateAgentUnBuild"        //代理解绑
	RPC_GateSendMsg             core.Rpc_Key = "RPC_GateSendMsg"             //代理发送消息
	RPC_GateSendMsgByGroup      core.Rpc_Key = "RPC_GateSendMsgByGroup"      //代理群发消息
	RPC_GateSendMsgByBroadcast  core.Rpc_Key = "RPC_GateSendMsgByBroadcast"  //代理广播消息
	RPC_GateAgentClose          core.Rpc_Key = "RPC_GateAgentClose"          //代理关闭
)

type IGateModule interface {
	core.IModule
	GetOptions() (options IOptions)
	GetLocalRouteMgrComp() ILocalRouteMgrComp
	RegisterRemoteRoute(comId uint16, sId, sType string) (result string, err string)
	UnRegisterRemoteRoute(comId uint16, sType, sId string)
	RegisterLocalRoute(comId uint16, f func(session core.IUserSession, msg proto.IMessage) (code core.ErrorCode, err string))
	UnRegisterLocalRoute(comId uint16, f func(session core.IUserSession, msg proto.IMessage) (code core.ErrorCode, err string))
	OnRoute(a IAgent, msg proto.IMessage) (core.ErrorCode, error)
	Connect(a IAgent)
	DisConnect(a IAgent)
	CloseAgent(sId string) (result string, err string)
	SendMsg(sId string, msg proto.IMessage) (result int, err string)                //发送消息到用户
	SendMsgByGroup(aIds []string, msg proto.IMessage) (result []string, err string) //发送消息到用户组
	SendMsgByBroadcast(msg proto.IMessage) (result int, err string)                 //发送消息到全服
}

type IAgentMgrComp interface {
	core.IModuleComp
	Connect(a IAgent)
	DisConnect(a IAgent)
	SendMsg(aId string, msg proto.IMessage) (result int, err string)                //发送消息到用户
	SendMsgByGroup(aIds []string, msg proto.IMessage) (result []string, err string) //发送消息到用户组
	SendMsgByBroadcast(msg proto.IMessage) (result int, err string)                 //发送消息到全服
	Close(aId string) (result string, err string)
}

type ILocalRouteMgrComp interface {
	core.IModuleComp
	SetNewSession(f func(module IGateModule, data map[string]interface{}) (s core.IUserSession, err error))
	RegisterRoute(comId uint16, f func(session core.IUserSession, msg proto.IMessage) (code core.ErrorCode, err string))
	UnRegisterRoute(comId uint16, f func(session core.IUserSession, msg proto.IMessage) (code core.ErrorCode, err string))
	OnRoute(agent IAgent, msg proto.IMessage) (code core.ErrorCode, err error)
}

type IRemoteRouteMgrComp interface {
	core.IModuleComp
	SetNewSession(f func(service base.IClusterService, data map[string]interface{}) (s core.IUserSession, err error))
	RegisterRoute(comId uint16, sId, sType string) (result string, err string)
	UnRegisterRoute(comId uint16, sType, sId string)
	OnRoute(agent IAgent, msg proto.IMessage) (code core.ErrorCode, err error)
}

type ICustomRouteComp interface {
	core.IModuleComp
	RegisterRoute(route core.CustomRoute, msgs map[uint16][]uint16) (result string, err string)
	RegisterRouteFunc(route core.CustomRoute, f func(a IAgent, msg proto.IMessage))
	OnRoute(agent IAgent, msg proto.IMessage) (code core.ErrorCode, err error)
}

type IConn interface {
	Read(b []byte) (n int, err error)
	Write(b []byte) (n int, err error)
	RemoteAddr() net.Addr
	Close()
}

type IAgent interface {
	Id() string
	IP() string
	GetSessionData() map[string]interface{}
	OnInit(gate IGateModule, coon IConn, agent IAgent) (err error)
	WriteMsg(msg proto.IMessage) (err error)
	OnRecover(msg proto.IMessage)
	RevNum() int64
	SendNum() int64
	IsClosed() bool
	OnRun()
	OnClose() //主动关闭接口
	Destory() //关闭销毁
}
