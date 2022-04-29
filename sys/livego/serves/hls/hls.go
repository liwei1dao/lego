package hls

import (
	"net"
	"sync"

	"github.com/liwei1dao/lego/sys/livego/core"
)

type Server struct {
	server   core.IServer
	listener net.Listener
	conns    *sync.Map
}
