package kafka

import (
	"sync"

	"github.com/liwei1dao/lego/core"
	"github.com/liwei1dao/lego/sys/log"
	lcore "github.com/liwei1dao/lego/sys/rpcl/core"
)

func NewKafkaConnPool(sys lcore.ISys, log log.ILogger, config *lcore.Config) (cpool *KafkaConnPool, err error) {
	cpool = &KafkaConnPool{
		sys: sys,
		log: log,
	}
	cpool.service, err = newService(sys, config)
	return
}

type KafkaConnPool struct {
	sys         lcore.ISys
	log         log.ILogger
	service     *Service
	clientMapMu sync.RWMutex
	clients     map[string]*Client
}

func (this *KafkaConnPool) Start() (err error) {

	return
}

func (this *KafkaConnPool) GetClient(node *core.ServiceNode) (client lcore.IConnClient, err error) {
	return
}
func (this *KafkaConnPool) Close() (err error) {

	return
}
