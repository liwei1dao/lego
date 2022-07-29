package connpool

import (
	"fmt"

	"github.com/liwei1dao/lego/sys/rpcl/connpool/kafka"
	"github.com/liwei1dao/lego/sys/rpcl/connpool/tcp"
	lcore "github.com/liwei1dao/lego/sys/rpcl/core"
)

//创建连接池对象
func NewConnPool(sys lcore.ISys, config *lcore.Config) (comm lcore.IConnPool, err error) {
	switch config.ConnectType {
	case lcore.Tcp:
		comm, err = tcp.NewTcpConnPool(sys, config)
		break
	case lcore.Kafka:
		comm, err = kafka.NewKafkaConnPool(sys, config)
		break
	default:
		err = fmt.Errorf("not support ConnectType:%v", config.ConnectType)
	}
	return
}
