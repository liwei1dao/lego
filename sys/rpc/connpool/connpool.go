package connpool

import (
	"fmt"

	"github.com/liwei1dao/lego/sys/log"
	"github.com/liwei1dao/lego/sys/rpc/connpool/kafka"
	"github.com/liwei1dao/lego/sys/rpc/connpool/tcp"
	"github.com/liwei1dao/lego/sys/rpc/rpccore"
)

//创建连接池对象
func NewConnPool(sys rpccore.ISys, log log.ILogger, config *rpccore.Config) (comm rpccore.IConnPool, err error) {
	switch config.ConnectType {
	case rpccore.Tcp:
		comm, err = tcp.NewTcpConnPool(sys, log, config)
		break
	case rpccore.Kafka:
		comm, err = kafka.NewKafkaConnPool(sys, log, config)
		break
	default:
		err = fmt.Errorf("not support ConnectType:%v", config.ConnectType)
	}
	return
}
