package cpool

import (
	"time"

	"github.com/liwei1dao/lego/core"
	lcore "github.com/liwei1dao/lego/sys/lrpc/core"
)

const (
	// ReaderBuffsize is used for bufio reader.
	ReaderBuffsize = 1024
	// WriterBuffsize is used for bufio writer.
	WriterBuffsize = 1024

	// WriteChanSize is used for response.
	WriteChanSize = 1024 * 1024
)

// 链接对象
type IConn interface {
}

// 通讯池
type ICPool interface {
	Start() error
	Heartbeat() []byte
	GetClient(node core.IServiceNode) (conn IConn, err error)
	AddClient(conn IConn, node core.IServiceNode) (err error)
	CloseClient(node core.IServiceNode) (err error)
	Handle(client IConn, message lcore.IMessage) //接收到远程消息
	Close() error
}

// 连接配置信息
type Config struct {
	Endpoints         []string      //节点信息
	ConnectionTimeout time.Duration //连接超时
	ReadTimeout       time.Duration //读取超时
	WriteTimeout      time.Duration //写入超时
	KeepAlivePeriod   time.Duration //保持活跃时期
	Username          string        //用户名
	Password          string        //密码
	Vsersion          string        //版本
}
