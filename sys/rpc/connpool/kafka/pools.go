package kafka

import (
	"sync"

	"github.com/liwei1dao/lego/core"
	"github.com/liwei1dao/lego/sys/log"
	"github.com/liwei1dao/lego/sys/rpc/rpccore"
)

func NewKafkaConnPool(sys rpccore.ISys, log log.ILogger, config *rpccore.Config) (cpool *KafkaConnPool, err error) {
	cpool = &KafkaConnPool{
		sys: sys,
		log: log,
	}
	cpool.service, err = newService(sys, config)
	return
}

type KafkaConnPool struct {
	sys         rpccore.ISys
	log         log.ILogger
	service     *Service
	clientMapMu sync.RWMutex
	clients     map[string]*Client
}

func (this *KafkaConnPool) Start() (err error) {

	return
}

func (this *KafkaConnPool) GetClient(node *core.ServiceNode) (client rpccore.IConnClient, err error) {
	return
}
func (this *KafkaConnPool) AddClient(client rpccore.IConnClient, node *core.ServiceNode) (err error) {
	return
}
func (this *KafkaConnPool) Close() (err error) {

	return
}
