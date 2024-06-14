package rpcx

import (
	"context"
	"fmt"
	"reflect"
	"runtime"
	"sync"

	"github.com/liwei1dao/lego/sys/log"
)

var typeOfError = reflect.TypeOf((*error)(nil)).Elem()
var typeOfContext = reflect.TypeOf((*context.Context)(nil)).Elem()

type seqKey struct{}
type methodType struct {
	sync.Mutex // protects counters
	method     reflect.Method
	ArgType    reflect.Type
	ReplyType  reflect.Type
	// numCalls   uint
}

type functionType struct {
	sync.Mutex // protects counters
	fn         reflect.Value
	ArgType    reflect.Type
	ReplyType  reflect.Type
}

type service struct {
	name     string                   // name of service
	rcvr     reflect.Value            // receiver of methods for the service
	typ      reflect.Type             // type of the receiver
	method   map[string]*methodType   // registered methods
	function map[string]*functionType // registered functions
}

func (this *service) call(ctx context.Context, mtype *methodType, argv, replyv reflect.Value) (err error) {
	defer func() {
		if r := recover(); r != nil {
			buf := make([]byte, 4096)
			n := runtime.Stack(buf, false)
			buf = buf[:n]

			err = fmt.Errorf("[service internal error]: %v, method: %s, argv: %+v, stack: %s",
				r, mtype.method.Name, argv.Interface(), buf)
			log.Errorf("err%v", err)
		}
	}()

	function := mtype.method.Func
	// Invoke the method, providing a new value for the reply.
	returnValues := function.Call([]reflect.Value{this.rcvr, reflect.ValueOf(ctx), argv, replyv})
	// The return value for the method is an error.
	errInter := returnValues[0].Interface()
	if errInter != nil {
		return errInter.(error)
	}
	return nil
}

// 执行注册方法
func (this *service) callForFunction(ctx context.Context, ft *functionType, argv, replyv reflect.Value) (err error) {
	defer func() {
		if r := recover(); r != nil {
			buf := make([]byte, 4096)
			n := runtime.Stack(buf, false)
			buf = buf[:n]

			// log.Errorf("failed to invoke service: %v, stacks: %s", r, string(debug.Stack()))
			err = fmt.Errorf("[service internal error]: %v, function: %s, argv: %+v, stack: %s",
				r, runtime.FuncForPC(ft.fn.Pointer()), argv.Interface(), buf)
			log.Errorf("err:%v", err)
		}
	}()

	// Invoke the function, providing a new value for the reply.
	returnValues := ft.fn.Call([]reflect.Value{reflect.ValueOf(ctx), argv, replyv})
	// The return value for the method is an error.
	errInter := returnValues[0].Interface()
	if errInter != nil {
		return errInter.(error)
	}

	return nil
}
