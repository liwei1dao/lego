package lcore

import (
	"errors"
	"time"
)

var (
	ErrServerClosed          = errors.New("http: Server closed")                                   //服务关闭
	ErrMetaKVMissing         = errors.New("wrong metadata lines. some keys or values are missing") //解析Meta对象错误
	ErrUnsupportedCompressor = errors.New("unsupported compressor")                                //解压缩错误
	ErrXClientNoServer       = errors.New("can not found any server")
	ErrUnsupportedCodec      = errors.New("unsupported codec")
)

type contextKey struct {
	name string
}

func (k *contextKey) String() string { return "rpcx context value " + k.name }

var (
	RemoteConnContextKey = &contextKey{"remote-conn"} //远程服务连接对象
)

type ClientState int32

const (
	ClientClose      ClientState = iota //关闭状态
	ClientShakeHands                    //握手状态
	ClientRuning                        //运行中
	ClientCloseing                      //关闭中
)

type ConnectType int //通信类型
const (
	Tcp   ConnectType = iota //Tcp  连接对象
	Kafka                    //Kafka 连接
	Nats                     //Nats  连接
	http                     //http   连接
)

type SelectMode int //选择器类型
const (
	RandomSelect       SelectMode = iota //随机选择器
	RoundRobin                           //轮询选择器
	WeightedRoundRobin                   //权重轮询选择器
	WeightedICMP                         //网络质量选择器
	RuleRobin                            //规则选择器 默认
)

// 连接配置信息
type CPoolConfig struct {
	Endpoints         []string      //节点信息
	ConnectionTimeout time.Duration //连接超时
	ReadTimeout       time.Duration //读取超时
	WriteTimeout      time.Duration //写入超时
	KeepAlivePeriod   time.Duration //保持活跃时期
	Username          string        //用户名
	Password          string        //密码
	Vsersion          string        //版本
}
