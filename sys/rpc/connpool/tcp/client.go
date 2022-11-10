package tcp

import (
	"bufio"
	"crypto/tls"
	"net"
	"runtime"
	"sync"
	"sync/atomic"
	"time"

	"github.com/liwei1dao/lego/core"
	"github.com/liwei1dao/lego/sys/rpc/protocol"
	"github.com/liwei1dao/lego/sys/rpc/rpccore"
)

func newClient(pool *TcpConnPool, config *rpccore.Config, conn net.Conn) (client *Client, err error) {
	client = &Client{
		pool:   pool,
		config: config,
		conn:   conn,
		hbeat:  0,
		state:  int32(rpccore.ClientShakeHands),
	}
	go client.serveConn()
	return
}

type Client struct {
	pool        *TcpConnPool
	config      *rpccore.Config
	node        *core.ServiceNode
	conn        net.Conn
	closeSignal chan bool
	hbeat       int32 //心跳发出次数
	state       int32 //状态 0 未启动 1 运行 2 关闭中
	wg          sync.WaitGroup
}

func (this *Client) ServiceNode() *core.ServiceNode {
	return this.node
}

func (this *Client) SetServiceNode(node *core.ServiceNode) {
	this.node = node
}

func (this *Client) ResetHbeat() {
	atomic.StoreInt32(&this.hbeat, 0)
}
func (this *Client) State() rpccore.ClientState {
	return rpccore.ClientState(atomic.LoadInt32(&this.state))
}

func (this *Client) Start() {
	atomic.StoreInt32(&this.state, int32(rpccore.ClientRuning))
	return
}

func (this *Client) Write(msg []byte) (err error) {
	_, err = this.conn.Write(msg)
	if err != nil {
		this.pool.log.Errorf("send msg err:%v", err)
	}
	return
}

func (this *Client) Close() (err error) {
	if atomic.CompareAndSwapInt32(&this.state, int32(rpccore.ClientRuning), int32(rpccore.ClientClose)) {
		this.conn.Close()
	}
	return
}

func (this *Client) serveConn() {
	defer func() {
		if err := recover(); err != nil {
			const size = 64 << 10
			buf := make([]byte, size)
			ss := runtime.Stack(buf, false)
			if ss > size {
				ss = size
			}
			buf = buf[:ss]
			this.pool.log.Errorf("serving %s panic error: %s, stack:\n %s", this.conn.RemoteAddr(), err, buf)
		}
		this.pool.log.Debugf("server closed conn: %v", this.conn.RemoteAddr().String())
	}()

	if tlsConn, ok := this.conn.(*tls.Conn); ok {
		if d := this.config.ReadTimeout; d != 0 {
			this.conn.SetReadDeadline(time.Now().Add(d))
		}
		if d := this.config.WriteTimeout; d != 0 {
			this.conn.SetWriteDeadline(time.Now().Add(d))
		}
		if err := tlsConn.Handshake(); err != nil {
			this.pool.log.Errorf("rpcx: TLS handshake error from %s: %v", this.conn.RemoteAddr(), err)
			return
		}
	}

	r := bufio.NewReaderSize(this.conn, ReaderBuffsize)
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
		go this.pool.sys.Handle(this, req)
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
			if err = this.Write(this.pool.sys.Heartbeat()); err != nil {
				this.pool.log.Errorf("err:%v", err)
				go this.pool.CloseClient(this.node)
			}
			if atomic.LoadInt32(&this.hbeat) > 3 {
				this.pool.log.Errorf("heartbeat exception !")
				go this.pool.CloseClient(this.node)
			}
		case <-this.closeSignal:
			break locp
		}
	}
	timer.Stop()
	this.wg.Done()
}
