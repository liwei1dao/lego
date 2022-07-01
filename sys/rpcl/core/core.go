package core

import (
	"context"
	"reflect"
	"sync"
	"time"

	"github.com/liwei1dao/lego/sys/log"
)

var TypeOfError = reflect.TypeOf((*error)(nil)).Elem()
var TypeOfContext = reflect.TypeOf((*context.Context)(nil)).Elem()

// SelectMode defines the algorithm of selecting a services from candidates.
type SelectMode int

const (
	RandomSelect       SelectMode = iota //随机选择器
	RoundRobin                           //轮询选择器
	WeightedRoundRobin                   //权重轮询选择器
	WeightedICMP                         //网络质量选择器
)

type (
	ISys interface {
		log.Ilogf
		UpdateInterval() time.Duration //更新间隔
		GetServers() (servers map[string]bool)
	}
	//路由
	IRoute interface {
	}
	ISelector interface {
		Select(ctx context.Context, route IRoute, serviceMethod string) string // SelectFunc
		UpdateServer(servers map[string]string)
	}
	IRPC interface {
		Go(ctx context.Context, servicePath, serviceMethod string) (call *Call, err error)
	}
	ServiceDiscovery interface {
		GetServices() []*KVPair
		WatchService() chan []*KVPair
		RemoveWatcher(ch chan []*KVPair)
		Clone(servicePath string) (ServiceDiscovery, error)
		Close()
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
