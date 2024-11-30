package tcppool

import "github.com/liwei1dao/lego/sys/lrpc/protocol"

type IPool interface {
	Options() *Options
	In() chan<- protocol.IMessage
	Out() <-chan protocol.IMessage
	CloseClient(path string) (err error)
}

var (
	defsys IPool
)
