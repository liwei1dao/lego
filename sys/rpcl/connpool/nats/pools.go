package nats

import (
	"bytes"
	"context"
	"sync"
	"time"

	"github.com/liwei1dao/lego/core"
	lcore "github.com/liwei1dao/lego/sys/rpcl/core"
	"github.com/liwei1dao/lego/sys/rpcl/protocol"
	"github.com/nats-io/nats.go"
)

func NewKafkaConnPool(sys lcore.ISys, config *lcore.Config) (cpool *NatsConnPool, err error) {
	cpool = &NatsConnPool{
		sys:    sys,
		config: config,
	}
	err = cpool.init()
	return
}

type NatsConnPool struct {
	sys         lcore.ISys
	config      *lcore.Config
	conn        *nats.Conn
	subs        *nats.Subscription
	clientMapMu sync.RWMutex
	clients     map[string]lcore.IConnClient
}

func (this *NatsConnPool) init() (err error) {
	if this.conn, err = nats.Connect(this.config.Endpoints[0]); err != nil {
		return
	}
	this.subs, err = this.conn.SubscribeSync(this.sys.ServiceNode().GetNodePath())
	return
}
func (this *NatsConnPool) GetClient(node *core.ServiceNode) (client lcore.IConnClient, err error) {
	var (
		ok bool
	)
	this.clientMapMu.RLock()
	client, ok = this.clients[node.Addr]
	this.clientMapMu.RUnlock()
	if !ok {
		if client, err = newClient(this, this.config, node); err == nil {
			if err = this.sys.ShakehandsRequest(context.Background(), client); err == nil {
				this.clientMapMu.Lock()
				this.clients[node.Addr] = client
				this.clientMapMu.Unlock()
				client.Start()
			}
		}
	}
	return
}
func (this *NatsConnPool) Start() (err error) {
	go this.run()
	return
}
func (this *NatsConnPool) Close() (err error) {
	return
}
func (this *NatsConnPool) CloseClient(node *core.ServiceNode) (err error) {
	var (
		client lcore.IConnClient
		ok     bool
	)
	this.clientMapMu.RLock()
	client, ok = this.clients[node.Addr]
	this.clientMapMu.RUnlock()
	if ok {
		this.clientMapMu.Lock()
		delete(this.clients, node.Addr)
		this.clientMapMu.Unlock()
		err = client.Close()
	}
	return
}
func (this *NatsConnPool) run() {
	var (
		err     error
		m       *nats.Msg
		message *protocol.Message
		client  lcore.IConnClient
	)
locp:
	for {
		m, err = this.subs.NextMsg(time.Minute)
		if err != nil && err == nats.ErrTimeout {
			continue
		} else if err != nil {
			break locp
		}
		if message, err = protocol.Read(bytes.NewReader(m.Data)); err != nil {
			this.sys.Errorf("err:%v", err)
			continue
		}
		if client, err = this.GetClient(message.From()); err != nil {
			this.sys.Errorf("err:%v", err)
			continue
		}
		go this.sys.Handle(client, message)
	}
	this.sys.Warnf("connpool nats service run() exit!:%v", err)
}
