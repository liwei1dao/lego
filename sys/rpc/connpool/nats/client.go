package nats

import (
	"sync"
	"sync/atomic"
	"time"

	"github.com/liwei1dao/lego/core"
	"github.com/liwei1dao/lego/sys/rpc"
	"github.com/nats-io/nats.go"
)

func newClient(pool *NatsConnPool, options *Options, snode core.IServiceNode) (client *Client, err error) {
	client = &Client{
		pool:    pool,
		options: options,
		node:    snode,
		hbeat:   0,
		state:   0,
	}
	if client.conn, err = nats.Connect(options.NatsAddr); err != nil {
		return
	}
	return
}

type Client struct {
	pool        *NatsConnPool
	options     *Options
	node        core.IServiceNode
	conn        *nats.Conn
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

func (this *Client) State() rpc.ClientState {
	return rpc.ClientState(atomic.LoadInt32(&this.state))
}

func (this *Client) Start() {
	atomic.StoreInt32(&this.state, 1)
	this.wg.Add(1)
	go this.heartbeat()
	return
}

func (this *Client) Write(msg []byte) (err error) {
	err = this.conn.Publish(this.node.Path(), msg)
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
	timer = time.NewTicker(time.Duration(this.options.KeepAlivePeriod) * time.Second)
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
