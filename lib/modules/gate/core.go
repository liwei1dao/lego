package gate

import (
	"lego/core"
	"lego/sys/proto"
	"net"
)

const ( //Rpc

	Rpc_GateRouteRegister core.Rpc_Key = "GateRouteRegister" //网关路由注册
	Rpc_GateRoute         core.Rpc_Key = "GateRoute"         //网关路由
	Rpc_GateAgentsIsKeep  core.Rpc_Key = "GateAgentsIsKeep"  //校验代理是否还在
	RPC_GateAgentBuild    core.Rpc_Key = "GateAgentBuild"    //代理绑定
	RPC_GateAgentUnBuild  core.Rpc_Key = "GateAgentUnBuild"  //代理解绑
	RPC_GateAgentSendMsg  core.Rpc_Key = "GateAgentSendMsg"  //代理发送消息
	RPC_GateAgentClose    core.Rpc_Key = "GateAgentClose"    //代理关闭
)

type IGateModule interface {
	core.IModule
	//需重构处理  内部函数为重构代码
	RegisterLocalRoute(comId uint16, f func(session core.IUserSession, msg proto.IMessage) (code int, err string))
	UnRegisterLocalRoute(comId uint16, f func(session core.IUserSession, msg proto.IMessage))
	OnRoute(a IAgent, msg proto.IMessage) (code int, err string)
	Connect(a IAgent)
	DisConnect(a IAgent)
	CloseAgent(sId string) (result string, err string)
	SendMsg(sId string, msg proto.IMessage) (result int, err string)
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
	OnClose() 		//主动关闭接口
	OnCloseWait()   //主动关闭接口 等待链接销毁
	Destory() 		//关闭销毁
}
