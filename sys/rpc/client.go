package rpc

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/liwei1dao/lego/core"
	"github.com/liwei1dao/lego/sys/discovery"
	"github.com/liwei1dao/lego/sys/log"
	"github.com/liwei1dao/lego/sys/rpc/protocol"
	"github.com/liwei1dao/lego/sys/rpc/rpccore"
)

func NewXClient(servicePath string, discovery discovery.IDiscovery, option *Options) (client *Client, err error) {
	client = &Client{
		options:   option,
		discovery: discovery,
		pending:   make(map[uint64]*MessageCall),
	}
	pairs := discovery.GetServices()
	sort.Slice(pairs, func(i, j int) bool {
		return strings.Compare(pairs[i].Key, pairs[j].Key) <= 0
	})
	servers := make(map[string]core.IServiceNode, len(pairs))
	for _, p := range pairs {
		servers[p.Key], _ = core.NewServiceNode(p.Value)
	}
	filterByStateAndGroup(servers)
	client.servers = servers
	ch := client.discovery.WatchService()
	if ch != nil {
		client.ch = ch
		go client.watch(ch)
	}
	return
}

type Client struct {
	options      *Options
	discovery    discovery.IDiscovery
	selector     rpccore.ISelector
	cpool        rpccore.IConnPool
	mu           sync.RWMutex
	servers      map[string]core.IServiceNode
	ch           chan []*discovery.KVPair
	pendingmutex sync.Mutex
	seq          uint64
	pending      map[uint64]*MessageCall
}

func (this *Client) ServiceNode() core.IServiceNode {
	return this.options.ServiceNode
}

func (c *Client) watch(ch chan []*discovery.KVPair) {
	for pairs := range ch {
		sort.Slice(pairs, func(i, j int) bool {
			return strings.Compare(pairs[i].Key, pairs[j].Key) <= 0
		})
		servers := make(map[string]core.IServiceNode, len(pairs))
		for _, p := range pairs {
			servers[p.Key], _ = core.NewServiceNode(p.Value)
		}
		c.mu.Lock()
		filterByStateAndGroup(servers)
		c.servers = servers

		if c.selector != nil {
			c.selector.UpdateServer(servers)
		}
		c.mu.Unlock()
	}
}

// 同步执行
func (this *Client) Call(ctx context.Context, servicePath string, serviceMethod string, req interface{}, reply interface{}) (err error) { //同步调用 等待结果
	seq := new(uint64)
	ctx = rpccore.WithValue(ctx, rpccore.CallSeqKey, seq)
	stime := time.Now()
	// this.options.Log.Debug("Call Start", log.Field{Key: "servicePath", Value: servicePath}, log.Field{Key: "serviceMethod", Value: serviceMethod}, log.Field{Key: "req", Value: req})
	defer func() {
		this.options.Log.Debug("RPC Call",
			log.Field{Key: "t", Value: time.Since(stime).Milliseconds()},
			log.Field{Key: "servicePath", Value: servicePath},
			log.Field{Key: "serviceMethod", Value: serviceMethod},
			log.Field{Key: "req", Value: req},
			log.Field{Key: "reply", Value: reply},
		)
	}()
	var call *MessageCall
	call, err = this.call(ctx, servicePath, serviceMethod, req, reply)
	select {
	case <-ctx.Done(): // cancel by context
		this.pendingmutex.Lock()
		call := this.pending[*seq]
		delete(this.pending, *seq)
		this.pendingmutex.Unlock()
		if call != nil {
			call.Error = ctx.Err()
			call.done(this.options.Log)
		}
		return ctx.Err()
	case call := <-call.Done:
		err = call.Error
	}
	return
}

// 异步执行 异步返回
func (this *Client) Go(ctx context.Context, servicePath string, serviceMethod string, req interface{}, reply interface{}) (call *MessageCall, err error) { //异步调用 异步返回
	seq := new(uint64)
	ctx = rpccore.WithValue(ctx, rpccore.CallSeqKey, seq)
	stime := time.Now()
	// this.options.Log.Debug("Go start!", log.Field{Key: "servicePath", Value: servicePath}, log.Field{Key: "serviceMethod", Value: serviceMethod}, log.Field{Key: "req", Value: req})
	defer func() {
		this.options.Log.Debug("RPC Go",
			log.Field{Key: "t", Value: time.Since(stime).Milliseconds()},
			log.Field{Key: "servicePath", Value: servicePath},
			log.Field{Key: "serviceMethod", Value: serviceMethod},
			log.Field{Key: "req", Value: req},
			log.Field{Key: "reply", Value: reply},
		)
	}()
	call, err = this.call(ctx, servicePath, serviceMethod, req, reply)
	return
}

func (this *Client) Broadcast(ctx context.Context, servicePath string, serviceMethod string, args interface{}) (err error) {
	return
}

// 执行远程服务---------------------------------------------------------------------------------------

func filterByStateAndGroup(servers map[string]core.IServiceNode) {
	for k, v := range servers {
		if v.State() == "inactive" {
			delete(servers, k)
		}
	}
}

func (this *Client) call(ctx context.Context, servicePath string, serviceMethod string, args interface{}, reply interface{}) (call *MessageCall, err error) {
	call = new(MessageCall)
	call.ServicePath = servicePath
	call.ServiceMethod = serviceMethod
	call.Done = make(chan *MessageCall, 10)
	call.Args = args
	call.Reply = reply
	this.pendingmutex.Lock()
	seq := this.seq
	this.seq++
	this.pending[seq] = call
	this.pendingmutex.Unlock()
	if cseq, ok := ctx.Value(rpccore.CallSeqKey).(*uint64); ok {
		*cseq = seq
	}
	var client rpccore.IConnClient
	if client, err = this.getclient(ctx, servicePath); err != nil {
		return
	}
	err = this.send(client, call, seq)
	return
}

// 获取请求消息对象
func (this *Client) getMessage(serviceMethod string, args interface{}, reply interface{}) (call *MessageCall, req *protocol.Message, err error) {
	var data []byte
	call = new(MessageCall)
	call.ServiceMethod = serviceMethod
	call.Done = make(chan *MessageCall, 10)
	call.Args = this.ServiceNode() //自己发起握手 需要传递本服务的节点信息
	call.Reply = reply
	req = protocol.GetPooledMsg()
	req.SetVersion(this.options.ProtoVersion)
	if call.Reply != nil {
		this.pendingmutex.Lock()
		seq := this.seq
		this.seq++
		this.pending[seq] = call
		this.pendingmutex.Unlock()
		req.SetSeq(seq)
		req.SetOneway(true)
	} else {
		req.SetOneway(false)
	}

	req.SetServiceMethod(call.ServiceMethod)
	req.SetFrom(this.ServiceNode())
	req.SetMessageType(rpccore.Request)
	req.SetSerializeType(this.options.SerializeType)
	data, err = codecs[this.options.SerializeType].Marshal(call.Args)
	if err != nil {
		return
	}
	if len(data) > 1024 && this.options.CompressType != rpccore.CompressNone {
		req.SetCompressType(this.options.CompressType)
	}
	req.SetPayload(data)
	return
}

func (this *Client) getclient(ctx context.Context, servicePath string) (client rpccore.IConnClient, err error) {
	nodes := this.selector.Select(ctx, servicePath)
	if nodes == nil || len(nodes) == 0 {
		err = fmt.Errorf("no found any node:%s", servicePath)
		this.options.Log.Errorln(err)
		return
	}
	if client, err = this.cpool.GetClient(nodes[0]); err != nil {
		this.options.Log.Errorln(err)
		return
	}
	return
}

func (this *Client) send(client rpccore.IConnClient, call *MessageCall, seq uint64) (err error) {
	var (
		data    []byte
		allData *[]byte
	)

	req := protocol.GetPooledMsg()
	req.SetVersion(this.options.ProtoVersion)
	req.SetMessageType(rpccore.Request)
	req.SetSeq(seq)
	if call.Reply == nil {
		req.SetOneway(true)
	}
	if call.Metadata != nil {
		req.SetMetadata(call.Metadata)
	}

	req.SetServiceMethod(call.ServiceMethod)
	req.SetFrom(this.options.ServiceNode)
	req.SetSerializeType(this.options.SerializeType)
	codec := codecs[this.options.SerializeType]
	if codec == nil {
		err = rpccore.ErrUnsupportedCodec
		return
	}
	data, err = codec.Marshal(call.Args)
	if err != nil {
		return
	}
	if len(data) > 1024 && this.options.CompressType != rpccore.CompressNone {
		req.SetCompressType(this.options.CompressType)
	}
	req.SetPayload(data)
	allData = req.EncodeSlicePointer()
	defer func() {
		protocol.PutData(allData)
		protocol.FreeMsg(req)
	}()
	err = client.Write(*allData)
	return
}
