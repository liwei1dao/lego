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
)
