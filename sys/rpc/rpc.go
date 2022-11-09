package rpc

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/liwei1dao/lego/core"
	"github.com/liwei1dao/lego/sys/discovery"
	"github.com/liwei1dao/lego/sys/event"
	"github.com/liwei1dao/lego/sys/log"
	"github.com/liwei1dao/lego/sys/rpc/connpool"
	"github.com/liwei1dao/lego/sys/rpc/protocol"
	"github.com/liwei1dao/lego/sys/rpc/rpccore"
	"github.com/liwei1dao/lego/sys/rpc/selector"
)

var TypeOfError = reflect.TypeOf((*error)(nil)).Elem()
var TypeOfContext = reflect.TypeOf((*context.Context)(nil)).Elem()

func newSys(options *Options) (sys *rpc, err error) {
	sys = &rpc{
		options:    options,
		serviceMap: make(map[string]*Server),
		seq:        0,
		pending:    make(map[uint64]*MessageCall),
	}

	if options.ConnectType == rpccore.Tcp {
		options.ServiceNode.Addr = options.MessageEndpoints[0]
	}

	if sys.cpool, err = connpool.NewConnPool(sys, options.Log, &rpccore.Config{
		ConnectType: options.ConnectType,
		Endpoints:   options.MessageEndpoints,
	}); err != nil {
		return
	}
	if sys.discovery, err = discovery.NewSys(
		discovery.SetBasePath(options.ServiceNode.Tag),
		discovery.SetServiceNode(options.ServiceNode),
		discovery.SetStoreType(options.DiscoveryStoreType),
		discovery.SetEndpoints(options.DiscoveryEndpoints),
		discovery.SetUpdateInterval(time.Duration(options.DiscoveryInterval)*time.Second),
		discovery.SetCodec(codecs[rpccore.JSON]),
		discovery.SetDebug(options.Debug),
	); err != nil {
		return
	}
	sys.selector, err = selector.NewSelector(sys.discovery.GetServices())
	return
}

type rpc struct {
	options      *Options
	cpool        rpccore.IConnPool
	discovery    discovery.ISys
	selector     rpccore.ISelector
	heartbeat    []byte
	serviceMapMu sync.RWMutex
	serviceMap   map[string]*Server
	pendingmutex sync.Mutex
	seq          uint64
	pending      map[uint64]*MessageCall
}

func (this *rpc) Heartbeat() []byte {
	return this.heartbeat
}

//启动系统
func (this *rpc) Start() (err error) {
	if err = this.cpool.Start(); err != nil {
		return
	}
	if err = this.discovery.Start(); err != nil {
		return
	}
	this.selector.UpdateServer(this.discovery.GetServices())
	go func() { //监控服务发现
		var add, del, change []*core.ServiceNode
		for v := range this.discovery.WatchService() {
			add, del, change = this.selector.UpdateServer(v)
			if len(add) > 0 { //发现节点
				this.options.Log.Debug("发现节点!", log.Field{Key: "nodes", Value: add})
				event.TriggerEvent(Event_RpcDiscoverNewNodes, add)
			}
			if len(del) > 0 { //丢失节点
				this.options.Log.Debug("丢失节点!", log.Field{Key: "nodes", Value: del})
				event.TriggerEvent(Event_RpcLoseNodes, del)
			}
			if len(change) > 0 { //节点属性变化
				this.options.Log.Debug("节点属性变化!", log.Field{Key: "nodes", Value: change})
				event.TriggerEvent(Event_RpChangeNodes, change)
			}
		}
	}()
	return
}

//关闭系统
func (this *rpc) Close() (err error) {
	if err = this.discovery.Close(); err != nil {
		return
	}
	if err = this.cpool.Close(); err != nil {
		return
	}
	return
}

func (this *rpc) ServiceNode() *core.ServiceNode {
	return this.options.ServiceNode
}

//注册服务 批量注册
func (this *rpc) Register(rcvr interface{}) (err error) {
	err = this.register(rcvr)
	return
}

//注册服务
func (this *rpc) RegisterFunction(fn interface{}) (err error) {
	err = this.registerFunction(fn, "", false)
	return
}

//注册服务
func (this *rpc) RegisterFunctionName(name string, fn interface{}) (err error) {
	err = this.registerFunction(fn, name, true)
	return
}

//注销服务
func (this *rpc) UnRegister(name string) {
	this.serviceMapMu.Lock()
	delete(this.serviceMap, name)
	this.serviceMapMu.Unlock()
}

//同步执行
func (this *rpc) Call(ctx context.Context, servicePath string, serviceMethod string, req interface{}, reply interface{}) (err error) { //同步调用 等待结果
	seq := new(uint64)
	ctx = rpccore.WithValue(ctx, rpccore.CallSeqKey, seq)
	this.options.Log.Debug("Call Start", log.Field{Key: "servicePath", Value: servicePath}, log.Field{Key: "serviceMethod", Value: serviceMethod}, log.Field{Key: "req", Value: req})
	defer func() {
		this.options.Log.Debug("Call End", log.Field{Key: "servicePath", Value: servicePath}, log.Field{Key: "serviceMethod", Value: serviceMethod}, log.Field{Key: "req", Value: req}, log.Field{Key: "reply", Value: reply})
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

//异步执行 异步返回
func (this *rpc) ClentForCall(ctx context.Context, client rpccore.IConnClient, serviceMethod string, req interface{}, reply interface{}) (err error) { //异步调用 异步返回
	seq := new(uint64)
	ctx = rpccore.WithValue(ctx, rpccore.CallSeqKey, seq)
	this.options.Log.Debug("ClentForCall start!", log.Field{Key: "node", Value: client.ServiceNode()}, log.Field{Key: "serviceMethod", Value: serviceMethod}, log.Field{Key: "req", Value: req})
	defer func() {
		this.options.Log.Debug("ClentForGo end!", log.Field{Key: "node", Value: client.ServiceNode()}, log.Field{Key: "serviceMethod", Value: serviceMethod}, log.Field{Key: "req", Value: req}, log.Field{Key: "reply", Value: reply})
	}()
	var call *MessageCall
	call, err = this.call2(ctx, client, serviceMethod, req, reply)
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

//异步执行 异步返回
func (this *rpc) Go(ctx context.Context, servicePath string, serviceMethod string, req interface{}, reply interface{}) (call *MessageCall, err error) { //异步调用 异步返回
	seq := new(uint64)
	ctx = rpccore.WithValue(ctx, rpccore.CallSeqKey, seq)
	this.options.Log.Debug("Go start!", log.Field{Key: "servicePath", Value: servicePath}, log.Field{Key: "serviceMethod", Value: serviceMethod}, log.Field{Key: "req", Value: req})
	defer func() {
		this.options.Log.Debug("Go end!", log.Field{Key: "servicePath", Value: servicePath}, log.Field{Key: "serviceMethod", Value: serviceMethod}, log.Field{Key: "req", Value: req}, log.Field{Key: "reply", Value: reply})
	}()
	call, err = this.call(ctx, servicePath, serviceMethod, req, reply)
	return
}

//异步执行 异步返回
func (this *rpc) ClentForGo(ctx context.Context, client rpccore.IConnClient, serviceMethod string, req interface{}, reply interface{}) (call *MessageCall, err error) { //异步调用 异步返回
	seq := new(uint64)
	ctx = rpccore.WithValue(ctx, rpccore.CallSeqKey, seq)
	this.options.Log.Debug("ClentForGo start!", log.Field{Key: "node", Value: client.ServiceNode()}, log.Field{Key: "serviceMethod", Value: serviceMethod}, log.Field{Key: "req", Value: req})
	defer func() {
		this.options.Log.Debug("ClentForGo end!", log.Field{Key: "node", Value: client.ServiceNode()}, log.Field{Key: "serviceMethod", Value: serviceMethod}, log.Field{Key: "req", Value: req}, log.Field{Key: "reply", Value: reply})
	}()
	call, err = this.call2(ctx, client, serviceMethod, req, reply)
	return
}

func (this *rpc) Broadcast(ctx context.Context, servicePath string, serviceMethod string, args interface{}) (err error) {
	return
}

//接收到远程消息
func (this *rpc) Handle(client rpccore.IConnClient, message rpccore.IMessage) {
	this.options.Log.Debug("[handle] message", log.Field{Key: "Header", Value: message.PrintHeader()}, log.Field{Key: "ServiceMethod", Value: message.ServiceMethod()})
	if message.MessageType() == rpccore.Request { //请求消息
		if message.IsHeartbeat() { //心跳
			client.ResetHbeat()
			return
		}
		if message.IsShakeHands() {
			this.ShakehandsResponse(rpccore.NewContext(context.Background()), client, message)
			return
		}
		if res, _ := this.handleRequest(rpccore.NewContext(context.Background()), message); res != nil {
			if !message.IsOneway() { //需要回应
				if len(res.Payload()) > 1024 && res.CompressType() != rpccore.CompressNone {
					res.SetCompressType(res.CompressType())
				}
				data := res.EncodeSlicePointer()
				client.Write(*data)
				protocol.PutData(data)
			}
		}
	} else { //回应
		this.handleresponse(rpccore.NewContext(context.Background()), message)
	}
}

//反射批量注册 服务---------------------------------------------------------------------------------
func (this *rpc) register(rcvr interface{}) (err error) {
	typ := reflect.TypeOf(rcvr)
	vof := reflect.ValueOf(rcvr)
	if vof.IsValid() {
		err = errors.New("rcvr IsValid!")
		return
	}
	if vof.Kind() != reflect.Ptr {
		err = errors.New("rcvr no Ptr!")
		return
	}
	for m := 0; m < typ.NumMethod(); m++ {
		method := typ.Method(m)
		mtype := method.Type
		mname := method.Name
		if method.PkgPath != "" {
			continue
		}
		if mtype.NumIn() != 4 {
			continue
		}
		ctxType := mtype.In(1)
		if !ctxType.Implements(TypeOfContext) {
			continue
		}

		argType := mtype.In(2)
		if !isExportedOrBuiltinType(argType) {
			continue
		}
		// Third arg must be a pointer.
		replyType := mtype.In(3)
		if replyType.Kind() != reflect.Ptr {
			continue
		}
		// Reply type must be exported.
		if !isExportedOrBuiltinType(replyType) {
			continue
		}
		// Method needs one out.
		if mtype.NumOut() != 1 {
			continue
		}
		// The return type of the method must be error.
		if returnType := mtype.Out(0); returnType != TypeOfError {
			continue
		}
		this.serviceMapMu.Lock()
		this.serviceMap[mname] = &Server{Fn: vof.MethodByName(mname), ArgType: argType, ReplyType: replyType}
		this.serviceMapMu.Unlock()
		//注册类型池
		reflectTypePools.Init(argType)
		reflectTypePools.Init(replyType)
	}
	return
}
func (this *rpc) registerFunction(fn interface{}, name string, useName bool) (err error) {
	f, ok := fn.(reflect.Value)
	if !ok {
		f = reflect.ValueOf(fn)
	}
	if f.Kind() != reflect.Func {
		err = errors.New("function must be func or bound method")
		return
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
		err = fmt.Errorf("rpc.registerFunction: no func name for type:%s", f.Type().String())
		return
	}
	t := f.Type()
	if t.NumIn() != 3 {
		return fmt.Errorf("rpcx.registerFunction: has wrong number of ins: %s", f.Type().String())
	}
	if t.NumOut() != 1 {
		return fmt.Errorf("rpcx.registerFunction: has wrong number of outs: %s", f.Type().String())
	}

	ctxType := t.In(0)
	if !ctxType.Implements(TypeOfContext) {
		return fmt.Errorf("function %s must use context as  the first parameter", f.Type().String())
	}

	argType := t.In(1)
	if !isExportedOrBuiltinType(argType) {
		return fmt.Errorf("function %s parameter type not exported: %v", f.Type().String(), argType)
	}

	replyType := t.In(2)
	if replyType.Kind() != reflect.Ptr {
		return fmt.Errorf("function %s reply type not a pointer: %s", f.Type().String(), replyType)
	}
	if !isExportedOrBuiltinType(replyType) {
		return fmt.Errorf("function %s reply type not exported: %v", f.Type().String(), replyType)
	}

	if returnType := t.Out(0); returnType != TypeOfError {
		return fmt.Errorf("function %s returns %s, not error", f.Type().String(), returnType.String())
	}

	this.serviceMapMu.Lock()
	this.serviceMap[fname] = &Server{Fn: f, ArgType: argType, ReplyType: replyType}
	this.serviceMapMu.Unlock()
	this.options.Log.Debug("注册服务!", log.Field{Key: "func", Value: fname})
	//注册类型池
	reflectTypePools.Init(argType)
	reflectTypePools.Init(replyType)
	return
}

//执行远程服务---------------------------------------------------------------------------------------
func (this *rpc) call(ctx context.Context, servicePath string, serviceMethod string, args interface{}, reply interface{}) (call *MessageCall, err error) {
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

func (this *rpc) call2(ctx context.Context, client rpccore.IConnClient, serviceMethod string, args interface{}, reply interface{}) (call *MessageCall, err error) {
	call = new(MessageCall)
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
	err = this.send(client, call, seq)
	return
}

//获取请求消息对象
func (this *rpc) getMessage(serviceMethod string, args interface{}, reply interface{}) (call *MessageCall, req *protocol.Message, err error) {
	var data []byte
	call = new(MessageCall)
	call.ServiceMethod = serviceMethod
	call.Done = make(chan *MessageCall, 10)
	call.Args = this.ServiceNode() //自己发起握手 需要传递本服务的节点信息
	req = protocol.GetPooledMsg()
	req.SetVersion(this.options.ProtoVersion)
	if reply != nil {
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
	data, err = codecs[rpccore.ProtoBuffer].Marshal(call.Args)
	if err != nil {
		return
	}
	if len(data) > 1024 && this.options.CompressType != rpccore.CompressNone {
		req.SetCompressType(this.options.CompressType)
	}
	req.SetPayload(data)
	return
}

func (this *rpc) getclient(ctx context.Context, servicePath string) (client rpccore.IConnClient, err error) {
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

func (this *rpc) send(client rpccore.IConnClient, call *MessageCall, seq uint64) (err error) {
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

//处理响应------------------------------------------------------------------------------------------
func (this *rpc) handleresponse(ctx context.Context, res rpccore.IMessage) {
	var call *MessageCall
	seq := res.Seq()
	isServerMessage := (res.MessageType() == rpccore.Request && !res.IsHeartbeat() && res.IsOneway())
	if !isServerMessage {
		this.pendingmutex.Lock()
		call = this.pending[seq]
		delete(this.pending, seq)
		this.pendingmutex.Unlock()
	}
	switch {
	case call == nil:
		this.options.Log.Warnf("call is nil res:%v", res)
	case res.MessageStatusType() == rpccore.Error:
		if len(res.Metadata()) > 0 {
			call.ResMetadata = res.Metadata()
			call.Error = errors.New(res.Metadata()[rpccore.ServiceError])
		}
		if len(res.Payload()) > 0 {
			data := res.Payload()
			codec := codecs[res.SerializeType()]
			if codec != nil {
				_ = codec.Unmarshal(data, call.Reply)
			}
			call.done(this.options.Log)
		}
	default:
		data := res.Payload()
		if len(data) > 0 {
			codec := codecs[res.SerializeType()]
			if codec == nil {
				call.Error = rpccore.ErrUnsupportedCodec
			} else {
				err := codec.Unmarshal(data, call.Reply)
				if err != nil {
					call.Error = err
				}
			}
		}
		if len(res.Metadata()) > 0 {
			call.ResMetadata = res.Metadata()
		}
		call.done(this.options.Log)
	}
}

//处理服务消息---------------------------------------------------------------------------------------
func (this *rpc) handleRequest(ctx context.Context, req rpccore.IMessage) (res rpccore.IMessage, err error) {
	methodName := req.ServiceMethod()
	res = req.Clone()
	res.SetMessageType(rpccore.Response)
	this.serviceMapMu.RLock()
	service, ok := this.serviceMap[methodName]
	this.options.Log.Debug("handleRequest", log.Field{Key: "ServiceMethod", Value: req.ServiceMethod()}, log.Field{Key: "Form", Value: req.From()})
	this.serviceMapMu.RUnlock()
	if !ok {
		err = errors.New("rpcx: can't find service " + methodName)
		return handleError(res, err)
	}
	codec := codecs[req.SerializeType()]
	if codec == nil {
		err = fmt.Errorf("can not find codec for %d", req.SerializeType())
		return handleError(res, err)
	}

	argv := reflectTypePools.Get(service.ArgType)

	err = codec.Unmarshal(req.Payload(), argv)
	if err != nil {
		return handleError(res, err)
	}
	replyv := reflectTypePools.Get(service.ReplyType)
	if service.ArgType.Kind() != reflect.Ptr {
		err = service.Call(ctx, reflect.ValueOf(argv).Elem(), reflect.ValueOf(replyv))
	} else {
		err = service.Call(ctx, reflect.ValueOf(argv), reflect.ValueOf(replyv))
	}
	reflectTypePools.Put(service.ArgType, argv)
	if err != nil {
		if replyv != nil {
			data, err := codec.Marshal(replyv)
			reflectTypePools.Put(service.ReplyType, replyv)
			if err != nil {
				return handleError(res, err)
			}
			res.SetPayload(data)
		}
		return handleError(res, err)
	}
	if !req.IsOneway() {
		data, err := codec.Marshal(replyv)
		reflectTypePools.Put(service.ReplyType, replyv)
		if err != nil {
			return handleError(res, err)
		}
		res.SetPayload(data)
	} else if replyv != nil {
		reflectTypePools.Put(service.ReplyType, replyv)
	}
	this.options.Log.Debugf("server called service %+v for an request %+v", service, req)
	return
}
