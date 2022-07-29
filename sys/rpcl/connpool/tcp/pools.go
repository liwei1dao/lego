package tcp

import (
	"context"
	"net"
	"sync"

	"github.com/liwei1dao/lego/core"
	lcore "github.com/liwei1dao/lego/sys/rpcl/core"
)

const (
	// ReaderBuffsize is used for bufio reader.
	ReaderBuffsize = 1024
	// WriterBuffsize is used for bufio writer.
	WriterBuffsize = 1024

	// WriteChanSize is used for response.
	WriteChanSize = 1024 * 1024
)

func NewTcpConnPool(sys lcore.ISys, config *lcore.Config) (cpool *TcpConnPool, err error) {
	cpool = &TcpConnPool{
		sys:    sys,
		config: config,
	}
	return
}

type TcpConnPool struct {
	sys         lcore.ISys
	config      *lcore.Config
	doneChan    chan struct{}
	clientMapMu sync.RWMutex
	clients     map[string]lcore.IConnClient
}

func (this *TcpConnPool) Start() (err error) {
	var (
		ln net.Listener
	)
	if ln, err = net.Listen("tcp", this.config.Endpoints[0]); err != nil {
		this.sys.Errorf("err:%v", err)
	}
	go this.serveListener(ln)
	return
}

func (this *TcpConnPool) GetClient(node *core.ServiceNode) (client lcore.IConnClient, err error) {
	var (
		ok   bool
		conn net.Conn
	)
	this.clientMapMu.RLock()
	client, ok = this.clients[node.GetNodePath()]
	this.clientMapMu.RUnlock()
	if !ok {
		if conn, err = net.DialTimeout("tcp", node.Addr, this.config.ConnectionTimeout); err == nil {
			client, err = this.createClient(conn)
		}
	}
	return
}

func (this *TcpConnPool) createClient(conn net.Conn) (client lcore.IConnClient, err error) {
	if client, err = newClient(this, this.config, conn); err == nil {
		if err = this.sys.ShakehandsRequest(context.Background(), client); err == nil {
			this.clientMapMu.Lock()
			this.clients[client.ServiceNode().GetNodePath()] = client
			this.clientMapMu.Unlock()
			client.Start()
		}
	}
	return
}

func (this *TcpConnPool) serveListener(ln net.Listener) error {
	for {
		conn, e := ln.Accept()
		if e != nil {
			return e
		}
		if tc, ok := conn.(*net.TCPConn); ok {
			if this.config.KeepAlivePeriod > 0 {
				tc.SetKeepAlive(true)
				tc.SetKeepAlivePeriod(this.config.KeepAlivePeriod)
				tc.SetLinger(10)
			}
		}
		this.createClient(conn)
	}
}

func (this *TcpConnPool) CloseClient(node *core.ServiceNode) (err error) {
	var (
		client lcore.IConnClient
		ok     bool
	)
	this.clientMapMu.RLock()
	client, ok = this.clients[node.GetNodePath()]
	this.clientMapMu.RUnlock()
	if ok {
		this.clientMapMu.Lock()
		delete(this.clients, node.GetNodePath())
		this.clientMapMu.Unlock()
		err = client.Close()
	}
	return
}

func (this *TcpConnPool) Close() (err error) {

	return
}
