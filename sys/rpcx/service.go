package rpcx

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/rcrowley/go-metrics"
	"github.com/rpcxio/rpcx-consul/serverplugin"
	"github.com/smallnest/rpcx/client"
	"github.com/smallnest/rpcx/protocol"
	"github.com/smallnest/rpcx/server"
	"github.com/smallnest/rpcx/share"
)

func newService(options *Options) (sys *Service, err error) {
	sys = &Service{
		options:    options,
		metadata:   fmt.Sprintf("stag=%s&stype=%s&sid=%s&version=%s&addr=%s", options.ServiceTag, options.ServiceType, options.ServiceId, options.ServiceVersion, "tcp@"+options.ServiceAddr),
		server:     server.NewServer(),
		selectors:  make(map[string]ISelector),
		clients:    make(map[string]net.Conn),
		clientmeta: make(map[string]string),
		pending:    make(map[uint64]*client.Call),
	}

	r := &serverplugin.ConsulRegisterPlugin{
		ServiceAddress: "tcp@" + options.ServiceAddr,
		ConsulServers:  options.ConsulServers,
		BasePath:       options.ServiceTag,
		Metrics:        metrics.NewRegistry(),
		UpdateInterval: time.Minute,
	}
	if err = r.Start(); err != nil {
		return
	}
	sys.server.Plugins.Add(r)
	sys.server.Plugins.Add(sys)
	sys.RegisterFunctionName(RpcX_ShakeHands, sys.RpcxShakeHands) //注册握手函数
	return
}

type Service struct {
	options       *Options
	metadata      string
	server        *server.Server
	selectormutex sync.RWMutex
	selectors     map[string]ISelector
	clientmutex   sync.RWMutex
	clients       map[string]net.Conn
	clientmeta    map[string]string
	mutex         sync.Mutex // protects following
	seq           uint64
	pending       map[uint64]*client.Call
}

// RPC 服务启动
func (this *Service) Start() (err error) {
	go func() {
		if err = this.server.Serve("tcp", this.options.ServiceAddr); err != nil {
			this.options.Log.Warnf("rpcx server exit:%v", err)
		}
	}()
	return
}

// 服务停止
func (this *Service) Stop() (err error) {
	err = this.server.Close()
	return
}

// 获取服务集群列表
func (this *Service) GetServiceTags() []string {
	this.selectormutex.RLock()
	tags := make([]string, len(this.selectors))
	n := 0
	for k, _ := range this.selectors {
		tags[n] = k
		n++
	}
	this.selectormutex.RUnlock()
	return tags
}

// 注册RPC 服务
func (this *Service) RegisterFunction(fn interface{}) (err error) {
	err = this.server.RegisterFunction(this.options.ServiceType, fn, this.metadata)
	return
}

// 注册RPC 服务
func (this *Service) RegisterFunctionName(name string, fn interface{}) (err error) {
	err = this.server.RegisterFunctionName(this.options.ServiceType, name, fn, this.metadata)
	return
}

// 注销 暂时不处理
func (this *Service) UnregisterAll() (err error) {
	// err = this.server.UnregisterAll()
	return
}

// 同步调用远程服务
func (this *Service) Call(ctx context.Context, servicePath string, serviceMethod string, args interface{}, reply interface{}) (err error) {
	var (
		done *client.Call
		conn net.Conn
	)
	seq := new(uint64)
	ctx = context.WithValue(ctx, seqKey{}, seq)
	if conn, done, err = this.call(&ctx, this.options.ServiceTag, servicePath, serviceMethod, args, reply, make(chan *client.Call, 1)); err != nil {
		return
	}
	select {
	case <-ctx.Done(): // cancel by context
		this.mutex.Lock()
		call := this.pending[*seq]
		delete(this.pending, *seq)
		this.mutex.Unlock()
		if call != nil {
			call.Error = ctx.Err()
			call.Done <- call
		}
		return ctx.Err()
	case call := <-done.Done:
		err = call.Error
		meta := ctx.Value(share.ResMetaDataKey)
		if meta != nil && len(call.ResMetadata) > 0 {
			resMeta := meta.(map[string]string)
			for k, v := range call.ResMetadata {
				resMeta[k] = v
			}
			resMeta[share.ServerAddress] = conn.RemoteAddr().String()
		}
	}
	return
}

// 广播调用
func (this *Service) Broadcast(ctx context.Context, servicePath string, serviceMethod string, args interface{}, reply interface{}) (err error) {
	err = this.broadcast(ctx, this.options.ServiceTag, servicePath, serviceMethod, args, reply)
	return
}

// 异步调用 远程服务
func (this *Service) Go(ctx context.Context, servicePath string, serviceMethod string, args interface{}, reply interface{}, done chan *client.Call) (_call *client.Call, err error) {
	_, _call, err = this.call(&ctx, this.options.ServiceTag, servicePath, serviceMethod, args, reply, done)
	return
}

// 跨服 同步调用 远程服务
func (this *Service) AcrossClusterCall(ctx context.Context, clusterTag string, servicePath string, serviceMethod string, args interface{}, reply interface{}) (err error) {
	var (
		done *client.Call
		conn net.Conn
	)
	seq := new(uint64)
	ctx = context.WithValue(ctx, seqKey{}, seq)
	if conn, done, err = this.call(&ctx, clusterTag, servicePath, serviceMethod, args, reply, make(chan *client.Call, 1)); err != nil {
		return
	}
	select {
	case <-ctx.Done(): // cancel by context
		this.mutex.Lock()
		call := this.pending[*seq]
		delete(this.pending, *seq)
		this.mutex.Unlock()
		if call != nil {
			call.Error = ctx.Err()
			call.Done <- call
		}
		return ctx.Err()
	case call := <-done.Done:
		err = call.Error
		meta := ctx.Value(share.ResMetaDataKey)
		if meta != nil && len(call.ResMetadata) > 0 {
			resMeta := meta.(map[string]string)
			for k, v := range call.ResMetadata {
				resMeta[k] = v
			}
			resMeta[share.ServerAddress] = conn.RemoteAddr().String()
		}
	}
	return
}

// 跨集群 广播
func (this *Service) AcrossClusterBroadcast(ctx context.Context, clusterTag string, servicePath string, serviceMethod string, args interface{}, reply interface{}) (err error) {
	err = this.broadcast(ctx, clusterTag, servicePath, serviceMethod, args, reply)
	return
}

// 跨服 异步调用 远程服务
func (this *Service) AcrossClusterGo(ctx context.Context, clusterTag, servicePath string, serviceMethod string, args interface{}, reply interface{}, done chan *client.Call) (_call *client.Call, err error) {
	_, _call, err = this.call(&ctx, clusterTag, servicePath, serviceMethod, args, reply, done)
	return
}

// 全集群广播
func (this *Service) ClusterBroadcast(ctx context.Context, servicePath string, serviceMethod string, args interface{}, reply interface{}) (err error) {
	err = this.clusterbroadcast(ctx, servicePath, serviceMethod, args, reply)
	return
}

// 监控rpc连接收到的请求消息 处理消息回调请求
func (this *Service) PostReadRequest(ctx context.Context, r *protocol.Message, e error) error {
	if isCallMessage := (r.MessageType() == protocol.Request); isCallMessage {
		var (
			stag     string
			selector ISelector
			ok       bool
		)
		req_metadata := r.Metadata
		this.options.Log.Debugf("PreReadRequest  ServicePath:%s ServicePath:%s Metadata:%v ", r.ServicePath, r.ServiceMethod, r.Metadata)
		if stag, ok = req_metadata[ServiceClusterTag]; ok {
			this.selectormutex.RLock()
			selector, ok = this.selectors[stag]
			this.selectormutex.RUnlock()
			if !ok {
				selector = newSelector(this.options.Log, stag, nil)
				this.selectormutex.Lock()
				this.selectors[stag] = selector
				this.selectormutex.Unlock()
			}

			if addr, ok := req_metadata[ServiceAddrKey]; ok {
				if _, ok = this.clientmeta[addr]; !ok {
					if smeta, ok := req_metadata[ServiceMetaKey]; ok {
						servers := make(map[string]string)
						this.clientmutex.Lock()
						this.clientmeta[addr] = smeta
						this.clients[addr] = ctx.Value(server.RemoteConnContextKey).(net.Conn)
						for k, v := range this.clientmeta {
							servers[k] = v
						}
						this.clientmutex.Unlock()
						selector.UpdateServer(servers)
						this.options.Log.Debugf("fond new node addr:%s smeta:%s \n", addr, smeta)
					}
				}
			}
		}
		return nil
	}
	e = errors.New("is callMessage")
	seq := r.Seq()
	this.mutex.Lock()
	call := this.pending[seq]
	delete(this.pending, seq)
	this.mutex.Unlock()
	switch {
	case call == nil:
		this.options.Log.Errorf("callmessage no found call:%d", seq)
	case r.MessageStatusType() == protocol.Error:
		if len(r.Metadata) > 0 {
			call.ResMetadata = r.Metadata
			call.Error = errors.New(r.Metadata[protocol.ServiceError])
		}
		if len(r.Payload) > 0 {
			data := r.Payload
			codec := share.Codecs[r.SerializeType()]
			if codec != nil {
				_ = codec.Decode(data, call.Reply)
			}
		}
		call.Done <- call
	default:
		data := r.Payload
		if len(data) > 0 {
			codec := share.Codecs[r.SerializeType()]
			if codec == nil {
				call.Error = errors.New(client.ErrUnsupportedCodec.Error())
			} else {
				err := codec.Decode(data, call.Reply)
				if err != nil {
					call.Error = err
				}
			}
		}
		if len(r.Metadata) > 0 {
			call.ResMetadata = r.Metadata
		}
		call.Done <- call
	}
	return nil
}

// 客户端配置 AutoConnect 的默认连接函数
func (this *Service) RpcxShakeHands(ctx context.Context, args *ServiceNode, reply *ServiceNode) error {
	// this.Debugf("RpcxShakeHands:%+v", ctx.Value(share.ReqMetaDataKey).(map[string]string))
	return nil
}

// 执行远程调用
func (this *Service) call(ctx *context.Context, clusterTag string, servicePath string, serviceMethod string, args interface{}, reply interface{}, done chan *client.Call) (conn net.Conn, _call *client.Call, err error) {
	var (
		spath      []string
		clientaddr string
		metadata   map[string]string
		selector   client.Selector
		ok         bool
	)
	if clusterTag == "" || servicePath == "" {
		err = fmt.Errorf("clusterTag:%s servicePath:%s no cant null", clusterTag, servicePath)
		return
	}
	metadata = map[string]string{
		ServiceClusterTag: this.options.ServiceTag,
		CallRoutRulesKey:  servicePath,
		ServiceAddrKey:    "tcp@" + this.options.ServiceAddr,
		ServiceMetaKey:    this.metadata,
	}
	spath = strings.Split(servicePath, "/")
	*ctx = context.WithValue(*ctx, share.ReqMetaDataKey, map[string]string{
		CallRoutRulesKey: servicePath,
		ServiceAddrKey:   "tcp@" + this.options.ServiceAddr,
		ServiceMetaKey:   this.metadata,
	})
	this.selectormutex.RLock()
	selector, ok = this.selectors[clusterTag]
	this.selectormutex.RUnlock()
	if !ok {
		err = fmt.Errorf("on found serviceTag:%s", clusterTag)
		return
	}

	if clientaddr = selector.Select(*ctx, spath[0], serviceMethod, args); clientaddr == "" {
		err = fmt.Errorf("on found servicePath:%s", servicePath)
		return
	}
	this.clientmutex.RLock()
	if conn, ok = this.clients[clientaddr]; !ok {
		err = fmt.Errorf("on found clientaddr:%s", clientaddr)
		this.clientmutex.RUnlock()
		return
	}
	this.clientmutex.RUnlock()
	_call = new(client.Call)
	_call.ServicePath = servicePath
	_call.ServiceMethod = serviceMethod
	_call.Args = args
	_call.Reply = reply
	if done == nil {
		done = make(chan *client.Call, 10) // buffered.
	} else {
		if cap(done) == 0 {
			log.Panic("rpc: done channel is unbuffered")
		}
	}
	_call.Done = done
	this.send(*ctx, conn, spath[0], serviceMethod, metadata, _call)
	return
}

// 执行远程调用
func (this *Service) broadcast(ctx context.Context, clusterTag string, servicePath string, serviceMethod string, args interface{}, reply interface{}) (err error) {
	var (
		spath       []string
		clientaddrs []string
		metadata    map[string]string
		selector    ISelector
		ok          bool
	)
	if servicePath == "" {
		err = errors.New("servicePath no cant null")
		return
	}
	metadata = map[string]string{
		ServiceClusterTag: this.options.ServiceTag,
		CallRoutRulesKey:  servicePath,
		ServiceAddrKey:    "tcp@" + this.options.ServiceAddr,
		ServiceMetaKey:    this.metadata,
	}
	spath = strings.Split(servicePath, "/")
	ctx = context.WithValue(ctx, share.ReqMetaDataKey, map[string]string{
		CallRoutRulesKey: servicePath,
		ServiceAddrKey:   "tcp@" + this.options.ServiceAddr,
		ServiceMetaKey:   this.metadata,
	})
	this.selectormutex.RLock()
	selector, ok = this.selectors[clusterTag]
	this.selectormutex.RUnlock()
	if !ok {
		err = fmt.Errorf("on found serviceTag:%s", clusterTag)
		return
	}
	if clientaddrs = selector.Find(ctx, spath[0], serviceMethod, args); clientaddrs == nil || len(clientaddrs) == 0 {
		err = fmt.Errorf("on found servicePath:%s", servicePath)
		return
	}
	l := len(clientaddrs)
	done := make(chan error, l)
	for _, v := range clientaddrs {
		go func(addr string) {
			this.clientmutex.RLock()
			conn, ok := this.clients[addr]
			if !ok {
				done <- fmt.Errorf("on found clientaddr:%s", addr)
				this.clientmutex.RUnlock()
				return
			}
			this.clientmutex.RUnlock()
			_call := new(client.Call)
			_call.ServicePath = servicePath
			_call.ServiceMethod = serviceMethod
			_call.Args = args
			_call.Reply = reply
			_call.Done = make(chan *client.Call, 10)
			seq := new(uint64)
			ctx = context.WithValue(ctx, seqKey{}, seq)
			this.send(ctx, conn, spath[0], serviceMethod, metadata, _call)
			select {
			case <-ctx.Done(): // cancel by context
				this.mutex.Lock()
				call := this.pending[*seq]
				delete(this.pending, *seq)
				this.mutex.Unlock()
				if call != nil {
					call.Error = ctx.Err()
					call.Done <- call
				}
				done <- ctx.Err()
			case call := <-_call.Done:
				err = call.Error
				meta := ctx.Value(share.ResMetaDataKey)
				if meta != nil && len(call.ResMetadata) > 0 {
					resMeta := meta.(map[string]string)
					for k, v := range call.ResMetadata {
						resMeta[k] = v
					}
					resMeta[share.ServerAddress] = conn.RemoteAddr().String()
				}
				done <- nil
			}
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
	return err
}

// 全集群广播
func (this *Service) clusterbroadcast(ctx context.Context, servicePath string, serviceMethod string, args interface{}, reply interface{}) (err error) {
	var (
		spath       []string
		clientaddrs []string
		metadata    map[string]string
	)
	if servicePath == "" {
		err = errors.New("servicePath no cant null")
		return
	}
	metadata = map[string]string{
		ServiceClusterTag: this.options.ServiceTag,
		CallRoutRulesKey:  servicePath,
		ServiceAddrKey:    "tcp@" + this.options.ServiceAddr,
		ServiceMetaKey:    this.metadata,
	}
	spath = strings.Split(servicePath, "/")
	ctx = context.WithValue(ctx, share.ReqMetaDataKey, map[string]string{
		CallRoutRulesKey: servicePath,
		ServiceAddrKey:   "tcp@" + this.options.ServiceAddr,
		ServiceMetaKey:   this.metadata,
	})
	addrs := make(map[string]struct{})
	this.selectormutex.RLock()
	for _, v := range this.selectors {
		if clientaddrs = v.Find(ctx, spath[0], serviceMethod, args); clientaddrs != nil && len(clientaddrs) >= 0 {
			for _, v1 := range clientaddrs {
				addrs[v1] = struct{}{}
			}
		}
	}
	this.selectormutex.RUnlock()
	if len(addrs) == 0 {
		err = fmt.Errorf("on found service:%s", spath[0])
		return
	}

	l := len(addrs)
	if l > 0 {
		done := make(chan error, l)
		for v, _ := range addrs {
			seq := new(uint64)
			tempctx := context.WithValue(ctx, seqKey{}, seq)
			go func(_ctx context.Context, addr string) {
				this.clientmutex.RLock()
				conn, ok := this.clients[addr]
				if !ok {
					done <- fmt.Errorf("on found clientaddr:%s", addr)
					this.clientmutex.RUnlock()
					return
				}
				this.clientmutex.RUnlock()
				_call := new(client.Call)
				_call.ServicePath = servicePath
				_call.ServiceMethod = serviceMethod
				_call.Args = args
				_call.Reply = reply
				_call.Done = make(chan *client.Call, 10)

				this.send(_ctx, conn, spath[0], serviceMethod, metadata, _call)
				seq, _ := _ctx.Value(seqKey{}).(*uint64)
				select {
				case <-_ctx.Done(): // cancel by context
					this.mutex.Lock()
					call := this.pending[*seq]
					delete(this.pending, *seq)
					this.mutex.Unlock()
					if call != nil {
						call.Error = _ctx.Err()
						call.Done <- call
					}
					done <- _ctx.Err()
				case call := <-_call.Done:
					err = call.Error
					meta := _ctx.Value(share.ResMetaDataKey)
					if meta != nil && len(call.ResMetadata) > 0 {
						resMeta := meta.(map[string]string)
						for k, v := range call.ResMetadata {
							resMeta[k] = v
						}
						resMeta[share.ServerAddress] = conn.RemoteAddr().String()
					}
					done <- nil
				}
			}(tempctx, v)
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

// 发送远程调用请求
func (this *Service) send(ctx context.Context, conn net.Conn, servicePath string, serviceMethod string, metadata map[string]string, call *client.Call) {
	defer func() {
		if call.Error != nil {
			call.Done <- call
		}
	}()
	serializeType := this.options.SerializeType
	codec := share.Codecs[serializeType]
	if codec == nil {
		call.Error = client.ErrUnsupportedCodec
		return
	}
	data, err := codec.Encode(call.Args)
	if err != nil {
		call.Error = err
		return
	}

	this.mutex.Lock()
	seq := this.seq
	this.seq++
	this.pending[seq] = call
	this.mutex.Unlock()
	if cseq, ok := ctx.Value(seqKey{}).(*uint64); ok {
		*cseq = seq
	}
	req := protocol.GetPooledMsg()
	req.SetMessageType(protocol.Request)
	req.SetSeq(seq)
	req.SetOneway(true)
	req.SetSerializeType(this.options.SerializeType)
	req.ServicePath = servicePath
	req.ServiceMethod = serviceMethod
	req.Metadata = metadata
	req.Payload = data

	b := req.EncodeSlicePointer()
	if _, err = conn.Write(*b); err != nil {
		call.Error = err
		this.mutex.Lock()
		delete(this.pending, seq)
		this.mutex.Unlock()
		return
	}
	protocol.PutData(b)
	protocol.FreeMsg(req)
	return
}
