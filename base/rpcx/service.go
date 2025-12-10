package rpcx

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/liwei1dao/lego/core"
	"github.com/liwei1dao/lego/core/cbase"
	"github.com/liwei1dao/lego/sys/cron"
	"github.com/liwei1dao/lego/sys/event"
	"github.com/liwei1dao/lego/sys/log"

	"github.com/rcrowley/go-metrics"
	etcdclient "github.com/rpcxio/rpcx-etcd/client"
	"github.com/rpcxio/rpcx-etcd/serverplugin"
	"github.com/smallnest/rpcx/client"
	"github.com/smallnest/rpcx/server"
)

/*
XRPC 服务对象
*/
type RPCXService struct {
	cbase.ServiceBase
	option      *Options
	serviceNode *core.ServiceNode
	server      *server.Server
	mu          sync.RWMutex
	clients     map[string]client.XClient //其他集群客户端
}

func (this *RPCXService) GetTag() string {
	return this.option.Setting.Tag
}
func (this *RPCXService) GetId() string {
	return this.option.Setting.Id
}
func (this *RPCXService) GetType() string {
	return this.option.Setting.Type
}

func (this *RPCXService) GetVersion() string {
	return this.option.Version
}
func (this *RPCXService) GetSettings() core.ServiceSttings {
	return this.option.Setting
}

func (this *RPCXService) Configure(option ...Option) {
	this.option = newOptions(option...)
	this.serviceNode, _ = core.NewServiceNode(fmt.Sprintf("%s=%s&%s=%s&%s=%s&%s=%s",
		core.ServiceNode_Tag, this.option.Setting.Tag,
		core.ServiceNode_Type, this.option.Setting.Type,
		core.ServiceNode_Id, this.option.Setting.Id,
		core.ServiceNode_Addr, this.option.RPCXConfig.Addr,
	))
}

func (this *RPCXService) Init(service core.IService) (err error) {
	err = this.ServiceBase.Init(service)
	this.server = server.NewServer()
	this.clients = make(map[string]client.XClient)
	r := &serverplugin.EtcdV3RegisterPlugin{
		ServiceAddress: "tcp@" + this.option.RPCXConfig.Addr,
		EtcdServers:    this.option.RPCXConfig.ETCDServers, // options.ETCDServers,
		BasePath:       this.option.Setting.Tag,
		Metrics:        metrics.NewRegistry(),
		UpdateInterval: time.Duration(this.option.RPCXConfig.UpdateInterval) * time.Second,
	}
	if err = r.Start(); err != nil {
		return
	}
	this.server.Plugins.Add(r)
	go func() {
		if err = this.server.Serve("tcp", fmt.Sprintf(":%d", this.option.RPCXConfig.Port)); err != nil {
			log.Warnf("rpcx server exit:%v", err)
		}
	}()
	return
}

func (this *RPCXService) InitSys() {
	if err := log.OnInit(this.option.Setting.Sys["log"]); err != nil {
		panic(fmt.Sprintf("sys log Init err:%v", err))
	} else {
		log.Infof("sys log Init success !")
	}
	if err := event.OnInit(this.option.Setting.Sys["event"]); err != nil {
		log.Panicf(fmt.Sprintf("sys event Init err:%v", err))
	} else {
		log.Infof("sys event Init success !")
	}
}

func (this *RPCXService) Destroy() (err error) {
	cron.Close()
	err = this.ServiceBase.Destroy()
	return
}

// 注册服务方法 自定义方法名称
func (this *RPCXService) Register(name string, fn interface{}) (err error) {
	err = this.server.RegisterFunctionName(this.GetType(), name, fn, this.serviceNode.Value())
	return
}

// 同步调用
func (this *RPCXService) RpcCall(ctx context.Context, servicePath string, serviceMethod string, args interface{}, reply interface{}) (err error) {
	var (
		_client client.XClient
	)
	if _client, err = this.getclient(&ctx, servicePath); err != nil {
		return
	}
	err = _client.Call(ctx, serviceMethod, args, reply)
	return
}

// 异步调用
func (this *RPCXService) RpcGo(ctx context.Context, servicePath string, serviceMethod string, args interface{}, reply interface{}, done chan *client.Call) (call *client.Call, err error) {
	var (
		_client client.XClient
	)
	if _client, err = this.getclient(&ctx, servicePath); err != nil {
		return
	}
	return _client.Go(ctx, string(serviceMethod), args, reply, done)
}

// 异步调用
func (this *RPCXService) RpcBroadcast(ctx context.Context, servicePath string, serviceMethod string, args interface{}, reply interface{}) (err error) {
	var (
		_client client.XClient
	)
	if _client, err = this.getclient(&ctx, servicePath); err != nil {
		return
	}
	err = _client.Broadcast(ctx, serviceMethod, args, reply)
	return
}

// 获取目标客户端
func (this *RPCXService) getclient(ctx *context.Context, servicePath string) (c client.XClient, err error) {
	if servicePath == "" {
		err = errors.New("service no cant null")
		return
	}
	var (
		ok bool
		d  client.ServiceDiscovery
	)
	this.mu.RLock()
	c, ok = this.clients[servicePath]
	this.mu.RUnlock()
	if !ok {
		if d, err = etcdclient.NewEtcdV3Discovery(this.option.Setting.Tag, servicePath, this.option.RPCXConfig.ETCDServers, false, nil); err != nil {
			return
		}
		c = client.NewBidirectionalXClient(servicePath, client.Failfast, client.RoundRobin, d, client.DefaultOption, nil)
		this.mu.Lock()
		this.clients[servicePath] = c
		this.mu.Unlock()
	}
	return
}
