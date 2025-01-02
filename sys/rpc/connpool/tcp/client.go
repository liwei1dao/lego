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
	"github.com/liwei1dao/lego/sys/rpc"
)

func newClient(pool *TcpConnPool, options *Options, conn net.Conn) (client *Client, err error) {
	client = &Client{
		pool:    pool,
		options: options,
		conn:    conn,
		hbeat:   0,
		state:   int32(rpc.ClientShakeHands),
	}
	go client.serveConn()
	return
}

type Client struct {
	options     *Options
	pool        *TcpConnPool
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
func (this *Client) State() rpc.ClientState {
	return rpc.ClientState(atomic.LoadInt32(&this.state))
}

func (this *Client) Start() {
	atomic.StoreInt32(&this.state, int32(rpc.ClientRuning))
	return
}

func (this *Client) Write(msg []byte) (err error) {
	_, err = this.conn.Write(msg)
	if err != nil {
		this.options.Log.Errorf("send msg err:%v", err)
	}
	return
}

func (this *Client) Close() (err error) {
	if atomic.CompareAndSwapInt32(&this.state, int32(rpc.ClientRuning), int32(rpc.ClientClose)) {
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
			this.options.Log.Errorf("serving %s panic error: %s, stack:\n %s", this.conn.RemoteAddr(), err, buf)
		}
		this.options.Log.Debugf("server closed conn: %v", this.conn.RemoteAddr().String())
	}()

	if tlsConn, ok := this.conn.(*tls.Conn); ok {
		if this.options.ReadTimeout != 0 {
			this.conn.SetReadDeadline(time.Now().Add(time.Duration(this.options.ReadTimeout) * time.Second))
		}
		if this.options.WriteTimeout != 0 {
			this.conn.SetWriteDeadline(time.Now().Add(time.Duration(this.options.WriteTimeout) * time.Second))
		}
		if err := tlsConn.Handshake(); err != nil {
			this.options.Log.Errorf("rpcx: TLS handshake error from %s: %v", this.conn.RemoteAddr(), err)
			return
		}
	}

	r := bufio.NewReaderSize(this.conn, ReaderBuffsize)
locp:
	for {
		t0 := time.Now()
		if this.options.ReadTimeout > 0 {
			this.conn.SetReadDeadline(t0.Add(time.Duration(this.options.ReadTimeout) * time.Second))
		}
		if err := this.pool.host.ReadProtocol(this, r); err != nil {
			go this.pool.CloseClient(this.node)
			break locp
		}
	}
}

func (this *Client) heartbeat() {
	var (
		timer *time.Ticker
		err   error
	)
	timer = time.NewTicker(time.Duration(this.options.KeepAlivePeriod) * time.Second)
locp:
	for {
		select {
		case <-timer.C:
			if err = this.Write(this.pool.host.Heartbeat()); err != nil {
				this.options.Log.Errorf("err:%v", err)
				go this.pool.CloseClient(this.node)
			}
			if atomic.LoadInt32(&this.hbeat) > 3 {
				this.options.Log.Errorf("heartbeat exception !")
				go this.pool.CloseClient(this.node)
			}
		case <-this.closeSignal:
			break locp
		}
	}
	timer.Stop()
	this.wg.Done()
}
