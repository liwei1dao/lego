package m_comps

import (
	"github.com/liwei1dao/lego/core"
	"github.com/liwei1dao/lego/sys/proto"
)

type (
	IMComp_GateComp interface {
		core.IModuleComp
		ReceiveMsg(session core.IUserSession, msg proto.IMessage) (code core.ErrorCode, err string)
	}
	IMComp_GateCompOptions interface {
		core.IModuleOptions
		GetDefGoroutine() int     //网关工作池默认协程数
		GetGateMaxGoroutine() int //网关工作池最大协程数
		GetHandleTimeOut() int    //工作池处理超时时间
	}
)
