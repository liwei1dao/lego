package kafka

import (
	"github.com/liwei1dao/lego/sys/kafka"
)

func newClient(pool *KafkaConnPool, servicePath string) (client *Client, err error) {
	return
}

type Client struct {
	pool  *KafkaConnPool
	kafka kafka.ISys
}

func (this *Client) run() {
	for v := range this.kafka.Consumer_Messages() {
		this.pool.log.Debugf("Client receive:%v", v)
	}
}
