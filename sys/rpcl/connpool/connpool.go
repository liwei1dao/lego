package connpool

import (
	"github.com/liwei1dao/lego/sys/rpcl/connpool/kafka"
	gcore "github.com/liwei1dao/lego/sys/rpcl/core"
)

//创建连接池对象
func NewConnPool(sys gcore.ISys, config *gcore.Config) (comm gcore.IConnPool, err error) {
	switch config.ConnectType {
	case gcore.Kafka:
		comm, err = kafka.NewKafkaConnPool(sys, config)
		break
	}
	return
}
