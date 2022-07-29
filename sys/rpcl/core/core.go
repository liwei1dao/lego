package core

import (
	"context"
	"time"

	"github.com/liwei1dao/lego/core"
	"github.com/liwei1dao/lego/sys/log"
	"github.com/liwei1dao/lego/sys/rpcl/protocol"
)

type ConnectType int //通信类型
const (
	Http  ConnectType = iota //http  连接对象
	Kafka                    //Kafka 连接
	Nats                     //Nats  连接
	Tcp                      //Tcp   连接
)

type SelectMode int //选择器类型
const (
	RandomSelect       SelectMode = iota //随机选择器
	RoundRobin                           //轮询选择器
	WeightedRoundRobin                   //权重轮询选择器
	WeightedICMP                         //网络质量选择器
	RuleRobin                            //规则选择器 默认
)

//系统对象
type ISys interface {
	log.Ilogf
	GetNodePath() string                                  //服务节点路径
	Handle(client IConnClient, message *protocol.Message) //接收到远程消息
}

//路由
type ICodec interface {
	Marshal(v interface{}) ([]byte, error)
	Unmarshal(data []byte, v interface{}) error
}

//选择器
type ISelector interface {
	Select(ctx context.Context, servicePath string) []*core.ServiceNode
	UpdateServer(servers []*core.ServiceNode)
}

//连接对象池
type IConnPool interface {
	Start() error
	GetClient(server *core.ServiceNode) (client IConnClient, err error)
	Close() error
}

type IConnClient interface {
	Write(msg []byte) (err error)
}

//连接配置信息
type Config struct {
	ConnectType       ConnectType   //通信类型
	Endpoints         []string      //节点信息
	ConnectionTimeout time.Duration //连接超时
	ReadTimeout       time.Duration //读取超时
	WriteTimeout      time.Duration //写入超时
	Username          string        //用户名
	Password          string        //密码
	Vsersion          string        //版本
}
