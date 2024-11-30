package rpc

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"runtime"
	"strings"
	"sync"

	"github.com/liwei1dao/lego/core"
	"github.com/liwei1dao/lego/sys/discovery"
	"github.com/liwei1dao/lego/sys/log"
	"github.com/liwei1dao/lego/sys/rpc/protocol"
	"github.com/liwei1dao/lego/sys/rpc/rpccore"
)

func newService(cpool rpccore.IConnPool, options *Options) (sys *Service, err error) {
	sys = &Service{
		options:    options,
		cpool:      cpool,
		serviceMap: make(map[string]*Server),
	}
	return
}

type Service struct {
	options      *Options
	cpool        rpccore.IConnPool
	heartbeat    []byte
	serviceMapMu sync.RWMutex
	serviceMap   map[string]*Server
}

func (this *Service) Heartbeat() []byte {
	return this.heartbeat
}

// 启动系统
func (this *Service) Start() (err error) {
	if err = discovery.Start(); err != nil {
		return
	}
	return
}

// 关闭系统
func (this *Service) Close() (err error) {
	if err = discovery.Stop(); err != nil {
		return
	}
	return
}

func (this *Service) ServiceNode() core.IServiceNode {
	return this.options.ServiceNode
}

// 注册服务 批量注册
func (this *Service) Register(rcvr interface{}) (err error) {
	err = this.register(rcvr)
	return
}

// 注册服务
func (this *Service) RegisterFunction(fn interface{}) (err error) {
	err = this.registerFunction(fn, "", false)
	return
}

// 注册服务
func (this *Service) RegisterFunctionName(name string, fn interface{}) (err error) {
	err = this.registerFunction(fn, name, true)
	return
}

// 注销服务
func (this *Service) UnRegister(name string) {
	this.serviceMapMu.Lock()
	delete(this.serviceMap, name)
	this.serviceMapMu.Unlock()
}

// 接收到远程消息
func (this *Service) Handle(client rpccore.IConnClient, message rpccore.IMessage) {
	defer func() {
		if r := recover(); r != nil {
			buf := make([]byte, 1024)
			buf = buf[:runtime.Stack(buf, true)]
			this.options.Log.Errorf("failed to handle the request: %v， stacks: %s", r, buf)
		}
	}()
	ctx := rpccore.WithValue(context.Background(), rpccore.RemoteConnContextKey, client)
	// this.options.Log.Debug("[handle] message", log.Field{Key: "Header", Value: message.PrintHeader()}, log.Field{Key: "ServiceMethod", Value: message.ServiceMethod()})
	if message.MessageType() == rpccore.Request { //请求消息
		if message.IsHeartbeat() { //心跳
			client.ResetHbeat()
			return
		}
		if message.IsShakeHands() {
			this.ShakehandsResponse(ctx, client, message)
			return
		}
		cancelFunc := parseServerTimeout(ctx, message)
		if cancelFunc != nil {
			defer cancelFunc()
		}
		if res, _ := this.handleRequest(ctx, message); res != nil {
			if !message.IsOneway() { //需要回应
				if len(res.Payload()) > 1024 && res.CompressType() != rpccore.CompressNone {
					res.SetCompressType(res.CompressType())
				}
				data := res.EncodeSlicePointer()
				client.Write(*data)
				protocol.PutData(data)
			}
		}
	}
}

// 反射批量注册 服务---------------------------------------------------------------------------------
func (this *Service) register(rcvr interface{}) (err error) {
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
func (this *Service) registerFunction(fn interface{}, name string, useName bool) (err error) {
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
		err = fmt.Errorf("Service.registerFunction: no func name for type:%s", f.Type().String())
		return
	}
	t := f.Type()
	if t.NumIn() != 3 {
		return fmt.Errorf("Servicex.registerFunction: has wrong number of ins: %s", f.Type().String())
	}
	if t.NumOut() != 1 {
		return fmt.Errorf("Servicex.registerFunction: has wrong number of outs: %s", f.Type().String())
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

// 处理服务消息---------------------------------------------------------------------------------------
func (this *Service) handleRequest(ctx context.Context, req rpccore.IMessage) (res rpccore.IMessage, err error) {
	methodName := req.ServiceMethod()
	res = req.Clone()
	res.SetMessageType(rpccore.Response)
	this.serviceMapMu.RLock()
	service, ok := this.serviceMap[methodName]
	this.options.Log.Debug("handleRequest", log.Field{Key: "ServiceMethod", Value: req.ServiceMethod()}, log.Field{Key: "Form", Value: req.From()})
	this.serviceMapMu.RUnlock()
	if !ok {
		err = errors.New("Servicex: can't find service " + methodName)
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
	return
}
