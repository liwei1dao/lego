package kafka

import (
	"sync"

	"github.com/liwei1dao/lego/core"
	gcore "github.com/liwei1dao/lego/sys/rpcl/core"
)

func NewKafkaConnPool(sys gcore.ISys, config *gcore.Config) (cpool *KafkaConnPool, err error) {
	cpool = &KafkaConnPool{
		sys: sys,
	}
	cpool.service, err = newService(sys, config)
	return
}

type KafkaConnPool struct {
	sys         gcore.ISys
	service     *Service
	clientMapMu sync.RWMutex
	clients     map[string]*Client
}

func (this *KafkaConnPool) Start() (err error) {

	return
}

func (this *KafkaConnPool) GetClient(server *core.ServiceNode) (client gcore.IConnClient, err error) {
	return
}

func (this *KafkaConnPool) Close() (err error) {

	return
}
