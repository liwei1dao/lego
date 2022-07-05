package rpcl

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"runtime"
	"strings"
	"sync"
	"time"
	"unicode"
	"unicode/utf8"

	jsoniter "github.com/json-iterator/go"
	"github.com/liwei1dao/lego/sys/rpcl/connect"
	"github.com/liwei1dao/lego/sys/rpcl/core"
	"github.com/liwei1dao/lego/utils/codec"
)

func newSys(options Options) (sys *RPCL, err error) {
	sys = &RPCL{
		options:    options,
		decoder:    options.Decoder,
		encoder:    options.Encoder,
		serviceMap: make(map[string]*core.Server),
	}
	if sys.decoder == nil {
		sys.decoder = &codec.Decoder{DefDecoder: jsoniter.Unmarshal}
	}
	if sys.encoder == nil {
		sys.encoder = &codec.Encoder{DefEncoder: jsoniter.Marshal}
	}
	if sys.connect, err = connect.NewConnect(sys); err != nil {
		return
	}
	return
}

type RPCL struct {
	options      Options
	node         *core.ServiceNode
	decoder      codec.IDecoder
	encoder      codec.IEncoder
	connect      core.IConnect
	serviceMapMu sync.RWMutex
	serviceMap   map[string]*core.Server
}

//启动系统
func (this *RPCL) Start() (err error) {
	return
}

//关闭系统
func (this *RPCL) Close() (err error) {
	return
}
func (this *RPCL) GetUpdateInterval() time.Duration {
	return this.options.UpdateInterval
}

func (this *RPCL) GetBasePath() string {
	return this.options.BasePath
}

func (this *RPCL) GetNodePath() string {
	return fmt.Sprintf("%s/%s/%s", this.options.BasePath, this.options.ServiceType, this.options.ServiceId)
}

func (this *RPCL) GetServiceNode() *core.ServiceNode {
	return this.node
}

func (this *RPCL) Decoder() codec.IDecoder {
	return this.decoder
}
func (this *RPCL) Encoder() codec.IEncoder {
	return this.encoder
}

func (this *RPCL) GetKafkaAddr() []string {
	return this.options.Kafka_Addr
}
func (this *RPCL) GetKafkaVersion() string {
	return this.options.Kafka_Version
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
func (this *RPCL) Call(ctx context.Context, route core.IRoute, serviceMethod string, args interface{}, reply interface{}) (err error) { //同步调用 等待结果

	return nil
}

//异步执行 异步返回
func (this *RPCL) Go(ctx context.Context, route core.IRoute, serviceMethod string, args interface{}, reply interface{}, done chan *core.Call) (call *core.Call, err error) { //异步调用 异步返回
	// sid := this.selector.Select(ctx, route, serviceMethod)
	return nil, nil
}

//同步执行 无返回
func (this *RPCL) GoNR(ctx context.Context, route core.IRoute, serviceMethod string, args interface{}) (err error) { //异步调用 无返回
	// sid := this.selector.Select(ctx, route, serviceMethod)
	return nil
}

//接收到远程消息
func (this *RPCL) Receive(message *core.Message) {

}

///日志***********************************************************************
func (this *RPCL) Debug() bool {
	return this.options.Debug
}
func (this *RPCL) Debugf(format string, a ...interface{}) {
	if this.options.Debug {
		this.options.Log.Debugf("[SYS RPCL] "+format, a)
	}
}
func (this *RPCL) Infof(format string, a ...interface{}) {
	if this.options.Debug {
		this.options.Log.Infof("[SYS RPCL] "+format, a)
	}
}
func (this *RPCL) Warnf(format string, a ...interface{}) {
	if this.options.Debug {
		this.options.Log.Warnf("[SYS RPCL] "+format, a)
	}
}
func (this *RPCL) Errorf(format string, a ...interface{}) {
	if this.options.Debug {
		this.options.Log.Errorf("[SYS RPCL] "+format, a)
	}
}
func (this *RPCL) Panicf(format string, a ...interface{}) {
	if this.options.Debug {
		this.options.Log.Panicf("[SYS RPCL] "+format, a)
	}
}
func (this *RPCL) Fatalf(format string, a ...interface{}) {
	if this.options.Debug {
		this.options.Log.Fatalf("[SYS RPCL] "+format, a)
	}
}

//---------------------------------------------------------------------------------
//反射批量注册 服务
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
		if !ctxType.Implements(core.TypeOfContext) {
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
		if returnType := mtype.Out(0); returnType != core.TypeOfError {
			continue
		}
		this.serviceMapMu.Lock()
		this.serviceMap[mname] = &core.Server{Fn: vof.MethodByName(mname), ArgType: argType, ReplyType: replyType}
		this.serviceMapMu.Unlock()
		//注册类型池
		reflectTypePools.Init(argType)
		reflectTypePools.Init(replyType)
	}
	return
}

//注册服务
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
	if !ctxType.Implements(core.TypeOfContext) {
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

	if returnType := t.Out(0); returnType != core.TypeOfError {
		return fmt.Errorf("function %s returns %s, not error", f.Type().String(), returnType.String())
	}

	this.serviceMapMu.Lock()
	this.serviceMap[fname] = &core.Server{Fn: f, ArgType: argType, ReplyType: replyType}
	this.serviceMapMu.Unlock()

	//注册类型池
	reflectTypePools.Init(argType)
	reflectTypePools.Init(replyType)
	return
}

//处理请求
func (this *RPCL) handleRequest(ctx context.Context, req *core.Message) (res *core.Message, err error) {
	methodName := req.ServiceMethod
	res = req.Clone()
	res.SetMessageType(core.Response)
	this.serviceMapMu.RLock()
	service, ok := this.serviceMap[methodName]
	this.Debugf("server get service %+v for an request %+v", service, req)
	this.serviceMapMu.RUnlock()
	if ok {
		argv := reflectTypePools.Get(service.ArgType)
		err = this.decoder.Decoder(req.Payload, argv)
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
				data, err := this.encoder.Encoder(replyv)
				reflectTypePools.Put(service.ReplyType, replyv)
				if err != nil {
					return handleError(res, err)
				}
				res.Payload = data
			}
			return handleError(res, err)
		}
		if !req.IsOneway() {
			data, err := this.encoder.Encoder(replyv)
			reflectTypePools.Put(service.ReplyType, replyv)
			if err != nil {
				return handleError(res, err)
			}
			res.Payload = data
		} else if replyv != nil {
			reflectTypePools.Put(service.ReplyType, replyv)
		}
		this.Debugf("server called service %+v for an request %+v", service, req)
	}
	return
}

func isExportedOrBuiltinType(t reflect.Type) bool {
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return isExported(t.Name()) || t.PkgPath() == ""
}

func isExported(name string) bool {
	rune, _ := utf8.DecodeRuneInString(name)
	return unicode.IsUpper(rune)
}

func handleError(res *core.Message, err error) (*core.Message, error) {
	res.SetMessageStatusType(core.Error)
	if res.Metadata == nil {
		res.Metadata = make(map[string]string)
	}
	res.Metadata[core.ServiceError] = err.Error()
	return res, err
}
