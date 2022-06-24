package kafka

import (
	"github.com/liwei1dao/lego/sys/kafka"
	"github.com/liwei1dao/lego/sys/rpcl/core"
)

type Service struct {
	sys   core.ISys
	kafka kafka.ISys
}

func (this *Service) run() {
	for v := range this.kafka.Consumer_Messages() {
		this.sys.Debugf("Client receive:%v", v)
	}
}
