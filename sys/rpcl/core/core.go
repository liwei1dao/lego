package core

import (
	"context"
	"fmt"
	"reflect"
	"runtime"
	"sync"
	"time"

	"github.com/liwei1dao/lego/sys/log"
	"github.com/liwei1dao/lego/utils/codec"
)

var TypeOfError = reflect.TypeOf((*error)(nil)).Elem()
var TypeOfContext = reflect.TypeOf((*context.Context)(nil)).Elem()

type NewConnectType int //通信类型
const (
	Http  NewConnectType = iota //http  连接对象
	Kafka                       //Kafka 连接
	Nats                        //Nats  连接
	Tcp                         //Tcp   连接
)

type SelectMode int //选择器类型
const (
	RandomSelect       SelectMode = iota //随机选择器
	RoundRobin                           //轮询选择器
	WeightedRoundRobin                   //权重轮询选择器
	WeightedICMP                         //网络质量选择器
	RuleRobin                            //规则选择器 默认
)

type (
	ISys interface {
		log.Ilogf
		Decoder() codec.IDecoder
		Encoder() codec.IEncoder
		GetUpdateInterval() time.Duration //更新间隔
		GetBasePath() string              //根节点路径
		GetNodePath() string              //节点路径
		GetServiceNode() *ServiceNode     //获取服务节点信息
		GetKafkaAddr() []string           //kafka 连接地址
		GetKafkaVersion() string          //kafka版本号
		Receive(message *Message)         //接收到远程消息
	}
	//路由
	IRoute interface {
	}
	//选择器
	ISelector interface {
		Select(ctx context.Context, route IRoute, serviceMethod string) string // SelectFunc
		UpdateServer(servers map[string]string)
	}
	//发现
	IDiscovery interface {
		GetServices() []*ServiceNode
		Close()
	}
	IConnect interface {
		Start() (err error)
		Close() (err error)
		Go(ctx context.Context, servicePath, serviceMethod string) (call *Call, err error)
	}

	//服务对象
	Server struct {
		sync.Mutex
		Fn        reflect.Value //执行方法
		ArgType   reflect.Type  //请求参数类型
		ReplyType reflect.Type  //返回数据类型
		IsActive  bool          //是否激活
	}

	//异步返回结构
	Call struct {
		Args  interface{} //请求参数
		Reply interface{} //返回参数
		Error error       //错误信息
		Done  chan *Call
	}
)

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
