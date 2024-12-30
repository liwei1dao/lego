package connpool

import (
	"context"

	"github.com/liwei1dao/lego/core"
)

type ClientState int32

const (
	ClientClose      ClientState = iota //关闭状态
	ClientShakeHands                    //握手状态
	ClientRuning                        //运行中
	ClientCloseing                      //关闭中
)

type IConnClient interface {
	ServiceNode() core.IServiceNode
	SetServiceNode(node core.IServiceNode)
	State() ClientState
	Start()
	ResetHbeat()
	Write(msg []byte) (err error)
	Close() (err error)
}

// 连接对象池
type IConnPool interface {
	Start() error
	GetClient(node core.IServiceNode) (client IConnClient, err error)
	AddClient(client IConnClient, node core.IServiceNode) (err error)
	Close() error
}

// 主机对象
type IBodyHost interface {
	ShakehandsRequest(ctx context.Context, client IConnClient) (err error)
}
