package rpc

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/liwei1dao/lego/core"
	"github.com/liwei1dao/lego/sys/discovery"
	"github.com/liwei1dao/lego/sys/log"
	"github.com/liwei1dao/lego/sys/rpc/rpccore"
	"github.com/smallnest/rpcx/share"
)

var (
	// ErrXClientShutdown xclient is shutdown.
	ErrXClientShutdown = errors.New("xClient is shut down")
	// ErrXClientNoServer selector can't found one server.
	ErrXClientNoServer = errors.New("can not found any server")
	// ErrServerUnavailable selected server is unavailable.
	ErrServerUnavailable = errors.New("selected server is unavailable")
)

func NewLClient(servicePath string, discovery discovery.IDiscovery, options *Options) (client *LClient) {
	client = &LClient{
		options:     options,
		servicePath: servicePath,
		discovery:   discovery,
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

type LClient struct {
	options     *Options
	servicePath string
	discovery   discovery.IDiscovery
	selector    rpccore.ISelector
	isShutdown  bool
	auth        string
	cpool       rpccore.IConnPool
	mu          sync.RWMutex
	servers     map[string]core.IServiceNode
	ch          chan []*discovery.KVPair
}

func (this *LClient) ServiceNode() core.IServiceNode {
	return this.options.ServiceNode
}

func (this *LClient) watch(ch chan []*discovery.KVPair) {
	for pairs := range ch {
		sort.Slice(pairs, func(i, j int) bool {
			return strings.Compare(pairs[i].Key, pairs[j].Key) <= 0
		})
		servers := make(map[string]core.IServiceNode, len(pairs))
		for _, p := range pairs {
			servers[p.Key], _ = core.NewServiceNode(p.Value)
		}
		this.mu.Lock()
		filterByStateAndGroup(servers)
		this.servers = servers

		if this.selector != nil {
			this.selector.UpdateServer(servers)
		}
		this.mu.Unlock()
	}
}

// 同步执行
func (this *LClient) Call(ctx context.Context, serviceMethod string, req interface{}, reply interface{}) (err error) { //同步调用 等待结果
	seq := new(uint64)
	ctx = rpccore.WithValue(ctx, rpccore.CallSeqKey, seq)
	stime := time.Now()
	// this.options.Log.Debug("Call Start", log.Field{Key: "this.servicePath", Value: this.servicePath}, log.Field{Key: "serviceMethod", Value: serviceMethod}, log.Field{Key: "req", Value: req})
	defer func() {
		this.options.Log.Debug("RPC Call",
			log.Field{Key: "t", Value: time.Since(stime).Milliseconds()},
			log.Field{Key: "this.servicePath", Value: this.servicePath},
			log.Field{Key: "serviceMethod", Value: serviceMethod},
			log.Field{Key: "req", Value: req},
			log.Field{Key: "reply", Value: reply},
		)
	}()
	var call *MessageCall
	call, err = this.call(ctx, serviceMethod, req, reply)
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
func (this *LClient) Go(ctx context.Context, serviceMethod string, req interface{}, reply interface{}, done chan *MessageCall) (call *MessageCall, err error) { //异步调用 异步返回
	if this.isShutdown {
		return nil, ErrXClientShutdown
	}

	if this.auth != "" {
		metadata := ctx.Value(rpccore.ReqMetaDataKey)
		if metadata == nil {
			metadata = map[string]string{}
			ctx = context.WithValue(ctx, share.ReqMetaDataKey, metadata)
		}
		m := metadata.(map[string]string)
		m[share.AuthKey] = this.auth
	}

	ctx = setServerTimeout(ctx)

	if share.Trace {
		log.Debugf("select a client for %s.%s, req: %+v in case of xclient Go", this.servicePath, serviceMethod, req)
	}
	_, client, err := this.selectClient(ctx, this.servicePath, serviceMethod, req)
	if err != nil {
		return nil, err
	}
	if share.Trace {
		log.Debugf("selected a client %s for %s.%s, args: %+v in case of xclient Go", client.ServiceNode().Addr(), this.servicePath, serviceMethod, req)
	}
	return client.Go(ctx, serviceMethod, req, reply, done), nil
}

func (this *LClient) Broadcast(ctx context.Context, serviceMethod string, args interface{}) (err error) {
	return
}

func (this *LClient) Close() (err error) {
	return
}
func (this *LClient) selectClient(ctx context.Context, servicePath, serviceMethod string, args interface{}) (k string, client IClinet, err error) {
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

func setServerTimeout(ctx context.Context) context.Context {
	if deadline, ok := ctx.Deadline(); ok {
		metadata := ctx.Value(share.ReqMetaDataKey)
		if metadata == nil {
			metadata = map[string]string{}
			ctx = context.WithValue(ctx, share.ReqMetaDataKey, metadata)
		}
		m := metadata.(map[string]string)
		m[share.ServerTimeout] = fmt.Sprintf("%d", time.Until(deadline).Milliseconds())
	}

	return ctx
}
