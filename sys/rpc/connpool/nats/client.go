package nats

import (
	"sync"
	"sync/atomic"
	"time"

	"github.com/liwei1dao/lego/core"
	"github.com/liwei1dao/lego/sys/rpc/rpccore"
	"github.com/nats-io/nats.go"
)

func newClient(pool *NatsConnPool, config *rpccore.Config, snode *core.ServiceNode) (client *Client, err error) {
	client = &Client{
		pool:   pool,
		config: config,
		node:   snode,
		hbeat:  0,
		state:  0,
	}
	if client.conn, err = nats.Connect(config.Endpoints[0]); err != nil {
		return
	}
	return
}

type Client struct {
	pool        *NatsConnPool
	config      *rpccore.Config
	node        *core.ServiceNode
	conn        *nats.Conn
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

func (this *Client) State() rpccore.ClientState {
	return rpccore.ClientState(atomic.LoadInt32(&this.state))
}

func (this *Client) Start() {
	atomic.StoreInt32(&this.state, 1)
	this.wg.Add(1)
	go this.heartbeat()
	return
}

func (this *Client) Write(msg []byte) (err error) {
	err = this.conn.Publish(this.node.GetNodePath(), msg)
	if err != nil {
		this.pool.log.Errorf("send msg err:%v", err)
	}
	err = this.conn.Flush()
	return
}

func (this *Client) Close() (err error) {
	this.conn.Close()
	if !atomic.CompareAndSwapInt32(&this.state, 1, 2) {
		return
	} else {
		this.closeSignal <- true
		this.wg.Wait()
	}
	return
}

func (this *Client) ResetHbeat() {
	atomic.StoreInt32(&this.hbeat, 0)
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
