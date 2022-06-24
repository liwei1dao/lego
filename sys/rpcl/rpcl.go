package rpcl

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"runtime"
	"strings"
	"sync"
	"unicode"
	"unicode/utf8"

	"github.com/liwei1dao/lego/sys/rpcl/core"
)

func newSys(options Options) (sys ISys, err error) {

	return
}

type RPCL struct {
	options      Options
	selector     core.ISelector
	serviceMapMu sync.RWMutex
	serviceMap   map[string]*core.Server
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
func (this *RPCL) UnRegister(name string) (err error) {
	return
}

//同步执行
func (this *RPCL) Call(ctx context.Context, route core.IRoute, serviceMethod string) (err error) { //同步调用 等待结果
	// sid := this.selector.Select(ctx, route, serviceMethod)
	return nil
}

//异步执行 异步返回
func (this *RPCL) Go(ctx context.Context, route core.IRoute, serviceMethod string) (call *core.Call, err error) { //异步调用 异步返回
	// sid := this.selector.Select(ctx, route, serviceMethod)
	return nil, nil
}

//同步执行 无返回
func (this *RPCL) GoNR(ctx context.Context, route core.IRoute, serviceMethod string) (err error) { //异步调用 无返回
	// sid := this.selector.Select(ctx, route, serviceMethod)
	return nil
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
