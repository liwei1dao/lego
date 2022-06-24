package core

import (
	"context"
	"reflect"
	"sync"

	"github.com/liwei1dao/lego/sys/log"
)

var TypeOfError = reflect.TypeOf((*error)(nil)).Elem()
var TypeOfContext = reflect.TypeOf((*context.Context)(nil)).Elem()

type (
	ISys interface {
		log.Ilogf
	}
	//路由
	IRoute interface {
	}
	ISelector interface {
		Select(ctx context.Context, route IRoute, serviceMethod string) string // SelectFunc
		UpdateServer(servers map[string]string)
	}
	IRPC interface {
		Go(ctx context.Context, sId, serviceMethod string) (call *Call, err error)
	}
	//服务对象
	Server struct {
		sync.Mutex
		Fn        reflect.Value
		ArgType   reflect.Type
		ReplyType reflect.Type
	}
	//异步返回结构
	Call struct {
		Args  interface{} //请求参数
		Reply interface{} //返回参数
		Error error       //错误信息
		Done  chan *Call
	}
)

type (
	//服务节点信息
	ServiceNode struct {
		Tag     string  `json:"Tag"`     //服务集群标签
		Type    string  `json:"Type"`    //服务类型
		Id      string  `json:"Id"`      //服务Id
		Version string  `json:"Version"` //服务版本
		Addr    string  `json:"Addr"`    //服务端地址
		Weight  float64 `json:"Weight"`  //服务负载权重
	}
)
