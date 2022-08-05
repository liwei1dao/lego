package rpc

import (
	"context"
	"fmt"
	"reflect"
	"runtime"
	"sync"

	"github.com/liwei1dao/lego/sys/log"
)

//异步返回结构
type MessageCall struct {
	ServicePath   string
	ServiceMethod string
	Metadata      map[string]string
	ResMetadata   map[string]string
	Args          interface{} //请求参数
	Reply         interface{} //返回参数
	Error         error       //错误信息
	Done          chan *MessageCall
}

func (call *MessageCall) done(log log.Ilogf) {
	select {
	case call.Done <- call:
		// ok
	default:
		log.Debugf("rpc: discarding Call reply due to insufficient Done chan capacity")

	}
}

//服务对象
type Server struct {
	sync.Mutex
	Fn        reflect.Value //执行方法
	ArgType   reflect.Type  //请求参数类型
	ReplyType reflect.Type  //返回数据类型
	IsActive  bool          //是否激活
}

func (this *Server) Call(ctx context.Context, argv, replyv reflect.Value) (err error) {
	defer func() {
		if r := recover(); r != nil {
			buf := make([]byte, 4096)
			n := runtime.Stack(buf, false)
			buf = buf[:n]

			err = fmt.Errorf("%v, argv: %+v, stack: %s",
				r, argv.Interface(), buf)
		}
	}()
	// Invoke the method, providing a new value for the reply.
	returnValues := this.Fn.Call([]reflect.Value{reflect.ValueOf(ctx), argv, replyv})
	// The return value for the method is an error.
	errInter := returnValues[0].Interface()
	if errInter != nil {
		return errInter.(error)
	}
	return nil
}
