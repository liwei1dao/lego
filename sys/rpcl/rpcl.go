package rpcl

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"runtime"
	"strings"
	"sync"

	"github.com/liwei1dao/lego/core"
	"github.com/liwei1dao/lego/sys/rpcl/connpool"
	gcore "github.com/liwei1dao/lego/sys/rpcl/core"
	"github.com/liwei1dao/lego/sys/rpcl/protocol"
)

var TypeOfError = reflect.TypeOf((*error)(nil)).Elem()
var TypeOfContext = reflect.TypeOf((*context.Context)(nil)).Elem()

func newSys(options Options) (sys *RPCL, err error) {
	sys = &RPCL{
		options:    options,
		serviceMap: make(map[string]*Server),
	}
	if sys.cpool, err = connpool.NewConnPool(sys, options.Config); err != nil {
		return
	}
	return
}

type RPCL struct {
	options      Options
	node         *core.ServiceNode
	selector     gcore.ISelector
	cpool        gcore.IConnPool
	serviceMapMu sync.RWMutex
	serviceMap   map[string]*Server
	pendingmutex sync.Mutex
	seq          uint64
	pending      map[uint64]*MessageCall
}

//启动系统
func (this *RPCL) Start() (err error) {
	return
}

//关闭系统
func (this *RPCL) Close() (err error) {
	return
}

func (this *RPCL) GetNodePath() string {
	return fmt.Sprintf("%s/%s/%s", this.options.ServiceNode.Tag, this.options.ServiceNode.Type, this.options.ServiceNode.Id)
}

func (this *RPCL) GetServiceNode() *core.ServiceNode {
	return this.node
}

//注册服务 批量注册
func (this *RPCL) Register(rcvr interface{}) (err error) {
	err = this.register(rcvr)
	return
}

//注册服务
func (this *RPCL) RegisterFunction(fn interface{}) (err error) {
	err = this.registerFunction(fn, "", false)
	return
}

//注册服务
func (this *RPCL) RegisterFunctionName(name string, fn interface{}) (err error) {
	err = this.registerFunction(fn, name, true)
	return
}

//注销服务
func (this *RPCL) UnRegister(name string) {
	this.serviceMapMu.Lock()
	delete(this.serviceMap, name)
	this.serviceMapMu.Unlock()
}

//同步执行
func (this *RPCL) Call(ctx context.Context, servicePath string, serviceMethod string, args interface{}, reply interface{}) (err error) { //同步调用 等待结果
	
	return nil
}

//异步执行 异步返回
func (this *RPCL) Go(ctx context.Context, servicePath string, serviceMethod string, args interface{}, reply interface{}) (call *MessageCall, err error) { //异步调用 异步返回
	return nil, nil
}

//同步执行 无返回
func (this *RPCL) GoNR(ctx context.Context, servicePath string, serviceMethod string, args interface{}) (err error) { //异步调用 无返回

	return nil
}

func (this *RPCL) Broadcast(ctx context.Context, servicePath string, serviceMethod string, args interface{}) (err error) {
	return
}

//接收到远程消息
func (this *RPCL) Handle(client gcore.IConnClient, message *protocol.Message) {
	if message.MessageType() == protocol.Request { //请求消息
		if res, _ := this.handleRequest(context.Background(), message); res != nil {
			if !message.IsOneway() { //需要回应
				if len(res.Payload) > 1024 && res.CompressType() != protocol.CompressNone {
					res.SetCompressType(res.CompressType())
				}
				data := res.EncodeSlicePointer()
				client.Write(*data)
				protocol.PutData(data)
			}
		}
	} else { //回应
		this.handleresponse(context.Background(), message)
	}
}

//反射批量注册 服务---------------------------------------------------------------------------------
func (this *RPCL) register(rcvr interface{}) (err error) {
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
func (this *RPCL) registerFunction(fn interface{}, name string, useName bool) (err error) {
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
		err = fmt.Errorf("rpcl.registerFunction: no func name for type:%s", f.Type().String())
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

	//注册类型池
	reflectTypePools.Init(argType)
	reflectTypePools.Init(replyType)
	return
}

//执行远程服务---------------------------------------------------------------------------------------
func (this *RPCL) call(ctx context.Context, seq uint64, call *MessageCall) (err error) {
	var (
		targets []*core.ServiceNode
		client  gcore.IConnClient
		data    []byte
		allData *[]byte
	)
	targets = this.selector.Select(ctx, call.ServicePath)
	if targets == nil || len(targets) == 0 {
		err = ErrXClientNoServer
		return
	}
	req := protocol.GetPooledMsg()
	req.SetMessageType(protocol.Request)
	req.SetSeq(seq)
	if call.Reply == nil {
		req.SetOneway(true)
	}
	if call.Metadata != nil {
		req.Metadata = call.Metadata
	}

	req.ServicePath = call.ServicePath
	req.ServiceMethod = call.ServiceMethod

	codec := codecs[this.options.SerializeType]
	if codec == nil {
		err = ErrUnsupportedCodec
		return
	}
	data, err = codec.Marshal(call.Args)
	if err != nil {
		return
	}
	if len(data) > 1024 && this.options.CompressType != protocol.CompressNone {
		req.SetCompressType(this.options.CompressType)
	}
	req.Payload = data
	allData = req.EncodeSlicePointer()
	defer func() {
		protocol.PutData(allData)
		protocol.FreeMsg(req)
	}()
	if client, err = this.cpool.GetClient(targets[0]); err != nil {
		return
	}
	err = client.Write(*allData)
	return
}

//处理响应------------------------------------------------------------------------------------------
func (this *RPCL) handleresponse(ctx context.Context, res *protocol.Message) {
	var call *MessageCall
	seq := res.Seq()
	isServerMessage := (res.MessageType() == protocol.Request && !res.IsHeartbeat() && res.IsOneway())
	if !isServerMessage {
		this.pendingmutex.Lock()
		call = this.pending[seq]
		delete(this.pending, seq)
		this.pendingmutex.Unlock()
	}
	switch {
	case call == nil:
		this.Warnf("call is nil res:%v", res)
	case res.MessageStatusType() == protocol.Error:
		if len(res.Metadata) > 0 {
			call.ResMetadata = res.Metadata
			call.Error = errors.New(res.Metadata[ServiceError])
		}
		if len(res.Payload) > 0 {
			data := res.Payload
			codec := codecs[res.SerializeType()]
			if codec != nil {
				_ = codec.Unmarshal(data, call.Reply)
			}
			call.done(this)
		}
	default:
		data := res.Payload
		if len(data) > 0 {
			codec := codecs[res.SerializeType()]
			if codec == nil {
				call.Error = ErrUnsupportedCodec
			} else {
				err := codec.Unmarshal(data, call.Reply)
				if err != nil {
					call.Error = err
				}
			}
		}
		if len(res.Metadata) > 0 {
			call.ResMetadata = res.Metadata
		}
		call.done(this)
	}
}

//处理服务消息---------------------------------------------------------------------------------------
func (this *RPCL) handleRequest(ctx context.Context, req *protocol.Message) (res *protocol.Message, err error) {
	methodName := req.ServiceMethod
	res = req.Clone()
	res.SetMessageType(protocol.Response)
	this.serviceMapMu.RLock()
	service, ok := this.serviceMap[methodName]
	this.Debugf("server get service %+v for an request %+v", service, req)
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

	err = codec.Unmarshal(req.Payload, argv)
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
			res.Payload = data
		}
		return handleError(res, err)
	}
	if !req.IsOneway() {
		data, err := codec.Marshal(replyv)
		reflectTypePools.Put(service.ReplyType, replyv)
		if err != nil {
			return handleError(res, err)
		}
		res.Payload = data
	} else if replyv != nil {
		reflectTypePools.Put(service.ReplyType, replyv)
	}
	this.Debugf("server called service %+v for an request %+v", service, req)
	return
}

func handleError(res *protocol.Message, err error) (*protocol.Message, error) {
	res.SetMessageStatusType(protocol.Error)
	if res.Metadata == nil {
		res.Metadata = make(map[string]string)
	}
	res.Metadata[ServiceError] = err.Error()
	return res, err
}
