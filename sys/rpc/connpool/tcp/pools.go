package tcp

import (
	"context"
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/liwei1dao/lego/core"
	"github.com/liwei1dao/lego/sys/log"
	"github.com/liwei1dao/lego/sys/rpc"
)

const (
	// ReaderBuffsize is used for bufio reader.
	ReaderBuffsize = 1024
	// WriterBuffsize is used for bufio writer.
	WriterBuffsize = 1024

	// WriteChanSize is used for response.
	WriteChanSize = 1024 * 1024
)

func newSys(options *Options) (cpool *TcpConnPool, err error) {
	cpool = &TcpConnPool{
		options: options,
		clients: make(map[string]rpc.IConnClient),
	}
	return
}

type TcpConnPool struct {
	options     *Options
	host        rpc.IBodyHost
	doneChan    chan struct{}
	clientMapMu sync.RWMutex
	clients     map[string]rpc.IConnClient
}

func (this *TcpConnPool) Start(host rpc.IBodyHost) (err error) {
	var (
		ln net.Listener
	)
	this.host = host
	if ln, err = net.Listen("tcp", this.options.ListterAddr); err != nil {
		this.options.Log.Errorf("err:%v", err)
	}
	go this.serveListener(ln)
	this.options.Log.Debug("TcpConnPool Start Listen !", log.Field{Key: "Endpoints", Value: this.options.ListterAddr})
	return
}

func (this *TcpConnPool) GetClient(node core.IServiceNode) (client rpc.IConnClient, err error) {
	var (
		ok   bool
		conn net.Conn
	)
	this.clientMapMu.RLock()
	client, ok = this.clients[node.Path()]
	this.clientMapMu.RUnlock()
	if !ok {
		if conn, err = net.DialTimeout("tcp", node.Addr(), time.Duration(this.options.ConnectionTimeout)*time.Second); err != nil {
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
func (this *TcpConnPool) createClient(conn net.Conn, node core.IServiceNode) (client rpc.IConnClient, err error) {
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

func (this *TcpConnPool) serveListener(ln net.Listener) {
	for {
		conn, e := ln.Accept()
		if e != nil {
			this.options.Log.Error("serveListener Err!", log.Field{Key: "err", Value: e.Error()})
			return
		}
		if tc, ok := conn.(*net.TCPConn); ok {
			if this.options.KeepAlivePeriod > 0 {
				tc.SetKeepAlive(true)
				tc.SetKeepAlivePeriod(time.Second * time.Duration(this.options.KeepAlivePeriod))
				tc.SetLinger(10)
			}
		}
		if _, err := newClient(this, this.options, conn); err != nil {
			this.options.Log.Error("newClient Err!", log.Field{Key: "err", Value: err.Error()})
		}
	}
}

func (this *TcpConnPool) AddClient(client rpc.IConnClient, node core.IServiceNode) (err error) {
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

func (this *TcpConnPool) CloseClient(node core.IServiceNode) (err error) {
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

func (this *TcpConnPool) Close() (err error) {

	return
}
