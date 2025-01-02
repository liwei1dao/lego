package lgrpc

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"sync"

	"github.com/liwei1dao/lego/sys/rpc"
	"github.com/liwei1dao/lego/utils"
)

func newServer() *Server {
	return &Server{
		plugins:    &pluginContainer{},
		serviceMap: make(map[string]*service),
	}
}

type Server struct {
	plugins      PluginContainer
	serviceMapMu sync.RWMutex
	serviceMap   map[string]*service
}

func (this *Server) Register(name string, fn interface{}, metadata string) error {
	err := this.registerFunction(name, fn)
	if err != nil {
		return err
	}
	return this.plugins.DoRegister(name, fn, metadata)
}

func (this *Server) UnregisterAll() error {
	this.serviceMapMu.RLock()
	defer this.serviceMapMu.RUnlock()
	var es []error
	for k := range this.serviceMap {
		err := this.plugins.DoUnregister(k)
		if err != nil {
			es = append(es, err)
		}
	}

	if len(es) > 0 {
		return utils.NewMultiError(es)
	}
	return nil
}

// 注册服务
func (this *Server) registerFunction(name string, fn interface{}) (err error) {
	f, ok := fn.(reflect.Value)
	if !ok {
		f = reflect.ValueOf(fn)
	}
	if f.Kind() != reflect.Func {
		err = errors.New("function must be func or bound method")
		return
	}

	if name == "" {
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
	this.serviceMap[name] = &service{Fn: f, ArgType: argType, ReplyType: replyType}
	this.serviceMapMu.Unlock()
	//注册类型池
	reflectTypePools.Init(argType)
	reflectTypePools.Init(replyType)
	return
}

// 处理请求
func (this *Server) handleRequest(ctx context.Context, req rpc.IMessage) (res rpc.IMessage, err error) {
	servicename := req.GetService()
	res = req.Clone()
	res.SetMessageType(rpc.Response)
	this.serviceMapMu.RLock()
	service, ok := this.serviceMap[servicename]
	this.serviceMapMu.RUnlock()
	if !ok {
		err = errors.New("rpcx: can't find service " + servicename)
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
