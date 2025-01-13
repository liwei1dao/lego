package nats

import (
	"bytes"
	"context"
	"net"
	"sync"
	"time"

	"github.com/liwei1dao/lego/core"
	"github.com/liwei1dao/lego/sys/log"
	"github.com/liwei1dao/lego/sys/rpc"
	"github.com/liwei1dao/lego/sys/rpc/protocol"

	"github.com/nats-io/nats.go"
)

func newSys(options *Options) (cpool *NatsConnPool, err error) {
	cpool = &NatsConnPool{
		options: options,
		clients: make(map[string]rpc.IConnClient),
	}
	return
}

type NatsConnPool struct {
	options     *Options
	host        rpc.IBodyHost
	conn        *nats.Conn
	subs        *nats.Subscription
	clientMapMu sync.RWMutex
	clients     map[string]rpc.IConnClient
}

func (this *NatsConnPool) init() (err error) {
	if this.conn, err = nats.Connect(this.options.NatsAddr); err != nil {
		return
	}
	this.subs, err = this.conn.SubscribeSync(this.options.ServiceNode().Path())
	return
}

// 创建远程连接客户端
func (this *NatsConnPool) createClient(conn net.Conn, node core.IServiceNode) (client rpc.IConnClient, err error) {
	if client, err = newClient(this, this.options, conn); err != nil {
		this.options.Log.Errorln(err)
		return
	}
	if err = this.host.ShakehandsRequest(context.Background(), client); err != nil {
		this.options.Log.Errorln(err)
		return
	}
	this.AddClient(client, node)
	return
}

func (this *NatsConnPool) GetClient(node core.IServiceNode) (client rpc.IConnClient, err error) {
	var (
		ok bool
	)
	this.clientMapMu.RLock()
	client, ok = this.clients[node.Path()]
	this.clientMapMu.RUnlock()
	if !ok {
		if client, err = newClient(this, this.options); err != nil {
			this.options.Log.Errorln(err)
			return
		}
		client.SetServiceNode(node)
		if err = this.host.ShakehandsRequest(context.Background(), client); err != nil {
			this.options.Log.Errorln(err)
			return
		}
		this.options.Log.Debug("CreateClient Succ!", log.Field{Key: "node", Value: node})
		this.clientMapMu.Lock()
		this.clients[node.Path()] = client
		this.clientMapMu.Unlock()
		client.Start()
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
func (this *NatsConnPool) CloseClient(node core.IServiceNode) (err error) {
	var (
		client rpc.IConnClient
		ok     bool
	)
	this.clientMapMu.RLock()
	client, ok = this.clients[node.Path()]
	this.clientMapMu.RUnlock()
	if ok {
		this.clientMapMu.Lock()
		delete(this.clients, node.Path())
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
		client  rpc.IConnClient
	)
locp:
	for {
		m, err = this.subs.NextMsg(time.Minute)
		if err != nil && err == nats.ErrTimeout {
			continue
		} else if err != nil {
			break locp
		}

		if client, err = this.GetClient(message.From()); err != nil {
			this.options.Log.Errorf("err:%v", err)
			continue
		}
		go this.host.ReadProtocol(client, bytes.NewReader(m.Data))
	}
	this.options.Log.Warnf("connpool nats service run() exit!:%v", err)
}
