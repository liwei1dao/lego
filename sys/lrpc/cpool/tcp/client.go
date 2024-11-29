package tcp

import (
	"bufio"
	"crypto/tls"
	"net"
	"sync"
	"sync/atomic"
	"time"

	"github.com/liwei1dao/lego"
	"github.com/liwei1dao/lego/core"
	"github.com/liwei1dao/lego/sys/log"
	lcore "github.com/liwei1dao/lego/sys/lrpc/core"
	"github.com/liwei1dao/lego/sys/lrpc/cpool"
	"github.com/liwei1dao/lego/sys/lrpc/protocol"
)

func newClient(pool cpool.ICPool, config *cpool.Config, conn net.Conn) (client *Client, err error) {
	client = &Client{
		pool:   pool,
		config: config,
		conn:   conn,
		hbeat:  0,
		state:  int32(lcore.ClientShakeHands),
	}
	go client.serveConn()
	return
}

type Client struct {
	pool        cpool.ICPool
	config      *cpool.Config
	node        core.IServiceNode
	conn        net.Conn
	closeSignal chan bool
	hbeat       int32 //心跳发出次数
	state       int32 //状态 0 未启动 1 运行 2 关闭中
	wg          sync.WaitGroup
}

func (this *Client) ServiceNode() core.IServiceNode {
	return this.node
}

func (this *Client) SetServiceNode(node core.IServiceNode) {
	this.node = node
}

func (this *Client) ResetHbeat() {
	atomic.StoreInt32(&this.hbeat, 0)
}
func (this *Client) State() lcore.ClientState {
	return lcore.ClientState(atomic.LoadInt32(&this.state))
}

func (this *Client) Start() {
	atomic.StoreInt32(&this.state, int32(lcore.ClientRuning))
	return
}

func (this *Client) Write(msg []byte) (err error) {
	_, err = this.conn.Write(msg)
	if err != nil {
		log.Errorf("send msg err:%v", err)
	}
	return
}

func (this *Client) Close() (err error) {
	if atomic.CompareAndSwapInt32(&this.state, int32(lcore.ClientRuning), int32(lcore.ClientClose)) {
		this.conn.Close()
	}
	return
}

func (this *Client) serveConn() {
	defer lego.Recover("lrpc.serveConn")
	if tlsConn, ok := this.conn.(*tls.Conn); ok {
		if d := this.config.ReadTimeout; d != 0 {
			this.conn.SetReadDeadline(time.Now().Add(d))
		}
		if d := this.config.WriteTimeout; d != 0 {
			this.conn.SetWriteDeadline(time.Now().Add(d))
		}
		if err := tlsConn.Handshake(); err != nil {
			log.Errorf("rpcx: TLS handshake error from %s: %v", this.conn.RemoteAddr(), err)
			return
		}
	}

	r := bufio.NewReaderSize(this.conn, cpool.ReaderBuffsize)
locp:
	for {
		t0 := time.Now()
		if this.config.ReadTimeout > 0 {
			this.conn.SetReadDeadline(t0.Add(this.config.ReadTimeout))
		}
		req := protocol.GetPooledMsg()
		err := req.Decode(r)
		if err != nil {
			go this.pool.CloseClient(this.node)
			break locp
		}
		go this.pool.Handle(this, req)
	}
}

func (this *Client) heartbeat() {
	var (
		timer *time.Ticker
		err   error
	)
	timer = time.NewTicker(this.config.KeepAlivePeriod)
locp:
	for {
		select {
		case <-timer.C:
			if err = this.Write(this.pool.Heartbeat()); err != nil {
				log.Errorf("err:%v", err)
				go this.pool.CloseClient(this.node)
			}
			if atomic.LoadInt32(&this.hbeat) > 3 {
				log.Errorf("heartbeat exception !")
				go this.pool.CloseClient(this.node)
			}
		case <-this.closeSignal:
			break locp
		}
	}
	timer.Stop()
	this.wg.Done()
}
