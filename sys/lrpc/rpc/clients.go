package rpc

import (
	"context"
	"errors"
	"strings"
	"sync"

	"github.com/liwei1dao/lego/sys/log"
	"github.com/smallnest/rpcx/client"
	"github.com/smallnest/rpcx/share"
)

// 异步返回结构
type Call struct {
	ServicePath   string
	ServiceMethod string
	Metadata      map[string]string
	ResMetadata   map[string]string
	Args          interface{} //请求参数
	Reply         interface{} //返回参数
	Error         error       //错误信息
	Done          chan *Call
}

func (this *Call) done() {
	select {
	case this.Done <- this:
	default:
		log.Debugf("rpc: 由于 Done chan 容量不足而丢弃呼叫回复!")
	}
}

// 客户端列表
type Clients struct {
	options *Options
	mu      sync.RWMutex
	clients map[string]*Client
}

// 异步执行 异步返回
func (this *Clients) Go(ctx context.Context, servicePath string, serviceMethod string, req interface{}, reply interface{}) (call *Call, err error) { //异步调用 异步返回

	return
}

// 获取目标客户端
func (this *Clients) getclient(ctx *context.Context, servicePath string) (c *Client, err error) {
	if servicePath == "" {
		err = errors.New("servicePath no cant null")
		return
	}
	var (
		spath []string
		d     client.ServiceDiscovery
		ok    bool
	)
	spath = strings.Split(servicePath, "/")

	this.mu.RLock()
	c, ok = this.clients[spath[0]]
	this.mu.RUnlock()
	if !ok {
		if d, err = discovery.NewEtcdDiscovery(clusterTag, spath[0], this.options.ETCDServers, false, nil); err != nil {
			return
		}
		c = client.NewBidirectionalXClient(spath[0], client.Failfast, client.RoundRobin, d, client.DefaultOption, this.msgChan)
		cluster.Mu.Lock()
		cluster.clients[spath[0]] = c
		cluster.Mu.Unlock()
		c.GetPlugins().Add(this)
		// if this.options.RpcxStartType == RpcxStartByClient && this.options.AutoConnect {
		c.SetSelector(newSelector(this.options.Log, clusterTag, this.UpdateServer))
		// } else {
		// 	c.SetSelector(newSelector(this.options.Log, clusterTag, nil))
		// }
	}

	*ctx = context.WithValue(*ctx, share.ReqMetaDataKey, map[string]string{
		ServiceClusterTag: this.options.ServiceNode.Tag(),
		CallRoutRulesKey:  servicePath,
		ServiceAddrKey:    "tcp@" + this.options.ServiceNode.Addr(),
		ServiceMetaKey:    this.metadata,
	})
	return
}
