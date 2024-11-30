package tcppool

import (
	"context"
	"fmt"
	"net"
	"sync"

	"github.com/liwei1dao/lego/core"
	"github.com/liwei1dao/lego/sys/log"
	"github.com/liwei1dao/lego/sys/lrpc/protocol"
)

func newsys(options *Options) (cpool *TcpConnPool, err error) {
	cpool = &TcpConnPool{
		options: options,
		clients: make(map[string]*Client),
		msgpack: make(chan protocol.IMessage, 100),
	}
	return
}

type TcpConnPool struct {
	options     *Options
	doneChan    chan struct{}
	clientMapMu sync.RWMutex
	clients     map[string]*Client
	msgpack     chan protocol.IMessage
}

func (this *TcpConnPool) Options() *Options {
	return this.options
}

func (this *TcpConnPool) Start() (err error) {
	var (
		ln net.Listener
	)
	if ln, err = net.Listen("tcp", fmt.Sprintf(":%d", this.options.Prot)); err != nil {
		this.options.Log.Errorf("err:%v", err)
	}
	go this.serveListener(ln)
	this.options.Log.Debug("TcpConnPool Start Listen !", log.Field{Key: "Endpoints", Value: this.options})
	return
}

func (this *TcpConnPool) GetClient(node *core.ServiceNode) (client *Client, err error) {
	var (
		ok   bool
		conn net.Conn
	)
	this.clientMapMu.RLock()
	client, ok = this.clients[node.Path()]
	this.clientMapMu.RUnlock()
	if !ok {
		if conn, err = net.DialTimeout("tcp", node.Addr(), this.options.ConnectionTimeout); err != nil {
			this.options.Log.Error("TcpConnPool GetClient Dial Err!", log.Field{Key: "add", Value: node.Addr}, log.Field{Key: "err", Value: err.Error()})
			return
		}
		if client, err = this.createClient(conn, node); err != nil {
			this.options.Log.Error("TcpConnPool createClient Err!", log.Field{Key: "err", Value: err.Error()})
			return
		}
	}
	return
}

// 创建远程连接客户端
func (this *TcpConnPool) createClient(conn net.Conn, node core.IServiceNode) (client *Client, err error) {
	if client, err = newClient(this, conn); err != nil {
		this.options.Log.Errorln(err)
		return
	}
	if err = this.ShakehandsRequest(context.Background(), client); err != nil {
		this.options.Log.Errorln(err)
		return
	}
	this.AddClient(client, node)
	return
}

func (this *TcpConnPool) serveListener(ln net.Listener) error {
	for {
		conn, e := ln.Accept()
		if e != nil {
			return e
		}
		if tc, ok := conn.(*net.TCPConn); ok {
			if this.options.KeepAlivePeriod > 0 {
				tc.SetKeepAlive(true)
				tc.SetKeepAlivePeriod(this.options.KeepAlivePeriod)
				tc.SetLinger(10)
			}
		}
		if _, err := newClient(this, conn); err != nil {
			this.options.Log.Error("newClient Err!", log.Field{Key: "err", Value: err.Error()})
		}
	}
}

func (this *TcpConnPool) ShakehandsRequest(ctx context.Context, client *Client) (err error) {
	return
}

func (this *TcpConnPool) AddClient(client *Client, node core.IServiceNode) (err error) {
	var (
		ok bool
	)
	this.clientMapMu.RLock()
	_, ok = this.clients[node.Path()]
	this.clientMapMu.RUnlock()
	if !ok {
		this.options.Log.Debug("AddClient Succ!", log.Field{Key: "node", Value: node})
		client.SetServiceNode(node)
		this.clientMapMu.Lock()
		this.clients[client.ServiceNode().Path()] = client
		this.clientMapMu.Unlock()
		client.Start()
	} else {
		err = fmt.Errorf("%v client already exists", node)
		this.options.Log.Errorln(err)
	}
	return
}

func (this *TcpConnPool) CloseClient(path string) (err error) {
	var (
		client *Client
		ok     bool
	)
	this.clientMapMu.RLock()
	client, ok = this.clients[path]
	this.clientMapMu.RUnlock()
	if ok {
		this.clientMapMu.Lock()
		delete(this.clients, path)
		this.clientMapMu.Unlock()
		err = client.Close()
	}
	return
}

func (this *TcpConnPool) Close() (err error) {
	
	return
}

func (this *TcpConnPool) In() chan<- protocol.IMessage {
	return this.msgpack
}
func (this *TcpConnPool) Out() <-chan protocol.IMessage {
	return this.msgpack
}
