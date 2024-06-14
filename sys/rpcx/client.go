package rpcx

import (
	"context"
	"errors"
	"fmt"
	"net"
	"reflect"
	"runtime"
	"strings"
	"sync"
	"time"
	"unicode"
	"unicode/utf8"

	consulclient "github.com/rpcxio/rpcx-consul/client"
	"github.com/smallnest/rpcx/client"
	"github.com/smallnest/rpcx/protocol"
	"github.com/smallnest/rpcx/share"
)

func newClient(options *Options) (sys *Client, err error) {
	sys = &Client{
		options:        options,
		metadata:       fmt.Sprintf("stag=%s&stype=%s&sid=%s&version=%s&addr=%s", options.ServiceTag, options.ServiceType, options.ServiceId, options.ServiceVersion, "tcp@"+options.ServiceAddr),
		clusterClients: make(map[string]*clusterClients),
		conns:          make(map[string]net.Conn),
		serviceMap:     make(map[string]*service),
		msgChan:        make(chan *protocol.Message, 1000),
	}
	return
}

type clusterClients struct {
	Mu      sync.RWMutex
	clients map[string]client.XClient //其他集群客户端
}

type Client struct {
	options        *Options
	metadata       string
	writeTimeout   time.Duration
	AsyncWrite     bool
	clusterMu      sync.RWMutex
	clusterClients map[string]*clusterClients //其他集群客户端
	connsMapMu     sync.RWMutex
	conns          map[string]net.Conn
	// connectMapMu   sync.RWMutex
	// connecting     map[string]struct{}
	serviceMapMu sync.RWMutex
	serviceMap   map[string]*service
	msgChan      chan *protocol.Message // 接收rpcXServer推送消息
}

// DoMessage 服务端消息处理
func (this *Client) DoMessage() {
	for msg := range this.msgChan {
		go func(req *protocol.Message) {
			if req.ServicePath != "" && req.ServiceMethod != "" {
				this.options.Log.Debugf("DoMessage :%v", req)
				addr, ok := req.Metadata[ServiceAddrKey]
				if !ok {
					this.options.Log.Errorf("Metadata no found ServiceAddrKey!")
					return
				}
				conn, ok := this.conns[addr]
				if !ok {
					this.options.Log.Errorf("no found conn addr:%s", addr)
					return
				}
				res, _ := this.handleRequest(context.Background(), req)
				this.sendResponse(conn, req, res)
			}
		}(msg)
	}
}

// 启动RPC 服务 接收消息处理
func (this *Client) Start() (err error) {
	go this.DoMessage()
	return
}

// 停止RPC 服务
func (this *Client) Stop() (err error) {
	this.clusterMu.Lock()
	for _, v := range this.clusterClients {
		v.Mu.Lock()
		for _, v1 := range v.clients {
			v1.Close()
		}
		v.Mu.Unlock()
	}
	this.clusterMu.RUnlock()
	close(this.msgChan) //关闭消息处理
	return
}

// 获取服务集群列表
func (this *Client) GetServiceTags() []string {
	this.clusterMu.RLock()
	tags := make([]string, len(this.clusterClients))
	n := 0
	for k, _ := range this.clusterClients {
		tags[n] = k
		n++
	}
	this.clusterMu.RUnlock()
	return tags
}

// 注册Rpc 服务
func (this *Client) RegisterFunction(fn interface{}) (err error) {
	_, err = this.registerFunction(this.options.ServiceType, fn, "", false)
	if err != nil {
		return err
	}
	return
}

// 注册Rpc 服务
func (this *Client) RegisterFunctionName(name string, fn interface{}) (err error) {
	_, err = this.registerFunction(this.options.ServiceType, fn, name, true)
	if err != nil {
		return err
	}
	return
}

// 注销 暂不处理
func (this *Client) UnregisterAll() (err error) {
	return nil
}

// 同步调用
func (this *Client) Call(ctx context.Context, servicePath string, serviceMethod string, args interface{}, reply interface{}) (err error) {
	var (
		_client client.XClient
	)
	if _client, err = this.getclient(&ctx, this.options.ServiceTag, servicePath); err != nil {
		return
	}
	err = _client.Call(ctx, serviceMethod, args, reply)
	return
}

// 异步调用
func (this *Client) Go(ctx context.Context, servicePath string, serviceMethod string, args interface{}, reply interface{}, done chan *client.Call) (call *client.Call, err error) {
	var (
		_client client.XClient
	)
	if _client, err = this.getclient(&ctx, this.options.ServiceTag, servicePath); err != nil {
		return
	}
	return _client.Go(ctx, string(serviceMethod), args, reply, done)
}

// 异步调用
func (this *Client) Broadcast(ctx context.Context, servicePath string, serviceMethod string, args interface{}, reply interface{}) (err error) {
	var (
		_client client.XClient
	)
	if _client, err = this.getclient(&ctx, this.options.ServiceTag, servicePath); err != nil {
		return
	}
	err = _client.Broadcast(ctx, serviceMethod, args, reply)
	return
}

// 跨集群 同步调用
func (this *Client) AcrossClusterCall(ctx context.Context, clusterTag string, servicePath string, serviceMethod string, args interface{}, reply interface{}) (err error) {
	var (
		_client client.XClient
	)
	if _client, err = this.getclient(&ctx, clusterTag, servicePath); err != nil {
		return
	}
	err = _client.Call(ctx, serviceMethod, args, reply)
	return
}

// 跨集群 异步调用
func (this *Client) AcrossClusterGo(ctx context.Context, clusterTag string, servicePath string, serviceMethod string, args interface{}, reply interface{}, done chan *client.Call) (call *client.Call, err error) {
	var (
		_client client.XClient
	)
	if _client, err = this.getclient(&ctx, clusterTag, servicePath); err != nil {
		return
	}

	return _client.Go(ctx, string(serviceMethod), args, reply, done)
}

// 跨集群 广播
func (this *Client) AcrossClusterBroadcast(ctx context.Context, clusterTag string, servicePath string, serviceMethod string, args interface{}, reply interface{}) (err error) {
	var (
		_client client.XClient
	)
	if _client, err = this.getclient(&ctx, clusterTag, servicePath); err != nil {
		return
	}
	return _client.Broadcast(ctx, serviceMethod, args, reply)
}

func (this *Client) ClusterBroadcast(ctx context.Context, servicePath string, serviceMethod string, args interface{}, reply interface{}) (err error) {
	if servicePath == "" {
		err = errors.New("servicePath no cant null")
		return
	}
	var (
		spath   []string
		clients []client.XClient
	)
	spath = strings.Split(servicePath, "/")
	ctx = context.WithValue(ctx, share.ReqMetaDataKey, map[string]string{
		ServiceClusterTag: this.options.ServiceTag,
		CallRoutRulesKey:  servicePath,
		ServiceAddrKey:    "tcp@" + this.options.ServiceAddr,
		ServiceMetaKey:    this.metadata,
	})
	clients = make([]client.XClient, 0)
	this.clusterMu.RLock()
	for _, v := range this.clusterClients {
		v.Mu.RLock()
		if _client, ok := v.clients[spath[0]]; ok {
			clients = append(clients, _client)
		}
		v.Mu.RUnlock()
	}
	this.clusterMu.RUnlock()
	l := len(clients)
	if l > 0 {
		done := make(chan error, l)
		for _, v := range clients {
			go func(c client.XClient) {
				done <- c.Broadcast(ctx, serviceMethod, args, reply)
			}(v)
		}
		timeout := time.NewTimer(time.Minute)
	check:
		for {
			select {
			case err = <-done:
				l--
				if l == 0 || err != nil { // all returns or some one returns an error
					break check
				}
			case <-timeout.C:
				err = errors.New(("timeout"))
				break check
			}
		}
		timeout.Stop()
	} else {
		err = errors.New("on found any service")
	}
	return
}

// 监控服务发现，发现没有连接上的额服务端 就连接上去
func (this *Client) UpdateServer(servers map[string]*ServiceNode) {
	for _, v := range servers {
		this.clusterMu.RLock()
		cluster, ok := this.clusterClients[v.ServiceTag]
		this.clusterMu.RUnlock()
		if ok {
			cluster.Mu.RLock()
			_, ok = cluster.clients[v.ServiceType]
			cluster.Mu.RUnlock()
			if ok {
				continue
			}
		}
		//没有建立客户端 主动发起握手
		if err := this.Call(context.Background(), fmt.Sprintf("%s/%s", v.ServiceType, v.ServiceId), RpcX_ShakeHands, &ServiceNode{
			ServiceTag:  this.options.ServiceTag,
			ServiceId:   this.options.ServiceId,
			ServiceType: this.options.ServiceType,
			ServiceAddr: this.options.ServiceAddr},
			&ServiceNode{}); err != nil {
			this.options.Log.Errorf("ShakeHands new node addr:%s err:%v", v.ServiceAddr, err)
		} else {
			this.options.Log.Debugf("UpdateServer addr:%s ", v.ServiceAddr)
		}
	}
}

// 监控连接建立
func (this *Client) ClientConnected(conn net.Conn) (net.Conn, error) {
	addr := "tcp@" + conn.RemoteAddr().String()
	this.connsMapMu.Lock()
	this.conns[addr] = conn
	this.connsMapMu.Unlock()
	this.options.Log.Debugf("ClientConnected addr:%v", addr)
	return conn, nil
}

// 监听连接关闭
func (this *Client) ClientConnectionClose(conn net.Conn) error {
	addr := "tcp@" + conn.RemoteAddr().String()
	this.connsMapMu.Lock()
	delete(this.conns, addr)
	this.connsMapMu.Unlock()
	this.options.Log.Debugf("ClientConnectionClose addr:%v", addr)
	return nil
}

// 获取目标客户端
func (this *Client) getclient(ctx *context.Context, clusterTag string, servicePath string) (c client.XClient, err error) {
	if servicePath == "" {
		err = errors.New("servicePath no cant null")
		return
	}
	var (
		spath   []string
		cluster *clusterClients
		d       *consulclient.ConsulDiscovery
		ok      bool
	)
	spath = strings.Split(servicePath, "/")
	this.clusterMu.RLock()
	cluster, ok = this.clusterClients[clusterTag]
	this.clusterMu.RUnlock()
	if !ok {
		cluster = &clusterClients{clients: make(map[string]client.XClient)}
		this.clusterMu.Lock()
		this.clusterClients[clusterTag] = cluster
		this.clusterMu.Unlock()
	}
	cluster.Mu.RLock()
	c, ok = cluster.clients[spath[0]]
	cluster.Mu.RUnlock()
	if !ok {
		if d, err = consulclient.NewConsulDiscovery(clusterTag, spath[0], this.options.ConsulServers, nil); err != nil {
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
		ServiceClusterTag: this.options.ServiceTag,
		CallRoutRulesKey:  servicePath,
		ServiceAddrKey:    "tcp@" + this.options.ServiceAddr,
		ServiceMetaKey:    this.metadata,
	})
	return
}

// 注册服务方法
func (this *Client) registerFunction(servicePath string, fn interface{}, name string, useName bool) (string, error) {
	this.serviceMapMu.Lock()
	defer this.serviceMapMu.Unlock()

	ss := this.serviceMap[servicePath]
	if ss == nil {
		ss = new(service)
		ss.name = servicePath
		ss.function = make(map[string]*functionType)
	}

	f, ok := fn.(reflect.Value)
	if !ok {
		f = reflect.ValueOf(fn)
	}
	if f.Kind() != reflect.Func {
		return "", errors.New("function must be func or bound method")
	}

	fname := runtime.FuncForPC(reflect.Indirect(f).Pointer()).Name()
	if fname != "" {
		i := strings.LastIndex(fname, ".")
		if i >= 0 {
			fname = fname[i+1:]
		}
	}
	if useName {
		fname = name
	}
	if fname == "" {
		errorStr := "rpcx.registerFunction: no func name for type " + f.Type().String()
		this.options.Log.Errorf(errorStr)
		return fname, errors.New(errorStr)
	}

	t := f.Type()
	if t.NumIn() != 3 {
		return fname, fmt.Errorf("rpcx.registerFunction: has wrong number of ins: %s", f.Type().String())
	}
	if t.NumOut() != 1 {
		return fname, fmt.Errorf("rpcx.registerFunction: has wrong number of outs: %s", f.Type().String())
	}

	// First arg must be context.Context
	ctxType := t.In(0)
	if !ctxType.Implements(typeOfContext) {
		return fname, fmt.Errorf("function %s must use context as  the first parameter", f.Type().String())
	}

	argType := t.In(1)
	if !isExportedOrBuiltinType(argType) {
		return fname, fmt.Errorf("function %s parameter type not exported: %v", f.Type().String(), argType)
	}

	replyType := t.In(2)
	if replyType.Kind() != reflect.Ptr {
		return fname, fmt.Errorf("function %s reply type not a pointer: %s", f.Type().String(), replyType)
	}
	if !isExportedOrBuiltinType(replyType) {
		return fname, fmt.Errorf("function %s reply type not exported: %v", f.Type().String(), replyType)
	}

	// The return type of the method must be error.
	if returnType := t.Out(0); returnType != typeOfError {
		return fname, fmt.Errorf("function %s returns %s, not error", f.Type().String(), returnType.String())
	}

	// Install the methods
	ss.function[fname] = &functionType{fn: f, ArgType: argType, ReplyType: replyType}
	this.serviceMap[servicePath] = ss

	// init pool for reflect.Type of args and reply
	reflectTypePools.Init(argType)
	reflectTypePools.Init(replyType)
	return fname, nil
}

// 处理远程服务请求
func (this *Client) handleRequest(ctx context.Context, req *protocol.Message) (res *protocol.Message, err error) {
	serviceName := req.ServicePath
	methodName := req.ServiceMethod

	res = req.Clone()
	res.SetMessageType(protocol.Response)
	this.serviceMapMu.RLock()
	service := this.serviceMap[serviceName]
	this.serviceMapMu.RUnlock()
	if service == nil {
		err = errors.New("rpcx: can't find service " + serviceName)
		return handleError(res, err)
	}
	mtype := service.method[methodName]
	if mtype == nil {
		if service.function[methodName] != nil { // check raw functions
			return this.handleRequestForFunction(ctx, req)
		}
		err = errors.New("rpcx: can't find method " + methodName)
		return handleError(res, err)
	}
	// get a argv object from object pool
	argv := reflectTypePools.Get(mtype.ArgType)

	codec := share.Codecs[req.SerializeType()]
	if codec == nil {
		err = fmt.Errorf("can not find codec for %d", req.SerializeType())
		return handleError(res, err)
	}
	err = codec.Decode(req.Payload, argv)
	if err != nil {
		return handleError(res, err)
	}

	// and get a reply object from object pool
	replyv := reflectTypePools.Get(mtype.ReplyType)
	if mtype.ArgType.Kind() != reflect.Ptr {
		err = service.call(ctx, mtype, reflect.ValueOf(argv).Elem(), reflect.ValueOf(replyv))
	} else {
		err = service.call(ctx, mtype, reflect.ValueOf(argv), reflect.ValueOf(replyv))
	}
	reflectTypePools.Put(mtype.ArgType, argv)
	if err != nil {
		if replyv != nil {
			data, err := codec.Encode(replyv)
			// return reply to object pool
			reflectTypePools.Put(mtype.ReplyType, replyv)
			if err != nil {
				return handleError(res, err)
			}
			res.Payload = data
		}
		return handleError(res, err)
	}

	if !req.IsOneway() {
		data, err := codec.Encode(replyv)
		// return reply to object pool
		reflectTypePools.Put(mtype.ReplyType, replyv)
		if err != nil {
			return handleError(res, err)
		}
		res.Payload = data
	} else if replyv != nil {
		reflectTypePools.Put(mtype.ReplyType, replyv)
	}
	this.options.Log.Debugf("server called service %+v for an request %+v", service, req)
	return res, nil
}

// 处理远程服务请求 for 方法
func (this *Client) handleRequestForFunction(ctx context.Context, req *protocol.Message) (res *protocol.Message, err error) {
	res = req.Clone()

	res.SetMessageType(protocol.Response)

	serviceName := req.ServicePath
	methodName := req.ServiceMethod
	this.serviceMapMu.RLock()
	service := this.serviceMap[serviceName]
	this.serviceMapMu.RUnlock()
	if service == nil {
		err = errors.New("rpcx: can't find service  for func raw function")
		return handleError(res, err)
	}
	mtype := service.function[methodName]
	if mtype == nil {
		err = errors.New("rpcx: can't find method " + methodName)
		return handleError(res, err)
	}

	argv := reflectTypePools.Get(mtype.ArgType)

	codec := share.Codecs[req.SerializeType()]
	if codec == nil {
		err = fmt.Errorf("can not find codec for %d", req.SerializeType())
		return handleError(res, err)
	}

	err = codec.Decode(req.Payload, argv)
	if err != nil {
		return handleError(res, err)
	}

	replyv := reflectTypePools.Get(mtype.ReplyType)

	if mtype.ArgType.Kind() != reflect.Ptr {
		err = service.callForFunction(ctx, mtype, reflect.ValueOf(argv).Elem(), reflect.ValueOf(replyv))
	} else {
		err = service.callForFunction(ctx, mtype, reflect.ValueOf(argv), reflect.ValueOf(replyv))
	}

	reflectTypePools.Put(mtype.ArgType, argv)

	if err != nil {
		reflectTypePools.Put(mtype.ReplyType, replyv)
		return handleError(res, err)
	}

	if !req.IsOneway() {
		data, err := codec.Encode(replyv)
		reflectTypePools.Put(mtype.ReplyType, replyv)
		if err != nil {
			return handleError(res, err)
		}
		res.Payload = data
	} else if replyv != nil {
		reflectTypePools.Put(mtype.ReplyType, replyv)
	}

	return res, nil
}

// 发送远程服务调用 回应消息
func (this *Client) sendResponse(conn net.Conn, req, res *protocol.Message) {
	if len(res.Payload) > 1024 && req.CompressType() != protocol.None {
		res.SetCompressType(req.CompressType())
	}
	data := res.EncodeSlicePointer()
	if this.writeTimeout != 0 {
		conn.SetWriteDeadline(time.Now().Add(this.writeTimeout))
	}
	conn.Write(*data)
	protocol.PutData(data)
}

// 请求错误 封装回应消息
func handleError(res *protocol.Message, err error) (*protocol.Message, error) {
	res.SetMessageStatusType(protocol.Error)
	if res.Metadata == nil {
		res.Metadata = make(map[string]string)
	}
	res.Metadata[protocol.ServiceError] = err.Error()
	return res, err
}

// 服务注册 类型判断
func isExportedOrBuiltinType(t reflect.Type) bool {
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	// PkgPath will be non-empty even for an exported type,
	// so we need to check the type name as well.
	return isExported(t.Name()) || t.PkgPath() == ""
}

func isExported(name string) bool {
	rune, _ := utf8.DecodeRuneInString(name)
	return unicode.IsUpper(rune)
}
