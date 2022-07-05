package kafka

import (
	"context"

	"github.com/liwei1dao/lego/sys/kafka"
	"github.com/liwei1dao/lego/sys/rpcl/core"
)

func newClient(sys core.ISys, servicePath string) (client *Client, err error) {
	return
}

type Client struct {
	sys   core.ISys
	kafka kafka.ISys
}

func (this *Client) run() {
	for v := range this.kafka.Consumer_Messages() {
		this.sys.Debugf("Client receive:%v", v)
	}
}

func (this *Client) Go(ctx context.Context, serviceMethod string, args interface{}, reply interface{}, done chan *core.Call) *core.Call {
	return nil
}
